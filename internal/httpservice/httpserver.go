package httpservice

import (
    "encoding/json"
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/pkg/errors"
    logger "github.com/shengkehua/xlog4go"
    "github.com/swaggo/gin-swagger"
    "github.com/swaggo/gin-swagger/swaggerFiles"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/internal/define"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/internal/httpservice/apimodel"
    _ "gitlab.kingbay-tech.com/engine-lottery/magneto/internal/httpservice/docs"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/internal/httpservice/middleware"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/internal/httpservice/services"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/pkg/googollee/go-socket.io"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/pkg/helpers/common"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/pkg/helpers/kafkaproducer"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/pkg/helpers/utils"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/pkg/token-verification"
    "net/http"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "github.com/importcjj/sensitive"
)

type HServer struct {
    RunDirPath string
}

//type HandlerFunc func(w http.ResponseWriter, r *http.Request)

func NewHTTPServer() *HServer {
    s := &HServer{}

    return s
}

var SocketioServer, _ = socketio.NewServer(nil)

var MsgFilter = sensitive.New()
// @title Swagger test_swag API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host
// @BasePath /
func (s *HServer) InitHttpServer() error {

    defer func() {
        err := recover()
        if err != nil {
            logger.Error("magneto panic err: %s", err)
            stackInfo := utils.GetStackInfo()
            //utils.CallSlack(stackInfo, define.SLACK_PANIC_CHANNEL, define.SLACK_PANIC_SENDFROM)
            logger.Error("magneto panic stackinfo: %s", stackInfo)
        }
    }()

    serverAddr := fmt.Sprintf("%s:%d", define.Cfg.TIP, define.Cfg.Http_server_port)

    //init
    verification.GetInstance().Init(
        define.Cfg.JavaThriftGW.Host,
        define.Cfg.JavaThriftGW.Port,
        define.Cfg.JavaThriftGW.Key,
        define.Cfg.JavaThriftGW.TokenExpire,
    )

    runDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
    txtPath := fmt.Sprintf("%s/../config/list.txt", runDir)
    //cv.MessageFilter.InitChatViolation(txtPath)
    MsgFilter.LoadWordDict(txtPath)

    router := s.Router()
    gin.SetMode(gin.DebugMode)
    router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    return router.Run(serverAddr)
}

func (s *HServer) Router() *gin.Engine {
    r := gin.Default()
    //r.Use(middleware.Cors())
    //r.Use(middleware.SetReqKey())
    //r.Use(middleware.RespLog())
    r.Use(middleware.PanicHandle())
    s.RunDirPath, _ = filepath.Abs(filepath.Dir(os.Args[0]))
    templeteDir := fmt.Sprintf("%s/../web/templates/*", s.RunDirPath)
    r.LoadHTMLGlob(templeteDir)
    staticDir := fmt.Sprintf("%s/../web/static", s.RunDirPath)

    magneto := r.Group("/magneto")
    {
        magneto.Static("/static", staticDir)
        wsClient := magneto.Group("/wscli")
        {
            wsClient.GET("", services.WebSocketClient)
        }
        fmt.Println("web socket client url: http://127.0.0.1:<:port>/wscli")

        rpc := magneto.Group("/rpc")
        {
            rpc.GET("/dashboard", services.Dashboard)
        }

        v1 := magneto.Group("/v1")
        {
            handshakeController(v1)
        }
    }

    return r
}

func getFrontUserInfoFromCtx(ctx *gin.Context) (*verification.FrontEndToken, error) {
    u, exists := ctx.Get("UserInfo")
    userInfo := verification.TokenResolved{}
    if exists {
        userInfo = u.(verification.TokenResolved)
    }else{
        return nil, errors.New("not found user")
    }

    if userInfo.Eor != nil{
        return nil, userInfo.Eor
    }
    d := userInfo.Data.(verification.FrontEndToken)
    return &d, nil
}

