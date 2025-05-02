package id1_client

import (
	"fmt"
	"slices"
	"sync"
)

type Id1ClientMock struct {
	listeners  map[string]func(Command)
	disconnect chan bool
	keys       map[string][]byte
	cmdIn      chan Command
	mu         *sync.Mutex
}

func NewId1ClientMock() Id1Client {

	cmdIn := make(chan Command, 8)
	listeners := map[string]func(Command){}
	mockKeys := map[string][]byte{}

	mockClient := Id1ClientMock{
		listeners: listeners,
		keys:      mockKeys,
		mu:        &sync.Mutex{},
		cmdIn:     cmdIn,
	}
	client := &mockClient

	go func() {
		for {
			cmd := <-client.cmdIn
			for _, v := range client.listeners {
				v(cmd)
			}
		}
	}()
	return client
}

func (t Id1ClientMock) Authenticate(id string, privateKey string) error {
	return nil
}

func (t Id1ClientMock) Connect() (chan bool, error) {
	t.disconnect = make(chan bool)
	return t.disconnect, nil
}

func (t Id1ClientMock) Close() {
	t.disconnect <- true
}

func (t Id1ClientMock) AddListener(listener func(cmd Command), listenerId string) string {
	t.mu.Lock()
	defer t.mu.Unlock()
	(t.listeners)[listenerId] = listener
	return listenerId
}

func (t Id1ClientMock) RemoveListener(listenerId string) {

}

func (t Id1ClientMock) Send(cmd Command) error {
	t.Exec(cmd)
	return nil
}

func (t Id1ClientMock) Exec(cmd Command) ([]byte, error) {
	t.cmdIn <- cmd
	switch cmd.Op {
	case Get:
		return t.Get(cmd.Key)
	case Del:
		return []byte{}, t.Del(cmd.Key)
	case Set:
		return []byte{}, t.Set(cmd.Key, cmd.Data)
	case Add:
		return []byte{}, t.Set(cmd.Key, cmd.Data)
	case Mov:
		return []byte{}, t.Mov(cmd.Key, K(string(cmd.Data)))
	}
	return []byte{}, nil
}

func (t Id1ClientMock) Get(key Id1Key) ([]byte, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	keys := t.keys

	if len(keys[key.String()]) == 0 {
		return keys[key.String()], fmt.Errorf("not found")
	} else {
		return keys[key.String()], nil
	}

}

func (t Id1ClientMock) Del(key Id1Key) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	keys := t.keys
	keys[key.String()] = []byte{}
	return nil
}
func (t Id1ClientMock) Set(key Id1Key, data []byte) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	keys := t.keys
	keys[key.String()] = data
	return nil
}

func (t Id1ClientMock) Add(key Id1Key, data []byte) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	keys := t.keys
	keys[key.String()] = slices.Concat(keys[key.String()], data)
	return nil
}

func (t Id1ClientMock) Mov(key, tgtKey Id1Key) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	keys := t.keys
	keys[tgtKey.String()] = keys[key.String()]
	keys[key.String()] = []byte{}
	return nil
}

func (t Id1ClientMock) List(key Id1Key, options ListOptions) (map[string][]byte, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	keys := t.keys
	return keys, nil

}
