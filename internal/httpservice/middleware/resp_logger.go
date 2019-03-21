package middleware

import (
    "github.com/gin-gonic/gin"
    logger "github.com/shengkehua/xlog4go"
    "bytes"
    "fmt"
    "strings"
)

type bodyLogWriter struct {
    gin.ResponseWriter
    body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
    w.body.Write(b)
    return w.ResponseWriter.Write(b)
}

func RespLog() gin.HandlerFunc {
    return func(c *gin.Context) {
        if strings.Contains(c.GetHeader("Content-type"), "application/json") {
            blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
            c.Writer = blw
            c.Next()
            msg := ""
            key, exist := c.Get("SetReqKey")
            if exist {
                msg = fmt.Sprintf("%s|%v", msg, key)
            }

            ui, exist := c.Get("UserInfo")
            if exist {
                msg = fmt.Sprintf("%s|%v", msg, ui)
            }
            msg = fmt.Sprintf("%s|%s", msg, blw.body.String())
            logger.Info(msg)
        }
    }
}
