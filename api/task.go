package api

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/nihileon/ticktak/dal"
    "github.com/nihileon/ticktak/log"
    "github.com/nihileon/ticktak/models"
    "strconv"
)

const (
    UninitializedState = iota
    ActiveNotDeleted
    DoneNotDeleted
    ExpiredNotDeleted
    ActiveOrExpiredDeleted
    DoneDeleted
    UpperBoundState
)

const (
    UninitializedPriority = iota
    P1
    P2
    P3
    P4
    UpperBoundPriority
)

type AddTaskReq struct {
    Title    string `json:"title" binding:"required"`
    State    uint   `json:"state"`
    Priority uint   `json:"priority"`
    Content  string `json:"content" binding:"required"`
}
type AddTaskResp struct {
    Id int64 `json:"id"`
}

func AddTask(c *gin.Context) {
    user, ok := AuthReq(c)
    if !ok {
        return
    }
    req := &AddTaskReq{}
    if err := c.BindJSON(req); err != nil {
        doResp(c, nil, err)
        return
    }
    if req.State <= UninitializedState || req.State >= UpperBoundState {
        req.State = ActiveNotDeleted
    }
    if req.Priority <= UninitializedPriority || req.Priority >= UpperBoundPriority {
        req.Priority = P2
    }

    log.GetLogger().Info("Add task, user: %s, title: %v", user, req.Title)
    insertTask := &models.TaskInsert{
        Username: user,
        Title:    req.Title,
        State:    req.State,
        Priority: req.Priority,
        Content:  req.Content,
    }

    id, err := dal.InsertTask(dal.FetchSession(), insertTask)
    resp := &AddTaskResp{
        Id: id,
    }
    doResp(c, resp, err)

}

func GetStatePriorityInfo(c *gin.Context) *models.TaskStatePriorityUpdate {
    info := &models.TaskStatePriorityUpdate{
        State:    UninitializedState,
        Priority: UninitializedPriority,
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
        return
    }
    page := GetPageInfo(c)
    statePriority := GetStatePriorityInfo(c)
    s := statePriority.State
    if s <= UninitializedState || s >= UpperBoundState {
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
        return
    }
    page := GetPageInfo(c)
    statePriority := GetStatePriorityInfo(c)
    p := statePriority.Priority
    if p <= UninitializedPriority || p >= UpperBoundPriority {
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
    err := dal.UpdateTaskState(dal.FetchSession(), req.ID, req.State)
    if err != nil {
        doResp(c, nil, err)
    }
    doResp(c, "update state successfully", nil)
}

func ChangeTaskPriority(c *gin.Context) {
    user, ok := AuthReq(c)
    if !ok {
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
    }
    doResp(c, "update priority successfully", nil)
}
