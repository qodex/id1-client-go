package id1_client

import (
	"fmt"
	"strings"
)

type Id1Key struct {
	Id       string
	Name     string
	Parent   string
	Pub      bool
	Segments []string
}

func (t Id1Key) String() string {
	str := strings.Join(t.Segments, "/")
	if len(t.Segments) == 1 {
		str = str + "/"
	}
	str = strings.Trim(str, "\n")
	str = strings.Trim(str, " ")
	return str
}

func K(s string) Id1Key {
	k := Id1Key{}
	if len(s) == 0 {
		return k
	}
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, " ", "")
	s = strings.Trim(s, "/")

	k.Segments = strings.Split(s, "/")
	k.Id = k.Segments[0]
	k.Name = k.Segments[len(k.Segments)-1]

	if len(k.Segments) > 1 {
		k.Parent = strings.Join(k.Segments[:len(k.Segments)-1], "/")
	}
	if len(k.Segments) > 1 {
		k.Pub = k.Segments[1] == "pub"
	}

	return k
}

func KK(segments ...any) Id1Key {
	strSegments := []string{}
	for _, seg := range segments {
		if s, ok := seg.(string); ok {
			strSegments = append(strSegments, s)
		}
		if i, ok := seg.(int); ok {
			strSegments = append(strSegments, fmt.Sprintf("%d", i))
		}
		if stringer, ok := seg.(fmt.Stringer); ok {
			strSegments = append(strSegments, stringer.String())
		}
	}
	return K(strings.Join(strSegments, "/"))
}
