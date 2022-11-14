package beego

import (
	"github.com/raozhaofeng/beego/db/define"
	"time"
)

// Redis Redis配置文件
type Redis struct {
	Network         string        //	网络
	Server          string        //	地址
	Port            int           //	端口
	Pass            string        //	密码
	Dbname          int           // 	库名
	ConnectTimeout  time.Duration //	连接超时时间
	ReadTimeout     time.Duration //	读取超时时间
	WriteTimeout    time.Duration //	写入超时时间
	MaxOpenConn     int           // 	设置最大连接数
	ConnMaxIdleTime time.Duration // 	空闲连接超时
	MaxIdleConn     int           // 	最大空闲连接数
	Wait            bool          // 	如果超过最大连接数是否等待
}

// Logs 日志配置文件
type Logs struct {
	OutputPaths []string //	日志输出路径
}

// Config 配置文件
type Config struct {
	Debug    bool             // 是否调试
	Database *define.Database // 数据库配置
	Redis    *Redis           // 缓存配置文件
	Logs     *Logs            // 日志配置
}
