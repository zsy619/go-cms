package biz

import (
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
	"haedu.gov.cn/cms/app/biz/bmodel"
	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/app/dal/query"
	"haedu.gov.cn/tools/xgeneric"
)

// var (
// 	Cache_ApiArticleCategoryFind map[string][]map[string]interface{}
// 	Cache_ApiArticleFind         map[string][]map[string]interface{}
// )

// func init() {
// 	Cache_ApiArticleCategoryFind = make(map[string][]map[string]interface{})
// 	Cache_ApiArticleFind = make(map[string][]map[string]interface{})
// }

type ApiArticle struct{}

func NewApiArticle() *ApiArticle {
	return &ApiArticle{}
}

/**
 * @description: CategoryFind 获取栏目列表
 * @param {string} channel_name 频道名称
 * @return {*}
 */
func (this *ApiArticle) CategoryFind(channel_name string) ([]map[string]interface{}, int64, error) {
	cacheKey := fmt.Sprintf("ApiArticleCategoryFind_%s", channel_name)
	if found, item := ApiCache.Get(cacheKey); found {
		find := item.([]map[string]interface{})
		return find, int64(len(find)), nil
	}
	// if outChannel, ok := Cache_ApiArticleCategoryFind[channel_name]; ok {
	// 	return outChannel, int64(len(outChannel)), nil
	// }
	outChannel := []map[string]interface{}{}
	_, do := query.CmsArticleCategoryDo()
	sqlSelect := "b.`name` as channel_name,b.title as channel_title,a.category_id,a.parent_id,a.site_id,a.channel_id,a.title,a.call_index,a.link_url,a.img_url,a.sort_id"
	sql := "SELECT " + sqlSelect + " FROM cms_article_category a LEFT JOIN cms_site_channel b ON a.channel_id=b.channel_id WHERE a.is_deleted=0 AND a.`status`=2 AND a.`is_show`=1 AND b.`name`=? ORDER BY a.sort_id"
	err := do.UnderlyingDB().Raw(sql, channel_name).Scan(&outChannel).Error
	if err != nil {
		return nil, 0, err
	}
	ApiCache.Set(cacheKey, outChannel, 1800)
	return outChannel, int64(len(outChannel)), nil
}

/**
 * @description: CategoryOne 获取栏目详情
 * @param {int64} category_id 栏目ID
 * @param {string} call_index 栏目别名
 * @return {*}
 */
func (this *ApiArticle) CategoryOne(category_id int64, call_index string) (*bmodel.CategoryOneModel, error) {
	cacheKey := fmt.Sprintf("ApiArticleCategoryOne_%d_%s", category_id, call_index)
	if found, item := ApiCache.Get(cacheKey); found {
		return item.(*bmodel.CategoryOneModel), nil
	}
	outChannel := &bmodel.CategoryOneModel{}
	_, do := query.CmsArticleCategoryDo()
	sqlSelect := "b.`name` as channel_name,b.title as channel_title,a.category_id,a.channel_id,a.title,a.call_index,a.link_url,a.img_url,a.seo_title,a.seo_keyword,a.seo_description,a.content"
	sql := "SELECT " + sqlSelect + " FROM cms_article_category a LEFT JOIN cms_site_channel b ON a.channel_id=b.channel_id WHERE a.is_deleted=0 AND a.`status`=2 AND a.`is_show`=1" +
		xgeneric.IFF(category_id > 0, " AND a.category_id="+strconv.FormatInt(category_id, 10), "") +
		xgeneric.IFF(call_index != "", " AND a.call_index='"+call_index+"'", "")
	err := do.UnderlyingDB().Raw(sql).Scan(&outChannel).Error
	if err != nil {
		return nil, err
	}
	ApiCache.Set(cacheKey, outChannel, 1800)
	return outChannel, nil
}

