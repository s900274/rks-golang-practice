/*
监控守护进程

20160606
@aosen
*/
package monitor

import (
	"fmt"
	"log"
	"sync"
	"time"
    "strings"
    "encoding/json"
)

func InitMonitor() {
	initMonitor()
}

const (
	ADD = 0
	MAX = 1
)

const (
	//时间周期
	CYCLE = 10 * time.Second
)

var mr *Monitor

//服务给monitor传输的数据
type Msg struct {
	src   interface{}
	mtype int
	key   string
	value int64
}

//每次reset前将infodata=sumdata
type Monitor struct {
    //用于展示的业务统计数据
    StatInfo []*StatData
	//给info接口提供的数据
	InfoData *Data
	//用于累加
	SumData *Data
	//用于异步数据通信
	queue chan *Msg
	//读写锁
	mlocker *sync.RWMutex
}

func copyValue(src *map[string]int64, dest *map[string]int64) {

    if nil != dest {
        *dest = *src
    }
}

func addValue(inter interface{}, key string, v int64) {

    src, ok := inter.(map[string]int64)
    if !ok {
        return
    } else {
        src[key] += v
    }
}

func getValue(inter interface{}, key string) interface{} {

    src, ok := inter.(map[string]int64)
    if !ok {
        return nil
    }

    if val, ok := src[key]; !ok {
        return nil
    } else {
        return val
    }
}

//ture 赋值成功， false 赋值失败
func setValue(inter interface{}, key string, v int64) bool {
    src, ok := inter.(map[string]int64)
    if !ok {
        return false
    }

    src[key] = v
    return true
}

func less(inter interface{}, key string, v int64) bool {

    src, ok := inter.(map[string]int64)
    if !ok {
        return false
    }

    if val, ok := src[key]; !ok {
        return true
    } else if val < v {
        return true
    } else {
        return false
    }
}

func addMaxValue(inter interface{}, key string, v int64) {
	if less(inter, key, v) {
		setValue(inter, key, v)
	}
}

func initData(d *Data) {
	d.Count = make(map[string]int64)
	d.Timeused = make(map[string]int64)
    d.TimeMax = make(map[string]int64)
	d.Resource = make(map[string]int64)
	d.Err = make(map[string]int64)
}

//初始化监控
func initMonitor() {
	infodata := &Data{}
	sumdata := &Data{}
	//初始化，全置0
	initData(infodata)
	initData(sumdata)
	mr = &Monitor{
		InfoData: infodata,
		SumData:  sumdata,
		queue:    make(chan *Msg, 10000000),
		mlocker:  new(sync.RWMutex),
	}
	go mr.run()
}

func (self *Monitor) rest() {
	initData(self.SumData)
}

func (self *Monitor) copydata() {

	copyValue(&self.SumData.Count, &self.InfoData.Count)
	copyValue(&self.SumData.Timeused, &self.InfoData.Timeused)
    copyValue(&self.SumData.TimeMax, &self.InfoData.TimeMax)
	copyValue(&self.SumData.Err, &self.InfoData.Err)
    copyValue(&self.SumData.Resource, &self.InfoData.Resource)

    self.StatInfo = make([]*StatData, 0)

    for k, v := range self.InfoData.Count {
        stat := &StatData{}
        names := strings.Split(k, "_")
        stat.Type = names[1]
        stat.Cmd = names[2]
        stat.Cnt = v
        if _, ok := self.InfoData.Timeused[k]; !ok {
            stat.Avg = 0
        } else {
            stat.Avg = self.InfoData.Timeused[k] / v / int64(time.Millisecond)
        }
        if _, ok := self.InfoData.Err[k]; !ok {
            stat.Err = 0
        } else {
            stat.Err = self.InfoData.Err[k]
        }
        if _, ok := self.InfoData.TimeMax[k]; !ok {
            stat.Max = 0
        } else {
            stat.Max = self.InfoData.TimeMax[k] / int64(time.Millisecond)
        }

        self.StatInfo = append(self.StatInfo, stat)
    }
}

func (self *Monitor) run() {
	defer func() {
		fmt.Errorf("msg=[monitor run panic]\n")
		if err := recover(); err != nil {
			fmt.Errorf("msg=[maybach monitor panic in run] detail=[%v]\n", err)
		}
	}()
	//每30秒重置一次
	go func() {
		for {
			time.Sleep(CYCLE)
			start := time.Now().UnixNano()
			mr.mlocker.Lock()
			//记录chan大小
			setValue(mr.SumData.Resource, "MonitorChanSize", int64(len(mr.queue)))
            setValue(mr.InfoData.Resource, "MonitorChanSize", int64(len(mr.queue)))
			mr.copydata()
			mr.rest()
			mr.mlocker.Unlock()
			end := time.Now().UnixNano()
			cost := (end - start)/1000000
			log.Printf("msg=[maybach monitor copy] timeused=[%v]\n", cost)
		}
	}()
	for {
		msg := <-mr.queue
		mr.mlocker.Lock()
		switch msg.mtype {
		case ADD:

			addValue(msg.src, msg.key, msg.value)
		case MAX:
			addMaxValue(msg.src, msg.key, msg.value)
		}
		mr.mlocker.Unlock()
	}
}

