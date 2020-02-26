package cache

import (
	"bookzone/util/log"
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"strconv"
	"time"
	"bookzone/conf"
)

var (
	redisPool 		*redis.Pool
)

func init() {
	var err error
	redisHostKey := "redis_host"
	redisPortKey := "redis_port"
	maxIdleKey := "max_idle"
	maxActiveKey := "max_active"
	idleTimeoutKey := "idle_timeout_min"
	redisHost := conf.GlobalCfg.Section("redis").Key(redisHostKey).String()
	redisPort := conf.GlobalCfg.Section("redis").Key(redisPortKey).String()
	maxIdleStr := conf.GlobalCfg.Section("redis").Key(maxIdleKey).String()
	maxActiveStr := conf.GlobalCfg.Section("redis").Key(maxActiveKey).String()
	idleTimeoutStr := conf.GlobalCfg.Section("redis").Key(idleTimeoutKey).String()
	maxIdle, err := strconv.Atoi(maxIdleStr)
	if err != nil {
		maxIdle = 20
	}
	maxAvtive, err := strconv.Atoi(maxActiveStr)
	if err != nil {
		maxAvtive = 20
	}
	idleTimeout, err := strconv.Atoi(idleTimeoutStr)
	if err != nil {
		idleTimeout = 10
	}

	redisPool = &redis.Pool{
		MaxIdle: maxIdle,
		MaxActive: maxAvtive,
		IdleTimeout: time.Duration(idleTimeout) * time.Minute,
		Dial: func() (redis.Conn, error) {
			return redis.Dial(
				"tcp",
				fmt.Sprintf("%s:%s", redisHost, redisPort),
				redis.DialReadTimeout(1 * time.Second),
				redis.DialWriteTimeout(1 * time.Second),
				redis.DialConnectTimeout(1 * time.Second))
		},
	}
}

func redisDo(cmd string, key interface{}, args ...interface{}) (interface{}, error) {
	connection := redisPool.Get()
	defer connection.Close()
	if err := connection.Err(); err != nil {
		return nil, err
	}

	params := make([]interface{}, 0)
	params = append(params, key)

	if len(args) > 0 {
		for _, v := range args {
			params = append(params, v)
		}
	}

	return connection.Do(cmd, params...)
}

func WriteString(key string, value string, expire int64) error {
	_, err := redisDo("SET", key, value)
	if err != nil {
		log.Errorf(err.Error())
		return err
	}
	log.Infof("redis set %s:%s", key, value)

	if expire == 0 {
		return nil
	}
	_, err = redisDo("EXPIRE", key, expire)
	if err != nil {
		log.Errorf(err.Error())
		return err
	}
	return nil
}

func ReadString(key string) (string, error) {
	ret , err := redisDo("GET", key)
	if err != nil {
		log.Errorf(err.Error())
		return "", err
	} else {
		str, _ := redis.String(ret, err)
		return str, nil
	}
}

func WriteStruct(key string, obj interface{}, expire int64) error {
	ret, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	return WriteString(key, string(ret), expire)
}

func ReadStruct(key string, obj interface{}) error {
	var data string
	var err error
	if data, err = ReadString(key); err != nil {
		log.Errorf(err.Error())
		return err
	}

	err = json.Unmarshal([]byte(data), obj)
	return err
}

func WriteList(key string, list interface{}, total int) error {
	listKey := key + "_list"
	countKey := key + "_list"
	data, err := json.Marshal(list)
	if err != nil {
		log.Errorf(err.Error())
		return err
	}

	err = WriteString(countKey, strconv.Itoa(total), 0)
	if err != nil {
		return err
	}

	err = WriteString(listKey, string(data), 0)
	if err != nil {
		return err
	}
	return nil
}

func ReadList(key string, list interface{}) (int, error) {
	var countStr string
	var dataStr string
	var err error

	listKey := key + "_list"
	countKey := key + "_list"

	if countStr, err = ReadString(countKey); err != nil {
		return 0, err
	}
	count, err := strconv.Atoi(countStr)
	if err != nil {
		return 0, err
	}

	if dataStr, err = ReadString(listKey); err != nil {
		return count, err
	}

	err = json.Unmarshal([]byte(dataStr), list)
	if err != nil {
		return count, err
	}

	return count, nil
}