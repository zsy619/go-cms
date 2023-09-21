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
func (c *XsbmController) Index() {
	channelName := "zsjy"
	channelModel, err := biz.NewApiChannel().Find(channelName, 0)
	if err != nil {
		c.Ctx.WriteString(err.Error())
		c.StopRun()
	}
	c.Data["channel"] = channelModel
	c.Data["channelName"] = channelName
	c.TplName = c.GetView(www.DefaultSite.Template, "plgn.xsbm.html")
}

// Save 线上报名
// @router /plugin/xsbm/save [post]
func (c *XsbmController) Save() {
	captcha := c.GetSafeString("captcha")
	if controllers.VerifyCode(captcha) == false {
		c.JSONError("图形验证码错误")
		return
	}

	mdl := model.PlgOnlineRegister{}
	if err := c.ParseForm(&mdl); err != nil {
		logs.Error("Save", err.Error())
		c.JSONError(err.Error())
	}
	mdl.IP = c.Ctx.Input.IP()

	if err := biz.NewPlgOnlineRegister().Save(&mdl); err != nil {
		logs.Error("Save", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("在线报名成功", nil)
}
