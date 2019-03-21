package downstream

import (
    "fmt"
    "time"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/internal/define"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/pkg/helpers/ucm"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/pkg/thrift.git/lib/go/thrift"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/pkg/bc/engine/mars"
)


func InitMarsClient(apiService string, stream define.DownStream) error {

    cycle, err := time.ParseDuration(define.Cfg.Mars.Cycle)
    if err != nil {
        return err
    }

    timeout, err := time.ParseDuration(stream.Timeout)
    if err != nil {
        return err
    }

    client := func(addr string, timeout time.Duration) (interface{}, error) {
        var transport thrift.TTransport
        var err error

        transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
        transport, err = thrift.NewTSocketTimeout(addr, timeout)
        if err != nil {
            return nil, fmt.Errorf("Open socket failed, err:%s", err.Error())
        }

        transport, err = transportFactory.GetTransport(transport)
        if err != nil {
            return nil, fmt.Errorf("Get transport failed, err:%s", err.Error())
        }
        if err = transport.Open(); err != nil {
            transport.Close()
            return nil, fmt.Errorf("Open transport failed, err:%s", err.Error())
        }

        protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
        client := mars.NewMarsServiceClientFactory(transport, protocolFactory)
        return client, nil
    }

    isopen := func(c interface{}) bool {
        cli := c.(*mars.MarsServiceClient)
        return cli.Transport.IsOpen()
    }

    down := func(c interface{}) {
        cli := c.(*mars.MarsServiceClient)
        cli.Transport.Close()
    }

    cp, err := ucm.NewChanConnPool(
        &ucm.ConnPoolReq{
            Addrs:             stream.Servers,
            ConnSvrTimeOut:    timeout,
            CheckAddrCycle:    cycle,
            HealthyThreshold:  stream.Healthythreshold,
            MaxCooldownTime:   stream.Maxcooldowntime,
            MinHealthyRatio:   stream.Minhealthyratio,
            RetryTimes:        stream.Retrytimes,
            Size:              stream.Connsize,
            MaxConnBeFalseCnt: stream.Maxconnfalsecnt,
            Create:            client,
            IsOpen:            isopen,
            Down:              down,
        },
    )
    if err != nil {
        return err
    }

    Register(apiService, cp)
    return nil
}
