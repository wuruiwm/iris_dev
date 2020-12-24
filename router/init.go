package router

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"iris_dev/config"
)

func HttpInit(){
	//实例化iris
	app := iris.New()
	//注册模板路径
	app.RegisterView(iris.HTML("./app/view",".html"))
	//路由设置
	app = router(app)
	//启动服务
	_ = app.Run(iris.Addr(getHttpString()))
}

func getHttpString()string{
	port := config.GetString("server_port")
	return fmt.Sprintf("0.0.0.0:%s",port)
}
