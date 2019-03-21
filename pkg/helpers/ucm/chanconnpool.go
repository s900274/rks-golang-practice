/*
分布式连接池管理器
支持健康检查

20160417
@aosen
*/
package ucm

import (
	"errors"
	"math/rand"
	"sync"
	"time"
    "fmt"
	logger "github.com/shengkehua/xlog4go"
	"reflect"
)

type ChanConnPool struct {
	BaseConn
	//默认最大单连接出错数，如果超过这个出错数就close这个连接
	MaxConnBeFalseCnt int
	//池
	idle   map[string]chan interface{}
	locker *sync.RWMutex
}

//初始化thrift连接池, 重复初始化会生成新的thriftcp
func NewChanConnPool(req *ConnPoolReq) (*ChanConnPool, error) {
	l := len(req.Addrs)
	if l < 0 || req.RetryTimes < 0 {
		return nil, errors.New("len(addrs) < 0 or retrytimes < 0")
	}
	baseconn := BaseConn{
		r:                rand.New(rand.NewSource(time.Now().UnixNano())),
		Addrs:            map[string]*addrAttr{},
		alocker:          new(sync.RWMutex),
		HealthyThreshold: req.HealthyThreshold,
		MaxCooldownTime:  req.MaxCooldownTime,
		MinHealthyRatio:  req.MinHealthyRatio,
		ConnSvrTimeout:   req.ConnSvrTimeOut,
		RetryTimes:       req.RetryTimes,
		CheckAddrCycle:   req.CheckAddrCycle,
		MapConn:          map[interface{}]*connAttr{},
		mlocker:          new(sync.RWMutex),
		Create:           req.Create,
		IsOpen:           req.IsOpen,
		Down:             req.Down,
	}
	//地址赋值
	for _, addr := range req.Addrs {
		baseconn.AddAddr(addr)
	}
	//新建一个thriftcp
	cp := &ChanConnPool{
		BaseConn:          baseconn,
		MaxConnBeFalseCnt: req.GetMaxConnBeFalseCnt(),
		locker:            new(sync.RWMutex),
	}
	//初始化idle
	cp.initIdle(req.Size)
	//生成queue
	/*
		for addr := range cp.GetAllAddr() {
			for i := 0; i < req.Size; i++ {
				client, err := cp.newConn(addr, cp.ConnSvrTimeout)
				if err != nil {
					continue
				}
				cp.SetMapConnAddr(client, addr)
				cp.inIdle(addr, client)
			}
		}
	*/
	go cp.Check()
	return cp, nil
}

//新建一个连接
func (self *ChanConnPool) newConn(addr string, timeout time.Duration) (interface{}, error) {
	for i := 0; i < 1+self.RetryTimes; i++ {
        fmt.Println("xxxxxxxxxx newConn  xxxxxx addr :%s", addr)
		client, err := self.Create(addr, timeout)
		if err != nil {
			continue
		} else {
			return client, nil
		}
	}
	return nil, errors.New("connect server fail.")
}



//获取当前活跃连接数
func (self *ChanConnPool) getActive() int {
	active := self.LenMapConn()
	for addr := range self.getIdle() {
		active = active - self.lenIdle(addr)
	}
	return active
}

//获取某地址的总连接数

//idle的操作
//新建一个idle
func (self *ChanConnPool) initIdle(size int) {
	idle := map[string]chan interface{}{}
	for addr := range self.GetAllAddr() {
		idle[addr] = make(chan interface{}, size)
	}
	self.idle = idle
}

//新增一個addr idle into idleMap
func (self *ChanConnPool) intoIdle(addr string, size int) {
	self.idle[addr] = make(chan interface{}, size)
}

//获取idle
func (self *ChanConnPool) getIdle() map[string]chan interface{} {
	return self.idle
}

//获取idle
func (self *ChanConnPool) GetIdle() map[string]chan interface{} {
	return self.getIdle()
}

