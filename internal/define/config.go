package define

type ServiceConfig struct {
    LogFile          string `flag:"logfile" cfg:"logfile" toml:"logfile"`
    Http_server_port int    `flag:"http_server_port" cfg:"http_server_port" toml:"http_server_port"`

    TIP      string `flag:"tip" cfg:"tip"`
    TPort    int    `flag:"tport" cfg:"tport"`
    TTimeout string `flag:"ttimeout" cfg:"ttimeout"`
    Mars     DownStream

    Downstreams   map[string]config_section
    Apollo        Apollo_settting
    ApolloDefault DefaultConfiguration

    RedisCfg RedisConfig

    JavaThriftGW JavaThrift

    KafkaProducerCfg KafkaProducer

    KafkaConsumerCfg map[string]KafkaConsumer
}

type Apollo_settting struct {
    ConfigFileName string
    ApinameListKey string
    AppId          string
    Cluster        string
    NamespaceName  string
    ConfigServer   string
}

type config_section map[string]interface{}

type DownStream struct {
    Servers          []string `flag:"servers" cfg:"servers"`
    Healthythreshold int64    `flag:"healthythreshold" cfg:"healthythreshold"`
    Maxcooldowntime  int64    `flag:"maxcooldowntime" cfg:"maxcooldowntime"`
    Minhealthyratio  float64  `flag:"minhealthyratio" cfg:"minhealthyratio"`
    Connsize         int      `flag:"connsize" cfg:"connsize"`
    Timeout          string   `flag:"timeout" cfg:"timeout"`
    Cycle            string   `flag:"cycle" cfg:"cycle"`
    Retrytimes       int      `flag:"retrytimes" cfg:"retrytimes"`
    Maxconnfalsecnt  int      `flag:"maxconnfalsecnt" cfg:"maxconnfalsecnt"`
    CheckAliveConnSessionCycle string   `flag:"checkaliveconnsessioncycle" cfg:"checkaliveconnsessioncycle"`
    ConnSessionAliveTime       string   `flag:"connsessionalivetime" cfg:"connsessionalivetime"`
}

type RedisConfig struct {
    Redis_svr           []string // redis server
    Redis_conn_timeout  int      // redis连接超时 毫秒
    Redis_read_timeout  int      // redis读超时 毫秒
    Redis_write_timeout int      // redis写超时 毫秒
    Redis_max_idle      int      // 最大空闲连接
    Redis_max_active    int      // 最大活动连接
    Redis_expire_second int      // redis数据过期时间 秒 线上配置10分钟 600
}

type JavaThrift struct {
    Host string
    Port string
    Key string
    TokenExpire int
}


// Kafka Producer Config
type KafkaProducer struct {
    BrokerList   string
    BatchNum     int
    PartitionNum int
    ProduceNum   int
    DialTimeout     int
    WriteTimeout    int
    ReadTimeout     int
    ReturnError     bool
    ReturnSuccess   bool
    FlushFrequency  int
    ChannelBufferSize int

}

type KafkaConsumer struct {
    Topic                   string
    Group                   string
    ProcessTimeout          int
    CommitInterval          int
    RetryTimes              int
    MetaMaxRetry            int
    ChannelSize             int
    HttpTimeout             int
    ZookeeperTimeout        int
    MetaRefreshFrequency    int
    ZookeeperChroot         string
    RetryInterval           int
    ZookeeperAddresses      []string
}