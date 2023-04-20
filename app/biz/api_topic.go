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

type ApiTopic struct{}

func NewApiTopic() *ApiTopic {
	return &ApiTopic{}
}

func (this *ApiTopic) find(cackeKeyPrefix string, limit int, siteId, channelId int64) ([]*bizmodel.ApiTopicListModel, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	cacheKey := fmt.Sprintf("%s_%d_%d_%d", cackeKeyPrefix, limit, siteId, channelId)
	if found, item := ApiCache.Get(cacheKey); found {
		list := item.([]*bizmodel.ApiTopicListModel)
		logs.Debug("TopicFind[Cache]::", "cacheKey", cacheKey, "Topic", list)
		return list, int64(len(list)), nil
	}

	list := make([]*bizmodel.ApiTopicListModel, 0)
	mdl, do := query.CmsTopicDo()
	if siteId > 0 {
		do = do.Where(mdl.SiteID.Eq(siteId))
	}
	if channelId > 0 {
		do = do.Where(mdl.ChannelID.Eq(channelId))
	}
	do = do.Where(mdl.Status.Eq(int32(StatusPass)))
	do = do.Select(mdl.TopicID, mdl.SiteID, mdl.ChannelID, mdl.Name, mdl.Title, mdl.ImgUrl1, mdl.ImgUrl2,
		mdl.SeoTitle, mdl.SeoKeyword, mdl.SeoDescription, mdl.SortID, mdl.Click, mdl.Template)
	if cackeKeyPrefix == "ApiTopic_Find" {
		do = do.Order(mdl.SortID, mdl.TopicID)
	} else {
		do = do.Order(mdl.TopicID.Desc(), mdl.SortID)
	}
	err := do.Scan(&list)
	if err == nil {
		ApiCache.Set(cacheKey, list, 1800)
	}
	return list, int64(len(list)), err
}

/**
 * @description: 获取专题列表
 * @param {int} limit 限制数量
 * @param {*} siteId 站点ID
 * @param {int64} channelId 栏目ID
 * @return {*}
 */
func (this *ApiTopic) Find(limit int, siteId, channelId int64) ([]*bizmodel.ApiTopicListModel, int64, error) {
	return this.find("ApiTopic_Find", limit, siteId, channelId)
}

/**
 * @description: 获取最新专题列表
 * @param {int} limit 限制数量
 * @param {*} siteId 站点ID
 * @param {int64} channelId 栏目ID
 * @return {*}
 */
func (this *ApiTopic) FindNew(limit int, siteId, channelId int64) ([]*bizmodel.ApiTopicListModel, int64, error) {
	return this.find("ApiTopic_FindNew", limit, siteId, channelId)
}

/**
 * @description: Click 点击数+1
 * @param {int64} tag_id 专题ID
 * @return {*}
 */
func (this *ApiTopic) Click(tag_id int64) error {
	mdl, do := query.CmsTagDo()
	do.Where(mdl.TagID.Eq(tag_id), mdl.Status.Eq(int32(StatusPass))).Updates(map[string]interface{}{
		mdl.Click.ColumnName().String():      gorm.Expr("click + ?", 1),
		mdl.UpdateTime.ColumnName().String(): time.Now(),
	})
	return nil
}

/**
 * @description: ArticlePaginate 获取文章分页列表
 * @param {*} page 页码
 * @param {int} limit 每页数量
 * @param {string} topic_name 专题名称
 * @param {int64} channel_id 频道ID
 * @param {string} channel_name 频道名称
 * @param {int64} category_id 栏目ID
 * @param {string} call_index 栏目别名
 * @param {string} keyword 关键词：按标题、摘要进行搜索
 * @param {int} is_top 是否置顶
 * @param {int} is_red 是否推荐
 * @param {int} is_hot 是否热门
 * @param {int} is_slide 是否幻灯片
 * @param {int} is_search 是否搜索
 * @param {string} order_by 排序字段，为空则默认按sort_id排序，可选值：sort_id,publish_time
 * @return {*}
 */
func (this *ApiTopic) ArticlePaginate(page, limit int, topic_name string, channel_id int64, channel_name string, category_id int64, call_index string, keyword string, is_top, is_red, is_hot, is_slide, is_search int, order_by string) ([]*bizmodel.ApiArticleListModel, int64, error) {
	order_by = xgeneric.IFF(order_by == "", "sort_id", order_by)
	where := xgeneric.IFF(channel_id <= 0, "", " And b.channel_id="+strconv.FormatInt(channel_id, 10)) +
		xgeneric.IFF(channel_name == "", "", " And c.name='"+channel_name+"'") +
		xgeneric.IFF(topic_name == "", "", " And locate(',"+topic_name+",',concat(',',a.topic,','))") +
		xgeneric.IFF(category_id <= 0, "", " And b.category_id="+strconv.FormatInt(category_id, 10)) +
		xgeneric.IFF(call_index == "", "", " And b.call_index='"+call_index+"'") +
		xgeneric.IFF(is_top < 0, "", " And a.is_top="+strconv.Itoa(is_hot)) +
		xgeneric.IFF(is_red < 0, "", " And a.is_red="+strconv.Itoa(is_hot)) +
		xgeneric.IFF(is_hot < 0, "", " And a.is_hot="+strconv.Itoa(is_hot)) +
		xgeneric.IFF(is_slide < 0, "", " And a.is_slide="+strconv.Itoa(is_slide)) +
		xgeneric.IFF(is_search < 0, "", " And b.is_search="+strconv.Itoa(is_search)) +
		xgeneric.IFF(keyword == "", "", " And (a.title LIKE '%"+keyword+"%' OR a.summary LIKE '%"+keyword+"%')")

	sqlSelectCount := "COUNT(1) as count"
	sqlCount := "SELECT " + sqlSelectCount + " FROM cms_article a" +
		" LEFT JOIN cms_article_category b ON a.category_id=b.category_id" +
		" LEFT JOIN cms_site_channel c ON a.channel_id = c.channel_id" +
		" WHERE a.`status`=2 AND b.`status`=2" +
		where

	var count int64

	_, do := query.CmsArticleDo()
	if err := do.UnderlyingDB().Raw(sqlCount).Pluck("count", &count).Error; err != nil {
		return nil, 0, err
	}

	sql := bizmodel.ApiArticleListModel_Table +
		where +
		" ORDER BY a.is_top DESC,a." + order_by +
		" LIMIT ? OFFSET ?"
	list := make([]*bizmodel.ApiArticleListModel, 0)
	err := do.UnderlyingDB().Raw(sql, limit, (page-1)*limit).Scan(&list).Error

	return list, count, err
}

func (this *ApiTopic) ArtilceTop(limit int, topic_name string) ([]*bizmodel.ApiArticleListModel, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	cacheKey := fmt.Sprintf("%s_%d_%s", "TopicArtilceTop", limit, topic_name)
	if found, item := ApiCache.Get(cacheKey); found {
		list := item.([]*bizmodel.ApiArticleListModel)
		logs.Debug("TopicFind[Cache]::", "cacheKey", cacheKey, "Topic", list)
		return list, int64(len(list)), nil
	}
	list, _, err := NewApiTopic().ArticlePaginate(1, limit, topic_name, 0, "", 0, "", "", 0, 0, 0, 0, 0, "")
	if err != nil {
		return []*bizmodel.ApiArticleListModel{}, 0, err
	}
	ApiCache.Set(cacheKey, list, 1800)
	return list, int64(len(list)), nil
}
