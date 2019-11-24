package api

import "github.com/gin-gonic/gin"

const (
    HomePageUrl = "http://aliyun.nihil.top"
)



func RedirectToCloud(c *gin.Context) {
    w := c.Writer
    w.Header().Set("Location", HomePageUrl)
    w.WriteHeader(302)
    return
}

func GetCurrUser(c *gin.Context) {
    user, ok := AuthReq(c)
    if !ok {
        return
    }
    type UserInfo struct {
        User string `json:"user"`
    }
    info := &UserInfo{
        User: user,
    }
    doResp(c, info, nil)
}
