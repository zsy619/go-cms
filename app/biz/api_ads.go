package biz

import (
	"fmt"
	"strconv"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"gorm.io/gorm"
	"haedu.gov.cn/cms/app/biz/bizmodel"
	"haedu.gov.cn/cms/app/dal/query"
	"haedu.gov.cn/tools/xgeneric"
)

type ApiAds struct{}

func NewApiAds() *ApiAds {
	return &ApiAds{}
}

/**
* @description: Find 获取广告列表
* @param {int} limit 获取数量
* @param {int64} category_id 广告分类ID
* @param {string} call_index 广告分类标识
* @return {*}
 */
func (this *ApiAds) Find(limit int, category_id int64, call_index string) ([]*bizmodel.ApiAdsListModel, int64, error) {
	cacheKey := fmt.Sprintf("AdFind::%d::%d::%s", limit, category_id, call_index)
	if found, item := ApiCache.Get(cacheKey); found {
		Ads := item.([]*bizmodel.ApiAdsListModel)
		logs.Debug("AdFindByCategory[Cache]::", "cacheKey", cacheKey, "Ads", Ads)
		return Ads, int64(len(Ads)), nil
	}
	outAd := []*bizmodel.ApiAdsListModel{}

	_, AdDo := query.CmsAdsDo()
	sqlSelect := "a.ad_id,a.site_id,a.channel_id,a.category_id,b.title as category_title,a.title,a.link_url,a.target,a.click,a.img_url1,a.img_url2,a.is_lock,a.is_red,a.is_hot,a.is_slide,a.begin_time,a.end_time"
	sql := "SELECT " + sqlSelect + " FROM cms_ads a LEFT JOIN cms_ads_category b ON a.category_id = b.category_id WHERE a.`status`=2 and NOW() between a.begin_time and a.end_time " +
		xgeneric.IFF(call_index == "", "", " AND b.call_index = '"+call_index+"'") +
		xgeneric.IFF(category_id <= 0, "", " AND b.category_id = "+strconv.FormatInt(category_id, 10)) +
		" ORDER BY a.is_top desc,a.sort_id ASC"
	if limit > 0 {
		sql += fmt.Sprintf(" LIMIT %d", limit)
	}
	err := AdDo.UnderlyingDB().Raw(sql).Scan(&outAd).Error
	if err == nil {
		ApiCache.Set(cacheKey, outAd, 1800)
	}
	return outAd, int64(len(outAd)), err
}

/**
 * @description: Paginate 获取广告列表
 * @param {*} page 页码
 * @param {int} limit 获取数量
 * @param {int64} category_id 广告分类ID
 * @param {string} call_index 广告分类标识
 * @return {*}
 */
func (this *ApiAds) Paginate(page, limit int, category_id int64, call_index string) ([]*bizmodel.ApiAdsListModel, int64, error) {
	outAd := []*bizmodel.ApiAdsListModel{}
	_, AdDo := query.CmsAdsDo()
	sqlSelectRow := "a.ad_id,a.site_id,a.channel_id,a.category_id,b.title as category_title,a.title,a.link_url,a.target,a.click,a.img_url1,a.img_url2,a.is_lock,a.is_red,a.is_hot,a.is_slide,a.begin_time,a.end_time"
	sqlRow := "SELECT " + sqlSelectRow + " FROM cms_ads a LEFT JOIN cms_ads_category b ON a.category_id = b.category_id WHERE a.`status`=2" +
		xgeneric.IFF(call_index == "", "", " AND b.call_index = '"+call_index+"'") +
		xgeneric.IFF(category_id <= 0, "", " AND b.category_id = "+strconv.FormatInt(category_id, 10)) +
		" ORDER BY a.is_top desc,a.sort_id ASC"

	sqlSelectCount := "count(1) as count"
	sqlCount := "SELECT " + sqlSelectCount + " FROM cms_ads a LEFT JOIN cms_ads_category b ON a.category_id = b.category_id WHERE a.`status`=2 and NOW() between a.begin_time and a.end_time " +
		xgeneric.IFF(call_index == "", "", " AND b.call_index = '"+call_index+"'") +
		xgeneric.IFF(category_id <= 0, "", " AND b.category_id = "+strconv.FormatInt(category_id, 10))

	var count int64
	AdDo.UnderlyingDB().Raw(sqlCount).Pluck("count", &count)

	sqlRow += fmt.Sprintf(" LIMIT %d,%d", (page-1)*limit, limit)
	err := AdDo.UnderlyingDB().Raw(sqlRow).Scan(&outAd).Error
	return outAd, count, err
}

/**
 * @description: Click 点击数+1
 * @param {int64} ads_id 广告ID
 * @return {*}
 */
func (this *ApiAds) Click(ads_id int64) error {
	mdl, do := query.CmsAdsDo()
	do.Where(mdl.AdsID.Eq(ads_id), mdl.Status.Eq(2)).Updates(map[string]interface{}{
		mdl.Click.ColumnName().String():      gorm.Expr("click + ?", 1),
		mdl.UpdateTime.ColumnName().String(): time.Now(),
	})
	return nil
}