func socketHandler(ctx *gin.Context) {

    SocketioServer.On(define.EVENT_CONNECT, func(so socketio.Socket) {
        // Assign the socket to a global variable
        source := strings.ToLower(so.Request().FormValue("source"))
        if source != "pc"{
            source = "h5"
        }

        tokenInfo, err := getFrontUserInfoFromCtx(ctx)
        if err != nil || tokenInfo == nil{
            errMsg := "token fail"
            if err != nil{
                errMsg = err.Error()
            }

            so.Emit(define.EVENT_CHAT_MESSAGE, apimodel.WSRespFmt(
                nil,
                define.ERR_CHECKTOKEN_ERROR,
                errMsg,
                "",
                "",
                0,
            ))
            so.Disconnect()
            return
        }
        logger.Debug(fmt.Sprintf("tokenInfo: %v", tokenInfo))

        canSpeak, err := services.GetUserSpeak(tokenInfo)

        if err != nil{
            logger.Error("get UserSpeak fail, %s", err.Error())
        }

        so.Emit(define.EVENT_USER_SPEAK, apimodel.WSRespFmt(
            canSpeak,
            define.ERR_OK,
            "",
            "",
            define.EVENT_USER_SPEAK,
            define.BROADCAST_TYPE_ONLYYOU,
        ))

        // Concat the platId and memberId as room name => platId:memberId
        //roomName := fmt.Sprintf(define.PLAT_ROOM_NAME, tokenInfo.PlatInfoId , tokenInfo.MemberId)
        roomName := fmt.Sprintf(define.PLAT_ROOM_NAME, tokenInfo.PlatInfoId, source)
        logger.Debug("%v join room : %v", so.Id(), roomName)

        // Join room
        so.Join(roomName)

        logger.Debug("Rooms : %v", SocketioServer.Rooms())
        // listen the chat message event
        so.On(define.EVENT_CHAT_MESSAGE, func(msg string) {
            // Log the message body
            logger.Debug("Request body : %v", msg)

            // Decode request json body
            msgObj := apimodel.ChatMessageRequest{}
            common.Json2Struct(msg, &msgObj)


            msgObj.Msg = strings.Replace(msgObj.Msg, "<", "", -1)
            msgObj.Msg = strings.Replace(msgObj.Msg, ">", "", -1)
            // Swearing Filter
            //oldMsg := msgObj.Msg
            //msgObj.Msg = cv.MessageFilter.WordsFilter(msgObj.Msg)
            msgObj.Msg = MsgFilter.Replace(msgObj.Msg, 42)

            // If the message is empty then do nothing
            if msgObj.Msg == "" {
                return
            }

            respMsg := apimodel.ChatMessageResponse{}
            respMsg.Msg = msgObj.Msg
            respMsg.SetAccount(tokenInfo.MemberAccount)

            msgString := apimodel.WSRespFmt(
                respMsg,
                define.ERR_OK,
                "",
                msgObj.UUID,
                define.EVENT_CHAT_MESSAGE,
                define.BROADCAST_TYPE_ALL,
            )

            // If the sending message is different between original message (be filtered) then do not broadcast
            //if oldMsg != msgObj.Msg {
            //    so.Emit(define.EVENT_CHAT_MESSAGE, msgString)
            //    return
            //}
            platId := strconv.FormatInt(int64(tokenInfo.PlatInfoId), 10)
            MulticastPlatMessage([]string{platId}, msgString)
        })

        so.On(define.EVENT_BET_ORDER, func(msg string) {

            msgObj := apimodel.BetOrderRequest{}
            common.Json2Struct(msg, &msgObj)
            respMsg := apimodel.BetOrderResponse{}
            respMsg.Msg = msgObj.Data
            respMsg.SetAccount(tokenInfo.MemberAccount)

            //send all client
            UnicastMessage(roomName, apimodel.WSRespFmt(
                respMsg,
                define.ERR_OK,
                "",
                msgObj.UUID,
                define.EVENT_BET_ORDER,
                define.BROADCAST_TYPE_UNICAST,
            ))
        })

        so.On(define.EVENT_CALL_CUSTOMER_SERVICE, func(msg string) {
            // TODO: 呼叫客服
        })

        // listen the disconnect event
        so.On(define.EVENT_DISCONNECT, func() {
            logger.Debug("%v on disconnect", so.Id())
        })
    })

    // Socket io error handler
    SocketioServer.On(define.EVENT_ERROR, func(so socketio.Socket, err error) {
        logger.Error("Failed : %v", err)
    })

    SocketioServer.ServeHTTP( ctx.Writer, ctx.Request)

    //[GIN-debug] [WARNING] Headers were already written. Wanted to override status code 200 with 400
    ctx.Abort()
}


