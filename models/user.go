package models

type UserInsert struct {
    Username    string `json:"username" binding:"required"`
    Password    string `json:"password" binding:"required"`
    Description string `json:"description"`
    CreateTime  string `json:"create_time"`
    ModifyTime  string `json:"modify_time"`
}

type LoginUpdate struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type LoginSelect struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type UserSelect struct {
    Username    string `json:"username"`
    Description string `json:"description"`
    CreateTime  string `json:"create_time"`
    ModifyTime  string `json:"modify_time"`
}
