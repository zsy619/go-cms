package biz

import (
	"fmt"
	"strconv"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"gorm.io/gorm"
	"haedu.gov.cn/cms/app/biz/bmodel"
	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/app/dal/query"
	"haedu.gov.cn/tools/xgeneric"
	"haedu.gov.cn/tools/xstring"
)

// ApiArticle 文章
type ApiArticle struct{}

// NewApiArticle 实例化
func NewApiArticle() *ApiArticle {
	return &ApiArticle{}
}

/**
 * @description: CategoryNav 获取栏目导航
 * @param {string} channel_name 频道名称
 * @param {in64} channel_id 频道ID
 * @param {string} call_index 栏目别名
 * @param {in64} category_id 栏目ID
 * @param {int64} article_id 文章ID
 * @return {*}
 */
func (this *ApiArticle) CategoryNav(channel_name string, channel_id int64, call_index string, category_id int64, article_id int64) ([]*bmodel.ApiCategoryNav, error) {
	if article_id > 0 && category_id <= 0 {
		cacheKey := fmt.Sprintf("ApiArticleCategoryNav_%d", article_id)
		if found, item := ApiCache.Get(cacheKey); found {
			if oks, ok := xstring.ToInt64Array(item.(string), "#"); ok && len(oks) == 2 {
				channel_id = oks[0]
				category_id = oks[1]
			}
		} else {
			mdl, do := query.CmsArticleDo()
			if article, err := do.Where(mdl.ArticleID.Eq(article_id)).Select(mdl.ChannelID, mdl.CategoryID, mdl.SiteID).First(); err != nil {
			} else {
				category_id = article.CategoryID
				channel_id = article.ChannelID
				ApiCache.Set(cacheKey, fmt.Sprintf("%d#%d", channel_id, category_id), 1800)
			}
		}
	}
	if (channel_name == "" && channel_id <= 0) && (call_index != "" || category_id > 0) {
		cacheKey := fmt.Sprintf("ApiArticleCategoryNav_%s_%d", call_index, category_id)
		if found, item := ApiCache.Get(cacheKey); found {
			channel_name = item.(string)
		} else {
			_, do := query.CmsSiteChannelDo()
			sql := `SELECT a.name FROM cms_site_channel a LEFT JOIN cms_article_category b ON a.channel_id = b.channel_id`
			if call_index != "" {
				sql += ` WHERE b.call_index='` + call_index + `'`
			} else {
				sql += ` WHERE b.category_id=` + fmt.Sprintf("%d", category_id)
			}
			if err := do.Debug().UnderlyingDB().Raw(sql).Scan(&channel_name); err != nil {
			} else {
				ApiCache.Set(cacheKey, channel_name, 1800)
			}
		}
	}
	cacheKey := fmt.Sprintf("ApiArticleCategoryNav_%s_%d_%s_%d_%d", channel_name, channel_id, call_index, category_id, article_id)
	if found, item := ApiCache.Get(cacheKey); found {
		return item.([]*bmodel.ApiCategoryNav), nil
	}
	outResult := []*bmodel.ApiCategoryNav{}
	// 获取 频道信息
	{
		mdl, do := query.CmsSiteChannelDo()
		channel := &model.CmsSiteChannel{}
		if channel_name != "" {
			channel, _ = do.Or(mdl.Name.Eq(channel_name)).Select(mdl.ChannelID, mdl.Name, mdl.Title).First()
		} else {
			channel, _ = do.Or(mdl.ChannelID.Eq(channel_id)).Select(mdl.ChannelID, mdl.Name, mdl.Title).First()
		}
		if channel != nil {
			channelItem := &bmodel.ApiCategoryNav{
				Title:     channel.Title,
				CallIndex: channel.Name,
				LinkURL:   "",
				NavType:   "channel",
			}
			outResult = append(outResult, channelItem)
		}
	}
	// 获取 栏目信息
	{
		mdl, do := query.CmsArticleCategoryDo()
		category := &model.CmsArticleCategory{}
		if call_index != "" {
			category, _ = do.Where(mdl.CallIndex.Eq(call_index)).First()
		} else {
			category, _ = do.Where(mdl.CategoryID.Eq(category_id)).First()
		}
		if category != nil {
			categoryItem := &bmodel.ApiCategoryNav{
				Title:     category.Title,
				CallIndex: category.CallIndex,
				LinkURL:   category.LinkURL,
				NavType:   "category",
			}
			outResult = append(outResult, categoryItem)
		}
	}

	ApiCache.Set(cacheKey, outResult, 1800)

	return outResult, nil
}

