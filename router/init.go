package router

import (
	"github.com/kataras/iris/v12"

)

func HttpInit(){
	app := iris.New()
	app.Get("/", func(ctx iris.Context){
		ctx.WriteString("pong111211")
	})
	_ = app.Run(iris.Addr(":8088"))
}
