package util

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/charmbracelet/log"
)

func Logger() {
	logger := log.NewWithOptions(os.Stdout, log.Options{
		ReportTimestamp: false,
		ReportCaller:    false,
	})

	log.SetDefault(logger)
}

func GenerateStacktrace(temp interface{}) string {
	stack := make([]uintptr, 64)
	currStack := runtime.Callers(2, stack)
	// stack is inverse
	stack = invertStack(stack)

	// strs := strings.Split((runtime.FuncForPC(reflect.ValueOf(temp).Pointer()).Name()), ".")

	var ret string
	for _,stackPC := range stack[:currStack] {
		fun := runtime.FuncForPC(stackPC)

		name := strings.Split(fun.Name(), "/");
		ret = fmt.Sprintf("%v > %v", ret, name[len(name) - 1])
	}

	return ret
}

func invertStack(arr []uintptr) []uintptr {
	out := make([]uintptr, len(arr))

	for i := 0; i < len(arr); i++ {
		if arr[i] != 0 {
			out[len(arr) - i - 1] = arr[i]
		}
	}

	var cutOff = 0;
	for i,e := range out {
		if e != 0 {
			cutOff = i
			break
		}
	}

	out = out[cutOff:]

	log.Infof("In %v\nOut %v", arr, out)

	return out
}

