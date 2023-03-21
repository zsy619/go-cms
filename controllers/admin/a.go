package admin

import "github.com/beego/beego/v2/server/web"

var (
	GlobalAdminId   int64  // 管理员id
	GlobalAuthFlag  int    // 1:管理员 2:学校
	GlobalAdminName string // 管理员名称
	GlobalRealName  string // 管理员名称
)

func init() {
	web.Router("cms/admin/login", &LoginController{}, "*:AdminLogin")
	web.Router("cms/admin/login/verify", &LoginController{}, "*:AdminLoginVerify")

	web.Router("admin/index", &IndexController{}, "*:Index")
	web.Router("admin/site/index", &SiteController{}, "*:Index")
	web.Router("admin/site/channel", &SiteController{}, "*:Channel")

	web.AutoPrefix("admin", &LinkController{})
}
