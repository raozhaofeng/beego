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
func InitializationLocales(rds redis.Conn, localesList map[int64]map[string]string) {
	LocalesManager = &Locales{}
	// 初始化语言
	for adminId, locales := range localesList {
		LocalesManager.SetAdminLocalesAll(rds, adminId, locales)
	}
}

// SetAdminLocalesAll 设置管理所有语言
func (c *Locales) SetAdminLocalesAll(rds redis.Conn, adminId int64, locales map[string]string) {
	for key, val := range locales {
		_, _ = rds.Do("HSET", c.LocalesRedisName(adminId), key, val)
	}
}

// GetAdminLocales 获取管理语言值
func (c *Locales) GetAdminLocales(rds redis.Conn, adminId int64, localeKey string) string {
	data, err := redis.String(rds.Do("HGET", c.LocalesRedisName(adminId), localeKey))
	if err != nil {
		return localeKey
	}
	return data
}

// GetAdminLocalesAll 获取管理所有语言
func (c *Locales) GetAdminLocalesAll(rds redis.Conn, adminId int64) map[string]string {
	locales, _ := redis.Strings(rds.Do("HGETALL", c.LocalesRedisName(adminId)))

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
func (c *Locales) LocalesRedisName(adminId int64) string {
	adminIdStr := strconv.FormatInt(adminId, 10)
	return LocalesRedisName + "_" + adminIdStr
}
