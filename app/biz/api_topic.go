package biz

import (
	"fmt"
	"haedu.gov.cn/cms/app/lib"
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

func (this *ApiTopic) get(cackeKeyPrefix string, limit int, site_id int64, site_flag string, channel_id int64) ([]*bizmodel.ApiTopicModel, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	cacheKey := fmt.Sprintf("%s_%d_%d_%s_%d", cackeKeyPrefix, limit, site_id, site_flag, channel_id)
	/*if found, item := ApiCache.Get(cacheKey); found {
		list := item.([]*bizmodel.ApiTopicModel)
		logs.Debug("TopicFind[Cache]::", "cacheKey", cacheKey, "Topic", list)
		return list, int64(len(list)), nil
	}*/
	if found, item := lib.TopicCache.Get(cacheKey); found {
		list := item.([]*bizmodel.ApiTopicModel)
		logs.Debug("TopicList[Cache]::", "cacheKey", cacheKey, "TopicList", list)
		return list, int64(len(list)), nil
	}

	list := make([]*bizmodel.ApiTopicModel, 0)
	_, do := query.CmsTopicDo()
	field := `a.topic_id,a.site_id,a.channel_id,a.name,a.title,a.img_url1,a.img_url2,a.seo_title,a.seo_keyword,a.seo_description,a.sort_id,a.click,a.template,b.flag as site_flag`
	sql := `SELECT ` + field + ` FROM cms_topic a` +
		` LEFT JOIN cms_site b ON a.site_id=b.site_id` +
		` WHERE a.status=2` +
		xgeneric.IFF(site_flag == "", "", " AND b.flag = '"+site_flag+"'") +
		xgeneric.IFF(site_id <= 0, "", " AND b.site_id = "+strconv.FormatInt(site_id, 10)) +
		xgeneric.IFF(channel_id <= 0, "", " AND a.channel_id = "+strconv.FormatInt(channel_id, 10))
	if cackeKeyPrefix == "ApiTopic_Get" {
		sql += ` ORDER BY a.sort_id,a.topic_id`
	} else {
		sql += ` ORDER BY a.topic_id desc,a.sort_id`
	}
	err := do.UnderlyingDB().Debug().Raw(sql).Scan(&list).Error
	if err == nil && len(list) > 0 {
		// ApiCache.Set(cacheKey, list, 1800)
		lib.TopicCache.Set(cacheKey, list)
	}
	return list, int64(len(list)), err
}

/**
 * @description: 获取专题列表
 * @param {int} limit 限制数量
 * @param {int64} site_id 站点ID
 * @param {string} site_flag 站点标识
 * @param {int64} channel_id 栏目ID
 * @return {*}
 */
func (this *ApiTopic) Get(limit int, site_id int64, site_flag string, channel_id int64) ([]*bizmodel.ApiTopicModel, int64, error) {
	return this.get("ApiTopic_Get", limit, site_id, site_flag, channel_id)
}

/**
 * @description: 获取最新专题列表
 * @param {int64} site_id 站点ID
 * @param {string} site_flag 站点标识
 * @param {int64} channel_id 栏目ID
 * @return {*}
 */
func (this *ApiTopic) GetNew(limit int, site_id int64, site_flag string, channel_id int64) ([]*bizmodel.ApiTopicModel, int64, error) {
	return this.get("ApiTopic_GetNew", limit, site_id, site_flag, channel_id)
}

/**
 * @description: Click 点击数+1
 * @param {int64} topic_id 专题ID
 * @return {*}
 */
func (this *ApiTopic) Click(topic_id int64) error {
	mdl, do := query.CmsTopicDo()
	_, err := do.Where(mdl.TopicID.Eq(topic_id), mdl.Status.Eq(int32(StatusPass))).Updates(map[string]interface{}{
		mdl.Click.ColumnName().String():      gorm.Expr("click + ?", 1),
		mdl.UpdateTime.ColumnName().String(): time.Now(),
	})
	return err
}

/**
 * @description: 获取专题详情
 * @param {int64} topic_id 专题ID
 * @param {string} name 专题名称
 * @return {*}
 */
func (this *ApiTopic) Find(topic_id int64, name string) (*bizmodel.ApiTopicModel, error) {
	cacheKey := fmt.Sprintf("ApiTopic_Find_%d_%s", topic_id, name)
	if found, item := ApiCache.Get(cacheKey); found {
		model := item.(*bizmodel.ApiTopicModel)
		logs.Debug("TopicFind[Cache]::", "cacheKey", cacheKey, "Topic", model)
		return model, nil
	}
	field := `a.topic_id,a.site_id,a.channel_id,a.name,a.title,a.img_url1,a.img_url2,a.seo_title,a.seo_keyword,a.seo_description,a.sort_id,a.click,a.template,b.flag as site_flag`
	sql := `SELECT ` + field + ` FROM cms_topic a` +
		` LEFT JOIN cms_site b ON a.site_id=b.site_id` +
		` WHERE a.status=2` +
		xgeneric.IFF(topic_id <= 0, "", " AND a.topic_id = "+strconv.FormatInt(topic_id, 10)) +
		xgeneric.IFF(name == "", "", " AND a.name = '"+name+"'")
	_, do := query.CmsTopicDo()
	model := &bizmodel.ApiTopicModel{}
	err := do.UnderlyingDB().Debug().Raw(sql).Scan(model).Error
	if err == nil && model.Name != "" {
		ApiCache.Set(cacheKey, model, 2400)
	}
	return model, err
}

/**
 * @description: ArticlePaginate 获取文章分页列表
 * @param {*} page 页码
 * @param {int} limit 每页数量
 * @param {int64} site_id 站点ID
 * @param {string} site_flag 站点标识
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
func (this *ApiTopic) ArticlePaginate(page, limit int, topic_name string, site_id int64, site_flag string, channel_id int64, channel_name string, category_id int64, call_index string, keyword string, is_top, is_red, is_hot, is_slide, is_search int, order_by string) ([]*bizmodel.ApiArticleListModel, int64, error) {
	order_by = xgeneric.IFF(order_by == "", "sort_id", order_by)
	where := xgeneric.IFF(site_flag == "", "", " AND d.flag = '"+site_flag+"'") +
		xgeneric.IFF(site_id <= 0, "", " AND d.site_id = "+strconv.FormatInt(site_id, 10)) +
		xgeneric.IFF(channel_id <= 0, "", " And b.channel_id="+strconv.FormatInt(channel_id, 10)) +
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
		" LEFT JOIN cms_site d ON a.site_id = d.site_id" +
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
	cacheKey := fmt.Sprintf("%s_%d_%s", "ApiTopic_ArtilceTop", limit, topic_name)
	if found, item := ApiCache.Get(cacheKey); found {
		list := item.([]*bizmodel.ApiArticleListModel)
		logs.Debug("ArtilceTop[Cache]::", "cacheKey", cacheKey, "Topic", list)
		return list, int64(len(list)), nil
	}
	list, _, err := NewApiTopic().ArticlePaginate(1, limit, topic_name, 0, "", 0, "", 0, "", "", 0, 0, 0, 0, 0, "")
	if err != nil {
		return []*bizmodel.ApiArticleListModel{}, 0, err
	}
	ApiCache.Set(cacheKey, list, 1800)
	return list, int64(len(list)), nil
}
