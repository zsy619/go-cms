package admin

import (
	"github.com/beego/beego/v2/core/logs"
	"github.com/zsy619/tools/xjson"

	"haedu.gov.cn/cms/app/cms/domain"
	"haedu.gov.cn/cms/app/cms/service"
	"haedu.gov.cn/cms/controllers/admin/vmodel"
)

// Album 相册管理
func (ctrl *CommonController) Album() {
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

func (ctrl *CommonController) AlbumSearch() {
	page, limit := ctrl.GetPagingParameters()
	tableName := ctrl.GetSafeString("tableName")
	title := ctrl.GetSafeString("title")
	ext := ctrl.GetSafeString("ext")
	list, count, err := service.NewCmsAlbum().AlbumSearch(page, limit, tableName, title, ext, GlobalAdminId, GlobalRoleType)
	if err != nil {
		logs.Error("AlbumSearch", err.Error())
	}
	ctrl.JSONPageSuccess(list, count)
}

func (ctrl *CommonController) AlbumPaginate() {
	tableName := ctrl.GetSafeString("tableName")
	recordId, _ := ctrl.GetInt64("recordId")
	typeId, _ := ctrl.GetInt32("typeId")
	list, count, err := service.NewCmsAlbum().AlbumPaginate(1, 99999, tableName, recordId, typeId, GlobalAdminId, GlobalRoleType)
	if err != nil {
		logs.Error("AlbumPaginate", err.Error())
	}
	ctrl.JSONPageSuccess(list, count)
}

func (ctrl *CommonController) AlbumSave() {
	mdl := domain.CmsAlbum{}
	if err := ctrl.ParseForm(&mdl); err != nil {
		logs.Error("AlbumSave", err.Error())
		ctrl.JSONError(err.Error())
	}
	mdl.CreateID = int32(GlobalAdminId)
	if err := service.NewCmsAlbum().AlbumSave(&mdl); err != nil {
		logs.Error("AlbumSave", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("保存成功", nil)
}

func (ctrl *CommonController) AlbumSaveShow() {
	mdls := vmodel.Article_AlbumSaveShowModel{}
	data := ctrl.Ctx.Input.RequestBody
	logs.Debug("AlbumSaveShow", string(data))
	if err := xjson.Unmarshal(ctrl.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("AlbumSaveShow", err.Error())
		ctrl.JSONError(err.Error())
	}
	for _, mdl := range mdls.AlbumIds {
		if err := service.NewCmsAlbum().AlbumSaveShow(mdl, mdls.Show); err != nil {
			logs.Error("AlbumSaveShow", err.Error())
			ctrl.JSONError(err.Error())
			return
		}
	}
	ctrl.JSONSuccess("保存成功", nil)
}

func (ctrl *CommonController) AlbumSaveBatch() {
	mdls := []vmodel.Article_AlbumSaveBatchdModel{}
	data := ctrl.Ctx.Input.RequestBody
	logs.Debug("AlbumSaveBatch", string(data))
	if err := xjson.Unmarshal(ctrl.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("AlbumSaveBatch", err.Error())
		ctrl.JSONError(err.Error())
	}
	service := service.NewCmsAlbum()
	for _, mdl := range mdls {
		if err := service.AlbumSaveInfo(mdl.AlbumID, mdl.Title, mdl.LinkURL, mdl.Click, mdl.SortID, mdl.Remark); err != nil {
			logs.Error("AlbumSaveBatch", err.Error())
			ctrl.JSONError(err.Error())
			return
		}
	}
	ctrl.JSONSuccess("保存成功", nil)
}

func (ctrl *CommonController) AlbumDestory() {
	albumId, _ := ctrl.GetInt64("albumId")
	if err := service.NewCmsAlbum().AlbumDestory(albumId); err != nil {
		logs.Error("AlbumDestory", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("删除成功", nil)
}
