package admin

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/tools/xjson"

	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/controllers/admin/vmodel"
)

func (ctrl *WeixinController) AccountPaginate() {
	page, limit := ctrl.GetPagingParameters()
	name := ctrl.GetSafeString("name")
	status, _ := ctrl.GetInt32("status")
	list, count, err := biz.NewWeixinAccount().AccountPaginate(page, limit, name, status)
	if err != nil {
		logs.Error("AccountPaginate", err.Error())
	}
	ctrl.JSONPageSuccess(list, count)
}

func (ctrl *WeixinController) AccountEdit() {
	accountId, _ := ctrl.GetInt64("accountId")
	clone, _ := ctrl.GetInt("clone")
	mdl, err := biz.NewWeixinAccount().AccountFind(accountId)
	if err != nil {
		mdl = &model.WeixinAccount{
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
	mdl := model.WeixinAccount{}
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
	if err := biz.NewWeixinAccount().AccountSave(&mdl); err != nil {
		logs.Error("AccountSave", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("保存成功", nil)
}

func (ctrl *WeixinController) AccountDestory() {
	accountId, _ := ctrl.GetInt64("accountId")
	if err := biz.NewWeixinAccount().AccountDestory(accountId); err != nil {
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
		if err := biz.NewWeixinAccount().AccountChangeStatus(accountId, mdl.Status); err != nil {
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
	fmt.Println("AccountSaveSortId", string(data))
	if err := xjson.Unmarshal(ctrl.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("AccountSaveSortId", err.Error())
		ctrl.JSONError(err.Error())
	}
	for _, mdl := range mdls {
		if err := biz.NewWeixinAccount().AccountSaveSortId(mdl.AccountId, int32(mdl.SortId)); err != nil {
			logs.Error("AccountSaveSortId", err.Error())
			ctrl.JSONError(err.Error())
			return
		}
	}
	ctrl.JSONSuccess("保存成功", nil)
}
