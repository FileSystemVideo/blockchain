package log

import "testing"

func TestLog(t *testing.T) {
	addr := "localhost"
	Critical("Admin server Running on %s", addr)
}