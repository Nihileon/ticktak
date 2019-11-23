package sqlc

import "fmt"

type valType interface{}

type Cond interface {
    ToString() string
}

type EqCond struct {
    field string
    val   valType
}

func Equal(f string, v valType) *EqCond {
    return &EqCond{
        field: f,
        val:   v,
    }
}

func (c *EqCond) ToString() string {
    return fmt.Sprintf("%s = '%s'",
        c.field,
        MysqlRealEscapeString(c.val))
}

type BtCond struct {
    field string
    lval  valType
    rval  valType
}

func Between(f string, l, r valType) *BtCond {
    return &BtCond{
        field: f,
        lval:  l,
        rval:  r,
    }
}

func (c *BtCond) ToString() string {
    return fmt.Sprintf("%s BETWEEN '%s' AND '%s'",
        c.field,
        MysqlRealEscapeString(c.lval),
        MysqlRealEscapeString(c.rval))
}

type InCond struct {
    field string
    vals  []valType
}

func In(f string, val ...valType) *InCond {
    inCond := &InCond{
        field: f,
        vals:  []valType{},
    }
    for _, v := range val {
        inCond.vals = append(inCond.vals, v)
    }
    return inCond
}

func (c *InCond) ToString() string {
    var sql, sep string
    for _, val := range c.vals {
        sql += sep + MysqlRealEscapeString(val)
        sep = ","
    }
    return fmt.Sprintf("%s IN (%s)",
        c.field,
        sql)
}

type LikeCond struct {
    field string
    val   valType
}

func Like(f string, v valType) *LikeCond {
    return &LikeCond{
        field: f,
        val:   v,
    }
}

func (c *LikeCond) ToString() string {
    return fmt.Sprintf("%s LIKE '%s'",
        c.field,
        MysqlRealEscapeString(c.val))
}

type ExpCond struct {
    expression string
}

func Exp(exp string) *ExpCond {
    return &ExpCond{
        expression: exp,
    }
}

func (c *ExpCond) ToString() string {
    return c.expression
}
