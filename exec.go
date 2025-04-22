package id1_client

func (t *id1ClientHttp) Exec(cmd Command) ([]byte, error) {
	switch cmd.Op {
	case Get:
		return t.Get(cmd.Key)
	case Set:
		return []byte{}, t.Set(cmd.Key, cmd.Data)
	case Add:
		return []byte{}, t.Add(cmd.Key, cmd.Data)
	case Mov:
		return []byte{}, t.Mov(cmd.Key, K(string(cmd.Data)))
	case Del:
		return []byte{}, t.Del(cmd.Key)
	default:
		return []byte{}, ErrUnexpected
	}
}