/**
 * @description: CategoryFind 获取栏目列表
 * @param {string} channel_name 频道名称
 * @return {*}
 */
func (this *ApiArticle) CategoryFind(channel_name string) ([]*bmodel.ApiCategoryFindModel, int64, error) {
	cacheKey := fmt.Sprintf("ApiArticleCategoryFind_%s", channel_name)
	if found, item := ApiCache.Get(cacheKey); found {
		find := item.([]*bmodel.ApiCategoryFindModel)
		return find, int64(len(find)), nil
	}
	outChannel := make([]*bmodel.ApiCategoryFindModel, 0)
	_, do := query.CmsArticleCategoryDo()
	sqlSelect := "b.`name` as channel_name,b.title as channel_title,a.category_id,a.parent_id,a.site_id,a.channel_id,a.title,a.call_index,a.class_layer,a.link_url,a.img_url1,a.img_url2,a.sort_id,a.is_show,a.is_search,a.is_deleted"
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
func (this *ApiArticle) CategoryOne(category_id int64, call_index string) (*bmodel.ApiCategoryOneModel, error) {
	cacheKey := fmt.Sprintf("ApiArticleCategoryOne_%d_%s", category_id, call_index)
	if found, item := ApiCache.Get(cacheKey); found {
		return item.(*bmodel.ApiCategoryOneModel), nil
	}
	outChannel := &bmodel.ApiCategoryOneModel{}
	_, do := query.CmsArticleCategoryDo()
	sqlSelect := "b.`name` as channel_name,b.title as channel_title,a.category_id,a.parent_id,a.site_id,a.channel_id,a.title,a.call_index,a.class_layer,a.link_url,a.img_url1,a.img_url2,a.sort_id,a.is_show,a.is_search,a.is_deleted,a.seo_title,a.seo_keyword,a.seo_description,a.content"
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
func (this *ApiArticle) Find(limit int, channel_id, category_id int64, call_index string, is_top, is_red, is_hot, is_slide int, order_by string) ([]*bmodel.ApiArticleListModel, int64, error) {
	order_by = xgeneric.IFF(order_by == "", "sort_id", order_by)
	cacheKey := fmt.Sprintf("ApiArticleFind_%s_%d_%d_%d_%d_%d_%d_%s", call_index, channel_id, is_top, is_red, is_hot, is_slide, limit, order_by)
	if found, item := ApiCache.Get(cacheKey); found {
		find := item.([]*bmodel.ApiArticleListModel)
		return find, int64(len(find)), nil
	}
	outArticle := make([]*bmodel.ApiArticleListModel, 0)
	_, do := query.CmsArticleDo()
	sqlSelect := "a.article_id,a.site_id,a.channel_id,a.category_id,a.title,a.sub_title,a.ico_url1,a.ico_url2,a.call_index,a.source,a.author,a.link_url,a.img_url1,a.img_url2,a.seo_title,a.seo_keyword,a.seo_description,a.tags,a.summary,a.click,a.is_lock,a.is_comment,a.like_count,a.is_top,a.is_hot,a.is_slide,a.static_url,a.publish_time"
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
 * @param {int} is_search 是否搜索
 * @param {string} order_by 排序字段，为空则默认按sort_id排序，可选值：sort_id,publish_time
 * @return {*}
 */
func (this *ApiArticle) Paginate(page, limit int, channel_id, category_id int64, call_index string, keyword string, is_top, is_red, is_hot, is_slide, is_search int, order_by string) ([]*bmodel.ApiArticleListModel, int64, error) {
	order_by = xgeneric.IFF(order_by == "", "sort_id", order_by)
	_, do := query.CmsArticleDo()
	sqlSelectCount := "COUNT(1) as count"
	sqlCount := "SELECT " + sqlSelectCount + " FROM cms_article a LEFT JOIN cms_article_category b ON a.category_id=b.category_id LEFT JOIN cms_site_channel c ON a.channel_id = c.channel_id WHERE a.`status`=2 AND b.`status`=2" +
		xgeneric.IFF(channel_id <= 0, "", " And b.channel_id="+strconv.FormatInt(channel_id, 10)) +
		xgeneric.IFF(category_id <= 0, "", " And b.category_id="+strconv.FormatInt(category_id, 10)) +
		xgeneric.IFF(call_index == "", "", " And b.call_index='"+call_index+"'") +
		xgeneric.IFF(is_top < 0, "", " And a.is_top="+strconv.Itoa(is_hot)) +
		xgeneric.IFF(is_red < 0, "", " And a.is_red="+strconv.Itoa(is_hot)) +
		xgeneric.IFF(is_hot < 0, "", " And a.is_hot="+strconv.Itoa(is_hot)) +
		xgeneric.IFF(is_slide < 0, "", " And a.is_slide="+strconv.Itoa(is_slide)) +
		xgeneric.IFF(is_search < 0, "", " And b.is_search="+strconv.Itoa(is_search)) +
		xgeneric.IFF(keyword == "", "", " And (a.title LIKE '%"+keyword+"%' OR a.summary LIKE '%"+keyword+"%')")
	var count int64
	if err := do.UnderlyingDB().Raw(sqlCount).Pluck("count", &count).Error; err != nil {
		return nil, 0, err
	}

	sqlSelect := `a.article_id,a.site_id,a.channel_id,a.category_id,a.title,a.sub_title,a.call_index,a.ico_url1,a.ico_url2,a.source,a.author,a.link_url,a.img_url1,a.img_url2,a.seo_title,a.seo_keyword,a.seo_description,a.tags,a.summary,a.click,a.is_lock,a.is_comment,a.like_count,a.is_top,a.is_hot,a.is_slide,a.static_url,a.publish_time`
	sqlSelect += `,b.call_index as category_call_index,b.title as category_title,b.link_url as category_link_url,c.name as channel_name,c.title as channel_title`
	sql := "SELECT " + sqlSelect + " FROM cms_article a LEFT JOIN cms_article_category b ON a.category_id=b.category_id LEFT JOIN cms_site_channel c ON a.channel_id = c.channel_id WHERE a.`status`=2 AND b.`status`=2" +
		xgeneric.IFF(channel_id <= 0, "", " And b.channel_id="+strconv.FormatInt(channel_id, 10)) +
		xgeneric.IFF(call_index == "", "", " And b.call_index='"+call_index+"'") +
		xgeneric.IFF(is_top < 0, "", " And a.is_top="+strconv.Itoa(is_hot)) +
		xgeneric.IFF(is_red < 0, "", " And a.is_red="+strconv.Itoa(is_hot)) +
		xgeneric.IFF(is_hot < 0, "", " And a.is_hot="+strconv.Itoa(is_hot)) +
		xgeneric.IFF(is_slide < 0, "", " And a.is_slide="+strconv.Itoa(is_slide)) +
		xgeneric.IFF(is_search < 0, "", " And b.is_search="+strconv.Itoa(is_search)) +
		xgeneric.IFF(keyword == "", "", " And (a.title LIKE '%"+keyword+"%' OR a.summary LIKE '%"+keyword+"%')") +
		" ORDER BY a.is_top DESC,a." + order_by +
		" LIMIT ? OFFSET ?"
	outArticle := make([]*bmodel.ApiArticleListModel, 0)
	err := do.UnderlyingDB().Raw(sql, limit, (page-1)*limit).Scan(&outArticle).Error
	return outArticle, count, err
}

/**
 * @description: Get 根据article_id获取文章详情、相册、附件
 * @param {string} call_index 文章调用别名
 * @param {int64} article_id 文章id
 * @return {*}
 */
func (this *ApiArticle) Get(call_index string, article_id int64) (*bmodel.ApiArticleOneModel, []*model.CmsAlbum, []*model.CmsAttach, error) {
	mdl, do := query.CmsArticleDo()
	if call_index != "" {
		if article_id <= 0 {
			if err := do.Where(mdl.CallIndex.Eq(call_index), mdl.Status.Eq(2)).Pluck(mdl.ArticleID, &article_id); err != nil {
				logs.Error("Get", err.Error())
			}
		}
	}
	field := `a.*,b.call_index as category_call_index,b.title as category_title,b.link_url as category_link_url,c.name as channel_name,c.title as channel_title`
	sql := `SELECT ` + field + ` FROM cms_article a LEFT JOIN cms_article_category b ON a.category_id = b.category_id LEFT JOIN cms_site_channel c ON a.channel_id = c.channel_id`
	sql += ` WHERE a.article_id=? AND a.is_deleted=0 AND a.status=2 AND b.status=2`
	artilce := &bmodel.ApiArticleOneModel{}
	if err := do.Debug().UnderlyingDB().Raw(sql, article_id).Scan(&artilce).Error; err != nil {
		logs.Error("Get", err.Error())
		return artilce, []*model.CmsAlbum{}, []*model.CmsAttach{}, err
	}

	albumMdl, alblumDo := query.CmsAlbumDo()
	articleAlbum, _ := alblumDo.Where(albumMdl.TableName_.Eq("article"), albumMdl.RecordID.Eq(article_id), albumMdl.IsShow.Eq(1)).Order(albumMdl.SortID).Find()
	if articleAlbum == nil {
		articleAlbum = []*model.CmsAlbum{}
	}
	attachMdl, attachDo := query.CmsAttachDo()
	articleAttach, _ := attachDo.Where(attachMdl.TableName_.Eq("article"), attachMdl.RecordID.Eq(article_id), attachMdl.IsShow.Eq(1)).Order(attachMdl.SortID).Find()
	if articleAttach == nil {
		articleAttach = []*model.CmsAttach{}
	}
	return artilce, articleAlbum, articleAttach, nil
}

/**
 * @description: Get 根据article_id获取文章上一个、下一个
 * @param {string} call_index 栏目调用别名
 * @param {int64} category_id 栏目id
 * @param {int64} article_id 文章id
 * @return {*}
 */
func (this *ApiArticle) PrevNext(call_index string, category_id, article_id int64) (prev *bmodel.ApiArticleOneModel, next *bmodel.ApiArticleOneModel) {
	_, do := query.CmsArticleDo()
	field := `a.*,b.call_index as category_call_index,b.title as category_title,b.link_url as category_link_url,c.name as channel_name,c.title as channel_title`
	sql := `SELECT ` + field + ` FROM cms_article a LEFT JOIN cms_article_category b ON a.category_id = b.category_id LEFT JOIN cms_site_channel c ON a.channel_id = c.channel_id`
	sql += ` WHERE a.is_deleted=0 AND a.status=2 AND b.status=2` +
		xgeneric.IFF(category_id <= 0, "", " AND b.category_id="+strconv.FormatInt(category_id, 10)) +
		xgeneric.IFF(call_index == "", "", " AND b.call_index='"+call_index+"'") +
		` AND a.article_id<? ORDER BY a.sort_id DESC LIMIT 1`
	prev = &bmodel.ApiArticleOneModel{}
	if err := do.Debug().UnderlyingDB().Raw(sql, article_id).Scan(&prev).Error; err != nil {
		logs.Error("PrevNext", err.Error())
	}
	sql = `SELECT ` + field + ` FROM cms_article a LEFT JOIN cms_article_category b ON a.category_id = b.category_id LEFT JOIN cms_site_channel c ON a.channel_id = c.channel_id`
	sql += ` WHERE a.is_deleted=0 AND a.status=2 AND b.status=2` +
		xgeneric.IFF(category_id <= 0, "", " AND b.category_id="+strconv.FormatInt(category_id, 10)) +
		xgeneric.IFF(call_index == "", "", " AND b.call_index='"+call_index+"'") +
		` AND a.article_id>? ORDER BY a.sort_id ASC LIMIT 1`
	next = &bmodel.ApiArticleOneModel{}
	if err := do.Debug().UnderlyingDB().Raw(sql, article_id).Scan(&next).Error; err != nil {
		logs.Error("PrevNext", err.Error())
	}
	return prev, next
}

/**
 * @description: Article 获取文章详情
 * @param {string} call_index 文章调用别名
 * @param {int64} article_id 文章id
 * @return {*}
 */
func (this *ApiArticle) Article(call_index string, article_id int64) (*bmodel.ApiArticleOneModel, error) {
	_, do := query.CmsArticleDo()
	field := `a.*,b.call_index as category_call_index,b.title as category_title,b.link_url as category_link_url,c.name as channel_name,c.title as channel_title`
	sql := `SELECT ` + field + ` FROM cms_article a LEFT JOIN cms_article_category b ON a.category_id = b.category_id LEFT JOIN cms_site_channel c ON a.channel_id = c.channel_id`
	artilce := &bmodel.ApiArticleOneModel{}
	if call_index != "" {
		sql += ` WHERE a.call_index=? AND a.is_deleted=0 AND a.status=2 AND b.status=2`
		if err := do.Debug().UnderlyingDB().Raw(sql, call_index).Scan(&artilce).Error; err != nil {
			logs.Error("Get", err.Error())
			return artilce, err
		}
		return artilce, nil
	}
	sql += ` WHERE a.article_id=? AND a.is_deleted=0 AND a.status=2 AND b.status=2`
	if err := do.Debug().UnderlyingDB().Raw(sql, article_id).Scan(&artilce).Error; err != nil {
		logs.Error("Get", err.Error())
		return artilce, err
	}
	return artilce, nil
}

/**
 * @description: Album 获取文章相册列表
 * @param {string} call_index 文章调用别名
 * @param {int64} article_id 文章id
 * @return {*}
 */
func (this *ApiArticle) Album(call_index string, article_id int64, type_id int32) ([]*model.CmsAlbum, error) {
	if call_index != "" {
		article, articleDo := query.CmsArticleDo()
		articleDo.Where(article.CallIndex.Eq(call_index), article.Status.Eq(2)).Pluck(article.ArticleID, &article_id)
	}
	albumMdl, alblumDo := query.CmsAlbumDo()
	return alblumDo.Where(albumMdl.TableName_.Eq("article"), albumMdl.RecordID.Eq(article_id), albumMdl.TypeID.Eq(type_id), albumMdl.IsShow.Eq(1)).Order(albumMdl.SortID).Find()
}

/**
 * @description: Attach 获取文章附件列表
 * @param {string} call_index 文章调用别名
 * @param {int64} article_id 文章id
 * @return {*}
 */
func (this *ApiArticle) Attach(call_index string, article_id int64, type_id int32) ([]*model.CmsAttach, error) {
	if call_index != "" {
		article, articleDo := query.CmsArticleDo()
		articleDo.Where(article.CallIndex.Eq(call_index), article.Status.Eq(2)).Pluck(article.ArticleID, &article_id)
	}
	attachMdl, attachDo := query.CmsAttachDo()
	return attachDo.Where(attachMdl.TableName_.Eq("article"), attachMdl.RecordID.Eq(article_id), attachMdl.TypeID.Eq(type_id), attachMdl.IsShow.Eq(1)).Order(attachMdl.SortID).Find()
}

/**
 * @description: Click 点击数+1
 * @param {string} call_index 文章调用别名
 * @param {int64} article_id 文章id
 * @return {*}
 */
func (this *ApiArticle) Click(call_index string, article_id int64) error {
	mdl, do := query.CmsArticleDo()
	if call_index != "" {
		do.Where(mdl.CallIndex.Eq(call_index), mdl.Status.Eq(2)).Updates(map[string]interface{}{
			mdl.Click.ColumnName().String():      gorm.Expr("click + ?", 1),
			mdl.UpdateTime.ColumnName().String(): time.Now(),
		})
	}
	do.Where(mdl.ArticleID.Eq(article_id), mdl.Status.Eq(2)).Updates(map[string]interface{}{
		mdl.Click.ColumnName().String():      gorm.Expr("click + ?", 1),
		mdl.UpdateTime.ColumnName().String(): time.Now(),
	})
	return nil
}

/**
 * @description: Like 点赞数+1
 * @param {string} call_index 文章调用别名
 * @param {int64} article_id 文章id
 * @return {*}
 */
func (this *ApiArticle) Like(call_index string, article_id int64) error {
	mdl, do := query.CmsArticleDo()
	do.Where(mdl.ArticleID.Eq(article_id), mdl.Status.Eq(2)).Updates(map[string]interface{}{
		mdl.LikeCount.ColumnName().String():  gorm.Expr("like_count + ?", 1),
		mdl.UpdateTime.ColumnName().String(): time.Now(),
	})
	return nil
}

/**
 * @description: AlbumClick 点击数+1
 * @param {int64} article_id 文章id
 * @param {int64} ablum_id 图片id
 * @return {*}
 */
func (this *ApiArticle) AlbumClick(article_id, ablum_id int64) error {
	mdl, do := query.CmsAlbumDo()
	do.Where(mdl.AlbumID.Eq(ablum_id), mdl.IsShow.Eq(1)).Updates(map[string]interface{}{
		mdl.Click.ColumnName().String():      gorm.Expr("click + ?", 1),
		mdl.UpdateTime.ColumnName().String(): time.Now(),
	})
	return nil
}
