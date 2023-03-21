package admin

type LinkController struct{ BaseController }

// Index 链接管理
// @router /admin/link/index [get]
func (c *LinkController) Index() {
	c.display()
}

// List 链接列表
func (c *LinkController) List() {
}

// Category 链接分类
// @router /admin/link/category [get]
func (c *LinkController) Category() {
	c.display()
}
