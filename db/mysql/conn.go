package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/raozhaofeng/beego/db/define"
	"strings"
	"time"
)

// InitializationDb 初始化Db
func InitializationDb(conf *define.Database) *sql.DB {
	connStr := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", conf.User, conf.Pass, conf.Network, conf.Server, conf.Port, conf.Dbname)

	//	连接附加参数
	var confParamsList []string
	for k, v := range conf.Params {
		confParamsList = append(confParamsList, fmt.Sprintf("%v=%v", k, v))
	}

	if len(confParamsList) > 0 {
		connStr = connStr + "?" + strings.Join(confParamsList, "&")
	}
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		panic(err)
	}

	//	设置连接周期3600s, conn 连接数30 空闲连接时间3600s 空闲连接数10
	if conf.ConnMaxLifeTime == 0 {
		conf.ConnMaxLifeTime = 60 * 60
	}
	if conf.MaxOpenConn == 0 {
		conf.MaxOpenConn = 30
	}
	if conf.ConnMaxIdleTime == 0 {
		conf.ConnMaxIdleTime = 60 * 60
	}
	if conf.MaxIdleConn == 0 {
		conf.ConnMaxIdleTime = 10
	}

	db.SetConnMaxLifetime(conf.ConnMaxLifeTime * time.Second)
	db.SetMaxOpenConns(conf.MaxOpenConn)
	db.SetConnMaxIdleTime(conf.ConnMaxIdleTime * time.Second)
	db.SetMaxIdleConns(conf.MaxIdleConn)

	return db
}
