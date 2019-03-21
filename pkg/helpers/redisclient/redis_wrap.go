/*
   redis访问接口的包装，内部采取连接池实现
*/
package redisclient

import (
	"fmt"
	"rks-golang-practice/pkg/helpers/utils"
	"time"
	"rks-golang-practice/pkg/redigo/redis"
	logger "github.com/shengkehua/xlog4go"
	"errors"
    "strconv"
)

var RedisCli *RedisClient

type RedisClient struct {
	Servers        []string
	ConnTimeoutMs  int
	WriteTimeoutMs int
	ReadTimeoutMs  int

	MaxIdle      int
	MaxActive    int
	IdleTimeoutS int
	Password     string

	current_index int
	pool          *redis.Pool
}

func (client *RedisClient) Close() {
	client.pool.Close()
}

func (client *RedisClient) Init() error {
	if len(client.Servers) == 0 {
		return fmt.Errorf("invalid Redis config servers:%s", client.Servers)
	}

	client.pool = &redis.Pool{
		MaxIdle:        client.MaxIdle,
		IdleTimeout:    time.Duration(client.IdleTimeoutS) * time.Second,
		MaxActive:      client.MaxActive,
		ConnTimeoutMs:  client.ConnTimeoutMs,
		ReadTimeoutMs:  client.ReadTimeoutMs,
		WriteTimeoutMs: client.WriteTimeoutMs,
		Password:       client.Password,
		Dial: func() (redis.Conn, error) {
			var c redis.Conn
			var err error
			for i := 0; i < len(client.Servers)+1; i++ {
				//随机挑选一个IP
				index := utils.RandIntn(len(client.Servers))
				client.current_index = index
				c, err = redis.DialTimeout("tcp", client.Servers[index],
					time.Duration(client.ConnTimeoutMs)*time.Millisecond,
					time.Duration(client.ReadTimeoutMs)*time.Millisecond,
					time.Duration(client.WriteTimeoutMs)*time.Millisecond)
				if err != nil {
					logger.Warn("warning=[redis_connect_failed] num=[%d] server=[%s] err=[%s]",
						i, client.Servers[index], err.Error())
					continue
				}
				//支持密码认证
				if len(client.Password) > 0 {
					if _, err_pass := c.Do("AUTH", client.Password); err_pass != nil {
						c.Close()
					}
				}
				if err == nil {
					break
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			res, err := c.Do("PING")
			logger.Warn("retry happend,check connection invalidity:%v", res)
			return err
		},
	}
	//wgl
	err := client.pool.Init(client.Servers)
	if err != nil {
		logger.Error("error=[init_connect_redis_failed] err=[%s]", err.Error())
	}
	return nil
}

func (client *RedisClient) Mget(key []interface{}) ([]string, error) {

	conn := client.pool.Get(true)
	defer func() {

		conn.Close()
	}()

	value, err := redis.Strings(conn.Do("MGET", key...))
	if err != nil {
		if err == redis.ErrNil {
			logger.Warn("error=[redis_mget_failed] server=[%s] err=[%s]",
				conn.GetAddr(), err.Error())
			return nil, err
		} else {
			logger.Error("error=[redis_mget_failed] server=[%s] err=[%s]",
				conn.GetAddr(), err.Error())
		}

		conn_second := client.pool.Get(false)
		defer conn_second.Close()

		value, err = redis.Strings(conn_second.Do("MGET", key))
		if err != nil {
			if err == redis.ErrNil {
				logger.Warn("second error=[redis_mget_failed] server=[%s]  err=[%s]",
					conn_second.GetAddr(), err.Error())
			} else {
				logger.Error("second error=[redis_mget_failed] server=[%s]  err=[%s]",
					conn_second.GetAddr(), err.Error())
			}
			return nil, err
		}
	}

	return value, nil
}

func (client *RedisClient) Set(key string, value interface{}) error {
	conn := client.pool.Get(true)
	defer conn.Close()

	_, err := conn.Do("SET", key, value)
	if err != nil {
		logger.Error("error=[redis_set_failed] server=[%s] key=[%s] err=[%s]",
			conn.GetAddr(), key, err.Error())

		conn_second := client.pool.Get(false)
		defer conn_second.Close()
		_, err = conn_second.Do("SET", key, value)
		if err != nil {
			logger.Error("second error=[redis_set_failed] server=[%s] key=[%s] err=[%s]",
				conn_second.GetAddr(), key, err.Error())
			return err
		}
	}

	return nil
}

func (client *RedisClient) SetEx(key, value string, livetime int) error {
	conn := client.pool.Get(true)
	defer conn.Close()

	reply ,err := conn.Do("SET", key, value)
	if err != nil {
		logger.Error("error=[redis_set_failed] server=[%s] key=[%s] err=[%s]",
			conn.GetAddr(), key, err.Error())

		conn_second := client.pool.Get(false)
		defer conn_second.Close()
		reply, err = conn_second.Do("SET", key, value)
		if err != nil {
			logger.Error("second error=[redis_set_failed] server=[%s] key=[%s] err=[%s]",
				conn_second.GetAddr(), key, err.Error())
			return err
		}
	}
	client.Expire(key, livetime)

	if reply == "OK" {
		return nil
	} else {
		return errors.New("redisclient: unexpected reply of set")
	}
}

func (client *RedisClient) Get(key string) (string, error) {
	conn := client.pool.Get(true)
	defer conn.Close()

	value, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		if err == redis.ErrNil {
			logger.Warn("error=[redis_get_failed] server=[%s] key=[%s] err=[%s]",
				conn.GetAddr(), key, err.Error())
			return "", err
		} else {
			logger.Error("error=[redis_get_failed] server=[%s] key=[%s] err=[%s]",
				conn.GetAddr(), key, err.Error())
		}

		conn_second := client.pool.Get(false)
		defer conn_second.Close()

		value, err = redis.Bytes(conn_second.Do("GET", key))
		if err != nil {
			if err == redis.ErrNil {
				logger.Warn("second error=[redis_get_failed] server=[%s] key=[%s] err=[%s]",
					conn_second.GetAddr(), key, err.Error())
			} else {
				logger.Error("second error=[redis_get_failed] server=[%s] key=[%s] err=[%s]",
					conn_second.GetAddr(), key, err.Error())
			}
			return "" , err
		}
	}

	return string(value), nil
}

func (client *RedisClient) Rpush(key string, value string) (int64, error) {
	conn := client.pool.Get(true)
	defer conn.Close()

	list_len, err := redis.Int64(conn.Do("RPUSH", key, value))
	if err != nil {
		logger.Error("error=[redis_rpush_failed] server=[%s] key=[%s] err=[%s]",
			conn.GetAddr(), key, err.Error())

		conn_second := client.pool.Get(false)
		defer conn_second.Close()

		list_len, err = redis.Int64(conn_second.Do("RPUSH", key, value))
		if err != nil {
			logger.Error("second error=[redis_rpush_failed] server=[%s] key=[%s] err=[%s]",
				conn_second.GetAddr(), key, err.Error())
			return -1, err
		}
	}

	return list_len, nil
}

func (client *RedisClient) Lpop(key string) (string, error) {
	conn := client.pool.Get(true)
	defer conn.Close()

	value, err := redis.String(conn.Do("LPOP", key))
	if err != nil {
		if err == redis.ErrNil {
			logger.Warn("error=[redis_lpop_failed] server=[%s] key=[%s] err=[%s]",
				conn.GetAddr(), key, err.Error())
			return "", err
		} else {
			logger.Error("error=[redis_lpop_failed] server=[%s] key=[%s] err=[%s]",
				conn.GetAddr(), key, err.Error())
		}

		conn_second := client.pool.Get(false)
		defer conn_second.Close()

		value, err = redis.String(conn_second.Do("LPOP", key))
		if err != nil {
			if err == redis.ErrNil {
				logger.Warn("second error=[redis_lpop_failed] server=[%s] key=[%s] err=[%s]",
					conn_second.GetAddr(), key, err.Error())
			} else {
				logger.Error("second error=[redis_lpop_failed] server=[%s] key=[%s] err=[%s]",
					conn_second.GetAddr(), key, err.Error())
			}
			return "", err
		}
	}

	return value, nil
}

func (client *RedisClient) Llen(key string) (int64, error) {
	conn := client.pool.Get(true)
	defer conn.Close()

	value, err := redis.Int64(conn.Do("LLEN", key))
	if err != nil {
		logger.Error("error=[redis_llen_failed] server=[%s] key=[%s] err=[%s]",
			conn.GetAddr(), key, err.Error())

		conn_second := client.pool.Get(false)
		defer conn_second.Close()

		value, err = redis.Int64(conn_second.Do("LLEN", key))
		if err != nil {
			logger.Error("second error=[redis_llen_failed] server=[%s] key=[%s] err=[%s]",
				conn_second.GetAddr(), key, err.Error())
			return -1, err
		}
	}
	return value, nil
}


func (client *RedisClient) DelSig(key string) (int64, error) {
	conn := client.pool.Get(true)
	defer conn.Close()

	value, err := redis.Int64(conn.Do("DEL", key))
	if err != nil {
		logger.Error("error=[redis_del_failed] server=[%s] keys=[%v] err=[%s]",
			conn.GetAddr(), key, err.Error())

		conn_second := client.pool.Get(false)
		defer conn_second.Close()

		value, err = redis.Int64(conn_second.Do("DEL", key))
		if err != nil {
			logger.Error("second error=[redis_del_failed] server=[%s] keys=[%v] err=[%s]",
				conn_second.GetAddr(), key, err.Error())
			return -1, err
		}
	}

	return value, nil
}

func (client *RedisClient) Del(keys []interface{}) (int64, error) {
	conn := client.pool.Get(true)
	defer conn.Close()

	value, err := redis.Int64(conn.Do("DEL", keys...))
	if err != nil {
		logger.Error("error=[redis_del_failed] server=[%s] keys=[%v] err=[%s]",
			conn.GetAddr(), keys, err.Error())

		conn_second := client.pool.Get(false)
		defer conn_second.Close()

		value, err = redis.Int64(conn_second.Do("DEL", keys...))
		if err != nil {
			logger.Error("second error=[redis_del_failed] server=[%s] keys=[%v] err=[%s]",
				conn_second.GetAddr(), keys, err.Error())
			return -1, err
		}
	}

	return value, nil
}

func (client *RedisClient) Hmset(key string, value []interface{}) (string, error) {
	conn := client.pool.Get(true)
	defer conn.Close()

	//input_params := make([]interface{}, len(value)+1)
	var input_params []interface{}
	input_params = append(input_params, key)
	input_params = append(input_params, value...)
	res, err := redis.String(conn.Do("HMSET", input_params...))
	if err != nil {
		logger.Error("error=[redis_hmset_failed] server=[%s] key=[%s] err=[%s]",
			conn.GetAddr(), key, err.Error())

		conn_second := client.pool.Get(false)
		defer conn_second.Close()

		res, err = redis.String(conn_second.Do("HMSET", input_params...))
		if err != nil {
			logger.Error("second error=[redis_hmset_failed] server=[%s] key=[%s] err=[%s]",
				conn_second.GetAddr(), key, err.Error())
			return "", err
		}
	}

	return res, nil
}

func (client *RedisClient) Hmget(data []interface{}) ([]string, error) {

	conn := client.pool.Get(true)
	defer conn.Close()
	res, err := redis.Strings(conn.Do("HMGET", data...))

	if err != nil {
		logger.Error("error=[redis_hmget_failed] server=[%s] err=[%s]",
			conn.GetAddr(), err.Error())
		conn_second := client.pool.Get(false)
		defer conn_second.Close()
		res, err = redis.Strings(conn_second.Do("HMGET", data...))
		if err != nil {
			logger.Error("error=[redis_hmget_failed] server=[%s] err=[%s]",
				conn_second.GetAddr(), err.Error())
			return nil, err
		}
	}

	return res, err
}

func (client *RedisClient) Hdel(key string, value []interface{}) (int64, error) {
	conn := client.pool.Get(true)
	defer conn.Close()

	//input_params := make([]interface{}, len(value)+1)
	var input_params []interface{}
	input_params = append(input_params, key)
	input_params = append(input_params, value...)
	res, err := redis.Int64(conn.Do("HDEL", input_params...))
	if err != nil {
		logger.Error("error=[redis_hdel_failed] server=[%s] key=[%s] err=[%s]",
			conn.GetAddr(), key, err.Error())

		conn_second := client.pool.Get(false)
		defer conn_second.Close()

		res, err = redis.Int64(conn_second.Do("HDEL", input_params...))
		if err != nil {
			logger.Error("second error=[redis_hdel_failed] server=[%s] key=[%s] err=[%s]",
				conn_second.GetAddr(), key, err.Error())
			return -1, err
		}
	}

	return res, nil
}

func (client *RedisClient) Hkeys(key string) ([]string, error) {
	conn := client.pool.Get(true)
	defer conn.Close()

	res, err := redis.Strings(conn.Do("HKEYS", key))
	if err != nil {
		logger.Error("error=[redis_hkeys_failed] server=[%s] key=[%s] err=[%s]",
			conn.GetAddr(), key, err.Error())

		conn_second := client.pool.Get(false)
		defer conn_second.Close()

		res, err = redis.Strings(conn_second.Do("HKEYS", key))
		if err != nil {
			logger.Error("second error=[redis_hkeys_failed] server=[%s] key=[%s] err=[%s]",
				conn_second.GetAddr(), key, err.Error())
			return nil, err
		}
	}

	return res, nil
}

func (client *RedisClient) Exists(key string) (int64, error) {
	conn := client.pool.Get(true)
	defer conn.Close()

	res, err := redis.Int64(conn.Do("EXISTS", key))
	if err != nil {
		logger.Error("error=[redis_exists_failed] server=[%s] key=[%s] err=[%s]",
			conn.GetAddr(), key, err.Error())

		conn_second := client.pool.Get(false)
		defer conn_second.Close()

		res, err = redis.Int64(conn_second.Do("EXISTS", key))
		if err != nil {
			logger.Error("second error=[redis_exists_failed] server=[%s] key=[%s] err=[%s]",
				conn_second.GetAddr(), key, err.Error())
			return 0, err
		}
	}

	return res, nil
}

func (client *RedisClient) Sismember(key, elem string) (int64, error) {
	conn := client.pool.Get(true)
	defer conn.Close()

	res, err := redis.Int64(conn.Do("SISMEMBER", key, elem))
	if err != nil {
		logger.Error("error=[redis_sismember_failed] server=[%s] key=[%s] err=[%s]",
			conn.GetAddr(), key, err.Error())

		conn_second := client.pool.Get(false)
		defer conn_second.Close()

		res, err = redis.Int64(conn_second.Do("SISMEMBER", key, elem))
		if err != nil {
			logger.Error("second error=[redis_sismember_failed] server=[%s] key=[%s] err=[%s]",
				conn_second.GetAddr(), key, err.Error())
			return 0, err
		}
	}

	return res, nil
}

func (client *RedisClient) MulHexist(key string, fields []string) (interface{}, error) {

	var err error
	var reply interface{}
	conn := client.pool.Get(true)
	defer conn.Close()

	var res []interface{}
	for _, field := range fields {
		var input_params []interface{}
		input_params = append(input_params, key)
		input_params = append(input_params, field)
		err = conn.Send("HEXISTS", input_params...)
	}
	conn.Flush()

	for i := 0; i < len(fields); i = i + 1 {

		reply, err = conn.Receive()
		if err == nil {
			res = append(res, reply)
		} else {
			break
		}
	}

	if err != nil {
		logger.Error("error=[redis_hexists_failed] server=[%s] err=[%s]",
			conn.GetAddr(), err.Error())

		res = append(res[:0], res[len(res):])

		conn_second := client.pool.Get(false)
		defer conn_second.Close()

		for _, field := range fields {
			var input_params []interface{}
			input_params = append(input_params, key)
			input_params = append(input_params, field)
			err = conn_second.Send("HEXISTS", input_params...)
		}

		conn_second.Flush()
		for range fields {
			reply, err = conn_second.Receive()
			if err == nil {
				res = append(res, reply)
			} else {
				break
			}
		}

		if err != nil {
			logger.Error("error=[redis_hexists_failed] server=[%s] err=[%s]",
				conn_second.GetAddr(), err.Error())
			return nil, err
		}
	}
	return res, err
}

func (client *RedisClient) Hgetall(key string) ([]string, error) {
	conn := client.pool.Get(true)
	defer conn.Close()

	res, err := redis.Strings(conn.Do("HGETALL", key))
	if err != nil {
		logger.Error("error=[redis_hgetall_failed] server=[%s] key=[%s] err=[%s]",
			conn.GetAddr(), key, err.Error())

		conn_second := client.pool.Get(false)
		defer conn_second.Close()

		res, err = redis.Strings(conn_second.Do("HGETALL", key))
		if err != nil {
			logger.Error("second error=[redis_hgetall_failed] server=[%s] key=[%s] err=[%s]",
				conn_second.GetAddr(), key, err.Error())
			return nil, err
		}
	}

	return res, nil
}

func (client *RedisClient) Sadd(key string, value []interface{}) (int64, error) {
	conn := client.pool.Get(true)
	defer conn.Close()

	var input_params []interface{}
	input_params = append(input_params, key)
	input_params = append(input_params, value...)
	res, err := redis.Int64(conn.Do("SADD", input_params...))
	if err != nil {
		logger.Error("error=[redis_sadd_failed] server=[%s] key=[%s] err=[%s]",
			conn.GetAddr(), key, err.Error())

		conn_second := client.pool.Get(false)
		defer conn_second.Close()

		res, err = redis.Int64(conn_second.Do("SADD", input_params...))
		if err != nil {
			logger.Error("second error=[redis_sadd_failed] server=[%s] key=[%s] err=[%s]",
				conn_second.GetAddr(), key, err.Error())
			return -1, err
		}
	}

	return res, nil
}

func (client *RedisClient) Smembers(key string) ([]string, error) {
	conn := client.pool.Get(true)
	defer conn.Close()

	res, err := redis.Strings(conn.Do("SMEMBERS", key))
	if err != nil {
		logger.Error("error=[redis_smembers_failed] server=[%s] key=[%s] err=[%s]",
			conn.GetAddr(), key, err.Error())

		conn_second := client.pool.Get(false)
		defer conn_second.Close()

		res, err = redis.Strings(conn_second.Do("SMEMBERS", key))
		if err != nil {
			logger.Error("second error=[redis_smembers_failed] server=[%s] key=[%s] err=[%s]",
				conn_second.GetAddr(), key, err.Error())
			return nil, err
		}
	}

	return res, nil
}

func (client *RedisClient) GetData(cmd string, key string, value ...interface{}) (interface{}, error) {

	conn := client.pool.Get(true)
	defer func() {
		conn.Close()
	}()

	var input_params []interface{}
	input_params = append(input_params, key)
	for _, v := range value {
		input_params = append(input_params, v)
	}
	res, err := conn.Do(cmd, input_params...)

	//res, err := conn.Do(cmd, key)
	if err != nil {
		if err == redis.ErrNil {
			logger.Warn("error=[redis_get_failed] server=[%s] err=[%s]",
				conn.GetAddr(), err.Error())
			return nil, err
		} else {
			logger.Error("error=[redis_get_failed] server=[%s] err=[%s]",
				conn.GetAddr(), err.Error())
		}

		conn_second := client.pool.Get(false)
		defer conn_second.Close()
		res, err = conn_second.Do(cmd, input_params...)
		if err != nil {
			if err == redis.ErrNil {
				logger.Warn("error=[redis_get_failed] server=[%s] err=[%s]",
					conn_second.GetAddr(), err.Error())
			} else {
				logger.Error("error=[redis_get_failed] server=[%s] err=[%s]",
					conn_second.GetAddr(), err.Error())
			}
			return nil, err
		}
	}
	return res, err
}

func (client *RedisClient) GetDataWithPipeline(cmds []string, keys []string) (interface{}, error) {

	var err error
	var reply interface{}

	conn := client.pool.Get(true)
	defer func() {
		conn.Close()
	}()

	var res []interface{}
	for idx, c := range cmds {
		err = conn.Send(c, keys[idx])
	}
	conn.Flush()

	for i := 0; i < len(cmds); i = i + 1 {

		reply, err = conn.Receive()
		if err == nil {
			res = append(res, reply)
		} else {
			break
		}
	}

	if err != nil {
		logger.Error("error=[redis_%s_failed] server=[%s] err=[%s]", cmds,
			conn.GetAddr(), err.Error())

		res = append(res[:0], res[len(res):])
		conn_second := client.pool.Get(false)
		defer conn_second.Close()

		for idx, c := range cmds {
			conn_second.Send(c, keys[idx])
		}
		conn_second.Flush()
		for range cmds {
			reply, err = conn_second.Receive()
			if err == nil {
				res = append(res, reply)
			} else {
				break
			}
		}

		if err != nil {
			logger.Error("error=[redis_%s_failed] server=[%s] err=[%s]", cmds,
				conn_second.GetAddr(), err.Error())
			return nil, err
		}
	}
	return res, err
}

func (client *RedisClient) SetData(cmd string, key string, value []interface{}) (interface{}, error) {

	conn := client.pool.Get(true)
	defer func() {
		conn.Close()
	}()

	var input_params []interface{}
	input_params = append(input_params, key)
	input_params = append(input_params, value...)
	res, err := conn.Do(cmd, input_params...)

	if err != nil {
		logger.Error("error=[redis_%s_failed] server=[%s] key=[%s],err=[%s]", cmd,
			conn.GetAddr(), key, err.Error())

		conn_second := client.pool.Get(false)
		defer conn_second.Close()
		res, err = conn_second.Do(cmd, input_params...)
		if err != nil {
			logger.Error("error=[redis_%s_failed] server=[%s] key=[%s],err=[%s]", cmd,
				conn_second.GetAddr(), key, err.Error())
			return nil, err
		}
	}
	return res, err
}

func (client *RedisClient) SetDataWithPipeline(cmd []string, key []string, value [][]interface{}) (interface{}, error) {

	var err error
	var reply interface{}
	conn := client.pool.Get(true)
	defer conn.Close()

	var res []interface{}
	for idx, c := range cmd {

		var input_params []interface{}
		input_params = append(input_params, key[idx])

		for i := 0; i < len(value[idx]); i = i + 1 {
			input_params = append(input_params, value[idx][i])
		}
		err = conn.Send(c, input_params...)
	}
	conn.Flush()

	for i := 0; i < len(cmd); i = i + 1 {

		reply, err = conn.Receive()
		if err == nil {
			res = append(res, reply)
		} else {
			break
		}
	}

	if err != nil {
		logger.Error("error=[redis_%s_failed] server=[%s] err=[%s]", cmd,
			conn.GetAddr(), err.Error())

		res = append(res[:0], res[len(res):])

		conn_second := client.pool.Get(false)
		defer conn_second.Close()

		for idx, c := range cmd {
			var input_params []interface{}
			input_params = append(input_params, key[idx])

			for i := 0; i < len(value[idx]); i = i + 1 {
				input_params = append(input_params, value[idx][i])
			}

			err = conn_second.Send(c, input_params...)
		}
		conn_second.Flush()
		for range cmd {
			reply, err = conn_second.Receive()
			if err == nil {
				res = append(res, reply)
			} else {
				break
			}
		}

		if err != nil {
			logger.Error("error=[redis_%s_failed] server=[%s] err=[%s]", cmd,
				conn_second.GetAddr(), err.Error())
			return nil, err
		}
	}
	return res, err
}

func (client *RedisClient) Hexists(key, field string) (bool, error) {

	conn := client.pool.Get(true)
	defer conn.Close()
	res, err := redis.Bool(conn.Do("HEXISTS", key, field))

	if err != nil {
		logger.Error("error=[redis_hexists_failed] server=[%s] err=[%s]",
			conn.GetAddr(), err.Error())

		conn_second := client.pool.Get(false)
		defer conn_second.Close()
		res, err = redis.Bool(conn_second.Do("HEXISTS", key, field))
		if err != nil {
			logger.Error("error=[redis_hexists_failed] server=[%s] err=[%s]",
				conn_second.GetAddr(), err.Error())
			return false, err
		}
	}

	return res, err
}

func (client *RedisClient) Hget(key, field string) (string, error) {

	conn := client.pool.Get(true)
	defer conn.Close()
	res, err := redis.String(conn.Do("HGET", key, field))

	if err != nil {
		logger.Error("error=[redis_hget_failed] server=[%s] err=[%s]",
			conn.GetAddr(), err.Error())

		conn_second := client.pool.Get(false)
		defer conn_second.Close()
		res, err = redis.String(conn_second.Do("HGET", key, field))
		if err != nil {
			logger.Error("error=[redis_hget_failed] server=[%s] err=[%s]",
				conn_second.GetAddr(), err.Error())
			return "", err
		}
	}

	return res, err
}

func (client *RedisClient) Rename(oldName, newName string) (string, error) {

	conn := client.pool.Get(true)
	defer conn.Close()
	res, err := redis.String(conn.Do("RENAME"))

	if err != nil {
		logger.Error("error=[redis_rename_failed] server=[%s] err=[%s]",
			conn.GetAddr(), err.Error())
		conn_second := client.pool.Get(false)
		defer conn_second.Close()
		res, err = redis.String(conn_second.Do("RENAME"))
		if err != nil {
			logger.Error("error=[redis_rename_failed] server=[%s] err=[%s]",
				conn_second.GetAddr(), err.Error())
			return "", err
		}
	}
	return res, err
}
func (client *RedisClient) Expire(key string, expiresecond int) error {
	c := client.pool.Get(true)
	defer c.Close()

	reply, err := redis.Int(c.Do("EXPIRE", key, expiresecond))
	if err != nil {
		return err
	}

	if reply == 1 {
		return nil
	} else {
		return errors.New("redisclient: unexpected reply of expire")
	}
}

func (rc *RedisClient) ZAdd(key, value string, score int64) error {
	c := rc.pool.Get(true)
	defer c.Close()

	_, err := redis.Int(c.Do("ZADD", key, score, value))
	if err != nil {
		return err
	}

	return nil

}

func (rc *RedisClient) RPush(key , value string) error {
	c := rc.pool.Get(true)
	defer c.Close()

	_ ,	err := c.Do("RPUSH" , key , value)
	if err != nil {
		return err
	}
	return nil
}

func (rc *RedisClient) LRem(key string , count int64 , value string) (int64, error) {
	c := rc.pool.Get(true)
	defer c.Close()

	reply , err := redis.Int64(c.Do("LREM" , key , count , value))
	if err != nil {
		return reply , err
	}

	return reply , nil
}

func (rc *RedisClient) ZRange(key string, start, stop int) ([]string, error) {
    c := rc.pool.Get(true)
    defer c.Close()

    reply, err := redis.Strings(c.Do("ZRANGE", key, start, stop))
    if err != nil {
        return nil, err
    }

    return reply, nil
}

func (rc *RedisClient) ZRangeWithScores(key string, start, stop int) (map[string]string, error) {
    c := rc.pool.Get(true)
    defer c.Close()

    reply, err := redis.StringMap(c.Do("ZRANGE", key, start, stop, "WITHSCORES"))
    if err != nil {
        return nil, err
    }

    return reply, nil
}

func (rc *RedisClient) ZRangeByScore(key string, min, max int64, minopen, maxopen bool) ([]string, error) {
    c := rc.pool.Get(true)
    defer c.Close()

    minstr := strconv.FormatInt(min, 10)
    maxstr := strconv.FormatInt(max, 10)
    if minopen {
        minstr = "(" + strconv.FormatInt(min, 10)
    }

    if maxopen {
        maxstr = "(" + strconv.FormatInt(max, 10)
    }

    reply, err := redis.Strings(c.Do("ZRANGEBYSCORE", key, minstr, maxstr))
    if err != nil {
        return nil, err
    }

    return reply, nil
}

func (rc *RedisClient) ZRemRangeByScore(key string , min , max int64 ) (int64 , error) {
	c := rc.pool.Get(true)
	defer c.Close()

	minstr := strconv.FormatInt(min, 10)
	maxstr := strconv.FormatInt(max, 10)

	count , err := c.Do("ZREMRANGEBYSCORE", key, minstr, maxstr)
	if err != nil {
		return 0 , err
	}

	return count.(int64) , nil
}

func (rc *RedisClient) ZRangeByScoreWithScores(key string, min, max int, minopen, maxopen bool) (map[string]string, error) {
    c := rc.pool.Get(true)
    defer c.Close()

    minstr := strconv.FormatInt(int64(min), 10)
    maxstr := strconv.FormatInt(int64(max), 10)
    if minopen {
        minstr = "(" + strconv.FormatInt(int64(min), 10)
    }

    if maxopen {
        maxstr = "(" + strconv.FormatInt(int64(max), 10)
    }

    reply, err := redis.StringMap(c.Do("ZRANGEBYSCORE", key, minstr, maxstr, "WITHSCORES"))
    if err != nil {
        return nil, err
    }

    return reply, nil
}

func (rc *RedisClient) ZCount(key string, min, max int64) (int64, error) {
    c := rc.pool.Get(true)
    defer c.Close()

    minstr := strconv.FormatInt(min, 10)
    maxstr := strconv.FormatInt(max, 10)

    if 0 == min {
        minstr = "-inf"
    }

    if -1 == max {
        maxstr = "+inf"
    }

    reply, err := redis.Int64(c.Do("ZCOUNT", key, minstr, maxstr))
    if err != nil {
        return reply, err
    }

    return reply, nil
}

func (rc *RedisClient) Setnx(key, value string) error {
    c := rc.pool.Get(true)
    defer c.Close()

    reply, err := redis.Int(c.Do("SETNX", key, value))
    if err != nil {
        return err
    }

    if reply == 1 {
        return nil
    } else {
        return errors.New("redisclient: setnx fail of key exist")
    }
}

func (rc *RedisClient) SetnxEx(key , value string ,  livetime int) error {
	c := rc.pool.Get(true)
	defer c.Close()

	reply, err := redis.Int(c.Do("SETNX", key, value))
	if err != nil {
		rc.Expire(key, livetime)
		return err
	}

	if reply == 1 {
		rc.Expire(key, livetime)
		return nil
	} else {
		return errors.New("redisclient: setnx fail of key exist")
	}
}

func (client *RedisClient) Ttl(key string) (int, error) {
    c := client.pool.Get(true)
    defer c.Close()

    reply, err := redis.Int(c.Do("TTL", key))
    if err != nil {
        return 0, err
    }

    return reply, nil
}

func (rc *RedisClient) Incr(key string) (int64, error) {
    c := rc.pool.Get(true)
    defer c.Close()

    reply, err := redis.Int64(c.Do("INCR", key))
    if err != nil {
        return -1, err
    }

    if reply >= 1 {
        return reply, nil
    } else {
        return -1, errors.New("redisclient: unexpected reply of incr")
    }
}
