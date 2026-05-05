package admin

import (
	"github.com/beego/beego/v2/core/logs"
	"github.com/zsy619/tools/xjson"

	"haedu.gov.cn/cms/app/cms/domain"
	"haedu.gov.cn/cms/app/cms/service"
	lib "haedu.gov.cn/cms/app/tool"
	"haedu.gov.cn/cms/controllers/admin/vmodel"
)

type LinkController struct{ BaseController }

// Index 链接管理
// @router /admin/link/index [get]
func (ctrl *LinkController) Index() {
	// 根据站点权限查询站点列表
	siteList, categoryList, _ := service.NewCmsLink().SiteCategoryGet(GlobalRoleId, GlobalRoleType)
	ctrl.Data["categoryList"] = categoryList
	ctrl.Data["siteList"] = siteList
	// 获取角色权限
	roleMap := ctrl.RolePowerGet("links_index")
	ctrl.Data["roleMap"] = roleMap
	ctrl.display()
}

func (ctrl *LinkController) LinkEdit() {
	siteId, _ := ctrl.GetInt64("siteId")
	linkId, _ := ctrl.GetInt64("linkId")
	clone, _ := ctrl.GetInt("clone")
	mdl, err := service.NewCmsLink().LinkFind(linkId)
	if err != nil {
		mdl = &domain.CmsLink{
			SortID: 99,
			Target: "_blank",
			SiteID: siteId,
		}
	}
	// 是否克隆
	if clone == 1 {
		mdl.LinkID = 0
	}
	ctrl.Data["mdl"] = mdl
	_, categoryList, _ := service.NewCmsLink().SiteCategoryGet(GlobalRoleId, GlobalRoleType)
	ctrl.Data["categoryList"] = categoryList
	// 获取角色权限
	roleMap := ctrl.RolePowerGet("links_index")
	ctrl.Data["roleMap"] = roleMap
	ctrl.display()
}

