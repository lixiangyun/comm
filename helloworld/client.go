package main

import (
	"log"
	"time"

	"github.com/lixiangyun/comm"
)

const (
	CLIENT_CERT = "./crt/client.crt"
	CLIENT_KEY  = "./crt/client.pem"
)

func clienthandler(c *comm.Client, reqid uint32, body []byte) {
	log.Println(string(body))
}

func Client() {

	client := comm.NewClient(IP + ":" + PORT)
	client.TlsEnable(CA, CLIENT_CERT, CLIENT_KEY)
	client.RegHandler(0, clienthandler)
	client.Start(1, 10)

	sendbuf := []byte("hello world from client!")
	err := client.SendMsg(0, sendbuf)
	if err != nil {
		log.Println(err.Error())
		return
	}

	time.Sleep(time.Second)

	client.Stop()
	client.Wait()
}
