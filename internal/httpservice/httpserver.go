package httpservice

import (
    "rks-golang-practice/internal/define"
    "fmt"
    "github.com/gin-gonic/gin"
    "rks-golang-practice/internal/httpservice/services"
    _ "rks-golang-practice/internal/httpservice/docs"
    "github.com/swaggo/gin-swagger"
    "github.com/swaggo/gin-swagger/swaggerFiles"
    "rks-golang-practice/internal/httpservice/middleware"
)

type HServer struct {
}


func NewHTTPServer() *HServer {
    s := &HServer{}
    return s
}

// ********* Controller Start 寫在這裡 ********************

func exampleController(v1 *gin.RouterGroup) *gin.RouterGroup {

    c:= services.ExampleController{}
    example := v1.Group("/example")
    {
        example.GET("/testget/:name", c.TestGet)
        example.POST("/testpost", c.TestPOST)
    }
    return v1
}

//每一個Contoller, 都是一個function



// ********* Controller End寫在這裡 ************************

func (s *HServer) Router() *gin.Engine {
    r := gin.Default()
    r.Use(middleware.Cors())
    r.Use(middleware.SetReqKey())
    r.Use(middleware.RespLog())
    r.Use(middleware.PanicHandle())

    v1 := r.Group("/v1")
    {
        exampleController(v1)
        // ControllerFunction(v1)
    }
    return r
}

// @title Swagger httpSwagApiExample API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host
// @BasePath /
func (s *HServer) InitHttpServer() error {
    serverAddr := fmt.Sprintf("%s:%d", define.Cfg.HttpServerIp, define.Cfg.HttpServerPort)

    router := s.Router()
    gin.SetMode(gin.DebugMode)
    router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    return router.Run(serverAddr)
}

