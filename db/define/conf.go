package define

import "time"

// Database 数据库配置
type Database struct {
	Drive           string                 //	驱动
	User            string                 //	用户名
	Pass            string                 //	密码
	Dbname          string                 // 	库名
	Network         string                 //	网络
	Server          string                 //	地址
	Port            int                    //	端口
	ConnMaxLifeTime time.Duration          //	最大连接周期
	MaxOpenConn     int                    //	设置最大连接数
	ConnMaxIdleTime time.Duration          //	最大空闲连接时间
	MaxIdleConn     int                    //	最大空闲连接数
	Params          map[string]interface{} //	连接参数
}
