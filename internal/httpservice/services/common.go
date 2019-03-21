package services

import (
    "gitlab.kingbay-tech.com/engine-lottery/magneto/pkg/bc/base"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/pkg/helpers/common"
    "github.com/satori/go.uuid"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/internal/define"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/pkg/bc/engine/mars"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/pkg/helpers/utils"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/pkg/helpers/ucm"
    "fmt"
    "errors"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/pkg/helpers/monitor"
    "reflect"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/pkg/downstream"
    logger "github.com/shengkehua/xlog4go"

)


type MarsClient struct {
}

func (s *MarsClient) SetMarsRequest(apiName, body string, action int64, version int8) *base.MarsRequest {
    header := &base.MarsHeader{Version: version, Timestamp: common.GetTimeNowInNs()}
    return &base.MarsRequest{Apiname: apiName, Action: action, Header: header, Body: body}
}
func (s *MarsClient) SetMarsRequestForJava(apiName, body string, action int64, version int8) *base.MarsRequest {
    header := &base.MarsHeader{Version: version, Timestamp: common.GetTimeNowInMs() - 30000}
    return &base.MarsRequest{Apiname: apiName, Action: action, Header: header, Body: body}
}


func (s *MarsClient) SetTrace() *base.Trace {
    u1 := uuid.Must(uuid.NewV4())
    return &base.Trace{LogId: u1.String(), Caller: define.CALLER_NAME}
}


func (s *MarsClient)ProcessDispatch(req *base.MarsRequest, trace *base.Trace) (r *base.MarsResponse, err error) {
    //TODO 調用網關 範例
    defer func() {
        err := recover()
        if nil != err {
            stackInfo := utils.GetStackInfo()
            utils.CallSlack(stackInfo, define.SLACK_PANIC_CHANNEL, define.SLACK_PANIC_SENDFROM)
        }
    }()

    var (
        pool *ucm.ChanConnPool
        cli interface{}
        err2 error
        ok bool
        health bool = true
    )

    // debug
    logger.Debug("step in processDispatch trace=[%s:%s] apiname=[%v] action=[%d]", trace.Caller, trace.LogId, req.Apiname, req.Action)
    defer func() {
        if err2 != nil {
            // 统计请求错误数
            reqKey := fmt.Sprintf("Request#%v#%v", req.GetApiname(), req.GetAction())
            monitor.AddError(reqKey)
            logger.Error("step out processDispatch trace=[%s:%s] apiname=[%v] action=[%d] err=[%s]", trace.Caller, trace.LogId, req.Apiname, req.Action, err2.Error())
        } else {
            logger.Debug("step out processDispatch trace=[%s:%s] apiname=[%v] action=[%d]", trace.Caller, trace.LogId, req.Apiname, req.Action)
        }
    }()

    // 释放连接资源
    defer func() {
        if (pool != nil) && (cli != nil) {
            pool.Put(cli, health)
            if !health {
                logger.Warn("cli type=[%v]", reflect.TypeOf(cli))
            }
        }
    }()


    var downType string

    //取得網關通道
    pool, ok = downstream.DownStreamMgr[define.GW_APINAME]
    if !ok {
        logger.Error("downType error downType=[%s]", downType)
        err2 = errors.New(define.ErrMsgMap[define.ERR_DOWNSTREAM_ERROR])
        r = base.NewMarsResponse()
        return s.createResp(define.ERR_DOWNSTREAM_ERROR, r)
    }
    if pool == nil {
        err2 = errors.New(define.ErrMsgMap[define.ERR_CONN_POOL_NIL])
        r =  base.NewMarsResponse()
        return s.createResp(define.ERR_CONN_POOL_NIL, r)
    }

    cli, err2 = pool.Get()
    if err2 != nil {
        logger.Error("Get connection fail. detail=[%s] downType=[%s]", err2.Error(), downType)
        return
    }

    //嘗試發送請求
    for i := 0; i <= 1; i++ {
        logger.Debug("trace=[%s:%s] before call downstream. apiname=[%v] action=[%s]", trace.Caller, trace.LogId, req.Apiname, req.Action)
        client := cli.(*mars.MarsServiceClient)
        if r, err2 = client.PublicInterface(req, trace); err2 != nil {
            health = false
        }

        logger.Debug("trace=[%s:%s] after call downstream. apiname=[%v] action=[%s]", trace.Caller, trace.LogId, req.Apiname, req.Action)

        if health {
            break
        } else {
            logger.Debug("trace=[%s:%s] ReConn apiname=[%v] action=[%X]", trace.Caller, trace.LogId, req.Apiname, req.Action)
            cli , err = pool.ReConn(cli)
            if err != nil {
                logger.Error("trace=[%s:%s] ReConn failed apiname=[%v] action=[%s] err=[%v]", trace.Caller, trace.LogId, req.Apiname, req.Action, err)
                break
            }
            health = true
        }
    }
    return
}

func (s *MarsClient) createResp(code int32, resp *base.MarsResponse) (*base.MarsResponse, error) {

    resp.Code = code
    resp.Msg = define.ErrMsgMap[code]
    if resp.Code != define.ERR_OK {
        logger.Error("errno=[%d] errmsg=[%s]", resp.Code, resp.Msg)
    }
    return resp, nil
}