package admin

import (
	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/biz/bmodel"
)

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
		if finder == nil || err != nil {
			finder = &bmodel.Weixin_RuleModel{
				RequestType: 1,
				SortID:      99,
				Name:        "文本回复",
			}
		}
		c.Data["mdl"] = finder
	}
	c.display()
}
