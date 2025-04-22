package id1_client

import (
	"net/http"
	"net/url"
)

func (t *id1ClientHttp) Mov(src, tgt Id1Key) error {
	url := url.URL{
		Scheme: t.url.Scheme,
		Path:   src.String(),
		Host:   t.url.Host,
	}
	req, _ := http.NewRequest(http.MethodPatch, url.String(), nil)
	req.Header.Add("X-Move-To", tgt.String())
	if len(t.token) > 0 {
		req.Header.Add("Authorization", t.token)
	}
	return t.do(req)
}
