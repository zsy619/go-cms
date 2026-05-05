package admin

import (
	"fmt"
	"html/template"
	"regexp"
	"strings"

	"github.com/beego/beego/v2/core/logs"

	"haedu.gov.cn/cms/app/cms/domain"
	"haedu.gov.cn/cms/app/cms/service"
	"haedu.gov.cn/cms/controllers"
	"haedu.gov.cn/cms/controllers/admin/vmodel"
	"haedu.gov.cn/cms/global"
)

// BaseController 后台管理基础控制器
// 继承自controllers.BaseController,提供后台管理的公共功能
type BaseController struct{ controllers.BaseController }

// Prepare 后台控制器前置处理
// 功能: 初始化XSRF防护,加载用户Session,设置权限数据
func (ctrl *BaseController) Prepare() {
	ctrl.BaseController.Prepare()
	ctrl.EnableXSRF = true
	ctrl.XSRFExpire = 3600
	ctrl.Data["xsrfdata"] = template.HTML(ctrl.XSRFFormHTML())
	ctrl.Data["xsrf_token"] = ctrl.XSRFToken()

	if GlobalAdminId == 0 {
		user := ctrl.GetSession("user").(*domain.CmsAdmin)
		GlobalAdminId = user.UserID
		GlobalUserType = int(user.UserType)
		GlobalAdminName = user.UserName
		GlobalRealName = user.RealName
		GlobalRoleId = user.RoleID
		GlobalRoleType = user.RoleType
		logs.Debug("后台用户Session初始化: adminId=%d, roleType=%s", GlobalAdminId, GlobalRoleType)
	}
	ctrl.Data["roleId"] = GlobalRoleId
}

// Finish 后台控制器后置处理
func (ctrl *BaseController) Finish() {
	// 预留后置处理逻辑
}

// display 渲染带后台布局的模板
// @param tpl 模板路径(可选,默认使用 controller/action.html)
func (ctrl *BaseController) display(tpl ...string) {
	var tplname string
	if len(tpl) > 0 {
		tplname = tpl[0] + ".html"
	} else {
		tplname = "admin/" + ctrl.ControllerName + "/" + ctrl.ActionName + ".html"
	}
	ctrl.Layout = "admin/layout/layout.html"
	ctrl.TplName = tplname
}

// displayNoLayout 渲染不带布局的模板(用于特殊页面)
func (ctrl *BaseController) displayNoLayout(tpl ...string) {
	var tplname string
	if len(tpl) > 0 {
		tplname = tpl[0] + ".html"
	} else {
		tplname = "admin/" + ctrl.ControllerName + "/" + ctrl.ActionName + ".html"
	}
	ctrl.TplName = tplname
}

// IsLogin 检查用户是否已登录
// @return int64 登录用户ID,未登录返回0
func (ctrl *BaseController) IsLogin() int64 {
	id := ctrl.GetSession(`adminId`)
	if id == nil {
		return 0
	} else {
		switch id := id.(type) {
		case int64:
			rt := id
			return rt
		default:
			return 0
		}
	}
}

// RolePowerGet 获取当前用户在指定导航的操作权限
// @param navName 导航名称,如"article_index"
// @return vmodel.RoleAction 权限操作对象
func (ctrl *BaseController) RolePowerGet(navName string) vmodel.RoleAction {
	roleAction := vmodel.RoleAction{}
	if global.IsSuper(GlobalRoleType) {
		roleAction.IsSuccess = true
		roleAction.IsSuper = true
		roleAction.IsHasAdd = true
		roleAction.IsHasAudit = true
		roleAction.IsHasEdit = true
		roleAction.IsHasDelete = true
		roleAction.IsHasView = true
		roleAction.IsHasAttach = true
		roleAction.IsHasAlbum = true
	} else {
		mdl, err := service.NewCmsAdmin().RolePower(GlobalRoleId, navName)
		if err != nil {
			logs.Error("获取角色权限失败: roleId=%d, navName=%s, error=%v", GlobalRoleId, navName, err)
		}
		if len(mdl.Action) > 0 {
			action := strings.Split(mdl.Action, ",")
			if len(action) > 0 {
				roleAction.IsSuccess = true
				for i := 0; i < len(action); i++ {
					switch action[i] {
					case "Add":
						roleAction.IsHasAdd = true
					case "Audit":
						roleAction.IsHasAudit = true
					case "Edit":
						roleAction.IsHasEdit = true
					case "Delete":
						roleAction.IsHasDelete = true
					case "View":
						roleAction.IsHasView = true
					case "Attach":
						roleAction.IsHasAttach = true
					case "Album":
						roleAction.IsHasAlbum = true
					}
				}
			}
		}
	}
	return roleAction
}

// CheckPasswordRole 检查密码是否符合安全规则
// 规则: 长度8-16位,必须包含数字、小写字母、大写字母、特殊字符
// @param ps 密码字符串
// @return error 验证失败返回错误信息
func CheckPasswordRole(ps string) error {
	if len(ps) < 8 || len(ps) > 16 {
		return fmt.Errorf("密码长度应为8~16位")
	}
	num := `[0-9]{1}`
	a_z := `[a-z]{1}`
	A_Z := `[A-Z]{1}`
	symbol := `[!@#~$%^&*()+|_]{1}`
	if b, err := regexp.MatchString(num, ps); !b || err != nil {
		return fmt.Errorf("密码中必须包含数字")
	}
	if b, err := regexp.MatchString(a_z, ps); !b || err != nil {
		return fmt.Errorf("密码中必须包含小写字母")
	}
	if b, err := regexp.MatchString(A_Z, ps); !b || err != nil {
		return fmt.Errorf("密码中必须包含大写字母")
	}
	if b, err := regexp.MatchString(symbol, ps); !b || err != nil {
		return fmt.Errorf("密码中必须包含特殊字符")
	}
	return nil
}