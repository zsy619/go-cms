package admin

import (
	"github.com/beego/beego/v2/core/logs"
	"github.com/zsy619/tools/xjson"

	"haedu.gov.cn/cms/app/cms/domain"
	"haedu.gov.cn/cms/app/cms/service"
	"haedu.gov.cn/cms/controllers/admin/vmodel"
)

func (ctrl *WeixinController) AccountPaginate() {
	page, limit := ctrl.GetPagingParameters()
	name := ctrl.GetSafeString("name")
	status, _ := ctrl.GetInt32("status")
	list, count, err := service.NewWeixinAccount().AccountPaginate(page, limit, name, status)
	if err != nil {
		logs.Error("AccountPaginate", err.Error())
	}
	ctrl.JSONPageSuccess(list, count)
}

func (ctrl *WeixinController) AccountEdit() {
	accountId, _ := ctrl.GetInt64("accountId")
	clone, _ := ctrl.GetInt("clone")
	mdl, err := service.NewWeixinAccount().AccountFind(accountId)
	if err != nil {
		mdl = &domain.WeixinAccount{
			SortID: 99,
		}
	}
	// 是否克隆
	if clone == 1 {
		mdl.AccountID = 0
	}
	ctrl.Data["mdl"] = mdl
	ctrl.display()
}

func (ctrl *WeixinController) AccountSave() {
	mdl := domain.WeixinAccount{}
	if err := ctrl.ParseForm(&mdl); err != nil {
		logs.Error("AccountSave", err.Error())
		ctrl.JSONError(err.Error())
	}
	if mdl.AccountID == 0 {
		mdl.CreateID = int32(GlobalAdminId)
		mdl.CreateName = GlobalAdminName
	} else {
		mdl.UpdateID = int32(GlobalAdminId)
		mdl.UpdateName = GlobalAdminName
	}
	if err := service.NewWeixinAccount().AccountSave(&mdl); err != nil {
		logs.Error("AccountSave", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("保存成功", nil)
}

func (ctrl *WeixinController) AccountDestory() {
	accountId, _ := ctrl.GetInt64("accountId")
	if err := service.NewWeixinAccount().AccountDestory(accountId); err != nil {
		logs.Error("AccountDestory", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("删除成功", nil)
}

func (ctrl *WeixinController) AccountChangeStatus() {
	var mdl vmodel.Account_ChangeStatusModel
	if err := xjson.Unmarshal(ctrl.Ctx.Input.RequestBody, &mdl); err != nil {
		logs.Error("AccountChangeStatus", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	for _, accountId := range mdl.AccountIds {
		if err := service.NewWeixinAccount().AccountChangeStatus(accountId, mdl.Status); err != nil {
			logs.Error("AccountChangeStatus", err.Error())
			ctrl.JSONError(err.Error())
			return
		}
	}
	ctrl.JSONSuccess("更改状态成功", nil)
}

func (ctrl *WeixinController) AccountSaveSortId() {
	mdls := []vmodel.Account_SaveSortIdModel{}
	data := ctrl.Ctx.Input.RequestBody
	logs.Debug("AccountSaveSortId", string(data))
	if err := xjson.Unmarshal(ctrl.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("AccountSaveSortId", err.Error())
		ctrl.JSONError(err.Error())
	}
	for _, mdl := range mdls {
		if err := service.NewWeixinAccount().AccountSaveSortId(mdl.AccountId, int32(mdl.SortId)); err != nil {
			logs.Error("AccountSaveSortId", err.Error())
			ctrl.JSONError(err.Error())
			return
		}
	}
	ctrl.JSONSuccess("保存成功", nil)
}
