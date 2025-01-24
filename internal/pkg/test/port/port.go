package port

import (
	"net"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func GetFreePort(t *testing.T) int {
	t.Helper()

	a, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		t.Fatal(cmp.Diff(err, nil))
	}

	l, err := net.ListenTCP("tcp", a)
	if err != nil {
		t.Fatal(cmp.Diff(err, nil))
	}
	defer l.Close()

	return l.Addr().(*net.TCPAddr).Port
}
