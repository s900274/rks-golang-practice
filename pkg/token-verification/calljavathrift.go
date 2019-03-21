package verification

import (
    "gitlab.kingbay-tech.com/engine-lottery/magneto/pkg/thrift.git/lib/go/thrift"
    logger "github.com/shengkehua/xlog4go"
    "net"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/pkg/bc/engine/mars"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/pkg/bc/base"
    "fmt"
    "time"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/pkg/helpers/common"
    "strings"
    "crypto/md5"
)

type CallJavaThrift struct {
    Host string
    Port string
    ApiName string
    Key string
    TokenExpireSecond int


}

var JAVA_TOKEN_ACTION_MAP = map[string] int64 {
    "FRONT": 200004,
    "FRONT_PC": 100001,
    "OWNER": 200002,
    "MASTER": 200003,
}

var JAVA_TOKEN_APINAME_MAP = map[string] string {
    "FRONT": "uaa",
    "FRONT_PC": "token",
    "OWNER": "uaa",
    "MASTER": "uaa",
}

func (s *CallJavaThrift) Init(host string, port string, key string, tokenExpire int) {
    s.Host = host
    s.Port = port
    s.Key = key
    //fhA4hUoBcReZ8bJddPKkqCE42sn0PzoX
    s.TokenExpireSecond = tokenExpire
    s.ApiName = "uaa"
}


// @Param Authorization header string true "token"
// @Param UserAccount header string true "Account"
// @Param DF_Origin header string flase "Origin"
type FrontRequest struct {
    Authorization string `json:"authorization"`
    MemberName string `json:"memberName"`
    Origin string `json:"origin"`
}


type JavaFrontResponse struct {
    MemberId int64 `json:"memberId"` //会员id
    MemberName string `json:"memberName"`//会员账号
    PlatInfoId int64 `json:"platInfoId"`//平台商id
    PlatAccount string `json:"platAccount"` //平台账号
    PlatName string `json:"platName"` //平台名稱
    AgentId int64 `json:"agentId"` //代理id
    AgentAccount string `json:"agentAccount"` //代理账号
    AgentName string `json:"agentName"` //代理名称
    LevelId int64 `json:"levelId"`//层级id
    LevelName string `json:"levelName"`//层级名称
}

type MasterAndOwnerRequest struct {
    Authorization string `json:"authorization"`
    UserName string `json:"username"`
    Origin string `json:"origin"`
}

type JavaMasterResponse struct {
    UserId int64 `json:"userId"`
    UserName string `json:"userName"`
}

type JavaOwnerResponse struct {
    UserId int64 `json:"userId"`//用户id
    Username string `json:"userName"`//用户账号
    PlatInfoId int64 `json:"platInfoId"`//平台商id
    PlatName string `json:"platName"`//平台账号
    LoginType string `json:"loginType"`//agnt 代理、plat平台
}

func Md5Encode(data string) string {
    Digest := md5.Sum([]byte(data))
    return fmt.Sprintf("%x", Digest)
}

func (s *CallJavaThrift) genHeader() *base.MarsHeader{
    ts := time.Now().Add(-30 * time.Second).UnixNano() / int64(time.Millisecond)
    headers := &base.MarsHeader{
        Version:   1,
        Timestamp: ts,
    }
    return headers
}

func (s *CallJavaThrift) assignSign(thriftReq *base.MarsRequest) (*base.MarsRequest) {
    signString := fmt.Sprintf("key=%s:version=%d:apiname=%s:action=%d:timestamp=%d:body=%s",
       s.Key,
       thriftReq.Header.GetVersion(),
       thriftReq.GetApiname(),
       thriftReq.GetAction(),
       thriftReq.Header.GetTimestamp(),
       thriftReq.GetBody())

    thriftReq.Header.Signature = strings.ToUpper(Md5Encode(signString))
    return thriftReq
}

func (s *CallJavaThrift) CheckFrontToken(token string, account string, origin string) (r *base.MarsResponse, err error) {

    var req FrontRequest
    req.Authorization = token
    req.Origin = origin
    req.MemberName = account
    body, _ := common.Struct2Json(req)
    logger.Debug("Send java check token: %s", body)
    thriftReq := &base.MarsRequest{
        Apiname: s.ApiName,
        Action:  JAVA_TOKEN_ACTION_MAP["FRONT"],
        Header:  s.genHeader(),
        Body:    body,
    }
    thriftReq = s.assignSign(thriftReq)
    return s.SendJavaGWThrift(thriftReq)
}

func (s *CallJavaThrift) CheckPCFrontToken(token string, account string, origin string) (r *base.MarsResponse, err error) {

    var req FrontRequest
    req.Authorization = token
    req.Origin = origin
    req.MemberName = account
    body, _ := common.Struct2Json(req)
    logger.Debug("Send java h5 check token: %s", body)
    thriftReq := &base.MarsRequest{
        Apiname: JAVA_TOKEN_APINAME_MAP["FRONT_PC"],
        Action:  JAVA_TOKEN_ACTION_MAP["FRONT_PC"],
        Header:  s.genHeader(),
        Body:    body,
    }
    thriftReq = s.assignSign(thriftReq)
    return s.SendJavaGWThrift(thriftReq)
}

func (s *CallJavaThrift) CheckMsterToken(token string, account string, origin string) (r *base.MarsResponse, err error) {

    var req MasterAndOwnerRequest
    req.Authorization = token
    req.Origin = origin
    req.UserName = account
    body, _ := common.Struct2Json(req)

    thriftReq := &base.MarsRequest{
        Apiname: s.ApiName,
        Action:  JAVA_TOKEN_ACTION_MAP["MASTER"],
        Header:  s.genHeader(),
        Body:    body,
    }
    thriftReq = s.assignSign(thriftReq)
    return s.SendJavaGWThrift(thriftReq)
}

func (s *CallJavaThrift) CheckOwnerToken(token string, account string, origin string) (r *base.MarsResponse, err error) {

    var req MasterAndOwnerRequest
    req.Authorization = token
    req.Origin = origin
    req.UserName = account
    body, _ := common.Struct2Json(req)

    thriftReq := &base.MarsRequest{
        Apiname: s.ApiName,
        Action:  JAVA_TOKEN_ACTION_MAP["OWNER"],
        Header:  s.genHeader(),
        Body:    body,
    }
    thriftReq = s.assignSign(thriftReq)
    return s.SendJavaGWThrift(thriftReq)
}


func (s *CallJavaThrift)SendJavaGWThrift(req *base.MarsRequest) (r *base.MarsResponse, err error) {
    tSocket, err := thrift.NewTSocket(net.JoinHostPort(s.Host, s.Port))
    if err != nil {
        logger.Error("tSocket java thrift error: %s", err.Error())
    }

    transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
    transport, err := transportFactory.GetTransport(tSocket)
    protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

    client := mars.NewMarsServiceClientFactory(transport, protocolFactory)

    if err := transport.Open(); err != nil {
        logger.Error("Error opening java thrift: %s:%s", s.Host, s.Port)
    }
    defer transport.Close()

    trace := &base.Trace{
        Caller: "toulouse",
        LogId:  fmt.Sprintf("%d", time.Now().UnixNano()/1000000),
    }
    return client.PublicInterface(req, trace)
}
