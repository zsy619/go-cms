package admin

import (
	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/cms/app/biz"
)

// Response 消息记录
func (c *WeixinController) Response() {
	{
		list, _, _ := biz.NewWeixinAccount().AccountPaginate(1, 9999, "", -1)
		c.Data["accountList"] = list
	}
	c.display()
}

// ContentPaginate 消息记录分页
func (c *WeixinController) ContentPaginate() {
	page, limit := c.GetPagingParameters()
	accountId, _ := c.GetInt64("account_id")
	list, total, err := biz.NewWeixinContent().ContentPaginate(page, limit, accountId)
	if err != nil {
		logs.Error("ContentPaginate", err.Error())
		c.JSONPageError(err.Error(), list, total)
	}
	c.JSONPageSuccess(list, total)
}
