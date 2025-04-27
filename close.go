package id1_client

func (t id1ClientHttp) Close() {
	if t.conn != nil {
		t.conn.Close()
	}
}
