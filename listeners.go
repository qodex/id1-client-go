package id1_client

import "fmt"

func (t *id1ClientHttp) AddListener(listener func(cmd Command), listenerId string) string {
	t.mu.Lock()
	if len(listenerId) == 0 {
		listenerId = fmt.Sprintf("%d", len(t.listeners))
	}
	t.listeners[listenerId] = listener
	t.mu.Unlock()
	return listenerId
}

func (t *id1ClientHttp) RemoveListener(listenerId string) {
	t.mu.Lock()
	if len(listenerId) == 0 {
		listenerId = fmt.Sprintf("%d", len(t.listeners))
	}
	delete(t.listeners, listenerId)
	t.mu.Unlock()
}

func (t *id1ClientHttp) notifyListeners() {
	for {
		cmd := <-t.cmdIn
		t.mu.Lock()
		for _, listener := range t.listeners {
			go listener(cmd)
		}
		t.mu.Unlock()
	}
}
