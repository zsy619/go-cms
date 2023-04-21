package biz

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/cms/app/biz/bizmodel"
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
func (this *ApiSite) Default() (*bizmodel.ApiSiteModel, error) {
	cacheKey := "ApiSite_Default"
	if found, item := ApiCache.Get(cacheKey); found {
		return item.(*bizmodel.ApiSiteModel), nil
	}
	site, siteDo := query.CmsSiteDo()
	find := &bizmodel.ApiSiteModel{}
	err := siteDo.Where(site.IsDefault.Is(true), site.IsDeleted.Is(false)).Scan(&find)
	if err == nil {
		// 获取默认模板
		if find.DirPath == "" {
			mdl, do := query.CmsThemeDo()
			theme, err := do.Where(mdl.IsDefault.Is(true)).First()
			if err != nil {
			} else {
				find.DirPath = theme.Name
			}
		}
		ApiCache.Set(cacheKey, find, 1800)
	} else {
		find = &bizmodel.ApiSiteModel{}
	}
	return find, nil
}

/**
 * @description: Get 获取站点信息
 * @param {int64} site_id 站点ID
 * @return {*}
 */
func (this *ApiSite) Get(site_id int64) (*model.CmsSite, error) {
	cacheKey := "ApiSite_Get_" + fmt.Sprintf("%d", site_id)
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
func (this *ApiSite) ChannelFind(site_id int64) ([]*bizmodel.ApiChannelFindModel, int64, error) {
	cacheKey := "ApiSite_ChannelFind_" + fmt.Sprintf("%d", site_id)
	if found, item := ApiCache.Get(cacheKey); found {
		logs.Debug("ApiCache")
		find := item.([]*bizmodel.ApiChannelFindModel)
		return find, int64(len(find)), nil
	}
	outChannel := []*bizmodel.ApiChannelFindModel{}
	mdl, do := query.CmsSiteChannelDo()
	err := do.Where(mdl.SiteID.Eq(site_id)).Select(mdl.ChannelID, mdl.ParentID, mdl.Title, mdl.Name, mdl.Kind, mdl.ClassLayer, mdl.ImgUrl1, mdl.ImgUrl2, mdl.SortID, mdl.IsAlbum, mdl.IsAttach, mdl.IsSpec).Order(mdl.SortID).Scan(&outChannel)
	if err != nil {
		return nil, 0, err
	} else {
		ApiCache.Set(cacheKey, outChannel, 1800)
	}
	return outChannel, int64(len(outChannel)), nil
}

func (this *ApiSite) NavFind(site_id int64, channel_id int64) ([]*bizmodel.ApiNavFindModel, int64, error) {
	return nil, 0, nil
}
