package admin

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/app/lib"
	"haedu.gov.cn/cms/controllers/admin/vmodel"
	"haedu.gov.cn/tools/xjson"
)

func (c *WeixinController) MenuFind() {
	accountId, _ := c.GetInt64("accountId")
	list, count, err := biz.NewWeixinMenu().MenuPaginate(1, 99999, accountId)
	if err != nil {
		logs.Error("MenuFind", err.Error())
	}
	c.JSONPage(lib.CodeSuccess, "", list, count)
}

func (c *WeixinController) MenuEdit() {
	accountId, _ := c.GetInt64("accountId")
	if accountId <= 0 {
		c.Abort("404")
		c.StopRun()
		return
	}
	c.Data["accountId"] = accountId
	menuId, _ := c.GetInt64("menuId")
	parentId, _ := c.GetInt64("parentId")
	mdl, err := biz.NewWeixinMenu().MenuFind(menuId)
	if err != nil {
		mdl = &model.WeixinMenu{
			AccountID: accountId,
			ParentID:  parentId,
			SortID:    99,
			Type:      "view",
		}
	}
	c.Data["mdl"] = mdl
	c.display()
}

func (c *WeixinController) MenuSave() {
	mdl := model.WeixinMenu{}
	if err := c.ParseForm(&mdl); err != nil {
		logs.Error("MenuSave", err.Error())
		c.JSONError(err.Error())
	}
	if err := biz.NewWeixinMenu().MenuSave(&mdl); err != nil {
		logs.Error("MenuSave", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("保存成功", nil)
}

func (c *WeixinController) MenuSaveSortId() {
	mdls := []vmodel.Menu_SaveSortIdModel{}
	data := c.Ctx.Input.RequestBody
	fmt.Println("MenuSaveSortId", string(data))
	if err := xjson.Unmarshal(c.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("MenuSaveSortId", err.Error())
		c.JSONError(err.Error())
	}
	for _, mdl := range mdls {
		if err := biz.NewWeixinMenu().MenuSaveSortId(mdl.MenuId, int32(mdl.SortId)); err != nil {
			logs.Error("MenuSaveSortId", err.Error())
			c.JSONError(err.Error())
			return
		}
	}
	c.JSONSuccess("保存成功", nil)
}

func (c *WeixinController) MenuDestory() {
	menuId, _ := c.GetInt64("menuId")
	if err := biz.NewWeixinMenu().MenuDestory(menuId); err != nil {
		logs.Error("MenuDestory", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("删除成功", nil)
}
