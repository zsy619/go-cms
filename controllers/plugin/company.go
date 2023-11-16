package plugin

import "haedu.gov.cn/cms/controllers/www"

type CompanyController struct{ www.BaseController }

func (ctrl *CompanyController) List() {
	ctrl.TplName = ctrl.GetView(www.DefaultSite.Template, www.DefaultSite.Template+".companyList.html")
}

func (ctrl *CompanyController) Detail() {
	cId, _ := ctrl.GetInt("cId", 0)
	ctrl.Data["cId"] = cId
	ctrl.TplName = ctrl.GetView(www.DefaultSite.Template, www.DefaultSite.Template+".companyDetail.html")
}
