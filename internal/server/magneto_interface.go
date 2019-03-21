package server

import (
    "gitlab.kingbay-tech.com/engine-lottery/magneto/pkg/downstream"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/internal/define"
    "crypto/md5"
    "encoding/hex"
    "strings"
    "bytes"
    "strconv"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/pkg/bc/base"
    logger "github.com/shengkehua/xlog4go"
    //"github.com/axgle/mahonia"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/pkg/helpers/ucm"
    "errors"
    //"golang.org/x/text/encoding"
    "reflect"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/pkg/helpers/monitor"
    "fmt"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/pkg/bc/engine/mars"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/pkg/helpers/utils"

)

func (s *MarsServer)PublicInterface(req *base.MarsRequest, trace *base.Trace) (r *base.MarsResponse, err error) {

    defer func() {
        err := recover()
        if nil != err {
            stackInfo := utils.GetStackInfo()
            msg := fmt.Sprintf(">>> ```%s``` \n\n `%v`", stackInfo, err)
            utils.CallSlack(msg, define.SLACK_PANIC_CHANNEL, define.SLACK_PANIC_SENDFROM)
        }
    }()
    //TODO 個服務實作 switch

    switch req.Action{
    case 1:

    default:

    }

   return r, nil
}

func (s *MarsServer)processDispatch(req *base.MarsRequest, trace *base.Trace) (r *base.MarsResponse, err error) {
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
    logger.Debug("step in processDispatch trace=[%s:%s] apiname=[%v] action=[%s]", trace.Caller, trace.LogId, req.Apiname, req.Action)
    defer func() {
        if err2 != nil {
            // 统计请求错误数
            reqKey := fmt.Sprintf("Request#%v#%v", req.GetApiname(), req.GetAction())
            monitor.AddError(reqKey)
            logger.Error("step out processDispatch trace=[%s:%s] apiname=[%v] action=[%s] err=[%s]", trace.Caller, trace.LogId, req.Apiname, req.Action, err2.Error())
        } else {
            logger.Debug("step out processDispatch trace=[%s:%s] apiname=[%v] action=[%s]", trace.Caller, trace.LogId, req.Apiname, req.Action)
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


func (s *MarsServer) GetConnStats() map[string] define.ConnStats{
    var apiConnStats = make(map[string] define.ConnStats)

    for apiNmae, connObj := range downstream.DownStreamMgr{
        cliConnCnt := make(map[string]define.PoolStatus)
        connStats := define.ConnStats{}
        connStats.Addrs = make(map[string] define.AddrStats)
        connStats.NoInAddrs = make(map[string] define.PoolStatus)

        logger.Debug("> GetConnStats, %s Start\n", apiNmae)
        //取得 這個ApiName的所有Conn
        for cli, _ := range connObj.MapConn {
            cliAddr, inPool, _, _ := connObj.GetMapConn(cli)
            p :=0
            u :=0
            switch inPool {
            case true:
                p ++
            default:
                u ++
            }
            if poolStruct, ok := cliConnCnt[cliAddr]; ok {
                poolStruct.InPool += p
                poolStruct.OnUse += u
                cliConnCnt[cliAddr] = poolStruct
            }else{
                _ps := define.PoolStatus{
                    InPool: p,
                    OnUse: u,
                }
                cliConnCnt[cliAddr] = _ps
            }
        }

        //整理 所有的Addr資訊
        for addr, _ := range connObj.GetAllAddr(){
            health, errCnt, _, connCnt := connObj.GetAddrAttr(addr)
            addrObj := define.AddrStats{
                Health: health,
                ErrCount: errCnt,
                ConnCount: connCnt,
            }

            if val, ok := cliConnCnt[addr]; ok{
                addrObj.InPoolCount = val.InPool
                addrObj.OnUseCount = val.OnUse
            }else{
                addrObj.InPoolCount = -1
                addrObj.OnUseCount = -1
            }
            connStats.Addrs[addr] = addrObj
        }

        //有在pool中, 但沒有出現在Addr的
        for cliAddr, v := range cliConnCnt {
            //找不存在
            if _, ok := connStats.Addrs[cliAddr]; !ok {
                connStats.NoInAddrs[cliAddr] = v
            }
        }

        //fmt.Println(common.Struct2Json(connStats))
        apiConnStats[apiNmae] = connStats
    }
    return apiConnStats
}

func (s *MarsServer) CheckSignature(req *base.MarsRequest) int32 {
    //暫時將驗證移除
    return define.ERR_OK

    keyStr := "fhA4hUoBcReZ8bJddPKkqCE42sn0PzoX"
    buf := bytes.Buffer{}
    buf.WriteString("version=")
    buf.WriteString(strconv.FormatInt(int64(req.GetHeader().GetVersion()), 10))
    buf.WriteString(":apiname=")
    buf.WriteString(req.Apiname)
    buf.WriteString(":action=")
    buf.WriteString(strconv.FormatInt(req.GetAction() , 10 ))
    buf.WriteString(":timestamp=")
    buf.WriteString(strconv.FormatInt(req.GetHeader().GetTimestamp(), 10))
    buf.WriteString(":body=")
    //body := mahonia.NewDecoder("gbk").ConvertString(string(req.GetBody()))
    body := req.GetBody()
    buf.WriteString(body)

    cypt := md5.New()
    cypt.Write(buf.Bytes())
    cipher_one := cypt.Sum(nil)

    cypt.Reset()
    buf.Reset()
    buf.WriteString(keyStr)
    buf.WriteString(":")
    buf.WriteString(hex.EncodeToString(cipher_one))

    cypt.Write(buf.Bytes())
    cipher_hex := cypt.Sum(nil)
    cipher_text := strings.ToUpper(hex.EncodeToString(cipher_hex))

    if strings.EqualFold(cipher_text, req.GetHeader().GetSignature()) {
        return define.ERR_OK
    }

    return define.ERR_SIGNATURE_FAIL
}

func (s *MarsServer) internalErrorResponse() (r *base.MarsResponse, err error) {

    r = base.NewMarsResponse()
    r.Code = define.ERR_INTERNAL_ERROR
    r.Msg = define.ErrMsgMap[define.ERR_INTERNAL_ERROR]
    err = nil

    return
}

func (s *MarsServer) createResp(code int32, resp *base.MarsResponse) (*base.MarsResponse, error) {

    resp.Code = code
    resp.Msg = define.ErrMsgMap[code]
    if resp.Code != define.ERR_OK {
        logger.Error("errno=[%d] errmsg=[%s]", resp.Code, resp.Msg)
    }
    return resp, nil
}
