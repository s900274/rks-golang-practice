package kafkaproducer

import (
    "encoding/json"
    "github.com/Shopify/sarama"
    logger "github.com/shengkehua/xlog4go"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/internal/define"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/pkg/helpers/utils"
    "math/rand"
    "strconv"
    "strings"
    "time"
)

var (
    Jobchan      = make(chan *KfkJobData, define.Cfg.KafkaProducerCfg.ChannelBufferSize)
    ProducerList = make([]*SimAsyncProducer, 0)
)

type KfkJobData struct {
    Topic string
    Key   string
    Value string
}

func NewProducerConfiguration() *sarama.Config {
    config := sarama.NewConfig()
    config.Net.DialTimeout = time.Duration(define.Cfg.KafkaProducerCfg.DialTimeout) * time.Millisecond
    config.Net.WriteTimeout = time.Duration(define.Cfg.KafkaProducerCfg.WriteTimeout) * time.Millisecond
    config.Net.ReadTimeout = time.Duration(define.Cfg.KafkaProducerCfg.ReadTimeout) * time.Millisecond

    config.Producer.Return.Errors = define.Cfg.KafkaProducerCfg.ReturnError
    config.Producer.Return.Successes = define.Cfg.KafkaProducerCfg.ReturnSuccess
    config.Producer.Flush.Messages = define.Cfg.KafkaProducerCfg.BatchNum
    config.Producer.Flush.MaxMessages = define.Cfg.KafkaProducerCfg.BatchNum * 3
    config.Producer.Flush.Frequency = time.Duration(define.Cfg.KafkaProducerCfg.FlushFrequency) * time.Millisecond
    config.Producer.Partitioner = sarama.NewHashPartitioner
    config.Producer.RequiredAcks = sarama.WaitForAll
    //使用snappy压缩
    config.Producer.Compression = sarama.CompressionSnappy

    config.Producer.Partitioner = sarama.NewManualPartitioner
    config.ChannelBufferSize = define.Cfg.KafkaProducerCfg.ChannelBufferSize
    return config
}

type SimAsyncProducer struct {
    BrokerList      string
    PartitionNum    int
    config          *sarama.Config
    asyncProducer   sarama.AsyncProducer
    RandomPartition *rand.Rand
}

func NewSimAsyncProducer(brokerlist string, partitionnum int) (*SimAsyncProducer, error) {

    s1 := rand.NewSource(time.Now().UnixNano())

    producer := &SimAsyncProducer{
        BrokerList:      brokerlist,
        PartitionNum:    partitionnum,
        RandomPartition: rand.New(s1),
    }
    producer.config = NewProducerConfiguration()
    var err error
    producer.asyncProducer, err = sarama.NewAsyncProducer(strings.Split(brokerlist, ","), producer.config)
    if err != nil {
        logger.Error("Failed to NewSimAsyncProducer;err:%s", err.Error())
        return nil, err
    }
    return producer, nil
}

func AddNewTopic(brokerlist string, partitionnum int, topic string) (error) {

    config := sarama.NewConfig()

    config.Version = sarama.V1_0_0_0

    admin, err := sarama.NewClusterAdmin(strings.Split(brokerlist, ","), config)

    if err != nil {
        logger.Error("Failed to Connect kafka broker;err: %v", err)
        return err
    }

    detail := sarama.TopicDetail{NumPartitions: int32(partitionnum), ReplicationFactor: 1}
    err = admin.CreateTopic(topic, &detail, false)
    if err != nil {
        logger.Error("Failed to add topic;err: %v", err)
        return err
    }
    return nil
}

func (sap *SimAsyncProducer) Start() {
    go sap.work()
}

func (sap *SimAsyncProducer) work() {
    defer func() {
        err := recover()
        if err != nil {
            stackInfo := utils.GetStackInfo()
            utils.CallSlack(stackInfo, define.SLACK_PANIC_CHANNEL, define.SLACK_PANIC_SENDFROM)
            logger.Error("magneto panic stackinfo: %s", stackInfo)
        }
    }()

    for {
        select {
        case success := <-sap.asyncProducer.Successes():
            key, _ := success.Key.Encode()
            value, _ := success.Value.Encode()
            logger.Info("%s Topic asyncProducer success key:%s value:%v", success.Topic, string(key), string(value))
        case err := <-sap.asyncProducer.Errors():
            logger.Error("asyncProducer err:%v", err)
        case msg := <-Jobchan:
            logger.Info("%s Topic record to Producer.", msg.Topic)
            sap.produce(msg)
        }
    }
}

func (sap *SimAsyncProducer) produce(msg *KfkJobData) {

    logger.Debug("msg:%s", msg)
    message := &sarama.ProducerMessage{
        Topic:     msg.Topic,
        Key:       sarama.StringEncoder(DigToString(msg.Key)),
        Value:     sarama.ByteEncoder(msg.Value),
        Partition: int32(sap.RandomPartition.Intn(define.Cfg.KafkaProducerCfg.PartitionNum)),
    }
    sap.asyncProducer.Input() <- message
}

func (sap *SimAsyncProducer) Stop() {
    logger.Info("stop SimAsyncProducer")
    sap.asyncProducer.AsyncClose()
}

func DigToString(a interface{}) string {
    if v, p := a.(json.Number); p {
        return v.String()
    }
    if v, p := a.(int); p {
        return strconv.Itoa(v)
    }
    if v, p := a.(float64); p {
        return strconv.FormatFloat(v, 'f', -1, 64)
    }
    if v, p := a.(float32); p {
        return strconv.FormatFloat(float64(v), 'f', -1, 32)
    }
    if v, p := a.(int16); p {
        return strconv.Itoa(int(v))
    }
    if v, p := a.(uint); p {
        return strconv.Itoa(int(v))
    }
    if v, p := a.(int32); p {
        return strconv.Itoa(int(v))
    }
    if v, p := a.(int64); p {
        return strconv.Itoa(int(v))
    }
    if v, p := a.(string); p {
        return v
    }
    return "error"
}
