package comm

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"log"
	"net"
	"sync"
)

// 服务端消息处理函数
type ServerHandler func(*Server, uint32, []byte)

// server端控制结构
type Server struct {
	taskexit chan bool // 调度线程退出信号
	tasknum  int       // 调度线程数量

	socket  net.Conn //
	conn    *connect
	handler map[uint32]ServerHandler
	wait    sync.WaitGroup
}

// 服务端监听资源结构
type Listen struct {
	listen    net.Listener
	tlsconfig *tls.Config
}

// 监听地址
func NewListen(addr string) *Listen {

	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	return &Listen{listen: listen}
}

func (l *Listen) TlsEnable(ca, cert, key string) error {

	//这里读取的是根证书
	buf, err := ioutil.ReadFile(ca)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	cert_ca_pool := x509.NewCertPool()
	cert_ca_pool.AppendCertsFromPEM(buf)

	//加载服务端证书
	cert_cfg, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	l.tlsconfig = &tls.Config{
		Certificates: []tls.Certificate{cert_cfg},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    cert_ca_pool,
	}

	return nil
}

// 分配一个服务端实例
func (l *Listen) Accept() (*Server, error) {

	conn, err := l.listen.Accept()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	s := new(Server)
	s.handler = make(map[uint32]ServerHandler, 100)

	if l.tlsconfig != nil {
		s.socket = tls.Server(conn, l.tlsconfig)
	} else {
		s.socket = conn
	}

	return s, nil
}

// 服务端处理的调度任务
func msgprocess_server(s *Server) {
	defer s.wait.Done()

	for {

		var msg Header

		select {
		case msg = <-s.conn.RecvBuf:
		case <-s.taskexit:
			{
				return
			}
		}

		fun, b := s.handler[msg.ReqID]
		if b == false {
			log.Println("can not found [", msg.ReqID, "] handler!")
		} else {
			fun(s, msg.ReqID, msg.Body)
		}
	}
}

// 启动消息处理任务
func (s *Server) Start(num, buflen int) error {

	s.conn = NewConnect(s.socket, buflen)
	s.wait.Add(num)
	s.tasknum = num
	s.taskexit = make(chan bool, num)
	for i := 0; i < num; i++ {
		go msgprocess_server(s)
	}
	return nil
}

// 主动停止服务端处理
func (s *Server) Stop() {
	s.conn.Close()
}

// 等待资源销毁
func (s *Server) Wait() {
	s.conn.Wait()
	for i := 0; i < s.tasknum; i++ {
		s.taskexit <- true
	}
	s.wait.Wait()
}

// 发送消息
func (s *Server) SendMsg(reqid uint32, body []byte) error {
	var msg Header

	msg.ReqID = reqid
	msg.Body = make([]byte, len(body))
	copy(msg.Body, body)

	return s.conn.SendMsg(msg)
}

// 注册消息处理函数
func (s *Server) RegHandler(reqid uint32, fun ServerHandler) error {
	_, b := s.handler[reqid]
	if b == true {
		return errors.New("handler id has been register!")
	}
	s.handler[reqid] = fun
	return nil
}
