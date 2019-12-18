package api

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/nihileon/ticktak/dal"
    "github.com/nihileon/ticktak/log"
    "github.com/nihileon/ticktak/models"
    "strconv"
)

type AddTaskReq struct {
    Title    string `json:"title" binding:"required"`
    State    uint   `json:"state"`
    Priority uint   `json:"priority"`
    Content  string `json:"content" binding:"required"`
    Tag      string `json:"tag"`
    DDLTime  string `json:"ddl_time"`
}
type AddTaskResp struct {
    Id int64 `json:"id"`
}

func AddTask(c *gin.Context) {
    user, ok := AuthReq(c)
    if !ok {
        doResp(c, "You don't have permissions.", nil)
        return
    }
    req := &AddTaskReq{}
    if err := c.BindJSON(req); err != nil {
        doResp(c, nil, err)
        return
    }
    if req.State <= models.UninitializedState || req.State >= models.UpperBoundState {
        req.State = models.ActiveNotDeleted
    }
    if req.Priority <= models.UninitializedPriority || req.Priority >= models.UpperBoundPriority {
        req.Priority = models.P2
    }

    log.GetLogger().Info("Add task, user: %s, title: %v", user, req.Title)
    insertTask := &models.TaskInsert{
        Username: user,
        Title:    req.Title,
        State:    req.State,
        Priority: req.Priority,
        Content:  req.Content,
        DDLTime:  req.DDLTime,
        Tag:      req.Tag,
    }

    session := dal.FetchSession()
    if err := session.Begin(); err != nil {
        doResp(c, nil, err)
        return
    }
    defer session.Rollback()

    id, err := dal.InsertTask(session, insertTask)
    if err != nil {
        doResp(c, nil, err)
        return
    }
    if req.State == models.DoneNotDeleted || req.State == models.DoneDeleted {
        err = dal.UpdateDoneTime(session, id)
        if err != nil {
            doResp(c, nil, err)
            return
        }
    }
    resp := &AddTaskResp{
        Id: id,
    }
    err = session.Commit()
    doResp(c, resp, err)
}

func getStatePriorityInfo(c *gin.Context) *models.TaskStatePriorityUpdate {
    info := &models.TaskStatePriorityUpdate{
        State:    models.UninitializedState,
        Priority: models.UninitializedPriority,
    }
    params := c.Request.URL.Query()

    stateV := params["state"]
    if len(stateV) > 0 {
        v, err := strconv.Atoi(stateV[0])
        if err == nil && v > 0 {
            info.State = uint(v)
        }
    }
    priorityV := params["priority"]
    if len(priorityV) > 0 {
        v, err := strconv.Atoi(priorityV[0])
        if err == nil && v > 0 {
            info.Priority = uint(v)
        }
    }
    return info
}

func GetTasksByUsername(c *gin.Context) {
    user, ok := AuthReq(c)
    if !ok {
        doResp(c, "You don't have permissions.", nil)
        return
    }
    page := GetPageInfo(c)
    log.GetLogger().Info("Get Task by username, user: %s, page: %v", user, page)
    count, tasks, err := dal.SelectTasksByUsername(dal.FetchSession(), user, page)
    if err != nil {
        doRespWithCount(c, count, nil, err)
        return
    }
    doRespWithCount(c, count, tasks, nil)
}

func GetTasksByUsernameState(c *gin.Context) {
    user, ok := AuthReq(c)
    if !ok {
        doResp(c, "You don't have permissions.", nil)
        return
    }
    page := GetPageInfo(c)
    statePriority := getStatePriorityInfo(c)
    s := statePriority.State
    if s <= models.UninitializedState || s >= models.UpperBoundState {
        doRespWithCount(c, 0, nil, fmt.Errorf("without a correct state"))
        return
    }
    log.GetLogger().Info("Get Task by username, user: %s, page: %v", user, page)
    count, tasks, err := dal.SelectTasksByUsernameState(dal.FetchSession(), user, s, page)
    if err != nil {
        doRespWithCount(c, count, nil, err)
        return
    }
    doRespWithCount(c, count, tasks, nil)
}

func GetTasksByUsernamePriority(c *gin.Context) {
    user, ok := AuthReq(c)
    if !ok {
        doResp(c, "You don't have permissions.", nil)
        return
    }
    page := GetPageInfo(c)
    statePriority := getStatePriorityInfo(c)
    p := statePriority.Priority
    if p <= models.UninitializedPriority || p >= models.UpperBoundPriority {
        doRespWithCount(c, 0, nil, fmt.Errorf("without a correct priority"))
        return
    }
    log.GetLogger().Info("Get Task by username, user: %s, page: %v", user, page)
    count, tasks, err := dal.SelectTasksByUsernamePriority(dal.FetchSession(), user, p, page)
    if err != nil {
        doRespWithCount(c, count, nil, err)
        return
    }
    doRespWithCount(c, count, tasks, nil)
}

