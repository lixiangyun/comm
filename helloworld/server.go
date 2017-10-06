package main

import (
	"comm"
	"log"
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
		log.Println("new server instance.")

		server.RegHandler(0, serverhandler)

		go func() {
			server.Start(1, 10)
			server.Wait()
			log.Println("free server instance.")
		}()
	}
}
