package plugin

import "haedu.gov.cn/cms/controllers/www"

type AirkeynoteController struct {
	www.BaseController
}

func (c *AirkeynoteController) List() {
	c.TplName = c.GetView(www.DefaultSite.Template, www.DefaultSite.Template+".airKeynoteList.html")
}

func (c *AirkeynoteController) Detail() {
	aId, _ := c.GetInt("aId", 0)
	c.Data["aId"] = aId
	c.TplName = c.GetView(www.DefaultSite.Template, www.DefaultSite.Template+".airKeynoteDetail.html")
}
