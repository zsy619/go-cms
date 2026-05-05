package admin

import (
	"net/http"

	"github.com/beego/beego/v2/core/logs"
	"github.com/zsy619/tools/xcache"
	"github.com/zsy619/tools/xcas"
	"github.com/zsy619/tools/xgeneric"

	"haedu.gov.cn/cms/app/cms/domain"
	"haedu.gov.cn/cms/app/cms/service"
	lib "haedu.gov.cn/cms/app/tool"
	"haedu.gov.cn/cms/controllers"
	"haedu.gov.cn/cms/controllers/admin/vmodel"
)

type LoginController struct{ controllers.BaseController }

// AdminLogin 管理员登录页面渲染
// @router cms/admin/login [get]
func (ctrl *LoginController) AdminLogin() {
	ctrl.Data["captcha"] = "/captcha"
	ctrl.TplName = "admin/login/login.html"
}

// SavaAdminState 保存管理员登录状态到Session和全局变量
// @param user *domain.CmsAdmin 管理员用户实体
func (ctrl *LoginController) SavaAdminState(user *domain.CmsAdmin) {
	ctrl.SetSession("adminId", user.UserID)
	ctrl.SetSession("adminAccount", user.UserName)
	ctrl.SetSession("adminName", user.UserName)
	ctrl.SetSession("realName", user.RealName)
	ctrl.SetSession("adminState", 1)
	ctrl.SetSession("adminLevel", 1)
	ctrl.SetSession("user", user)

	GlobalAdminId = user.UserID
	GlobalUserType = int(user.UserType)
	GlobalAdminName = user.UserName
	GlobalRealName = user.RealName
	GlobalIsAudit = user.IsAudit
	GlobalRoleId = user.RoleID
	GlobalRoleType = user.RoleType
}

// AdminLoginVerify 管理员登录验证处理
// @router cms/admin/login/verify [post]
func (ctrl *LoginController) AdminLoginVerify() {
	result := vmodel.LoginResult{
		Code:    0,
		Message: "登陆完成，载入中...",
		Url:     "",
	}

	captcha := ctrl.GetSafeString("captcha")
	if !controllers.VerifyCode(captcha) {
		result.Code = 1
		result.Message = "验证码错误"
		ctrl.Data["json"] = &result
		ctrl.ServeJSON()
		return
	}

	username := ctrl.GetSafeString("username")
	password := ctrl.GetSafeString("password")
	logs.Debug("登录尝试: username=%s, captcha=%s", username, captcha)

	adminDo := service.NewCmsAdmin()
	user, err := adminDo.Login(username, password, 0, service.LoginAll)
	if err != nil {
		logs.Error("登录失败: %v", err)
		result.Code = 2
		result.Message = "账号密码错误"
		ctrl.Data["json"] = &result
		ctrl.ServeJSON()
		return
	}
	ctrl.SavaAdminState(user)
	adminDo.LoginLog(user.UserID, username, "AdminLoginVerify", "", "", "OK", ctrl.GetClientIp())
	result.Url = "/admin/index"
	ctrl.Data["json"] = &result
	ctrl.ServeJSON()
}

// Logout 管理员退出登录
// @router cms/admin/logout [get]
func (ctrl *LoginController) Logout() {
	ctrl.DestroySession()
	ctrl.Redirect("/admin/login", 302)
}

// School 学校用户登录入口
func (ctrl *LoginController) School() {
	ctrl.login(lib.LoginCasPathOfSchool, lib.LoginPathOfSchool, "school")
}

// Admin 管理员登录入口
func (ctrl *LoginController) Admin() {
	ctrl.login(lib.LoginCasPathOfAdmin, lib.LoginPathOfAdmin, "admin")
}

// login 处理CAS单点登录的核心逻辑
// @param loginCasPath CAS登录重定向路径
// @param loginPath 服务端回调路径
// @param kind 登录类型(school/admin)
func (ctrl *LoginController) login(loginCasPath, loginPath string, kind string) {
	xcache.SetDiskvString("cas_login", kind)
	GlobalAuthFlag = kind
	ticket := ctrl.GetSafeString("ticket")
	logs.Debug("CAS ticket: %s", ticket)

	if ticket == "" {
		ctrl.Redirect(loginCasPath, http.StatusFound)
		return
	}

	url := lib.CaseServiceValidatePath + "?service=" + loginPath + "&ticket=" + ticket
	logs.Debug("CAS验证URL: %s", url)

	serviceResponse, err := xcas.CasVersion2ServiceValidateAction(url)
	if err != nil {
		logs.Error("CAS验证请求失败: %v", err)
		ctrl.Redirect(loginCasPath, http.StatusFound)
		ctrl.StopRun()
		return
	}

	if serviceResponse.Failure != nil {
		logs.Warning("CAS验证失败: %s", serviceResponse.Failure.Message)
		ctrl.Ctx.ResponseWriter.Write([]byte(serviceResponse.Failure.Message))
		ctrl.StopRun()
	}

	if serviceResponse.Success != nil {
		kindx := ""
		for _, attribute := range serviceResponse.Success.Attributes.UserAttributes.Attributes {
			if attribute.Name == xcas.UserTypeField {
				kindx = attribute.Value
				break
			}
		}
		if kindx != kind {
			ctrl.Ctx.ResponseWriter.Write([]byte("invalid user type"))
			logs.Warning("用户类型不匹配: expected=%s, got=%s", kind, kindx)
			ctrl.StopRun()
		}

		account := serviceResponse.Success.User
		userDo := service.NewCmsAdmin()
		user, err := userDo.FindByAccount(account)
		if user == nil || user.UserName == "" || err != nil {
			logs.Warning("用户不存在，将创建新用户: account=%s, error=%v", account, err)
			user = &domain.CmsAdmin{
				UserName:   account,
				NickName:   account,
				RealName:   account,
				Email:      account + "@hnzhjypt.com",
				UserNumber: account,
				Mobile:     account,
				UserType:   xgeneric.IFF[int32](kind == "school", 2, 1),
				Enabled:    true,
				IsAudit:    1,
				Remark:     "CAS登录创建",
			}
			err := userDo.AdminSave(user)
			if err != nil {
				logs.Error("创建用户失败: %v", err)
				ctrl.Ctx.ResponseWriter.Write([]byte("创建用户失败"))
				ctrl.StopRun()
			}
		}

		user, err = userDo.FindByAccount(account)
		if err != nil {
			logs.Error("获取用户信息失败: %v", err)
			ctrl.Ctx.ResponseWriter.Write([]byte("未能获取用户信息"))
			ctrl.StopRun()
		}

		logs.Info("用户登录成功: userId=%d, account=%s", user.UserID, account)
		ctrl.SavaAdminState(user)
		userDo.LoginLog(user.UserID, account, "AdminLoginCase", "", "", "OK", ctrl.GetClientIp())
		url = "/admin/index"
		ctrl.Redirect(url, http.StatusFound)
		ctrl.StopRun()
	}
}
