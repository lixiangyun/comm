package main

import (
	"comm"
	"log"
	"time"
)

const (
	IP   = "localhost"
	PORT = "6565"
)

var serverTable []*comm.Server
var serverstat comm.Stat

// server消息处理handler
func serverhandler(s *comm.Server, reqid uint32, body []byte) {

	err := s.SendMsg(reqid, body)
	if err != nil {
		log.Println(err.Error())
		serverstat.AddCnt(0, 1, 0)

		return
	}

	serverstat.AddCnt(1, 1, 0)
	serverstat.AddSize(len(body), len(body))
}

// server端消息处理的统计显示
func netstat_server() {

	laststat := serverstat

	for {

		time.Sleep(time.Second)

		tempstat := serverstat
		tempstat.Sub(laststat)

		log.Printf("Recv %d kTPS \t %.3f MB/s \r\n",
			tempstat.RecvCnt/1024,
			float32(tempstat.RecvSize/(1024*1024)))

		log.Printf("Send %d kTPS \t %.3f MB/s \r\n",
			tempstat.SendCnt/1024,
			float32(tempstat.SendSize)/(1024*1024))

		laststat = serverstat
	}
}

// 启动server端，并且监听
func Server() {

	list := comm.NewListen(":" + PORT)
	go netstat_server()
	for {
		server, err := list.Accept()
		if err != nil {
			log.Println(err.Error())
			return
		}
		server.RegHandler(0, serverhandler)
		server.Start(1)
	}
}
