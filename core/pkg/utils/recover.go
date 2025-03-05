package utils

import "log"

func Recover() {
	if r := recover(); r != nil {
		log.Fatalf("recovered from panic: %s", r)
	}
}
