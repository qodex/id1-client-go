package id1_client

import "fmt"

func (t id1ClientHttp) AddListener(listener func(cmd Command), listenerId string) string {
	t.mu.Lock()
	listeners := *t.listeners
	if len(listenerId) == 0 {
		listenerId = fmt.Sprintf("%d", len(listeners))
	}
	listeners[listenerId] = listener
	t.mu.Unlock()
	return listenerId
}

func (t id1ClientHttp) RemoveListener(listenerId string) {
	t.mu.Lock()
	listeners := *t.listeners
	if len(listenerId) == 0 {
		listenerId = fmt.Sprintf("%d", len(listeners))
	}
	delete(listeners, listenerId)
	t.mu.Unlock()
}

func (t id1ClientHttp) notifyListeners() {
	for {
		cmd := <-t.cmdIn
		t.mu.Lock()
		listeners := *t.listeners
		for _, listener := range listeners {
			go listener(cmd)
		}
		t.mu.Unlock()
	}
}
