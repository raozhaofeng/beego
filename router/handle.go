package router

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// ControllerFunc 控制器方法名
type ControllerFunc func(w http.ResponseWriter, r *http.Request, p httprouter.Params)

type Handle struct {
	Name       string         //	名称
	Module     string         //	模块
	Method     string         // 	类型
	Route      string         // 	路由
	Controller ControllerFunc //	控制器
	RouteAuth  bool           //	路由验证
	TokenAuth  bool           //	Token验证
}

// NewHandle 创建路由
func NewHandle(module, name, method, route string, controllerFunc ControllerFunc, routeAuth, tokenAuth bool) *Handle {
	return &Handle{
		Module: module, Name: name, Method: method, Route: route, Controller: controllerFunc, RouteAuth: routeAuth, TokenAuth: tokenAuth,
	}
}

// NewHomeHandle 前台路由
func NewHomeHandle(module, name, method, route string, controllerFunc ControllerFunc, tokenAuth bool) *Handle {
	return &Handle{
		Module: module, Name: name, Method: method, Route: route, Controller: controllerFunc, RouteAuth: false, TokenAuth: tokenAuth,
	}
}

// SetRouteHandle 设置路由函数
func (c *Router) SetRouteHandle(routeHandle []*Handle) *Router {
	for _, handle := range routeHandle {
		handleController := c.middlewareToken(&Handle{
			Module:    handle.Module,
			Name:      handle.Name,
			Method:    handle.Method,
			Route:     handle.Route,
			RouteAuth: handle.RouteAuth,
			TokenAuth: handle.TokenAuth,
		}, httprouter.Handle(handle.Controller))

		//	添加路由函数
		c.httpRouter.Handle(handle.Method, handle.Route, handleController)
	}
	return c
}

// middlewareToken Token中间件
func (c *Router) middlewareToken(handle *Handle, next httprouter.Handle) httprouter.Handle {
	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		crossDomainRequest(writer, request)
		var claims *Claims
		if handle.TokenAuth {
			rds := c.RedisManager.Get()
			defer rds.Close()

			//	验证Token
			claims = TokenManager.Verify(rds, request)
			if claims == nil {
				c.StatusUnauthorized(writer)
				return
			}

			//	如果需要验证路由
			if handle.RouteAuth && !TokenManager.AuthRouter(rds, claims.AdminId, handle.Route) {
				c.StatusUnauthorized(writer)
				return
			}

			//	复制Token信息传递
			ctx := context.WithValue(request.Context(), ClaimsKey, claims)
			request = request.WithContext(ctx)
		}

		//	正常返回
		c.AccessLogsFunc(handle, request, claims)
		next(writer, request, params)
	}
}

// StatusUnauthorized 返回没有权限
func (c *Router) StatusUnauthorized(writer http.ResponseWriter) {
	writer.WriteHeader(http.StatusUnauthorized)
	_, _ = writer.Write([]byte("401 Unauthorized"))
}
