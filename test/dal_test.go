package test

import (
    "fmt"
    "github.com/nihileon/ticktak/config"
    "github.com/nihileon/ticktak/dal"
    "github.com/nihileon/ticktak/log"
    "github.com/nihileon/ticktak/models"
    "math/rand"
    "testing"
    "time"
)

var d struct {
    config *config.Config
    usr    string
    pwd    string
}

func TestMain(m *testing.M) {
    var err error
    log.Init()
    d.config, err = config.InitConfig("../conf/conf.yaml")
    if err != nil {
        panic(fmt.Errorf("init config error: %s", err))
    }

    err = dal.InitDB(d.config.MysqlDSN)
    if err != nil {
        panic(err)
    }

    err = dal.InitKV(d.config.RedisAddr, d.config.MemoryOrRedis)
    if err != nil {
        panic(err)
    }

    rand.Seed(time.Now().UnixNano())
    d.usr = "testUser" + string(rand.Int31()%10000)
    d.pwd = "emmm"

    if m.Run() != 0 {
        panic("error")
    }
}

func TestInsertUser(t *testing.T) {
    user := &models.UserInsert{
        Username:    d.usr,
        Password:    d.pwd,
        Description: "emmm",
    }
    _, err := dal.InsertUser(dal.FetchSession(), user)
    if err != nil {
        panic(err)
    }
}

func TestLoginCheck(t *testing.T) {
    user := &models.LoginSelect{
        Username: d.usr,
        Password: d.pwd,
    }
    success, _, err := dal.LoginCheck(dal.FetchSession(), user)
    if err != nil {
        panic(err)
    }

    if !success {
        panic("error")
    }
    user.Password = "wrong password"
    success, _, err = dal.LoginCheck(dal.FetchSession(), user)
    if err == nil || success {
        panic("error")
    }
}

func TestSetGet(t *testing.T) {
    k := string(rand.Int63())
    v := string(rand.Int63())
    dal.KVs.Set(k, v, 2*time.Second)
    nv, err := dal.KVs.Get(k)
    if err != nil {
        panic(err)
    }
    if nv != v {
        panic("error")
    }
}

func TestTaskInsert(t *testing.T) {
    user := &models.TaskInsert{
        Username: d.usr,
        Title:    "emmm test task",
        State:    1,
        Priority: 1,
        Content:  "lallalal ",
        Tag:      "",
        DoneTime: "",
        DDLTime:  "",
    }
    id, err := dal.InsertTask(dal.FetchSession(), user)
    if err != nil {
        panic(err)
    }

    _, err = dal.SelectTasksByTaskID(dal.FetchSession(), id)
    if err != nil {
        panic(err)
    }

}
