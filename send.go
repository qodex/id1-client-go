package id1_client

import "fmt"

func (t *id1ClientHttp) Send(cmd Command) error {
	if t.conn == nil {
		return fmt.Errorf("not connected")
	}
	t.cmdOut <- cmd
	return nil
}
