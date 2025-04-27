package id1_client

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/golang-jwt/jwt"
)

func (t *id1ClientHttp) Authenticate(id string, privateKeyPEM string) error {
	t.id = &id
	t.privateKey = &privateKeyPEM

	url := t.url
	url.Path = fmt.Sprintf("%s/auth", *t.id)
	req, _ := http.NewRequest(http.MethodGet, url.String(), nil)

	if res, err := t.doRes(req); err == nil {
		return ErrUnexpected
	} else if !errors.Is(err, ErrNotAuthenticated) {
		return err
	} else if resBodyBase64, err := io.ReadAll(res.Body); err != nil {
		return err
	} else if encryptedSecret, err := base64.StdEncoding.DecodeString(string(resBodyBase64)); err != nil {
		return err
	} else if secret, err := decrypt(encryptedSecret, privateKeyPEM); err != nil {
		return err
	} else {
		claims := jwt.StandardClaims{
			Subject: id,
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		jwtToken, _ := token.SignedString(secret)
		bearerToken := fmt.Sprintf("Bearer %s", jwtToken)
		t.token = &bearerToken
		return nil
	}
}

func decrypt(message []byte, privateKeyPEM string) ([]byte, error) {
	if block, _ := pem.Decode([]byte(privateKeyPEM)); block == nil {
		return []byte{}, fmt.Errorf("bad privateKeyPEM")
	} else if privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
		return []byte{}, err
	} else if decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, message); err != nil {
		return []byte{}, err
	} else {
		return decrypted, nil
	}
}
