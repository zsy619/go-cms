package plugin

import "haedu.gov.cn/cms/controllers/www"

type AirkeynoteController struct{ www.BaseController }

func (ctrl *AirkeynoteController) List() {
	ctrl.TplName = ctrl.GetView(www.DefaultSite.Template, www.DefaultSite.Template+".airKeynoteList.html")
}

func (ctrl *AirkeynoteController) Detail() {
	aId, _ := ctrl.GetInt("aId", 0)
	ctrl.Data["aId"] = aId
	ctrl.TplName = ctrl.GetView(www.DefaultSite.Template, www.DefaultSite.Template+".airKeynoteDetail.html")
}
