package www

import "github.com/beego/beego/v2/server/web"

func init() {
	web.AutoPrefix("/", &IndexController{})
	web.AutoPrefix("/", &WeixinController{})

	web.Router("/api/link/find", &ApiLinkController{}, "*:Find")
	web.Router("/api/link/paginate", &ApiLinkController{}, "*:Paginate")
	web.Router("/api/link/click", &ApiLinkController{}, "*:Click")

	web.Router("/api/site/default", &ApiSiteController{}, "*:Default")
	web.Router("/api/site/get", &ApiSiteController{}, "*:Get")
	web.Router("/api/channel/find", &ApiSiteController{}, "*:ChannelFind")

	web.Router("/api/category/nav", &ApiArticleController{}, "*:CategoryNav")
	web.Router("/api/category/find", &ApiArticleController{}, "*:CategoryFind")
	web.Router("/api/category/one", &ApiArticleController{}, "*:CategoryOne")
	web.Router("/api/article/find", &ApiArticleController{}, "*:Find")
	web.Router("/api/article/paginate", &ApiArticleController{}, "*:Paginate")
	web.Router("/api/article/one", &ApiArticleController{}, "*:One")
	web.Router("/api/article/article", &ApiArticleController{}, "*:Article")
	web.Router("/api/article/album", &ApiArticleController{}, "*:Album")
	web.Router("/api/article/attach", &ApiArticleController{}, "*:Attach")
	web.Router("/api/article/click", &ApiArticleController{}, "*:Click")
	web.Router("/api/article/like", &ApiArticleController{}, "*:Like")

	web.Router("/api/online/register", &ApiOnlineRegisterController{}, "*:Save")

	web.Router("/api/cache/clear", &ApiCacheController{}, "*:Clear")

	web.Router("/", &IndexController{}, "*:Index")

	web.AutoPrefix("/", &XxgkController{}) // 学校概况
	web.AutoPrefix("/", &ZsjyController{}) // 招生就业
	web.AutoPrefix("/", &XyzxController{}) // 校园资讯
	web.AutoPrefix("/", &JyjxController{}) // 教育教学
	web.AutoPrefix("/", &LxwmController{}) // 联系我们
}
