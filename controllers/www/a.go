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

	web.Router("/api/category/find", &ApiArticleController{}, "*:CategoryFind")
	web.Router("/api/article/find", &ApiArticleController{}, "*:Find")
	web.Router("/api/article/paginate", &ApiArticleController{}, "*:Paginate")
	web.Router("/api/article/one", &ApiArticleController{}, "*:One")
	web.Router("/api/article/article", &ApiArticleController{}, "*:Article")
	web.Router("/api/article/album", &ApiArticleController{}, "*:Album")
	web.Router("/api/article/attach", &ApiArticleController{}, "*:Attach")
	web.Router("/api/article/click", &ApiArticleController{}, "*:Click")
	web.Router("/api/article/like", &ApiArticleController{}, "*:Like")

	web.Router("/api/online/register", &ApiOnlineRegisterController{}, "*:Find")

	web.Router("/api/cache/clear", &ApiCacheController{}, "*:Clear")

	web.Router("/", &IndexController{}, "*:Index")
}
