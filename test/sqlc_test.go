package test

import (
    "fmt"
    sqlc2 "github.com/nihileon/ticktak/sqlc"
    "testing"
)

func TestSelect(t *testing.T) {
    sqlc := sqlc2.NewSQLc("my_table")
    sqlc.And(
        sqlc2.Equal("name", "testName")).And(
        sqlc2.Between("age", 1, 19)).And(
        sqlc2.In("sex", "male", "female")).Ext(
        sqlc2.Limit(0, 10))

    type selectFields struct {
        id   uint64 `json:"f_id"`
        name string `json:"f_name"`
    }
    selectInputs := selectFields{}
    sql, err := sqlc.ToSelect(selectInputs)
    if err != nil {
        t.Errorf("select failed, %s", err)
    }
    fmt.Printf("select sql: %s\n", sql)
}

func TestInsert(t *testing.T) {
    sqlc := sqlc2.NewSQLc("my_table")

    type inputFields struct {
        Id   uint64  `json:"f_id"`
        Age  float32 `json:"f_age"`
        Name string  `json:"f_name"`
    }

    input := inputFields{
        Id:   234234,
        Age:  19.3,
        Name: "testName",
    }
    sql, err := sqlc.ToInsert(input)
    if err != nil {
        t.Errorf("insert failed, %s", err)
    }
    fmt.Printf("insert SQL: %s\n", sql)
}

func TestUpdate(t *testing.T) {
    sqlc := sqlc2.NewSQLc("my_table")

    type inputFields struct {
        Age  float32 `json:"f_age"`
        Name string  `json:"f_name"`
    }

    sqlc.And(sqlc2.Equal("id", 232))
    input := inputFields{
        Age:  19.3,
        Name: "testName",
    }
    sql, err := sqlc.ToUpdate(input)
    if err != nil {
        t.Errorf("update failed, %s", err)
    }
    fmt.Printf("update SQL: %s\n", sql)

}

func TestDelete(t *testing.T) {
    sqlc := sqlc2.NewSQLc("my_table")

    sqlc.And(sqlc2.Equal("id", 32434))
    sql, err := sqlc.ToDelete()
    if err != nil {
        t.Errorf("delete failed, %s", err)
    }
    fmt.Printf("delete sql: %s\n", sql)
}
