package funcs

/**
 * @description: 标签云
 * @param {int64} siteId 站点ID
 * @param {int} limit 条数
 * @return {*}
 */
func TagNew(siteId int64, limit int) []interface{} {
	if limit <= 0 {
		limit = 5
	}
	return []interface{}{}
}
