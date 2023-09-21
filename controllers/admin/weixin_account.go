package admin

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/controllers/admin/vmodel"
	"haedu.gov.cn/tools/xjson"
)

func (c *WeixinController) AccountPaginate() {
	page, limit := c.GetPagingParameters()
	name := c.GetSafeString("name")
	status, _ := c.GetInt32("status")
	list, count, err := biz.NewWeixinAccount().AccountPaginate(page, limit, name, status)
	if err != nil {
		logs.Error("AccountPaginate", err.Error())
	}
	c.JSONPageSuccess(list, count)
}

func (c *WeixinController) AccountEdit() {
	accountId, _ := c.GetInt64("accountId")
	clone, _ := c.GetInt("clone")
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
	c.Data["mdl"] = mdl
	c.display()
}

func (c *WeixinController) AccountSave() {
	mdl := model.WeixinAccount{}
	if err := c.ParseForm(&mdl); err != nil {
		logs.Error("AccountSave", err.Error())
		c.JSONError(err.Error())
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
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("保存成功", nil)
}

func (c *WeixinController) AccountDestory() {
	accountId, _ := c.GetInt64("accountId")
	if err := biz.NewWeixinAccount().AccountDestory(accountId); err != nil {
		logs.Error("AccountDestory", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("删除成功", nil)
}

func (c *WeixinController) AccountChangeStatus() {
	var mdl vmodel.Account_ChangeStatusModel
	if err := xjson.Unmarshal(c.Ctx.Input.RequestBody, &mdl); err != nil {
		logs.Error("AccountChangeStatus", err.Error())
		c.JSONError(err.Error())
		return
	}
	for _, accountId := range mdl.AccountIds {
		if err := biz.NewWeixinAccount().AccountChangeStatus(accountId, mdl.Status); err != nil {
			logs.Error("AccountChangeStatus", err.Error())
			c.JSONError(err.Error())
			return
		}
	}
	c.JSONSuccess("更改状态成功", nil)
}

func (c *WeixinController) AccountSaveSortId() {
	mdls := []vmodel.Account_SaveSortIdModel{}
	data := c.Ctx.Input.RequestBody
	fmt.Println("AccountSaveSortId", string(data))
	if err := xjson.Unmarshal(c.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("AccountSaveSortId", err.Error())
		c.JSONError(err.Error())
	}
	for _, mdl := range mdls {
		if err := biz.NewWeixinAccount().AccountSaveSortId(mdl.AccountId, int32(mdl.SortId)); err != nil {
			logs.Error("AccountSaveSortId", err.Error())
			c.JSONError(err.Error())
			return
		}
	}
	c.JSONSuccess("保存成功", nil)
}
