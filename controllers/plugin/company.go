package plugin

import "haedu.gov.cn/cms/controllers/www"

type CompanyController struct {
	www.BaseController
}

func (c *CompanyController) List() {
	c.TplName = c.GetView(www.DefaultSite.Template, www.DefaultSite.Template+".companyList.html")
}

func (c *CompanyController) Detail() {
	cId, _ := c.GetInt("cId", 0)
	c.Data["cId"] = cId
	c.TplName = c.GetView(www.DefaultSite.Template, www.DefaultSite.Template+".companyDetail.html")
}
