package main

import (
	"log"
	"os"
)

// 进程启动入口
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
