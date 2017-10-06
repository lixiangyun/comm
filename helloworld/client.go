package main

import (
	"comm"
	"log"
	"time"
)

func clienthandler(c *comm.Client, reqid uint32, body []byte) {
	log.Println(string(body))
}

func Client() {

	client := comm.NewClient(IP + ":" + PORT)
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
