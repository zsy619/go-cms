package www

type IndexController struct {
	BaseController
}

func (c *IndexController) Index() {
	c.TplName = c.GetView(DefatulSite.Template, "index.html")
	// c.display()
}
