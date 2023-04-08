package admin

import (
	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/dal/model"
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
		outModel = &biz.ContentSubscribeOrDefault{
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
	mdl := biz.ContentSubscribeOrDefault{}
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
