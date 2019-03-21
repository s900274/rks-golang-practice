package verification

import (
    "gitlab.kingbay-tech.com/engine-lottery/magneto/pkg/helpers/redisclient"
    "fmt"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/pkg/redigo/redis"
    logger "github.com/shengkehua/xlog4go"
    "errors"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/pkg/helpers/common"
)

type TokenResolved struct {
    Data interface{} `json:"data"`
    Eor error `json:"eor"`
}


type FrontEndToken struct {
    MemberId int64 `json:"memberId"`
    MemberAccount string `json:"memberAccount"`//会员账号
    PlatInfoId int64 `json:"platInfoId"`//平台商id
    PlatAccount string `json:"platAccount"`//平台账号
    PlatName string `json:"platName"` //平台名稱
    AgentId int64 `json:"agentId"` //代理id
    AgentName string `json:"agentName"` //代理名称
    AgentAccount string `json:"agentAccount"`//代理账号
    LevelId int64 `json:"levelId"`//层级id
    LevelName string `json:"levelName"`//层级名称
}


type MasterToken struct {
    UserId int64 `json:"userId"`
    UserAccount string `json:"userAccount"`
}


type OwnerToken struct {
    UserId int64 `json:"userId"`//用户id
    UserAccount string `json:"userAccount"`//用户账号
    PlatInfoId int64 `json:"platInfoId"`//平台商id
    PlatAccount string `json:"platAccount"`//平台账号
    LoginType string `json:"loginType"`//agnt 代理、plat平台
}

var instance *CallJavaThrift

//singleton
func GetInstance() *CallJavaThrift {
    if instance == nil {
        //只有一開始init, 先不考慮線程問題
        instance = &CallJavaThrift{}
    }
    return instance
}

const (
    FRONT_TOKEN_PREKEY = "magneto-FE-"
    MASTER_TOKEN_PREKEY = "magneto-MS-"
    OWNER_TOKEN_PREKEY= "magneto-OW-"
)

type GetTokenFactoryer interface {
    GetFrontTokenInfo(token string, account string, origin string, fromCache bool)(interface{}, error)
}

func NewGetTokenFactoryer(from string) (s GetTokenFactoryer) {
    switch from {
    case "H5":
        s = GetH5TokenInfoInstance()
    case "PC":
        s = GetPCTokenInfoInstance()
    }
    return s
}

type GetH5TokenInfo struct{
}

type GetPCTokenInfo struct{
}

var h5instance *GetH5TokenInfo
func GetH5TokenInfoInstance() *GetH5TokenInfo {
    if h5instance == nil {
        //只有一開始init, 先不考慮線程問題
        h5instance = &GetH5TokenInfo{}
    }
    return h5instance
}

var pcinstance *GetPCTokenInfo
func GetPCTokenInfoInstance() *GetPCTokenInfo {
    if pcinstance == nil {
        //只有一開始init, 先不考慮線程問題
        pcinstance = &GetPCTokenInfo{}
    }
    return pcinstance
}


