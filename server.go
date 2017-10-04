package comm

import (
	"errors"
	"log"
	"net"
	"sync"
)

type ServerHandler func(*Server, uint32, []byte)

type Server struct {
	socket  net.Conn
	conn    *connect
	handler map[uint32]ServerHandler
	wait    sync.WaitGroup
}

type Listen struct {
	listen net.Listener
}

func NewListen(addr string) *Listen {

	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	return &Listen{listen: listen}
}

func (l *Listen) Accept() (*Server, error) {

	conn, err := l.listen.Accept()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	s := new(Server)
	s.socket = conn
	s.handler = make(map[uint32]ServerHandler, 100)

	return s, nil
}

func msgprocess_server(s *Server) {
	defer s.wait.Done()

	for {

		msg, err := s.conn.RecvBuf()
		if err != nil {
			log.Println(err.Error())
			return
		}

		fun, b := s.handler[msg.ReqID]
		if b == false {
			log.Println("can not found [", msg.ReqID, "] handler!")
		} else {
			fun(s, msg.ReqID, msg.Body)
		}
	}
}

func (s *Server) Start(num int) error {

	s.conn = NewConnect(s.socket, 1000)
	s.wait.Add(num)

	for i := 0; i < num; i++ {
		go msgprocess_server(s)
	}

	return nil
}

func (s *Server) Stop() {
	s.conn.Close()
	s.wait.Wait()
}

func (s *Server) SendMsg(reqid uint32, body []byte) error {
	var msg Header

	msg.ReqID = reqid
	msg.Body = make([]byte, len(body))
	copy(msg.Body, body)

	return s.conn.SendBuf(msg)
}

func (s *Server) RegHandler(reqid uint32, fun ServerHandler) error {

	_, b := s.handler[reqid]
	if b == true {
		return errors.New("handler id has been register!")
	}
	s.handler[reqid] = fun

	return nil
}
