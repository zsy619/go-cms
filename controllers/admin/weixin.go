package admin

import (
	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/biz/bizmodel"
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

// RuleSave 保存规则
func (c *WeixinController) RuleSave() {
	mdl := bizmodel.Weixin_RuleModel{}
	if err := c.ParseForm(&mdl); err != nil {
		logs.Error("RuleSave", err.Error())
		c.JSONError(err.Error())
	}
	err := biz.NewWeixinRequest().RuleSave(&mdl)
	if err != nil {
		logs.Error("RuleSave", err.Error())
		c.JSONError(err.Error())
	}
	c.JSONSuccess("保存成功", nil)
}
