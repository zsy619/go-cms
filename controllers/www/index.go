package www

type IndexController struct {
	BaseController
}

func (c *IndexController) Index() {
	c.TplName = c.getView(DefatulSite.Template, "index.html")
	// c.display()
}
