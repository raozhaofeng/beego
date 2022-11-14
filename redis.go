package beego

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
)

// InitializationConnPool 初始化缓存
func InitializationConnPool(conf *Redis) *redis.Pool {
	var rds *redis.Pool
	//	如果没有设置,那么直接返回nil
	if conf == nil || conf.Server == "" {
		return nil
	}

	//	设置默认最大连接数
	maxActive := 1000
	if conf.MaxOpenConn > 0 {
		maxActive = conf.MaxOpenConn
	}
	//	设置最大空闲数
	maxIdle := 100
	if conf.MaxIdleConn > 0 {
		maxIdle = conf.MaxIdleConn
	}
	//	设置最大空闲超时时间
	var idleTimeout time.Duration = 30
	if conf.ConnMaxIdleTime > 0 {
		idleTimeout = conf.ConnMaxIdleTime
	}
	//	连接超时时间
	if conf.ConnectTimeout == 0 {
		conf.ConnectTimeout = 30
	}
	//	读取超时时间
	if conf.ReadTimeout == 0 {
		conf.ReadTimeout = 30
	}
	//	写入超时时间
	if conf.WriteTimeout == 0 {
		conf.WriteTimeout = 30
	}

	rds = &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: idleTimeout * time.Second,
		Wait:        conf.Wait,
		Dial: func() (redis.Conn, error) {
			host := fmt.Sprintf("%s:%d", conf.Server, conf.Port)
			conn, err := redis.Dial(
				conf.Network,
				host,
				redis.DialPassword(conf.Pass),
				redis.DialDatabase(conf.Dbname),
				redis.DialConnectTimeout(conf.ConnectTimeout*time.Second),
				redis.DialReadTimeout(conf.ReadTimeout*time.Second),
				redis.DialWriteTimeout(conf.WriteTimeout*time.Second),
			)
			if err != nil {
				return nil, err
			}
			return conn, nil
		},
	}
	return rds
}
