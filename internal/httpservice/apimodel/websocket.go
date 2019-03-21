package apimodel

import (
    "fmt"
)

type TokenVerification struct {
    UserToken string `json:"userToken"`
    UserAccount string `json:"userAccount"`
    Origin string `json:"origin"`
    From string `json:"from"` //來源: ex: front(前端) / owner(業主) / master(主控)
}

type Account struct {
    UserAccount string `json:"userAccount"`
}

func (s *Account)SetAccount(account string)  {

    s.UserAccount = fmt.Sprintf("%s***", account[0:len(account)-3])
}

type ChatMessageRequest struct {
    Msg string `json:"msg"`
    UUID string `json:"uuid"`
    Type string `json:"type"`
}

type BetOrderRequest struct {
    Data interface{} `json:"data"`
    UUID string `json:"uuid"`
    Type string `json:"type"`
}


type BetOrderResponse struct {
    Account
    Msg interface{} `json:"msg"`
}

type ChatMessageResponse struct {
    Account
    Msg string `json:"msg"`
}

//func (s *ChatMessageResponse)SetAccount(account string)  {
//    result := account
//    l := len(result)
//    if len(account) <= 6{
//        result = "******"
//    }
//
//    s.UserAccount = fmt.Sprintf("***%s***", result[3:l-3])
//}

//眾籌
type CrowdfundingRequest struct {
    Type string `json:"type"`
    UUID string `json:"uuid"`
    Data interface{} `json:"data"`
}

//眾籌
type StartCrowdfundingRequest struct {
    CrowdfundingNum string `json:"crowdfundingNum"`
}