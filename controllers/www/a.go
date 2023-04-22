package www

import "github.com/beego/beego/v2/server/web"

func init() {
	web.AutoPrefix("/", &IndexController{})

	{
		// 微信公众号
		web.Router("/weixin/mp/index ", &WechatMpController{}, "GET:Signature")
		web.Router("/weixin/mp/index ", &WechatMpController{}, "POST:Message")
		web.Router("/wechat/mp/tooauth2 ", &MpWebAuthController{}, "*:ToOauth2")
		web.Router("/wechat/mp/redirect_uri ", &MpWebAuthController{}, "*:RedirectUri")
	}

	{
		web.Router("/api/link/find", &ApiLinkController{}, "*:Find")
		web.Router("/api/link/paginate", &ApiLinkController{}, "*:Paginate")
		web.Router("/api/link/click", &ApiLinkController{}, "*:Click")

		web.Router("/api/site/default", &ApiSiteController{}, "*:Default")
		web.Router("/api/site/get", &ApiSiteController{}, "*:Get")
		web.Router("/api/channel/find", &ApiSiteController{}, "*:ChannelFind")
		web.Router("/api/site/menu", &ApiSiteController{}, "*:Menu")

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
		web.Router("/api/article/album/click", &ApiArticleController{}, "*:AlbumClick")
		web.Router("/api/article/prev_next", &ApiArticleController{}, "*:PrevNext")

		web.Router("/api/ads/find", &ApiAdsController{}, "*:Find")
		web.Router("/api/ads/find/new", &ApiAdsController{}, "*:FindNew")
		web.Router("/api/ads/paginate", &ApiAdsController{}, "*:Paginate")
		web.Router("/api/ads/click", &ApiAdsController{}, "*:Click")

		web.Router("/api/tag/find", &ApiTagController{}, "*:Find")
		web.Router("/api/tag/find/new", &ApiTagController{}, "*:FindNew")
		web.Router("/api/tag/click", &ApiTagController{}, "*:Click")
		web.Router("/api/tag/article/paginate", &ApiTagController{}, "*:ArticlePaginate")
	}

	web.Router("/api/online/register", &ApiOnlineRegisterController{}, "*:Save")

	web.Router("/api/cache/clear", &ApiCacheController{}, "*:Clear")

	web.Router("/", &IndexController{}, "*:Index")

	{
		web.Router("/:flag/:name", &ChannelController{}, "*:Channel")            // 频道首页
		web.Router("/:flag/:name/:category", &ChannelController{}, "*:Category") // 频道分类

		// web.Router("/:name", &ChannelController{}, "*:Index")                      // 频道首页
		// web.Router("/:name/:category", &ChannelController{}, "*:Category")         // 频道分类
	}

	{
		web.Router("/topic/:name", &TopicController{}, "*:Index") // 专题
	}

	// web.Include(&ArticleController{})

	web.Router("/article/detail/:article_id", &ArticleController{}, "*:Detail")   // 文章详情
	web.Router("/article/teacher/:article_id", &ArticleController{}, "*:Teacher") // 教师风采
	web.Router("/article/search/:keyword", &ArticleController{}, "*:Search")      // 搜索页面
	web.Router("/article/search", &ArticleController{}, "*:Search")               // 搜索页面

	web.AutoPrefix("/", &XxgkController{}) // 学校概况
	web.AutoPrefix("/", &ZsjyController{}) // 招生就业
	web.AutoPrefix("/", &XyzxController{}) // 校园资讯
	web.AutoPrefix("/", &JyjxController{}) // 教育教学
	web.AutoPrefix("/", &LxwmController{}) // 联系我们
	web.AutoPrefix("/", &JgszController{}) // 机构设置
}
