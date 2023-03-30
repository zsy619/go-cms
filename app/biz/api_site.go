package biz

import (
	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/app/dal/query"
)

// var (
// 	Cache_ApiSiteDefault     *model.CmsSite
// 	Cache_ApiSiteGet         map[int64]*model.CmsSite
// 	Cache_ApiSiteChannelFind map[int64][]map[string]interface{}
// )

// func init() {
// 	Cache_ApiSiteDefault = nil
// 	Cache_ApiSiteGet = make(map[int64]*model.CmsSite)
// 	Cache_ApiSiteChannelFind = make(map[int64][]map[string]interface{})
// }

type ApiSite struct{}

func NewApiSite() *ApiSite {
	return &ApiSite{}
}

/**
 * @description: Default 获取站点信息
 * @return {*}
 */
func (this *ApiSite) Default() (*model.CmsSite, error) {
	cacheKey := "ApiSiteDefault"
	if found, item := ApiCache.Get(cacheKey); found {
		logs.Debug("ApiCache")
		return item.(*model.CmsSite), nil
	}
	// if Cache_ApiSiteDefault != nil {
	// 	logs.Debug("Cache_ApiSiteDefault")
	// 	return Cache_ApiSiteDefault, nil
	// }
	site, siteDo := query.CmsSiteDo()
	find, err := siteDo.Where(site.IsDefault.Is(true), site.IsDeleted.Is(false)).First()
	if err == nil {
		// Cache_ApiSiteDefault = find
		ApiCache.Set(cacheKey, find, 1800)
	} else {
		find = &model.CmsSite{}
	}
	return find, nil
}

/**
 * @description: Get 获取站点信息
 * @param {int64} site_id 站点ID
 * @return {*}
 */
func (this *ApiSite) Get(site_id int64) (*model.CmsSite, error) {
	cacheKey := "ApiSiteGet_" + string(site_id)
	if found, item := ApiCache.Get(cacheKey); found {
		logs.Debug("ApiCache")
		return item.(*model.CmsSite), nil
	}
	// if v, ok := Cache_ApiSiteGet[site_id]; ok {
	// 	logs.Debug("Cache_ApiSiteGet")
	// 	return v, nil
	// }
	out, err := NewCmsSite().SiteOne(site_id)
	if err != nil {
		return nil, err
	} else {
		// Cache_ApiSiteGet[site_id] = out
		ApiCache.Set(cacheKey, out, 1800)
	}
	return out, nil
}

/**
 * @description: ChannelFind 获取站点频道
 * @param {int64} site_id 站点ID
 * @return {*}
 */
func (this *ApiSite) ChannelFind(site_id int64) ([]map[string]interface{}, int64, error) {
	cacheKey := "ApiSiteChannelFind_" + string(site_id)
	if found, item := ApiCache.Get(cacheKey); found {
		logs.Debug("ApiCache")
		find := item.([]map[string]interface{})
		return find, int64(len(find)), nil
	}
	// if find, ok := Cache_ApiSiteChannelFind[site_id]; ok {
	// 	logs.Debug("Cache_ApiSiteChannelFind")
	// 	return find, int64(len(find)), nil
	// }
	outChannel := []map[string]interface{}{}
	mdl, do := query.CmsSiteChannelDo()
	err := do.Where(mdl.SiteID.Eq(site_id)).Select(mdl.ChannelID, mdl.ParentID, mdl.Title, mdl.Name, mdl.Kind, mdl.ClassLayer, mdl.SortID, mdl.IsAlbum, mdl.IsAttach, mdl.IsSpec).Order(mdl.SortID).Scan(&outChannel)
	if err != nil {
		return nil, 0, err
	} else {
		// Cache_ApiSiteChannelFind[site_id] = outChannel
		ApiCache.Set(cacheKey, outChannel, 1800)
	}
	return outChannel, int64(len(outChannel)), nil
}
