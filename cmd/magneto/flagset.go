package main

import (
    //"time"
    "flag"
)

func Flagset() *flag.FlagSet {
    flagSet := flag.NewFlagSet("uranus", flag.ExitOnError)
    flagSet.String("config", "", "path to config file")

    //global flag
    flagSet.String("logfile", "", "log output file")
    flagSet.Int("http-server-port", 9918, "maybach http port")

    flagSet.String("tip", "", "thrift server ip")
    flagSet.Int("tport", 15248, "thrift server port")
    flagSet.Int("ttimeout", 1000, "thrift server timeout")

    // downstream flag
    flagSet.String("servers", "", "lbspp proxy server info")
    flagSet.Int64("healthythreshold", 10, "healthy threshold")
    flagSet.Int64("maxcooldowntime", 30, "max cooldown time")
    flagSet.Float64("minhealthyratio", 1.0, "min healthy ratio")
    flagSet.Int("connsize", 200, "connsize")
    flagSet.String("timeout", "30000ms", "timeout")
    flagSet.String("cycle", "1000ms", "healthy check cycle")
    flagSet.Int("retrytimes", 0, "retry times")
    flagSet.Int("maxconnfalsecnt", 1, "max connection false cnt")
    flagSet.String("checkaliveconnsessioncycle", "5m", "check alive session cycle")
    flagSet.String("connsessionalivetime", "10m", "session alive time")

    return flagSet
}
