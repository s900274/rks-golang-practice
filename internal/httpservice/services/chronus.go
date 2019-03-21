package services

import (
    "gitlab.kingbay-tech.com/engine-lottery/magneto/pkg/token-verification"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/internal/httpservice/apimodel"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/internal/define"
    "fmt"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/pkg/helpers/common"
    "github.com/pkg/errors"
)

func GetUserSpeak(tokenInfo *verification.FrontEndToken) (bool, error){


    reqData := apimodel.GetUserSpeakRequest{}
    reqData.OwnerId = tokenInfo.PlatInfoId
    reqData.LevelId = tokenInfo.LevelId
    canSpeak := true
    reqData.Status = &canSpeak

    marsCli := MarsClient{}
    reqBody, _ := common.Struct2Json(reqData)
    marsReq := marsCli.SetMarsRequest(define.CHRONUS_APINAME, reqBody, define.CHRONUS_GET_USERSPEAK, define.CHRONUS_VERSION)
    trace := marsCli.SetTrace()
    fmt.Println(common.Struct2Json(marsReq))

    resp, err := marsCli.ProcessDispatch(marsReq, trace)

    if err != nil{
        return false, err
    }

    if resp.Code != 0{
        return false, errors.New(resp.Msg)
    }

    if resp.Data == nil{
        return false, errors.New("not found data")
    }

    result := []apimodel.GetUserSpeakResponse{}
    common.Json2Struct(*resp.Data , &result)
    if len(result) == 0{
        return false, nil
    }
    return true, nil
}
