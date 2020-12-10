package admin

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"iris_dev/app/model"
	"iris_dev/app/response"
)

//文章列表
func ArticleList(c iris.Context){
	//获取参数
	keyword := c.URLParamTrim("keyword")
	categoryId := c.URLParamIntDefault("category_id",0)
	page := c.URLParamIntDefault("page",1)
	limit := c.URLParamIntDefault("limit",15)

	//where条件
	where := make(map[string]interface{})
	if categoryId != 0 {
		where["a.category_id"] = categoryId
	}

	//查询数据
	article,err := model.ArticleList(page,limit,where,keyword)
	if err != nil{
		response.Error(c,"获取失败")
		return
	}
	response.Success(c,"获取成功",article)
	return
}

//文章创建
func ArticleCreate(c iris.Context){
	defer response.RecoverError(c,"创建失败")
	//获取参数
	title := c.FormValue("title")
	content := c.FormValue("content")

	//参数校验
	if title == ""{
		panic("请输入标题")
	}
	if content == ""{
		panic("请输入内容")
	}

	err := model.ArticleCreate(title,content)
	if err != nil{
		panic("创建失败")
	}
	response.Success(c,"创建成功",nil)
}

//文章修改
func ArticleUpdate(c iris.Context){
	defer response.RecoverError(c,"修改失败")
	//获取参数
	id := c.URLParamIntDefault("id",0)
	title := c.FormValue("title")
	content := c.FormValue("content")

	//参数校验
	if id == 0{
		panic("id不能为空")
	}
	if title == ""{
		panic("请输入标题")
	}
	if content == ""{
		panic("请输入内容")
	}

	err := model.ArticleUpdate(id,title,content)
	if err != nil{
		panic("修改失败")
	}
	response.Success(c,"修改成功",nil)
}

//文章删除
func ArticleDelete(c iris.Context){
	defer response.RecoverError(c,"删除失败")
	//获取参数
	id := c.URLParamIntDefault("id",0)

	//参数校验
	if id == 0{
		panic("id不能为空")
	}

	err := model.ArticleDelete(id)
	if err != nil {
		panic("删除失败")
	}
	response.Success(c,"删除成功",nil)
}

//文章详情
func ArticleDetail(c iris.Context){
	defer response.RecoverError(c,"数据不存在,请重试")
	//获取参数
	id := c.URLParamIntDefault("id",0)
	fmt.Println(id)

	//参数校验
	if id == 0{
		panic("id不能为空")
	}

	article,err := model.ArticleDetail(id)
	if err != nil {
		panic("数据不存在")
	}
	response.Success(c,"获取成功",article)
}