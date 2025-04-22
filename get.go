package id1_client

import (
	"io"
	"net/http"
	"net/url"
)

func (t *id1ClientHttp) Get(key Id1Key) ([]byte, error) {
	url := url.URL{
		Scheme: t.url.Scheme,
		Path:   key.String(),
		Host:   t.url.Host,
	}
	req, _ := http.NewRequest(http.MethodGet, url.String(), nil)
	if len(t.token) > 0 {
		req.Header.Add("Authorization", t.token)
	}
	if res, err := http.DefaultClient.Do(req); err != nil {
		return []byte{}, err
	} else if err := httpStatusErr(res.StatusCode); err != nil {
		return []byte{}, err
	} else if body, err := io.ReadAll(res.Body); err != nil {
		return []byte{}, err
	} else {
		return body, nil
	}
}
