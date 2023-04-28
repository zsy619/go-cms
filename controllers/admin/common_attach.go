package admin

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/controllers/admin/vmodel"
	"haedu.gov.cn/tools/xjson"
)

// Attach 附件管理
func (c *CommonController) Attach() {
	tableName := c.GetString("tableName")
	recordId, _ := c.GetInt64("recordId")
	typeId, _ := c.GetInt32("typeId")
	if tableName == "" || recordId <= 0 {
		c.JSONError("参数错误")
		return
	}
	c.Data["tableName"] = tableName
	c.Data["recordId"] = recordId
	c.Data["typeId"] = typeId
	c.display()
}

func (c *CommonController) AttachPaginate() {
	tableName := c.GetString("tableName")
	recordId, _ := c.GetInt64("recordId")
	typeId, _ := c.GetInt32("typeId")
	list, count, err := biz.NewCmsAttach().AttachPaginate(1, 99999, tableName, recordId, typeId)
	if err != nil {
		logs.Error("AttachPaginate", err.Error())
	}
	c.JSONPageSuccess(list, count)
}

func (c *CommonController) AttachSave() {
	mdl := model.CmsAttach{}
	if err := c.ParseForm(&mdl); err != nil {
		logs.Error("AttachSave", err.Error())
		c.JSONError(err.Error())
	}
	if err := biz.NewCmsAttach().AttachSave(&mdl); err != nil {
		logs.Error("AttachSave", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("保存成功", nil)
}

func (c *CommonController) AttachSaveShow() {
	mdls := vmodel.Article_AttachSaveShowModel{}
	data := c.Ctx.Input.RequestBody
	fmt.Println("AttachSaveShow", string(data))
	if err := xjson.Unmarshal(c.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("AttachSaveShow", err.Error())
		c.JSONError(err.Error())
	}
	for _, mdl := range mdls.AttachIds {
		if err := biz.NewCmsAttach().AttachSaveShow(mdl, mdls.Show); err != nil {
			logs.Error("AttachSaveShow", err.Error())
			c.JSONError(err.Error())
			return
		}
	}
	c.JSONSuccess("保存成功", nil)
}

func (c *CommonController) AttachSaveBatch() {
	mdls := []vmodel.Article_AttachSaveBatchdModel{}
	data := c.Ctx.Input.RequestBody
	fmt.Println("AttachSaveBatch", string(data))
	if err := xjson.Unmarshal(c.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("AttachSaveBatch", err.Error())
		c.JSONError(err.Error())
	}
	service := biz.NewCmsAttach()
	for _, mdl := range mdls {
		if err := service.AttachSaveInfo(mdl.AttachID, mdl.Title, mdl.Point, mdl.Click, mdl.SortID, mdl.Remark); err != nil {
			logs.Error("AttachSaveBatch", err.Error())
			c.JSONError(err.Error())
			return
		}
	}
	c.JSONSuccess("保存成功", nil)
}

func (c *CommonController) AttachDestory() {
	attachId, _ := c.GetInt64("attachId")
	if err := biz.NewCmsAttach().AttachDestory(attachId); err != nil {
		logs.Error("AttachDestory", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("删除成功", nil)
}
