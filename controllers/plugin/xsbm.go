package plugin

import (
	"haedu.gov.cn/cms/app/biz"
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
	c.TplName = c.GetView(www.DefatulSite.Template, "plgn.xsbm.html")
}
