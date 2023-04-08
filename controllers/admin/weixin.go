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
	list, _, _ := biz.NewWeixinAccount().AccountPaginate(1, 9999, "", -1)
	c.Data["accountList"] = list
	c.Data["request_type"] = 6
	c.display()
}

// 默认回复
func (c *WeixinController) Default() {
	list, _, _ := biz.NewWeixinAccount().AccountPaginate(1, 9999, "", -1)
	c.Data["accountList"] = list
	c.Data["request_type"] = 0
	c.display()
}

// 文本回复
func (c *WeixinController) Text() {
	list, _, _ := biz.NewWeixinAccount().AccountPaginate(1, 9999, "", -1)
	c.Data["accountList"] = list
	c.Data["request_type"] = 1
	c.display()
}

// 文本回复
func (c *WeixinController) TextEdit() {
	{
		list, _, _ := biz.NewWeixinAccount().AccountPaginate(1, 9999, "", -1)
		c.Data["accountList"] = list
		c.Data["request_type"] = 1
	}
	{
		ruleId, _ := c.GetInt64("rule_id")
		c.Data["rule_id"] = ruleId
		finder, err := biz.NewWeixinRequest().RuleFind(ruleId)
		if err != nil {
			finder = &biz.RuleFindModel{
				RequestType: 1,
			}
		}
		c.Data["mdl"] = finder
	}
	c.display()
}

// 图文回复
func (c *WeixinController) Picture() {
	list, _, _ := biz.NewWeixinAccount().AccountPaginate(1, 9999, "", -1)
	c.Data["accountList"] = list
	c.Data["request_type"] = 2
	c.display()
}

// 语音回复
func (c *WeixinController) Sound() {
	list, _, _ := biz.NewWeixinAccount().AccountPaginate(1, 9999, "", -1)
	c.Data["accountList"] = list
	c.Data["request_type"] = 3
	c.display()
}

// 消息记录
func (c *WeixinController) Response() {
	c.display()
}
