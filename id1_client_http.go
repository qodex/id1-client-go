package id1_client

import (
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type id1ClientHttp struct {
	url        *url.URL
	id         *string
	privateKey *string
	conn       *websocket.Conn
	token      *string
	cmdIn      chan Command
	cmdOut     chan Command
	listeners  *map[string]func(Command)
	mu         *sync.Mutex
}

func NewHttpClient(apiEndpoint string) (Id1Client, error) {
	client := id1ClientHttp{
		cmdIn:     make(chan Command, 256),
		cmdOut:    make(chan Command, 256),
		listeners: &map[string]func(Command){},
		mu:        &sync.Mutex{},
	}
	if url, err := url.Parse(apiEndpoint); err != nil {
		return &client, err
	} else {
		client.url = url
		return &client, err
	}
}

func (t id1ClientHttp) do(req *http.Request) error {
	_, err := t.doRes(req)
	return err
}

func (t id1ClientHttp) doRes(req *http.Request) (*http.Response, error) {
	if t.token != nil && len(*t.token) > 0 {
		req.Header.Add("Authorization", *t.token)
	}
	select {
	default:
		if res, err := http.DefaultClient.Do(req); err != nil {
			return res, err
		} else {
			return res, httpStatusErr(res.StatusCode)
		}
	case <-time.After(time.Second):
		return nil, ErrTimeout
	}
}

func (t id1ClientHttp) Send(cmd Command) error {
	t.cmdOut <- cmd
	return nil
}

func (t id1ClientHttp) Set(key Id1Key, data []byte) error {
	_, err := t.Exec(Command{Op: Set, Key: key, Data: data})
	return err
}

func (t id1ClientHttp) Add(key Id1Key, data []byte) error {
	_, err := t.Exec(Command{Op: Add, Key: key, Data: data})
	return err
}

func (t id1ClientHttp) Get(key Id1Key) ([]byte, error) {
	data, err := t.Exec(Command{Op: Get, Key: key})
	return data, err
}

func (t id1ClientHttp) Del(key Id1Key) error {
	_, err := t.Exec(Command{Op: Del, Key: key})
	return err
}

func (t id1ClientHttp) Mov(src, tgt Id1Key) error {
	_, err := t.Exec(Command{Op: Mov, Key: src, Data: []byte(tgt.String())})
	return err
}

func (t id1ClientHttp) List(key Id1Key, options ListOptions) (map[string][]byte, error) {
	if data, err := t.Exec(Command{Op: Get, Key: K(key.String() + "*"), Args: options.Map()}); err != nil {
		return map[string][]byte{}, err
	} else {
		return decodeList(data), nil
	}
}

func (t id1ClientHttp) Close() {
	if t.conn != nil {
		t.conn.Close()
	}
}
