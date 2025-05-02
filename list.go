package id1_client

import (
	"bytes"
	"encoding/base64"
	"log"
)

func decodeList(data []byte) map[string][]byte {
	list := map[string][]byte{}
	lines := bytes.Split(data, []byte("\n"))
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
	return list
}
