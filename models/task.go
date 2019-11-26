package models

type TaskInsert struct {
    Username   string `json:"username" binding:"required"`
    Title      string `json:"title" binding:"required"`
    State      uint   `json:"state"`
    Priority   uint   `json:"priority"`
    Content    string `json:"content"`
    CreateTime string `json:"create_time"`
    ModifyTime string `json:"modify_time"`
}

type TaskSelect struct {
    Id         int64  `json:"id"`
    Username   string `json:"username"`
    Title      string `json:"title"`
    State      uint   `json:"state"`
    Priority   uint   `json:"priority"`
    Content    string `json:"content"`
    ModifyTime string `json:"modify_time"`
}

type TaskUpdate struct {
    Title      string `json:"title"`
    State      uint   `json:"state"`
    Priority   uint   `json:"priority"`
    Content    string `json:"content"`
    ModifyTime string `json:"modify_time"`
}

type TaskUsernameUpdate struct {
    Username string `json:"username"`
}

type TaskStatePriorityUpdate struct {
    State      uint   `json:"state"`
    Priority   uint   `json:"priority"`
    ModifyTime string `json:"modify_time"`
}

type TaskStatePriorityInfo struct {
    State    uint `json:"state"`
    Priority uint `json:"priority"`
}
