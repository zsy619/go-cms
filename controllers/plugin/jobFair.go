package plugin

import "haedu.gov.cn/cms/controllers/www"

type JobfairController struct{ www.BaseController }

func (ctrl *JobfairController) List() {
	ctrl.TplName = ctrl.GetView(www.DefaultSite.Template, www.DefaultSite.Template+".jobFairList.html")
}

func (ctrl *JobfairController) Detail() {
	jId, _ := ctrl.GetInt("jId", 0)
	ctrl.Data["jId"] = jId
	ctrl.TplName = ctrl.GetView(www.DefaultSite.Template, www.DefaultSite.Template+".jobFairDetail.html")
}

func (ctrl *JobfairController) DateList() {
	ctrl.TplName = ctrl.GetView(www.DefaultSite.Template, www.DefaultSite.Template+".dateList.html")
}
