package router

import "github.com/kataras/iris/v12"

//设置路由
func router(r *iris.Application)*iris.Application{
	r = adminRouter(r)
	return r
}