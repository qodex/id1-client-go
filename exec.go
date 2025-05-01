package id1_client

import (
	"bytes"
	"net/http"
	"net/url"
)

func (t id1ClientHttp) Exec(cmd Command) ([]byte, error) {
	switch cmd.Op {
	case Get:
		return t.Get(cmd.Key)
	case Set:
		args := url.Values{}
		for arg := range cmd.Args {
			args.Set(arg, cmd.Args[arg])
		}
		url := url.URL{
			Scheme:   t.url.Scheme,
			Path:     cmd.Key.String(),
			RawQuery: args.Encode(),
			Host:     t.url.Host,
		}
		req, _ := http.NewRequest(http.MethodPost, url.String(), bytes.NewReader(cmd.Data))
		return []byte{}, t.do(req)
	case Add:
		return []byte{}, t.Add(cmd.Key, cmd.Data)
	case Mov:
		return []byte{}, t.Mov(cmd.Key, K(string(cmd.Data)))
	case Del:
		return []byte{}, t.Del(cmd.Key)
	default:
		return []byte{}, ErrUnexpected
	}
}
