package id1_client

import (
	"fmt"
	"net/url"
	"slices"
	"strings"
)

func (t Command) Bytes() []byte {
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
	bytes := slices.Concat([]byte(command), []byte("\n"), t.Data)
	return bytes
}

func (t Command) String() string {
	return string(t.Bytes())
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
