package admin

import (
	"github.com/kataras/iris/v12"
)

func ChatIndex(c iris.Context){
	_ = c.View("admin/chat/index.html")
}