package main

import (
	"database/sql"
	"fmt"
	"github.com/raozhaofeng/beego"
	"github.com/raozhaofeng/beego/db"
	"github.com/raozhaofeng/beego/db/define"
	"github.com/raozhaofeng/beego/example/admin"
	"github.com/raozhaofeng/beego/router"
	"net/http"
)

type Test struct {
	define.Db
}

func NewTest(tx *sql.Tx) *Test {
	return &Test{db.Manager.NewInterfaceDb(tx).Table("user_assets")}
}

func main() {
	app := beego.NewBeeGo("./").InitializationLocales(map[int64]map[string]map[string]string{
		1: {"zh-CN": {"name": "名称"}},
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

	userAssetsModel2 := NewTest(nil)
	userAssetsModel2.AndWhere("user_id=?", 1).AndWhere("assets_id=?", 1).QueryRow(func(row *sql.Row) {
		fmt.Println(11)
	})

	//_, err := NewTest(nil).Field("admin_id", "user_id", "code", "created_at").
	//	Args(1, 0, "8888", time.Now().Unix()).Insert()
	//if err != nil {
	//	panic(err)
	//}

	//	启动监听
	_ = app.Router.ServeFiles("assets").
		SetRouteHandle(admin.Router()).
		InitializationToken(adminTokenParams, adminRolesRouter).
		SetAccessLogsFunc(accessLogsFunc).
		ListenAndServe("0.0.0.0:9090")
}
