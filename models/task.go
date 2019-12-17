package models

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

type TaskInsert struct {
    Username   string `json:"username" binding:"required"`
    Title      string `json:"title" binding:"required"`
    State      uint   `json:"state"`
    Priority   uint   `json:"priority"`
    Content    string `json:"content"`
    CreateTime string `json:"create_time"`
    ModifyTime string `json:"modify_time"`
    Tag        string `json:"tag"`
    DoneTime   string `json:"done_time"`
    DDLTime    string `json:"ddl_time"`
}

type TaskSelect struct {
    Id         int64  `json:"id"`
    Username   string `json:"username"`
    Title      string `json:"title"`
    State      uint   `json:"state"`
    Priority   uint   `json:"priority"`
    Content    string `json:"content"`
    ModifyTime string `json:"modify_time"`
    Tag        string `json:"tag"`
    DoneTime   string `json:"done_time"`
    DDLTime    string `json:"ddl_time"`
}

type TaskUpdate struct {
    ID         int64  `json:"id" binding:"required"`
    Title      string `json:"title"`
    State      uint   `json:"state"`
    Priority   uint   `json:"priority"`
    Content    string `json:"content"`
    ModifyTime string `json:"modify_time"`
    Tag        string `json:"tag"`
    DDLTime    string `json:"ddl_time"`
}

type TaskTagSelect struct {
    Tag string `json:"tag"`
}

type TaskUsernameUpdate struct {
    Username string `json:"username"`
}

type TaskStatePriorityUpdate struct {
    State      uint   `json:"state"`
    Priority   uint   `json:"priority"`
    ModifyTime string `json:"modify_time"`
}

type TaskDoneTimeUpdate struct {
    DoneTime string `json:"done_time"`
}

type TaskStatePriorityInfo struct {
    State    uint `json:"state"`
    Priority uint `json:"priority"`
}

type TaskTagInfo struct {
    Tag string `json:"tag"`
}
