package kafkaconsumer

import (
    "fmt"
    "github.com/Shopify/sarama"
    "github.com/satori/go.uuid"
    logger "github.com/shengkehua/xlog4go"
    "github.com/wvanbergen/kafka/consumergroup"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/internal/define"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/internal/httpservice"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/pkg/helpers/utils"
    "log"
    "os"
    "time"
)

type KfkJobData struct {
    Topic   string
    Key     string
    Value   string
}

type SimConsumer struct {
    Config        *consumergroup.Config
    Consumer      *consumergroup.ConsumerGroup
    Offsets       map[string]map[int32]int64
}

func NewConsumerConfiguration(topic string) *consumergroup.Config {
    // consumer config
    config := consumergroup.NewConfig()
    config.Offsets.Initial              = sarama.OffsetNewest
    config.Offsets.ProcessingTimeout    = time.Duration(define.Cfg.KafkaConsumerCfg[topic].ProcessTimeout) * time.Millisecond
    config.Offsets.CommitInterval       = time.Duration(define.Cfg.KafkaConsumerCfg[topic].CommitInterval) * time.Millisecond
    config.Consumer.Retry.Backoff       = time.Duration(define.Cfg.KafkaConsumerCfg[topic].RetryInterval) * time.Millisecond
    config.ChannelBufferSize            = define.Cfg.KafkaConsumerCfg[topic].ChannelSize
    config.Metadata.Retry.Max           = define.Cfg.KafkaConsumerCfg[topic].MetaMaxRetry
    config.Metadata.Retry.Backoff       = time.Duration(define.Cfg.KafkaConsumerCfg[topic].MetaRefreshFrequency) * time.Millisecond
    config.Zookeeper.Chroot             = define.Cfg.KafkaConsumerCfg[topic].ZookeeperChroot

    return config
}

func NewKafkaConsumer(topic string) (*SimConsumer, error) {

    sarama.Logger = log.New(os.Stdout, "", log.Ltime)

    consumerGroup := &SimConsumer{}

    consumerGroup.Config = NewConsumerConfiguration(topic)

    logger.Debug("=== Start Subscribe: ", topic)

    groupName := fmt.Sprintf(define.CONSUMER_GROUP, uuid.Must(uuid.NewV4()))
    //machineId, _ := machineid.ID()
    //groupName := fmt.Sprintf(define.CONSUMER_GROUP, machineId)
    // join to consumer group
    var err error
    consumerGroup.Consumer, err = consumergroup.JoinConsumerGroup(groupName, []string{topic}, define.Cfg.KafkaConsumerCfg[topic].ZookeeperAddresses, consumerGroup.Config)
    if err != nil {
        return nil, err
    }

    return consumerGroup, err
}

func (sap *SimConsumer) Start() {
    defer func() {
        err := recover()
        if err != nil {
            logger.Error("magneto panic err: %s", err)
            stackInfo := utils.GetStackInfo()
            utils.CallSlack(stackInfo, define.SLACK_PANIC_CHANNEL, define.SLACK_PANIC_SENDFROM)
            logger.Error("magneto panic stackinfo: %s", stackInfo)
        }
    }()

    sap.Offsets = make(map[string]map[int32]int64)

    for {
        select {
        case msg := <-sap.Consumer.Messages():
            if sap.Offsets[msg.Topic] == nil {
                sap.Offsets[msg.Topic] = make(map[int32]int64)
            }
            if sap.Offsets[msg.Topic][msg.Partition] != 0 && sap.Offsets[msg.Topic][msg.Partition] != msg.Offset-1 {
                sarama.Logger.Println("Unexpected offset on %s:%d. Expected %d, found %d, diff %d.\n", msg.Topic, msg.Partition, sap.Offsets[msg.Topic][msg.Partition]+1, msg.Offset, msg.Offset-sap.Offsets[msg.Topic][msg.Partition]+1)
                // continue
            }

            logger.Debug("Topic: %s", msg.Topic)
            logger.Debug("Partition: %v", msg.Partition)
            logger.Debug("Value: %s", string(msg.Value))
            logger.Debug("Timestamp: %v", msg.Timestamp)
            httpservice.MessageConsumer(string(msg.Key), string(msg.Value))

            sap.Offsets[msg.Topic][msg.Partition] = msg.Offset

            // commit to zookeeper that message is read
            // this prevent read message multiple times after restart
            err := sap.Consumer.CommitUpto(msg)
            if err != nil {
                fmt.Println("Error commit zookeeper: ", err.Error())
            }
        }
    }
}
