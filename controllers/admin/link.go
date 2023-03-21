package admin

import (
	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/lib"
)

type LinkController struct{ BaseController }

// Index 链接管理
// @router /admin/link/index [get]
func (c *LinkController) Index() {
	list, _, _ := biz.NewCmsLink().CategoryPaginate(1, 99999, -1, -1, "", "")
	c.Data["categoryList"] = list
	c.display()
}

// LinkPaginate 列表
// @router /admin/link/linkpaginate [get]
func (c *LinkController) LinkPaginate() {
	page, limit := c.GetPagingParameters()
	categoryId, _ := c.GetInt64("categoryId")
	title := c.GetString("title")
	callIndex := c.GetString("callIndex")
	list, count, _ := biz.NewCmsLink().LinkPaginate(page, limit, -1, -1, categoryId, title, callIndex)
	c.JSONPaging(lib.CodeSuccess, "", list, count)
}

// Category 链接分类
// @router /admin/link/category [get]
func (c *LinkController) Category() {
	c.display()
}

func (c *LinkController) CategoryEdit() {
	categoryId, _ := c.GetInt64("categoryId")
	mdl, _ := biz.NewCmsLink().CategoryFind(categoryId)
	c.Data["mdl"] = mdl
	c.display()
}

// CategoryPaginate 列表
// @router /admin/link/categorypaginate [get]
func (c *LinkController) CategoryPaginate() {
	page, limit := c.GetPagingParameters()
	title := c.GetString("title")
	callIndex := c.GetString("callIndex")
	list, count, _ := biz.NewCmsLink().CategoryPaginate(page, limit, -1, -1, title, callIndex)
	c.JSONPaging(lib.CodeSuccess, "", list, count)
}
