package socketioservice

import (
    "github.com/gin-gonic/gin"
    "github.com/googollee/go-socket.io"
    logger "github.com/shengkehua/xlog4go"
    "log"
)

var SocketioServer *socketio.Server

func SocketioServerInit() (error) {
    SocketioServer, err := socketio.NewServer(nil)
    if nil != err {
        logger.Error("Socket Io Server Start failed, err :%v", err)
        return err
    }

    SocketioServer.On("connection", func(so socketio.Socket) {
        log.Println("on connection")
        so.Join("chat")
        so.On("chat message", func(msg string) {
            log.Println("emit:", so.Emit("chat message", msg))
            SocketioServer.BroadcastTo("chat", "chat message", msg)
        })

        so.On("disconnection", func() {
            log.Println("on disconnect")
        })
    })

    SocketioServer.On("error", func(so socketio.Socket, err error) {
        logger.Error("Failed : %v", err)
    })


    return nil
}

// Handler initializes the prometheus middleware.
func Handler() gin.HandlerFunc {

    logger.Debug("YOYOYOYOYOYO")

    return func(c *gin.Context) {
        if SocketioServer == nil {
            logger.Error("Error nillllll")
        }
        origin := c.GetHeader("Origin")
        c.Header("Access-Control-Allow-Credentials", "true")
        c.Header("Access-Control-Allow-Origin", origin)
        SocketioServer.ServeHTTP(c.Writer, c.Request)
    }
}