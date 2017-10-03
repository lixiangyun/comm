package main

import (
	"comm"
	"log"
	"os"
	"sync"
	"time"
)

const (
	IP   = "localhost"
	PORT = "6565"
)

const (
	MIN_BODY_SIZE = 8
	MAX_BODY_SIZE = comm.MAX_BUF_SIZE / 2
)

var flag chan int
var servertable []*comm.Server
var client *comm.Client

func serverhandler(s *comm.Server, reqid uint32, body []byte) {

	err := s.SendMsg(reqid, body)
	if err != nil {
		log.Println(err.Error())
		return
	}

	recvmsgcnt++
	recvmsgsize += len(body)

	sendmsgcnt++
	sendmsgsize += len(body)
}

var recvmsgcnt int
var recvmsgsize int

var sendmsgcnt int
var sendmsgsize int

var sendbuflen = MIN_BODY_SIZE

type banchmark struct {
	sendbuflen  int
	sendmsgsize int
	sendmsgcnt  int
}

var bexit bool
var banchmarktest [32]banchmark

func netstat_client(exit *sync.WaitGroup) {

	defer exit.Done()

	num := 0

	log.Println("banch mark test start...")

	lastrecvmsgcnt := recvmsgcnt
	lastrecvmsgsize := recvmsgsize

	lastsendmsgcnt := sendmsgcnt
	lastsendmsgsize := sendmsgsize

	for {

		time.Sleep(time.Second)

		log.Printf("Recv %d cnt/s , %.3f MB/s \r\n",
			recvmsgcnt-lastrecvmsgcnt,
			float32(recvmsgsize-lastrecvmsgsize)/(1024*1024))

		log.Printf("Send %d cnt/s , %.3f MB/s \r\n",
			sendmsgcnt-lastsendmsgcnt,
			float32(sendmsgsize-lastsendmsgsize)/(1024*1024))

		banchmarktest[num].sendbuflen = sendbuflen
		banchmarktest[num].sendmsgsize = recvmsgsize - lastrecvmsgsize
		banchmarktest[num].sendmsgcnt = sendmsgcnt - lastsendmsgcnt

		if sendbuflen*2 < MAX_BODY_SIZE {
			sendbuflen = sendbuflen * 2
		} else {
			sendbuflen = MIN_BODY_SIZE
		}

		num++

		lastrecvmsgcnt = recvmsgcnt
		lastrecvmsgsize = recvmsgsize

		lastsendmsgcnt = sendmsgcnt
		lastsendmsgsize = sendmsgsize

		if num >= len(banchmarktest) {
			log.Println("banch mark test end.")
			break
		}
	}

	for _, v := range banchmarktest {
		log.Printf("bufLen %d , cnt %d , size %.3f MB/s \r\n",
			v.sendbuflen,
			v.sendmsgcnt,
			float32(v.sendmsgsize)/(1024*1024))
	}

	bexit = true
}

func netstat_server() {

	lastrecvmsgcnt := recvmsgcnt
	lastrecvmsgsize := recvmsgsize

	lastsendmsgcnt := sendmsgcnt
	lastsendmsgsize := sendmsgsize

	for {

		time.Sleep(time.Second)

		log.Printf("Recv %d cnt/s , %.3f MB/s \r\n",
			recvmsgcnt-lastrecvmsgcnt,
			float32(recvmsgsize-lastrecvmsgsize)/(1024*1024))

		log.Printf("Send %d cnt/s , %.3f MB/s \r\n",
			sendmsgcnt-lastsendmsgcnt,
			float32(sendmsgsize-lastsendmsgsize)/(1024*1024))

		lastrecvmsgcnt = recvmsgcnt
		lastrecvmsgsize = recvmsgsize

		lastsendmsgcnt = sendmsgcnt
		lastsendmsgsize = sendmsgsize
	}
}

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

var recv_no uint64

func clienthandler(c *comm.Client, reqid uint32, body []byte) {

	recvmsgcnt++
	recvmsgsize += len(body)

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

	for {

		version++
		comm.PutUint64(version, sendbuf[0:])

		err := client.SendMsg(0, sendbuf[0:sendbuflen])
		if err != nil {
			log.Println(err.Error())
			return
		}

		sendmsgcnt++
		sendmsgsize += sendbuflen

		if bexit == true {
			break
		}
	}

	exit.Wait()

	client.Stop()
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
