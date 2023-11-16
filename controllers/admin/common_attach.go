package admin

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/tools/xjson"

	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/controllers/admin/vmodel"
)

// Attach 附件管理
func (ctrl *CommonController) Attach() {
	tableName := ctrl.GetSafeString("tableName")
	recordId, _ := ctrl.GetInt64("recordId")
	typeId, _ := ctrl.GetInt32("typeId")
	if tableName == "" || recordId <= 0 {
		ctrl.JSONError("参数错误")
		return
	}
	ctrl.Data["tableName"] = tableName
	ctrl.Data["recordId"] = recordId
	ctrl.Data["typeId"] = typeId
	ctrl.display()
}

func (ctrl *CommonController) AttachPaginate() {
	tableName := ctrl.GetSafeString("tableName")
	recordId, _ := ctrl.GetInt64("recordId")
	typeId, _ := ctrl.GetInt32("typeId")
	list, count, err := biz.NewCmsAttach().AttachPaginate(1, 99999, tableName, recordId, typeId, GlobalAdminId, GlobalRoleType)
	if err != nil {
		logs.Error("AttachPaginate", err.Error())
	}
	ctrl.JSONPageSuccess(list, count)
}

func (ctrl *CommonController) AttachSave() {
	mdl := model.CmsAttach{}
	if err := ctrl.ParseForm(&mdl); err != nil {
		logs.Error("AttachSave", err.Error())
		ctrl.JSONError(err.Error())
	}
	mdl.CreateID = int32(GlobalAdminId)
	if err := biz.NewCmsAttach().AttachSave(&mdl); err != nil {
		logs.Error("AttachSave", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("保存成功", nil)
}

func (ctrl *CommonController) AttachSaveShow() {
	mdls := vmodel.Article_AttachSaveShowModel{}
	data := ctrl.Ctx.Input.RequestBody
	fmt.Println("AttachSaveShow", string(data))
	if err := xjson.Unmarshal(ctrl.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("AttachSaveShow", err.Error())
		ctrl.JSONError(err.Error())
	}
	for _, mdl := range mdls.AttachIds {
		if err := biz.NewCmsAttach().AttachSaveShow(mdl, mdls.Show); err != nil {
			logs.Error("AttachSaveShow", err.Error())
			ctrl.JSONError(err.Error())
			return
		}
	}
	ctrl.JSONSuccess("保存成功", nil)
}

func (ctrl *CommonController) AttachSaveBatch() {
	mdls := []vmodel.Article_AttachSaveBatchdModel{}
	data := ctrl.Ctx.Input.RequestBody
	fmt.Println("AttachSaveBatch", string(data))
	if err := xjson.Unmarshal(ctrl.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("AttachSaveBatch", err.Error())
		ctrl.JSONError(err.Error())
	}
	service := biz.NewCmsAttach()
	for _, mdl := range mdls {
		if err := service.AttachSaveInfo(mdl.AttachID, mdl.Title, mdl.Point, mdl.Click, mdl.SortID, mdl.Remark); err != nil {
			logs.Error("AttachSaveBatch", err.Error())
			ctrl.JSONError(err.Error())
			return
		}
	}
	ctrl.JSONSuccess("保存成功", nil)
}

func (ctrl *CommonController) AttachDestory() {
	attachId, _ := ctrl.GetInt64("attachId")
	if err := biz.NewCmsAttach().AttachDestory(attachId); err != nil {
		logs.Error("AttachDestory", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("删除成功", nil)
}
