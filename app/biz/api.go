package biz

func CleanCahe() {
	Cache_ApiArticleCategoryFind = make(map[string][]map[string]interface{})
	Cache_ApiArticleFind = make(map[string][]map[string]interface{})
}
