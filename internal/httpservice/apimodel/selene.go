package apimodel



//已發起眾籌詳情
//Action: 0x20000B(2097163)
type CrowdFundingDetailReq struct {
    UserId          int64  `json:"userId"`          //用戶編號
    CrowdfundingNum string `json:"crowdfundingNum"` //眾籌編號 varchar(20)
}


//已發起眾籌詳情
type CrowdFundingDetailResp struct {
    CrowdFundingNum  string  `json:"crowdFundingNum"`  //眾籌編號
    LotteryId        int64   `json:"lotteryId"`        //彩種id
    LotteryName      string  `json:"lotteryName"`      //彩種名稱
    PCode            string  `json:"pCode"`            //期號
    FounderId        int64   `json:"founderId"`        //發起人id
    FounderAccount   string  `json:"founderAccount"`   //發起人帳號
    ProfitPercentage float64 `json:"profitPercentage"` //盈利统计 百份比
    RemainCopies     int64   `json:"remainCopies"`     //剩餘份數
    TotalJackpot     int64   `json:"totalJackpot"`     //累積獎金統計
    Ranking          int64   `json:"ranking"`          //眾籌排名
    Isfounder        bool    `json:"isfounder"`        //是否為發起人
    IssueAlias       string  `json:"issueAlias"`       //说明显示
    //--------------Detail-------------------------------//
    CopyAmount         int64 `json:"copyAmount"`         //每份多少錢(分)
    TotalCopies        int64 `json:"totalCopies"`        //總份數
    CopyCommission     int64 `json:"copyCommission"`     //佣金/份(分)
    SelfpurchaseCopies int64 `json:"selfpurchaseCopies"` //發起者的自購幾份
    AddCopies          int64 `json:"addCopies"`          //已參於份數
    StartTime          int64 `json:"startTime"`          //發起時間(ms)
    EndTime            int64 `json:"endTime"`            //截止時間(ms)
}