func ChangeTaskState(c *gin.Context) {
    user, ok := AuthReq(c)
    if !ok {
        doResp(c, "You don't have permissions.", nil)
        return
    }
    type updateState struct {
        ID    int64 `json:"id" binding:"required"`
        State uint  `json:"state" binding:"required"`
    }
    req := &updateState{}
    if err := c.BindJSON(req); err != nil {
        doResp(c, nil, err)
        return
    }
    log.GetLogger().Infof("user:%s, update task %u, state %u", user, req.ID, req.State)

    session := dal.FetchSession()
    if err := session.Begin(); err != nil {
        doResp(c, nil, err)
        return
    }
    defer session.Rollback()

    task, err := dal.SelectTasksByTaskID(session, req.ID)
    if err != nil {
        doResp(c, nil, err)
        return
    }
    if task.State != models.DoneNotDeleted && task.State != models.DoneDeleted {
        if req.State == models.DoneDeleted || req.State == models.DoneNotDeleted {
            err = dal.UpdateDoneTime(session, req.ID)
            if err != nil {
                doResp(c, nil, err)
                return
            }
        }
    }

    err = dal.UpdateTaskState(session, req.ID, req.State)
    if err != nil {
        doResp(c, nil, err)
        return
    }

    err = session.Commit()
    doResp(c, "update state successfully", nil)
}

func ChangeTaskPriority(c *gin.Context) {
    user, ok := AuthReq(c)
    if !ok {
        doResp(c, "You don't have permissions.", nil)
        return
    }
    type updatePriority struct {
        ID       int64 `json:"id" binding:"required"`
        Priority uint  `json:"priority" binding:"required"`
    }
    req := &updatePriority{}
    if err := c.BindJSON(req); err != nil {
        doResp(c, nil, err)
        return
    }
    log.GetLogger().Infof("user:%s, update task %u, priority %u", user, req.ID, req.Priority)
    err := dal.UpdateTaskPriority(dal.FetchSession(), req.ID, req.Priority)
    if err != nil {
        doResp(c, nil, err)
        return
    }
    doResp(c, "update priority successfully", nil)
}

func TaskModify(c *gin.Context) {
    user, ok := AuthReq(c)
    if !ok {
        doResp(c, "You don't have permissions.", nil)
        return
    }
    req := &models.TaskUpdate{}
    if err := c.BindJSON(req); err != nil {
        doResp(c, nil, err)
        return
    }

    log.GetLogger().Infof("user:%s, update task %u, state %u", user, req.ID, req.State)

    session := dal.FetchSession()
    if err := session.Begin(); err != nil {
        doResp(c, nil, err)
        return
    }
    defer session.Rollback()

    task, err := dal.SelectTasksByTaskID(session, req.ID)
    if err != nil {
        doResp(c, nil, err)
        return
    }
    if task.State != models.DoneNotDeleted && task.State != models.DoneDeleted {
        if req.State == models.DoneDeleted || req.State == models.DoneNotDeleted {
            err = dal.UpdateDoneTime(session, req.ID)
            if err != nil {
                doResp(c, nil, err)
                return
            }
        }
    }

    err = dal.UpdateTask(session, req.ID, req)
    if err != nil {
        doResp(c, nil, err)
        return
    }

    err = session.Commit()
    doResp(c, "task modify successfully", nil)
}

func GetTaskTagsByUsername(c *gin.Context) {
    user, ok := AuthReq(c)
    if !ok {
        doResp(c, "You don't have permissions.", nil)
        return
    }
    page := GetPageInfo(c)
    log.GetLogger().Info("Get Task tags by username, user: %s, page: %v", user, page)
    count, tags, err := dal.SelectTaskTagsByUsername(dal.FetchSession(), user, page)
    if err != nil {
        doRespWithCount(c, count, nil, err)
        return
    }
    doRespWithCount(c, count, tags, nil)
}

func UpdateTaskStateIfExpired(c *gin.Context) {
    user, ok := AuthReq(c)
    if !ok {
        doResp(c, "You don't have permissions.", nil)
        return
    }
    log.GetLogger().Info("Update task state if expired by username, user: %s", user)
    err := dal.UpdateTaskStateIfExpired(dal.FetchSession(), user)
    if err != nil && "Exec affect 0 row" == err.Error() {
        doResp(c, "update state if expired", nil)
        return
    }
    if err != nil {
        doResp(c, "update state if expired", err)
        return
    }
    doResp(c, "update state if expired", nil)
}
