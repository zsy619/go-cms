package funcs

import (
	"fmt"

	"haedu.gov.cn/cms/app/biz"
)

/**
 * @description: 文章链接
 * @param {*} flag 站点标识
 * @param {*} name 频道名称
 * @param {string} call_index 栏目别名
 * @param {int64} article_id 文章ID
 * @return {*}
 */
func UrlForArticleExt(flag, name, call_index string, article_id int64) string {
	return fmt.Sprintf("/%s/%s/%s/%d", flag, name, call_index, article_id)
}

/**
 * @description: 文章链接
 * @param {string} call_index 栏目别名
 * @param {int64} article_id 文章ID
 * @param {string} url 自定义链接
 * @return {*}
 */
func UrlForArticle(call_index string, article_id int64, url string) string {
	if url != "" {
		return url
	}
	return fmt.Sprintf("/article/%s/%d", call_index, article_id)
}

/**
 * @description: 栏目链接
 * @param {*} flag 站点标识
 * @param {*} name 频道名称
 * @param {string} call_index 栏目别名
 * @return {*}
 */
func UrlForCategoryExt(flag, name, call_index string) string {
	find, err := biz.NewApiArticle().CategoryFind(0, call_index)
	if err != nil {
		return ""
	}
	if find.LinkURL != "" {
		return find.LinkURL
	}
	return fmt.Sprintf("/%s/%s/%s", flag, name, call_index)
}

/**
 * @description: 栏目链接
 * @param {string} call_index 栏目别名
 * @return {*}
 */
func UrlForCategory(call_index string) string {
	find, err := biz.NewApiArticle().CategoryFind(0, call_index)
	if err != nil {
		return ""
	}
	if find.LinkURL != "" {
		return find.LinkURL
	}
	return fmt.Sprintf("/%s/%s/%s", find.SiteFlag, find.ChannelName, find.CallIndex)
}

/**
 * @description: 频道链接
 * @param {string} name 频道名称
 * @return {*}
 */
func UrlForChannel(name string) string {
	find, err := biz.NewApiChannel().Find(name, 0)
	if err != nil {
		return ""
	}
	if find.LinkURL != "" {
		return find.LinkURL
	}
	return fmt.Sprintf("/%s/%s", find.SiteFlag, name)
}
