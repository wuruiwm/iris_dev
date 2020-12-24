package router

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/websocket"
	"iris_dev/app/controller/admin"
	ws "iris_dev/app/websocket"
)

//后台路由
func adminRouter(r *iris.Application)*iris.Application{
	adminGroup := r.Party("/admin")
	adminGroup.Post("/login",admin.Login).Name = "后台登录"
	//adminGroup.Use(middleware.AdminAuth) //后台用户鉴权中间件
	{
		article := adminGroup.Party("/article")
		{
			article.Get("/",admin.ArticleList).Name = "文章列表"
			article.Post("/", admin.ArticleCreate).Name = "文章创建"
			article.Put("/:id", admin.ArticleUpdate).Name = "文章编辑"
			article.Delete("/:id", admin.ArticleDelete).Name = "文章删除"
			article.Get("/:id", admin.ArticleDetail).Name = "文章详情"
		}
		chat := adminGroup.Party("/chat")
		{
			chat.Get("/",admin.ChatIndex).Name = "聊天室模板页"
			chat.Get("/ws/:id",websocket.Handler(ws.Chat())).Name = "聊天室websocket"
		}
	}
	return r
}