package keyformat

import (
    "time"
    "fmt"
)


// 生成计数key
func FormatGridKeyCnt(grid, key string, business, period, ct int64) string {

    if 0 == len(grid) {
        return ""
    }

    if 0 == len(key) {
        return ""
    }

    stime := make_timestamp_by_period(period, ct)
    if 0 == len(stime) {
        return ""
    }

    s := fmt.Sprintf("SDS_%s##CNT_%s_%d_%s", key, grid, business, stime)
    return s
}

// 生成求和key
func FormatGridKeySum(grid, key string, business, period, ct int64) string {
    if 0 == len(grid) {
        return ""
    }

    if 0 == len(key) {
        return ""
    }

    stime := make_timestamp_by_period(period, ct)
    if 0 == len(stime) {
        return ""
    }

    s := fmt.Sprintf("SDS_%s##SUM_%s_%d_%s", key, grid, business, stime)
    return s
}

// 生成去重计数key
func FormatGridKeyUnqCnt(grid, key string, business, period, ct int64) string {
    if 0 == len(grid) {
        return ""
    }

    if 0 == len(key) {
        return ""
    }

    //stime := make_timestamp_by_period(period, ct)
    //if 0 == len(stime) {
    //    return ""
    //}

    s := fmt.Sprintf("SDS_%s##UNQCNT_%s_%d", key, grid, business)
    return s
}

// 生成本轮内出现的格子key
func FormatKeyGrids(key string, business, period, ct int64) string {

    if 0 == len(key) {
        return ""
    }

    stime := make_timestamp_by_period(period, ct)
    if 0 == len(stime) {
        return ""
    }

    s := fmt.Sprintf("SDS_%s##GRIDS_%d_%s", key, business, stime)
    return s
}

// 生成出现的格子key
func FormatKeyGridsNoTs(key string, business, period, ct int64) string {

    if 0 == len(key) {
        return ""
    }

    s := fmt.Sprintf("SDS_%s##GRIDS_%d", key, business)
    return s
}

// 生成本轮内出现的格子key，按城市划分
func FormatKeyGridsSepByCity(key string, business, period, ct, cityid int64) string {

    if cityid < 0 {
        return ""
    }

    if 0 == len(key) {
        return ""
    }

    stime := make_timestamp_by_period(period, ct)
    if 0 == len(stime) {
        return ""
    }

    s := fmt.Sprintf("SDS_%s##GRIDS_%d_%s_%d", key, business, stime, cityid)
    return s
}

// 生成本轮内出现的城市ID
func FormatKeyCities(key string, business, period, ct int64) string {

    if 0 == len(key) {
        return ""
    }

    stime := make_timestamp_by_period(period, ct)
    if 0 == len(stime) {
        return ""
    }

    s := fmt.Sprintf("SDS_%s##CITIES_%d_%s", key, business, stime)
    return s
}

// 根据周期和时间戳生成本轮内的时间字符串
func make_timestamp_by_period(period, ct int64) string {

    if ct <= 0 {
        return ""
    }
    stime := ""
    ctime := time.Unix(ct, 0)
    switch period {
    case 1:
        // 以秒为周期粒度
        stime = fmt.Sprintf("%04d%02d%02d%02d%02d%02d", ctime.Year(), ctime.Month(), ctime.Day(), ctime.Hour(), ctime.Minute(), ctime.Second())
    case 30:
        // 以30秒为周期粒度,将1分钟分成上半分钟和下半分钟
        half_min := "fh"
        if ctime.Second() >= 30 {
            half_min = "sh"
        }
        stime = fmt.Sprintf("%04d%02d%02d%02d%02d%s", ctime.Year(), ctime.Month(), ctime.Day(), ctime.Hour(), ctime.Minute(), half_min)
    case 60:
        // 以60秒为周期粒度
        stime = fmt.Sprintf("%04d%02d%02d%02d%02d", ctime.Year(), ctime.Month(), ctime.Day(), ctime.Hour(), ctime.Minute())
    default:
        return ""
    }
    return stime
}