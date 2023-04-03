package www

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/tools/xstring"
)

// ArticleController 文章控制器
type ArticleController struct{ BaseController }

// Detail 文章详情
// @router /article/detail/:article_id:int64 [get]
func (c *ArticleController) Detail() {
	particle_id := c.Ctx.Input.Param(":article_id") // 获取路由参数
	fmt.Println("particle_id:", particle_id)
	article_id := xstring.ToInt64(particle_id)
	if article_id <= 0 {
		c.Abort("404")
	}
	// 获取文章详情
	article, album, attach, err := c.ArticleGet("", article_id)
	if err != nil {
		logs.Error("Detail:", err.Error())
	}
	c.Data["article_id"] = article_id
	c.Data["title"] = article.Title
	c.Data["article"] = article
	c.Data["album"] = album
	c.Data["attach"] = attach
	c.displayNoLayout()
}

// Teacher 教师风采
// @router /article/teacher/:article_id:int64 [get]
func (c *ArticleController) Teacher() {
	particle_id := c.Ctx.Input.Param(":article_id") // 获取路由参数
	fmt.Println("particle_id:", particle_id)
	article_id := xstring.ToInt64(particle_id)
	if article_id <= 0 {
		c.Abort("404")
	}
	// 获取文章详情
	article, album, attach, err := c.ArticleGet("", article_id)
	if err != nil {
		logs.Error("Teacher:", err.Error())
	}
	c.Data["article_id"] = article_id
	c.Data["title"] = article.Title
	c.Data["article"] = article
	c.Data["album"] = album
	c.Data["attach"] = attach
	c.displayNoLayout()
}

// Search 文章搜索
// @router /article/search/:keyword [get]
func (c *ArticleController) Search() {
	keyword := c.Ctx.Input.Param(":keyword") // 获取路由参数
	c.Data["keyword"] = keyword
	c.displayNoLayout()
}
