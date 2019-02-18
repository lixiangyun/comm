package main

import (
	"log"

	"github.com/lixiangyun/comm"
)

const (
	IP          = "127.0.0.1"
	PORT        = "6565"
	CA          = "./crt/ca.crt"
	SERVER_CERT = "./crt/server.crt"
	SERVER_KEY  = "./crt/server.pem"
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

	listen.TlsEnable(CA, SERVER_CERT, SERVER_KEY)

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
