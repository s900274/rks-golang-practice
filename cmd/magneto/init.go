package main

import (
    "errors"
    "flag"
    "github.com/BurntSushi/toml"
    "github.com/mreiferson/go-options"
    logger "github.com/shengkehua/xlog4go"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/internal/define"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/internal/httpservice"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/internal/kafkaconsumer"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/pkg/helpers/kafkaproducer"
    "log"
    "os"
    "strings"
    "sync"
    "github.com/zouyx/agollo"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/pkg/helpers/common"
    "reflect"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/pkg/downstream"
    "fmt"
)


func initConf(flagSet *flag.FlagSet) error {

    flagSet.Parse(os.Args[1:])
    if len(flagSet.Lookup("config").Value.(flag.Getter).Get().(string)) > 0 {
        configFile := flagSet.Lookup("config").Value.String()
        if configFile != "" {
            _, err := toml.DecodeFile(configFile, &define.Cfg)
            if err != nil {
                log.Fatalf("ERROR: failed to load config file %s - %s\n", configFile, err.Error())
                return err
            }

            for tag, down_cfg := range define.Cfg.Downstreams {
                switch {
                case strings.EqualFold(tag, "mars"):
                    options.Resolve(&define.Cfg.Mars, flagSet, down_cfg)
                }
            }

        } else {
            log.Fatalln("ERROR: config file is nil")
            err := errors.New("ERROR: config file is nil")
            return err
        }
    } else {
        log.Fatalln("ERROR: no config param given")
        err := errors.New("ERROR: no config param given")
        return err
    }

    return nil
}

func initLogger() error {
    err := logger.SetupLogWithConf(define.Cfg.LogFile)
    return err
}

func initKafkaProducer() error {
    for tag, _ := range define.Cfg.KafkaConsumerCfg {
        _ = kafkaproducer.AddNewTopic(define.Cfg.KafkaProducerCfg.BrokerList, define.Cfg.KafkaProducerCfg.PartitionNum, tag)
        logger.Debug("Add topic %v for %v partitions", tag, define.Cfg.KafkaProducerCfg.PartitionNum)
    }

    for i := 0; i < define.Cfg.KafkaProducerCfg.ProduceNum; i++ {
        producer, err := kafkaproducer.NewSimAsyncProducer(define.Cfg.KafkaProducerCfg.BrokerList, define.Cfg.KafkaProducerCfg.PartitionNum)
        if err != nil {
            logger.Error("NewSimAsyncProducer failed;err:%s", err.Error())
            return err
        }
        producer.Start()
        kafkaproducer.ProducerList = append(kafkaproducer.ProducerList, producer)
    }
    return nil
}

func initKafkaConsumer() error {
    logger.Debug("init kafka consumer")
    logger.Debug("init kafka topic cfg : %v", define.Cfg.KafkaConsumerCfg)
    for tag, _ := range define.Cfg.KafkaConsumerCfg {

        logger.Debug("init kafka topic : %v", tag)
        consumerGroup, err := kafkaconsumer.NewKafkaConsumer(tag)
        if err != nil {
            logger.Error("New Kafka Consumer failed;err:%s", err.Error())
            return err
        }
        go consumerGroup.Start()
    }
    return nil
}

func initDownStreamsClient() error {

    // init apollo config
    readyConfig := &agollo.AppConfig{
        AppId:         define.Cfg.Apollo.AppId,
        Cluster:       define.Cfg.Apollo.Cluster,
        NamespaceName: define.Cfg.Apollo.NamespaceName,
        Ip:            define.Cfg.Apollo.ConfigServer,
    }
    agollo.InitCustomConfig(func() (*agollo.AppConfig, error) {
        return readyConfig, nil
    })

    // listen apollo
    agollo.Start()
    logger.Info("Apollo Server: %s", define.Cfg.Apollo.ConfigServer)

    // get all apiName keys
    defApiNameStr, _ := common.Struct2Json(define.Cfg.ApolloDefault.ApiNameList)
    apiNmaeStr := agollo.GetStringValue(define.Cfg.Apollo.ApinameListKey, defApiNameStr)
    var apiNameList []string
    _ = common.Json2Struct(apiNmaeStr, &apiNameList)

    //初始化 public service連接池
    for _, element := range apiNameList {
        stream := define.Cfg.Mars
        serverInfo := define.ServiceInfo{}
        defEleStr, _ := common.Struct2Json(define.Cfg.ApolloDefault.ApiInfo[element])
        eleStr := agollo.GetStringValue(element, defEleStr)
        _ = common.Json2Struct(eleStr, &serverInfo)
        stream.Servers = serverInfo.Servers
        if err := downstream.InitMarsClient(element, stream); err != nil {
            return err
        }
    }

    go func() {
        for {
            event := agollo.ListenChangeEvent()
            changeEvent := <-event

            apiNmaeStr := agollo.GetStringValue(define.Cfg.Apollo.ApinameListKey, "[]")
            var apiNameList []string
            _ = common.Json2Struct(apiNmaeStr, &apiNameList)
            logger.Info("> do apollo change apinameKeys: %v", apiNmaeStr)

            for _, apiName := range apiNameList {
                info := changeEvent.Changes[apiName]
                //0: add
                //1: change
                //2: del
                if info.ChangeType == 1 {
                    serverInfo := &define.ServiceInfo{}
                    newDataStr := info.NewValue
                    _ = common.Json2Struct(newDataStr, serverInfo)
                    go ChangeApiServiceFromApollo(apiName, serverInfo)
                }
            }
        }
    }()

    keys := reflect.ValueOf(downstream.DownStreamMgr).MapKeys()
    fmt.Println(keys)

    return nil
}

func RunHttpServer()  {
    var wg sync.WaitGroup

    wg.Add(1)
    go func() {
        defer wg.Done()
        HServer := httpservice.NewHTTPServer()
        err := HServer.InitHttpServer()
        if nil != err {
            logger.Error("HTTPServerStart failed, err :%v", err)
            return
        }
    }()
    wg.Wait()
}