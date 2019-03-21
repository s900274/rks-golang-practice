package middleware

import (
    "gitlab.kingbay-tech.com/engine-lottery/magneto/pkg/token-verification"
    "github.com/gin-gonic/gin"
    "reflect"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/pkg/helpers/common"
    "github.com/pkg/errors"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/internal/httpservice/apimodel"
    "strings"
)

type THfn func(string, string, string, bool) (interface{}, error)


var TOKEN_H5_HANDLE_MAP = map[string]THfn{
    "FRONT":  verification.NewGetTokenFactoryer("H5").GetFrontTokenInfo,
    "MASTER": verification.GetMasterTokenInfo,
    "OWNER":  verification.GetOwnerTokenInfo,
}

var TOKEN_PC_HANDLE_MAP = map[string]THfn{
    "FRONT":  verification.NewGetTokenFactoryer("PC").GetFrontTokenInfo,
    "MASTER": verification.GetMasterTokenInfo,
    "OWNER":  verification.GetOwnerTokenInfo,
}

func getTokenVerificationParamsFromQuery(queryValue map[string][]string) (*apimodel.TokenVerification, error) {

    keys := []string{"userToken", "origin", "userAccount", "from"}

    queryKeys := make([]string, 0, len(queryValue))
    for k, _ := range queryValue {
        queryKeys = append(queryKeys, k)
    }

    result := &apimodel.TokenVerification{}
    tv := reflect.ValueOf(result)
    tv = tv.Elem()

    for i := 0; i < tv.NumField(); i++ {
        jskey := tv.Type().Field(i).Tag.Get("json")
        //query string key沒在 keys中, 代表參數錯誤
        if common.StringInSlice(jskey, keys) == false {
            return nil, errors.New("not found key: " + jskey)
        }

        // 有給key, 沒給value
        if len(queryValue[jskey]) == 0 {
            return nil, errors.New(" key: " + jskey + " is empty")
        }

        //set value
        f := tv.FieldByName(tv.Type().Field(i).Name)
        f.SetString(queryValue[jskey][0])
    }

    return result, nil
}

func CheckToken() gin.HandlerFunc {
    return func(c *gin.Context) {
        tv, err := getTokenVerificationParamsFromQuery(c.Request.URL.Query())
        result := verification.TokenResolved{}
        userKey := "UserInfo"
        if err != nil {
            result.Eor = err
            c.Set(userKey, result)
            return
        }

        from := tv.From
        account := tv.UserAccount
        token := tv.UserToken
        origin := tv.Origin
        var (
            u   interface{}
        )
        if strings.ToLower(c.Query("source")) == "pc" {
            u, err = TOKEN_PC_HANDLE_MAP[from](token, account, origin, true)
        } else { //default h5
            u, err = TOKEN_H5_HANDLE_MAP[from](token, account, origin, true)
        }
        if err != nil {
            result.Eor = err
            c.Set(userKey, result)
            return
        }

        switch from {
        case "FRONT":
            result.Data = *u.(*verification.FrontEndToken)
        case "MASTER":
            result.Data = *u.(*verification.MasterToken)
        case "OWNER":
            result.Data = *u.(*verification.OwnerToken)
        }

        c.Set("UserInfo", result)
        c.Next()
        return
    }
}