//in idle, 只有c不在池中 并且是mapconn的一个key，才能执行inIdle
func (self *ChanConnPool) inIdle(addr string, c interface{}) {
	if self.lenIdle(addr) >= self.capIdle(addr) {
		//这种情况不会发生
		return
	}
	if queue, ok := self.getIdle()[addr]; ok {
		//c不在池中，并且mapconn有c这个key
		queue <- c
		self.SetMapConnState(c, true)
		//self.idlelocker.Unlock()
	}
}

//out idle
func (self *ChanConnPool) outIdle(addr string) (interface{}, error) {
	if queue, ok := self.getIdle()[addr]; ok {
		select {
		case client := <-queue:
			if !self.IsOpen(client) {
				//如果连接已经关闭，则新生成一个连接
				//如果新生成连接失败，则返回错误，并将之前关闭的连接回队列里，防止池泄露
				addr, err := self.GetRAddr(self.GetAllAddr())
				if err != nil {
					logger.Debug("xxxxxxxxxx[Test Conn] do inIdle")
					self.inIdle(addr, client)
					//如果返回错，说明没有可用ip地址了
					return nil, err
				}
				cli, err := self.newConn(addr, self.ConnSvrTimeout)
				if err != nil {
					logger.Debug("xxxxxxxxxx[Test Conn] do inIdle")
					self.inIdle(addr, client)
					return nil, err
				} else {
					logger.Debug("xxxxxxxxxx[Test Conn] do reSet ConnAddr")
                    fmt.Println("xxxxxxxx addr :%s addrConnCount: %d", addr, self.GetAddrConnCount(addr))
					self.DeleteMapConn(client)
					self.SetMapConnAddr(cli, addr)
					return cli, nil
				}
			} else {
				self.SetMapConnState(client, false)
				return client, nil
			}
		default:
			//当前该addr的总连接小于容量
			if self.GetAddrConnCount(addr) < self.capIdle(addr) {
				cli, err := self.newConn(addr, self.ConnSvrTimeout)
				logger.Debug("xxxxxxxxxx[Test Conn] do check new conn")
				if err == nil {
					logger.Debug("xxxxxxxxxx[Test Conn] do set new conn")
					self.SetMapConnAddr(cli, addr)
					self.IncAddrConnCount(addr)
					return cli, nil
				}else {
					return nil, errors.New("New Conn Fail, Please Checkout Client Conn, Addr: " + addr)
				}
			}
			return nil, errors.New("Connection Pool Is Empty.")
		}
	}
	return nil, errors.New("Addr Illegal.")
}

//获取某个addr下的idle长度
func (self *ChanConnPool) lenIdle(addr string) int {
	if queue, ok := self.getIdle()[addr]; ok {
		l := len(queue)
		return l
	}
	return 0
}

//获取某个addr下的idle容量
func (self *ChanConnPool) capIdle(addr string) int {
	if queue, ok := self.getIdle()[addr]; ok {
		c := cap(queue)
		return c
	}
	return 0
}

func (self *ChanConnPool) GetLenIdle(addr string) int {
	return self.lenIdle(addr)

}

//获取连接, 在超时时间内如果没有client返回则报错
func (self *ChanConnPool) Get() (interface{}, error) {
	self.locker.Lock()
	defer self.locker.Unlock()
	//随机获取一个可用addr
	if addr, err := self.GetRAddr(self.GetAllAddr()); err != nil {
		logger.Debug("xxxxxxxxxx[Test Conn] no addr")
		return nil, err
	} else {
		cli, err := self.outIdle(addr)
		if err != nil {
			logger.Debug("xxxxxxxxxx[Test Conn] do IncAddrCount")
			self.IncAddrCount(addr)
			return nil, err
		} else {
			return cli, nil
		}
	}
}

