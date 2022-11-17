package beego

import (
	"github.com/gomodule/redigo/redis"
	"strconv"
)

const (
	LocalesRedisName = "_locales"
)

type Locales struct {
}

// InitializationLocales 初始化语言
func InitializationLocales(rds redis.Conn, localesList map[int64]map[string]map[string]string) {
	LocalesManager = &Locales{}
	// 初始化语言
	for adminId, aliasLocales := range localesList {
		for alias, locales := range aliasLocales {
			LocalesManager.SetAdminLocalesAll(rds, adminId, alias, locales)
		}
	}
}

// SetAdminLocalesAll 设置管理所有语言
func (c *Locales) SetAdminLocalesAll(rds redis.Conn, adminId int64, alias string, locales map[string]string) {
	for key, val := range locales {
		_, err := rds.Do("HSET", c.LocalesRedisName(adminId, alias), key, val)
		if err != nil {
			panic(err)
		}
	}
}

// SetAdminLocales 设置管理语言
func (c *Locales) SetAdminLocales(rds redis.Conn, adminId int64, alias string, localeKey string, localVal any) {
	_, err := rds.Do("HSET", c.LocalesRedisName(adminId, alias), localeKey, localVal)
	if err != nil {
		panic(err)
	}
}

// GetAdminLocales 获取管理语言值
func (c *Locales) GetAdminLocales(rds redis.Conn, adminId int64, alias string, localeKey string) string {
	data, err := redis.String(rds.Do("HGET", c.LocalesRedisName(adminId, alias), localeKey))
	if err != nil {
		return localeKey
	}
	return data
}

// GetAdminLocalesAll 获取管理所有语言
func (c *Locales) GetAdminLocalesAll(rds redis.Conn, adminId int64, alias string) map[string]string {
	locales, _ := redis.Strings(rds.Do("HGETALL", c.LocalesRedisName(adminId, alias)))

	data := map[string]string{}
	localeKey := ""
	for i, locale := range locales {
		if i%2 == 0 {
			localeKey = locale
		} else {
			data[localeKey] = locale
		}
	}
	return data
}

// LocalesRedisName 本地语言缓存名称
func (c *Locales) LocalesRedisName(adminId int64, alias string) string {
	adminIdStr := strconv.FormatInt(adminId, 10)
	return LocalesRedisName + "_" + adminIdStr + "_" + alias
}
