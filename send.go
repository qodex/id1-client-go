package id1_client

func (t id1ClientHttp) Send(cmd Command) error {
	t.cmdOut <- cmd
	return nil
}
