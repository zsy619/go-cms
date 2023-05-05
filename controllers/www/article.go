package www

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/tools/xstring"
)

// ArticleController 文章控制器
type ArticleController struct{ BaseController }

// // Index 文章首页
// // @router /:flag:string/:name:string/:call_index:string/:article_id:int64 [get]
// func (c *ArticleController) Index() {
// 	flag := c.Ctx.Input.Param(":flag")             // 频道名称
// 	name := c.Ctx.Input.Param(":name")             // 频道名称
// 	call_index := c.Ctx.Input.Param(":call_index") // 栏目别名
// 	article_id := c.Ctx.Input.Param(":article_id") // 文章ID
// 	fmt.Println("flag:", flag)
// 	fmt.Println("name:", name)
// 	fmt.Println("call_index:", call_index)
// 	fmt.Println("article_id:", article_id)
// }

// Detail 文章详情
// @router /article/:call_index:string/:article_id:int64 [get]
func (c *ArticleController) Detail() {
	call_index := c.Ctx.Input.Param(":call_index") // 栏目别名
	fmt.Println("call_index:", call_index)
	particle_id := c.Ctx.Input.Param(":article_id") // 获取路由参数
	fmt.Println("particle_id:", particle_id)
	article_id := xstring.ToInt64(particle_id)
	if article_id <= 0 {
		c.Abort("404")
	}
	// 获取文章详情
	articleModel, albumModel, attachModel, err := c.ArticleFind("", article_id)
	if err != nil {
		logs.Error("Detail:", err.Error())
	}
	c.Data["article_id"] = article_id
	c.Data["title"] = articleModel.Title
	c.Data["article"] = articleModel
	c.Data["album"] = albumModel
	c.Data["attach"] = attachModel

	if articleModel.Template == "" {
		if articleModel.TmplDtl != "" {
			articleModel.Template = articleModel.TmplDtl
		}
	}
	if articleModel.Template == "" {
		c.TplName = c.GetView(DefatulSite.Template, "article.html")
	} else {
		if xstring.HasSuffix(articleModel.Template, ".html", ".htm", ".tpl") == false {
			articleModel.Template += ".html"
		}
		c.TplName = c.GetView(DefatulSite.Template, articleModel.Template)
	}
}

// Search 文章搜索
// @router /article/search/:keyword [get]
func (c *ArticleController) Search() {
	keyword := c.Ctx.Input.Param(":keyword") // 获取路由参数
	c.Data["keyword"] = keyword
	c.TplName = c.GetView(DefatulSite.Template, "search.html")
}
