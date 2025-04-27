package id1_client

import (
	"net/http"
	"net/url"
)

func (t id1ClientHttp) Del(key Id1Key) error {
	url := url.URL{
		Scheme: t.url.Scheme,
		Path:   key.String(),
		Host:   t.url.Host,
	}
	req, err := http.NewRequest(http.MethodDelete, url.String(), nil)
	if err != nil {
		return err
	}
	return t.do(req)
}
