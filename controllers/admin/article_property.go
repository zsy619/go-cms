package admin

import (
	"fmt"
	"strconv"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/tools/xjson"

	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/controllers/admin/vmodel"
)

func (ctrl *ArticleController) PropertyEdit() {
	channelId, _ := ctrl.GetInt64("channelId")
	propertyId, _ := ctrl.GetInt64("propertyId")
	if propertyId <= 0 || channelId <= 0 {
		ctrl.Abort("404")
		ctrl.StopRun()
		return
	}
	mdl, err := biz.NewCmsArticle().PropertyFind(propertyId)
	if err != nil {
		mdl = &model.CmsArticleProperty{
			SortID:     99,
			CreateName: GlobalAdminName,
			CreateTime: time.Now(),
		}
	}
	ctrl.Data["mdl"] = mdl
	// 获取角色权限
	roleMap := ctrl.RolePowerGet("channel_" + strconv.FormatInt(channelId, 10) + "_article")
	ctrl.Data["roleMap"] = roleMap
	ctrl.display()
}

func (ctrl *ArticleController) PropertySave() {
	mdl := model.CmsArticleProperty{}
	if err := ctrl.ParseForm(&mdl); err != nil {
		logs.Error("PropertySave", err.Error())
		ctrl.JSONError(err.Error())
	}
	if mdl.PropertyID > 0 {
		mdl.UpdateID = int32(GlobalAdminId)
		mdl.UpdateName = GlobalAdminName
	} else {
		mdl.CreateID = int32(GlobalAdminId)
		mdl.CreateName = GlobalAdminName
	}

	if err := biz.NewCmsArticle().PropertySave(&mdl); err != nil {
		logs.Error("PropertySave", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("保存成功", nil)
}

func (ctrl *ArticleController) PropertyPaginate() {
	page, _ := ctrl.GetPagingParameters()
	limit := 99999
	parentId, _ := ctrl.GetInt64("parentId")
	articleId, _ := ctrl.GetInt64("articleId")
	status, _ := ctrl.GetInt32("status")
	callIndex := ctrl.GetSafeString("callIndex")
	title := ctrl.GetSafeString("title")

	list, count, err := biz.NewCmsArticle().PropertyPaginate(page, limit, parentId, articleId, status, callIndex, title)
	if err != nil {
		logs.Error("PropertyPaginate", err.Error())
	}
	ctrl.JSONPageSuccess(list, count)
}

func (ctrl *ArticleController) PropertyDestroy() {
	propertyId, _ := ctrl.GetInt64("propertyId")
	if err := biz.NewCmsArticle().PropertyDestroy(propertyId); err != nil {
		logs.Error("PropertyDestroy", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("删除成功", nil)
}

func (ctrl *ArticleController) PropertyChangeStatus() {
	var mdl vmodel.Property_ChangeStatusModel
	if err := xjson.Unmarshal(ctrl.Ctx.Input.RequestBody, &mdl); err != nil {
		logs.Error("PropertyChangeStatus", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	for _, propertyId := range mdl.PropertyIds {
		if err := biz.NewCmsArticle().PropertyChangeStatus(propertyId, mdl.Status); err != nil {
			logs.Error("PropertyChangeStatus", err.Error())
			ctrl.JSONError(err.Error())
			return
		}
	}
	ctrl.JSONSuccess("更改状态成功", nil)
}

func (ctrl *ArticleController) PropertySaveSortId() {
	mdls := []vmodel.Property_SaveSortIdModel{}
	data := ctrl.Ctx.Input.RequestBody
	fmt.Println("PropertySaveSortId", string(data))
	if err := xjson.Unmarshal(ctrl.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("PropertySaveSortId", err.Error())
		ctrl.JSONError(err.Error())
	}
	for _, mdl := range mdls {
		if err := biz.NewCmsArticle().PropertySaveSortId(mdl.PropertyId, int32(mdl.SortId)); err != nil {
			logs.Error("PropertySaveSortId", err.Error())
			ctrl.JSONError(err.Error())
			return
		}
	}
	ctrl.JSONSuccess("保存成功", nil)
}
