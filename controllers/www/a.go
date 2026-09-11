package www

import (
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
)

const (
	WwwPrefix    = "www"
	ApiPrefix    = "api"
	WechatPrefix = "wechat"
)

func init() {
	logs.Debug("www 开始注册路由")

	web.AutoPrefix("/", &IndexController{})

	{
		web.Router("/:flag/:name", &ChannelController{}, "*:Channel")            // 频道首页
		web.Router("/:flag/:name/:category", &ChannelController{}, "*:Category") // 频道分类

		// web.Router("/:flag/:name/:category/:article_id", &ArticleController{}, "*:Index") // 文章详情
		web.Router("/article/:call_index/:article_id", &ArticleController{}, "*:Detail") // 文章详情
		web.Router("/article/search/:keyword", &ArticleController{}, "*:Search")         // 搜索页面
		web.Router("/article/search", &ArticleController{}, "*:Search")                  // 搜索页面
	}

	{
		web.Router("/:flag/topic/:name", &TopicController{}, "*:Index") // 专题
		web.Router("/:flag/tag/:name", &TagController{}, "*:Index")     // 标签
	}

	{
		// 微信公众号
		web.Router(WechatPrefix+"/mp/index ", &WechatMpController{}, "GET:Signature")
		web.Router(WechatPrefix+"/mp/index ", &WechatMpController{}, "POST:Message")
		web.Router(WechatPrefix+"/mp/tooauth2 ", &WechatMpWebAuthController{}, "*:ToOauth2")
		web.Router(WechatPrefix+"/mp/redirect_uri ", &WechatMpWebAuthController{}, "*:RedirectUri")
	}

	{
		web.Router(ApiPrefix+"/ads/get", &ApiAdsController{}, "*:Get")
		web.Router(ApiPrefix+"/ads/get/new", &ApiAdsController{}, "*:GetNew")
		web.Router(ApiPrefix+"/ads/paginate", &ApiAdsController{}, "*:Paginate")
		web.Router(ApiPrefix+"/ads/click", &ApiAdsController{}, "*:Click")

		web.Router(ApiPrefix+"/link/get", &ApiLinkController{}, "*:Get")
		web.Router(ApiPrefix+"/link/get/new", &ApiLinkController{}, "*:GetNew")
		web.Router(ApiPrefix+"/link/paginate", &ApiLinkController{}, "*:Paginate")
		web.Router(ApiPrefix+"/link/click", &ApiLinkController{}, "*:Click")

		web.Router(ApiPrefix+"/site/default", &ApiSiteController{}, "*:Default")
		web.Router(ApiPrefix+"/site/domain", &ApiSiteController{}, "*:FindDomain")
		web.Router(ApiPrefix+"/site/find", &ApiSiteController{}, "*:Find")
		web.Router(ApiPrefix+"/site/find/:site_id", &ApiSiteController{}, "*:Find")
		web.Router(ApiPrefix+"/channel/get", &ApiSiteController{}, "*:ChannelGet")
		web.Router(ApiPrefix+"/site/menu", &ApiSiteController{}, "*:Menu")
		web.Router(ApiPrefix+"/site/menu/:site_id", &ApiSiteController{}, "*:Menu")
		web.Router(ApiPrefix+"/site/menu/flag", &ApiSiteController{}, "*:MenuFlag")
		web.Router(ApiPrefix+"/site/menu/flag/:site_flag", &ApiSiteController{}, "*:MenuFlag")
		web.Router(ApiPrefix+"/site/language/list", &ApiSiteController{}, "*:SiteLanguageList")
		web.Router(ApiPrefix+"/site/language/find", &ApiSiteController{}, "*:SiteLanguageBySiteID")
		web.Router(ApiPrefix+"/site/language/by-code", &ApiSiteController{}, "*:SiteLanguageByCode")

		web.Router(ApiPrefix+"/category/nav", &ApiArticleController{}, "*:CategoryNav")
		web.Router(ApiPrefix+"/category/get", &ApiArticleController{}, "*:CategoryGet")
		web.Router(ApiPrefix+"/category/get/:channel_name", &ApiArticleController{}, "*:CategoryGet")
		web.Router(ApiPrefix+"/category/find", &ApiArticleController{}, "*:CategoryFind")
		web.Router(ApiPrefix+"/article/get", &ApiArticleController{}, "*:Get")
		web.Router(ApiPrefix+"/article/get/new", &ApiArticleController{}, "*:GetNew")
		web.Router(ApiPrefix+"/article/paginate", &ApiArticleController{}, "*:Paginate")
		web.Router(ApiPrefix+"/article/find", &ApiArticleController{}, "*:Find")
		web.Router(ApiPrefix+"/article/article", &ApiArticleController{}, "*:Article")
		web.Router(ApiPrefix+"/article/album", &ApiArticleController{}, "*:Album")
		web.Router(ApiPrefix+"/article/attach", &ApiArticleController{}, "*:Attach")
		web.Router(ApiPrefix+"/article/click", &ApiArticleController{}, "*:Click")
		web.Router(ApiPrefix+"/article/like", &ApiArticleController{}, "*:Like")
		web.Router(ApiPrefix+"/article/album/click", &ApiArticleController{}, "*:AlbumClick")
		web.Router(ApiPrefix+"/article/prev_next", &ApiArticleController{}, "*:PrevNext")
		web.Router(ApiPrefix+"/article/property", &ApiArticleController{}, "*:Property")

		web.Router(ApiPrefix+"/tag/get", &ApiTagController{}, "*:Get")
		web.Router(ApiPrefix+"/tag/get/new", &ApiTagController{}, "*:GetNew")
		web.Router(ApiPrefix+"/tag/click", &ApiTagController{}, "*:Click")
		web.Router(ApiPrefix+"/tag/article/paginate", &ApiTagController{}, "*:ArticlePaginate")

		web.Router(ApiPrefix+"/topic/get", &ApiTopicController{}, "*:Get")
		web.Router(ApiPrefix+"/topic/get/new", &ApiTopicController{}, "*:GetNew")
		web.Router(ApiPrefix+"/topic/click", &ApiTopicController{}, "*:Click")
		web.Router(ApiPrefix+"/topic/article/paginate", &ApiTopicController{}, "*:ArticlePaginate")

		web.Router(ApiPrefix+"/plg/online/register/save", &PlgOnlineRegisterController{}, "*:Save")
	}

	web.Router(ApiPrefix+"/cache/clear", &ApiCacheController{}, "*:Clear")

	web.Router("/", &IndexController{}, "*:Index")

	web.Router("/sse", &SSEController{}, "*:Message")

	logs.Debug("www 结束注册路由")

	InitWechatMpVerifyRouter()
}
