package controllers

import "github.com/beego/beego/v2/server/web"

func init() {
	web.Router("/captcha", &CaptchaController{}, "get:GenerateHandler")     // 生成图形验证码
	web.Router("/captcha/verify", &CaptchaController{}, "get:VerifyHandle") // 验证

	web.Router("error/404", &ErrorController{}, "*:Error404")
	web.Router("error/500", &ErrorController{}, "*:Error500")

	web.Router("/help", &HelpController{}, "*:Index")
	web.Router("/help/:chn/:page", &HelpController{}, "*:Index")
}
