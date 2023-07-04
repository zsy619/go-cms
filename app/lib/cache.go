package lib

import "haedu.gov.cn/tools/xcache"

type CacheItemModel struct {
	Key     string `json:"key"`     // Key 缓存键
	Len     int    `json:"len"`     // Len 缓存长度
	Note    string `json:"note"`    // Note 缓存说明
	Expired int64  `json:"expired"` // Expired 过期时间

	cache *xcache.ExpiredMap `json:"-"` // cache 缓存对象
}

// NewCacheItemModel 创建缓存项
func NewCacheItemModel(key string, note string, expired int64) *CacheItemModel {
	return &CacheItemModel{
		Key:     key,
		Note:    note,
		Expired: expired,
		cache:   xcache.NewExpiredMap(),
	}
}

// Set 设置缓存
func (this *CacheItemModel) Set(cacheKey interface{}, data interface{}, expired ...int64) bool {
	var _expired int64
	if len(expired) > 0 {
		_expired = expired[0]
	} else {
		_expired = this.Expired
	}
	return this.cache.Set(cacheKey, data, _expired)
}

// Get 获取缓存
func (this *CacheItemModel) Get(cacheKey interface{}) (bool, interface{}) {
	return this.cache.Get(cacheKey)
}

// Remove 删除缓存
func (this *CacheItemModel) Remove(cacheKey interface{}) {
	this.cache.Remove(cacheKey)
}

// Length 获取缓存长度
func (this *CacheItemModel) Length() int {
	return this.cache.Length()
}

// Reset 重置缓存
func (this *CacheItemModel) Reset() {
	this.cache.Clear()
	this.cache.Close()
	this.cache = xcache.NewExpiredMap()
}

var (
	ArticleCache       *CacheItemModel // ArticleCache 文章列表缓存
	NoticeCache        *CacheItemModel // NoticeCache 系统公告列表缓存
	AdsCache           *CacheItemModel // AdsCache 广告列表缓存
	LinkCache          *CacheItemModel // LinkCache 链接列表缓存
	TagCache           *CacheItemModel // TagCache 标签列表缓存
	TopicCache         *CacheItemModel // TopicCache 专题列表缓存
	SiteCache          *CacheItemModel // SiteCache 站点列表缓存
	ChannelCache       *CacheItemModel // ChannelCache 频道列表缓存
	WechatAccountCache *CacheItemModel // WechatAccountCache 微信账号列表缓存
	WechatVerifyCache  *CacheItemModel // WechatVerifyCache 微信校验文件列表缓存
)

// init 初始化
func init() {
	ArticleCache = NewCacheItemModel("ArticleCache", "文章列表缓存", 60*20)
	NoticeCache = NewCacheItemModel("NoticeCache", "系统公告列表缓存", 60*20)
	AdsCache = NewCacheItemModel("AdsCache", "广告列表缓存", 60*20)
	LinkCache = NewCacheItemModel("LinkCache", "链接列表缓存", 60*20)
	TagCache = NewCacheItemModel("TagCache", "标签列表缓存", 60*20)
	TopicCache = NewCacheItemModel("TopicCache", "专题列表缓存", 60*20)
	SiteCache = NewCacheItemModel("SiteCache", "站点列表缓存", 60*20)
	ChannelCache = NewCacheItemModel("ChannelCache", "频道列表缓存", 60*20)
	WechatAccountCache = NewCacheItemModel("WechatAccountCache", "微信账号列表缓存", 60*20)
	WechatVerifyCache = NewCacheItemModel("WechatVerifyCache", "微信校验文件列表缓存", 60*20)
}
