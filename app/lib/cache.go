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
	ArticleCache        *CacheItemModel // ArticleCache 文章列表缓存
	ArticleGetNewCache  *CacheItemModel // ArticleGetNewCache 最新文章列表缓存
	CategoryNavCache    *CacheItemModel // CategoryNavCache 栏目导航缓存
	CategoryGetCache    *CacheItemModel // CategoryGetCache 栏目列表缓存
	CategoryFindCache   *CacheItemModel // CategoryFindCache 栏目详情缓存
	NoticeCache         *CacheItemModel // NoticeCache 系统公告列表缓存
	AdsCache            *CacheItemModel // AdsCache 广告列表缓存
	LinkCache           *CacheItemModel // LinkCache 链接列表缓存
	TagCache            *CacheItemModel // TagCache 标签列表缓存
	TagArticleCache     *CacheItemModel // TagArticleCache 标签文章列表缓存
	TagFindCache        *CacheItemModel // TagFindCache 标签详情缓存
	TopicCache          *CacheItemModel // TopicCache 专题列表缓存
	TopicFindCache      *CacheItemModel // TopicFindCache 专题详情缓存
	TopicArticleCache   *CacheItemModel // TopicArticleCache 专题文章列表缓存
	SiteCache           *CacheItemModel // SiteCache 站点列表缓存
	SiteFindCache       *CacheItemModel // SiteFindCache 站点信息缓存
	NavGetByFlagCache   *CacheItemModel // NavGetByFlagCache 站点导航缓存
	NavGetCache         *CacheItemModel // NavGetCache 导航菜单缓存
	NavCategoryGetCache *CacheItemModel // NavCategoryGetCache 站点栏目缓存
	ChannelCache        *CacheItemModel // ChannelCache 频道列表缓存
	ChannelGetCache     *CacheItemModel // ChannelGetCache 站点频道缓存
	WechatAccountCache  *CacheItemModel // WechatAccountCache 微信账号列表缓存
	WechatVerifyCache   *CacheItemModel // WechatVerifyCache 微信校验文件列表缓存
)

// init 初始化
func init() {
	CategoryNavCache = NewCacheItemModel("CategoryNavCache", "栏目导航缓存", 60*100)
	CategoryGetCache = NewCacheItemModel("CategoryGetCache", "栏目列表缓存", 60*100)
	CategoryFindCache = NewCacheItemModel("CategoryFindCache", "栏目详情缓存", 60*100)
	ArticleCache = NewCacheItemModel("ArticleCache", "文章列表缓存", 60*20)
	ArticleGetNewCache = NewCacheItemModel("ArticleGetNewCache", "最新文章列表缓存", 60*20)
	NoticeCache = NewCacheItemModel("NoticeCache", "系统公告列表缓存", 60*90)
	AdsCache = NewCacheItemModel("AdsCache", "广告列表缓存", 60*80)
	LinkCache = NewCacheItemModel("LinkCache", "链接列表缓存", 60*70)
	TagCache = NewCacheItemModel("TagCache", "标签列表缓存", 60*60)
	TagFindCache = NewCacheItemModel("TagFindCache", "标签详情缓存", 60*60)
	TagArticleCache = NewCacheItemModel("TagArticleCache", "标签文章列表缓存", 60*50)
	TopicCache = NewCacheItemModel("TopicCache", "专题列表缓存", 60*30)
	TopicFindCache = NewCacheItemModel("TopicFindCache", "专题详情缓存", 60*30)
	TopicArticleCache = NewCacheItemModel("TopicArticleCache", "专题文章列表缓存", 60*30)
	SiteCache = NewCacheItemModel("SiteCache", "站点列表缓存", 60*50)
	SiteFindCache = NewCacheItemModel("SiteFindCache", "站点信息缓存", 60*50)
	NavGetByFlagCache = NewCacheItemModel("NavGetByFlagCache", "站点导航缓存", 60*50)
	NavGetCache = NewCacheItemModel("NavGetCache", "导航菜单缓存", 60*50)
	NavCategoryGetCache = NewCacheItemModel("NavCategoryGetCache", "站点栏目缓存", 60*50)
	ChannelCache = NewCacheItemModel("ChannelCache", "频道列表缓存", 60*50)
	ChannelGetCache = NewCacheItemModel("ChannelGetCache", "站点频道缓存", 60*50)
	WechatAccountCache = NewCacheItemModel("WechatAccountCache", "微信账号列表缓存", 60*90)
	WechatVerifyCache = NewCacheItemModel("WechatVerifyCache", "微信校验文件列表缓存", 60*90)
}
