package main

import (
	"log"
)

func initLogger() {
	log.SetFlags(0)
	log.Println("Server is running")
}

func pnc(err error) {
	if err != nil {
		panic(err)
	}
}
