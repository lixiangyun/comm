package main

import (
	"log"
	"sync"
	"time"

	"github.com/lixiangyun/comm"
)

const (
	MIN_BODY_SIZE = 8
	MAX_BODY_SIZE = comm.MAX_BUF_SIZE / 2
)

var flag chan int
var clientStat comm.Stat

// 消息发送的body大小
var sendbuflen = MIN_BODY_SIZE

var bexit bool
var banchmarktest [32]comm.Stat
var banchmarkbuflen [32]int

// 客户端消息发送、接收 统计显示
func netstat_client(exit *sync.WaitGroup) {

	defer exit.Done()

	log.Println("banchmark start.")

	num := 0
	laststat := clientStat

	for {

		time.Sleep(time.Second)

		tempstat := clientStat
		tempstat.Sub(laststat)

		log.Printf("Recv %d kTPS \t %.3f MB/s \r\n",
			tempstat.RecvCnt/1024,
			float32(tempstat.RecvSize/(1024*1024)))

		log.Printf("Send %d kTPS \t %.3f MB/s \r\n",
			tempstat.SendCnt/1024,
			float32(tempstat.SendSize)/(1024*1024))

		banchmarktest[num] = tempstat
		banchmarkbuflen[num] = sendbuflen

		if sendbuflen*2 < MAX_BODY_SIZE {
			sendbuflen = sendbuflen * 2
		} else {
			break
		}

		num++

		laststat = clientStat

		if num >= len(banchmarktest) {
			break
		}
	}

	log.Println("banchmark end.")

	for idx, tempstat := range banchmarktest {

		log.Printf("BufLen %d \t Send %d kTPS \t %.3f MB/s \r\n",
			banchmarkbuflen[idx],
			tempstat.SendCnt/1024,
			float32(tempstat.SendSize)/(1024*1024))

	}

	bexit = true
}

// 客户端消息处理handler
func clienthandler(c *comm.Client, reqid uint32, body []byte) {
	clientStat.AddCnt(0, 1, 0)
	clientStat.AddSize(0, len(body))
}

// 客户端启动、退出函数
func Client() {

	var exit sync.WaitGroup

	flag = make(chan int)

	// 启动客户端，并且注册消息处理函数
	client := comm.NewClient(IP + ":" + PORT)
	client.RegHandler(0, clienthandler)
	client.Start(4, 1000)

	exit.Add(1)
	// 创建统计协程
	go netstat_client(&exit)

	var sendbuf [comm.MAX_BUF_SIZE]byte

	for bexit != true {
		// 发送消息，并且进行统计
		err := client.SendMsg(0, sendbuf[0:sendbuflen])
		if err != nil {
			log.Println(err.Error())
			return
		}

		clientStat.AddCnt(1, 0, 0)
		clientStat.AddSize(sendbuflen, 0)
	}

	// 等待协程退出
	exit.Wait()

	// 销毁client资源
	client.Stop()
	client.Wait()
}