/**
 * @description: Find 获取文章列表
 * @param {int} limit 获取数量
 * @param {int64} channel_id 频道ID
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
func (this *ApiArticle) Find(limit int, channel_id, category_id int64, call_index string, is_top, is_red, is_hot, is_slide int, order_by string, is_cache bool) ([]map[string]interface{}, int64, error) {
	order_by = xgeneric.IFF(order_by == "", "sort_id", order_by)
	cacheKey := fmt.Sprintf("ApiArticleFind_%s_%d_%d_%d_%d_%d_%d_%s", call_index, channel_id, is_top, is_red, is_hot, is_slide, limit, order_by)
	if is_cache {
		if found, item := ApiCache.Get(cacheKey); found {
			find := item.([]map[string]interface{})
			return find, int64(len(find)), nil
		}
		// if outArticle, ok := Cache_ApiArticleFind[cacheKey]; ok {
		// 	return outArticle, int64(len(outArticle)), nil
		// }
	}
	outArticle := []map[string]interface{}{}
	_, do := query.CmsArticleDo()
	sqlSelect := "a.article_id,a.site_id,a.channel_id,a.category_id,a.title,a.sub_title,a.ico_url,a.call_index,a.source,a.author,a.link_url,a.img_url,a.seo_title,a.seo_keyword,a.seo_description,a.tags,a.summary,a.click,a.is_lock,a.is_comment,a.like_count,a.is_top,a.is_hot,a.is_slide,a.static_url,a.publish_time"
	sql := "SELECT " + sqlSelect + " FROM cms_article a LEFT JOIN cms_article_category b ON a.category_id=b.category_id WHERE a.`status`=2 AND b.`status`=2" +
		xgeneric.IFF(channel_id <= 0, "", " And b.channel_id="+strconv.FormatInt(channel_id, 10)) +
		xgeneric.IFF(category_id <= 0, "", " And b.category_id="+strconv.FormatInt(category_id, 10)) +
		xgeneric.IFF(call_index == "", "", " And b.call_index='"+call_index+"'") +
		xgeneric.IFF(is_top <= 0, "", " And a.is_top="+strconv.Itoa(is_hot)) +
		xgeneric.IFF(is_red <= 0, "", " And a.is_red="+strconv.Itoa(is_hot)) +
		xgeneric.IFF(is_hot <= 0, "", " And a.is_hot="+strconv.Itoa(is_hot)) +
		xgeneric.IFF(is_slide <= 0, "", " And a.is_slide="+strconv.Itoa(is_slide)) +
		" ORDER BY a.is_top DESC,a." + order_by +
		" LIMIT ?"
	err := do.UnderlyingDB().Raw(sql, limit).Scan(&outArticle).Error
	if err == nil {
		ApiCache.Set(cacheKey, outArticle, 1800)
		// Cache_ApiArticleFind[cacheKey] = outArticle
	}
	return outArticle, int64(len(outArticle)), err
}

/**
 * @description: Paginate 获取文章分页列表
 * @param {*} page 页码
 * @param {int} limit 每页数量
 * @param {int64} channel_id 频道ID
 * @param {int64} category_id 栏目ID
 * @param {string} call_index 栏目别名
 * @param {string} keyword 关键词：按标题、摘要进行搜索
 * @param {int} is_top 是否置顶
 * @param {int} is_red 是否推荐
 * @param {int} is_hot 是否热门
 * @param {int} is_slide 是否幻灯片
 * @param {string} order_by 排序字段，为空则默认按sort_id排序，可选值：sort_id,publish_time
 * @return {*}
 */
func (this *ApiArticle) Paginate(page, limit int, channel_id, category_id int64, call_index string, keyword string, is_top, is_red, is_hot, is_slide int, order_by string) ([]map[string]interface{}, int64, error) {
	order_by = xgeneric.IFF(order_by == "", "sort_id", order_by)
	_, do := query.CmsArticleDo()
	sqlSelectCount := "COUNT(1) as count"
	sqlCount := "SELECT " + sqlSelectCount + " FROM cms_article a LEFT JOIN cms_article_category b ON a.category_id=b.category_id WHERE a.`status`=2 AND b.`status`=2" +
		xgeneric.IFF(channel_id <= 0, "", " And b.channel_id="+strconv.FormatInt(channel_id, 10)) +
		xgeneric.IFF(category_id <= 0, "", " And b.category_id="+strconv.FormatInt(category_id, 10)) +
		xgeneric.IFF(call_index == "", "", " And b.call_index='"+call_index+"'") +
		xgeneric.IFF(is_top <= 0, "", " And a.is_top="+strconv.Itoa(is_hot)) +
		xgeneric.IFF(is_red <= 0, "", " And a.is_red="+strconv.Itoa(is_hot)) +
		xgeneric.IFF(is_hot <= 0, "", " And a.is_hot="+strconv.Itoa(is_hot)) +
		xgeneric.IFF(is_slide <= 0, "", " And a.is_slide="+strconv.Itoa(is_slide)) +
		xgeneric.IFF(keyword == "", "", " And (a.title LIKE '%"+keyword+"%' OR a.summary LIKE '%"+keyword+"%')")
	var count int64
	if err := do.UnderlyingDB().Raw(sqlCount).Pluck("count", &count).Error; err != nil {
		return nil, 0, err
	}

	sqlSelect := "a.article_id,a.site_id,a.channel_id,a.category_id,a.title,a.sub_title,a.call_index,a.ico_url,a.source,a.author,a.link_url,a.seo_title,a.seo_keyword,a.seo_description,a.tags,a.summary,a.click,a.is_lock,a.is_comment,a.like_count,a.is_top,a.is_hot,a.is_slide,a.static_url,a.publish_time"
	sql := "SELECT " + sqlSelect + " FROM cms_article a LEFT JOIN cms_article_category b ON a.category_id=b.category_id WHERE a.`status`=2 AND b.`status`=2" +
		xgeneric.IFF(channel_id <= 0, "", " And b.channel_id="+strconv.FormatInt(channel_id, 10)) +
		xgeneric.IFF(call_index == "", "", " And b.call_index='"+call_index+"'") +
		xgeneric.IFF(is_top <= 0, "", " And a.is_top="+strconv.Itoa(is_hot)) +
		xgeneric.IFF(is_red <= 0, "", " And a.is_red="+strconv.Itoa(is_hot)) +
		xgeneric.IFF(is_hot <= 0, "", " And a.is_hot="+strconv.Itoa(is_hot)) +
		xgeneric.IFF(is_slide <= 0, "", " And a.is_slide="+strconv.Itoa(is_slide)) +
		xgeneric.IFF(keyword == "", "", " And (a.title LIKE '%"+keyword+"%' OR a.summary LIKE '%"+keyword+"%')") +
		" ORDER BY a.is_top DESC,a." + order_by +
		" LIMIT ? OFFSET ?"
	outArticle := []map[string]interface{}{}
	err := do.UnderlyingDB().Raw(sql, limit, (page-1)*limit).Scan(&outArticle).Error
	return outArticle, count, err
}

