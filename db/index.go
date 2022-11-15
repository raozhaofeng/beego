package db

import (
	"database/sql"
	"github.com/raozhaofeng/beego/db/define"
	"github.com/raozhaofeng/beego/db/mysql"
)

// Manager 数据管理
var Manager *Pool

// Pool 数据池
type Pool struct {
	db *sql.DB
}

// NewInterfaceDb 创建新的接口Db
func (c *Pool) NewInterfaceDb(tx *sql.Tx) define.Db {
	if tx == nil {
		return mysql.NewDb(tx, c.db)
	}
	return mysql.NewDb(tx, nil)
}

// GetDb 获取Db
func (c *Pool) GetDb() *sql.DB {
	return c.db
}

// GetTx 获取Tx
func (c *Pool) GetTx() *sql.Tx {
	tx, _ := c.db.Begin()
	return tx
}

// InitializationDb 初始化数据库
func InitializationDb(conf *define.Database) {
	Manager = &Pool{mysql.InitializationDb(conf)}
}
