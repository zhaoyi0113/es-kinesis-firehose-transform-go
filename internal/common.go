package internal

import "log"

func FailOnError(err error, message string) {
	if err != nil {
		log.Println(message, err)
		panic(err)
	}
}
