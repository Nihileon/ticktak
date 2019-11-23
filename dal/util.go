package dal

import (
    "github.com/nihileon/ticktak/sqlc"
    "time"
)

func timeStampNow() string {
    return time.Now().Format("2006-01-02 15:04:05")
}

type PageInfo struct {
    PageNum   int `json:"pn"`
    RecordNum int `json:"rn"`
}

func (p *PageInfo) ToLimit() *sqlc.LimitEx {
    return sqlc.Limit((p.PageNum-1)*p.RecordNum, p.RecordNum)
}

const (
    DefaultRecordNum = 10
)