//将连接放回连接池
func (self *ChanConnPool) Put(c interface{}, safe bool) {
	self.locker.Lock()
	defer self.locker.Unlock()
	realsafe := true
	//获取端对应的addr, 存在并且连接不在池中
	if addr, in, falsecnt, ok := self.GetMapConn(c); !in && ok {
		//如果安全就重置该连接的错误次数, 反之+1
		if !safe {
			logger.Debug("xxxxxxxxxx[Test Conn] do AddConnFalseCnt")
			self.AddConnFalseCnt(c)
		} else {
			self.ResetConnFalseCnt(c)
		}
		//如果超过该连接的不安全阈值，那就是真不安全，需要关闭该连接，并将该addr 错误次数 +1
		if falsecnt > self.MaxConnBeFalseCnt-1 {
			realsafe = false
		}
		//如果是不安全的连接，增加不安全计数，如果安全减少不安全计数
		if realsafe {
			self.DecAddrCount(addr)
		} else {
			logger.Debug("xxxxxxxxxx[Test Conn] do IncAddrCount")
			self.IncAddrCount(addr)
		}
		//如果 已经关闭|连接不安全|addr不存在(被reset)，都会重启连接
		if !self.IsOpen(c) || !realsafe || !self.IsInAddrs(addr) {
			logger.Debug("xxxxxxxxxx[Test Conn] do close conn addr")
			//确保连接关闭
			self.Down(c)
			self.DeleteMapConn(c)
			self.DecAddrConnCount(addr)
		} else {
			self.inIdle(addr, c)
		}
	}
	return
}

//該地址是否在Addr列表中
func (self *ChanConnPool) IsInAddrs(addr string) bool{
	result := false
	addrs := reflect.ValueOf(self.GetAllAddr()).MapKeys()
	for _, v := range addrs{
		if addr == v.String() {
			result = true
			break
		}
	}
	return result
}


//获取总的活跃连接数
func (self *ChanConnPool) ActiveCount() int {
	self.locker.RLock()
	defer self.locker.RUnlock()
	return self.getActive()
}

//获取连接对应的addr
func (self *ChanConnPool) GetAddrByConn(c interface{}) (string, bool) {
	self.locker.RLock()
	defer self.locker.RUnlock()
	addr, _, _, ok := self.GetMapConn(c)
	return addr, ok
}

//如果连接出问题，可以通过ReConn 来进行重连，获取新的连接，将旧的连接从mapconn中删除
//重连需要加互斥锁，防止开发者恶意通过现有合法端获取多个端,造成合法连接过多
func (self *ChanConnPool) ReConn(c interface{}) (interface{}, error) {
	self.locker.Lock()
	defer self.locker.Unlock()
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
		/**** 2018/6/6 jim add *****/
		//增加錯誤連線次數
		self.IncAddrCount(addr)
		//close 與client連線
		self.Down(c)
		//刪除在MapConn連線
		self.DeleteMapConn(c)
		//Addr.ConnCount --
		self.DecAddrConnCount(addr)
		/**** 2018/6/6 jim add *****/
		return nil, errors.New("reconnect fail.")
	} else {
		//更新mapconn，先删除旧端，再加新端
		self.DeleteMapConn(c)
		self.SetMapConnAddr(cli, addr)
		return cli, nil
	}
}

//pop在池裡的無效Addr, conn也同時斷掉
func (self *ChanConnPool) popIdleByInvalidAddr(addrs []string) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("popIdleByInvalidAddr, err: %s", err )
		}
	}()
	for _ , addr := range addrs{
		if queue, ok := self.getIdle()[addr]; ok {
            close(queue)
			for client := range queue {
				//fmt.Printf("> queue client:\n", client)
				self.Down(client)
				self.DeleteMapConn(client)
			}
		}
	}
}

//刪除idle[addr]
func (self *ChanConnPool) delIdleAddr(addrs []string) {
	for _ , addr := range addrs {
		if _, ok := self.idle[addr]; ok {
			self.locker.Lock()
			delete(self.idle, addr)
			self.locker.Unlock()
		}
	}
}

//create server conn and into pool
func (self *ChanConnPool) CreateServerIntoDownStream(addr string, size int) {
	self.intoIdle(addr, size)
	self.AddAddr(addr)
}

//DropOut a addr from pool
func (self *ChanConnPool) DropOutAddrFromIdle(addrs []string) {
	self.popIdleByInvalidAddr(addrs)
	self.delIdleAddr(addrs)
}

