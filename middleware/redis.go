package middleware

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/labstack/echo"
	"github.com/shibingli/realclouds_go/utils"
)

//DefaultRedis *
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

//Redis *
type Redis struct {
	RedisPool *redis.Pool
	Mutex     sync.RWMutex
}

//MwRedis Redis middleware
func (r *Redis) MwRedis(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		r.Mutex.Lock()
		defer r.Mutex.Unlock()
		c.Set("redis", r)
		return next(c)
	}
}

//Setex *
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

//Publish *
func (r *Redis) Publish(key string, msg interface{}) (err error) {
	key = strings.TrimSpace(key)
	conn := r.RedisPool.Get()
	defer conn.Close()
	if err = conn.Err(); err != nil {
		return
	}
	_, err = conn.Do("PUBLISH", key, msg)
	if nil != err {
		return
	}
	return
}

//ListenPubSubChannels *
func (r *Redis) ListenPubSubChannels(
	ctx context.Context,
	onStart func() error,
	onMessage func(channel string, data []byte) error,
	onPMessage func(pattern, channel string, data []byte) error,
	channels ...string) (err error) {

	return listenPubSubChannels(ctx, r.RedisPool, onStart, onMessage, onPMessage, channels...)
}

//ListenPubSubChannels *
func ListenPubSubChannels(
	ctx context.Context,
	rPool *redis.Pool,
	onStart func() error,
	onMessage func(channel string, data []byte) error,
	onPMessage func(pattern, channel string, data []byte) error,
	channels ...string) (err error) {

	return listenPubSubChannels(ctx, rPool, onStart, onMessage, onPMessage, channels...)
}

func listenPubSubChannels(
	ctx context.Context,
	rPool *redis.Pool,
	onStart func() error,
	onMessage func(channel string, data []byte) error,
	onPMessage func(pattern, channel string, data []byte) error,
	channels ...string) (err error) {

	conn := rPool.Get()
	defer conn.Close()
	if err = conn.Err(); err != nil {
		return
	}

	psc := redis.PubSubConn{
		Conn: conn,
	}

	if err := psc.Subscribe(redis.Args{}.AddFlat(channels)...); err != nil {
		return err
	}

	if err := psc.PSubscribe(redis.Args{}.AddFlat(channels)...); err != nil {
		return err
	}

	done := make(chan error, 1)

	go func() {
		for {
			switch n := psc.Receive().(type) {
			case error:
				done <- n
				return
			case redis.Message:
				if nil != onMessage {
					if err := onMessage(n.Channel, n.Data); err != nil {
						done <- err
						return
					}
				}
			case redis.PMessage:
				if nil != onPMessage {
					if err := onPMessage(n.Pattern, n.Channel, n.Data); err != nil {
						done <- err
						return
					}
				}
			case redis.Subscription:
				switch n.Count {
				case len(channels):
					if err := onStart(); err != nil {
						done <- err
						return
					}
				case 0:
					done <- nil
					return
				}
			}
		}
	}()

	ticker := time.NewTicker(time.Duration(5 * time.Second))
	defer ticker.Stop()

loop:
	for err == nil {
		select {
		case <-ticker.C:
			if err = psc.Ping(""); err != nil {
				break loop
			}
		case <-ctx.Done():
			break loop
		case err := <-done:
			return err
		}
	}

	if err := psc.Unsubscribe(); nil != err {
		done <- err
	}

	if err := psc.PUnsubscribe(); nil != err {
		done <- err
	}

	return <-done
}

//FlushDB **
func (r *Redis) FlushDB() (err error) {
	conn := r.RedisPool.Get()
	defer conn.Close()
	if err = conn.Err(); err != nil {
		return
	}
	_, err = conn.Do("FLUSHDB")
	if nil != err {
		return
	}
	return
}

//FlushAll *
func (r *Redis) FlushAll() (err error) {
	conn := r.RedisPool.Get()
	defer conn.Close()
	if err = conn.Err(); err != nil {
		return
	}
	_, err = conn.Do("FLUSHALL")
	if nil != err {
		return
	}
	return
}

//Del *
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

//Get *
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

//GetString *
func (r *Redis) GetString(key string) (string, error) {
	return redis.String(r.Get(key))
}

//GetInterface *
func (r *Redis) GetInterface(key string, inf interface{}) (str string, err error) {
	str, err = redis.String(r.Get(key))
	if nil != err {
		return
	}
	err = json.Unmarshal([]byte(str), &inf)
	return
}

//GetStrings *
func (r *Redis) GetStrings(key string) ([]string, error) {
	return redis.Strings(r.Get(key))
}

//GetStringMap *
func (r *Redis) GetStringMap(key string) (map[string]string, error) {
	return redis.StringMap(r.Get(key))
}

//GetInt *
func (r *Redis) GetInt(key string) (int, error) {
	return redis.Int(r.Get(key))
}

//GetInts *
func (r *Redis) GetInts(key string) ([]int, error) {
	return redis.Ints(r.Get(key))
}

//GetIntMap *
func (r *Redis) GetIntMap(key string) (map[string]int, error) {
	return redis.IntMap(r.Get(key))
}

//GetInt64 *
func (r *Redis) GetInt64(key string) (int64, error) {
	return redis.Int64(r.Get(key))
}

//GetInt64Map *
func (r *Redis) GetInt64Map(key string) (map[string]int64, error) {
	return redis.Int64Map(r.Get(key))
}

//GetUint64 *
func (r *Redis) GetUint64(key string) (uint64, error) {
	return redis.Uint64(r.Get(key))
}

//GetFloat64 *
func (r *Redis) GetFloat64(key string) (float64, error) {
	return redis.Float64(r.Get(key))
}

//GetBytes *
func (r *Redis) GetBytes(key string) ([]byte, error) {
	return redis.Bytes(r.Get(key))
}

//GetByteSlices *
func (r *Redis) GetByteSlices(key string) ([][]byte, error) {
	return redis.ByteSlices(r.Get(key))
}

//GetBool *
func (r *Redis) GetBool(key string) (bool, error) {
	return redis.Bool(r.Get(key))
}

//NewRedis *
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
			}

			return dial("tcp", host, password)
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
