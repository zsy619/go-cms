package admin

import (
	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/biz/bizmodel"
)

// 语音回复
func (ctrl *WeixinController) Sound() {
	list, _, _ := biz.NewWeixinAccount().AccountPaginate(1, 9999, "", -1)
	ctrl.Data["accountList"] = list
	ctrl.Data["request_type"] = 3

	ctrl.display()
}

func (ctrl *WeixinController) SoundEdit() {
	{
		list, _, _ := biz.NewWeixinAccount().AccountPaginate(1, 9999, "", -1)
		ctrl.Data["accountList"] = list
		ctrl.Data["request_type"] = 3
	}
	{
		ruleId, _ := ctrl.GetInt64("rule_id")
		ctrl.Data["rule_id"] = ruleId
		finder, err := biz.NewWeixinRequest().RuleFind(ruleId)
		if finder == nil || err != nil {
			finder = &bizmodel.Weixin_RuleModel{
				RequestType: 3,
				SortID:      99,
				Name:        "语音回复",
			}
		}
		ctrl.Data["mdl"] = finder
	}
	ctrl.display()
}
