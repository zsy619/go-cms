package funcs

import (
	"haedu.gov.cn/cms/app/cms/service"
	service_model "haedu.gov.cn/cms/app/cms/service/model"
)

/**
 * @description: ArticleNew 获取最新文章列表
 * @param {int} limit 获取数量，小于等于0时按6条处理
 * @param {int64} channel_id 频道ID
 * @param {string} channel_name 频道名称
 * @param {int64} category_id 栏目ID
 * @param {string} call_index 栏目别名
 * @param {int} is_top 是否置顶
 * @param {int} is_red 是否推荐
 * @param {int} is_hot 是否热门
 * @param {int} is_slide 是否幻灯片
 * @param {string} order_by 排序字段，为空则默认按sort_id排序，可选值：sort_id,publish_time
 * @return {*}
 */
func ArticleNewExtend(limit int, channel_id int64, channel_name string, category_id int64, call_index string, is_top, is_red, is_hot, is_slide int, order_by string) []*service_model.ApiArticleListModel {
	if limit <= 0 {
		limit = 6
	}
	find, _, err := service.NewApiArticle().ArticleGetNew(limit, channel_id, channel_name, category_id, call_index, is_top, is_red, is_hot, is_slide, order_by)
	if err != nil {
		find = []*service_model.ApiArticleListModel{}
	}
	return find
}

/**
 * @description: 获取最新文章列表
 * @param {int} limit 获取数量，小于等于0时按6条处理
 * @return {*}
 */
func ArticleNew(limit int) []*service_model.ApiArticleListModel {
	return ArticleNewExtend(limit, 0, "", 0, "", 0, 0, 0, 0, "")
}

/**
 * @description: ArticleTop 获取文章列表
 * @param {int} limit 获取数量，小于等于0时按6条处理
 * @param {int64} channel_id 频道ID
 * @param {string} channel_name 频道名称
 * @param {int64} category_id 栏目ID
 * @param {string} call_index 栏目别名
 * @param {int} is_top 是否置顶
 * @param {int} is_red 是否推荐
 * @param {int} is_hot 是否热门
 * @param {int} is_slide 是否幻灯片
 * @param {string} order_by 排序字段，为空则默认按sort_id排序，可选值：sort_id,publish_time
 * @return {*}
 */
func ArticleTopExtend(limit int, channel_id int64, channel_name string, category_id int64, call_index string, is_top, is_red, is_hot, is_slide int, order_by string) []*service_model.ApiArticleListModel {
	if limit <= 0 {
		limit = 6
	}
	find, _, err := service.NewApiArticle().ArticleGet(limit, channel_id, channel_name, category_id, call_index, is_top, is_red, is_hot, is_slide, order_by)
	if err != nil {
		find = []*service_model.ApiArticleListModel{}
	}
	return find
}

/**
 * @description: 获取文章列表
 * @param {int} limit 获取数量，小于等于0时按6条处理
 * @return {*}
 */
func ArticleTop(limit int) []*service_model.ApiArticleListModel {
	return ArticleTopExtend(limit, 0, "", 0, "", 0, 0, 0, 0, "")
}

/**
 * @description: 获取栏目导航
 * @param {string} channel_name 频道名称
 * @param {in64} channel_id 频道ID
 * @param {string} call_index 栏目别名
 * @param {in64} category_id 栏目ID
 * @param {int64} article_id 文章ID
 * @return {*}
 */
func CategoryNav(channel_name string, channel_id int64, call_index string, category_id int64, article_id int64) []*service_model.ApiCategoryNav {
	find, err := service.NewApiArticle().CategoryNav(channel_name, channel_id, call_index, category_id, article_id)
	if err != nil {
		find = []*service_model.ApiCategoryNav{}
	}
	return find
}

/**
 * @description: CategoryGet 获取栏目列表
 * @param {string} channel_name 频道名称
 * @return {*}
 */
func CategoryGet(channel_name string) []*service_model.ApiCategoryGetModel {
	if channel_name == "" {
		return []*service_model.ApiCategoryGetModel{}
	}
	find, _, err := service.NewApiArticle().CategoryGet(channel_name, "")
	if err != nil {
		find = []*service_model.ApiCategoryGetModel{}
	}
	return find
}

/**
 * @description: CategoryFind 获取栏目详情
 * @param {int64} category_id 栏目ID
 * @param {string} call_index 栏目别名
 * @return {*}
 */
func CategoryFind(category_id int64, call_index string) *service_model.ApiCategoryFindModel {
	find, err := service.NewApiArticle().CategoryFind(category_id, call_index)
	if err != nil {
		find = &service_model.ApiCategoryFindModel{}
	}
	return find
}

/**
 * @description: ArticleFind 根据article_id获取文章详情、相册、附件
 * @param {string} call_index 调用别名
 * @param {int64} article_id 文章id
 * @return {*}
 */
func ArticleFind(call_index string, article_id int64) *service_model.ApiArticleModel {
	find1, find2, find3, find4, err := service.NewApiArticle().ArticleFind(call_index, article_id)
	if err != nil {
		return &service_model.ApiArticleModel{
			Article:  &service_model.ApiArticleOneModel{},
			Album:    []*service_model.ApiAlbumModel{},
			Attach:   []*service_model.ApiAttachModel{},
			Property: []*service_model.ApiPropertyModel{},
		}
	}
	result := &service_model.ApiArticleModel{
		Article:  find1,
		Album:    find2,
		Attach:   find3,
		Property: find4,
	}
	return result
}

/**
 * @description: Get 根据article_id获取文章上一个、下一个
 * @param {string} call_index 栏目调用别名
 * @param {int64} category_id 栏目id
 * @param {int64} article_id 文章id
 * @return {*}
 */
func ArticlePrevNext(call_index string, category_id, article_id int64) *service_model.ApiArticlePrevNext {
	result := &service_model.ApiArticlePrevNext{
		Prev: &service_model.ApiArticlePrevNextModel{},
		Next: &service_model.ApiArticlePrevNextModel{},
	}
	prev, next := service.NewApiArticle().PrevNext(call_index, category_id, article_id)
	if prev != nil {
		result.Prev = prev
	}
	if next != nil {
		result.Next = next
	}

	return result
}

/**
 * @description: Article 获取文章详情
 * @param {string} call_index 调用别名
 * @param {int64} article_id 文章id
 * @return {*}
 */
func ArticleArticle(call_index string, article_id int64) *service_model.ApiArticleOneModel {
	find, err := service.NewApiArticle().Article(call_index, article_id)
	if err != nil {
		find = &service_model.ApiArticleOneModel{}
	}
	return find
}

/**
 * @description: 获取文章相册列表
 * @param {string} call_index 文章调用别名
 * @param {int64} article_id 文章id
 * @param {int32} type_id 分类
 * @return {*}
 */
func ArticleAlbum(call_index string, article_id int64, type_id int32) []*service_model.ApiAlbumModel {
	find, err := service.NewApiArticle().Album(call_index, article_id, type_id)
	if err != nil {
		find = []*service_model.ApiAlbumModel{}
	}
	return find
}

/**
 * @description: 获取文章附件列表
 * @param {string} call_index 文章调用别名
 * @param {int64} article_id 文章id
 * @param {int32} type_id 分类
 * @return {*}
 */
func ArticleAttach(call_index string, article_id int64, type_id int32) []*service_model.ApiAttachModel {
	find, err := service.NewApiArticle().Attach(call_index, article_id, type_id)
	if err != nil {
		find = []*service_model.ApiAttachModel{}
	}
	return find
}
