package id1_client

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
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

// export env_ids="one\ntwo"
// seg0/seg1/seg2/seg3 > $env_ids/arc/$1/$2/$timestamp = [one/arc/seg1/seg2/1746089763786, two/arc/seg1/seg2/1746089763786]
func (t Id1Key) Map(keymap string) []Id1Key {
	keys := []Id1Key{}

	if len(keymap) == 0 {
		keys = append(keys, t)
		return keys
	}

	kmask := K(keymap)
	segNumRE, _ := regexp.Compile(`^\$\d+$`)
	ts := fmt.Sprintf("%d", time.Now().UnixMicro())

	// eval $
	for i, seg := range kmask.Segments {
		segVal := "map_err"
		if !strings.HasPrefix(seg, "$") {
			segVal = seg
		} else if seg == "$timestamp" {
			segVal = ts
		} else if segNumRE.MatchString(seg) {
			if segNum, err := strconv.Atoi(seg[1:]); err == nil && len(t.Segments) > segNum {
				segVal = t.Segments[segNum]
			}
		} else if len(os.Getenv(seg[1:])) > 0 {
			segVal = os.Getenv(seg[1:])
		}
		kmask.Segments[i] = segVal
	}

	// eval multiline segs
	fanned := fanout(kmask.Segments)
	for _, segs := range fanned {
		keys = append(keys, K(strings.Join(segs, "/")))
	}

	return keys
}

func K(s string) Id1Key {
	k := Id1Key{}
	if len(s) == 0 {
		return k
	}
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

// if input is []string{"one\ntwo", "three\nfour"},
// then output should be [][]string{
// []string{"one", "three"},
// []string{"one", "four"},
// []string{"two", "three"},
// []string{"two", "four"},
// }
func fanout(input []string) [][]string {
	splitInputs := make([][]string, len(input))
	for i, str := range input {
		splitInputs[i] = strings.Split(str, "\n")
	}

	var build func(int, []string)
	var result [][]string

	build = func(index int, current []string) {
		if index == len(splitInputs) {
			combination := make([]string, len(current))
			copy(combination, current)
			result = append(result, combination)
			return
		}
		for _, val := range splitInputs[index] {
			build(index+1, append(current, val))
		}
	}

	build(0, []string{})
	return result
}
