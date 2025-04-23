package id1_client

type Id1Client interface {
	Authenticate(id string, privateKey string) error
	Connect() (chan bool, error)
	Close()
	AddListener(listener func(cmd Command), listenerId string) string
	RemoveListener(listenerId string)
	Send(cmd Command) error
	Exec(cmd Command) ([]byte, error)
	Get(key Id1Key) ([]byte, error)
	Del(key Id1Key) error
	Set(key Id1Key, data []byte) error
	Add(key Id1Key, data []byte) error
	Mov(key, tgtKey Id1Key) error
	List(key Id1Key, options ListOptions) (map[string][]byte, error)
}
