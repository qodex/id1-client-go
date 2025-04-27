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
	//log.Println("doRes---", (*req).URL.Path)
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
