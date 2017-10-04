package comm

type Stat struct {
	RecvCnt  int
	SendCnt  int
	ErrCnt   int
	RecvSize int
	SendSize int
}

func (s1 *Stat) Sub(s2 Stat) Stat {
	s1.SendCnt -= s2.SendCnt
	s1.RecvCnt -= s2.RecvCnt
	s1.ErrCnt -= s2.ErrCnt

	s1.SendSize -= s2.SendSize
	s1.RecvSize -= s2.RecvSize
	return *s1
}

func (s1 *Stat) Add(s2 Stat) Stat {
	s1.SendCnt += s2.SendCnt
	s1.RecvCnt += s2.RecvCnt
	s1.ErrCnt += s2.ErrCnt

	s1.SendSize += s2.SendSize
	s1.RecvSize += s2.RecvSize
	return *s1
}

func (s1 *Stat) AddCnt(send, recv, err int) Stat {
	s1.SendCnt += send
	s1.RecvCnt += recv
	s1.ErrCnt += err
	return *s1
}

func (s1 *Stat) AddSize(sendsize, recvsize int) Stat {
	s1.SendSize += sendsize
	s1.RecvSize += recvsize
	return *s1
}
