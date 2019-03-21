package services

import (
    "github.com/gin-gonic/gin"
    "rks-golang-practice/internal/httpservice/apimodel"
    "fmt"
    "github.com/gin-gonic/gin/binding"
    "net/http"
    "rks-golang-practice/internal/define"
)

type ExampleController struct {

}


// @Summary 用GET method 測試api
// @Description 用GET method 測試api
// @Accept  json
// @Produce  json
// @Tags EX1
//          參數名   方式 資料型態 必填    說明
// @Param Authorization header string true "token"
// @Param   name path string true   "名字1" default(Jim)
// @Param   old query int false "年紀" default(18)
// @Header token content-type string true "token"
// @Success 200 {string} apimodel.Response	"ok"
// @Failure 400 {object} apimodel.Response "Bad Request"
// @Router /v1/example/testget/{name} [get]
func (c *ExampleController)TestGet(ctx *gin.Context){

    var req  apimodel.ExampleGet
    if ctx.ShouldBindQuery(&req) != nil {
        ctx.JSON(http.StatusBadRequest,
            apimodel.ApiRespFmt(nil, define.ERR_PARAMS_ERROR, ""),
        )
        return
    }

    ctx.JSON(http.StatusOK, apimodel.ApiRespFmt(
        fmt.Sprintf("%s | %s", ctx.Param("name"), ctx.Query("old")),
        define.ERR_OK,
        "",
    ),)
    return
}

// @Summary 用POST method 測試api
// @Description 用POST method 測試api
// @Accept  json
// @Produce  json
// @Tags EX2
// @Content-Type application/json
// @Param   sendBody body apimodel.ExamplePost true  "body"
// @Success 200 {string} apimodel.Response	"ok"
// @Failure 400 {object} apimodel.Response "Bad Request"
// @Router /v1/example/testpost [post]
func (c *ExampleController)TestPOST(ctx *gin.Context){

    var req apimodel.ExamplePost
    if err := ctx.ShouldBindWith(&req, binding.JSON); err != nil {
        ctx.JSON(http.StatusBadRequest, apimodel.RespParamsError(err))
        return
    }

    ctx.JSON(http.StatusOK, apimodel.ApiRespFmt(
        req,
        define.ERR_OK,
        "",
    ),)
    return
}
