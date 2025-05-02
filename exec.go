package id1_client

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
)

var opHttpMethod = map[Op]string{
	Get: http.MethodGet,
	Set: http.MethodPost,
	Add: http.MethodPatch,
	Mov: http.MethodPatch,
	Del: http.MethodDelete,
}

func (t id1ClientHttp) Exec(cmd Command) ([]byte, error) {
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
	req, err := http.NewRequest(opHttpMethod[cmd.Op], url.String(), bytes.NewReader(cmd.Data))
	if err != nil {
		return []byte{}, err
	}

	if cmd.Op == Mov {
		req.Header.Add("X-Move-To", string(cmd.Data))
	}

	data := []byte{}
	if res, err := t.doRes(req); err != nil {
		return data, err
	} else if err := httpStatusErr(res.StatusCode); err != nil {
		return data, err
	} else if body, err := io.ReadAll(res.Body); err != nil {
		return data, err
	} else {
		data = body
	}

	return data, nil

}
