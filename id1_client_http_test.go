package id1_client

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

var apiUrl = "http://localhost:8080"

func setupTestId(t *testing.T) string {
	id := fmt.Sprintf("id%d", time.Now().UnixMilli())
	pubKey := KK(id, "pub", "key")
	id1, _ := NewHttpClient(apiUrl)
	if err := id1.Set(pubKey, []byte(testPublicKey)); err != nil {
		t.Errorf("set key err %s", err)
	}
	t.Cleanup(func() {
		if err := id1.Authenticate(id, testPrivateKey); err != nil {
			t.Errorf("auth error %s", err)
		} else if err := id1.Del(K(id)); err != nil {
			t.Errorf("cleanup err %s", err)
		}
	})
	return id
}

func TestAuth(t *testing.T) {
	id := setupTestId(t)
	id1, _ := NewHttpClient(apiUrl)
	if err := id1.Authenticate(id, testPrivateKey); err != nil {
		t.Errorf("auth error %s", err)
	}
}

func TestCRUMD(t *testing.T) {
	id := setupTestId(t)
	testKey := KK(id, "test", "one")
	testTargetKey := KK(id, "test", "two")
	id1, _ := NewHttpClient(apiUrl)
	if err := id1.Authenticate(id, testPrivateKey); err != nil {
		t.Errorf("auth err %s", err)
	} else if err := id1.Set(testKey, []byte("test")); err != nil {
		t.Errorf("set err %s", err)
	} else if err := id1.Add(testKey, []byte("123")); err != nil {
		t.Errorf("add err %s", err)
	} else if val, err := id1.Get(testKey); err != nil {
		t.Errorf("get err %s", err)
	} else if string(val) != "test123" {
		t.Errorf("value not expected: %s", string(val))
	} else if err := id1.Mov(testKey, testTargetKey); err != nil {
		t.Errorf("mov err %s", err)
	} else if _, err := id1.Get(testKey); err == nil {
		t.Errorf("mov err src still exists")
	} else if val, err := id1.Get(testTargetKey); err != nil {
		t.Errorf("mov, get tgt err %s", err)
	} else if string(val) != "test123" {
		t.Errorf("mov, tgt value not expected: %s", string(val))
	} else if list, err := id1.List(KK(id), ListOptions{Recursive: true, Keys: true}); err != nil || len(list) != 2 {
		t.Errorf("list err")
	}

	if err := id1.Del(testTargetKey); err != nil {
		t.Errorf("del %s err %s", testTargetKey, err)
	} else if _, err := id1.Get(testTargetKey); err == nil {
		t.Errorf("del %s err, value still exists", testTargetKey)
	}
}

func TestExec(t *testing.T) {
	id := setupTestId(t)
	testKey := KK(id, "test", "one")
	id1, _ := NewHttpClient(apiUrl)
	if err := id1.Authenticate(id, testPrivateKey); err != nil {
		t.Errorf("auth err %s", err)
	}
	cmdSet := Command{
		Op:   Set,
		Key:  testKey,
		Data: []byte("test"),
	}
	cmdGet := Command{
		Op:  Get,
		Key: testKey,
	}
	cmdDel := Command{
		Op:  Del,
		Key: testKey,
	}
	if _, err := id1.Exec(cmdSet); err != nil {
		t.Errorf("set err %s", err)
	}
	if data, err := id1.Exec(cmdGet); err != nil || string(data) != "test" {
		t.Errorf("get err %s", err)
	}
	if _, err := id1.Exec(cmdDel); err != nil {
		t.Errorf("del err %s", err)
	}
}

func TestSend(t *testing.T) {
	id := setupTestId(t)
	testKey := KK(id, "test", "one")
	id1, _ := NewHttpClient(apiUrl)
	if err := id1.Authenticate(id, testPrivateKey); err != nil {
		t.Errorf("auth err %s", err)
	} else if _, err := id1.Connect(); err != nil {
		t.Errorf("err connecting %s", err)
	} else if err := id1.Send(Command{
		Op:   Set,
		Key:  testKey,
		Data: []byte("test"),
	}); err != nil {
		t.Errorf("send err %s", err)
	}
	time.Sleep(time.Millisecond * 10)
	if data, err := id1.Get(testKey); err != nil || string(data) != "test" {
		t.Errorf("send err...")
	}
}

