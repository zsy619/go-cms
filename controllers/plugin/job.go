package plugin

import "haedu.gov.cn/cms/controllers/www"

type JobController struct {
	www.BaseController
}

func (c *JobController) List() {
	c.TplName = c.GetView(www.DefaultSite.Template, www.DefaultSite.Template+".jobList.html")
}

func (c *JobController) Detail() {
	jId, _ := c.GetInt("jId", 0)
	c.Data["jId"] = jId
	c.TplName = c.GetView(www.DefaultSite.Template, www.DefaultSite.Template+".jobDetail.html")
}
