package funcs

/**
 * @description: 最新文章
 * @param {int64} siteId 站点ID
 * @param {string} channelName 栏目名称
 * @param {string} categoryName 分类调用链接
 * @param {int} limit 条数
 * @return {*}
 */
func ArticleNew(siteId int64, channelName string, categoryName string, limit int) []interface{} {
	if limit <= 0 {
		limit = 5
	}
	return []interface{}{}
}

/**
 * @description: 置顶文章
 * @param {int64} siteId 站点ID
 * @param {string} channelName 栏目名称
 * @param {string} categoryName 分类调用链接
 * @param {int} limit 条数
 * @return {*}
 */
func ArticleTop(siteId int64, channelName string, categoryName string, limit int) []interface{} {
	if limit <= 0 {
		limit = 5
	}
	return []interface{}{}
}
