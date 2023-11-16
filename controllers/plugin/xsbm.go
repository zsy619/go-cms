package plugin

import (
	"github.com/beego/beego/v2/core/logs"

	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/controllers"
	"haedu.gov.cn/cms/controllers/www"
)

type XsbmController struct{ www.BaseController }

// Index 线上报名首页
// @router /plugin/xsbm/index
func (ctrl *XsbmController) Index() {
	channelName := "zsjy"
	channelModel, err := biz.NewApiChannel().Find(channelName, 0)
	if err != nil {
		ctrl.Ctx.WriteString(err.Error())
		ctrl.StopRun()
	}
	ctrl.Data["channel"] = channelModel
	ctrl.Data["channelName"] = channelName
	ctrl.TplName = ctrl.GetView(www.DefaultSite.Template, "plgn.xsbm.html")
}

// Save 线上报名
// @router /plugin/xsbm/save [post]
func (ctrl *XsbmController) Save() {
	captcha := ctrl.GetSafeString("captcha")
	if !controllers.VerifyCode(captcha) {
		ctrl.JSONError("图形验证码错误")
		return
	}

	mdl := model.PlgOnlineRegister{}
	if err := ctrl.ParseForm(&mdl); err != nil {
		logs.Error("Save", err.Error())
		ctrl.JSONError(err.Error())
	}
	mdl.IP = ctrl.Ctx.Input.IP()

	if err := biz.NewPlgOnlineRegister().Save(&mdl); err != nil {
		logs.Error("Save", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("在线报名成功", nil)
}