/**
 * @description: Get 根据article_id获取文章详情、相册、附件
 * @param {int64} article_id 文章id
 * @return {*}
 */
func (this *ApiArticle) Get(article_id int64) (*model.CmsArticle, []*model.CmsArticleAlbum, []*model.CmsArticleAttach, error) {
	mdl, do := query.CmsArticleDo()
	article, err := do.Where(mdl.ArticleID.Eq(article_id), mdl.Status.Eq(2)).First()
	if err != nil {
		return nil, nil, nil, err
	}

	albumMdl, alblumDo := query.CmsArticleAlbumDo()
	articleAlbum, _ := alblumDo.Where(albumMdl.ArticleID.Eq(article_id)).Order(albumMdl.SortID).Find()

	attachMdl, attachDo := query.CmsArticleAttachDo()
	articleAttach, _ := attachDo.Where(attachMdl.ArticleID.Eq(article_id)).Order(attachMdl.SortID).Find()

	return article, articleAlbum, articleAttach, nil
}

/**
 * @description: Article 获取文章详情
 * @param {int64} article_id 文章id
 * @return {*}
 */
func (this *ApiArticle) Article(article_id int64) (*model.CmsArticle, error) {
	mdl, do := query.CmsArticleDo()
	return do.Where(mdl.ArticleID.Eq(article_id), mdl.Status.Eq(2)).First()
}

/**
 * @description: Album 获取文章相册列表
 * @param {int64} article_id 文章id
 * @return {*}
 */
func (this *ApiArticle) Album(article_id int64) ([]*model.CmsArticleAlbum, error) {
	albumMdl, alblumDo := query.CmsArticleAlbumDo()
	return alblumDo.Where(albumMdl.ArticleID.Eq(article_id)).Order(albumMdl.SortID).Find()
}

/**
 * @description: Attach 获取文章附件列表
 * @param {int64} article_id 文章id
 * @return {*}
 */
func (this *ApiArticle) Attach(article_id int64) ([]*model.CmsArticleAttach, error) {
	attachMdl, attachDo := query.CmsArticleAttachDo()
	return attachDo.Where(attachMdl.ArticleID.Eq(article_id)).Order(attachMdl.SortID).Find()
}

/**
 * @description: Click 点击数+1
 * @param {int64} article_id 文章id
 * @return {*}
 */
func (this *ApiArticle) Click(article_id int64) error {
	mdl, do := query.CmsArticleDo()
	do.Where(mdl.ArticleID.Eq(article_id), mdl.Status.Eq(2)).Updates(map[string]interface{}{
		mdl.Click.ColumnName().String():      gorm.Expr("click + ?", 1),
		mdl.UpdateTime.ColumnName().String(): time.Now(),
	})
	return nil
}

/**
 * @description: Like 点赞数+1
 * @param {int64} article_id 文章id
 * @return {*}
 */
func (this *ApiArticle) Like(article_id int64) error {
	mdl, do := query.CmsArticleDo()
	do.Where(mdl.ArticleID.Eq(article_id), mdl.Status.Eq(2)).Updates(map[string]interface{}{
		mdl.LikeCount.ColumnName().String():  gorm.Expr("like_count + ?", 1),
		mdl.UpdateTime.ColumnName().String(): time.Now(),
	})
	return nil
}
