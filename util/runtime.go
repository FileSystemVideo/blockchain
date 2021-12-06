package util

import (
	"runtime"
	"strings"
)


func GetFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	arry := strings.Split(f.Name(), ".")
	return arry[len(arry)-1]
}
