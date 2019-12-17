package middlewares

import (
    "encoding/json"
    "errors"
    "fmt"
    "github.com/nihileon/ticktak/dal"
    "github.com/nihileon/ticktak/log"
    "github.com/nihileon/ticktak/models"
    "time"

    "github.com/dgrijalva/jwt-go"
    "github.com/gin-gonic/gin"
)

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

// JWTAuth 中间件，检查token
func JWTAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.Request.Header.Get("token")
        if token == "" {
            doResp(c, nil, fmt.Errorf("request without a token"))
            c.Abort()
            return
        }

        log.GetLogger().Info("get token: ", token)

        j := NewJWT()
        // parseToken 解析token包含的信息
        claims, err := j.ParseToken(token)
        if err != nil {
            if err == TokenExpired {
                doResp(c, nil, fmt.Errorf("authorization expired"))
                c.Abort()
                return
            }
            doResp(c, nil, err)
            c.Abort()
            return
        }
        // 继续交由下一个路由处理,并将解析出的信息传递下去
        c.Set("claims", claims)
    }
}

// JWT 签名结构
type JWT struct {
    SigningKey []byte
}

// 一些常量
var (
    TokenExpired     error  = errors.New("token is expired")
    TokenNotValidYet error  = errors.New("token not active yet")
    TokenMalformed   error  = errors.New("that's not even a token")
    TokenInvalid     error  = errors.New("couldn't handle this token:")
    SignKey          string = "nihileon"
)

func NewJWT() *JWT {
    return &JWT{
        []byte(GetSignKey()),
    }
}

func GetSignKey() string {
    return SignKey
}

func SetSignKey(key string) string {
    SignKey = key
    return SignKey
}

func (j *JWT) CreateToken(claims models.CustomClaims) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signed, err := token.SignedString(j.SigningKey)
    if err != nil {
        return "", err
    }
    bytes, err := json.Marshal(&claims)
    if err != nil {
        return "", err
    }
    err = dal.KVs.Set(signed, string(bytes), time.Duration(claims.ExpiresAt-time.Now().Unix())*time.Second)
    return signed, err
}

func (j *JWT) ParseToken(tokenString string) (*models.CustomClaims, error) {
    val, err := dal.KVs.Get(tokenString)
    claims := &models.CustomClaims{}
    err = json.Unmarshal([]byte(val), &claims)
    if val != "" && err == nil {
        return claims, nil
    }

    token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
        return j.SigningKey, nil
    })
    if err != nil {
        if ve, ok := err.(*jwt.ValidationError); ok {
            if ve.Errors&jwt.ValidationErrorMalformed != 0 {
                return nil, TokenMalformed
            } else if ve.Errors&jwt.ValidationErrorExpired != 0 {
                return nil, TokenExpired
            } else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
                return nil, TokenNotValidYet
            } else {
                return nil, TokenInvalid
            }
        }
        return nil, TokenInvalid
    }
    if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
        return claims, nil
    }
    return nil, TokenInvalid
}

func (j *JWT) RefreshToken(tokenString string) (string, error) {
    jwt.TimeFunc = func() time.Time {
        return time.Unix(0, 0)
    }
    token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
        return j.SigningKey, nil
    })
    if err != nil {
        return "", err
    }
    if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
        jwt.TimeFunc = time.Now
        claims.StandardClaims.ExpiresAt = time.Now().Add(12 * time.Hour).Unix()
        return j.CreateToken(*claims)
    }
    return "", TokenInvalid
}
