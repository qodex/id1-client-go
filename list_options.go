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

func MapListOptions(args map[string]string) ListOptions {
	opt := ListOptions{}
	if i, err := strconv.ParseInt(args["limit"], 10, 64); err == nil {
		opt.Limit = int(i)
	} else {
		opt.Limit = 1000
	}
	if i, err := strconv.ParseInt(args["size-limit"], 10, 64); err == nil {
		opt.SizeLimit = int(i)
	} else {
		opt.SizeLimit = 100 * MB
	}
	if i, err := strconv.ParseInt(args["total-size-limit"], 10, 64); err == nil {
		opt.TotalSizeLimit = int(i)
	} else {
		opt.TotalSizeLimit = 100 * MB
	}
	opt.Keys = args["keys"] == "true"
	opt.Recursive = args["recursive"] == "true"
	opt.Children = args["children"] == "true"
	return opt
}
