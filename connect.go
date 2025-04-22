package id1_client

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"net/url"
)

func (t *id1ClientHttp) Connect() error {
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
	if conn, _, err := websocket.DefaultDialer.Dial(url.String(), header); err != nil {
		return err
	} else {
		t.conn = conn
		go t.readWebsocket()
		go t.writeWebsocket()
		go t.notifyListeners()
		return nil
	}
}
