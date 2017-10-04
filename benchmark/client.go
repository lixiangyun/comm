package main

import (
	"comm"
	"log"
	"sync"
	"time"
)

var flag chan int
var clientStat comm.Stat

var sendbuflen = MIN_BODY_SIZE

var bexit bool
var banchmarktest [32]comm.Stat
var banchmarkbuflen [32]int

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
			sendbuflen = MIN_BODY_SIZE
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

var recv_no uint64

func clienthandler(c *comm.Client, reqid uint32, body []byte) {

	clientStat.AddCnt(0, 1, 0)
	clientStat.AddSize(0, len(body))

	version := comm.GetUint64(body)

	if version <= recv_no {
		log.Println("error! ", version, recv_no)
	} else {
		recv_no = version
	}
}

func Client() {

	var version uint64
	var exit sync.WaitGroup

	flag = make(chan int)

	client := comm.NewClient(IP + ":" + PORT)
	client.RegHandler(0, clienthandler)
	client.Start(1)

	exit.Add(1)
	go netstat_client(&exit)

	var sendbuf [comm.MAX_BUF_SIZE]byte

	for bexit != true {

		version++
		comm.PutUint64(version, sendbuf[0:])

		err := client.SendMsg(0, sendbuf[0:sendbuflen])
		if err != nil {
			log.Println(err.Error())
			return
		}

		clientStat.AddCnt(1, 0, 0)
		clientStat.AddSize(sendbuflen, 0)
	}

	exit.Wait()

	client.Stop()
}
