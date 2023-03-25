package admin

import (
	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/lib"
)

type WeixinController struct{ BaseController }

// 微信设置
func (c *WeixinController) Index() {
	c.Data["website"] = lib.C_LOCAL_DOMAIN_Backslash()
	c.display()
}

// 微信菜单
func (c *WeixinController) Menu() {
	list, _, _ := biz.NewWeixinAccount().AccountPaginate(1, 9999, "", -1)
	c.Data["accountList"] = list
	c.display()
}

// 关注回复
func (c *WeixinController) Subscribe() {
	c.display()
}

// 默认回复
func (c *WeixinController) Default() {
	c.display()
}

// 文本回复
func (c *WeixinController) Text() {
	c.display()
}

// 图文回复
func (c *WeixinController) Picture() {
	c.display()
}

// 语音回复
func (c *WeixinController) Sound() {
	c.display()
}

// 消息记录
func (c *WeixinController) Response() {
	c.display()
}
