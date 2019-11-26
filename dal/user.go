package dal

import (
    "fmt"
    "github.com/nihileon/ticktak/log"
    "github.com/nihileon/ticktak/models"
    "github.com/nihileon/ticktak/sqlc"
)

func LoginCheck(session *Session, req *models.LoginSelect) (bool, *models.UserSelect, error) {
    c := sqlc.NewSQLc(UserTable)
    c.And(sqlc.Equal("username", req.Username)).
        And(sqlc.Equal("password", req.Password))
    loginSelects := []models.UserSelect{}
    err := session.Select(c, &loginSelects)
    if err != nil {
        log.GetLogger().Info("can't find this user: %s, %s, error:", req.Username, req.Password, err)
        log.GetLogger().Error(err)
    }
    if len(loginSelects) != 1 {
        return false, nil, fmt.Errorf("select %d rows", len(loginSelects))
    }
    loginSelect := &loginSelects[0]
    return true, loginSelect, nil
}

func InsertUser(session *Session, user *models.UserInsert) (int64, error) {
    user.CreateTime = timeStampNow()
    user.ModifyTime = timeStampNow()
    c := sqlc.NewSQLc(UserTable)
    id, err := session.Insert(c, *user)
    return id, err
}

func UpdateLoginInfo(session *Session, currentUsername string, newUser *models.LoginUpdate) error {
    c := sqlc.NewSQLc(UserTable)
    c.And(sqlc.Equal("username", currentUsername))
    err := session.Update(c, *newUser)
    if err != nil {
        return err
    }
    return nil
}
