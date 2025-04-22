package id1_client

import (
	"bytes"
	"net/http"
	"net/url"
)

func (t *id1ClientHttp) Add(key Id1Key, data []byte) error {
	url := url.URL{
		Scheme: t.url.Scheme,
		Path:   key.String(),
		Host:   t.url.Host,
	}
	req, _ := http.NewRequest(http.MethodPatch, url.String(), bytes.NewReader(data))
	if len(t.token) > 0 {
		req.Header.Add("Authorization", t.token)
	}
	return t.do(req)
}
