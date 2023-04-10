package admin

import (
	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/biz/bmodel"
	"haedu.gov.cn/tools/xjson"
)

// 图文回复
func (c *WeixinController) Picture() {
	list, _, _ := biz.NewWeixinAccount().AccountPaginate(1, 9999, "", -1)
	c.Data["accountList"] = list
	c.Data["request_type"] = 2
	c.display()
}

func (c *WeixinController) PictureEdit() {
	{
		list, _, _ := biz.NewWeixinAccount().AccountPaginate(1, 9999, "", -1)
		c.Data["accountList"] = list
		c.Data["request_type"] = 2
	}
	{
		ruleId, _ := c.GetInt64("rule_id")
		c.Data["rule_id"] = ruleId
		finder, err := biz.NewWeixinRequest().RuleFind(ruleId)
		if finder == nil || err != nil {
			finder = &bmodel.Weixin_RuleModel{
				RequestType: 2,
				SortID:      99,
				Name:        "图文回复",
			}
		}
		c.Data["mdl"] = finder
	}
	c.display()
}

// PictureSave 保存图文回复
func (c *WeixinController) PictureSave() {
	mdl := bmodel.Weixin_PictureModel{}
	if err := xjson.Unmarshal(c.Ctx.Input.RequestBody, &mdl); err != nil {
		logs.Error("PictureSave", err.Error())
		c.JSONError(err.Error())
	}
	err := biz.NewWeixinRequest().PictureSave(&mdl)
	if err != nil {
		logs.Error("PictureSave", err.Error())
		c.JSONError(err.Error())
	}
	c.JSONSuccess("保存成功", nil)
}
