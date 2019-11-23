package api

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/nihileon/ticktak/models"
)

func GetDataByTime(c *gin.Context) {
    claims := c.MustGet("claims").(*models.CustomClaims)
    if claims == nil {
        doResp(c, nil, fmt.Errorf("error"))
    }
    doResp(c, claims, nil)
}
