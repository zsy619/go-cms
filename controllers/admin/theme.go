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

type ThemeController struct{ BaseController }

// Index 管理
// @router /admin/Theme/index [get]
func (c *ThemeController) Index() {
	list, _, _ := biz.NewCmsTheme().ThemePaginate(1, 9999, "", "")
	c.Data["theme"] = list
	c.display()
}

// ThemeEdit 编辑
// @router /admin/Theme/Edit [get]
func (c *ThemeController) Edit() {
	list, _, _ := biz.NewCmsTheme().ThemePaginate(1, 9999, "", "")
	c.Data["theme"] = list
	c.display()
}

// ThemeSave 保存
// @router /admin/Theme/ThemeSave [post]
func (c *ThemeController) ThemeSave() {
	mdl := model.CmsTheme{}
	if err := c.ParseForm(&mdl); err != nil {
		logs.Error("ThemeSave", err.Error())
		c.JSONError(err.Error())
	}
	if err := biz.NewCmsTheme().ThemeSave(&mdl); err != nil {
		logs.Error("ThemeSave", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("保存成功", nil)
}

// ThemeSaveSortId 保存排序
// @router /admin/Theme/ThemeSaveSortId [post]
func (c *ThemeController) ThemeSaveSortId() {
	mdls := []vmodel.Theme_SaveSortIdModel{}
	data := c.Ctx.Input.RequestBody
	fmt.Println("ThemeSaveSortId", string(data))
	if err := xjson.Unmarshal(c.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("ThemeSaveSortId", err.Error())
		c.JSONError(err.Error())
	}
	service := biz.NewCmsTheme()
	for _, mdl := range mdls {
		if err := service.ThemeSaveSortId(mdl.ThemeId, int32(mdl.SortId)); err != nil {
			logs.Error("ThemeSaveSortId", err.Error())
			c.JSONError(err.Error())
			return
		}
	}
	c.JSONSuccess("保存成功", nil)
}

// ThemeDestory 删除
// @router /admin/Theme/ThemeDestory [post]
func (c *ThemeController) ThemeDestory() {
	ThemeId, _ := c.GetInt64("ThemeId")
	if err := biz.NewCmsTheme().ThemeDestory(ThemeId); err != nil {
		logs.Error("ThemeDestory", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("删除成功", nil)
}

// ThemePaginate 列表
// @router /admin/Theme/ThemePaginate [get]
func (c *ThemeController) ThemePaginate() {
	page, limit := c.GetPagingParameters()
	title := c.GetString("title")
	name := c.GetString("name")
	list, count, _ := biz.NewCmsTheme().ThemePaginate(page, limit, name, title)
	c.JSONPage(lib.CodeSuccess, "", list, count)
}

// SetDefault 设置默认主题
// @router /admin/theme/setDefault [post]
func (c *ThemeController) SetDefault() {
	mdl := struct {
		Name string `json:"name"`
	}{}
	if err := xjson.Unmarshal(c.Ctx.Input.RequestBody, &mdl); err != nil {
		logs.Error("SetDefault", err.Error())
		c.JSONError(err.Error())
	}
	if mdl.Name == "" {
		c.JSONError("参数错误")
	}
	err := biz.NewCmsTheme().ThemeSetDefault(mdl.Name)
	if err != nil {
		c.JSONError(err.Error())
	}
	c.JSONSuccess("", nil)
}