func handshakeController(v1 *gin.RouterGroup) *gin.RouterGroup {

    hsGroup := v1.Group("/socket.io")
    {
        hsGroup.GET("/", middleware.CheckToken(), socketHandler)
        hsGroup.POST("/", middleware.CheckToken(), socketHandler)
        hsGroup.Handle("WS", "/", middleware.CheckToken(), socketHandler)
        hsGroup.Handle("WSS", "/", middleware.CheckToken(), socketHandler)
        hsGroup.GET("/unicast", unicastMessageHandler)
        hsGroup.GET("/multicast", multicastPlatMessageHandler)
        hsGroup.GET("/broadcast", broadcastMessageHandler)
        hsGroup.GET("/sessions", sessionCounter)
    }

    return v1
}

func sessionCounter(ctx *gin.Context) {

    ctx.JSON(http.StatusOK, gin.H{
        "code" : http.StatusOK,
        "sessions": SocketioServer.Count(),// session count
        "rooms": SocketioServer.Rooms(),
    })
    return
}

func messageProducer(roomName string, msgString string) {

    userInfo := &kafkaproducer.KfkJobData{
        Topic: define.KAFKA_TOPIC_CHATROOM,
        Key:    roomName,
        Value:  msgString,
    }
    kafkaproducer.Jobchan <- userInfo
}

func MessageConsumer(roomName string, msgString string) {

    var msgData = &apimodel.WSResponse{}
    err := json.Unmarshal([]byte(msgString), msgData)

    if err != nil {
        logger.Error("Decode message failed: %v", err)
    } else {
        SocketioServer.BroadcastTo(roomName, msgData.Event, msgString)
    }
}


func unicastMessageHandler(c *gin.Context) {
    roomName := fmt.Sprintf(define.ROOM_NAME, c.Request.FormValue("platId") ,c.Request.FormValue("memberId"))

    UnicastMessage(roomName, apimodel.WSRespFmt(
        "send to user MSG",
        define.ERR_OK,
        "",
        "",
        "",
        define.BROADCAST_TYPE_UNICAST,
    ))
}


func UnicastMessage(roomName string, msg string) {

    logger.Debug("UnicastMessage Room Name : %v", roomName)
    messageProducer(roomName, msg)
}

func multicastPlatMessageHandler(c *gin.Context) {

    var platSlice []string

    platSlice = append(platSlice, c.Request.FormValue("platId"))

    MulticastPlatMessage(platSlice,
        apimodel.WSRespFmt(
            "send to user MSG",
            define.ERR_OK,
            "",
            "",
            "",
            define.BROADCAST_TYPE_UNICAST,
        ))
}

func MulticastPlatMessage(platMap []string, msg string) {

    for _, v := range platMap {
        logger.Debug("Room Name : %v", v)

        roomH5 := fmt.Sprintf(define.PLAT_ROOM_NAME, v, "h5")
        roomPC := fmt.Sprintf(define.PLAT_ROOM_NAME, v, "pc")

        messageProducer(roomH5, msg)
        messageProducer(roomPC, msg)
    }

    //platConcat := strings.Join(platMap,"|")
    //
    //for _, v := range SocketioServer.Rooms() {
    //    logger.Debug("Room Name : %v", v)
    //    if m, _ := regexp.MatchString("^" + define.PLAT_KEY + "("+platConcat+")_(pc|h5)$", v); m {
    //        messageProducer(v, msg)
    //    }
    //}
}

func broadcastMessageHandler(c *gin.Context) {
    BroadcastMessage(apimodel.WSRespFmt(
        "send to user MSG",
        define.ERR_OK,
        "",
        "",
        "",
        define.BROADCAST_TYPE_ALL,
    ))
}

func BroadcastMessage(msg string) {
    for _, v := range SocketioServer.Rooms() {
        logger.Debug("Room Name : %v", v)
        messageProducer(v, msg)
    }
}
