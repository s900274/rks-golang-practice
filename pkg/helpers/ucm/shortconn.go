/*
通用短连接管理器
20160418
@aosen
*/
package ucm

import (
	"errors"
	"math/rand"
	"sync"
	"time"
)

type ShortConn struct {
	BaseConn
}

//初始化短连接管理器
func NewShortConn(req *ShortConnReq) (*ShortConn, error) {
	l := len(req.Addrs)
	if l < 0 || req.RetryTimes < 0 {
		return nil, errors.New("len(addrs) < 0 or retrytimes < 0")
	}
	baseconn := BaseConn{
		r:                rand.New(rand.NewSource(time.Now().UnixNano())),
		Addrs:            map[string]*addrAttr{},
		alocker:          new(sync.RWMutex),
		HealthyThreshold: req.GetHealthyThreshold(),
		MaxCooldownTime:  req.GetMaxCooldownTime(),
		MinHealthyRatio:  req.GetMinHealthyRatio(),
		ConnSvrTimeout:   req.GetTimeOut(),
		RetryTimes:       req.GetRetryTimes(),
		CheckAddrCycle:   req.GetCheckAddrCycle(),
		MapConn:          map[interface{}]*connAttr{},
		mlocker:          new(sync.RWMutex),
		Create:           req.Create,
		IsOpen:           req.IsOpen,
		Down:             req.Down,
	}
	//新建短连接管理器对象
	sc := &ShortConn{
		BaseConn: baseconn,
	}
	//地址赋值
	for _, addr := range req.Addrs {
		sc.AddAddr(addr)
	}
	go sc.Check()
	return sc, nil
}

//新建一个连接, 重试机制
func (self *ShortConn) newConn(addr string, timeout time.Duration) (interface{}, error) {
	for i := 0; i < 1+self.RetryTimes; i++ {
		client, err := self.Create(addr, timeout)
		if err != nil {
			continue
		} else {
			return client, nil
		}
	}
	return nil, errors.New("connect server fail.")
}

//获取一个新连接
func (self *ShortConn) Get() (interface{}, error) {
	//随机获取可用addr, 非线程安全
	addr, err := self.GetRAddr(self.GetAllAddr())
	if err != nil {
		return nil, err
	}
	//新建连接
	cli, err := self.newConn(addr, self.ConnSvrTimeout)
	if err != nil {
		self.IncAddrCount(addr)
		return nil, err
	} else {
		self.DecAddrCount(addr)
		self.SetMapConnAddr(cli, addr)
		return cli, err
	}
}

//回收连接，对于短连接就是关闭
func (self *ShortConn) Put(c interface{}, safe bool) {
	if _, _, _, ok := self.GetMapConn(c); ok {
		//如果是不安全的连接，增加不安全计数，如果安全减少不安全计数
		/*
			if safe {
				self.DecAddrCount(addr)
			} else {
				self.IncAddrCount(addr)
			}
		*/
		if self.IsOpen(c) {
			self.Down(c)
		}
		//之所以放在这里是防止开发之在管理器外关闭连接
		//就是无论如何都要从mapconn删掉
		self.DeleteMapConn(c)
	}
	return
}

//获取活跃连接数
func (self *ShortConn) ActiveCount() int {
	return self.LenMapConn()
}

//获取连接对应的addr
func (self *ShortConn) GetAddrByConn(c interface{}) (string, bool) {
	addr, _, _, ok := self.GetMapConn(c)
	return addr, ok
}

//如果连接出问题，可以通过ReConn 来进行重连，获取新的连接，将旧的连接从mapconn中删除
//重连需要加互斥锁，防止开发者恶意通过现有合法端获取多个端,造成合法连接过多
func (self *ShortConn) ReConn(c interface{}) (interface{}, error) {
	if c == nil {
		return nil, errors.New("The client is illegal")
	}
	//获取c对应的addr
	addr, in, _, ok := self.GetMapConn(c)
	//如果这个端不在mapconn中，则为非法端，
	if !ok {
		return nil, errors.New("The client is illegal.")
	}
	//如果c已经在池中，则不再为其生成新的端
	if in {
		return nil, errors.New("The client has in pool.")
	}
	var (
		cli interface{}
		err error
	)
	//新建客户端
	for i := 0; i < 1+self.RetryTimes; i++ {
		cli, err = self.Create(addr, self.ConnSvrTimeout)
		if err != nil {
			continue
		} else {
			break
		}
	}
	//如果失败返回错误，成功返回端，并更新mapconn
	if err != nil {
		return nil, errors.New("reconnect fail.")
	} else {
		//更新mapconn，先删除旧端，再加新端
		self.DeleteMapConn(c)
		self.SetMapConnAddr(cli, addr)
		return cli, nil
	}
}
