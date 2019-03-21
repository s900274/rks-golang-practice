package services

import (
    "gitlab.kingbay-tech.com/engine-lottery/magneto/internal/httpservice/apimodel"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/pkg/helpers/common"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/internal/define"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/pkg/token-verification"
    "github.com/pkg/errors"
    "fmt"
)

func GetCrowdfundingDetail(tokenInfo *verification.FrontEndToken, msgObj apimodel.CrowdfundingRequest) (*apimodel.CrowdFundingDetailResp, error ){
    s, _ := common.Struct2Json(msgObj.Data)
    reqData := apimodel.StartCrowdfundingRequest{}
    common.Json2Struct(s, &reqData)


    marsCli := MarsClient{}
    sendObj := apimodel.CrowdFundingDetailReq{}
    sendObj.UserId = tokenInfo.MemberId
    sendObj.CrowdfundingNum = reqData.CrowdfundingNum
    reqBody, _ := common.Struct2Json(sendObj)
    marsReq := marsCli.SetMarsRequest(define.SELENE_APINAME, reqBody, define.SELENE_CROWDFUNDING_DETAIL, define.SELENE_VERSION)
    trace := marsCli.SetTrace()
    fmt.Println(common.Struct2Json(marsReq))

    resp, err := marsCli.ProcessDispatch(marsReq, trace)

    if err != nil{
        return nil, err
    }

    if resp.Code != 0{
        return nil, errors.New(resp.Msg)
    }

    result := &apimodel.CrowdFundingDetailResp{}
    common.Json2Struct(*resp.Data , result)
    return result, nil

}


