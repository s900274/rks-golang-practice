package server


import (
    "fmt"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/pkg/thrift.git/lib/go/thrift"
    logger "github.com/shengkehua/xlog4go"
    "strconv"
    "time"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/internal/define"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/pkg/bc/engine/mars"
)

type MarsServer struct {
    tip         string
    tport       int64
    ttimeout    time.Duration // thrift 接口超时时间
}

func (s *MarsServer) getThriftListentAddr() string {
    return s.tip + ":" + strconv.FormatInt(int64(s.tport), 10)
}

func NewMarsServer() *MarsServer {
    s := new(MarsServer)
    return s
}

func InitMarsServer() (*MarsServer, error) {
    s := NewMarsServer()
    if s == nil {
        return nil, fmt.Errorf("NewMarsServer fail")
    }

    if err := s.initThriftConf(); err != nil {
        return nil, err
    }

    return s, nil
}

func (s *MarsServer) initThriftConf() error {
    s.tip = define.Cfg.TIP
    s.tport = int64(define.Cfg.TPort)
    var err error
    s.ttimeout, err = time.ParseDuration(define.Cfg.TTimeout)
    return err
}

func (s *MarsServer) StartService() error {
    taddr := s.getThriftListentAddr()
    tf := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
    pf := thrift.NewTBinaryProtocolFactoryDefault()
    ss, err := thrift.NewTServerSocketTimeout(taddr, s.ttimeout)
    if err != nil {
        logger.Fatal("msg=[start service fail] detail=[new thrift server fail] err=[%s] instance=[%s]", err.Error())
        return err
    }

    processor := mars.NewMarsServiceProcessor(s)
    server := thrift.NewTSimpleServer4(processor, ss, tf, pf)
    logger.Info("msg=[start test server succ] addr=[%s]", taddr)
    err = server.Serve()
    if nil != err {
        logger.Info("test server failed. err=[%s]", err.Error())
    } else {
        logger.Info("test server will stop. err=[nil]")
    }

    return nil
}

func RunMarsServerr() {
    s, err := InitMarsServer()
    if err != nil {
        logger.Fatal("msg=[init engindemo server fail] err=[%s]", err.Error())
        return
    }

    s.StartService()
}

