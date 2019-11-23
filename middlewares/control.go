package middlewares

import (
    "github.com/gin-gonic/gin"
)

// AllowControl .
func AllowControl() gin.HandlerFunc {
    return func(c *gin.Context) {
        if str := c.GetHeader("Origin"); len(str) > 0 {
            c.Header("Access-Control-Allow-Origin", str)
        } else {
            c.Header("Access-Control-Allow-Origin", "*")
        }
        c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PATCH, PUT, DELETE")
        c.Header("Access-Control-Allow-Credentials", "true")
        c.Header("Access-Control-Allow-Headers", c.Request.Header.Get("Access-Control-Request-Headers"))
        c.Next()
    }
}

// OPTIONSHandle .
func OPTIONSHandle(c *gin.Context) {
    if str := c.GetHeader("Origin"); len(str) > 0 {
        c.Header("Access-Control-Allow-Origin", str)
    } else {
        c.Header("Access-Control-Allow-Origin", "*")
    }
    c.Header("Access-Control-Allow-Methods", c.Request.Header.Get("Access-Control-Request-Method"))
    c.Header("Access-Control-Allow-Credentials", "true")
    c.Header("Access-Control-Allow-Headers", c.Request.Header.Get("Access-Control-Request-Headers"))
}