func Get() []string {
	mr.mlocker.RLock()
    retStr := make([]string, 0)
    for _, stat := range mr.StatInfo {
        jsonv, _ := json.Marshal(stat)
        retStr = append(retStr, string(jsonv))
    }

    jsonv, _ := json.Marshal(mr.InfoData.Resource)
    retStr = append(retStr, string(jsonv))

	mr.mlocker.RUnlock()
	return retStr
}

//+1 计数  增加某个key的计数
func AddCount(key string) {
	select {
	case mr.queue <- &Msg{mr.SumData.Count, ADD, key, int64(1)}:
	default:
		log.Println("msg=[AddCount fail] detail=[queue full]")
	}
}

//添加耗时
func AddTimeUsed(key string, value int64) {
	select {
	case mr.queue <- &Msg{mr.SumData.Timeused, ADD, key, value}:
	default:
		log.Println("msg=[AddCount fail] detail=[queue full]")
	}
}

//添加最大耗时
func AddMaxTimeUsed(key string, value int64) {
	select {
	case mr.queue <- &Msg{mr.SumData.TimeMax, MAX, key, value}:
	default:
		log.Println("msg=[AddCount fail] detail=[queue full]")
	}
}

//添加错误计数 +1
func AddError(key string) {
	select {
	case mr.queue <- &Msg{mr.SumData.Err, ADD, key, int64(1)}:
	default:
		log.Println("msg=[AddError fail] detail=[queue full]")
	}
}

//添加多个计数 +n
func AddCountMulti(key string, v int64) {
	select {
	case mr.queue <- &Msg{mr.SumData.Count, ADD, key, int64(v)}:
	default:
		log.Println("msg=[AddCountMulti fail] detail=[queue full]")
	}
}

//添加多个错误计数 +n
func AddErrorMulti(key string, v int64) {
	select {
	case mr.queue <- &Msg{mr.SumData.Err, ADD, key, int64(v)}:
	default:
		log.Println("msg=[AddErrorMulti fail] detail=[queue full]")
	}
}

//添加总计数及错误计数
func AddCountAndErrCnt(key string, errKey string, allCount int64, rightCount int64) {
	if allCount <= 0 {
		return
	}
	AddCountMulti(key, allCount)

	if rightCount >= allCount {
		return
	}
	AddErrorMulti(errKey, allCount-rightCount)
}

//直接赋值，不累加 true 赋值成功， false 赋值失败
func SetCount(key string, v int64) bool {
	mr.mlocker.Lock()
	ok1 := setValue(mr.InfoData.Count, key, v)
	ok2 := setValue(mr.SumData.Count, key, v)
	mr.mlocker.Unlock()
	if ok1 && ok2 {
		return true
	}
	return false
}

//直接赋值，不累加 true 赋值成功， false 赋值失败
func SetResourceCount(key string, v int64) bool {
	mr.mlocker.Lock()
	ok1 := setValue(mr.InfoData.Resource, key, v)
	ok2 := setValue(mr.SumData.Resource, key, v)
	mr.mlocker.Unlock()
	if ok1 && ok2 {
		return true
	}
	return false
}

//直接赋值，不累加 true 赋值成功， false 赋值失败
func SetTimeUsed(key string, v int64) bool {
	mr.mlocker.Lock()
	ok1 := setValue(mr.InfoData.Timeused, key, v)
	ok2 := setValue(mr.SumData.Timeused, key, v)
	mr.mlocker.Unlock()
	if ok1 && ok2 {
		return true
	}
	return false
}

//直接赋值，不累加 true 赋值成功， false 赋值失败
func SetErrCount(key string, v int64) bool {
	mr.mlocker.Lock()
	ok1 := setValue(mr.InfoData.Err, key, v)
	ok2 := setValue(mr.SumData.Err, key, v)
	mr.mlocker.Unlock()
	if ok1 && ok2 {
		return true
	}
	return false
}

//获取数据
func GetCount(key string) interface{} {
	mr.mlocker.RLock()
	ret := getValue(mr.InfoData.Count, key)
	mr.mlocker.RUnlock()
	return ret
}
