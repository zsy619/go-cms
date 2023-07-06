package biz

import (
	"fmt"
	"haedu.gov.cn/cms/app/lib"

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
	/*if found, item := ApiCache.Get(cacheKey); found {
		return item.(*bizmodel.ApiSiteModel), nil
	}*/
	if found, item := lib.SiteCache.Get(cacheKey); found {
		mdl := item.(*bizmodel.ApiSiteModel)
		logs.Debug("SiteList[Cache]::", "cacheKey", cacheKey, "SiteList", mdl)
		return mdl, nil
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
		// ApiCache.Set(cacheKey, find, 1800)
		lib.SiteCache.Set(cacheKey, find)
	} else {
		find = &bizmodel.ApiSiteModel{}
	}
	return find, nil
}

/**
 * @description: Find 获取站点信息
 * @param {int64} site_id 站点ID
 * @return {*}
 */
func (this *ApiSite) Find(site_id int64) (*model.CmsSite, error) {
	cacheKey := "ApiSite_Find_" + fmt.Sprintf("%d", site_id)
	if found, item := lib.SiteFindCache.Get(cacheKey); found {
		logs.Debug("ApiCache")
		return item.(*model.CmsSite), nil
	}
	out, err := NewCmsSite().SiteOne(site_id)
	if err != nil {
		return nil, err
	} else {
		if out != nil {
			lib.SiteFindCache.Set(cacheKey, out)
		}
	}
	return out, nil
}

/**
 * @description: ChannelGet 获取站点频道
 * @param {int64} site_id 站点ID
 * @return {*}
 */
func (this *ApiSite) ChannelGet(site_id int64) ([]*bizmodel.ApiChannelModel, int64, error) {
	cacheKey := "ApiSite_ChannelGet_" + fmt.Sprintf("%d", site_id)
	if found, item := lib.ChannelGetCache.Get(cacheKey); found {
		logs.Debug("ApiCache")
		find := item.([]*bizmodel.ApiChannelModel)
		return find, int64(len(find)), nil
	}
	outChannel := []*bizmodel.ApiChannelModel{}
	mdl, do := query.CmsSiteChannelDo()
	err := do.Where(mdl.SiteID.Eq(site_id)).Select(mdl.ChannelID, mdl.ParentID, mdl.Title, mdl.Name, mdl.Kind, mdl.ClassLayer, mdl.ImgUrl1, mdl.ImgUrl2, mdl.SortID, mdl.IsAlbum, mdl.IsAttach, mdl.IsSpec).Order(mdl.SortID).Scan(&outChannel)
	if err != nil {
		return nil, 0, err
	} else {
		if len(outChannel) > 0 {
			lib.ChannelGetCache.Set(cacheKey, outChannel, 1800)
		}
	}
	return outChannel, int64(len(outChannel)), nil
}

func (this *ApiSite) NavGetByFlag(site_flag string, channel_id int64) ([]*bizmodel.ApiNavModel, int64, error) {
	cacheKey := fmt.Sprintf("ApiSite_NavGetByFlag_%s_%d", site_flag, channel_id)
	if found, item := lib.NavGetByFlagCache.Get(cacheKey); found {
		find := item.([]*bizmodel.ApiNavModel)
		return find, int64(len(find)), nil
	}
	site_id := int64(0)
	siteMdl, siteDo := query.CmsSiteDo()
	siteErr := siteDo.Where(siteMdl.Flag.Eq(site_flag)).Pluck(siteMdl.SiteID, &site_id)
	if siteErr != nil {
		return nil, 0, siteErr
	}
	find, count, err := this.NavGet(site_id, channel_id)
	if err != nil {
		return nil, 0, err
	}
	if len(find) > 0 {
		lib.NavGetByFlagCache.Set(cacheKey, find)
	}
	return find, count, nil
}

/**
 * @description: NavGet 获取站点导航
 * @param {int64} site_id 站点ID
 * @param {int64} channel_id 频道ID
 * @return {*}
 */
func (this *ApiSite) NavGet(site_id int64, channel_id int64) ([]*bizmodel.ApiNavModel, int64, error) {
	cacheKey := fmt.Sprintf("ApiSite_NavGet_%d_%d", site_id, channel_id)
	if found, item := lib.NavGetCache.Get(cacheKey); found {
		find := item.([]*bizmodel.ApiNavModel)
		return find, int64(len(find)), nil
	}
	flag := ""
	siteMdl, siteDo := query.CmsSiteDo()
	siteErr := siteDo.Where(siteMdl.SiteID.Eq(site_id)).Pluck(siteMdl.Flag, &flag)
	if siteErr != nil {
		return nil, 0, siteErr
	}
	outNav := []*bizmodel.ApiNavModel{}
	mdl, do := query.CmsSiteChannelDo()
	err := do.Where(mdl.SiteID.Eq(site_id), mdl.ParentID.Eq(channel_id), mdl.Status.Eq(int32(StatusPass)), mdl.IsShow.Is(true)).Select(
		mdl.ChannelID.As("nav_id"), mdl.Title, mdl.Name, mdl.LinkURL, mdl.Target,
		mdl.ImgUrl1, mdl.ImgUrl2, mdl.SortID).Order(mdl.SortID).Scan(&outNav)
	if err != nil {
		return []*bizmodel.ApiNavModel{}, 0, err
	}
	for _, v := range outNav {
		v.Type = "channel"
		v.Children = []*bizmodel.ApiNavModel{}
		children, _, _ := this.NavGet(site_id, v.NavID)
		if len(children) > 0 {
			v.Children = append(v.Children, children...)
		}
		// 获取频道下的栏目
		categorys := this.NavCategoryGet(v.NavID, 0, flag, v.Name)
		if len(categorys) > 0 {
			v.Children = append(v.Children, categorys...)
		}
	}
	if len(outNav) > 0 {
		lib.NavGetCache.Set(cacheKey, outNav)
	}
	return outNav, int64(len(outNav)), nil
}

func (this *ApiSite) NavCategoryGet(channel_id int64, parent_id int64, flag, name string) []*bizmodel.ApiNavModel {
	cacheKey := fmt.Sprintf("ApiSite_NavCategoryGet_%d_%d", channel_id, parent_id)
	if found, item := lib.NavCategoryGetCache.Get(cacheKey); found {
		find := item.([]*bizmodel.ApiNavModel)
		return find
	}
	outNav := []*bizmodel.ApiNavModel{}
	mdl, do := query.CmsArticleCategoryDo()
	err := do.Where(mdl.ChannelID.Eq(channel_id), mdl.ParentID.Eq(parent_id), mdl.Status.Eq(int32(StatusPass)), mdl.IsShow.Is(true)).Select(
		mdl.CategoryID.As("nav_id"), mdl.Title, mdl.CallIndex.As("name"), mdl.LinkURL, mdl.Target,
		mdl.ImgUrl1, mdl.ImgUrl2, mdl.SortID).Order(mdl.SortID).Scan(&outNav)
	if err != nil {
		return []*bizmodel.ApiNavModel{}
	}
	for _, v := range outNav {
		v.Type = "category"
		// 如果链接为空，则自动拼接，否则使用自定义链接
		// 自动path：/站点标识/频道名称/栏目调用名
		if v.LinkURL == "" {
			v.LinkURL = "/" + flag + "/" + name + "/" + v.Name
		}
		v.Children = []*bizmodel.ApiNavModel{}
		children := this.NavCategoryGet(channel_id, v.NavID, flag, name)
		v.Children = append(v.Children, children...)
	}
	if len(outNav) > 0 {
		lib.NavCategoryGetCache.Set(cacheKey, outNav)
	}
	return outNav
}
