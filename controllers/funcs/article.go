package funcs

import (
	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/biz/bizmodel"
)

/**
 * @description: ArticleNew 获取最新文章列表
 * @param {int} limit 获取数量
 * @param {int64} channel_id 频道ID
 * @param {string} channel_name 频道名称
 * @param {int64} category_id 栏目ID
 * @param {string} call_index 栏目别名
 * @param {int} is_top 是否置顶
 * @param {int} is_red 是否推荐
 * @param {int} is_hot 是否热门
 * @param {int} is_slide 是否幻灯片
 * @param {string} order_by 排序字段，为空则默认按sort_id排序，可选值：sort_id,publish_time
 * @param {bool} is_cache 是否使用缓存
 * @return {*}
 */
func ArticleNewExt(limit int, channel_id int64, channel_name string, category_id int64, call_index string, is_top, is_red, is_hot, is_slide int, order_by string) []*bizmodel.ApiArticleListModel {
	if limit <= 0 {
		limit = 5
	}
	find, _, err := biz.NewApiArticle().ArticleGetNew(limit, channel_id, channel_name, category_id, call_index, is_top, is_red, is_hot, is_slide, order_by)
	if err != nil {
		find = []*bizmodel.ApiArticleListModel{}
	}
	return find
}

func ArticleNew(limit int) []*bizmodel.ApiArticleListModel {
	return ArticleNewExt(limit, 0, "", 0, "", 0, 0, 0, 0, "")
}

/**
 * @description: ArticleTop 获取文章列表
 * @param {int} limit 获取数量
 * @param {int64} channel_id 频道ID
 * @param {string} channel_name 频道名称
 * @param {int64} category_id 栏目ID
 * @param {string} call_index 栏目别名
 * @param {int} is_top 是否置顶
 * @param {int} is_red 是否推荐
 * @param {int} is_hot 是否热门
 * @param {int} is_slide 是否幻灯片
 * @param {string} order_by 排序字段，为空则默认按sort_id排序，可选值：sort_id,publish_time
 * @param {bool} is_cache 是否使用缓存
 * @return {*}
 */
func ArticleTopExt(limit int, channel_id int64, channel_name string, category_id int64, call_index string, is_top, is_red, is_hot, is_slide int, order_by string) []*bizmodel.ApiArticleListModel {
	if limit <= 0 {
		limit = 5
	}
	find, _, err := biz.NewApiArticle().ArticleGet(limit, channel_id, channel_name, category_id, call_index, is_top, is_red, is_hot, is_slide, order_by)
	if err != nil {
		find = []*bizmodel.ApiArticleListModel{}
	}
	return find
}

func ArticleTop(limit int) []*bizmodel.ApiArticleListModel {
	return ArticleTopExt(limit, 0, "", 0, "", 0, 0, 0, 0, "")
}
