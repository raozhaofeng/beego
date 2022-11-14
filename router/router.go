package router

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type Router struct {
	RedisManager   *redis.Pool                                           //	缓存池
	httpRouter     *httprouter.Router                                    // 路由实例
	AccessLogsFunc func(handle *Handle, r *http.Request, claims *Claims) //	访问日志
}

// NewRoute 创建路由
func NewRoute(redisManager *redis.Pool) *Router {
	httpRouter := httprouter.New()
	// 开启跨域请求
	httpRouter.GlobalOPTIONS = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		crossDomainRequest(writer, request)
		writer.WriteHeader(http.StatusNoContent)
	})

	// 全局异常拦截
	httpRouter.PanicHandler = func(writer http.ResponseWriter, request *http.Request, i interface{}) {
		fmt.Println(i)
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = writer.Write([]byte("500 Internal Server Error"))
	}

	return &Router{
		RedisManager: redisManager,
		httpRouter:   httpRouter,
	}
}

// SetAccessLogsFunc 设置访问日志函数
func (c *Router) SetAccessLogsFunc(fun func(handle *Handle, r *http.Request, claims *Claims)) *Router {
	c.AccessLogsFunc = fun
	return c
}

// ServeFiles 开启静态资源
func (c *Router) ServeFiles(filePath string) *Router {
	c.httpRouter.ServeFiles("/"+filePath+"/*filepath", http.Dir("./"+filePath))
	return c
}

// InitializationToken 初始化Token
func (c *Router) InitializationToken(tokenParamsList map[string]*TokenParams, adminRolesRouter map[int64][]string) *Router {
	rds := c.RedisManager.Get()
	defer rds.Close()

	InitializationToken(rds, tokenParamsList, adminRolesRouter)
	return c
}

// ListenAndServe 监听服务
func (c *Router) ListenAndServe(addr string) error {
	InitializeValidator()
	fmt.Println("Listen", addr, "Successful")
	return http.ListenAndServe(addr, c.httpRouter)
}

// crossDomainRequest 设置跨域请求
func crossDomainRequest(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Credentials", "true")
	writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE")
	writer.Header().Set("Access-Control-Allow-Headers", "Token, Token-Key, Content-Type, If-Match, If-Modified-Since, If-None-Match, If-Unmodified-Since, X-Requested-With")
	origin := request.Header.Get("origin")
	if origin == "" {
		origin = "*"
	}
	writer.Header().Set("Access-Control-Allow-Origin", origin)
}
