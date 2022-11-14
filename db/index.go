package db

import (
	"database/sql"
	"github.com/raozhaofeng/beego/db/define"
	"github.com/raozhaofeng/beego/db/mysql"
)

// InitializationDb 初始化数据库
func InitializationDb(conf *define.Database) *sql.DB {
	return mysql.InitializationDb(conf)
}
