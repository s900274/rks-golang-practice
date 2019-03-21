package middleware

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/internal/define"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/pkg/helpers/utils"
    logger "github.com/shengkehua/xlog4go"
)

func PanicHandle() gin.HandlerFunc {
    return func(c *gin.Context) {
    defer func() {
        err := recover()
        if nil != err {
            stackInfo := utils.GetStackInfo()
            msg := fmt.Sprintf(">>> ```%s``` \n\n `%v`", stackInfo, err)
            logger.Error("Panic : %v", msg)
            utils.CallSlack(msg, define.SLACK_PANIC_CHANNEL, define.SLACK_PANIC_SENDFROM)
        }
    }()
        c.Next()
    }
}

