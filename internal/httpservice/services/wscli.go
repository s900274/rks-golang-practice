package services

import (
    "github.com/gin-gonic/gin"
    "net/http"
    "fmt"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/internal/define"
)

func WebSocketClient(c *gin.Context) {

    c.HTML(http.StatusOK, "wscli.tmpl", gin.H{
        "port": fmt.Sprintf("%d", define.Cfg.Http_server_port),
    })


    return
}
