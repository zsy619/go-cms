package admin

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/tools/xgeneric"
	"haedu.gov.cn/tools/xjson"

	"haedu.gov.cn/cms/app/cms/domain"
	"haedu.gov.cn/cms/app/cms/service"
	service_model "haedu.gov.cn/cms/app/cms/service/model"
	"haedu.gov.cn/cms/controllers/admin/vmodel"
)

// ContentFindSubscribeOrDefault 获取关注回复与默认回复
// @router /admin/weixin/ContentFindSubscribeOrDefault [get]
func (ctrl *WeixinController) ContentFindSubscribeOrDefault() {
	accountId, _ := ctrl.GetInt64("account_id")
	requestType, _ := ctrl.GetInt32("request_type")

	outModel, err := service.NewWeixinRequest().ContentFindSubscribeOrDefault(accountId, requestType)
	if err != nil {
		logs.Error("ContentFindSubscribeOrDefault", err.Error())
		outModel = &service_model.Weixin_ContentSubscribeOrDefaultModel{
			AccountID:   accountId,
			RequestType: requestType,
			TextReply:   &domain.WeixinRequestContent{},
			ImageReply:  []*domain.WeixinRequestContent{},
			SoundReply:  &domain.WeixinRequestContent{},
		}
	}
	ctrl.JSONSuccess("获取数据", outModel)
}

// ContentSaveSubscribeOrDefault 保存关注回复与默认回复
// @router /admin/weixin/ContentSaveSubscribeOrDefault [post]
func (ctrl *WeixinController) ContentSaveSubscribeOrDefault() {
	mdl := service_model.Weixin_ContentSubscribeOrDefaultModel{}
	if err := xjson.Unmarshal(ctrl.Ctx.Input.RequestBody, &mdl); err != nil {
		logs.Error("ContentSubscribeOrDefault", err.Error())
		ctrl.JSONError(err.Error())
	}
	title := xgeneric.IFF(mdl.RequestType == 6, "关注回复", "默认回复")
	mdl.TextReply.Title = title
	if err := service.NewWeixinRequest().ContentSaveSubscribeOrDefault(&mdl); err != nil {
		logs.Error("ContentSaveSubscribeOrDefault", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("保存成功", nil)
}

func (ctrl *WeixinController) EmptyImageReply() {
	images := []*domain.WeixinRequestContent{}
	ctrl.JSONPageSuccess(images, int64(len(images)))
}

func (ctrl *WeixinController) RulePaginate() {
	page, limit := ctrl.GetPagingParameters()
	accountId, _ := ctrl.GetInt64("account_id")
	requestType, _ := ctrl.GetInt32("request_type")
	list, total, err := service.NewWeixinRequest().RulePaginate(page, limit, accountId, requestType)
	if err != nil {
		logs.Error("ContentPaginate", err.Error())
		ctrl.JSONPageError(err.Error(), list, total)
	}
	ctrl.JSONPageSuccess(list, total)
}

func (ctrl *WeixinController) RulePictureFind() {
	ruleId, _ := ctrl.GetInt64("rule_id")
	list, err := service.NewWeixinRequest().RulePictureFind(ruleId)
	if err != nil {
		logs.Error("RulePictureFind", err.Error())
		ctrl.JSONError(err.Error())
	}
	ctrl.JSONSuccess("获取数据", list)
}

func (ctrl *WeixinController) RulePaginateCount() {
	page, limit := ctrl.GetPagingParameters()
	accountId, _ := ctrl.GetInt64("account_id")
	requestType, _ := ctrl.GetInt32("request_type")
	list, total, err := service.NewWeixinRequest().RulePaginateCount(page, limit, accountId, requestType)
	if err != nil {
		logs.Error("RulePaginateCount", err.Error())
		ctrl.JSONPageError(err.Error(), list, total)
	}
	ctrl.JSONPageSuccess(list, total)
}

func (ctrl *WeixinController) RuleSaveSortId() {
	mdls := []vmodel.Rule_SaveSortIdModel{}
	data := ctrl.Ctx.Input.RequestBody
	fmt.Println("SiteSaveSortId", string(data))
	if err := xjson.Unmarshal(ctrl.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("SiteSaveSortId", err.Error())
		ctrl.JSONError(err.Error())
	}
	for _, mdl := range mdls {
		if err := service.NewWeixinRequest().RuleSaveSortId(mdl.RuleId, mdl.SortId); err != nil {
			logs.Error("RuleSaveSortId", err.Error())
			ctrl.JSONError(err.Error())
			return
		}
	}
	ctrl.JSONSuccess("保存成功", nil)
}

func (ctrl *WeixinController) RuleDestory() {
	ruleId, _ := ctrl.GetInt64("rule_id")
	if err := service.NewWeixinRequest().RuleDestory(ruleId); err != nil {
		logs.Error("RuleDestory", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("删除成功", nil)
}
