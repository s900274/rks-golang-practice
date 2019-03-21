package middleware

import (
    "github.com/gin-gonic/gin"
    "strings"
    "fmt"
)

func SetReqKey() gin.HandlerFunc {
    return func(c *gin.Context) {
        u := strings.Replace(c.Request.URL.String(), "/", "_", -1)
        m := c.Request.Method

        key := fmt.Sprintf("%s_%s", u, m)
        c.Set("SetReqKey", key)
    }
}
