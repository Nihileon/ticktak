package api

import (
    "fmt"
    "github.com/gin-gonic/gin"
)

func GetDataByTime(c *gin.Context) {
    user, ok := AuthReq(c)
    if !ok {
        return
    }
    if user == "" {
        doResp(c, nil, fmt.Errorf("error"))
    }
    doResp(c, user, nil)
}



