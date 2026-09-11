package admin

import (
	"strconv"

	"haedu.gov.cn/cms/app/cms/domain"
	"haedu.gov.cn/cms/app/cms/service"
	lib "haedu.gov.cn/cms/app/tool"
)

// LanguageController 语言设置控制器
type LanguageController struct{ BaseController }

// LanguagePage 语言设置页面
// @router /admin/site/language [get]
func (ctrl *LanguageController) LanguagePage() {
	roleMap := ctrl.RolePowerGet("site_language")
	ctrl.Data["roleMap"] = roleMap
	ctrl.display("admin/Site/Language")
}

// LanguageData 获取语言列表数据（不分页）
// @router /admin/site/language/data [get]
func (ctrl *LanguageController) LanguageData() {
	name := ctrl.GetStringTrim("name", "")
	code := ctrl.GetStringTrim("code", "")

	langService := service.NewCmsLanguage()
	list, err := langService.LanguageAll(name, code)
	if err != nil {
		list = []*domain.CmsLanguage{}
	}
	ctrl.JSONPage(lib.CodeSuccess, "", list, int64(len(list)))
}

// LanguageEdit 语言编辑页面
// @router /admin/site/language/edit [get]
func (ctrl *LanguageController) LanguageEdit() {
	id, _ := ctrl.GetInt64("id", 0)
	langService := service.NewCmsLanguage()

	if id > 0 {
		if lang, _ := langService.LanguageOne(id); lang != nil {
			ctrl.Data["language"] = lang
		}
	} else {
		ctrl.Data["language"] = &domain.CmsLanguage{}
	}

	ctrl.display("admin/Site/LanguageEdit")
}

// LanguageSave 保存语言
// @router /admin/site/language/save [post]
func (ctrl *LanguageController) LanguageSave() {
	id, _ := ctrl.GetInt64("id", 0)
	name := ctrl.GetStringTrim("name", "")
	code := ctrl.GetStringTrim("code", "")
	icon := ctrl.GetStringTrim("icon", "")
	sortId, _ := ctrl.GetInt("sort_id", 99)
	description := ctrl.GetStringTrim("description", "")
	status, _ := ctrl.GetInt("status", 1)
	isDefault := ctrl.GetStringTrim("is_default", "off") == "on"

	lang := &domain.CmsLanguage{
		LanguageID:  id,
		Name:        name,
		Code:        code,
		Icon:        icon,
		SortID:      int32(sortId),
		Description: description,
		Status:      int32(status),
		IsDefault:   isDefault,
	}
	lang.CreateID = int32(GlobalAdminId)
	lang.CreateName = GlobalAdminName
	lang.UpdateID = int32(GlobalAdminId)
	lang.UpdateName = GlobalAdminName

	langService := service.NewCmsLanguage()
	if err := langService.LanguageSave(lang); err != nil {
		ctrl.JSONError(err.Error())
		return
	}

	ctrl.JSONSuccess("保存成功", nil)
}

// LanguageDelete 删除语言
// @router /admin/site/language/delete [get]
func (ctrl *LanguageController) LanguageDelete() {
	ids := ctrl.GetStringTrim("ids", "")
	if ids == "" {
		ctrl.JSONError("参数丢失")
		return
	}

	langService := service.NewCmsLanguage()
	if err := langService.LanguageDelete(ids); err != nil {
		ctrl.JSONError(err.Error())
		return
	}

	ctrl.JSONSuccess("删除成功", nil)
}

// LanguageSort 语言排序
// @router /admin/site/language/sort [get]
func (ctrl *LanguageController) LanguageSort() {
	id, _ := ctrl.GetInt64("id", 0)
	sortId, _ := ctrl.GetInt("sort_id", 99)

	if id == 0 {
		ctrl.JSONError("参数丢失")
		return
	}

	langService := service.NewCmsLanguage()
	if err := langService.LanguageSort(id, int32(sortId)); err != nil {
		ctrl.JSONError(err.Error())
		return
	}

	ctrl.JSONSuccess("排序成功", nil)
}

// LanguageStatus 语言状态切换
// @router /admin/site/language/status [get]
func (ctrl *LanguageController) LanguageStatus() {
	id, _ := ctrl.GetInt64("id", 0)
	status, _ := ctrl.GetInt("status", 0)

	if id == 0 {
		ctrl.JSONError("参数丢失")
		return
	}

	langService := service.NewCmsLanguage()
	if err := langService.LanguageStatus(id, int32(status)); err != nil {
		ctrl.JSONError(err.Error())
		return
	}

	ctrl.JSONSuccess("状态更新成功", nil)
}

// LanguageInit 初始化语言数据
// @router /admin/site/language/init [get]
func (ctrl *LanguageController) LanguageInit() {
	langService := service.NewCmsLanguage()
	count, err := langService.LanguageInitData()
	if err != nil {
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("成功初始化 "+strconv.Itoa(count)+" 种语言", nil)
}
