package admin

import (
	"fmt"

	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/controllers"
	"haedu.gov.cn/cms/controllers/admin/vmodel"
)

type LoginController struct {
	controllers.BaseController
}

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
	GlobalAuthFlag = 1 // 1:管理员 2:学校
	GlobalAdminName = user.UserName
	GlobalRealName = user.RealName
}

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
