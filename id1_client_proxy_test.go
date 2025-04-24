package id1_client

import (
	"fmt"
	"testing"
)

func TestClientProxy(t *testing.T) {
	proxy := NewId1ClientProxy(NewId1ClientMock())
	proxy.AddPreprocessor(func(cmd *Command) error {
		if cmd.Op == Set {
			cmd.Key = KK(cmd.Key.String(), "good")
		}
		return nil
	})
	proxy.AddPreprocessor(func(cmd *Command) error {
		if cmd.Op == Del && cmd.Key.Name == "key" {
			return fmt.Errorf("don't delete keys")
		}
		return nil
	})
	proxy.AddPostprocessor(func(data []byte, err error) ([]byte, error) {
		if string(data) == "now" {
			data = []byte("please")
		}
		return data, err
	})

	if err := proxy.Del(K("test/pub/key")); err == nil {
		t.Errorf("expected error")
	}

	proxy.Set(K("test/me"), []byte("now"))
	if data, _ := proxy.Get(K("test/me")); len(data) > 0 {
		t.Errorf("expected none")
	}
	if data, _ := proxy.Get(K("test/me/good")); string(data) != "please" {
		t.Errorf("unexpected value %s", string(data))
	}

	proxy.AddPostprocessor(func(data []byte, err error) ([]byte, error) {
		if string(data) == "please" {
			data = []byte("bitte")
		}
		return data, err
	})
	if list, _ := proxy.List(K("test"), ListOptions{}); string(list["test/me/good"]) != "bitte" {
		t.Errorf("unexpected value %s", string(list["test/me/good"]))
	}
}
