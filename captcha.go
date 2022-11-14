package beego

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

const RedisName = "captcha_"

type RedisStore struct {
	Expire time.Duration
}

// Set 设置验证码数据
func (rs *RedisStore) Set(id string, digits []byte) {
	rds := RedisManager.Get()
	defer rds.Close()

	_, _ = rds.Do("SETEX", RedisName+id, rs.Expire.Seconds(), digits)
}

// Get 获取验证码数据
func (rs *RedisStore) Get(id string, clear bool) (digits []byte) {
	rds := RedisManager.Get()
	defer rds.Close()

	bs, _ := redis.Bytes(rds.Do("GET", RedisName+id))
	if clear {
		_, _ = rds.Do("DEL", RedisName+id)
	}
	return bs
}
