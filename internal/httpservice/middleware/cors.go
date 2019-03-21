package middleware

import "github.com/gin-gonic/gin"

func Cors() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization, Accept-Encoding, UserAccount, DF_Origin");
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, PATCH, DELETE, OPTIONS");
        c.Next()
    }
}

