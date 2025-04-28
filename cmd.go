package id1_client

import (
	"fmt"
	"net/url"
	"slices"
	"strings"
)

type Command struct {
	Op   Op
	Key  Id1Key
	Args map[string]string
	Data []byte
}

type Op int

const (
	Unknown Op = iota
	Set
	Add
	Get
	Del
	Mov
	List
)

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
	if op, ok := nameOp[s]; ok {
		return op
	} else {
		return Unknown
	}
}

func (t Command) Bytes() []byte {
	bytes := slices.Concat([]byte(t.String()), []byte("\n"), t.Data)
	return bytes
}

func (t Command) String() string {
	args := url.Values{}
	for arg := range t.Args {
		args.Set(arg, t.Args[arg])
	}
	url := url.URL{
		Scheme:   t.Op.String(),
		Path:     t.Key.String(),
		RawQuery: args.Encode(),
	}
	command := strings.ReplaceAll(url.String(), "//", "/")
	return command
}

func (t Command) IsEmpty() bool {
	return len(t.Key.String()) == 0
}

func (t Command) IsNotEmpty() bool {
	return !t.IsEmpty()
}

func ParseCommand(data []byte) (Command, error) {
	command := Command{}
	firstLineEnd := slices.Index(data, byte('\n'))
	if firstLineEnd < 0 {
		firstLineEnd = len(data)
		data = append(data, byte('\n'))
	}
	firstLine := string(data[0:firstLineEnd])
	command.Data = data[firstLineEnd+1:]

	if strings.HasPrefix(firstLine, "#") {
		return command, fmt.Errorf("not a command")
	}

	url, err := url.Parse(firstLine)
	if err != nil {
		return command, err
	}
	command.Op = op(url.Scheme)
	command.Key = K(url.Path)
	command.Args = map[string]string{}
	for k := range url.Query() {
		command.Args[k] = url.Query().Get(k)
	}
	return command, nil
}
