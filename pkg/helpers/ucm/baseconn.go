/*
连接管理器的父类

20160527
@aosen
*/
package ucm

import (
	"errors"
	"math/rand"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

const MaxRetryAfterCoolDown = 3

//地址属性
type addrAttr struct {
	//是否为不健康节点，
	health bool
	//危险计数器，当计数器达到阈值后将health置为false
	count int64
	//冷却时间点
	cooldowntime int64
	//地址当前总连接数
	conncount int
}

//连接属性
type connAttr struct {
	addr     string // 连接对应的地址
	inpool   bool   //是否在池中 true 在 false 不在
	falsecnt *int64 //出错数
}

type BaseConn struct {
	r *rand.Rand
	//连接列表, true为当前可用addr
	Addrs map[string]*addrAttr
	//合法地址列表
	al      []string
	alocker *sync.RWMutex
	//健康阈值，达到阈值将节点置为不健康节点 false
	HealthyThreshold int64
	//最大冷却时长
	MaxCooldownTime int64
	//最小健康比, 防止全部摘除
	MinHealthyRatio float64
	//连接超时时间
	ConnSvrTimeout time.Duration
	//重试次数
	RetryTimes int
	//健康检查周期
	CheckAddrCycle time.Duration
	//连接与地址的映射关系
	MapConn map[interface{}]*connAttr
	mlocker *sync.RWMutex
	//创建客户端连接方法
	Create func(addr string, timeout time.Duration) (interface{}, error)
	//连接 是否关闭
	IsOpen func(c interface{}) bool
	//关闭连接
	Down func(c interface{})

}

//根据客户端获取MapConn
func (self *BaseConn) GetMapConn(c interface{}) (string, bool, int, bool) {
	self.mlocker.RLock()
	defer self.mlocker.RUnlock()
	if connattr, ok := self.MapConn[c]; ok {
		return connattr.addr, connattr.inpool, int(*connattr.falsecnt), ok
	} else {
		return "", false, 0, false
	}
}

//设置MapConn 是否在池中的状态
func (self *BaseConn) SetMapConnState(c interface{}, in bool) {
	self.mlocker.Lock()
	defer self.mlocker.Unlock()
	if attr, ok := self.MapConn[c]; ok {
		attr.inpool = in
	}
}

//设置MapConn 的addr
func (self *BaseConn) SetMapConnAddr(c interface{}, addr string) {
	self.mlocker.Lock()
	defer self.mlocker.Unlock()
	if _, ok := self.MapConn[c]; !ok {
		falsecnt := int64(0)
		self.MapConn[c] = &connAttr{
			addr:     addr,
			inpool:   false,
			falsecnt: &falsecnt,
		}
	}
}

//连接错误数+1
func (self *BaseConn) AddConnFalseCnt(c interface{}) {
	self.mlocker.RLock()
	defer self.mlocker.RUnlock()
	if attr, ok := self.MapConn[c]; ok {
		atomic.AddInt64(attr.falsecnt, 1)
	}
}

//重置连接错误数
func (self *BaseConn) ResetConnFalseCnt(c interface{}) {
	self.mlocker.RLock()
	defer self.mlocker.RUnlock()
	if attr, ok := self.MapConn[c]; ok {
		atomic.StoreInt64(attr.falsecnt, int64(0))
	}
}

//修改MapConn
func (self *BaseConn) DeleteMapConn(c interface{}) {
	self.mlocker.Lock()
	defer self.mlocker.Unlock()
	delete(self.MapConn, c)
}

//获取MapConn长度
func (self *BaseConn) LenMapConn() int {
	self.mlocker.RLock()
	defer self.mlocker.RUnlock()
	l := len(self.MapConn)
	return l
}

//读addrs 的健康状态
func (self *BaseConn) GetAddr(k string) bool {
	self.alocker.RLock()
	defer self.alocker.RUnlock()
	attr, ok := self.Addrs[k]
	if !ok {
		return ok
	}
	return attr.health
}

//获取整个addrs字典
func (self *BaseConn) GetAllAddr() map[string]*addrAttr {
	self.alocker.RLock()
	defer self.alocker.RUnlock()
	addrs := self.Addrs
	return addrs
}



//添加addrs
func (self *BaseConn) AddAddr(addr string) {
	self.alocker.Lock()
	defer self.alocker.Unlock()
	self.Addrs[addr] = &addrAttr{true, 0, int64(0), 0}
	if self.al == nil {
		self.al = []string{}
	}
	self.al = append(self.al, addr)
	sort.Strings(self.al)
}

//刪除addrs
func (self *BaseConn) DelAddr(addr string) {
	self.alocker.Lock()
	defer self.alocker.Unlock()
	if _, ok := self.Addrs[addr]; ok {
		delete(self.Addrs, addr)
	}
	if self.al != nil {
		var na []string
		for _, v := range self.al {
			if v == addr {
				continue
			} else {
				na = append(na, v)
			}
		}
		self.al = na
	}
}

//添加addr的连接计数
func (self *BaseConn) IncAddrConnCount(addr string) {
	self.alocker.Lock()
	defer self.alocker.Unlock()
	if attr, ok := self.Addrs[addr]; ok {
		attr.conncount++
	}
}

//减少addr的连接计数
func (self *BaseConn) DecAddrConnCount(addr string) {
	self.alocker.Lock()
	defer self.alocker.Unlock()
	if attr, ok := self.Addrs[addr]; ok {
		if attr.count > 0 {
			attr.conncount--
		}
	}
}

//获取addr的连接计数
func (self *BaseConn) GetAddrConnCount(addr string) int {
	self.alocker.RLock()
	defer self.alocker.RUnlock()
	if attr, ok := self.Addrs[addr]; ok {
		return attr.conncount
	}
	return 0
}

//取得合法addr列表
func (self *BaseConn) Getal() []string{
	return self.al
}

//获取每个地址对应的所有连接数
func (self *BaseConn) GetConnCount() map[string]int {
	res := map[string]int{}
	for addr := range self.GetAllAddr() {
		res[addr] = self.GetAddrConnCount(addr)
	}
	return res
}

//设置addrs的冷去时间
func (self *BaseConn) SetAddrDown(addr string) {
	self.alocker.Lock()
	defer self.alocker.Unlock()
	if attr, ok := self.Addrs[addr]; ok {
		attr.cooldowntime = time.Now().Unix()
		attr.health = false
		attr.count = 0
	}
	al := []string{}
	for k, v := range self.Addrs {
		if v.health {
			al = append(al, k)
		}
	}
	sort.Strings(al)
	self.al = al
}

//增加addr的危险计数
func (self *BaseConn) IncAddrCount(addr string) {
	self.alocker.Lock()
	defer self.alocker.Unlock()
	if attr, ok := self.Addrs[addr]; ok {
		attr.count++
	}
}

//减少addr的危险计数
func (self *BaseConn) DecAddrCount(addr string) {
	self.alocker.Lock()
	defer self.alocker.Unlock()
	if attr, ok := self.Addrs[addr]; ok {
		if attr.count > 0 {
			attr.count--
		}
	}
}

//获取addr的err count
func (self *BaseConn) GetAddrCount(addr string) int64 {
	self.alocker.RLock()
	defer self.alocker.RUnlock()
	if attr, ok := self.Addrs[addr]; ok {
		return attr.count
	}
	return int64(0)
}

//获取冷却时间点
func (self *BaseConn) GetCoolDownTime(addr string) int64 {
	self.alocker.RLock()
	defer self.alocker.RUnlock()
	if attr, ok := self.Addrs[addr]; ok {
		return attr.cooldowntime
	}
	return int64(0)
}

//获取健康状态
func (self *BaseConn) GetAddrHealth(addr string) bool {
	self.alocker.RLock()
	defer self.alocker.RUnlock()
	if attr, ok := self.Addrs[addr]; ok {
		return attr.health
	}
	return true
}

//获取健康總數
func (self *BaseConn) GetHealthyCount() int {
	self.alocker.RLock()
	defer self.alocker.RUnlock()
	//获取所有的健康节点数
	healthcount := 0
	for _, attr := range self.Addrs {
		if attr.health {
			healthcount++
		}
	}
	return healthcount
}

//获取健康比
func (self *BaseConn) GetHealthyRatio() float64 {
	self.alocker.RLock()
	defer self.alocker.RUnlock()
	//获取所有的健康节点数
	healthcount := 0
	for _, attr := range self.Addrs {
		if attr.health {
			healthcount++
		}
	}
	return float64(healthcount) / float64(len(self.Addrs))
}

//重置addrs的冷却时间
func (self *BaseConn) ResetAddr(addr string) {
	self.alocker.Lock()
	defer self.alocker.Unlock()
	if attr, ok := self.Addrs[addr]; ok {
		attr.cooldowntime = int64(0)
		attr.count = self.HealthyThreshold - int64(MaxRetryAfterCoolDown)
		attr.health = true
	}
	al := []string{}
	for k, v := range self.Addrs {
		if v.health {
			al = append(al, k)
		}
	}
	sort.Strings(al)
	self.al = al
}

//返回可用IP地址, 返回错误说明没有可用ip
func (self *BaseConn) GetRAddr(addrs map[string]*addrAttr) (string, error) {
	self.alocker.Lock()
	defer self.alocker.Unlock()
	host_len := len(addrs)
	if host_len == 0 {
		return "", errors.New("No address available.")
	} else {
		if len(self.al) != 0 {
			return self.al[self.r.Intn(len(self.al))], nil
		}
	}
	return "", errors.New("addrs is empty.")
}

func (self *BaseConn) GetHealthy() map[string]bool {
	healthy := map[string]bool{}
	for k, v := range self.GetAllAddr() {
		healthy[k] = v.health
	}
	return healthy
}


//獲取addr 屬性
func (self *BaseConn) GetAddrAttr(addr string) (bool, int64, int64, int){
	attr := self.Addrs[addr]
	return attr.health, attr.count, attr.cooldowntime, attr.conncount
}

/*
checker
主要功能：
1.判断非健康连接是否超过了冷却周期，如果超过了，将属性count置0 health置true cooldowntime=0
2.判断健康节点中的属性count是否超过阈值，如果超过阈值将health置false，并设置冷却时间点
*/
func (self *BaseConn) Check() {
	for {
		time.Sleep(self.CheckAddrCycle)
		for addr := range self.GetAllAddr() {
			if !self.GetAddrHealth(addr) {
				//判断是否超过冷却周期
				if time.Now().Unix()-self.GetCoolDownTime(addr) > self.MaxCooldownTime {
					self.ResetAddr(addr)
				}
			} else {
				//超过阈值并且大于健康节点最小百分比，则置false，并设置冷却时间点，健康節點大於1, 才允許暫時下線
				if (self.GetAddrCount(addr) > self.HealthyThreshold) && self.GetHealthyCount() > 1 {
					if self.GetHealthyRatio() > self.MinHealthyRatio {
						self.SetAddrDown(addr)
					} else {
						self.ResetAddr(addr)
					}
				}
			}
		}
	}
}


