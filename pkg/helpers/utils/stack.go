package utils

import (
    "runtime"
    "fmt"
    "net/http"
    logger "github.com/shengkehua/xlog4go"
    "bytes"
    "net/url"
    "io/ioutil"
)

func stack() []byte {
    buf := make([]byte, 1024)
    n := runtime.Stack(buf, false)
    return buf[:n]
}

func GetStackInfo() string {

    stackInfo := stack()

    return fmt.Sprintf("%s", stackInfo)
}

type SpreadCmd struct {
    Uid       string
    Pwd       string
    SendFrom  string
    SendTo    []string
    Channel   string
    SendType  []int
    MsgSource string
}

func CallSlack(msg , channel , sendfrom string) {
    data := url.Values{}
    data.Set("token", "xoxp-206851584374-206851584454-206853563974-94770a3ca62816a44a59734d26ef85c9")
    data.Add("channel", channel)
    data.Add("username", sendfrom)
    data.Add("text", fmt.Sprintf("%s", msg))

    client := &http.Client{}
    r, err := http.NewRequest("POST", "https://slack.com/api/chat.postMessage", bytes.NewBufferString(data.Encode()))
    r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

    if err != nil {
        logger.Error( err.Error())
    }

    resp, err := client.Do(r)
    if err != nil {
        logger.Error(err.Error())
        return
    }

    _ , err = ioutil.ReadAll(resp.Body)
    if err != nil {
        logger.Error(err.Error())
    }
}
