package admin

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/biz/bizmodel"
	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/controllers/admin/vmodel"
	"haedu.gov.cn/tools/xgeneric"
	"haedu.gov.cn/tools/xjson"
)

// ContentFindSubscribeOrDefault 获取关注回复与默认回复
// @router /admin/weixin/ContentFindSubscribeOrDefault [get]
func (c *WeixinController) ContentFindSubscribeOrDefault() {
	accountId, _ := c.GetInt64("account_id")
	requestType, _ := c.GetInt32("request_type")

	outModel, err := biz.NewWeixinRequest().ContentFindSubscribeOrDefault(accountId, requestType)
	if err != nil {
		logs.Error("ContentFindSubscribeOrDefault", err.Error())
		outModel = &bizmodel.Weixin_ContentSubscribeOrDefaultModel{
			AccountID:   accountId,
			RequestType: requestType,
			TextReply:   &model.WeixinRequestContent{},
			ImageReply:  []*model.WeixinRequestContent{},
			SoundReply:  &model.WeixinRequestContent{},
		}
	}
	c.JSONSuccess("获取数据", outModel)
}

// ContentSaveSubscribeOrDefault 保存关注回复与默认回复
// @router /admin/weixin/ContentSaveSubscribeOrDefault [post]
func (c *WeixinController) ContentSaveSubscribeOrDefault() {
	mdl := bizmodel.Weixin_ContentSubscribeOrDefaultModel{}
	if err := xjson.Unmarshal(c.Ctx.Input.RequestBody, &mdl); err != nil {
		logs.Error("ContentSubscribeOrDefault", err.Error())
		c.JSONError(err.Error())
	}
	title := xgeneric.IFF(mdl.RequestType == 6, "关注回复", "默认回复")
	mdl.TextReply.Title = title
	if err := biz.NewWeixinRequest().ContentSaveSubscribeOrDefault(&mdl); err != nil {
		logs.Error("ContentSaveSubscribeOrDefault", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("保存成功", nil)
}

func (c *WeixinController) EmptyImageReply() {
	images := []*model.WeixinRequestContent{}
	c.JSONPageSuccess(images, int64(len(images)))
}

func (c *WeixinController) RulePaginate() {
	page, limit := c.GetPagingParameters()
	accountId, _ := c.GetInt64("account_id")
	requestType, _ := c.GetInt32("request_type")
	list, total, err := biz.NewWeixinRequest().RulePaginate(page, limit, accountId, requestType)
	if err != nil {
		logs.Error("ContentPaginate", err.Error())
		c.JSONPageError(err.Error(), list, total)
	}
	c.JSONPageSuccess(list, total)
}

func (c *WeixinController) RulePictureFind() {
	ruleId, _ := c.GetInt64("rule_id")
	list, err := biz.NewWeixinRequest().RulePictureFind(ruleId)
	if err != nil {
		logs.Error("RulePictureFind", err.Error())
		c.JSONError(err.Error())
	}
	c.JSONSuccess("获取数据", list)
}

func (c *WeixinController) RulePaginateCount() {
	page, limit := c.GetPagingParameters()
	accountId, _ := c.GetInt64("account_id")
	requestType, _ := c.GetInt32("request_type")
	list, total, err := biz.NewWeixinRequest().RulePaginateCount(page, limit, accountId, requestType)
	if err != nil {
		logs.Error("RulePaginateCount", err.Error())
		c.JSONPageError(err.Error(), list, total)
	}
	c.JSONPageSuccess(list, total)
}

func (c *WeixinController) RuleSaveSortId() {
	mdls := []vmodel.Rule_SaveSortIdModel{}
	data := c.Ctx.Input.RequestBody
	fmt.Println("SiteSaveSortId", string(data))
	if err := xjson.Unmarshal(c.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("SiteSaveSortId", err.Error())
		c.JSONError(err.Error())
	}
	for _, mdl := range mdls {
		if err := biz.NewWeixinRequest().RuleSaveSortId(mdl.RuleId, mdl.SortId); err != nil {
			logs.Error("RuleSaveSortId", err.Error())
			c.JSONError(err.Error())
			return
		}
	}
	c.JSONSuccess("保存成功", nil)
}

func (c *WeixinController) RuleDestory() {
	ruleId, _ := c.GetInt64("rule_id")
	if err := biz.NewWeixinRequest().RuleDestory(ruleId); err != nil {
		logs.Error("RuleDestory", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("删除成功", nil)
}
