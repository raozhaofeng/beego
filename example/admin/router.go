package admin

import (
	"github.com/raozhaofeng/beego/example/admin/controllers/index"
	"github.com/raozhaofeng/beego/router"
)

func Router() []*router.Handle {
	return []*router.Handle{
		{"生成Token", "Admin", "GET", "/admin/token", index.Token, false, false},
		{"验证Token", "Admin", "POST", "/admin/token/verify", index.Verify, true, true},
	}
}
