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
		if find.Template == "" {
			mdl, do := query.CmsThemeDo()
			theme, err := do.Where(mdl.IsDefault.Is(true)).First()
			if err != nil {
			} else {
				find.Template = theme.Name
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

/**
 * @description: NavFind 获取站点导航
 * @param {int64} site_id 站点ID
 * @param {int64} channel_id 频道ID
 * @return {*}
 */
func (this *ApiSite) NavFind(site_id int64, channel_id int64) ([]*bizmodel.ApiNavFindModel, int64, error) {
	cacheKey := "ApiSite_NavFind_" + fmt.Sprintf("%d", site_id) + "_" + fmt.Sprintf("%d", channel_id)
	if found, item := ApiCache.Get(cacheKey); found {
		find := item.([]*bizmodel.ApiNavFindModel)
		return find, int64(len(find)), nil
	}
	outNav := []*bizmodel.ApiNavFindModel{}
	mdl, do := query.CmsSiteChannelDo()
	err := do.Where(mdl.SiteID.Eq(site_id), mdl.ParentID.Eq(channel_id)).Select(
		mdl.ChannelID.As("nav_id"), mdl.Title, mdl.Name, mdl.LinkURL, mdl.Target,
		mdl.ImgUrl1, mdl.ImgUrl2, mdl.SortID).Order(mdl.SortID).Scan(&outNav)
	if err != nil {
		return []*bizmodel.ApiNavFindModel{}, 0, err
	}
	for _, v := range outNav {
		v.Type = "channel"
		v.Children = []*bizmodel.ApiNavFindModel{}
		children, _, _ := this.NavFind(site_id, v.NavID)
		if len(children) > 0 {
			v.Children = append(v.Children, children...)
		}
		// 获取频道下的栏目
		categorys := this.NavCategoryFind(v.NavID, 0)
		if len(categorys) > 0 {
			v.Children = append(v.Children, categorys...)
		}
	}
	ApiCache.Set(cacheKey, outNav, 1800)
	return outNav, int64(len(outNav)), nil
}

func (this *ApiSite) NavCategoryFind(channel_id int64, parent_id int64) []*bizmodel.ApiNavFindModel {
	cacheKey := "ApiSite_NavCategoryFind_" + fmt.Sprintf("%d", channel_id) + "_" + fmt.Sprintf("%d", parent_id)
	if found, item := ApiCache.Get(cacheKey); found {
		find := item.([]*bizmodel.ApiNavFindModel)
		return find
	}
	outNav := []*bizmodel.ApiNavFindModel{}
	mdl, do := query.CmsArticleCategoryDo()
	err := do.Where(mdl.ChannelID.Eq(channel_id), mdl.ParentID.Eq(parent_id), mdl.IsShow.Is(true)).Select(
		mdl.CategoryID.As("nav_id"), mdl.Title, mdl.CallIndex.As("name"), mdl.LinkURL, mdl.Target,
		mdl.ImgUrl1, mdl.ImgUrl2, mdl.SortID).Order(mdl.SortID).Scan(&outNav)
	if err != nil {
		return []*bizmodel.ApiNavFindModel{}
	}
	for _, v := range outNav {
		v.Type = "category"
		v.Children = []*bizmodel.ApiNavFindModel{}
		children := this.NavCategoryFind(channel_id, v.NavID)
		v.Children = append(v.Children, children...)
	}
	ApiCache.Set(cacheKey, outNav, 1800)
	return outNav
}
