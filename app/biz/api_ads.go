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

func (this *ApiAds) get(cackeKeyPrefix string, limit int, site_id int64, site_flag string, category_id int64, call_index string) ([]*bizmodel.ApiAdsModel, int64, error) {
	cacheKey := fmt.Sprintf("%s_%d_%d_%s_%d_%s", cackeKeyPrefix, limit, site_id, site_flag, category_id, call_index)
	if found, item := ApiCache.Get(cacheKey); found {
		list := item.([]*bizmodel.ApiAdsModel)
		logs.Debug("AdsFind[Cache]::", "cacheKey", cacheKey, "Ads", list)
		return list, int64(len(list)), nil
	}
	list := []*bizmodel.ApiAdsModel{}

	_, do := query.CmsAdsDo()
	sqlSelect := "a.ads_id,a.site_id,a.channel_id,a.category_id,b.title as category_title,a.title,a.link_url,a.target,a.click,a.img_url1,a.img_url2,a.is_lock,a.is_red,a.is_hot,a.is_slide,a.begin_time,a.end_time,c.flag as site_flag"
	sql := `SELECT ` + sqlSelect + ` FROM cms_ads a` +
		` LEFT JOIN cms_ads_category b ON a.category_id = b.category_id` +
		` LEFT JOIN cms_site c ON a.site_id = c.site_id` +
		` WHERE a.status=2 and NOW() between a.begin_time and a.end_time ` +
		xgeneric.IFF(site_flag == "", "", " AND c.flag = '"+site_flag+"'") +
		xgeneric.IFF(site_id <= 0, "", " AND b.site_id = "+strconv.FormatInt(site_id, 10)) +
		xgeneric.IFF(call_index == "", "", " AND b.call_index = '"+call_index+"'") +
		xgeneric.IFF(category_id <= 0, "", " AND b.category_id = "+strconv.FormatInt(category_id, 10))
	if cackeKeyPrefix == "ApiAds_GetNew" {
		sql += " ORDER BY a.begin_time DESC,a.sort_id ASC"
	} else {
		sql += " ORDER BY a.is_top DESC,a.sort_id ASC"
	}
	if limit > 0 {
		sql += fmt.Sprintf(" LIMIT %d", limit)
	}
	err := do.UnderlyingDB().Raw(sql).Scan(&list).Error
	if err == nil {
		ApiCache.Set(cacheKey, list, 2400)
	}
	return list, int64(len(list)), err
}

/**
* @description: Get 获取广告列表
* @param {int} limit 获取数量
* @param {int64} site_id 站点ID
* @param {string} site_flag 站点标识
* @param {int64} category_id 广告分类ID
* @param {string} call_index 广告分类标识
* @return {*}
 */
func (this *ApiAds) Get(limit int, site_id int64, site_flag string, category_id int64, call_index string) ([]*bizmodel.ApiAdsModel, int64, error) {
	return this.get("ApiAds_Get", limit, site_id, site_flag, category_id, call_index)
}

/**
* @description: GetNew 获取最新广告列表
* @param {int} limit 获取数量
* @param {int64} site_id 站点ID
* @param {string} site_flag 站点标识
* @param {int64} category_id 广告分类ID
* @param {string} call_index 广告分类标识
* @return {*}
 */
func (this *ApiAds) GetNew(limit int, site_id int64, site_flag string, category_id int64, call_index string) ([]*bizmodel.ApiAdsModel, int64, error) {
	return this.get("ApiAds_GetNew", limit, site_id, site_flag, category_id, call_index)
}

/**
 * @description: Paginate 获取广告列表
 * @param {int} page 页码
 * @param {int} limit 获取数量
 * @param {int64} site_id 站点ID
 * @param {string} site_flag 站点标识
 * @param {int64} category_id 广告分类ID
 * @param {string} call_index 广告分类标识
 * @return {*}
 */
func (this *ApiAds) Paginate(page, limit int, site_id int64, site_flag string, category_id int64, call_index string) ([]*bizmodel.ApiAdsModel, int64, error) {
	list := []*bizmodel.ApiAdsModel{}
	where := " WHERE a.`status`=2 and NOW() between a.begin_time and a.end_time" +
		xgeneric.IFF(site_flag == "", "", " AND c.flag = '"+site_flag+"'") +
		xgeneric.IFF(site_id <= 0, "", " AND a.site_id = "+strconv.FormatInt(site_id, 10)) +
		xgeneric.IFF(call_index == "", "", " AND b.call_index = '"+call_index+"'") +
		xgeneric.IFF(category_id <= 0, "", " AND b.category_id = "+strconv.FormatInt(category_id, 10))
	_, do := query.CmsAdsDo()
	sqlSelectRow := "a.ads_id,a.site_id,a.channel_id,a.category_id,b.title as category_title,a.title,a.link_url,a.target,a.click,a.img_url1,a.img_url2,a.is_lock,a.is_red,a.is_hot,a.is_slide,a.begin_time,a.end_time,c.flag as site_flag"
	sqlRow := "SELECT " + sqlSelectRow + " FROM cms_ads a" +
		" LEFT JOIN cms_ads_category b ON a.category_id = b.category_id" +
		" LEFT JOIN cms_site c ON a.site_id = c.site_id" +
		where +
		" ORDER BY a.is_top ASC,a.sort_id ASC"

	sqlSelectCount := "count(1) as count"
	sqlCount := "SELECT " + sqlSelectCount + " FROM cms_ads a" +
		" LEFT JOIN cms_ads_category b ON a.category_id = b.category_id" +
		" LEFT JOIN cms_site c ON a.site_id = c.site_id" +
		where

	var count int64
	do.UnderlyingDB().Raw(sqlCount).Pluck("count", &count)

	sqlRow += fmt.Sprintf(" LIMIT %d,%d", (page-1)*limit, limit)
	err := do.UnderlyingDB().Raw(sqlRow).Scan(&list).Error
	return list, count, err
}

/**
 * @description: Click 点击数+1
 * @param {int64} ads_id 广告ID
 * @return {*}
 */
func (this *ApiAds) Click(ads_id int64) error {
	mdl, do := query.CmsAdsDo()
	_, err := do.Where(mdl.AdsID.Eq(ads_id), mdl.Status.Eq(int32(StatusPass))).Updates(map[string]interface{}{
		mdl.Click.ColumnName().String():      gorm.Expr("click + ?", 1),
		mdl.UpdateTime.ColumnName().String(): time.Now(),
	})
	return err
}