func TestWebSocket(t *testing.T) {
	senderCount := 5
	eventCount := 5

	result := []Command{}
	cmdIn := make(chan Command, senderCount*eventCount)
	go func() {
		for {
			cmd := <-cmdIn
			result = append(result, cmd)
		}
	}()

	listenerId := setupTestId(t)
	listenerClient, _ := NewHttpClient(apiUrl)
	if err := listenerClient.Authenticate(listenerId, testPrivateKey); err != nil {
		t.Errorf("auth error %s", err)
	} else if _, err := listenerClient.Connect(); err != nil {
		t.Errorf("ws connect err %s", err)
	} else if err := listenerClient.Set(KK(listenerId, ".set"), []byte("*")); err != nil {
		t.Errorf("set err %s", err)
	} else {
		listenerClient.AddListener(func(cmd Command) {
			cmdIn <- cmd
		}, "1")
		listenerClient.Set(KK(listenerId, "test"), []byte("notification for self is not expected"))
	}
	time.Sleep(time.Millisecond * 10)
	for i := range senderCount {
		time.Sleep(time.Millisecond * 10)
		go func() {
			senderId := setupTestId(t)
			senderClient, _ := NewHttpClient(apiUrl)
			if err := senderClient.Authenticate(senderId, testPrivateKey); err != nil {
				t.Errorf("auth error %s", err)
			} else {
				for j := range eventCount {
					if err := senderClient.Set(KK(listenerId, fmt.Sprintf("test%d-%d", i, j)), fmt.Appendf(nil, "test%d-%d", i, j)); err != nil {
						t.Errorf("set err %s", err)
					}
				}
			}
			time.Sleep(time.Millisecond * 10)
		}()
	}

	time.Sleep(time.Second)

	if len(result) != senderCount*eventCount || result[0].Op != Set || len(result[0].Key.Segments) != 2 || !strings.HasPrefix(string(result[0].Data), "test") {
		t.Errorf("unexpected result %d", len(result))
	}
}

func TestWebSocketBreakPoint(t *testing.T) {
	accountCount := 0
	connectionCount := 0
	spacing := time.Millisecond * 10
	ids := []string{}
	for range accountCount {
		ids = append(ids, setupTestId(t))
		time.Sleep(spacing)
	}
	for _, id := range ids {
		for range connectionCount {
			go func() {
				id1, _ := NewHttpClient(apiUrl)
				if err := id1.Authenticate(id, testPrivateKey); err != nil {
					t.Errorf("auth error %s", err)
				} else if _, err := id1.Connect(); err != nil {
					t.Errorf("ws connect err %s", err)
				}
			}()
			time.Sleep(spacing)
		}
	}
}

var testPrivateKey = `-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQDhc8nghJEI1m6T1kfbBc9u0evVVgjI9Uw2wPVrvvg8/yGwlvUE
wUe1a2sA75M75MyFKNu14q7TsSwlmjKGZX76W+9gqP05a9MlDg1ev0gvg91l1L2q
FGrJ97Cfo7KZTZeKTvIJxbU3NqVgxi7ybdbgwxaDLWkvfarS5ftCdDzfEwIDAQAB
AoGAAUxy8b2hvjzJ90UJNexDJR7Fbn2n4Ex1a21SHJRFhshrb8J219erK5LaU2+N
2A751JyHI3WSTzrah77rIpYMFK9VArhwzRwgvnjyrg69EvLIyXp/9hy8/UQWKFm/
wHBl5nQSqbvP3iDgKqXqCbSNn6TbEpCusXElG7U+jPBVts0CQQDo0NHCLNBP7o5I
aInvoG5AXz1H3/ZGubihooRDxSo1CsyH/3s0laOdX36o+8N8gXaWLbmAhrxPLlQx
ZMKUnpQdAkEA9+dA8qq9P0FtRmNi8B327uduizYf1RUXY1StPUXN70q5DGDjA1J6
V6HBgh4UbugZfXY8G+zgeQxcjx1TiF947wJADvogBFXNsNav4IiZFwlDDnESCCWo
OjSIZB2IVLPCW1cugTE2Q9O8issx4r0Pflr1vgODA3mnc5CPaf4JZnYtIQJBAPPB
7wntuwIM2l8g8LL8M8d7xyWZhblm8MVaCLI8Bh9qEQTL68xjeCrcwcKowxy+mfnU
nYwz4hEEh6qtgmqQvf8CQDhXX7MX+aYHXHdjo3BZBkiVRlD3zSyRRtzuAigmPdoC
TDFpxmxIwBJUEpGJ0vAVhs6m9ouFh+1y0qVlNEcmlZY=
-----END RSA PRIVATE KEY-----
`

var testPublicKey = `-----BEGIN RSA PUBLIC KEY-----
MIGJAoGBAOFzyeCEkQjWbpPWR9sFz27R69VWCMj1TDbA9Wu++Dz/IbCW9QTBR7Vr
awDvkzvkzIUo27XirtOxLCWaMoZlfvpb72Co/Tlr0yUODV6/SC+D3WXUvaoUasn3
sJ+jsplNl4pO8gnFtTc2pWDGLvJt1uDDFoMtaS99qtLl+0J0PN8TAgMBAAE=
-----END RSA PUBLIC KEY-----
`
