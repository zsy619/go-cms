package admin

import (
	"github.com/beego/beego/v2/core/logs"
	"github.com/zsy619/tools/xjson"

	"haedu.gov.cn/cms/app/cms/domain"
	"haedu.gov.cn/cms/app/cms/service"
	"haedu.gov.cn/cms/controllers/admin/vmodel"
)

type NavController struct{ BaseController }

// Index 导航管理
// @router /admin/nav/index [get]
func (ctrl *NavController) Index() {
	// 获取角色权限
	roleMap := ctrl.RolePowerGet("nav")
	ctrl.Data["roleMap"] = roleMap
	ctrl.Data["roleId"] = GlobalRoleId
	ctrl.display()
}

// NavEdit 编辑
// @router /admin/nav/NavEdit [get]
func (ctrl *NavController) NavEdit() {
	navId, _ := ctrl.GetInt64("navId")
	mdl, err := service.NewCmsAdminNav().NavFind(navId)
	if err != nil {
		mdl = &domain.CmsAdminNav{
			SortID:     99,
			CreateName: GlobalAdminName,
			Type:       "System",
		}
	}
	ctrl.Data["mdl"] = mdl
	// 获取角色权限
	roleMap := ctrl.RolePowerGet("nav")
	ctrl.Data["roleMap"] = roleMap
	ctrl.Data["roleId"] = GlobalRoleId
	ctrl.display()
}

// NavSaveSortId 保存排序
// @router /admin/nav/NavSaveSortId [post]
func (ctrl *NavController) NavSaveSortId() {
	mdls := []vmodel.Nav_SaveSortIdModel{}
	data := ctrl.Ctx.Input.RequestBody
	logs.Debug("NavSaveSortId", string(data))
	if err := xjson.Unmarshal(ctrl.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("NavSaveSortId", err.Error())
		ctrl.JSONError(err.Error())
	}
	service := service.NewCmsAdminNav()
	for _, mdl := range mdls {
		if err := service.NavSaveSortId(mdl.NavId, int32(GlobalAdminId), GlobalAdminName, int32(mdl.SortId)); err != nil {
			logs.Error("NavSaveSortId", err.Error())
			ctrl.JSONError(err.Error())
			return
		}
	}
	ctrl.JSONSuccess("保存成功", nil)
}

// NavDestroy 删除
// @router /admin/nav/NavDestroy [post]
func (ctrl *NavController) NavDestroy() {
	navId, _ := ctrl.GetInt64("navId")
	if err := service.NewCmsAdminNav().NavDestroy(navId); err != nil {
		logs.Error("NavDestroy", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("删除成功", nil)
}

// NavTree 获取树形结构列表
// @router /admin/nav/NavTree [post]
func (ctrl *NavController) NavTree() {
	navId, _ := ctrl.GetInt64("navId")
	tree, err := service.NewCmsAdminNav().NavTree(navId)
	if err != nil {
		logs.Error("NavTree", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.Data["json"] = tree
	ctrl.ServeJSON()
	ctrl.StopRun()
}

// NavSave 保存导航详情
// @router /admin/nav/NavSave [post]
func (ctrl *NavController) NavSave() {
	mdl := domain.CmsAdminNav{}
	if err := ctrl.ParseForm(&mdl); err != nil {
		logs.Error("NavSave", err.Error())
		ctrl.JSONError(err.Error())
	}
	if mdl.NavID > 0 {
		mdl.UpdateID = int32(GlobalAdminId)
		mdl.UpdateName = GlobalAdminName
	} else {
		mdl.CreateID = int32(GlobalAdminId)
		mdl.CreateName = GlobalAdminName
	}
	if err := service.NewCmsAdminNav().NavSave(&mdl); err != nil {
		logs.Error("NavSave", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("保存成功", nil)
}
