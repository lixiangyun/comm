package comm

type Stat struct {
	RecvCnt  int // 接收统计计数
	SendCnt  int // 发送统计计数
	ErrCnt   int // 错误统计计数
	RecvSize int // 发送大小计数（单位：byte）
	SendSize int // 接收大小计数（单位：byte）
}

// 两个统计差值
func (s1 *Stat) Sub(s2 Stat) Stat {
	s1.SendCnt -= s2.SendCnt
	s1.RecvCnt -= s2.RecvCnt
	s1.ErrCnt -= s2.ErrCnt

	s1.SendSize -= s2.SendSize
	s1.RecvSize -= s2.RecvSize
	return *s1
}

// 两个统计之和
func (s1 *Stat) Add(s2 Stat) Stat {
	s1.SendCnt += s2.SendCnt
	s1.RecvCnt += s2.RecvCnt
	s1.ErrCnt += s2.ErrCnt

	s1.SendSize += s2.SendSize
	s1.RecvSize += s2.RecvSize
	return *s1
}

// 发送、接收、错误；计数统计接口
func (s1 *Stat) AddCnt(send, recv, err int) Stat {
	s1.SendCnt += send
	s1.RecvCnt += recv
	s1.ErrCnt += err
	return *s1
}

// size统计接口
func (s1 *Stat) AddSize(sendsize, recvsize int) Stat {
	s1.SendSize += sendsize
	s1.RecvSize += recvsize
	return *s1
}
