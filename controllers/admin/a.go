package admin

import (
	"fmt"

	"github.com/beego/beego/v2/server/web"

	"haedu.gov.cn/cms/app/lib"
)

var (
	GlobalAdminId    int64  // 管理员id
	GlobalAuthFlag   string // admin:管理员 school:学校
	GlobalUserType   int    // 用户类型
	GlobalAdminName  string // 管理员名称
	GlobalRealName   string // 管理员名称
	GlobalSchoolName string // 学校名称
	GlobalIsAudit    int32  // 是否审核
	GlobalRoleId     int64  // 管理员角色id
	GlobalRoleType   string // 管理员角色类型
)

const (
	AdminPreifx = "admin"
)

func init() {
	fmt.Println("admin 开始注册路由")

	loginNs := web.NewNamespace("cms",
		web.NSNamespace("admin",
			web.NSRouter("login", &LoginController{}, "*:AdminLogin"),
			web.NSRouter("logout", &LoginController{}, "*:Logout"),
			web.NSRouter("login/verify", &LoginController{}, "*:AdminLoginVerify"),
		),
	)
	web.AddNamespace(loginNs)

	adminNs := web.NewNamespace(AdminPreifx,
		web.NSAutoRouter(&ToolsController{}),
		web.NSAutoRouter(&FileController{}),
		web.NSAutoRouter(&UEditorController{}),
		web.NSRouter("index", &IndexController{}, "*:Index"),

		web.NSNamespace("site",
			web.NSRouter("edit", &SiteController{}, "get:SiteEdit"),
			web.NSRouter("data", &SiteController{}, "get:SiteData"),
			web.NSRouter("delete", &SiteController{}, "get:Delete"),
			web.NSRouter("save", &SiteController{}, "get:Save"),
		),

		web.NSAutoRouter(&CommonController{}),
		web.NSAutoRouter(&ThemeController{}),
		web.NSAutoRouter(&SiteController{}),
		web.NSAutoRouter(&IndexController{}),
		web.NSAutoRouter(&TagController{}),
		web.NSAutoRouter(&TopicController{}),
		web.NSAutoRouter(&LinkController{}),
		web.NSAutoRouter(&AdsController{}),
		web.NSAutoRouter(&AdminController{}),
		web.NSAutoRouter(&ArticleController{}),
		web.NSAutoRouter(&MenuController{}),
		web.NSAutoRouter(&NoticeController{}),
		web.NSAutoRouter(&NavController{}),
		web.NSAutoRouter(&CacheController{}),
		web.NSAutoRouter(&WeixinController{}),

		web.NSNamespace("weixin",
			web.NSRouter("/message/subscribe", &WeixinController{}, "*:Subscribe"),
			web.NSRouter("/message/default", &WeixinController{}, "*:Default"),
			web.NSRouter("/message/text", &WeixinController{}, "*:Text"),
			web.NSRouter("/message/picture", &WeixinController{}, "*:Picture"),
			web.NSRouter("/message/sound", &WeixinController{}, "*:Sound"),
			web.NSRouter("/message/response", &WeixinController{}, "*:Response"),
			web.NSRouter("/verify", &WeixinMpVerifyController{}, "*:Index"),
			web.NSRouter("/verify/list", &WeixinMpVerifyController{}, "*:List"),
			web.NSRouter("/verify/edit", &WeixinMpVerifyController{}, "*:Edit"),
			web.NSRouter("/verify/savesortid", &WeixinMpVerifyController{}, "*:SaveSortId"),
			web.NSRouter("/verify/destory", &WeixinMpVerifyController{}, "*:Destory"),
			web.NSRouter("/verify/upload", &WeixinMpVerifyController{}, "*:Upload"),
			web.NSRouter("/verify/refrshcache", &WeixinMpVerifyController{}, "*:RefrshCache"),
		),

		web.NSNamespace("plg",
			web.NSRouter("/register", &PlgOnlineRegisterController{}, "*:Index"),
			web.NSRouter("/register/paginate", &PlgOnlineRegisterController{}, "*:Paginate"),
			web.NSRouter("/register/edit", &PlgOnlineRegisterController{}, "*:Edit"),
			web.NSRouter("/register/save", &PlgOnlineRegisterController{}, "*:Save"),
			web.NSRouter("/register/destory", &PlgOnlineRegisterController{}, "*:Destory"),
			web.NSRouter("/register/read", &PlgOnlineRegisterController{}, "*:ChangeRead"),
		),
	)
	web.AddNamespace(adminNs)

	web.Router(lib.Url_School_Login, &LoginController{}, "*:School")
	web.Router(lib.Url_Admin_Login, &LoginController{}, "*:Admin")

	fmt.Println("admin 结束注册路由")
}
