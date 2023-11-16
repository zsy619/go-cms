package admin

import (
	"haedu.gov.cn/cms/app/lib"
)

type CacheController struct{ BaseController }

func (ctrl *CacheController) Index() {
	ctrl.display()
}

// GetList 获取缓存列表
func (ctrl *CacheController) GetList() {
	items := []*lib.CacheItemModel{
		// 导航相关缓存
		{Key: lib.NavGetCache.Key, Len: lib.NavGetCache.Length() + lib.NavGetByFlagCache.Length() + lib.NavCategoryGetCache.Length(), Note: lib.NavGetCache.Note, Expired: lib.NavGetCache.Expired},
		// 站点相关缓存
		{Key: lib.SiteCache.Key, Len: lib.SiteCache.Length() + lib.SiteFindCache.Length(), Note: lib.SiteCache.Note, Expired: lib.SiteCache.Expired},
		// 频道相关缓存
		{Key: lib.ChannelCache.Key, Len: lib.ChannelCache.Length() + lib.ChannelGetCache.Length(), Note: lib.ChannelCache.Note, Expired: lib.ChannelCache.Expired},
		// 栏目相关缓存
		{Key: lib.CategoryNavCache.Key, Len: lib.CategoryNavCache.Length() + lib.CategoryGetCache.Length() + lib.CategoryFindCache.Length(), Note: lib.CategoryNavCache.Note, Expired: lib.CategoryNavCache.Expired},
		// 内容相关缓存
		{Key: lib.ArticleCache.Key, Len: lib.ArticleCache.Length(), Note: lib.ArticleCache.Note + lib.ArticleGetNewCache.Note, Expired: lib.ArticleCache.Expired},
		// 系统通知公告缓存
		{Key: lib.NoticeCache.Key, Len: lib.NoticeCache.Length(), Note: lib.NoticeCache.Note, Expired: lib.NoticeCache.Expired},
		// 广告缓存
		{Key: lib.AdsCache.Key, Len: lib.AdsCache.Length(), Note: lib.AdsCache.Note, Expired: lib.AdsCache.Expired},
		// 链接缓存
		{Key: lib.LinkCache.Key, Len: lib.LinkCache.Length(), Note: lib.LinkCache.Note, Expired: lib.LinkCache.Expired},
		// 标签缓存
		{Key: lib.TagCache.Key, Len: lib.TagCache.Length() + lib.TagFindCache.Length() + lib.TagArticleCache.Length(), Note: lib.TagCache.Note, Expired: lib.TagCache.Expired},
		// 专题缓存
		{Key: lib.TopicCache.Key, Len: lib.TopicCache.Length() + lib.TopicFindCache.Length() + lib.TopicArticleCache.Length(), Note: lib.TopicCache.Note, Expired: lib.TopicCache.Expired},
		// 微信缓存
		{Key: lib.WechatAccountCache.Key, Len: lib.WechatAccountCache.Length() + lib.WechatVerifyCache.Length(), Note: lib.WechatAccountCache.Note, Expired: lib.WechatAccountCache.Expired},

		/*{Key: lib.CategoryNavCache.Key, Len: lib.CategoryNavCache.Length(), Note: lib.CategoryNavCache.Note, Expired: lib.CategoryNavCache.Expired},
		{Key: lib.CategoryGetCache.Key, Len: lib.CategoryGetCache.Length(), Note: lib.CategoryGetCache.Note, Expired: lib.CategoryGetCache.Expired},
		{Key: lib.CategoryFindCache.Key, Len: lib.CategoryFindCache.Length(), Note: lib.CategoryFindCache.Note, Expired: lib.CategoryFindCache.Expired},
		{Key: lib.ArticleCache.Key, Len: lib.ArticleCache.Length(), Note: lib.ArticleCache.Note, Expired: lib.ArticleCache.Expired},
		{Key: lib.ArticleGetNewCache.Key, Len: lib.ArticleGetNewCache.Length(), Note: lib.ArticleGetNewCache.Note, Expired: lib.ArticleGetNewCache.Expired},
		{Key: lib.NoticeCache.Key, Len: lib.NoticeCache.Length(), Note: lib.NoticeCache.Note, Expired: lib.NoticeCache.Expired},
		{Key: lib.AdsCache.Key, Len: lib.AdsCache.Length(), Note: lib.AdsCache.Note, Expired: lib.AdsCache.Expired},
		{Key: lib.LinkCache.Key, Len: lib.LinkCache.Length(), Note: lib.LinkCache.Note, Expired: lib.LinkCache.Expired},
		{Key: lib.TagCache.Key, Len: lib.TagCache.Length(), Note: lib.TagCache.Note, Expired: lib.TagCache.Expired},
		{Key: lib.TagFindCache.Key, Len: lib.TagFindCache.Length(), Note: lib.TagFindCache.Note, Expired: lib.TagFindCache.Expired},
		{Key: lib.TagArticleCache.Key, Len: lib.TagArticleCache.Length(), Note: lib.TagArticleCache.Note, Expired: lib.TagArticleCache.Expired},
		{Key: lib.TopicCache.Key, Len: lib.TopicCache.Length(), Note: lib.TopicCache.Note, Expired: lib.TopicCache.Expired},
		{Key: lib.TopicFindCache.Key, Len: lib.TopicFindCache.Length(), Note: lib.TopicFindCache.Note, Expired: lib.TopicFindCache.Expired},
		{Key: lib.TopicArticleCache.Key, Len: lib.TopicArticleCache.Length(), Note: lib.TopicArticleCache.Note, Expired: lib.TopicArticleCache.Expired},
		{Key: lib.SiteCache.Key, Len: lib.SiteCache.Length(), Note: lib.SiteCache.Note, Expired: lib.SiteCache.Expired},
		{Key: lib.SiteFindCache.Key, Len: lib.SiteFindCache.Length(), Note: lib.SiteFindCache.Note, Expired: lib.SiteFindCache.Expired},
		{Key: lib.NavGetByFlagCache.Key, Len: lib.NavGetByFlagCache.Length(), Note: lib.NavGetByFlagCache.Note, Expired: lib.NavGetByFlagCache.Expired},
		{Key: lib.NavGetCache.Key, Len: lib.NavGetCache.Length(), Note: lib.NavGetCache.Note, Expired: lib.NavGetCache.Expired},
		{Key: lib.NavCategoryGetCache.Key, Len: lib.NavCategoryGetCache.Length(), Note: lib.NavCategoryGetCache.Note, Expired: lib.NavCategoryGetCache.Expired},
		{Key: lib.ChannelCache.Key, Len: lib.ChannelCache.Length(), Note: lib.ChannelCache.Note, Expired: lib.ChannelCache.Expired},
		{Key: lib.ChannelGetCache.Key, Len: lib.ChannelGetCache.Length(), Note: lib.ChannelGetCache.Note, Expired: lib.ChannelGetCache.Expired},
		{Key: lib.WechatAccountCache.Key, Len: lib.WechatAccountCache.Length(), Note: lib.WechatAccountCache.Note, Expired: lib.WechatAccountCache.Expired},
		{Key: lib.WechatVerifyCache.Key, Len: lib.WechatVerifyCache.Length(), Note: lib.WechatVerifyCache.Note, Expired: lib.WechatVerifyCache.Expired},*/
	}
	ctrl.JSONPage(lib.CodeSuccess, "", items, int64(len(items)))
}

