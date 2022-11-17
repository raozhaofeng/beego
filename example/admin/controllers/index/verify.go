package index

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/raozhaofeng/beego"
	"github.com/raozhaofeng/beego/router"
	"net/http"
)

func Verify(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	claims := router.TokenManager.GetContextClaims(r)

	rds := beego.RedisManager.Get()
	defer rds.Close()

	fmt.Println(beego.LocalesManager.GetAdminLocales(rds, claims.AdminId, "zh-CN", "name"))

	router.SuccessJSON(w, claims)
}
