package biz

import "haedu.gov.cn/cms/app/dal/model"

func CleanCahe() {
	Cache_ApiArticleCategoryFind = make(map[string][]map[string]interface{})
	Cache_ApiArticleFind = make(map[string][]map[string]interface{})

	Cache_ApiSiteDefault = nil
	Cache_ApiSiteGet = make(map[int64]*model.CmsSite)
	Cache_ApiSiteChannelFind = make(map[int64][]map[string]interface{})

	Cache_LinkFindByCategory = make(map[string][]map[string]interface{})
}
