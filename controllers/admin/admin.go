package admin

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/tools/xcrypto"
	"haedu.gov.cn/tools/xjson"
	"haedu.gov.cn/tools/xstring"

	"haedu.gov.cn/cms/app/cms/domain"
	"haedu.gov.cn/cms/app/cms/service"
	lib "haedu.gov.cn/cms/app/tool"
	"haedu.gov.cn/cms/controllers/admin/vmodel"
	"haedu.gov.cn/cms/global"
)

type AdminController struct{ BaseController }

func (ctrl *AdminController) Index() {
	roleList, _, _ := service.NewCmsAdmin().RolePaginate(0, 99999, "")
	ctrl.Data["roleList"] = roleList
	// 获取角色权限
	roleMap := ctrl.RolePowerGet("user_manager")
	ctrl.Data["roleMap"] = roleMap
	ctrl.display()
}

func (ctrl *AdminController) AdminPaginate() {
	page, limit := ctrl.GetPagingParameters()
	roleId, _ := ctrl.GetInt64("roleId")
	realName := ctrl.GetSafeString("realName")
	userName := ctrl.GetSafeString("userName")
	list, count, _ := service.NewCmsAdmin().AdminPaginate(page, limit, roleId, realName, userName)
	ctrl.JSONPage(lib.CodeSuccess, "", list, count)
}

func (ctrl *AdminController) AdminEdit() {
	roleList, _, _ := service.NewCmsAdmin().RolePaginate(0, 99999, "")
	ctrl.Data["roleList"] = roleList
	userId, _ := ctrl.GetInt64("userId")
	mdl, err := service.NewCmsAdmin().AdminFind(userId)
	if err != nil {
		mdl = &domain.CmsAdmin{
			SortID: 99,
		}
	}
	ctrl.Data["mdl"] = mdl
	ctrl.display()
}

func (ctrl *AdminController) AdminSave() {
	mdl := domain.CmsAdmin{}
	if err := ctrl.ParseForm(&mdl); err != nil {
		logs.Error("AdminSave", err.Error())
		ctrl.JSONError(err.Error())
	}
	// 默认加密sm4加密
	if mdl.Password != "" {
		// 检查密码是否符合规则
		if psErr := CheckPasswordRole(mdl.Password); psErr != nil {
			ctrl.JSONError(psErr.Error())
		}
		// 加密后存储
		key := global.ReverseLowerString(mdl.UserName)
		mdl.PasswordSalt, _ = xstring.RandomHexStr(8)
		mdl.PasswordFormat = 1
		mdl.Password, _ = xcrypto.Sm4Encrypt(mdl.Password, key)
	}
	do := service.NewCmsAdmin()
	if err := do.AdminSave(&mdl); err != nil {
		logs.Error("AdminSave", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("保存成功", nil)
}

func (ctrl *AdminController) AdminDestory() {
	userId, _ := ctrl.GetInt64("userId")
	if err := service.NewCmsAdmin().AdminDestory(userId); err != nil {
		logs.Error("AdminDestory", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("删除成功", nil)
}

func (ctrl *AdminController) Log() {
	ctrl.display()
}

func (ctrl *AdminController) LogPaginate() {
	page, limit := ctrl.GetPagingParameters()
	userName := ctrl.GetSafeString("userName")
	list, count, _ := service.NewCmsAdmin().LogPaginate(page, limit, 0, userName)
	ctrl.JSONPage(lib.CodeSuccess, "", list, count)
}

func (ctrl *AdminController) Role() {
	// 获取角色权限
	roleMap := ctrl.RolePowerGet("user_role")
	ctrl.Data["roleMap"] = roleMap
	ctrl.display()
}

func (ctrl *AdminController) RolePaginate() {
	page, limit := ctrl.GetPagingParameters()
	name := ctrl.GetSafeString("name")
	list, count, _ := service.NewCmsAdmin().RolePaginate(page, limit, name)
	ctrl.JSONPage(lib.CodeSuccess, "", list, count)
}

func (ctrl *AdminController) RoleEdit() {
	roleId, _ := ctrl.GetInt64("roleId")
	mdl, err := service.NewCmsAdmin().RoleFind(roleId)
	if err != nil {
		mdl = &domain.CmsAdminRole{
			SortID: 99,
		}
	}
	ctrl.Data["mdl"] = mdl
	ctrl.display()
}

func (ctrl *AdminController) RoleSave() {
	mdl := domain.CmsAdminRole{}
	if err := ctrl.ParseForm(&mdl); err != nil {
		logs.Error("RoleSave", err.Error())
		ctrl.JSONError(err.Error())
	}
	actions := ctrl.Ctx.Request.FormValue("action")
	role_vals := map[string]string{}
	if err := xjson.Unmarshal([]byte(actions), &role_vals); err != nil {
		fmt.Println(err.Error())
	}
	do := service.NewCmsAdmin()
	if err := do.RoleSave(&mdl); err != nil {
		logs.Error("RoleSave", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	if err := do.RoleValSave(mdl.RoleID, role_vals); err != nil {
		logs.Error("RoleSave", err.Error())
		ctrl.JSONError(err.Error())
		return
	}

	// 保存站点权限
	siteValue := ctrl.Ctx.Request.FormValue("siteSelect")
	siteList := xstring.Split(siteValue, ",")
	if err := do.RoleSiteSave(mdl.RoleID, siteList); err != nil {
		logs.Error("RoleSave", err.Error())
		ctrl.JSONError(err.Error())
		return
	}

	ctrl.JSONSuccess("保存成功", nil)
}

func (ctrl *AdminController) RoleSaveSortId() {
	mdls := []vmodel.Role_SaveSortIdModel{}
	data := ctrl.Ctx.Input.RequestBody
	fmt.Println("RoleSaveSortId", string(data))
	if err := xjson.Unmarshal(ctrl.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("RoleSaveSortId", err.Error())
		ctrl.JSONError(err.Error())
	}
	for _, mdl := range mdls {
		if err := service.NewCmsAdmin().RoleSaveSortId(mdl.RoleId, int32(mdl.SortId)); err != nil {
			logs.Error("RoleSaveSortId", err.Error())
			ctrl.JSONError(err.Error())
			return
		}
	}
	ctrl.JSONSuccess("保存成功", nil)
}

func (ctrl *AdminController) RoleDestory() {
	roleId, _ := ctrl.GetInt64("roleId")
	if err := service.NewCmsAdmin().RoleDestory(roleId); err != nil {
		logs.Error("RoleDestory", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("删除成功", nil)
}

func (ctrl *AdminController) NavFind() {
	roleId, _ := ctrl.GetInt64("roleId")
	list, count, err := service.NewCmsAdmin().NavFind(roleId)
	if err != nil {
		logs.Error("NavFind", err.Error())
	}
	ctrl.JSONPage(lib.CodeSuccess, "", list, count)
}

func (ctrl *AdminController) RoleValueFind() {
	roleId, _ := ctrl.GetInt64("roleId")
	list, count, err := service.NewCmsAdmin().RoleValueFind(roleId)
	if err != nil {
		logs.Error("RoleValueFind", err.Error())
	}
	ctrl.JSONPage(lib.CodeSuccess, "", list, count)
}

func (ctrl *AdminController) RoleSiteFind() {
	roleId, _ := ctrl.GetInt64("roleId")
	list, count, err := service.NewCmsAdmin().RoleSiteFind(roleId)
	if err != nil {
		logs.Error("RoleSiteFind", err.Error())
	}
	ctrl.JSONPage(lib.CodeSuccess, "", list, count)
}
