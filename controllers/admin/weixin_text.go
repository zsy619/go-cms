package admin

import (
	"haedu.gov.cn/cms/app/cms/service"
	service_model "haedu.gov.cn/cms/app/cms/service/model"
)

// 文本回复
func (ctrl *WeixinController) Text() {
	list, _, _ := service.NewWeixinAccount().AccountPaginate(1, 9999, "", -1)
	ctrl.Data["accountList"] = list
	ctrl.Data["request_type"] = 1
	ctrl.display()
}

// 文本回复
func (ctrl *WeixinController) TextEdit() {
	{
		list, _, _ := service.NewWeixinAccount().AccountPaginate(1, 9999, "", -1)
		ctrl.Data["accountList"] = list
		ctrl.Data["request_type"] = 1
	}
	{
		ruleId, _ := ctrl.GetInt64("rule_id")
		ctrl.Data["rule_id"] = ruleId
		finder, err := service.NewWeixinRequest().RuleFind(ruleId)
		if finder == nil || err != nil {
			finder = &service_model.Weixin_RuleModel{
				RequestType: 1,
				SortID:      99,
				Name:        "文本回复",
			}
		}
		ctrl.Data["mdl"] = finder
	}
	ctrl.display()
}
