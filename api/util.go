package api

import (
    "github.com/gin-gonic/gin"
    "github.com/nihileon/ticktak/auth"
    "github.com/nihileon/ticktak/dal"
    "github.com/nihileon/ticktak/log"
    "strconv"
)

func GetPageInfo(c *gin.Context) *dal.PageInfo {
    page := &dal.PageInfo{
        PageNum:   1,
        RecordNum: dal.DefaultRecordNum,
    }
    params := c.Request.URL.Query()

    pnV := params["pn"]
    if len(pnV) > 0 {
        p, err := strconv.Atoi(pnV[0])
        if err == nil && p > 0 {
            page.PageNum = p
        }
    }

    rnV := params["rn"]
    if len(rnV) > 0 {
        r, err := strconv.Atoi(rnV[0])
        if err == nil && r >= 0 {
            page.RecordNum = r
        }
    }
    return page
}

func SecretAuth(c *gin.Context) (string, bool) {
    key := c.Query("secret_key")
    user := c.Query("secret_user")
    if key != SecretPass || user == "" {
        return "", false
    }
    log.GetLogger().Info("req auth by secret key, user: %s", user)
    return user, true
}

func AuthReq(c *gin.Context) (string, bool) {
    user, ok := SecretAuth(c)
    if ok {
        return user, true
    }
    user = auth.GetUser(c)
    if user == "" {
        log.GetLogger().Errorf("req sso auth failed, need login")
        c.JSON(401, nil)
        return "", false
    }
    return user, true
}

type Response struct {
    ErrCode int         `json:"error_code"`
    ErrMsg  string      `json:"error_message"`
    Data    interface{} `json:"data"`
    //PageInfo Page        `json:"page_info"`
}

type RespWithCount struct {
    ErrCode int         `json:"error_code"`
    ErrMsg  string      `json:"error_message"`
    Data    interface{} `json:"data"`
    Count   int         `json:"total"`
}

func doResp(c *gin.Context, data interface{}, err error) {
    errcode := 0
    errmsg := "success"
    if err != nil {
        log.GetLogger().Errorf("%v", err)
        errcode = -1
        errmsg = err.Error()
    }
    c.JSON(200, &Response{
        ErrCode: errcode,
        ErrMsg:  errmsg,
        Data:    data,
    })
}
func doRespWithCount(c *gin.Context, count int, data interface{}, err error) {
    errcode := 0
    errmsg := "success"
    if err != nil {
        log.GetLogger().Errorf("%v", err)
        errcode = -1
        errmsg = err.Error()
        count = 0
    } else {
        log.GetLogger().Info("do resp with count: %d", count)
    }
    c.JSON(200, &RespWithCount{
        ErrCode: errcode,
        ErrMsg:  errmsg,
        Data:    data,
        Count:   count,
    })
}
