package plugin

import "haedu.gov.cn/cms/controllers/www"

type JobController struct{ www.BaseController }

func (ctrl *JobController) List() {
	ctrl.TplName = ctrl.GetView(www.DefaultSite.Template, www.DefaultSite.Template+".jobList.html")
}

func (ctrl *JobController) Detail() {
	jId, _ := ctrl.GetInt("jId", 0)
	ctrl.Data["jId"] = jId
	ctrl.TplName = ctrl.GetView(www.DefaultSite.Template, www.DefaultSite.Template+".jobDetail.html")
}
