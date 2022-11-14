package beego

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitializationLogger 初始化日志
func InitializationLogger(outputPaths []string, debug bool) *zap.Logger {
	//	日志登陆,哪些是需要记录日志的
	level := zap.NewAtomicLevelAt(zap.DebugLevel)
	if !debug {
		level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	}

	if len(outputPaths) == 0 {
		outputPaths = []string{"stdout"}
	}

	var zapConf = zap.Config{
		Level:             level,  //	日志级别
		Development:       debug,  //	是否是开发环境。如果是开发模式，对DPanicLevel进行堆栈跟踪
		DisableCaller:     !debug, //	不显示调用函数的文件名称和行号。默认情况下，所有日志都显示。
		DisableStacktrace: !debug, //	是否禁用堆栈跟踪捕获。默认对Warn级别以上和生产error级别以上的进行堆栈跟踪。
		Sampling:          nil,    //	抽样策略。设置为nil禁用采样。
		Encoding:          "json", //	编码方式，支持json, console
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "msg",   //	输入信息的key名
			LevelKey:       "level", //	输出日志级别的key名
			TimeKey:        "time",  //	输出时间的key名
			NameKey:        "name",
			CallerKey:      "caller",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,     //	每行的分隔符。"\\n"
			EncodeLevel:    zapcore.LowercaseLevelEncoder, //	将日志级别字符串转化为小写
			EncodeTime:     zapcore.ISO8601TimeEncoder,    //	输出的时间格式
			EncodeDuration: zapcore.StringDurationEncoder, //	执行消耗的时间转化成浮点型的秒
			EncodeCaller:   zapcore.ShortCallerEncoder,    //	以包/文件:行号 格式化调用堆栈
			EncodeName:     zapcore.FullNameEncoder,       //	可选值。
		},
		OutputPaths:      outputPaths,              //	可以配置多个输出路径，路径可以是文件路径和stdout（标准输出）
		ErrorOutputPaths: []string{"stderr"},       //	错误输出路径（日志内部错误）
		InitialFields:    map[string]interface{}{}, //	每条日志中都会输出这些值
	}
	logger, _ := zapConf.Build()

	return logger
}
