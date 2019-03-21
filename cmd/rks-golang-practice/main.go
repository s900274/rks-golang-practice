package main

import (
    logger "github.com/shengkehua/xlog4go"
)

func main() {

    flagSet := Flagset()

    if err := initConf(flagSet); err != nil {
        logger.Error("%s", err.Error())
    }

    if err := initLogger(); err != nil {
        logger.Error("%s", err.Error())
    }
    defer logger.Close()

    RunHttpServer()

    logger.Info("Server exit")
}
