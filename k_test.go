package id1_client

import (
	"os"
	"testing"
)

func TestStar(t *testing.T) {
	keys := K("cmd/test2/123").Map("test1/*")
	if keys[0].String() != "test1/cmd/test2/123" {
		t.Errorf("err star %s", keys)
	}
}

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

func TestParent(t *testing.T) {
	k := K("test1/cmd/test2/123")
	if k.Parent().Name != "test2" {
		t.Errorf("unexpected Parent %s", k.Parent())
	}
}
