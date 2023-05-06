package admin

import (
	"fmt"

	"github.com/beego/beego/v2/server/web"
)

var (
	GlobalAdminId   int64  // 管理员id
	GlobalAuthFlag  int    // 1:管理员 2:学校
	GlobalAdminName string // 管理员名称
	GlobalRealName  string // 管理员名称
)

const (
	AdminPreifx = "admin"
)

func init() {
	fmt.Println("admin 开始注册路由")

	web.Router("cms/admin/login", &LoginController{}, "*:AdminLogin")
	web.Router("cms/admin/logout", &LoginController{}, "*:Logout")
	web.Router("cms/admin/login/verify", &LoginController{}, "*:AdminLoginVerify")

	web.AutoPrefix(AdminPreifx, &ToolsController{})
	web.AutoPrefix(AdminPreifx, &FileController{})
	web.AutoPrefix(AdminPreifx, &UEditorController{})

	web.Router(AdminPreifx+"/index", &IndexController{}, "*:Index")

	web.Router(AdminPreifx+"/site/edit", &SiteController{}, "get:SiteEdit")
	web.Router(AdminPreifx+"/site/data", &SiteController{}, "*:SiteData")
	web.Router(AdminPreifx+"/site/delete", &SiteController{}, "*:Delete")
	web.Router(AdminPreifx+"/site/save", &SiteController{}, "*:Save")

	web.AutoPrefix(AdminPreifx, &CommonController{})
	web.AutoPrefix(AdminPreifx, &ThemeController{})
	web.AutoPrefix(AdminPreifx, &SiteController{})
	web.AutoPrefix(AdminPreifx, &IndexController{})
	web.AutoPrefix(AdminPreifx, &TagController{})
	web.AutoPrefix(AdminPreifx, &TopicController{})
	web.AutoPrefix(AdminPreifx, &LinkController{})
	web.AutoPrefix(AdminPreifx, &AdsController{})
	web.AutoPrefix(AdminPreifx, &AdminController{})
	web.AutoPrefix(AdminPreifx, &ArticleController{})
	web.AutoPrefix(AdminPreifx, &MenuController{})

	web.AutoPrefix(AdminPreifx, &WeixinController{})
	web.Router(AdminPreifx+"/weixin/message/subscribe", &WeixinController{}, "*:Subscribe")
	web.Router(AdminPreifx+"/weixin/message/default", &WeixinController{}, "*:Default")
	web.Router(AdminPreifx+"/weixin/message/text", &WeixinController{}, "*:Text")
	web.Router(AdminPreifx+"/weixin/message/picture", &WeixinController{}, "*:Picture")
	web.Router(AdminPreifx+"/weixin/message/sound", &WeixinController{}, "*:Sound")
	web.Router(AdminPreifx+"/weixin/message/response", &WeixinController{}, "*:Response")
	web.Router(AdminPreifx+"/weixin/verify", &WeixinMpVerifyController{}, "*:Index")
	web.Router(AdminPreifx+"/weixin/verify/list", &WeixinMpVerifyController{}, "*:List")
	web.Router(AdminPreifx+"/weixin/verify/edit", &WeixinMpVerifyController{}, "*:Edit")
	web.Router(AdminPreifx+"/weixin/verify/saveSortId", &WeixinMpVerifyController{}, "*:SaveSortId")
	web.Router(AdminPreifx+"/weixin/verify/destory", &WeixinMpVerifyController{}, "*:Destory")
	web.Router(AdminPreifx+"/weixin/verify/upload", &WeixinMpVerifyController{}, "*:Upload")

	web.Router(AdminPreifx+"/plg/register", &PlgOnlineRegisterController{}, "*:Index")
	web.Router(AdminPreifx+"/plg/register/paginate", &PlgOnlineRegisterController{}, "*:Paginate")
	web.Router(AdminPreifx+"/plg/register/edit", &PlgOnlineRegisterController{}, "*:Edit")
	web.Router(AdminPreifx+"/plg/register/save", &PlgOnlineRegisterController{}, "*:Save")
	web.Router(AdminPreifx+"/plg/register/destory", &PlgOnlineRegisterController{}, "*:Destory")
	web.Router(AdminPreifx+"/plg/register/read", &PlgOnlineRegisterController{}, "*:ChangeRead")

	fmt.Println("admin 结束注册路由")
}
