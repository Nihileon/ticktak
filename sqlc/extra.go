package sqlc

import "fmt"

type Extra interface {
    ToString() string
}

type LimitEx struct {
    begin int
    end   int
}

func Limit(a, b int) *LimitEx {
    return &LimitEx{
        begin: a,
        end:   b,
    }
}

func (e *LimitEx) ToString() string {
    return fmt.Sprintf("LIMIT %d, %d",
        e.begin,
        e.end)
}

type DIYEx struct {
    diy string
}

func DIY(diy string) *DIYEx {
    return &DIYEx{
        diy: diy,
    }
}

func (e *DIYEx) ToString() string {
    return e.diy
}

type OrderByEx struct {
    rule string
}

func OrderBy(rule string) *OrderByEx {
    return &OrderByEx{
        rule: rule,
    }
}

func (e *OrderByEx) ToString() string {
    return "ORDER BY " + e.rule
}
