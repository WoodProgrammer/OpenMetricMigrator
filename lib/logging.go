package lib

import (
	"log"
	"runtime"
)

func LogErrorWithLine(err error, msg string) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		log.Fatalf("%s: %v (File: %s, Line: %d)", msg, err, file, line)
	} else {
		log.Fatalf("%s: %v", msg, err)
	}
}
