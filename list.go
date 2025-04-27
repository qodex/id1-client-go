package id1_client

import (
	"bytes"
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"net/url"
)

func (t id1ClientHttp) List(key Id1Key, options ListOptions) (map[string][]byte, error) {
	list := map[string][]byte{}
	params := url.Values{}
	for k := range options.Map() {
		params.Set(k, options.Map()[k])
	}
	url := url.URL{
		Scheme:   t.url.Scheme,
		Path:     key.String() + "*",
		Host:     t.url.Host,
		RawQuery: params.Encode(),
	}
	req, _ := http.NewRequest(http.MethodGet, url.String(), nil)
	if res, err := t.doRes(req); err != nil {
		return list, err
	} else if err := httpStatusErr(res.StatusCode); err != nil {
		return list, err
	} else if body, err := io.ReadAll(res.Body); err != nil {
		return list, err
	} else {
		lines := bytes.Split(body, []byte("\n"))
		for _, line := range lines {
			k := string(bytes.Split(line, []byte("="))[0])
			v := []byte{}
			if len(line) > len(k) {
				v = line[len(k)+1:]
			}
			if data, err := base64.StdEncoding.DecodeString(string(v)); err != nil {
				log.Printf("list item decode err %s", err)
			} else {
				list[k] = data
			}
		}
		return list, nil
	}
}
