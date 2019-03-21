/*
声称连接管理器的请求格式
20160418
@aosen
*/
package ucm

import "time"

const (
	//默认连接服务的超时时间
	DefaultConnSvrTimeOut = "100ms"
	//健康检查的时间周期
	DefaultCheckAddrCycle = "3s"
	//默认连接服务的重试次数
	DefaultRetryTimes = 1
	//连接池大小
	DefaultSize = 20
	//动态扩容的最大倍数
	DefaultMultiple = 1
	//默认健康阈值，达到阈值将节点置为不健康节点 false
	DefaultHealthyThreshold = 20
	//默认最大冷却时长
	DefaultMaxCooldownTime = 120
	//默认最小健康比, 防止全部摘除
	DefaultMinHealthyRatio = 0.8
	//默认最大单连接出错数，如果超过这个出错数就close这个连接
	DefaultMaxConnBeFalseCnt = 0
)

//初始化连接池请求格式
type ConnPoolReq struct {
	//服务地址列表
	Addrs []string
	//连接服务的超时时间
	ConnSvrTimeOut time.Duration
	//健康检查的时间周期
	CheckAddrCycle time.Duration
	//连接服务的重试次数
	RetryTimes int
	//最大连接数，连接池大小
	Size int
	//动态扩容的最大倍数, 步长，当连接池连接不够用的时候会动态扩容，每次增加Size大小,不超过Multiple*Size,如果不初始化这个值，默认为1
	Multiple int
	//健康阈值，达到阈值将节点置为不健康节点 false
	HealthyThreshold int64
	//最大冷却时长
	MaxCooldownTime int64
	//最小健康比, 防止全部摘除
	MinHealthyRatio float64
	//默认最大单连接出错数，如果超过这个出错数就close这个连接
	MaxConnBeFalseCnt int
	//创建连接生成客户端的方法, 返回的接口就是Thrift客户端
	Create func(addr string, timeout time.Duration) (interface{}, error)
	//连接 是否关闭
	IsOpen func(c interface{}) bool
	//关闭连接
	Down func(c interface{})
	//健康检查方法
	CheckAddrs func()
}

func (self *ConnPoolReq) GetConnSvrTimeOut() time.Duration {
	if int64(self.ConnSvrTimeOut) <= int64(0) {
		tm, _ := time.ParseDuration(DefaultConnSvrTimeOut)
		return tm
	}
	return self.ConnSvrTimeOut
}

func (self *ConnPoolReq) GetCheckAddrCycle() time.Duration {
	if int64(self.CheckAddrCycle) <= int64(0) {
		tm, _ := time.ParseDuration(DefaultCheckAddrCycle)
		return tm
	}
	return self.CheckAddrCycle
}

func (self *ConnPoolReq) GetRetryTimes() int {
	if self.RetryTimes <= 0 {
		return DefaultRetryTimes
	}
	return self.RetryTimes
}

func (self *ConnPoolReq) GetSize() int {
	if self.Size <= 0 {
		return DefaultSize
	}
	return self.Size
}

func (self *ConnPoolReq) GetMultiple() int {
	if self.Multiple <= 0 {
		return DefaultMultiple
	}
	return self.Multiple
}

func (self *ConnPoolReq) GetHealthyThreshold() int64 {
	if self.HealthyThreshold <= 0 {
		return DefaultHealthyThreshold
	}
	return self.HealthyThreshold
}

func (self *ConnPoolReq) GetMaxCooldownTime() int64 {
	if self.MaxCooldownTime <= 0 {
		return DefaultMaxCooldownTime
	}
	return self.MaxCooldownTime
}

func (self *ConnPoolReq) GetMinHealthyRatio() float64 {
	if self.MinHealthyRatio <= 0 && self.MinHealthyRatio > 1 {
		return DefaultMinHealthyRatio
	}
	return self.MinHealthyRatio
}

func (self *ConnPoolReq) GetMaxConnBeFalseCnt() int {
	if self.MaxConnBeFalseCnt < 0 {
		self.MaxConnBeFalseCnt = DefaultMaxConnBeFalseCnt
	}
	return self.MaxConnBeFalseCnt
}

//初始化短链接管理器请求格式
type ShortConnReq struct {
	//服务地址列表
	Addrs []string
	//连接服务的超时时间
	TimeOut time.Duration
	//连接服务的重试次数
	RetryTimes int
	//健康检查的时间周期
	CheckAddrCycle time.Duration
	//健康阈值，达到阈值将节点置为不健康节点 false
	HealthyThreshold int64
	//最大冷却时长
	MaxCooldownTime int64
	//最小健康比, 防止全部摘除
	MinHealthyRatio float64
	//创建连接生成客户端的方法, 返回的接口就是Thrift客户端
	Create func(addr string, timeout time.Duration) (interface{}, error)
	//连接 是否关闭
	IsOpen func(c interface{}) bool
	//关闭连接
	Down func(c interface{})
	//健康检查方法
	CheckAddrs func()
}

func (self *ShortConnReq) GetTimeOut() time.Duration {
	if int64(self.TimeOut) <= int64(0) {
		tm, _ := time.ParseDuration(DefaultConnSvrTimeOut)
		return tm
	}
	return self.TimeOut
}

func (self *ShortConnReq) GetCheckAddrCycle() time.Duration {
	if int64(self.CheckAddrCycle) <= int64(0) {
		tm, _ := time.ParseDuration(DefaultCheckAddrCycle)
		return tm
	}
	return self.CheckAddrCycle
}

func (self *ShortConnReq) GetRetryTimes() int {
	if self.RetryTimes <= 0 {
		return DefaultRetryTimes
	}
	return self.RetryTimes
}

func (self *ShortConnReq) GetHealthyThreshold() int64 {
	if self.HealthyThreshold <= 0 {
		return DefaultHealthyThreshold
	}
	return self.HealthyThreshold
}

func (self *ShortConnReq) GetMaxCooldownTime() int64 {
	if self.MaxCooldownTime <= 0 {
		return DefaultMaxCooldownTime
	}
	return self.MaxCooldownTime
}

func (self *ShortConnReq) GetMinHealthyRatio() float64 {
	if self.MinHealthyRatio <= 0 && self.MinHealthyRatio > 1 {
		return DefaultMinHealthyRatio
	}
	return self.MinHealthyRatio
}
