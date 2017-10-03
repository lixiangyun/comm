package comm

import (
	"errors"
	"log"
	"net"
	"sync"
)

const (
	MAX_BUF_SIZE = 128 * 1024 // 缓冲区大小(单位：byte)
	MAGIC_FLAG   = 0x98b7f30a // 校验魔术字
	MSG_HEAD_LEN = 3 * 4      // 消息头长度
)

type Header struct {
	ReqID uint32 // 请求ID
	Body  []byte // 传输内容
}

type msgHeader struct {
	Flag  uint32 // 魔术字
	ReqID uint32 // 请求ID
	Size  uint32 // 内容长度
	Body  []byte // 传输的内容
}

type connect struct {
	conn    net.Conn       // 链路结构
	wait    sync.WaitGroup // 同步等待退出
	exit    chan bool      // 退出通道
	sendbuf chan Header    // 发送缓冲队列
	recvbuf chan Header    // 接收缓冲队列
}

// 申请链路操作资源
func NewConnect(conn net.Conn, buflen int) *connect {

	c := new(connect)

	c.conn = conn
	c.sendbuf = make(chan Header, buflen)
	c.recvbuf = make(chan Header, buflen)
	c.exit = make(chan bool)

	c.wait.Add(2)

	go c.sendtask()
	go c.recvtask()

	return c
}

// 链路资源销毁操作
func (c *connect) Close() {
	c.conn.Close()
	c.wait.Done()
	close(c.recvbuf)
	close(c.sendbuf)
}

// 发送调度协成
func (c *connect) sendtask() {

	defer c.wait.Done()
	var buf [MAX_BUF_SIZE]byte

	for {

		var buflen int
		var msg Header

		select {
		case msg = <-c.sendbuf:
		case <-c.exit:
			{
				return
			}
		}

		size := len(msg.Body)
		tmpbuf := make([]byte, MSG_HEAD_LEN+size)

		PutUint32(MAGIC_FLAG, tmpbuf[0:])
		PutUint32(msg.ReqID, tmpbuf[4:])
		PutUint32(uint32(size), tmpbuf[8:])
		copy(tmpbuf[12:], msg.Body)

		tmpbuflen := len(tmpbuf)

		if tmpbuflen >= MAX_BUF_SIZE/2 {
			err := fullywrite(c.conn, tmpbuf[0:])
			if err != nil {
				log.Println(err.Error())
				return
			}
		} else {
			copy(buf[0:tmpbuflen], tmpbuf[0:])
			buflen = tmpbuflen
		}

		chanlen := len(c.sendbuf)

		for i := 0; i < chanlen; i++ {

			msg = <-c.sendbuf

			size = len(msg.Body)
			tmpbuf = make([]byte, MSG_HEAD_LEN+size)

			PutUint32(MAGIC_FLAG, tmpbuf[0:])
			PutUint32(msg.ReqID, tmpbuf[4:])
			PutUint32(uint32(size), tmpbuf[8:])
			copy(tmpbuf[12:], msg.Body)

			tmpbuflen = len(tmpbuf)

			copy(buf[buflen:buflen+tmpbuflen], tmpbuf[0:])
			buflen += tmpbuflen

			if buflen >= MAX_BUF_SIZE/2 {
				err := fullywrite(c.conn, buf[0:buflen])
				if err != nil {
					log.Println(err.Error())
					return
				}
				buflen = 0
			}
		}

		if buflen > 0 {
			err := fullywrite(c.conn, buf[0:buflen])
			if err != nil {
				log.Println(err.Error())
				return
			}
		}
	}
}

// 接收调度协成
func (c *connect) recvtask() {

	var buf [MAX_BUF_SIZE]byte
	var totallen int

	defer c.wait.Done()

	for {

		var lastindex int

		recvnum, err := c.conn.Read(buf[totallen:])
		if err != nil {
			log.Println(err.Error())
			return
		}

		totallen += recvnum

		for {

			if lastindex+MSG_HEAD_LEN > totallen {
				copy(buf[0:totallen-lastindex], buf[lastindex:totallen])
				totallen = totallen - lastindex
				break
			}

			var msg msgHeader

			msg.Flag = GetUint32(buf[lastindex:])
			msg.ReqID = GetUint32(buf[lastindex+4:])
			msg.Size = GetUint32(buf[lastindex+8:])

			bodybegin := lastindex + MSG_HEAD_LEN
			bodyend := bodybegin + int(msg.Size)

			if msg.Flag != MAGIC_FLAG {

				log.Println("msghead_0:", msg)
				log.Println("totallen:", totallen)
				log.Println("bodybegin:", bodybegin, " bodyend:", bodyend)
				log.Println("body:", buf[lastindex:bodyend])
				log.Println("bodyFull:", buf[0:totallen])
				log.Println("close connect.")

				c.conn.Close()
				return
			}

			if bodyend > totallen {

				copy(buf[0:totallen-lastindex], buf[lastindex:totallen])
				totallen = totallen - lastindex
				break
			}

			var tempmsg Header

			tempmsg.ReqID = msg.ReqID
			tempmsg.Body = make([]byte, len(buf[bodybegin:bodyend]))
			copy(tempmsg.Body, buf[bodybegin:bodyend])

			c.recvbuf <- tempmsg

			lastindex = bodyend
		}
	}
}

// 发送封装的接口
func fullywrite(conn net.Conn, buf []byte) error {

	totallen := len(buf)
	sendcnt := 0

	for {

		cnt, err := conn.Write(buf[sendcnt:])
		if err != nil {
			return err
		}

		if cnt <= 0 {
			return errors.New("conn write error!")
		}

		if cnt+sendcnt >= totallen {
			return nil
		}

		sendcnt += cnt
	}
}
