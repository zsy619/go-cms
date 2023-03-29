package biz

import (
	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/app/dal/query"
)

var (
	Cache_ApiSiteGet         map[int64]*model.CmsSite
	Cache_ApiSiteChannelFind map[int64][]map[string]interface{}
)

func init() {
	Cache_ApiSiteGet = make(map[int64]*model.CmsSite)
	Cache_ApiSiteChannelFind = make(map[int64][]map[string]interface{})
}

type ApiSite struct{}

func NewApiSite() *ApiSite {
	return &ApiSite{}
}

// Get 获取站点信息
func (this *ApiSite) Get(site_id int64) (*model.CmsSite, error) {
	if v, ok := Cache_ApiSiteGet[site_id]; ok {
		logs.Debug("Cache_ApiSiteGet")
		return v, nil
	}
	out, err := NewCmsSite().SiteOne(site_id)
	if err != nil {
		return nil, err
	} else {
		Cache_ApiSiteGet[site_id] = out
	}
	return out, nil
}

// ChannelFind 获取站点栏目
func (this *ApiSite) ChannelFind(site_id int64) ([]map[string]interface{}, int64, error) {
	if find, ok := Cache_ApiSiteChannelFind[site_id]; ok {
		logs.Debug("Cache_ApiSiteChannelFind")
		return find, int64(len(find)), nil
	}
	outChannel := []map[string]interface{}{}
	mdl, do := query.CmsSiteChannelDo()
	err := do.Where(mdl.SiteID.Eq(site_id)).Select(mdl.ChannelID, mdl.ParentID, mdl.Title, mdl.Name, mdl.Kind, mdl.ClassLayer, mdl.SortID, mdl.IsAlbum, mdl.IsAttach, mdl.IsSpec).Order(mdl.SortID).Scan(&outChannel)
	if err != nil {
		return nil, 0, err
	} else {
		Cache_ApiSiteChannelFind[site_id] = outChannel
	}
	return outChannel, int64(len(outChannel)), nil
}
