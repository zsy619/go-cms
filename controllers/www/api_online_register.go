package www

import (
	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/controllers"
)

type ApiOnlineRegisterController struct{ BaseController }

// @router /api/online/register [get]
func (this *ApiOnlineRegisterController) Save() {
	captcha := this.GetString("captcha")
	if controllers.VerifyCode(captcha) == false {
		this.JSONError("图形验证码错误")
		return
	}

	mdl := model.PlgOnlineRegister{}
	if err := this.ParseForm(&mdl); err != nil {
		logs.Error("Save", err.Error())
		this.JSONError(err.Error())
	}

	if err := biz.NewPlgOnlineRegister().Save(&mdl); err != nil {
		logs.Error("Save", err.Error())
		this.JSONError(err.Error())
		return
	}
	this.JSONSuccess("在线报名成功", nil)
}