func (ctrl *LinkController) LinkSave() {
	mdl := domain.CmsLink{}
	if err := ctrl.ParseForm(&mdl); err != nil {
		logs.Error("LinkSave", err.Error())
		ctrl.JSONError(err.Error())
	}
	if mdl.LinkID <= 0 {
		mdl.CreateID = int32(GlobalAdminId)
		mdl.CreateName = GlobalAdminName
	} else {
		mdl.UpdateID = int32(GlobalAdminId)
		mdl.UpdateName = GlobalAdminName
	}
	if err := service.NewCmsLink().LinkSave(&mdl); err != nil {
		logs.Error("LinkSave", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("保存成功", nil)
}

// LinkSaveSortId 保存排序
func (ctrl *LinkController) LinkSaveSortId() {
	mdls := []vmodel.Link_SaveSortIdModel{}
	data := ctrl.Ctx.Input.RequestBody
	logs.Debug("LinkSaveSortId", string(data))
	if err := xjson.Unmarshal(ctrl.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("LinkSaveSortId", err.Error())
		ctrl.JSONError(err.Error())
	}
	service := service.NewCmsLink()
	for _, mdl := range mdls {
		if err := service.LinkSaveSortId(mdl.LinkId, int32(mdl.SortId)); err != nil {
			logs.Error("LinkSaveSortId", err.Error())
			ctrl.JSONError(err.Error())
			return
		}
	}
	// biz.NewApiLink().InitCache()
	ctrl.JSONSuccess("保存成功", nil)
}

func (ctrl *LinkController) LinkDestory() {
	linkId, _ := ctrl.GetInt64("linkId")
	if err := service.NewCmsLink().LinkDestory(linkId); err != nil {
		logs.Error("LinkDestory", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("删除成功", nil)
}

func (ctrl *LinkController) LinkChangeStatus() {
	var mdl vmodel.Link_ChangeStatusModel
	if err := xjson.Unmarshal(ctrl.Ctx.Input.RequestBody, &mdl); err != nil {
		logs.Error("LinkChangeStatus", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	for _, linkid := range mdl.LinkIds {
		if err := service.NewCmsLink().LinkChangeStatus(linkid, mdl.Status); err != nil {
			logs.Error("LinkChangeStatus", err.Error())
			ctrl.JSONError(err.Error())
			return
		}
	}
	ctrl.JSONSuccess("更改状态成功", nil)
}

// LinkPaginate 列表
// @router /admin/link/linkpaginate [get]
func (ctrl *LinkController) LinkPaginate() {
	page, limit := ctrl.GetPagingParameters()
	siteId, _ := ctrl.GetInt64("siteId")
	categoryId, _ := ctrl.GetInt64("categoryId")
	status, _ := ctrl.GetInt32("status")
	title := ctrl.GetSafeString("title")
	callIndex := ctrl.GetSafeString("callIndex")

	siteIds := service.NewCmsLink().SiteIdsGet(siteId, GlobalRoleType, GlobalRoleId)
	list, count, _ := service.NewCmsLink().LinkPaginate(page, limit, -1, categoryId, title, callIndex, status, siteIds...)
	ctrl.JSONPage(lib.CodeSuccess, "", list, count)
}

// Category 链接分类
// @router /admin/link/category [get]
func (ctrl *LinkController) Category() {
	// 根据站点权限查询站点列表
	siteList, _, _ := service.NewCmsLink().SiteCategoryGet(GlobalRoleId, GlobalRoleType)
	ctrl.Data["siteList"] = siteList
	// 获取角色权限
	roleMap := ctrl.RolePowerGet("links_category")
	ctrl.Data["roleMap"] = roleMap
	ctrl.display()
}

func (ctrl *LinkController) CategoryEdit() {
	siteId, _ := ctrl.GetInt64("siteId")
	categoryId, _ := ctrl.GetInt64("categoryId")
	mdl, err := service.NewCmsLink().CategoryFind(categoryId)
	if err != nil {
		mdl = &domain.CmsLinkCategory{
			SortID: 99,
			SiteID: siteId,
		}
	}
	ctrl.Data["mdl"] = mdl
	siteList, _, _ := service.NewCmsLink().SiteCategoryGet(GlobalRoleId, GlobalRoleType)
	ctrl.Data["siteList"] = siteList
	ctrl.display()
}

func (ctrl *LinkController) CategorySave() {
	mdl := domain.CmsLinkCategory{}
	if err := ctrl.ParseForm(&mdl); err != nil {
		logs.Error("CategorySave", err.Error())
		ctrl.JSONError(err.Error())
	}
	if err := service.NewCmsLink().CategorySave(&mdl); err != nil {
		logs.Error("CategorySave", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("保存成功", nil)
}

func (ctrl *LinkController) CategorySaveSortId() {
	mdls := []vmodel.Category_SaveSortIdModel{}
	data := ctrl.Ctx.Input.RequestBody
	logs.Debug("CategorySaveSortId", string(data))
	if err := xjson.Unmarshal(ctrl.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("CategorySaveSortId", err.Error())
		ctrl.JSONError(err.Error())
	}
	for _, mdl := range mdls {
		if err := service.NewCmsLink().CategorySaveSortId(mdl.CategoryId, int32(mdl.SortId)); err != nil {
			logs.Error("CategorySaveSortId", err.Error())
			ctrl.JSONError(err.Error())
			return
		}
	}
	ctrl.JSONSuccess("保存成功", nil)
}

func (ctrl *LinkController) CategoryDestory() {
	categoryId, _ := ctrl.GetInt64("categoryId")
	if err := service.NewCmsLink().CategoryDestory(categoryId); err != nil {
		logs.Error("CategoryDestory", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("删除成功", nil)
}

// CategoryPaginate 列表
// @router /admin/link/categorypaginate [get]
func (ctrl *LinkController) CategoryPaginate() {
	page, limit := ctrl.GetPagingParameters()
	siteId, _ := ctrl.GetInt64("siteId")
	title := ctrl.GetSafeString("title")
	callIndex := ctrl.GetSafeString("callIndex")

	siteIds := service.NewCmsLink().SiteIdsGet(siteId, GlobalRoleType, GlobalRoleId)
	categoryList, count, _ := service.NewCmsLink().CategoryPaginate(page, limit, -1, title, callIndex, siteIds...)
	ctrl.JSONPage(lib.CodeSuccess, "", categoryList, count)
}
