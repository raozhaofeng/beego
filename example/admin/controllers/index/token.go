package index

import (
	"github.com/julienschmidt/httprouter"
	"github.com/raozhaofeng/beego"
	"github.com/raozhaofeng/beego/router"
	"net/http"
)

func Token(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	rds := beego.RedisManager.Get()
	defer rds.Close()

	tokenStr := router.TokenManager.Generate(rds, "8888", 1, 1)

	router.SuccessJSON(w, tokenStr)
}
