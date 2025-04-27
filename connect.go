package id1_client

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/websocket"
)

func (t id1ClientHttp) Connect() (chan bool, error) {
	if t.id == nil {
		return nil, fmt.Errorf("no id")
	}
	header := http.Header{}
	if t.token != nil {
		header.Add("Authorization", *t.token)
	}

	if t.conn != nil {
		t.conn.Close()
		t.conn = nil
	}
	url := url.URL{
		Scheme: strings.ReplaceAll(strings.ToLower(t.url.Scheme), "http", "ws"),
		Path:   fmt.Sprintf("%s/ws", *t.id),
		Host:   t.url.Host,
	}
	if conn, _, err := websocket.DefaultDialer.Dial(url.String(), header); err != nil {
		log.Printf("error connecting: %s", err)
		return nil, err
	} else {
		t.conn = conn
		disconnectSignal := make(chan bool)
		go t.readWebsocket(disconnectSignal)
		go t.writeWebsocket(disconnectSignal)
		go t.notifyListeners()
		return disconnectSignal, nil
	}
}
