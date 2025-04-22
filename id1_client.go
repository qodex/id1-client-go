package id1_client

import "fmt"

type Id1Client interface {
	Authenticate(id string, privateKey string) error
	Connect() error
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

type Op int

const (
	Set Op = iota
	Add
	Get
	Del
	Mov
	List
)

type Command struct {
	Op   Op
	Key  Id1Key
	Args map[string]string
	Data []byte
}

type ListOptions struct {
	Limit          int
	SizeLimit      int
	TotalSizeLimit int
	Keys           bool
	Recursive      bool
	Children       bool
}

func (t ListOptions) Map() map[string]string {
	args := map[string]string{}
	if t.Limit > 0 {
		args["limit"] = fmt.Sprintf("%d", t.Limit)
	}
	if t.SizeLimit > 0 {
		args["size-limit"] = fmt.Sprintf("%d", t.SizeLimit)
	}
	if t.TotalSizeLimit > 0 {
		args["total-size-limit"] = fmt.Sprintf("%d", t.TotalSizeLimit)
	}
	args["keys"] = fmt.Sprintf("%t", t.Keys)
	args["recursive"] = fmt.Sprintf("%t", t.Recursive)
	args["children"] = fmt.Sprintf("%t", t.Children)
	return args
}

var opName = map[Op]string{
	Set:  "set",
	Add:  "add",
	Get:  "get",
	Del:  "del",
	Mov:  "mov",
	List: "list",
}

var nameOp = map[string]Op{
	"set":  Set,
	"add":  Add,
	"get":  Get,
	"del":  Del,
	"mov":  Mov,
	"list": List,
}

func (t Op) String() string {
	return opName[t]
}

func op(s string) Op {
	return nameOp[s]
}