func (t *GetH5TokenInfo)GetFrontTokenInfo(token string, account string, origin string, fromCache bool) (interface{}, error) {
    tokenKey := fmt.Sprintf("%s%s", FRONT_TOKEN_PREKEY, token)
    redisData, err := redisclient.RedisCli.Get(tokenKey)
    //如果readis 沒資料 or 不從cache拿
    if err == redis.ErrNil || !fromCache{
        if err == redis.ErrNil{
            logger.Debug(" redis.ErrNil: %s", account)
        }

        // get data from java
        resp, err := GetInstance().CheckFrontToken(token, account, origin)

        switch {
        case err !=nil:
            return nil, err
        case resp.Code != 0:
            return nil, errors.New(fmt.Sprintf("%d, %s", resp.Code, resp.Msg))
        case resp.Data == nil:
            return nil, errors.New("java resp.body is nil")
        }

        frontS := JavaFrontResponse{}
        err = common.Json2Struct(*resp.Data, &frontS)
        if err !=nil{
            return nil, errors.New(fmt.Sprintf("parars toekn fail from java, body: %s", *resp.Data))
        }

        result := &FrontEndToken{
            MemberId: frontS.MemberId,
            MemberAccount: frontS.MemberName,
            PlatInfoId: frontS.PlatInfoId,
            PlatAccount: frontS.PlatAccount,
            PlatName: frontS.PlatName,
            AgentId: frontS.AgentId,
            AgentName:frontS.AgentName,
            AgentAccount:frontS.AgentAccount,
            LevelId:frontS.LevelId,
            LevelName:frontS.LevelName,
        }
        resBody, _ := common.Struct2Json(result)

        err = redisclient.RedisCli.SetEx(tokenKey, resBody, GetInstance().TokenExpireSecond)
        if err != nil{
            logger.Error("Set Redis Error, err: %s", err.Error())
        }
        return result, nil
    }else if(err != nil){
        return nil, err
    }

    fmt.Println("data from redis")
    result :=  &FrontEndToken{}

    err = common.Json2Struct(redisData, result)
    if err != nil{
        return nil, errors.New("Convrt tokenInfo fail from redis")
        logger.Error("Convrt tokenInfo fail from redis, redisData: %s", redisData )
    }
    //readis有資料
    //fmt.Println(redisData)
    //fmt.Println(err)

    return result, err
}


func (t *GetPCTokenInfo)GetFrontTokenInfo(token string, account string, origin string, fromCache bool) (interface{}, error) {

    tokenKey := fmt.Sprintf("%s%s", FRONT_TOKEN_PREKEY, token)
    redisData, err := redisclient.RedisCli.Get(tokenKey)
    //如果readis 沒資料 or 不從cache拿
    if err == redis.ErrNil || !fromCache{
        if err == redis.ErrNil{
            logger.Debug(" redis.ErrNil: %s", account)
        }

        // get data from java
        resp, err := GetInstance().CheckPCFrontToken(token, account, origin)

        switch {
        case err !=nil:
            return nil, err
        case resp.Code != 0:
            return nil, errors.New(fmt.Sprintf("%d, %s", resp.Code, resp.Msg))
        case resp.Data == nil:
            return nil, errors.New("java resp.body is nil")
        }

        frontS := JavaFrontResponse{}
        err = common.Json2Struct(*resp.Data, &frontS)
        if err !=nil{
            return nil, errors.New(fmt.Sprintf("parars toekn fail from java, body: %s", *resp.Data))
        }

        result := &FrontEndToken{
            MemberId: frontS.MemberId,
            MemberAccount: frontS.MemberName,
            PlatInfoId: frontS.PlatInfoId,
            PlatAccount: frontS.PlatAccount,
            PlatName: frontS.PlatName,
            AgentId: frontS.AgentId,
            AgentName:frontS.AgentName,
            AgentAccount:frontS.AgentAccount,
            LevelId:frontS.LevelId,
            LevelName:frontS.LevelName,
        }
        resBody, _ := common.Struct2Json(result)

        err = redisclient.RedisCli.SetEx(tokenKey, resBody, GetInstance().TokenExpireSecond)
        if err != nil{
            logger.Error("Set Redis Error, err: %s", err.Error())
        }
        return result, nil
    }else if(err != nil){
        return nil, err
    }

    fmt.Println("data from redis")
    result :=  &FrontEndToken{}

    err = common.Json2Struct(redisData, result)
    if err != nil{
        return nil, errors.New("Convrt tokenInfo fail from redis")
        logger.Error("Convrt tokenInfo fail from redis, redisData: %s", redisData )
    }
    //readis有資料
    //fmt.Println(redisData)
    //fmt.Println(err)

    return result, err
}



