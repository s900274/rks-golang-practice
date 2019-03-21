package main

import (
    "errors"
    "flag"
    "github.com/BurntSushi/toml"
    logger "github.com/shengkehua/xlog4go"
    "log"
    "os"
    "rks-golang-practice/internal/define"
    "rks-golang-practice/internal/httpservice"
    "sync"
)


func initConf(flagSet *flag.FlagSet) error {

    flagSet.Parse(os.Args[1:])
    configFile := flagSet.Lookup("config").Value.String()
    if configFile != "" {
        _, err := toml.DecodeFile(configFile, &define.Cfg)
        if err != nil {
            log.Fatalf("ERROR: failed to load config file %s - %s\n", configFile, err.Error())
            return err
        }

    } else {
        log.Fatalln("ERROR: config file is nil")
        err := errors.New("ERROR: config file is nil")
        return err
    }
    return nil
}

func initLogger() error {
    err := logger.SetupLogWithConf(define.Cfg.LogFile)
    return err
}


func RunHttpServer() {
    var wg sync.WaitGroup

    wg.Add(1)
    go func() {
        defer wg.Done()
        HServer := httpservice.NewHTTPServer()
        err := HServer.InitHttpServer()
        if nil != err {
            logger.Error("HTTP ServerStart failed, err :%v", err)
            return
        }
    }()
    wg.Wait()
}