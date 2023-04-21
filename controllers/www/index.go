package www

type IndexController struct {
	BaseController
}

func (c *IndexController) Index() {
	c.TplName = c.getView(DefatulSite.DirPath, "index.html")
	// c.display()
}