func GetMasterTokenInfo(token string, account string, origin string, fromCache bool) (interface{}, error) {
    tokenKey := fmt.Sprintf("%s%s", MASTER_TOKEN_PREKEY , token)
    redisData, err := redisclient.RedisCli.Get(tokenKey)

    //如果readis 沒資料 or 不從cache拿
    if err == redis.ErrNil || !fromCache{
        if err == redis.ErrNil{
            logger.Debug(" redis.ErrNil: %s", account)
        }
        // get data from java
        resp, err := GetInstance().CheckMsterToken(token, account, origin)

        switch {
        case err !=nil:
            return nil, err
        case resp.Code != 0:
            return nil, errors.New(fmt.Sprintf("%d, %s", resp.Code, resp.Msg))
        case resp.Data == nil:
            return nil, errors.New("java resp.body is nil")
        }

        masterS := JavaMasterResponse{}
        err = common.Json2Struct(*resp.Data, &masterS)
        if err !=nil{
            return nil, errors.New(fmt.Sprintf("parars toekn fail from java, body: %s", *resp.Data))
        }

        result :=  &MasterToken{
            UserId: masterS.UserId,
            UserAccount: masterS.UserName,
        }
        resBody, _ := common.Struct2Json(result)
        err = redisclient.RedisCli.SetEx(tokenKey, resBody, GetInstance().TokenExpireSecond)
        if err != nil{
            logger.Error("Set Redis Error, err: %s", err.Error())
        }
        return result, nil
    }else if(err != nil){
        return nil, err
    }

    fmt.Println("data from redis")
    result :=  &MasterToken{}

    err = common.Json2Struct(redisData, result)
    if err != nil{
        return nil, errors.New("Convrt tokenInfo fail from redis")
        logger.Error("Convrt tokenInfo fail from redis, redisData: %s", redisData )
    }
    //readis有資料
    //fmt.Println(redisData)
    //fmt.Println(err)

    return result, err
}


func GetOwnerTokenInfo(token string, account string, origin string, fromCache bool) (interface{}, error) {
    tokenKey := fmt.Sprintf("%s%s", OWNER_TOKEN_PREKEY , token)
    redisData, err := redisclient.RedisCli.Get(tokenKey)

    //如果readis 沒資料 or 不從cache拿
    if err == redis.ErrNil || !fromCache{
        if err == redis.ErrNil{
            logger.Debug(" redis.ErrNil: %s", account)
        }
        // get data from java
        resp, err := GetInstance().CheckOwnerToken(token, account, origin)

        switch {
        case err !=nil:
            return nil, err
        case resp.Code != 0:
            return nil, errors.New(fmt.Sprintf("%d, %s", resp.Code, resp.Msg))
        case resp.Data == nil:
            return nil, errors.New("java resp.body is nil")
        }
        ownerS := JavaOwnerResponse{}
        err = common.Json2Struct(*resp.Data, &ownerS)
        if err !=nil{
            return nil, errors.New(fmt.Sprintf("parars toekn fail from java, body: %s", *resp.Data))
        }

        result :=  &OwnerToken{
            UserId: ownerS.UserId,
            UserAccount: ownerS.Username,
            PlatInfoId: ownerS.PlatInfoId,
            PlatAccount: ownerS.PlatName,
            LoginType: ownerS.LoginType,
        }
        resBody, _ := common.Struct2Json(result)

        err = redisclient.RedisCli.SetEx(tokenKey, resBody, GetInstance().TokenExpireSecond)
        if err != nil{
            logger.Error("Set Redis Error, err: %s", err.Error())
        }
        return result, nil
    }else if(err != nil){
        return nil, err
    }


    result :=  &OwnerToken{}

    err = common.Json2Struct(redisData, result)
    if err != nil{
        return nil, errors.New("Convrt tokenInfo fail from redis")
        logger.Error("Convrt tokenInfo fail from redis, redisData: %s", redisData )
    }
    //readis有資料
    //fmt.Println(redisData)
    //fmt.Println(err)

    return result, err
}
