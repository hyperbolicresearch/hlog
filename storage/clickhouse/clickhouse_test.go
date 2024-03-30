package clickhouseservice

import "testing"

func TestConn(t *testing.T) {
	addrs := []string{"localhost:9000"}
	_, err := Conn(addrs)
	if err != nil {
		t.Error(err)
	}
}
