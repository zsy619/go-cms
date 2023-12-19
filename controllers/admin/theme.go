package admin

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/tools/xjson"

	"haedu.gov.cn/cms/app/cms/domain"
	"haedu.gov.cn/cms/app/cms/service"
	"haedu.gov.cn/cms/app/lib"
	"haedu.gov.cn/cms/controllers/admin/vmodel"
)

type ThemeController struct{ BaseController }

// Index 管理
// @router /admin/Theme/index [get]
func (ctrl *ThemeController) Index() {
	list, _, _ := service.NewCmsTheme().ThemePaginate(1, 9999, "", "")
	ctrl.Data["theme"] = list
	// 获取角色权限
	roleMap := ctrl.RolePowerGet("themes_index")
	ctrl.Data["roleMap"] = roleMap
	ctrl.display()
}

// ThemeEdit 编辑
// @router /admin/Theme/Edit [get]
func (ctrl *ThemeController) Edit() {
	list, _, _ := service.NewCmsTheme().ThemePaginate(1, 9999, "", "")
	ctrl.Data["theme"] = list
	// 获取角色权限
	roleMap := ctrl.RolePowerGet("themes_index")
	ctrl.Data["roleMap"] = roleMap
	ctrl.display()
}

// ThemeSave 保存
// @router /admin/Theme/ThemeSave [post]
func (ctrl *ThemeController) ThemeSave() {
	mdl := domain.CmsTheme{}
	if err := ctrl.ParseForm(&mdl); err != nil {
		logs.Error("ThemeSave", err.Error())
		ctrl.JSONError(err.Error())
	}
	if mdl.ThemeID <= 0 {
		mdl.CreateID = int32(GlobalAdminId)
		mdl.CreateName = GlobalAdminName
	} else {
		mdl.UpdateID = int32(GlobalAdminId)
		mdl.UpdateName = GlobalAdminName
	}
	if err := service.NewCmsTheme().ThemeSave(&mdl); err != nil {
		logs.Error("ThemeSave", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("保存成功", nil)
}

// ThemeSaveSortId 保存排序
// @router /admin/Theme/ThemeSaveSortId [post]
func (ctrl *ThemeController) ThemeSaveSortId() {
	mdls := []vmodel.Theme_SaveSortIdModel{}
	data := ctrl.Ctx.Input.RequestBody
	fmt.Println("ThemeSaveSortId", string(data))
	if err := xjson.Unmarshal(ctrl.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("ThemeSaveSortId", err.Error())
		ctrl.JSONError(err.Error())
	}
	service := service.NewCmsTheme()
	for _, mdl := range mdls {
		if err := service.ThemeSaveSortId(mdl.ThemeId, int32(mdl.SortId)); err != nil {
			logs.Error("ThemeSaveSortId", err.Error())
			ctrl.JSONError(err.Error())
			return
		}
	}
	ctrl.JSONSuccess("保存成功", nil)
}

// ThemeDestory 删除
// @router /admin/Theme/ThemeDestory [post]
func (ctrl *ThemeController) ThemeDestory() {
	ThemeId, _ := ctrl.GetInt64("ThemeId")
	if err := service.NewCmsTheme().ThemeDestory(ThemeId); err != nil {
		logs.Error("ThemeDestory", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("删除成功", nil)
}

// ThemePaginate 列表
// @router /admin/Theme/ThemePaginate [get]
func (ctrl *ThemeController) ThemePaginate() {
	page, limit := ctrl.GetPagingParameters()
	title := ctrl.GetSafeString("title")
	name := ctrl.GetSafeString("name")
	list, count, _ := service.NewCmsTheme().ThemePaginate(page, limit, name, title)
	ctrl.JSONPage(lib.CodeSuccess, "", list, count)
}

// SetDefault 设置默认主题
// @router /admin/theme/setDefault [post]
func (ctrl *ThemeController) SetDefault() {
	mdl := struct {
		Name string `json:"name"`
	}{}
	if err := xjson.Unmarshal(ctrl.Ctx.Input.RequestBody, &mdl); err != nil {
		logs.Error("SetDefault", err.Error())
		ctrl.JSONError(err.Error())
	}
	if mdl.Name == "" {
		ctrl.JSONError("参数错误")
	}
	err := service.NewCmsTheme().ThemeSetDefault(mdl.Name)
	if err != nil {
		ctrl.JSONError(err.Error())
	}
	ctrl.JSONSuccess("", nil)
}
