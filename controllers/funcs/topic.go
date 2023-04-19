package funcs

/**
 * @description: 最新专题
 * @param {int64} siteId 站点ID
 * @param {int} limit 条数
 * @return {*}
 */
func TopicNew(siteId int64, limit int) []interface{} {
	if limit <= 0 {
		limit = 5
	}
	return []interface{}{}
}
