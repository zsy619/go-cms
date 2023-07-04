package admin

import (
	"haedu.gov.cn/cms/app/lib"
)

type CacheController struct{ BaseController }

func (this *CacheController) Index() {
	this.display()
}

// GetList 获取缓存列表
func (this *CacheController) GetList() {
	items := []*lib.CacheItemModel{
		{Key: lib.ArticleCache.Key, Len: lib.ArticleCache.Length(), Note: lib.ArticleCache.Note, Expired: lib.ArticleCache.Expired},
		{Key: lib.NoticeCache.Key, Len: lib.NoticeCache.Length(), Note: lib.NoticeCache.Note, Expired: lib.NoticeCache.Expired},
		{Key: lib.AdsCache.Key, Len: lib.AdsCache.Length(), Note: lib.AdsCache.Note, Expired: lib.AdsCache.Expired},
		{Key: lib.LinkCache.Key, Len: lib.LinkCache.Length(), Note: lib.LinkCache.Note, Expired: lib.LinkCache.Expired},
		{Key: lib.TagCache.Key, Len: lib.TagCache.Length(), Note: lib.TagCache.Note, Expired: lib.TagCache.Expired},
		{Key: lib.TopicCache.Key, Len: lib.TopicCache.Length(), Note: lib.TopicCache.Note, Expired: lib.TopicCache.Expired},
		{Key: lib.SiteCache.Key, Len: lib.SiteCache.Length(), Note: lib.SiteCache.Note, Expired: lib.SiteCache.Expired},
		{Key: lib.ChannelCache.Key, Len: lib.ChannelCache.Length(), Note: lib.ChannelCache.Note, Expired: lib.ChannelCache.Expired},
		{Key: lib.WechatAccountCache.Key, Len: lib.WechatAccountCache.Length(), Note: lib.WechatAccountCache.Note, Expired: lib.WechatAccountCache.Expired},
		{Key: lib.WechatVerifyCache.Key, Len: lib.WechatVerifyCache.Length(), Note: lib.WechatVerifyCache.Note, Expired: lib.WechatVerifyCache.Expired},
	}
	this.JSONPage(lib.CodeSuccess, "", items, int64(len(items)))
}

// Reset 重置缓存
func (this *CacheController) Reset() {
	cacheKey := this.GetString("cacheKey")
	switch cacheKey {
	case "ArticleCache":
		lib.ArticleCache.Reset()
	case "NoticeCache":
		lib.NoticeCache.Reset()
	case "AdsCache":
		lib.AdsCache.Reset()
	case "LinkCache":
		lib.LinkCache.Reset()
	case "TagCache":
		lib.TagCache.Reset()
	case "TopicCache":
		lib.TopicCache.Reset()
	case "SiteCache":
		lib.SiteCache.Reset()
	case "ChannelCache":
		lib.ChannelCache.Reset()
	case "WechatAccountCache":
		lib.WechatAccountCache.Reset()
	case "WechatVerifyCache":
		lib.WechatVerifyCache.Reset()
	}
	this.JSONSuccess("重置成功", nil)
}
