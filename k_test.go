package id1_client

import (
	"os"
	"testing"
)

func TestMap(t *testing.T) {
	os.Setenv("env1", "ein\nzwei")
	os.Setenv("env2", "uno\ndos")
	keys := K("test1/cmd/test2/123").Map("$env2/$env1/arc/$3")
	if len(keys) != 4 || keys[0].String() != "uno/ein/arc/123" {
		t.Errorf("err map %s", keys)
	}
	noMap := K("test1/cmd/test2/123").Map("")
	if noMap[0].String() != "test1/cmd/test2/123" {
		t.Errorf("err no map res %s", noMap)
	}
}
