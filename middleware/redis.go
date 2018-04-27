package middleware

import (
	"encoding/json"
	"strings"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/labstack/echo"
	"github.com/shibingli/realclouds_go/utils"
)

func DefaultRedis() (*Redis, error) {

	host := utils.GetENV("REDIS_HOST")
	if len(host) == 0 {
		host = "127.0.0.1:6379"
	}

	password := utils.GetENV("REDIS_PASSWORD")
	if len(password) == 0 {
		password = ""
	}

	dataBase := utils.GetENV("REDIS_DB")
	if len(dataBase) == 0 {
		dataBase = ""
	}

	maxIdleConns, err := utils.GetENVToInt("REDIS_MAXIDLECONNS")
	if nil != err {
		maxIdleConns = 10
	}

	idleTimeout, err := utils.GetENVToInt("REDIS_IDLETIMEOUT")
	if nil != err {
		idleTimeout = 300
	}

	rPool, err := NewRedis(host, password, maxIdleConns, idleTimeout, dataBase)
	if nil != err {
		return nil, err
	}

	redis := &Redis{
		RedisPool: rPool,
	}

	return redis, nil
}

type Redis struct {
	RedisPool *redis.Pool
	Mutex     sync.RWMutex
}

//MwMySQL MySQL middleware
func (r *Redis) MwRedis(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		r.Mutex.Lock()
		defer r.Mutex.Unlock()
		c.Set("redis", r)
		return next(c)
	}
}

func (r *Redis) Setex(key string, exp int, value interface{}) (err error) {
	key = strings.TrimSpace(key)
	conn := r.RedisPool.Get()
	defer conn.Close()
	if err = conn.Err(); err != nil {
		return
	}
	_, err = conn.Do("SETEX", key, exp, value)
	if nil != err {
		return
	}
	return
}

func (r *Redis) Del(key string) (err error) {
	key = strings.TrimSpace(key)
	conn := r.RedisPool.Get()
	defer conn.Close()
	if err = conn.Err(); err != nil {
		return
	}
	_, err = conn.Do("DEL", key)
	if nil != err {
		return
	}
	return
}

func (r *Redis) Get(key string) (data interface{}, err error) {
	key = strings.TrimSpace(key)
	conn := r.RedisPool.Get()
	defer conn.Close()
	if err = conn.Err(); err != nil {
		return
	}
	data, err = conn.Do("GET", key)
	if nil != err {
		return
	}
	return
}

func (r *Redis) GetString(key string) (string, error) {
	return redis.String(r.Get(key))
}

func (r *Redis) GetInterface(key string, inf interface{}) (str string, err error) {
	str, err = redis.String(r.Get(key))
	if nil != err {
		return
	}
	err = json.Unmarshal([]byte(str), &inf)
	return
}

func (r *Redis) GetStrings(key string) ([]string, error) {
	return redis.Strings(r.Get(key))
}

func (r *Redis) GetStringMap(key string) (map[string]string, error) {
	return redis.StringMap(r.Get(key))
}

func (r *Redis) GetInt(key string) (int, error) {
	return redis.Int(r.Get(key))
}

func (r *Redis) GetInts(key string) ([]int, error) {
	return redis.Ints(r.Get(key))
}

func (r *Redis) GetIntMap(key string) (map[string]int, error) {
	return redis.IntMap(r.Get(key))
}

func (r *Redis) GetInt64(key string) (int64, error) {
	return redis.Int64(r.Get(key))
}

func (r *Redis) GetInt64Map(key string) (map[string]int64, error) {
	return redis.Int64Map(r.Get(key))
}

func (r *Redis) GetUint64(key string) (uint64, error) {
	return redis.Uint64(r.Get(key))
}

func (r *Redis) GetFloat64(key string) (float64, error) {
	return redis.Float64(r.Get(key))
}

func (r *Redis) GetBytes(key string) ([]byte, error) {
	return redis.Bytes(r.Get(key))
}

func (r *Redis) GetByteSlices(key string) ([][]byte, error) {
	return redis.ByteSlices(r.Get(key))
}

func (r *Redis) GetBool(key string) (bool, error) {
	return redis.Bool(r.Get(key))
}

func NewRedis(host, password string, maxIdle, idleTimeout int, db string) (rPool *redis.Pool, err error) {

	host = strings.TrimSpace(host)
	password = strings.TrimSpace(password)

	rPool = &redis.Pool{
		MaxIdle:     maxIdle,
		IdleTimeout: time.Duration(idleTimeout) * time.Second,
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
		Dial: func() (redis.Conn, error) {
			if 0 != len(db) {
				return dialWithDB("tcp", host, password, db)
			} else {
				return dial("tcp", host, password)
			}
		},
	}
	return
}

func dial(network, address, password string) (conn redis.Conn, err error) {
	conn, err = redis.Dial(network, address)
	if nil != err {
		return
	}
	if 0 != len(password) {
		if _, err = conn.Do("AUTH", password); nil != err {
			conn.Close()
			return
		}
	}
	return
}

func dialWithDB(network, address, password, DB string) (conn redis.Conn, err error) {
	conn, err = dial(network, address, password)
	if nil != err {
		return
	}
	if _, err = conn.Do("SELECT", DB); nil != err {
		conn.Close()
		return
	}
	return
}