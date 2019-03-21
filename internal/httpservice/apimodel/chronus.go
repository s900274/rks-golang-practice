package apimodel


//取得UserLevel發言狀態
//Action:0x400007
type GetUserSpeakRequest struct {
    OwnerId int64 `json:"ownerId"`
    LevelId int64 `json:"levelId"`
    Status *bool `json:"status"` //是否可發言
}

//取得UserLevel發言狀態
//Action:0x400007
type GetUserSpeakResponse struct {
    Id int64 `json:"id"`
    OwnerId int64 `json:"ownerId"`
    LevelId int64 `json:"levelId"`
    Status bool `json:"status"`
    TriggerTime int64 `json:"triggerTime"`
}