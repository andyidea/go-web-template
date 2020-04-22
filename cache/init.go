package cache

import (
	"github.com/gomodule/redigo/redis"
	"project/conf"
	"time"
)

var prefix string
var pool *redis.Pool

func Init() {
	prefix = conf.Config.Redis.Prefix
	pool = &redis.Pool{
		MaxIdle:     conf.Config.Redis.MaxConn,
		MaxActive:   conf.Config.Redis.MaxConn,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", conf.Config.Redis.Host+":"+conf.Config.Redis.Port,
				redis.DialPassword(conf.Config.Redis.Password),
				redis.DialDatabase(conf.Config.Redis.DB))
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func Prefix() string {
	return prefix
}

func SetString(key, val string) error {
	conn := pool.Get()
	_, err := conn.Do("SET", key, val)
	if err != nil {
		println(err.Error())
	}
	defer conn.Close()
	return err
}

func GetString(key string) (string, error) {
	conn := pool.Get()
	defer conn.Close()
	val, err := redis.String(conn.Do("GET", key))
	if err != nil {
		println(err.Error())
	}
	return val, err
}

func KeyExists(key string) (bool, error) {
	conn := pool.Get()
	defer conn.Close()
	exist, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		println(err.Error())
	}
	return exist, err
}

func Delete(keyName string) error {
	conn := pool.Get()
	defer conn.Close()
	_, err := conn.Do("DEL", keyName)
	if err != nil {
		println(err.Error())
	}
	return err
}

func Rename(oldName string, newName string) error {
	conn := pool.Get()
	defer conn.Close()
	_, err := conn.Do("RENAME", oldName, newName)
	if err != nil {
		println(err.Error())
	}
	return err
}
