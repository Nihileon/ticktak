package api

import (
    "github.com/gin-gonic/gin"
    "github.com/nihileon/ticktak/dal"
    "github.com/nihileon/ticktak/models"
)

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

func ChangeCurrentUser(c *gin.Context) {
    currentUser, ok := AuthReq(c)
    if !ok {
        return
    }
    req := &models.LoginUpdate{}
    if err := c.BindJSON(req); err != nil {
        doResp(c, nil, err)
        return
    }
    session := dal.FetchSession()
    if err := session.Begin(); err != nil {
        doResp(c, nil, err)
        return
    }
    defer session.Rollback()
    if err := dal.UpdateLoginInfo(session, currentUser, req); err != nil {
        doResp(c, nil, err)
        return
    }
    if req.Username == currentUser {
        err := session.Commit()
        doResp(c, "successfully change your username or password", err)
    }
    if err := dal.UpdateTaskUsername(session, currentUser, req.Username); err != nil {
        doResp(c, nil, err)
        return
    }
    err := session.Commit()
    doResp(c, "successfully change your username or password", err)
}
