package beego

import (
	"database/sql"
	"github.com/gomodule/redigo/redis"
	"github.com/raozhaofeng/beego/db"
	"github.com/raozhaofeng/beego/db/define"
	"github.com/raozhaofeng/beego/router"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	Conf           *Config     //	配置文件
	Logger         *zap.Logger //	日志对象
	RedisManager   *redis.Pool //	缓存配置
	DbManager      *sql.DB     //	数据库管理
	LocalesManager *Locales    //	本地语言管理
)

type BeeGo struct {
	Router *router.Router
}

// NewBeeGo 创建框架对象
func NewBeeGo(confPath string) *BeeGo {
	conf := ReadConfigFile(confPath)
	Conf = conf
	Logger = InitializationLogger(conf.Logs.OutputPaths, conf.Debug)
	RedisManager = InitializationConnPool(conf.Redis)
	DbManager = db.InitializationDb(conf.Database)

	//	启动服务
	return &BeeGo{
		Router: router.NewRoute(RedisManager),
	}
}

// InitializationLocales 初始化语言
func (c *BeeGo) InitializationLocales(localesList map[int64]map[string]string) *BeeGo {
	rds := RedisManager.Get()
	defer rds.Close()

	InitializationLocales(rds, localesList)
	return c
}

// ReadConfigFile 读取配置文件
func ReadConfigFile(confPath string) *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(confPath)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	return &Config{
		Debug: viper.GetBool("Debug"),
		Database: &define.Database{
			Drive:           viper.GetString("Database.Drive"),
			User:            viper.GetString("Database.User"),
			Pass:            viper.GetString("Database.Pass"),
			Dbname:          viper.GetString("Database.Dbname"),
			Network:         viper.GetString("Database.Network"),
			Server:          viper.GetString("Database.Server"),
			Port:            viper.GetInt("Database.Port"),
			ConnMaxLifeTime: 0,
			MaxOpenConn:     0,
			ConnMaxIdleTime: 0,
			MaxIdleConn:     0,
			Params:          map[string]any{},
		},
		Redis: &Redis{
			Network:         viper.GetString("Redis.Network"),
			Server:          viper.GetString("Redis.Server"),
			Port:            viper.GetInt("Redis.Port"),
			Pass:            viper.GetString("Redis.Pass"),
			Dbname:          viper.GetInt("Redis.Dbname"),
			ConnectTimeout:  30,
			ReadTimeout:     30,
			WriteTimeout:    30,
			MaxOpenConn:     0,
			ConnMaxIdleTime: 30,
			MaxIdleConn:     0,
			Wait:            false,
		},
		Logs: &Logs{
			OutputPaths: viper.GetStringSlice("Logs.OutputPaths"),
		},
	}
}
