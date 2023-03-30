package www

type IndexController struct {
	BaseController
}

func (c *IndexController) Index() {
	c.display()
}
