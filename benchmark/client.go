package main

import (
	"log"
	"time"

	"github.com/lixiangyun/comm"
)

const (
	MIN_BODY_SIZE = 8
)

var flag chan int
var clientStat comm.Stat

// 消息发送的body大小
var sendbuflen = MIN_BODY_SIZE

var bexit bool

// 客户端消息发送、接收 统计显示
func netstat_client() {

	log.Println("banchmark start.")

	laststat := clientStat

	for {

		time.Sleep(time.Second)

		tempstat := clientStat
		tempstat.Sub(laststat)

		log.Printf("Recv %d TPS \t %.3f kB/s \r\n",
			tempstat.RecvCnt,
			float32(tempstat.RecvSize/(1024)))

		log.Printf("Send %d TPS \t %.3f kB/s \r\n",
			tempstat.SendCnt,
			float32(tempstat.SendSize)/(1024))

		laststat = clientStat
	}
}

// 客户端消息处理handler
func clienthandler(c *comm.Client, reqid uint32, body []byte) {
	clientStat.AddCnt(0, 1, 0)
	clientStat.AddSize(0, len(body))
}

// 客户端启动、退出函数
func Client(addr string) {

	log.Println("connect : ", addr)

	// 启动客户端，并且注册消息处理函数
	client := comm.NewClient(addr)
	client.RegHandler(0, clienthandler)
	client.Start(1, 1000)

	// 创建统计协程
	go netstat_client()

	var sendbuf [comm.MAX_BUF_SIZE]byte

	for {
		// 发送消息，并且进行统计
		err := client.SendMsg(0, sendbuf[0:sendbuflen])
		if err != nil {
			log.Println(err.Error())
			return
		}

		clientStat.AddCnt(1, 0, 0)
		clientStat.AddSize(sendbuflen, 0)
	}

	// 销毁client资源
	client.Stop()
	client.Wait()
}
