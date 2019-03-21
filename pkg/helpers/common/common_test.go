package common

import (
	"testing"
	"time"
)

func TestTimestamp2str(t *testing.T) {
	println(Timestamp2str(time.Now().Unix()))
	println(time.Now().Format(TIME_FORMAT))
}

func TestGetHourFromTimeStr(t *testing.T) {
	println(GetHourFromTimeStr("2016-08-29 16:00:00"))
}
func TestGetName(t *testing.T) {
	println(GetName("/gulfstream/realtimeDriverStat/get_driver_loc        "))
	println(GetName("/gulfstream/realtimeDriverStat/"))
	println(GetName("/gulfstream/realtimeDriverStat"))
}
