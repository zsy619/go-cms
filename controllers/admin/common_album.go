package admin

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/controllers/admin/vmodel"
	"haedu.gov.cn/tools/xjson"
)

// Album 相册管理
func (c *CommonController) Album() {
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

func (c *CommonController) AlbumSearch() {
	page, limit := c.GetPagingParameters()
	tableName := c.GetString("tableName")
	title := c.GetString("title")
	ext := c.GetString("ext")
	list, count, err := biz.NewCmsAlbum().AlbumSearch(page, limit, tableName, title, ext, GlobalAdminId, GlobalRoleType)
	if err != nil {
		logs.Error("AlbumSearch", err.Error())
	}
	c.JSONPageSuccess(list, count)
}

func (c *CommonController) AlbumPaginate() {
	tableName := c.GetString("tableName")
	recordId, _ := c.GetInt64("recordId")
	typeId, _ := c.GetInt32("typeId")
	list, count, err := biz.NewCmsAlbum().AlbumPaginate(1, 99999, tableName, recordId, typeId, GlobalAdminId, GlobalRoleType)
	if err != nil {
		logs.Error("AlbumPaginate", err.Error())
	}
	c.JSONPageSuccess(list, count)
}

func (c *CommonController) AlbumSave() {
	mdl := model.CmsAlbum{}
	if err := c.ParseForm(&mdl); err != nil {
		logs.Error("AlbumSave", err.Error())
		c.JSONError(err.Error())
	}
	mdl.CreateID = int32(GlobalAdminId)
	if err := biz.NewCmsAlbum().AlbumSave(&mdl); err != nil {
		logs.Error("AlbumSave", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("保存成功", nil)
}

func (c *CommonController) AlbumSaveShow() {
	mdls := vmodel.Article_AlbumSaveShowModel{}
	data := c.Ctx.Input.RequestBody
	fmt.Println("AlbumSaveShow", string(data))
	if err := xjson.Unmarshal(c.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("AlbumSaveShow", err.Error())
		c.JSONError(err.Error())
	}
	for _, mdl := range mdls.AlbumIds {
		if err := biz.NewCmsAlbum().AlbumSaveShow(mdl, mdls.Show); err != nil {
			logs.Error("AlbumSaveShow", err.Error())
			c.JSONError(err.Error())
			return
		}
	}
	c.JSONSuccess("保存成功", nil)
}

func (c *CommonController) AlbumSaveBatch() {
	mdls := []vmodel.Article_AlbumSaveBatchdModel{}
	data := c.Ctx.Input.RequestBody
	fmt.Println("AlbumSaveBatch", string(data))
	if err := xjson.Unmarshal(c.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("AlbumSaveBatch", err.Error())
		c.JSONError(err.Error())
	}
	service := biz.NewCmsAlbum()
	for _, mdl := range mdls {
		if err := service.AlbumSaveInfo(mdl.AlbumID, mdl.Title, mdl.LinkURL, mdl.Click, mdl.SortID, mdl.Remark); err != nil {
			logs.Error("AlbumSaveBatch", err.Error())
			c.JSONError(err.Error())
			return
		}
	}
	c.JSONSuccess("保存成功", nil)
}

func (c *CommonController) AlbumDestory() {
	albumId, _ := c.GetInt64("albumId")
	if err := biz.NewCmsAlbum().AlbumDestory(albumId); err != nil {
		logs.Error("AlbumDestory", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("删除成功", nil)
}
