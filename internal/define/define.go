package define

var Cfg ServiceConfig

const (
    GW_APINAME  = "mars"
    CALLER_NAME = "golang-magneto"

    //注單系統相關
    SELENE_APINAME                  = "selene"
    SELENE_VERSION             int8 = 1
    SELENE_CROWDFUNDING_DETAIL      = 0x20000B //已發起眾籌詳情

    //chronus
    CHRONUS_APINAME            = "chronus"
    CHRONUS_VERSION       int8 = 1
    CHRONUS_GET_USERSPEAK      = 0x400007 //取得是否可說話

    SLACK_PANIC_CHANNEL  = "golangpanics"
    SLACK_PANIC_SENDFROM = "magneto"
)

const (
    EVENT_CHAT_MESSAGE          = "chat message"
    EVENT_BET_ORDER             = "bet order"
    EVENT_CALL_CUSTOMER_SERVICE = "customer service"
    EVENT_USER_SPEAK            = "user speak"
    EVENT_CONNECT               = "connection"
    EVENT_DISCONNECT            = "disconnection"
    EVENT_ERROR                 = "error"

    ROOM_NAME      = "%v:%v"
    PLAT_KEY       = "plat:"
    PLAT_ROOM_NAME = PLAT_KEY + "%v_%s"
    CONSUMER_GROUP = "CHATROOM_GROUP_%v"

    KAFKA_TOPIC_CHATROOM = "CHATROOM"
)

const (
    BROADCAST_TYPE_ONLYYOU = 0
    BROADCAST_TYPE_UNICAST = 1
    BROADCAST_TYPE_ALL     = 2
)
