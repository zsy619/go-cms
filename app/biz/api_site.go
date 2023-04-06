package biz

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/cms/app/biz/bmodel"
	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/app/dal/query"
)

// ApiSite 站点
type ApiSite struct{}

// NewApiSite 实例化
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
	site, siteDo := query.CmsSiteDo()
	find, err := siteDo.Where(site.IsDefault.Is(true), site.IsDeleted.Is(false)).First()
	if err == nil {
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
	cacheKey := "ApiSiteGet_" + fmt.Sprintf("%d", site_id)
	if found, item := ApiCache.Get(cacheKey); found {
		logs.Debug("ApiCache")
		return item.(*model.CmsSite), nil
	}
	out, err := NewCmsSite().SiteOne(site_id)
	if err != nil {
		return nil, err
	} else {
		ApiCache.Set(cacheKey, out, 1800)
	}
	return out, nil
}

/**
 * @description: ChannelFind 获取站点频道
 * @param {int64} site_id 站点ID
 * @return {*}
 */
func (this *ApiSite) ChannelFind(site_id int64) ([]*bmodel.ApiChannelFindModel, int64, error) {
	cacheKey := "ApiSiteChannelFind_" + fmt.Sprintf("%d", site_id)
	if found, item := ApiCache.Get(cacheKey); found {
		logs.Debug("ApiCache")
		find := item.([]*bmodel.ApiChannelFindModel)
		return find, int64(len(find)), nil
	}
	outChannel := []*bmodel.ApiChannelFindModel{}
	mdl, do := query.CmsSiteChannelDo()
	err := do.Where(mdl.SiteID.Eq(site_id)).Select(mdl.ChannelID, mdl.ParentID, mdl.Title, mdl.Name, mdl.Kind, mdl.ClassLayer, mdl.ImgUrl1, mdl.ImgUrl2, mdl.SortID, mdl.IsAlbum, mdl.IsAttach, mdl.IsSpec).Order(mdl.SortID).Scan(&outChannel)
	if err != nil {
		return nil, 0, err
	} else {
		ApiCache.Set(cacheKey, outChannel, 1800)
	}
	return outChannel, int64(len(outChannel)), nil
}

func (this *ApiSite) NavFind(site_id int64, channel_id int64) ([]*bmodel.ApiNavFindModel, int64, error) {
	return nil, 0, nil
}
