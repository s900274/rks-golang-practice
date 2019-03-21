package define

import "errors"

var (
    ErrMsgMap = make(map[int32]string)
    ErrMsgNil = errors.New("Message is nil")
)

const (
    ERR_OK                         = 0
    ERR_MOD_NOT_EXIST              = 4001
    ERR_SIGNATURE_FAIL             = 4002
    ERR_REQUEST_NO_TOKEN           = 4003
    ERR_DOWNSTREAM_ERROR           = 4004
    ERR_CONN_POOL_NIL              = 4005
    ERR_CONN_GET_FAIL              = 4006
    ERR_INTERNAL_ERROR             = 9999
    ERR_PARAMS_ERROR               = -10
    ERR_CHECKTOKEN_ERROR           = -11
    ERR_GETCROWDFUNDINGDETAIL_FAIL = -12
    ERR_ERROR                      = -9999
)

func init() {
    ErrMsgMap[ERR_OK] = "Success"
    ErrMsgMap[ERR_MOD_NOT_EXIST] = "Moudle is not exist"
    ErrMsgMap[ERR_SIGNATURE_FAIL] = "Signature fail"
    ErrMsgMap[ERR_REQUEST_NO_TOKEN] = "Request not have token"
    ErrMsgMap[ERR_DOWNSTREAM_ERROR] = "Downstream is not found"
    ErrMsgMap[ERR_CONN_POOL_NIL] = "Connect pool is nil"
    ErrMsgMap[ERR_CONN_GET_FAIL] = "Get a connection from pool failed"
    ErrMsgMap[ERR_INTERNAL_ERROR] = "Interval error"
    ErrMsgMap[ERR_PARAMS_ERROR] = "Please check params"
    ErrMsgMap[ERR_CHECKTOKEN_ERROR] = "Check token fail"
    ErrMsgMap[ERR_GETCROWDFUNDINGDETAIL_FAIL] = "Get crowdgunding fail"

    ErrMsgMap[ERR_ERROR] = "No handle error"
}