// Reset 重置缓存
func (ctrl *CacheController) Reset() {
	cacheKey := ctrl.GetSafeString("cacheKey")
	switch cacheKey {
	case "CategoryNavCache":
		lib.CategoryNavCache.Reset()
		lib.CategoryGetCache.Reset()
		lib.CategoryFindCache.Reset()
	case "ArticleCache":
		lib.ArticleCache.Reset()
		lib.ArticleGetNewCache.Reset()
	case "NoticeCache":
		lib.NoticeCache.Reset()
	case "AdsCache":
		lib.AdsCache.Reset()
	case "LinkCache":
		lib.LinkCache.Reset()
	case "TagCache":
		lib.TagCache.Reset()
		lib.TagFindCache.Reset()
		lib.TagArticleCache.Reset()
	case "TopicCache":
		lib.TopicCache.Reset()
		lib.TopicFindCache.Reset()
		lib.TopicArticleCache.Reset()
	case "SiteCache":
		lib.SiteCache.Reset()
		lib.SiteFindCache.Reset()
	case "NavGetCache":
		lib.NavGetCache.Reset()
		lib.NavGetByFlagCache.Reset()
		lib.NavCategoryGetCache.Reset()
	case "ChannelCache":
		lib.ChannelCache.Reset()
		lib.ChannelGetCache.Reset()
	case "WechatAccountCache":
		lib.WechatAccountCache.Reset()
		lib.WechatVerifyCache.Reset()
	}
	ctrl.JSONSuccess("重置成功", nil)
}
