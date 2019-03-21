package downstream

import (
	logger "github.com/shengkehua/xlog4go"
	"gitlab.kingbay-tech.com/engine-lottery/magneto/pkg/helpers/redisclient"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/internal/define"
)

func InitRedisClient() error {
    // 初始化redis连接池
	redisclient.RedisCli = &redisclient.RedisClient{
		Servers:        define.Cfg.RedisCfg.Redis_svr,
		ConnTimeoutMs:  define.Cfg.RedisCfg.Redis_conn_timeout,
		WriteTimeoutMs: define.Cfg.RedisCfg.Redis_write_timeout,
		ReadTimeoutMs:  define.Cfg.RedisCfg.Redis_read_timeout,
		MaxIdle:        define.Cfg.RedisCfg.Redis_max_idle,
		MaxActive:      define.Cfg.RedisCfg.Redis_max_active,
		IdleTimeoutS:   define.Cfg.RedisCfg.Redis_expire_second,
		Password:       "",
	}
	err := redisclient.RedisCli.Init()
	if err != nil {
		logger.Error("init redis failed, err:%s", err.Error())
		return err
	}
	return nil
}
