package main

import (
	"log"
	"os"
)

// 进程启动入口
func main() {
	args := os.Args
	if len(args) < 3 {
		log.Println("Usage: <-s/-c> <ip:port>")
		return
	}

	switch args[1] {
	case "-s":
		Server(args[2])
	case "-c":
		Client(args[2])
	}
}
