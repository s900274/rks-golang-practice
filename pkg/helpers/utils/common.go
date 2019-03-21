package utils



import (
    "crypto/md5"
    "encoding/hex"
    "hash/crc32"
    "hash/crc64"
    "math"
    "math/rand"
    "os"
    "path"
    "runtime"
    "sort"
    "strconv"
    "strings"
    "time"
)

func CallerName() string {
    var pc uintptr
    var file string
    var line int
    var ok bool
    if pc, file, line, ok = runtime.Caller(1); !ok {
        return ""
    }
    name := runtime.FuncForPC(pc).Name()
    res := "[" + path.Base(file) + ":" + strconv.Itoa(line) + "]" + name
    return res
}

var rand_gen = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandInt() int {
    return rand_gen.Int()
}

func RandIntn(max int) int {
    return rand_gen.Intn(max)
}

func NowInS() int64 {
    return time.Now().Unix()
}

func NowInNs() int64 {
    return time.Now().UnixNano()
}

func Abs(x int32) int32 {
    switch {
    case x < 0:
        return -x
    case x == 0:
        return 0 // return correctly abs(-0)
    }
    return x
}

func Distance(flat float64, flng float64,
tlat float64, tlng float64) (r int32) {
    distance := math.Sqrt((flat-tlat)*(flat-tlat) + (flng-tlng)*(flng-tlng))
    return int32(distance * 100000)
}

const TIME_FORMAT = "2006-01-02 15:04:05"
const INVALID_TIMEHOUR = "-1"

//timestamp格式化(int32->string(%Y-%M-%D %H:%M:%S))
//Timestamp2str(1341072000)    = "2012-07-01 00:00:00"
func Timestamp2str(timestamp int64) string {
    str_time := ""
    if timestamp > 0 {
        str_time = time.Unix(timestamp, 0).Format(TIME_FORMAT)
    }
    return str_time
}

//从"2012-07-01 10:00:00"时间串得到小时10
func GetHourFromTimeStr(timestr string) string {
    times := strings.Split(timestr, " ")
    if len(times) < 2 {
        return INVALID_TIMEHOUR
    }

    hms := strings.Split(times[1], ":")
    if len(hms) < 3 {
        return INVALID_TIMEHOUR
    }
    hour := strings.TrimPrefix(hms[0], "0")
    return hour
}

//取当前时间戳（微秒）
func GenTimeNowInMicros() int64 {
    return time.Now().UnixNano() / int64(time.Microsecond)
}

//获取当前时间，格式为:"2006-01-02 15:04:05"
func GetNowTime() string {
    return time.Now().Format("2006-01-02 15:04:05")
}

//获取路径串的路径名
//eg: 获取uri路径/gulfstream/realtimeDriverStat/get_driver_loc的路径名get_driver_loc

func GetName(path string) string {
    lst := strings.Split(path, "/")
    if len(lst) < 1 {
        return ""
    }
    name := lst[len(lst)-1]
    return name
}

//bool转int
//true->1; false->0
func Bool2Int(v bool) int {
    if v {
        return 1
    }
    return 0
}

//四舍五入，保留places位小数
func Round(val float64, places int) float64 {
    var t float64
    f := math.Pow10(places)
    x := val * f
    if math.IsInf(x, 0) || math.IsNaN(x) {
        return val
    }
    if x >= 0.0 {
        t = math.Ceil(x)
        if (t - x) > 0.50000000001 {
            t -= 1.0
        }
    } else {
        t = math.Ceil(-x)
        if (t + x) > 0.50000000001 {
            t -= 1.0
        }
        t = -t
    }
    x = t / f

    if !math.IsInf(x, 0) {
        return x
    }

    return t
}

func KeySortToStr(signParams map[string]string) string {
    var keys sort.StringSlice
    for key := range signParams {
        keys = append(keys, key)
    }

    keys.Sort()

    var signStr string
    for i := 0; i < len(keys); i++ {
        signStr = signStr + keys[i] + "=" + signParams[keys[i]]
    }

    return signStr
}

//md5加密
func Md5(str string) string {
    h := md5.New()
    h.Write([]byte(str))
    md5str := hex.EncodeToString(h.Sum(nil))
    return md5str
}

func CRC64(crcStr string) uint64 {
    table := crc64.MakeTable(crc64.ISO)
    return crc64.Checksum([]byte(crcStr), table)
}

func CRC32(crcStr string) uint32 {
    return crc32.ChecksumIEEE([]byte(crcStr))
}

func Exist(fileName string) bool {
    _, err := os.Stat(fileName)
    return err == nil || os.IsExist(err)
}
