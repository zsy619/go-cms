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
	web.Router("cms/admin/logout", &LoginController{}, "*:Logout")
	web.Router("cms/admin/login/verify", &LoginController{}, "*:AdminLoginVerify")

	web.AutoPrefix("admin", &ToolsController{})

	web.Router("admin/index", &IndexController{}, "*:Index")

	web.Router("admin/site/edit", &SiteController{}, "get:SiteEdit")
	web.Router("admin/site/data", &SiteController{}, "*:SiteData")
	web.Router("admin/site/delete", &SiteController{}, "*:Delete")
	web.Router("admin/site/save", &SiteController{}, "*:Save")

	web.AutoPrefix("admin", &SiteController{})
	web.AutoPrefix("admin", &IndexController{})
	web.AutoPrefix("admin", &LinkController{})
	web.AutoPrefix("admin", &AdminController{})
	web.AutoPrefix("admin", &ArticleController{})
	web.AutoPrefix("admin", &MenuController{})

	web.AutoPrefix("admin", &WeixinController{})
	web.Router("admin/weixin/message/subscribe", &WeixinController{}, "*:Subscribe")
	web.Router("admin/weixin/message/default", &WeixinController{}, "*:Default")
	web.Router("admin/weixin/message/text", &WeixinController{}, "*:Text")
	web.Router("admin/weixin/message/picture", &WeixinController{}, "*:Picture")
	web.Router("admin/weixin/message/sound", &WeixinController{}, "*:Sound")
	web.Router("admin/weixin/message/response", &WeixinController{}, "*:Response")

	web.Router("/admin/plg/register", &PlgOnlineRegisterController{}, "*:Index")
	web.Router("/admin/plg/register/paginate", &PlgOnlineRegisterController{}, "*:Paginate")
	web.Router("/admin/plg/register/edit", &PlgOnlineRegisterController{}, "*:Edit")
	web.Router("/admin/plg/register/save", &PlgOnlineRegisterController{}, "*:Save")
	web.Router("/admin/plg/register/destory", &PlgOnlineRegisterController{}, "*:Destory")
	web.Router("/admin/plg/register/read", &PlgOnlineRegisterController{}, "*:ChangeRead")
}
