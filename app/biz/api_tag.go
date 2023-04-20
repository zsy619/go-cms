package biz

import (
	"fmt"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"gorm.io/gorm"
	"haedu.gov.cn/cms/app/biz/bizmodel"
	"haedu.gov.cn/cms/app/dal/query"
)

type ApiTag struct{}

func NewApiTag() *ApiTag {
	return &ApiTag{}
}

func (this *ApiTag) find(cackeKeyPrefix string, limit int, siteId, channelId int64) ([]*bizmodel.ApiTagListModel, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	cacheKey := fmt.Sprintf("%s_%d_%d_%d", cackeKeyPrefix, limit, siteId, channelId)
	if found, item := ApiCache.Get(cacheKey); found {
		tags := item.([]*bizmodel.ApiTagListModel)
		logs.Debug("TagFind[Cache]::", "cacheKey", cacheKey, "Ads", tags)
		return tags, int64(len(tags)), nil
	}

	outTags := make([]*bizmodel.ApiTagListModel, 0)
	mdl, do := query.CmsTagDo()
	if siteId > 0 {
		do = do.Where(mdl.SiteID.Eq(siteId))
	}
	if channelId > 0 {
		do = do.Where(mdl.ChannelID.Eq(channelId))
	}
	do = do.Where(mdl.Status.Eq(int32(StatusPass)))
	do = do.Select(mdl.TagID, mdl.SiteID, mdl.ChannelID, mdl.Name, mdl.Title, mdl.ImgUrl1, mdl.ImgUrl2,
		mdl.SeoTitle, mdl.SeoKeyword, mdl.SeoDescription, mdl.SortID)
	if cackeKeyPrefix == "ApiTag_Find" {
		do = do.Order(mdl.SortID, mdl.TagID)
	} else {
		do = do.Order(mdl.TagID.Desc(), mdl.SortID)
	}
	err := do.Scan(&outTags)
	if err == nil {
		ApiCache.Set(cacheKey, outTags, 1800)
	}
	return outTags, int64(len(outTags)), err
}

/**
 * @description: 获取标签列表
 * @param {int} limit 限制数量
 * @param {*} siteId 站点ID
 * @param {int64} channelId 栏目ID
 * @return {*}
 */
func (this *ApiTag) Find(limit int, siteId, channelId int64) ([]*bizmodel.ApiTagListModel, int64, error) {
	return this.find("ApiTag_Find", limit, siteId, channelId)
}

/**
 * @description: 获取最新标签列表
 * @param {int} limit 限制数量
 * @param {*} siteId 站点ID
 * @param {int64} channelId 栏目ID
 * @return {*}
 */
func (this *ApiTag) FindNew(limit int, siteId, channelId int64) ([]*bizmodel.ApiTagListModel, int64, error) {
	return this.find("ApiTag_FindNew", limit, siteId, channelId)
}

/**
 * @description: Click 点击数+1
 * @param {int64} tag_id 标签ID
 * @return {*}
 */
func (this *ApiTag) Click(tag_id int64) error {
	mdl, do := query.CmsTagDo()
	do.Where(mdl.TagID.Eq(tag_id), mdl.Status.Eq(int32(StatusPass))).Updates(map[string]interface{}{
		mdl.Click.ColumnName().String():      gorm.Expr("click + ?", 1),
		mdl.UpdateTime.ColumnName().String(): time.Now(),
	})
	return nil
}
