package trerr

import (
	"fmt"
	"testing"
)

func TestErrTr(t *testing.T) {
	fmt.Println(TransError("copyright vote power not enough"))
}
