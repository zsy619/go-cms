package admin

import (
	"fmt"
	"net/http"

	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/app/lib"
	"haedu.gov.cn/cms/controllers"
	"haedu.gov.cn/cms/controllers/admin/vmodel"
	"haedu.gov.cn/tools/xcache"
	"haedu.gov.cn/tools/xcas"
	"haedu.gov.cn/tools/xgeneric"
)

type LoginController struct {
	controllers.BaseController
}

// AdminLogin 管理员登录
// @router cms/admin/login [get]
func (c *LoginController) AdminLogin() {
	c.Data["captcha"] = "/captcha"
	c.TplName = "admin/login/login.html"
}

func (c *LoginController) SavaAdminState(user *model.CmsAdmin) {
	c.SetSession("adminId", user.UserID)
	c.SetSession("adminAccount", user.UserName)
	c.SetSession("adminName", user.UserName)
	c.SetSession("realName", user.RealName)
	c.SetSession("adminState", 1)
	c.SetSession("adminLevel", 1)
	c.SetSession("user", user)

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
func (c *LoginController) AdminLoginVerify() {
	result := vmodel.LoginResult{
		Code:    0,
		Message: "登陆完成，载入中...",
		Url:     "",
	}

	captcha := c.GetString("captcha")
	if controllers.VerifyCode(captcha) == false {
		result.Code = 1
		result.Message = "验证码错误"
		c.Data["json"] = &result
		c.ServeJSON()
		return
	}

	username := c.GetString("username")
	password := c.GetString("password")
	fmt.Println(username, password, captcha)
	adminDo := biz.NewCmsAdmin()
	user, err := adminDo.Login(username, password, 0, biz.LoginAll)
	if err != nil {
		fmt.Println("登录错误：", err.Error())
		result.Code = 2
		result.Message = "账号密码错误"
		c.Data["json"] = &result
		c.ServeJSON()
		return
	}
	c.SavaAdminState(user)
	adminDo.LoginLog(user.UserID, username, "AdminLoginVerify", "", "", "OK", c.GetClientIp())
	result.Url = "/admin/index"
	c.Data["json"] = &result
	c.ServeJSON()
	return
}

// Logout 退出登录
// @router cms/admin/logout [get]
func (c *LoginController) Logout() {
	c.DestroySession()
	c.Redirect("/admin/login", 302)
}

// School学校登录
func (c *LoginController) School() {
	c.login(lib.LoginCasPathOfSchool, lib.LoginPathOfSchool, "school")
}

// Admin 管理员登录
func (c *LoginController) Admin() {
	c.login(lib.LoginCasPathOfAdmin, lib.LoginPathOfAdmin, "admin")
}

func (this *LoginController) login(loginCasPath, loginPath string, kind string) {
	// TOD：20220418 登录类型
	xcache.SetDiskvString("cas_login", kind)
	GlobalAuthFlag = kind
	ticket := this.GetString("ticket")
	fmt.Println("ticket: ", ticket)
	if ticket == "" {
		this.Redirect(loginCasPath, http.StatusFound)
		return
	}
	url := lib.CaseServiceValidatePath + "?service=" + loginPath + "&ticket=" + ticket
	logs.Debug("--->", url)
	fmt.Println("--->", url)
	serviceResponse, err := xcas.CasVersion2ServiceValidateAction(url)
	if err != nil {
		logs.Error(err)
		fmt.Println("err:", err)
		this.Redirect(loginCasPath, http.StatusFound)
		this.StopRun()
		return
	}
	// fmt.Println("--------------------------------------------------------:::", serviceResponse)
	if serviceResponse.Failure != nil {
		fmt.Println("error: ", serviceResponse.Failure.Message)
		this.Ctx.ResponseWriter.Write([]byte(serviceResponse.Failure.Message))
		this.StopRun()
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
				this.Ctx.ResponseWriter.Write([]byte("invalid user type"))
				fmt.Println("invalid user type")
				this.StopRun()
			}
		}

		account := serviceResponse.Success.User
		// 判断用户是否存在，不存在则创建
		userDo := biz.NewCmsAdmin()
		user, err := userDo.FindByAccount(account)
		if user == nil || user.UserName == "" || err != nil {
			fmt.Println("error: ", err.Error())
			// 创建用户
			user = &model.CmsAdmin{
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
				this.Ctx.ResponseWriter.Write([]byte("创建用户失败"))
				this.StopRun()
			}
		}
		user, err = userDo.FindByAccount(account)
		if err != nil {
			this.Ctx.ResponseWriter.Write([]byte("未能获取用户信息"))
			this.StopRun()
		}
		fmt.Println("user: ", user)
		// 保存登录状态
		this.SavaAdminState(user)
		userDo.LoginLog(user.UserID, account, "AdminLoginCase", "", "", "OK", this.GetClientIp())
		url = "/admin/index"
		this.Redirect(url, http.StatusFound)
		this.StopRun()
		return
	}
}
