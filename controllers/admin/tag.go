package admin

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/tools/xjson"

	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/app/lib"
	"haedu.gov.cn/cms/controllers/admin/vmodel"
	"haedu.gov.cn/cms/global"
)

type TagController struct{ BaseController }

// Index 标签管理
// @router /admin/tag/index [get]
func (ctrl *TagController) Index() {
	siteList, _ := biz.NewCmsTag().SiteGet(GlobalRoleId, GlobalRoleType)
	ctrl.Data["site"] = siteList
	// 获取角色权限
	roleMap := ctrl.RolePowerGet("tag")
	ctrl.Data["roleMap"] = roleMap
	ctrl.display()
}

// TagEdit 编辑
// @router /admin/tag/TagEdit [get]
func (ctrl *TagController) TagEdit() {
	tagId, _ := ctrl.GetInt64("tagId")
	clone, _ := ctrl.GetInt("clone")
	mdl, err := biz.NewCmsTag().TagFind(tagId)
	if err != nil {
		mdl = &model.CmsTag{
			SortID: 99,
		}
	}
	// 是否克隆
	if clone == 1 {
		mdl.TagID = 0
	}
	ctrl.Data["mdl"] = mdl
	siteList, _ := biz.NewCmsTag().SiteGet(GlobalRoleId, GlobalRoleType)
	ctrl.Data["site"] = siteList
	// 获取角色权限
	roleMap := ctrl.RolePowerGet("tag")
	ctrl.Data["roleMap"] = roleMap
	ctrl.display()
}

// tagSave 保存
// @router /admin/tag/tagSave [post]
func (ctrl *TagController) TagSave() {
	mdl := model.CmsTag{}
	if err := ctrl.ParseForm(&mdl); err != nil {
		logs.Error("TagSave", err.Error())
		ctrl.JSONError(err.Error())
	}
	if mdl.TagID <= 0 {
		mdl.CreateID = int32(GlobalAdminId)
		mdl.CreateName = GlobalAdminName
	} else {
		mdl.UpdateID = int32(GlobalAdminId)
		mdl.UpdateName = GlobalAdminName
	}
	if err := biz.NewCmsTag().TagSave(&mdl); err != nil {
		logs.Error("TagSave", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("保存成功", nil)
}

// tagSaveSortId 保存排序
// @router /admin/tag/tagSaveSortId [post]
func (ctrl *TagController) TagSaveSortId() {
	mdls := []vmodel.Tag_SaveSortIdModel{}
	data := ctrl.Ctx.Input.RequestBody
	fmt.Println("TagSaveSortId", string(data))
	if err := xjson.Unmarshal(ctrl.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("TagSaveSortId", err.Error())
		ctrl.JSONError(err.Error())
	}
	service := biz.NewCmsTag()
	for _, mdl := range mdls {
		if err := service.TagSaveSortId(mdl.TagId, int32(mdl.SortId)); err != nil {
			logs.Error("TagSaveSortId", err.Error())
			ctrl.JSONError(err.Error())
			return
		}
	}
	ctrl.JSONSuccess("保存成功", nil)
}

// tagDestory 删除
// @router /admin/tag/tagDestory [post]
func (ctrl *TagController) TagDestory() {
	tagId, _ := ctrl.GetInt64("tagId")
	if err := biz.NewCmsTag().TagDestory(tagId); err != nil {
		logs.Error("TagDestory", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("删除成功", nil)
}

// tagChangeStatus 更改状态
// @router /admin/tag/tagChangeStatus [post]
func (ctrl *TagController) TagChangeStatus() {
	var mdl vmodel.Tag_ChangeStatusModel
	if err := xjson.Unmarshal(ctrl.Ctx.Input.RequestBody, &mdl); err != nil {
		logs.Error("TagChangeStatus", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	for _, tagId := range mdl.TagIds {
		if err := biz.NewCmsTag().TagChangeStatus(tagId, mdl.Status); err != nil {
			logs.Error("TagChangeStatus", err.Error())
			ctrl.JSONError(err.Error())
			return
		}
	}
	ctrl.JSONSuccess("更改状态成功", nil)
}

// tagPaginate 列表
// @router /admin/tag/tagPaginate [get]
func (ctrl *TagController) TagPaginate() {
	page, limit := ctrl.GetPagingParameters()
	siteId, _ := ctrl.GetInt64("siteId")
	status, _ := ctrl.GetInt32("status")
	title := ctrl.GetSafeString("title")
	name := ctrl.GetSafeString("name")
	var siteIds []int64
	if siteId > 0 {
		siteIds = append(siteIds, siteId)
	} else {
		if !global.IsSuper(GlobalRoleType) {
			siteIdList, _, _ := biz.NewCmsAdmin().RoleSiteFind(GlobalRoleId)
			for _, item := range siteIdList {
				siteIds = append(siteIds, item.SiteID)
			}
		}
	}
	list, count, _ := biz.NewCmsTag().TagPaginate(page, limit, -1, name, title, status, siteIds...)
	ctrl.JSONPage(lib.CodeSuccess, "", list, count)
}
