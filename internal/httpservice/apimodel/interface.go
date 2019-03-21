package apimodel

import (
    "gitlab.kingbay-tech.com/engine-lottery/magneto/internal/define"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/pkg/helpers/common"
)

func WSRespFmt(data interface{}, code int32, msg string, uuid string, event string, broadcastType int) string {
    resp := &WSResponse{}
    if len(msg) == 0{
        msg = define.ErrMsgMap[code]
    }
    resp.Message = msg
    resp.Code = code
    resp.Data = data
    resp.UUId = uuid
    resp.Event = event
    resp.BroadcastType = broadcastType
    s , _ :=common.Struct2Json(resp)
    return s
}

type WSResponse struct {
    Code int32 `json:"code"`
    Data interface{} `json:"data"`
    Message string `json:"message"`
    UUId string `json:"UUId"`
    // 顯示在前端的樣式
    Display int64 `json:"display"`
    Event string `json:"event"`
    BroadcastType int `json:"broadcastType"`
}
