package admin

import (
	"fmt"

	"github.com/beego/beego/v2/server/web"
	"haedu.gov.cn/cms/controllers"
	"haedu.gov.cn/cms/controllers/admin/model"
)

type LoginController struct{ web.Controller }

func (c *LoginController) AdminLogin() {
	c.Data["captcha"] = "/captcha"
	c.TplName = "admin/login/login.html"
}

func (c *LoginController) AdminLoginVerify() {
	result := model.LoginResult{
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
	// ssoUser := sso.NewSsoUserModel(models.DB_SSO)
	// user, err := ssoUser.Login(username, password, 0, sso.LoginAll)
	// if err != nil {
	// 	fmt.Println("登录错误：", err.Error())
	// 	result.Code = 2
	// 	result.Message = "账号密码错误"
	// 	c.Data["json"] = &result
	// 	c.ServeJSON()
	// 	return
	// }
	// c.SetSession("a.id", strconv.FormatInt(user.Id, 10))
	// c.SetSession("a.user_name", username)
	// c.SetSession("a.real_name", user.RealName)
	// c.SetSession("a.id_card", user.IdCard)
	// c.SetSession("a.sex", strconv.Itoa(user.Sex))
	// c.SetSession("a.user_type", strconv.Itoa(user.UserType))
	result.Url = "/admin/index"
	c.Data["json"] = &result
	c.ServeJSON()
	return
}
