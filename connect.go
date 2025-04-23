package id1_client

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"net/url"
)

func (t *id1ClientHttp) Connect() (chan bool, error) {
	header := http.Header{
		"Authorization": []string{t.token},
	}

	if t.conn != nil {
		t.conn.Close()
		t.conn = nil
	}
	url := url.URL{
		Scheme: "ws",
		Path:   fmt.Sprintf("%s/ws", t.id),
		Host:   t.url.Host,
	}
	disconnectSignal := make(chan bool)
	if conn, _, err := websocket.DefaultDialer.Dial(url.String(), header); err != nil {
		return nil, err
	} else {
		t.conn = conn
		go t.readWebsocket(disconnectSignal)
		go t.writeWebsocket(disconnectSignal)
		go t.notifyListeners()
		return disconnectSignal, nil
	}
}
