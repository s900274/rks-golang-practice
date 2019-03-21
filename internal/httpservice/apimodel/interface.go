package apimodel

import (
    "fmt"
    "rks-golang-practice/internal/define"
)

func ApiRespFmt(data interface{}, code int32, msg string) *Response {
    resp := &Response{}
    if len(msg) == 0{
        msg = define.ErrMsgMap[code]
    }
    resp.Message = msg
    resp.Code = code
    resp.Data = data
    return resp
}

type Response struct {
    Code int32
    Message string
    Data interface{}
}

func RespParamsError(err error) *Response {
    msg := fmt.Sprintf("%s|%s", define.ErrMsgMap[define.ERR_PARAMS_ERROR], err.Error())
    return ApiRespFmt(
        nil,
        define.ERR_PARAMS_ERROR,
        msg,
    )
}