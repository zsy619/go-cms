package admin

import (
	"github.com/beego/beego/v2/core/logs"

	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/biz/bizmodel"
	"haedu.gov.cn/cms/app/lib"
)

type WeixinController struct{ BaseController }

// 微信设置
func (ctrl *WeixinController) Index() {
	ctrl.Data["website"] = lib.C_LOCAL_DOMAIN_Backslash()
	ctrl.display()
}

// 微信菜单
func (ctrl *WeixinController) Menu() {
	list, _, _ := biz.NewWeixinAccount().AccountPaginate(1, 9999, "", -1)
	ctrl.Data["accountList"] = list
	ctrl.display()
}

// 关注回复
func (ctrl *WeixinController) Subscribe() {
	list, _, _ := biz.NewWeixinAccount().AccountPaginate(1, 9999, "", -1)
	ctrl.Data["accountList"] = list
	ctrl.Data["request_type"] = 6
	ctrl.display()
}

// 默认回复
func (ctrl *WeixinController) Default() {
	list, _, _ := biz.NewWeixinAccount().AccountPaginate(1, 9999, "", -1)
	ctrl.Data["accountList"] = list
	ctrl.Data["request_type"] = 0
	ctrl.display()
}

// RuleSave 保存规则
func (ctrl *WeixinController) RuleSave() {
	mdl := bizmodel.Weixin_RuleModel{}
	if err := ctrl.ParseForm(&mdl); err != nil {
		logs.Error("RuleSave", err.Error())
		ctrl.JSONError(err.Error())
	}
	err := biz.NewWeixinRequest().RuleSave(&mdl)
	if err != nil {
		logs.Error("RuleSave", err.Error())
		ctrl.JSONError(err.Error())
	}
	ctrl.JSONSuccess("保存成功", nil)
}
