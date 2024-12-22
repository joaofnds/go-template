package util

import (
	"reflect"
	"runtime"
	"strings"
)

func FunctionName(function any) string {
	funcInfo := runtime.FuncForPC(reflect.ValueOf(function).Pointer())
	if funcInfo == nil {
		return "unknown function"
	}

	// Example: "app/user/user_adapter.(*PromProbe).Register-fm"
	fullName := funcInfo.Name()

	fullName = strings.TrimSuffix(fullName, "-fm")

	if strings.Contains(fullName, "(*") {
		fullName = strings.Replace(fullName, "(*", "", 1)
		fullName = strings.Replace(fullName, ")", "", 1)
	}

	return fullName
}
