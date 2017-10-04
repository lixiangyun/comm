package main

import (
	"log"
	"os"
)

func main() {

	args := os.Args

	if len(args) < 2 {
		log.Println("Usage: <-s/-c>")
		return
	}

	switch args[1] {
	case "-s":
		Server()
	case "-c":
		Client()
	}
}
