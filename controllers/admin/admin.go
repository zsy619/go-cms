package admin

import (
	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/lib"
)

type AdminController struct{ BaseController }

func (c *AdminController) Log() {
	c.display()
}

func (c *AdminController) LogPaginate() {
	page, limit := c.GetPagingParameters()
	userName := c.GetString("userName")
	list, count, _ := biz.NewCmsAdmin().LogPaginate(page, limit, 0, userName)
	c.JSONPaging(lib.CodeSuccess, "", list, count)
}
