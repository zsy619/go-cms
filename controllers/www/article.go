package www

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/tools/xstring"
)

// ArticleController 文章控制器
type ArticleController struct{ BaseController }

// Index 文章首页
// @router /:flag:string/:name:string/:call_index:string/:article_id:int64 [get]
func (ctrl *ArticleController) Index() {
	flag := ctrl.Ctx.Input.Param(":flag")             // 频道名称
	name := ctrl.Ctx.Input.Param(":name")             // 频道名称
	call_index := ctrl.Ctx.Input.Param(":call_index") // 栏目别名
	article_id := ctrl.Ctx.Input.Param(":article_id") // 文章ID
	fmt.Println("flag:", flag)
	fmt.Println("name:", name)
	fmt.Println("call_index:", call_index)
	fmt.Println("article_id:", article_id)
}

// Detail 文章详情
// @router /article/:call_index:string/:article_id:int64 [get]
func (ctrl *ArticleController) Detail() {
	call_index := ctrl.Ctx.Input.Param(":call_index") // 栏目别名
	fmt.Println("call_index:", call_index)
	particle_id := ctrl.Ctx.Input.Param(":article_id") // 获取路由参数
	fmt.Println("particle_id:", particle_id)
	article_id := xstring.ToInt64(particle_id)
	if article_id <= 0 {
		ctrl.Abort("404")
	}
	// 获取文章详情
	articleModel, albumModel, attachModel, propertyModel, err := ctrl.ArticleFind("", article_id)
	if err != nil {
		logs.Error("Detail:", err.Error())
	}
	ctrl.Data["article_id"] = article_id
	ctrl.Data["title"] = articleModel.Title
	ctrl.Data["article"] = articleModel
	ctrl.Data["album"] = albumModel
	ctrl.Data["attach"] = attachModel
	ctrl.Data["property"] = propertyModel

	if articleModel.Template == "" {
		if articleModel.TmplDtl != "" {
			articleModel.Template = articleModel.TmplDtl
		}
	}
	if articleModel.Template == "" {
		ctrl.TplName = ctrl.GetView(DefaultSite.Template, "article.html")
	} else {
		if xstring.HasSuffix(articleModel.Template, ".html", ".htm", ".tpl") == false {
			articleModel.Template += ".html"
		}
		ctrl.TplName = ctrl.GetView(DefaultSite.Template, articleModel.Template)
	}
}

// Search 文章搜索
// @router /article/search/:keyword [get]
func (ctrl *ArticleController) Search() {
	keyword := ctrl.Ctx.Input.Param(":keyword") // 获取路由参数
	ctrl.Data["keyword"] = keyword
	ctrl.TplName = ctrl.GetView(DefaultSite.Template, "search.html")
}
