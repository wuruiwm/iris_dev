package router

import (
	"github.com/kataras/iris/v12"
	"iris_dev/app/controller/admin"
)

//后台路由
func adminRouter(r *iris.Application)*iris.Application{
	adminGroup := r.Party("/admin")
	{
		article := adminGroup.Party("/article")
		{
			article.Get("/",admin.ArticleList).Name = "文章列表"
			article.Post("/", admin.ArticleCreate).Name = "文章创建"
			article.Put("/:id", admin.ArticleUpdate).Name = "文章编辑"
			article.Delete("/:id", admin.ArticleDelete).Name = "文章删除"
			article.Get("/:id", admin.ArticleDetail).Name = "文章详情"
		}
	}
	return r
}