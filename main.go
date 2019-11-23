package main

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/nihileon/ticktak/api"
    "github.com/nihileon/ticktak/auth"
    "github.com/nihileon/ticktak/dal"
    "github.com/nihileon/ticktak/log"
    "github.com/nihileon/ticktak/middlewares"
)

func main() {
    log.Init()

    // config
    config, err := InitConfig()
    if err != nil {
        panic(fmt.Errorf("init config error: %s", err))
    }

    // MySQL
    err = dal.InitDB(config.MysqlDSN)
    if err != nil {
        panic(err)
    }

    // Redis Or Map
    err = dal.InitKV(config.RedisAddr, config.MemoryOrRedis)
    if err != nil {
        panic(err)
    }

    r := gin.Default()
    r.Use(middlewares.AllowControl())

    r.POST("/register", auth.RegisterUser)
    r.POST("/login", auth.Login)

    r.Use(middlewares.JWTAuth())

    //r.GET("/", api.RedirectToCloud)

    r.OPTIONS("api/*pattern", middlewares.OPTIONSHandle)

    platformAPI := r.Group("api")
    platformAPI.GET("/dataByTime", api.GetDataByTime)

    r.Run(":8080")

}
