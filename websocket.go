package id1_client

import (
	"log"

	"github.com/gorilla/websocket"
)

func (t *id1ClientHttp) readWebsocket(disconnectSignal chan bool) {
	for {
		if _, message, err := t.conn.ReadMessage(); err != nil {
			disconnectSignal <- true
			break
		} else if cmd, err := ParseCommand(message); err != nil {
			log.Printf("unknown command: %s", string(message))
		} else {
			t.cmdIn <- cmd
		}
	}
}

func (t *id1ClientHttp) writeWebsocket(disconnectSignal chan bool) {
	for {
		cmd := <-t.cmdOut
		if err := t.conn.WriteMessage(websocket.BinaryMessage, cmd.Bytes()); err != nil {
			t.cmdIn <- cmd
			disconnectSignal <- true
			break
		}
	}
}
