package www

type IndexController struct{ BaseController }

func (ctrl *IndexController) Index() {
	ctrl.TplName = ctrl.GetView(DefaultSite.Template, "index.html")
	// c.display()
}
