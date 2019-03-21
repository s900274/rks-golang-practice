package common

/*
全局logid生成器
*/
import (
	"sync/atomic"
	"time"
)

type LogIdGentor int64

var lg LogIdGentor

func init() {
	lg = LogIdGentor(time.Now().UnixNano())
}

func (i *LogIdGentor) GetNextId() int64 {
	return atomic.AddInt64((*int64)(i), 1)
}

func GenLogId() int64 {
	return lg.GetNextId()
}
