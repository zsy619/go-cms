package admin

import (
	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/biz/bizmodel"
)

// 语音回复
func (c *WeixinController) Sound() {
	list, _, _ := biz.NewWeixinAccount().AccountPaginate(1, 9999, "", -1)
	c.Data["accountList"] = list
	c.Data["request_type"] = 3

	c.display()
}

func (c *WeixinController) SoundEdit() {
	{
		list, _, _ := biz.NewWeixinAccount().AccountPaginate(1, 9999, "", -1)
		c.Data["accountList"] = list
		c.Data["request_type"] = 3
	}
	{
		ruleId, _ := c.GetInt64("rule_id")
		c.Data["rule_id"] = ruleId
		finder, err := biz.NewWeixinRequest().RuleFind(ruleId)
		if finder == nil || err != nil {
			finder = &bizmodel.Weixin_RuleModel{
				RequestType: 3,
				SortID:      99,
				Name:        "语音回复",
			}
		}
		c.Data["mdl"] = finder
	}
	c.display()
}
