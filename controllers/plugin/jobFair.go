package plugin

import "haedu.gov.cn/cms/controllers/www"

type JobfairController struct {
	www.BaseController
}

func (c *JobfairController) List() {
	c.TplName = c.GetView(www.DefaultSite.Template, www.DefaultSite.Template+".jobFairList.html")
}

func (c *JobfairController) Detail() {
	jId, _ := c.GetInt("jId", 0)
	c.Data["jId"] = jId
	c.TplName = c.GetView(www.DefaultSite.Template, www.DefaultSite.Template+".jobFairDetail.html")
}

func (c *JobfairController) DateList() {
	c.TplName = c.GetView(www.DefaultSite.Template, www.DefaultSite.Template+".dateList.html")
}
