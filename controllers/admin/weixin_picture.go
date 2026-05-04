package admin

import (
	"github.com/beego/beego/v2/core/logs"
	"github.com/zsy619/tools/xjson"

	"haedu.gov.cn/cms/app/cms/service"
	service_model "haedu.gov.cn/cms/app/cms/service/model"
)

// 图文回复
func (ctrl *WeixinController) Picture() {
	list, _, _ := service.NewWeixinAccount().AccountPaginate(1, 9999, "", -1)
	ctrl.Data["accountList"] = list
	ctrl.Data["request_type"] = 2
	ctrl.display()
}

// PictureEdit 编辑图文回复
func (ctrl *WeixinController) PictureEdit() {
	{
		list, _, _ := service.NewWeixinAccount().AccountPaginate(1, 9999, "", -1)
		ctrl.Data["accountList"] = list
		ctrl.Data["request_type"] = 2
	}
	{
		ruleId, _ := ctrl.GetInt64("rule_id")
		ctrl.Data["rule_id"] = ruleId
		finder, err := service.NewWeixinRequest().RuleFind(ruleId)
		if finder == nil || err != nil {
			finder = &service_model.Weixin_RuleModel{
				RequestType: 2,
				SortID:      99,
				Name:        "图文回复",
			}
		}
		ctrl.Data["mdl"] = finder
	}
	ctrl.display()
}

// PictureSave 保存图文回复
func (ctrl *WeixinController) PictureSave() {
	mdl := service_model.Weixin_PictureModel{}
	if err := xjson.Unmarshal(ctrl.Ctx.Input.RequestBody, &mdl); err != nil {
		logs.Error("PictureSave", err.Error())
		ctrl.JSONError(err.Error())
	}
	err := service.NewWeixinRequest().PictureSave(&mdl)
	if err != nil {
		logs.Error("PictureSave", err.Error())
		ctrl.JSONError(err.Error())
	}
	ctrl.JSONSuccess("保存成功", nil)
}
