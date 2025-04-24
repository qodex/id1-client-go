package id1_client

import (
	"fmt"
	"strconv"
)

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

func (t *ListOptions) Parse(args map[string]string) {
	if i, err := strconv.ParseInt(args["limit"], 10, 64); err == nil {
		t.Limit = int(i)
	} else {
		t.Limit = 1000
	}
	if i, err := strconv.ParseInt(args["size-limit"], 10, 64); err == nil {
		t.SizeLimit = int(i)
	} else {
		t.SizeLimit = 100 * MB
	}
	if i, err := strconv.ParseInt(args["total-size-limit"], 10, 64); err == nil {
		t.TotalSizeLimit = int(i)
	} else {
		t.TotalSizeLimit = 100 * MB
	}
	t.Keys = args["keys"] == "true"
	t.Recursive = args["recursive"] == "true"
	t.Children = args["children"] == "true"
}
