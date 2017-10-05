package comm

import (
	"errors"
	"log"
	"net"
	"sync"
)

type ClientHandler func(c *Client, reqid uint32, body []byte)

type Client struct {
	addr    string
	conn    *connect
	handler map[uint32]ClientHandler
	wait    sync.WaitGroup
}

func NewClient(addr string) *Client {
	c := Client{addr: addr}
	c.handler = make(map[uint32]ClientHandler, 100)
	return &c
}

func (s *Client) RegHandler(reqid uint32, fun ClientHandler) error {
	_, b := s.handler[reqid]
	if b == true {
		return errors.New("channel has been register!")
	}
	s.handler[reqid] = fun
	return nil
}

func msgprocess_client(c *Client) {

	defer c.wait.Done()

	for {

		msg, err := c.conn.RecvBuf()
		if err != nil {
			log.Println(err.Error())
			return
		}

		fun, b := c.handler[msg.ReqID]
		if b == false {
			log.Println("can not found [", msg.ReqID, "] handler!")
		} else {
			fun(c, msg.ReqID, msg.Body)
		}
	}
}

func (c *Client) Start(num int) error {

	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return err
	}

	c.conn = NewConnect(conn, 10000)

	c.wait.Add(num)
	for i := 0; i < num; i++ {
		go msgprocess_client(c)
	}

	return nil
}

// client结构资源销毁
func (c *Client) Stop() {
	c.conn.Close()
	c.wait.Wait()
}

// 发送消息结构
func (c *Client) SendMsg(reqid uint32, body []byte) error {
	var msg Header

	msg.ReqID = reqid
	msg.Body = make([]byte, len(body))
	copy(msg.Body, body)

	return c.conn.SendBuf(msg)
}
