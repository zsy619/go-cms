package admin

import (
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/controllers/admin/vmodel"
	"haedu.gov.cn/tools/xjson"
)

type NavController struct {
	BaseController
}

// Index 导航管理
// @router /admin/nav/index [get]
func (c *NavController) Index() {
	// 获取角色权限
	roleMap := c.RolePowerGet("nav")
	c.Data["roleMap"] = roleMap
	c.Data["roleId"] = GlobalRoleId
	c.display()
}

// NavEdit 编辑
// @router /admin/nav/NavEdit [get]
func (c *NavController) NavEdit() {
	navId, _ := c.GetInt64("navId")
	mdl, err := biz.NewCmsAdminNav().NavFind(navId)
	if err != nil {
		mdl = &model.CmsAdminNav{
			SortID:     99,
			CreateName: GlobalAdminName,
			Type:       "System",
		}
	}
	c.Data["mdl"] = mdl
	// 获取角色权限
	roleMap := c.RolePowerGet("nav")
	c.Data["roleMap"] = roleMap
	c.Data["roleId"] = GlobalRoleId
	c.display()
}

// NavSaveSortId 保存排序
// @router /admin/nav/NavSaveSortId [post]
func (c *NavController) NavSaveSortId() {
	mdls := []vmodel.Nav_SaveSortIdModel{}
	data := c.Ctx.Input.RequestBody
	fmt.Println("NavSaveSortId", string(data))
	if err := xjson.Unmarshal(c.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("NavSaveSortId", err.Error())
		c.JSONError(err.Error())
	}
	service := biz.NewCmsAdminNav()
	for _, mdl := range mdls {
		if err := service.NavSaveSortId(mdl.NavId, int32(GlobalAdminId), GlobalAdminName, int32(mdl.SortId)); err != nil {
			logs.Error("NavSaveSortId", err.Error())
			c.JSONError(err.Error())
			return
		}
	}
	c.JSONSuccess("保存成功", nil)
}

// NavDestroy 删除
// @router /admin/nav/NavDestroy [post]
func (c *NavController) NavDestroy() {
	navId, _ := c.GetInt64("navId")
	if err := biz.NewCmsAdminNav().NavDestroy(navId); err != nil {
		logs.Error("NavDestroy", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("删除成功", nil)
}

// NavTree 获取树形结构列表
// @router /admin/nav/NavTree [post]
func (c *NavController) NavTree() {
	navId, _ := c.GetInt64("navId")
	tree, err := biz.NewCmsAdminNav().NavTree(navId)
	if err != nil {
		logs.Error("NavTree", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.Data["json"] = tree
	c.ServeJSON()
	c.StopRun()
}

// NavSave 保存导航详情
// @router /admin/nav/NavSave [post]
func (c *NavController) NavSave() {
	mdl := model.CmsAdminNav{}
	if err := c.ParseForm(&mdl); err != nil {
		logs.Error("NavSave", err.Error())
		c.JSONError(err.Error())
	}
	if mdl.NavID > 0 {
		mdl.UpdateID = int32(GlobalAdminId)
		mdl.UpdateName = GlobalAdminName
	} else {
		mdl.CreateID = int32(GlobalAdminId)
		mdl.CreateName = GlobalAdminName
	}
	if err := biz.NewCmsAdminNav().NavSave(&mdl); err != nil {
		logs.Error("NavSave", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("保存成功", nil)
}
