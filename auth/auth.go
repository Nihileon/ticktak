package auth

import (
    "fmt"
    jwt_go "github.com/dgrijalva/jwt-go"
    "github.com/gin-gonic/gin"
    "github.com/nihileon/ticktak/dal"
    "github.com/nihileon/ticktak/log"
    "github.com/nihileon/ticktak/middlewares"
    "github.com/nihileon/ticktak/models"
    "regexp"
    "time"
)

var validUsername = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

type Response struct {
    ErrCode int         `json:"error_code"`
    ErrMsg  string      `json:"error_message"`
    Data    interface{} `json:"data"`
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

func GetUser(c *gin.Context) string {
    claims := c.MustGet("claims").(*models.CustomClaims)
    if claims == nil {
        return ""
    }
    return claims.Username
}

func RegisterUser(c *gin.Context) {
    req := &models.UserInsert{}
    if err := c.BindJSON(req); err != nil {
        doResp(c, nil, err)
        return
    }
    if validUsername.MatchString(req.Username) == false {
        doResp(c, nil, fmt.Errorf("invalid username"))
        return
    }
    log.GetLogger().Info("add user: %s", req.Username)
    _, err := dal.InsertUser(dal.FetchSession(), req)
    if err != nil {
        doResp(c, nil, err)
    }
    doResp(c, "successfully registered", err)
}

type LoginResult struct {
    Token       string `json:"token"`
    LoginSelect models.UserSelect
}

func Login(c *gin.Context) {
    req := &models.LoginSelect{}
    if err := c.BindJSON(req); err != nil {
        doResp(c, nil, err)
        return
    }
    isPass, user, err := dal.LoginCheck(dal.FetchSession(), req)
    if isPass {
        generateToken(c, user)
        return
    }
    doResp(c, nil, err)
}

// 生成令牌
func generateToken(c *gin.Context, user *models.UserSelect) {
    j := &middlewares.JWT{
        []byte("nihileon"),
    }
    claims := models.CustomClaims{
        Username: user.Username,
        StandardClaims: jwt_go.StandardClaims{
            NotBefore: time.Now().Unix() - 1000,
            ExpiresAt: time.Now().Unix() + 3600*12,
            Issuer:    "nihileon",
        },
    }

    token, err := j.CreateToken(claims)

    if err != nil {
        doResp(c, nil, fmt.Errorf("login failed"))
        return
    }

    log.GetLogger().Info(token)

    data := LoginResult{
        LoginSelect: *user,
        Token:       token,
    }
    doResp(c, data, nil)
}
