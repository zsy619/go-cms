package admin

import (
	"fmt"
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

// AdminLogin 管理员登录
// @router cms/admin/login [get]
func (ctrl *LoginController) AdminLogin() {
	ctrl.Data["captcha"] = "/captcha"
	ctrl.TplName = "admin/login/login.html"
}

func (ctrl *LoginController) SavaAdminState(user *domain.CmsAdmin) {
	ctrl.SetSession("adminId", user.UserID)
	ctrl.SetSession("adminAccount", user.UserName)
	ctrl.SetSession("adminName", user.UserName)
	ctrl.SetSession("realName", user.RealName)
	ctrl.SetSession("adminState", 1)
	ctrl.SetSession("adminLevel", 1)
	ctrl.SetSession("user", user)

	GlobalAdminId = user.UserID
	GlobalUserType = int(user.UserType) // 1:管理员 2:学校
	GlobalAdminName = user.UserName
	GlobalRealName = user.RealName
	GlobalIsAudit = user.IsAudit
	GlobalRoleId = user.RoleID
	GlobalRoleType = user.RoleType
}

// AdminLoginVerify 管理员登录验证
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
	fmt.Println(username, password, captcha)
	adminDo := service.NewCmsAdmin()
	user, err := adminDo.Login(username, password, 0, service.LoginAll)
	if err != nil {
		fmt.Println("登录错误：", err.Error())
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

// Logout 退出登录
// @router cms/admin/logout [get]
func (ctrl *LoginController) Logout() {
	ctrl.DestroySession()
	ctrl.Redirect("/admin/login", 302)
}

// School学校登录
func (ctrl *LoginController) School() {
	ctrl.login(lib.LoginCasPathOfSchool, lib.LoginPathOfSchool, "school")
}

// Admin 管理员登录
func (ctrl *LoginController) Admin() {
	ctrl.login(lib.LoginCasPathOfAdmin, lib.LoginPathOfAdmin, "admin")
}

func (ctrl *LoginController) login(loginCasPath, loginPath string, kind string) {
	// TOD：20220418 登录类型
	xcache.SetDiskvString("cas_login", kind)
	GlobalAuthFlag = kind
	ticket := ctrl.GetSafeString("ticket")
	fmt.Println("ticket: ", ticket)
	if ticket == "" {
		ctrl.Redirect(loginCasPath, http.StatusFound)
		return
	}
	url := lib.CaseServiceValidatePath + "?service=" + loginPath + "&ticket=" + ticket
	logs.Debug("--->", url)
	fmt.Println("--->", url)
	serviceResponse, err := xcas.CasVersion2ServiceValidateAction(url)
	if err != nil {
		logs.Error(err)
		fmt.Println("err:", err)
		ctrl.Redirect(loginCasPath, http.StatusFound)
		ctrl.StopRun()
		return
	}
	// fmt.Println("--------------------------------------------------------:::", serviceResponse)
	if serviceResponse.Failure != nil {
		fmt.Println("error: ", serviceResponse.Failure.Message)
		ctrl.Ctx.ResponseWriter.Write([]byte(serviceResponse.Failure.Message))
		ctrl.StopRun()
	}
	if serviceResponse.Success != nil {
		// 类型检查
		{
			kindx := ""
			for _, attribute := range serviceResponse.Success.Attributes.UserAttributes.Attributes {
				if attribute.Name == xcas.UserTypeField {
					kindx = attribute.Value
					break
				}
			}
			if kindx != kind {
				ctrl.Ctx.ResponseWriter.Write([]byte("invalid user type"))
				fmt.Println("invalid user type")
				ctrl.StopRun()
			}
		}

		account := serviceResponse.Success.User
		// 判断用户是否存在，不存在则创建
		userDo := service.NewCmsAdmin()
		user, err := userDo.FindByAccount(account)
		if user == nil || user.UserName == "" || err != nil {
			fmt.Println("error: ", err.Error())
			// 创建用户
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
				fmt.Println("error: ", err.Error())
				ctrl.Ctx.ResponseWriter.Write([]byte("创建用户失败"))
				ctrl.StopRun()
			}
		}
		user, err = userDo.FindByAccount(account)
		if err != nil {
			ctrl.Ctx.ResponseWriter.Write([]byte("未能获取用户信息"))
			ctrl.StopRun()
		}
		fmt.Println("user: ", user)
		// 保存登录状态
		ctrl.SavaAdminState(user)
		userDo.LoginLog(user.UserID, account, "AdminLoginCase", "", "", "OK", ctrl.GetClientIp())
		url = "/admin/index"
		ctrl.Redirect(url, http.StatusFound)
		ctrl.StopRun()
		return
	}
}
