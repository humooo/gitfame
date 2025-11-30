package internal

import (
	"log"
)

func ProcessError(err error, errorString string) {
	if err != nil {
		log.Fatalf("[ERROR] %s:\n%v", errorString, err)
	}
}
