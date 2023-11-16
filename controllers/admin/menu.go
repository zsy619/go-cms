package admin

import "haedu.gov.cn/cms/app/biz"

type MenuController struct{ BaseController }

// MenuList 登录者导航
// @router /admin/menu/MenuList [get]
func (ctrl *MenuController) MenuList() {
	menuServicee := biz.NewMenu()
	menuList := menuServicee.MenuList(GlobalAdminId)
	ctrl.Data["json"] = menuList
	ctrl.ServeJSON()
}
