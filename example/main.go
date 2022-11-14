package main

import (
	"fmt"
	"github.com/raozhaofeng/beego"
	"github.com/raozhaofeng/beego/example/admin"
	"github.com/raozhaofeng/beego/router"
	"net/http"
)

func main() {
	app := beego.NewBeeGo("./").InitializationLocales(map[int64]map[string]string{
		1: {"name": "名称"},
	})

	//	日志回调方法
	accessLogsFunc := func(handle *router.Handle, r *http.Request, claims *router.Claims) {
		fmt.Println("访问记录", handle.Name)
	}

	adminTokenParams := map[string]*router.TokenParams{
		"8888": {Key: "123456", Only: true, Expire: 3600},
	}

	adminRolesRouter := map[int64][]string{
		1: {"/admin/token/verify"},
	}

	//	启动监听
	_ = app.Router.ServeFiles("assets").
		SetRouteHandle(admin.Router()).
		InitializationToken(adminTokenParams, adminRolesRouter).
		SetAccessLogsFunc(accessLogsFunc).
		ListenAndServe("0.0.0.0:9090")
}
