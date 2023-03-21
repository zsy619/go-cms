package admin

type SiteController struct{ BaseController }

// Index 站点管理列表
// @router /admin/site/index [get]
func (c *SiteController) Index() {
	c.display()
}

// Channel 站点栏目管理
// @router /admin/site/channel [get]
func (c *SiteController) Channel() {
	c.display()
}
