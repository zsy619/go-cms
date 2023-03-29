package www

import "github.com/beego/beego/v2/server/web"

func init() {
	web.AutoPrefix("/", &IndexController{})
	web.AutoPrefix("/", &WeixinController{})

	web.Router("/api/link/find", &ApiLinkController{}, "*:Find")
	web.Router("/api/link/paginate", &ApiLinkController{}, "*:Paginate")

	web.Router("/api/site/get", &ApiSiteController{}, "*:Get")
	web.Router("/api/channel/find", &ApiSiteController{}, "*:ChannelFind")
	web.Router("/api/category/find", &ApiArticleController{}, "*:CategoryFind")
	web.Router("/api/article/find", &ApiArticleController{}, "*:Find")
	web.Router("/api/article/paginate", &ApiArticleController{}, "*:Paginate")
}
