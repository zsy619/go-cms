package biz

import (
	"fmt"
	"strconv"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"gorm.io/gorm"
	"haedu.gov.cn/cms/app/biz/bizmodel"
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
func (this *ApiArticle) CategoryNav(channel_name string, channel_id int64, call_index string, category_id int64, article_id int64) ([]*bizmodel.ApiCategoryNav, error) {
	if article_id > 0 && category_id <= 0 {
		cacheKey := fmt.Sprintf("ApiArticle_CategoryNav_%d", article_id)
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
		cacheKey := fmt.Sprintf("ApiArticle_CategoryNav_%s_%d", call_index, category_id)
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
		return item.([]*bizmodel.ApiCategoryNav), nil
	}
	outResult := []*bizmodel.ApiCategoryNav{}
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
			channelItem := &bizmodel.ApiCategoryNav{
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
			categoryItem := &bizmodel.ApiCategoryNav{
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
func (this *ApiArticle) CategoryFind(channel_name string) ([]*bizmodel.ApiCategoryFindModel, int64, error) {
	cacheKey := fmt.Sprintf("ApiArticle_CategoryFind_%s", channel_name)
	if found, item := ApiCache.Get(cacheKey); found {
		list := item.([]*bizmodel.ApiCategoryFindModel)
		return list, int64(len(list)), nil
	}
	list := make([]*bizmodel.ApiCategoryFindModel, 0)
	_, do := query.CmsArticleCategoryDo()
	sqlSelect := "b.`name` as channel_name,b.title as channel_title,a.category_id,a.parent_id,a.site_id,a.channel_id,a.title,a.call_index,a.class_layer,a.link_url,a.img_url1,a.img_url2,a.sort_id,a.is_show,a.is_search,a.is_deleted"
	sql := "SELECT " + sqlSelect + " FROM cms_article_category a LEFT JOIN cms_site_channel b ON a.channel_id=b.channel_id WHERE a.is_deleted=0 AND a.`status`=2 AND a.`is_show`=1 AND b.`name`=? ORDER BY a.sort_id"
	err := do.UnderlyingDB().Raw(sql, channel_name).Scan(&list).Error
	if err != nil {
		return nil, 0, err
	}
	ApiCache.Set(cacheKey, list, 1800)
	return list, int64(len(list)), nil
}

/**
 * @description: CategoryOne 获取栏目详情
 * @param {int64} category_id 栏目ID
 * @param {string} call_index 栏目别名
 * @return {*}
 */
func (this *ApiArticle) CategoryOne(category_id int64, call_index string) (*bizmodel.ApiCategoryOneModel, error) {
	cacheKey := fmt.Sprintf("ApiArticle_CategoryOne_%d_%s", category_id, call_index)
	if found, item := ApiCache.Get(cacheKey); found {
		return item.(*bizmodel.ApiCategoryOneModel), nil
	}
	list := &bizmodel.ApiCategoryOneModel{}
	_, do := query.CmsArticleCategoryDo()
	sqlSelect := `b.name as channel_name,b.title as channel_title` +
		`,a.category_id,a.parent_id,a.site_id,a.channel_id,a.title,a.call_index,a.class_layer,a.link_url,a.target,a.img_url1,a.img_url2,a.sort_id,a.is_show,a.is_search,a.is_deleted,a.seo_title,a.seo_keyword,a.seo_description,a.content` +
		`,case when a.tmpl_cat='' then b.tmpl_cat else a.tmpl_cat end tmpl_cat,case when a.tmpl_lst='' then b.tmpl_lst else a.tmpl_lst end tmpl_lst,case when a.tmpl_dtl='' then b.tmpl_dtl else a.tmpl_dtl end tmpl_dtl`
	sql := `SELECT ` + sqlSelect + ` FROM cms_article_category a LEFT JOIN cms_site_channel b ON a.channel_id=b.channel_id WHERE a.is_deleted=0 AND a.status=2 AND a.is_show=1` +
		xgeneric.IFF(category_id > 0, ` AND a.category_id=`+strconv.FormatInt(category_id, 10), ``) +
		xgeneric.IFF(call_index != "", ` AND a.call_index='`+call_index+`'`, ``)
	err := do.UnderlyingDB().Raw(sql).Scan(&list).Error
	if err != nil {
		return nil, err
	}
	ApiCache.Set(cacheKey, list, 1800)
	return list, nil
}

/**
 * @description: Find 获取文章列表
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
 * @return {*}
 */
func (this *ApiArticle) Find(limit int, channel_id int64, channel_name string, category_id int64, call_index string, is_top, is_red, is_hot, is_slide int, order_by string) ([]*bizmodel.ApiArticleListModel, int64, error) {
	order_by = xgeneric.IFF(order_by == "", "a.sort_id", order_by)
	if limit <= 0 {
		limit = 10
	}
	cacheKey := fmt.Sprintf("ApiArticle_Find_%d_%d_%s_%d_%s_%d_%d_%d_%d_%s", limit, channel_id, channel_name, category_id, call_index, is_top, is_red, is_hot, is_slide, order_by)
	if found, item := ApiCache.Get(cacheKey); found {
		list := item.([]*bizmodel.ApiArticleListModel)
		return list, int64(len(list)), nil
	}
	list := make([]*bizmodel.ApiArticleListModel, 0)
	_, do := query.CmsArticleDo()

	sql := bizmodel.ApiArticleListModel_Table +
		xgeneric.IFF(channel_id <= 0, "", " And b.channel_id="+strconv.FormatInt(channel_id, 10)) +
		xgeneric.IFF(channel_name == "", "", " And c.name='"+channel_name+"'") +
		xgeneric.IFF(category_id <= 0, "", " And b.category_id="+strconv.FormatInt(category_id, 10)) +
		xgeneric.IFF(call_index == "", "", " And b.call_index='"+call_index+"'") +
		xgeneric.IFF(is_top <= 0, "", " And a.is_top="+strconv.Itoa(is_hot)) +
		xgeneric.IFF(is_red <= 0, "", " And a.is_red="+strconv.Itoa(is_hot)) +
		xgeneric.IFF(is_hot <= 0, "", " And a.is_hot="+strconv.Itoa(is_hot)) +
		xgeneric.IFF(is_slide <= 0, "", " And a.is_slide="+strconv.Itoa(is_slide)) +
		" ORDER BY a.is_top DESC," + order_by +
		" LIMIT ?"
	fmt.Println(order_by)
	err := do.UnderlyingDB().Raw(sql, limit).Scan(&list).Error
	if err == nil {
		ApiCache.Set(cacheKey, list, 1800)
	}
	return list, int64(len(list)), err
}

/**
 * @description: FindNew 获取最新文章列表
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
 * @return {*}
 */
func (this *ApiArticle) FindNew(limit int, channel_id int64, channel_name string, category_id int64, call_index string, is_top, is_red, is_hot, is_slide int, order_by string) ([]*bizmodel.ApiArticleListModel, int64, error) {
	order_by = xgeneric.IFF(order_by == "", "a.sort_id", order_by)
	if limit <= 0 {
		limit = 10
	}
	cacheKey := fmt.Sprintf("ApiArticle_FindNew_%d_%d_%s_%d_%s_%d_%d_%d_%d_%s", limit, channel_id, channel_name, category_id, call_index, is_top, is_red, is_hot, is_slide, order_by)
	if found, item := ApiCache.Get(cacheKey); found {
		find := item.([]*bizmodel.ApiArticleListModel)
		return find, int64(len(find)), nil
	}
	outArticle := make([]*bizmodel.ApiArticleListModel, 0)
	_, do := query.CmsArticleDo()
	sql := bizmodel.ApiArticleListModel_Table +
		xgeneric.IFF(channel_id <= 0, "", " And b.channel_id="+strconv.FormatInt(channel_id, 10)) +
		xgeneric.IFF(channel_name == "", "", " And c.name='"+channel_name+"'") +
		xgeneric.IFF(category_id <= 0, "", " And b.category_id="+strconv.FormatInt(category_id, 10)) +
		xgeneric.IFF(call_index == "", "", " And b.call_index='"+call_index+"'") +
		xgeneric.IFF(is_top <= 0, "", " And a.is_top="+strconv.Itoa(is_hot)) +
		xgeneric.IFF(is_red <= 0, "", " And a.is_red="+strconv.Itoa(is_hot)) +
		xgeneric.IFF(is_hot <= 0, "", " And a.is_hot="+strconv.Itoa(is_hot)) +
		xgeneric.IFF(is_slide <= 0, "", " And a.is_slide="+strconv.Itoa(is_slide)) +
		" ORDER BY a.publish_time DESC," + order_by +
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
 * @param {string} channel_name 频道名称
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
func (this *ApiArticle) Paginate(page, limit int, channel_id int64, channel_name string, category_id int64, call_index string, keyword string, is_top, is_red, is_hot, is_slide, is_search int, order_by string) ([]*bizmodel.ApiArticleListModel, int64, error) {
	order_by = xgeneric.IFF(order_by == "", "a.sort_id", order_by)
	where := xgeneric.IFF(channel_id <= 0, "", " And b.channel_id="+strconv.FormatInt(channel_id, 10)) +
		xgeneric.IFF(channel_name == "", "", " And c.name='"+channel_name+"'") +
		xgeneric.IFF(category_id <= 0, "", " And b.category_id="+strconv.FormatInt(category_id, 10)) +
		xgeneric.IFF(call_index == "", "", " And b.call_index='"+call_index+"'") +
		xgeneric.IFF(is_top < 0, "", " And a.is_top="+strconv.Itoa(is_hot)) +
		xgeneric.IFF(is_red < 0, "", " And a.is_red="+strconv.Itoa(is_hot)) +
		xgeneric.IFF(is_hot < 0, "", " And a.is_hot="+strconv.Itoa(is_hot)) +
		xgeneric.IFF(is_slide < 0, "", " And a.is_slide="+strconv.Itoa(is_slide)) +
		xgeneric.IFF(is_search < 0, "", " And b.is_search="+strconv.Itoa(is_search)) +
		xgeneric.IFF(keyword == "", "", " And (a.title LIKE '%"+keyword+"%' OR a.summary LIKE '%"+keyword+"%')")

	sqlSelectCount := "COUNT(1) as count"
	sqlCount := "SELECT " + sqlSelectCount + " FROM cms_article a" +
		" LEFT JOIN cms_article_category b ON a.category_id=b.category_id" +
		" LEFT JOIN cms_site_channel c ON a.channel_id = c.channel_id" +
		" WHERE a.`status`=2 AND b.`status`=2" +
		where

	var count int64

	_, do := query.CmsArticleDo()
	if err := do.UnderlyingDB().Raw(sqlCount).Pluck("count", &count).Error; err != nil {
		return nil, 0, err
	}

	sql := bizmodel.ApiArticleListModel_Table +
		where +
		" ORDER BY a.is_top DESC," + order_by +
		" LIMIT ? OFFSET ?"
	list := make([]*bizmodel.ApiArticleListModel, 0)
	err := do.UnderlyingDB().Raw(sql, limit, (page-1)*limit).Scan(&list).Error

	return list, count, err
}

/**
 * @description: Get 根据article_id获取文章详情、相册、附件
 * @param {string} call_index 文章调用别名
 * @param {int64} article_id 文章id
 * @return {*}
 */
func (this *ApiArticle) Get(call_index string, article_id int64) (*bizmodel.ApiArticleOneModel, []*model.CmsAlbum, []*model.CmsAttach, error) {
	mdl, do := query.CmsArticleDo()
	if call_index != "" {
		if article_id <= 0 {
			if err := do.Where(mdl.CallIndex.Eq(call_index), mdl.Status.Eq(2)).Pluck(mdl.ArticleID, &article_id); err != nil {
				logs.Error("Get", err.Error())
			}
		}
	}
	field := `a.*,b.call_index as category_call_index,b.title as category_title,b.link_url as category_link_url,c.name as channel_name,c.title as channel_title,case when b.tmpl_dtl='' then c.tmpl_dtl else b.tmpl_dtl end tmpl_dtl`
	sql := `SELECT ` + field + ` FROM cms_article a LEFT JOIN cms_article_category b ON a.category_id = b.category_id LEFT JOIN cms_site_channel c ON a.channel_id = c.channel_id`
	sql += ` WHERE a.article_id=? AND a.is_deleted=0 AND a.status=2 AND b.status=2`
	list := &bizmodel.ApiArticleOneModel{}
	if err := do.Debug().UnderlyingDB().Raw(sql, article_id).Scan(&list).Error; err != nil {
		logs.Error("Get", err.Error())
		return list, []*model.CmsAlbum{}, []*model.CmsAttach{}, err
	}

	albumMdl, alblumDo := query.CmsAlbumDo()
	album, _ := alblumDo.Where(albumMdl.TableName_.Eq("article"), albumMdl.RecordID.Eq(article_id), albumMdl.IsShow.Eq(1)).Order(albumMdl.SortID).Find()
	if album == nil {
		album = []*model.CmsAlbum{}
	}
	attachMdl, attachDo := query.CmsAttachDo()
	attach, _ := attachDo.Where(attachMdl.TableName_.Eq("article"), attachMdl.RecordID.Eq(article_id), attachMdl.IsShow.Eq(1)).Order(attachMdl.SortID).Find()
	if attach == nil {
		attach = []*model.CmsAttach{}
	}
	return list, album, attach, nil
}

/**
 * @description: Get 根据article_id获取文章上一个、下一个
 * @param {string} call_index 栏目调用别名
 * @param {int64} category_id 栏目id
 * @param {int64} article_id 文章id
 * @return {*}
 */
func (this *ApiArticle) PrevNext(call_index string, category_id, article_id int64) (prev *bizmodel.ApiArticleOneModel, next *bizmodel.ApiArticleOneModel) {
	_, do := query.CmsArticleDo()
	field := `a.*,b.call_index as category_call_index,b.title as category_title,b.link_url as category_link_url,c.name as channel_name,c.title as channel_title,case when b.tmpl_dtl='' then c.tmpl_dtl else b.tmpl_dtl end tmpl_dtl`
	sql := `SELECT ` + field + ` FROM cms_article a LEFT JOIN cms_article_category b ON a.category_id = b.category_id LEFT JOIN cms_site_channel c ON a.channel_id = c.channel_id`
	sql += ` WHERE a.is_deleted=0 AND a.status=2 AND b.status=2` +
		xgeneric.IFF(category_id <= 0, "", " AND b.category_id="+strconv.FormatInt(category_id, 10)) +
		xgeneric.IFF(call_index == "", "", " AND b.call_index='"+call_index+"'") +
		` AND a.article_id<? ORDER BY a.sort_id DESC LIMIT 1`
	prev = &bizmodel.ApiArticleOneModel{}
	if err := do.Debug().UnderlyingDB().Raw(sql, article_id).Scan(&prev).Error; err != nil {
		logs.Error("PrevNext", err.Error())
	}
	sql = `SELECT ` + field + ` FROM cms_article a LEFT JOIN cms_article_category b ON a.category_id = b.category_id LEFT JOIN cms_site_channel c ON a.channel_id = c.channel_id`
	sql += ` WHERE a.is_deleted=0 AND a.status=2 AND b.status=2` +
		xgeneric.IFF(category_id <= 0, "", " AND b.category_id="+strconv.FormatInt(category_id, 10)) +
		xgeneric.IFF(call_index == "", "", " AND b.call_index='"+call_index+"'") +
		` AND a.article_id>? ORDER BY a.sort_id ASC LIMIT 1`
	next = &bizmodel.ApiArticleOneModel{}
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
func (this *ApiArticle) Article(call_index string, article_id int64) (*bizmodel.ApiArticleOneModel, error) {
	_, do := query.CmsArticleDo()
	field := `a.*,b.call_index as category_call_index,b.title as category_title,b.link_url as category_link_url,c.name as channel_name,c.title as channel_title,case when b.tmpl_dtl='' then c.tmpl_dtl else b.tmpl_dtl end tmpl_dtl`
	sql := `SELECT ` + field + ` FROM cms_article a LEFT JOIN cms_article_category b ON a.category_id = b.category_id LEFT JOIN cms_site_channel c ON a.channel_id = c.channel_id`
	list := &bizmodel.ApiArticleOneModel{}
	if call_index != "" {
		sql += ` WHERE a.call_index=? AND a.is_deleted=0 AND a.status=2 AND b.status=2`
		if err := do.Debug().UnderlyingDB().Raw(sql, call_index).Scan(&list).Error; err != nil {
			logs.Error("Get", err.Error())
			return list, err
		}
		return list, nil
	}
	sql += ` WHERE a.article_id=? AND a.is_deleted=0 AND a.status=2 AND b.status=2`
	if err := do.Debug().UnderlyingDB().Raw(sql, article_id).Scan(&list).Error; err != nil {
		logs.Error("Get", err.Error())
		return list, err
	}
	return list, nil
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
	var err error
	if call_index != "" {
		_, err = do.Where(mdl.CallIndex.Eq(call_index), mdl.Status.Eq(int32(StatusPass))).Updates(map[string]interface{}{
			mdl.Click.ColumnName().String():      gorm.Expr("click + ?", 1),
			mdl.UpdateTime.ColumnName().String(): time.Now(),
		})
	} else {
		_, err = do.Where(mdl.ArticleID.Eq(article_id), mdl.Status.Eq(2)).Updates(map[string]interface{}{
			mdl.Click.ColumnName().String():      gorm.Expr("click + ?", 1),
			mdl.UpdateTime.ColumnName().String(): time.Now(),
		})
	}
	return err
}

/**
 * @description: Like 点赞数+1
 * @param {string} call_index 文章调用别名
 * @param {int64} article_id 文章id
 * @return {*}
 */
func (this *ApiArticle) Like(call_index string, article_id int64) error {
	mdl, do := query.CmsArticleDo()
	_, err := do.Where(mdl.ArticleID.Eq(article_id), mdl.Status.Eq(int32(StatusPass))).Updates(map[string]interface{}{
		mdl.LikeCount.ColumnName().String():  gorm.Expr("like_count + ?", 1),
		mdl.UpdateTime.ColumnName().String(): time.Now(),
	})
	return err
}

/**
 * @description: AlbumClick 点击数+1
 * @param {int64} article_id 文章id
 * @param {int64} ablum_id 图片id
 * @return {*}
 */
func (this *ApiArticle) AlbumClick(article_id, ablum_id int64) error {
	mdl, do := query.CmsAlbumDo()
	_, err := do.Where(mdl.AlbumID.Eq(ablum_id), mdl.IsShow.Eq(1)).Updates(map[string]interface{}{
		mdl.Click.ColumnName().String():      gorm.Expr("click + ?", 1),
		mdl.UpdateTime.ColumnName().String(): time.Now(),
	})
	return err
}
