package dal

import (
    "fmt"
    "github.com/nihileon/ticktak/models"
    "github.com/nihileon/ticktak/sqlc"
)

func InsertTask(session *Session, task *models.TaskInsert) (int64, error) {
    task.CreateTime = timeStampNow()
    task.ModifyTime = timeStampNow()
    c := sqlc.NewSQLc(TaskTable)
    id, err := session.Insert(c, *task)
    return id, err
}

func SelectTasksByUsernameState(session *Session, username string, state uint, page *PageInfo) (int, []models.TaskSelect, error) {
    c := sqlc.NewSQLc(TaskTable)
    SQL := fmt.Sprintf("id IN (SELECT id FROM %s WHERE username = '%s' AND state = '%d')",
        TaskTable,
        username,
        state)
    c.And(sqlc.Exp(SQL)).Ext(page.ToLimit())
    count, err := session.Count(c)
    if err != nil {
        return count, nil, err
    }
    tasks := []models.TaskSelect{}
    err = session.Select(c, &tasks)
    return count, tasks, err
}

func SelectTasksByUsernamePriority(session *Session, username string, priority uint, page *PageInfo) (int, []models.TaskSelect, error) {
    c := sqlc.NewSQLc(TaskTable)
    SQL := fmt.Sprintf("id IN (SELECT id FROM %s WHERE username = '%s' AND priority = '%d')",
        TaskTable,
        username,
        priority)
    c.And(sqlc.Exp(SQL)).Ext(page.ToLimit())
    count, err := session.Count(c)
    if err != nil {
        return count, nil, err
    }
    tasks := []models.TaskSelect{}
    err = session.Select(c, &tasks)
    return count, tasks, err
}

func SelectTasksByUsername(session *Session, username string, page *PageInfo) (int, []models.TaskSelect, error) {
    c := sqlc.NewSQLc(TaskTable)
    SQL := fmt.Sprintf("id IN (SELECT id FROM %s WHERE username = '%s')",
        TaskTable,
        username)
    c.And(sqlc.Exp(SQL)).Ext(page.ToLimit())
    count, err := session.Count(c)
    if err != nil {
        return count, nil, err
    }
    tasks := []models.TaskSelect{}
    err = session.Select(c, &tasks)
    return count, tasks, err
}

func UpdateTaskState(session *Session, id int64, state uint) error {
    task := &models.TaskStatePriorityUpdate{}
    task.ModifyTime = timeStampNow()
    task.State = state
    c := sqlc.NewSQLc(TaskTable)
    c.And(sqlc.Equal("id", id))
    err := session.Update(c, *task)
    return err
}

func UpdateTaskPriority(session *Session, id int64, priority uint) error {
    task := &models.TaskStatePriorityUpdate{}
    task.ModifyTime = timeStampNow()
    task.Priority = priority
    c := sqlc.NewSQLc(TaskTable)
    c.And(sqlc.Equal("id", id))
    err := session.Update(c, *task)
    return err
}

func UpdateTask(session *Session, id int64, task *models.TaskUpdate) error {
    task.ModifyTime = timeStampNow()
    c := sqlc.NewSQLc(TaskTable)
    c.And(sqlc.Equal("id", id))
    err := session.Update(c, *task)
    return err
}

func UpdateTaskUsername(session *Session, oldUsername, newUsername string) error {
    c := sqlc.NewSQLc(TaskTable)
    task := &models.TaskUsernameUpdate{
        Username: newUsername,
    }
    c.And(sqlc.Equal("username", oldUsername))
    err := session.Update(c, *task)
    return err
}

func DeleteTask(session *Session, id int64) error {
    c := sqlc.NewSQLc(TaskTable)
    c.And(sqlc.Equal("id", id))
    err := session.Delete(c)
    if err != nil {
        return err
    }
    return nil
}
