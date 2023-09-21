package admin

import (
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/app/lib"
	"haedu.gov.cn/cms/controllers/admin/vmodel"
	"haedu.gov.cn/cms/global"
	"haedu.gov.cn/tools/xcrypto"
	"haedu.gov.cn/tools/xjson"
	"haedu.gov.cn/tools/xstring"
)

type AdminController struct{ BaseController }

func (c *AdminController) Index() {
	roleList, _, _ := biz.NewCmsAdmin().RolePaginate(0, 99999, "")
	c.Data["roleList"] = roleList
	// 获取角色权限
	roleMap := c.RolePowerGet("user_manager")
	c.Data["roleMap"] = roleMap
	c.display()
}

func (c *AdminController) AdminPaginate() {
	page, limit := c.GetPagingParameters()
	roleId, _ := c.GetInt64("roleId")
	realName := c.GetSafeString("realName")
	userName := c.GetSafeString("userName")
	list, count, _ := biz.NewCmsAdmin().AdminPaginate(page, limit, roleId, realName, userName)
	c.JSONPage(lib.CodeSuccess, "", list, count)
}

func (c *AdminController) AdminEdit() {
	roleList, _, _ := biz.NewCmsAdmin().RolePaginate(0, 99999, "")
	c.Data["roleList"] = roleList
	userId, _ := c.GetInt64("userId")
	mdl, err := biz.NewCmsAdmin().AdminFind(userId)
	if err != nil {
		mdl = &model.CmsAdmin{
			SortID: 99,
		}
	}
	c.Data["mdl"] = mdl
	c.display()
}

func (c *AdminController) AdminSave() {
	mdl := model.CmsAdmin{}
	if err := c.ParseForm(&mdl); err != nil {
		logs.Error("AdminSave", err.Error())
		c.JSONError(err.Error())
	}
	// 默认加密sm4加密
	if mdl.Password != "" {
		// 检查密码是否符合规则
		if psErr := CheckPasswordRole(mdl.Password); psErr != nil {
			c.JSONError(psErr.Error())
		}
		// 加密后存储
		key := global.ReverseLowerString(mdl.UserName)
		mdl.PasswordSalt, _ = xstring.RandomHexStr(8)
		mdl.PasswordFormat = 1
		mdl.Password, _ = xcrypto.Sm4Encrypt(mdl.Password, key)
	}
	do := biz.NewCmsAdmin()
	if err := do.AdminSave(&mdl); err != nil {
		logs.Error("AdminSave", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("保存成功", nil)
}

func (c *AdminController) AdminDestory() {
	userId, _ := c.GetInt64("userId")
	if err := biz.NewCmsAdmin().AdminDestory(userId); err != nil {
		logs.Error("AdminDestory", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("删除成功", nil)
}

func (c *AdminController) Log() {
	c.display()
}

func (c *AdminController) LogPaginate() {
	page, limit := c.GetPagingParameters()
	userName := c.GetSafeString("userName")
	list, count, _ := biz.NewCmsAdmin().LogPaginate(page, limit, 0, userName)
	c.JSONPage(lib.CodeSuccess, "", list, count)
}

func (c *AdminController) Role() {
	// 获取角色权限
	roleMap := c.RolePowerGet("user_role")
	c.Data["roleMap"] = roleMap
	c.display()
}

func (c *AdminController) RolePaginate() {
	page, limit := c.GetPagingParameters()
	name := c.GetSafeString("name")
	list, count, _ := biz.NewCmsAdmin().RolePaginate(page, limit, name)
	c.JSONPage(lib.CodeSuccess, "", list, count)
}

func (c *AdminController) RoleEdit() {
	roleId, _ := c.GetInt64("roleId")
	mdl, err := biz.NewCmsAdmin().RoleFind(roleId)
	if err != nil {
		mdl = &model.CmsAdminRole{
			SortID: 99,
		}
	}
	c.Data["mdl"] = mdl
	c.display()
}

func (c *AdminController) RoleSave() {
	mdl := model.CmsAdminRole{}
	if err := c.ParseForm(&mdl); err != nil {
		logs.Error("RoleSave", err.Error())
		c.JSONError(err.Error())
	}
	actions := c.Ctx.Request.FormValue("action")
	role_vals := map[string]string{}
	if err := xjson.Unmarshal([]byte(actions), &role_vals); err != nil {
		fmt.Println(err.Error())
	}
	do := biz.NewCmsAdmin()
	if err := do.RoleSave(&mdl); err != nil {
		logs.Error("RoleSave", err.Error())
		c.JSONError(err.Error())
		return
	}
	if err := do.RoleValSave(mdl.RoleID, role_vals); err != nil {
		logs.Error("RoleSave", err.Error())
		c.JSONError(err.Error())
		return
	}

	// 保存站点权限
	var siteValue = c.Ctx.Request.FormValue("siteSelect")
	var siteList = xstring.Split(siteValue, ",")
	if err := do.RoleSiteSave(mdl.RoleID, siteList); err != nil {
		logs.Error("RoleSave", err.Error())
		c.JSONError(err.Error())
		return
	}

	c.JSONSuccess("保存成功", nil)
}

func (c *AdminController) RoleSaveSortId() {
	mdls := []vmodel.Role_SaveSortIdModel{}
	data := c.Ctx.Input.RequestBody
	fmt.Println("RoleSaveSortId", string(data))
	if err := xjson.Unmarshal(c.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("RoleSaveSortId", err.Error())
		c.JSONError(err.Error())
	}
	for _, mdl := range mdls {
		if err := biz.NewCmsAdmin().RoleSaveSortId(mdl.RoleId, int32(mdl.SortId)); err != nil {
			logs.Error("RoleSaveSortId", err.Error())
			c.JSONError(err.Error())
			return
		}
	}
	c.JSONSuccess("保存成功", nil)
}

func (c *AdminController) RoleDestory() {
	roleId, _ := c.GetInt64("roleId")
	if err := biz.NewCmsAdmin().RoleDestory(roleId); err != nil {
		logs.Error("RoleDestory", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("删除成功", nil)
}

func (c *AdminController) NavFind() {
	roleId, _ := c.GetInt64("roleId")
	list, count, err := biz.NewCmsAdmin().NavFind(roleId)
	if err != nil {
		logs.Error("NavFind", err.Error())
	}
	c.JSONPage(lib.CodeSuccess, "", list, count)
}

func (c *AdminController) RoleValueFind() {
	roleId, _ := c.GetInt64("roleId")
	list, count, err := biz.NewCmsAdmin().RoleValueFind(roleId)
	if err != nil {
		logs.Error("RoleValueFind", err.Error())
	}
	c.JSONPage(lib.CodeSuccess, "", list, count)
}

func (c *AdminController) RoleSiteFind() {
	roleId, _ := c.GetInt64("roleId")
	list, count, err := biz.NewCmsAdmin().RoleSiteFind(roleId)
	if err != nil {
		logs.Error("RoleSiteFind", err.Error())
	}
	c.JSONPage(lib.CodeSuccess, "", list, count)
}
