package main

import (
	"comm"
	"log"
	"os"
	"time"
)

const (
	IP   = "localhost"
	PORT = "6565"
)

func serverhandler(s *comm.Server, reqid uint32, body []byte) {
	log.Println(string(body))
	body = []byte("hello world from server!")
	err := s.SendMsg(reqid, body)
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func Server() {
	listen := comm.NewListen(":" + PORT)
	for {
		server, err := listen.Accept()
		if err != nil {
			log.Println(err.Error())
			return
		}
		server.RegHandler(0, serverhandler)
		server.Start(1)
	}
}

func clienthandler(c *comm.Client, reqid uint32, body []byte) {
	log.Println(string(body))
}

func Client() {

	client := comm.NewClient(IP + ":" + PORT)
	client.RegHandler(0, clienthandler)
	client.Start(1)
	defer client.Stop()

	sendbuf := []byte("hello world from client!")
	err := client.SendMsg(0, sendbuf)
	if err != nil {
		log.Println(err.Error())
		return
	}

	time.Sleep(time.Second)
}

func main() {
	args := os.Args
	if len(args) < 2 {
		log.Println("Usage: <-s/-c>")
	}
	switch args[1] {
	case "-s":
		Server()
	case "-c":
		Client()
	}
}
