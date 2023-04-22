package biz

import (
	"errors"
	"fmt"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/cms/app/biz/bizmodel"
	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/app/dal/query"
	"haedu.gov.cn/tools/xgeneric"
)

/**
 * @description: CmsArticle 文章
 * @return {*}
 */
type CmsArticle struct{}

/**
 * @description: NewCmsArticle 新建文章
 * @return {*}
 */
func NewCmsArticle() *CmsArticle {
	return &CmsArticle{}
}

/**
 * @description: ArticlePaginate 文章分页
 * @param {*} page 页码
 * @param {int} limit 每页数量
 * @param {*} channelId 频道ID
 * @param {int64} categoryId 分类ID
 * @param {*} title 标题
 * @param {string} callIndex 别名
 * @param {int32} status 状态
 * @return {*}
 */
func (this *CmsArticle) ArticlePaginate(page, limit int, channelId, categoryId int64, title, callIndex string, status int32) ([]*model.CmsArticle, int64, error) {
	mdl, do := query.CmsArticleDo()
	if channelId > 0 {
		do = do.Where(mdl.ChannelID.Eq(channelId))
	}
	if categoryId > 0 {
		do = do.Where(mdl.CategoryID.Eq(categoryId))
	}
	if title != "" {
		do = do.Where(mdl.Title.Like("%" + title + "%"))
	}
	if callIndex != "" {
		do = do.Where(mdl.CallIndex.Like("%" + callIndex + "%"))
	}
	if status >= 0 {
		do = do.Where(mdl.Status.Eq(status))
	}
	return do.Order(mdl.IsTop.Desc()).Order(mdl.SortID).FindByPage((page-1)*limit, limit)
}

/**
 * @description: ArticleFind 获取
 * @param {int64} articleId 文章ID
 * @return {*}
 */
func (this *CmsArticle) ArticleFind(articleId int64) (*model.CmsArticle, error) {
	mdl, do := query.CmsArticleDo()
	return do.Where(mdl.ArticleID.Eq(articleId)).First()
}

/**
 * @description: ArticleSave 保存或更新
 * @param {*model.CmsArticle} input 文章
 * @return {*}
 */
func (this *CmsArticle) ArticleSave(input *model.CmsArticle) error {
	mdl, do := query.CmsArticleDo()
	if input.CallIndex != "" {
		if count, _ := do.Where(mdl.ArticleID.Neq(input.ArticleID), mdl.CallIndex.Eq(input.CallIndex)).Count(); count > 0 {
			return errors.New("调用别名重复")
		}
	}
	{
		chnMdl, chnDo := query.CmsSiteChannelDo()
		var siteId int64
		chnDo.Where(chnMdl.ChannelID.Eq(input.ChannelID)).Pluck(chnMdl.SiteID, &siteId)
		input.SiteID = siteId
	}
	var err error
	input.UpdateTime = time.Now()
	input.IcoUrl2 = xgeneric.IFF(input.IcoUrl1 == "", "", input.IcoUrl2)
	input.ImgUrl2 = xgeneric.IFF(input.ImgUrl1 == "", "", input.ImgUrl2)
	if input.ArticleID <= 0 {
		input.CreateTime = time.Now()
		err = do.Create(input)
	} else {
		_, err = do.Where(mdl.ArticleID.Eq(input.ArticleID)).Updates(map[string]interface{}{
			mdl.ChannelID.ColumnName().String():      input.ChannelID,
			mdl.SiteID.ColumnName().String():         input.SiteID,
			mdl.CategoryID.ColumnName().String():     input.CategoryID,
			mdl.Title.ColumnName().String():          input.Title,
			mdl.SubTitle.ColumnName().String():       input.SubTitle,
			mdl.IcoUrl1.ColumnName().String():        input.IcoUrl1,
			mdl.IcoUrl2.ColumnName().String():        input.IcoUrl2,
			mdl.CallIndex.ColumnName().String():      input.CallIndex,
			mdl.Source.ColumnName().String():         input.Source,
			mdl.Author.ColumnName().String():         input.Author,
			mdl.LinkURL.ColumnName().String():        input.LinkURL,
			mdl.ImgUrl1.ColumnName().String():        input.ImgUrl1,
			mdl.ImgUrl2.ColumnName().String():        input.ImgUrl2,
			mdl.SeoTitle.ColumnName().String():       input.SeoTitle,
			mdl.SeoKeyword.ColumnName().String():     input.SeoKeyword,
			mdl.SeoDescription.ColumnName().String(): input.SeoDescription,
			mdl.Tags.ColumnName().String():           input.Tags,
			mdl.Summary.ColumnName().String():        input.Summary,
			mdl.Content.ColumnName().String():        input.Content,
			mdl.SortID.ColumnName().String():         input.SortID,
			mdl.IsLock.ColumnName().String():         input.IsLock,
			mdl.IsComment.ColumnName().String():      input.IsComment,
			mdl.IsTop.ColumnName().String():          input.IsTop,
			mdl.IsRed.ColumnName().String():          input.IsRed,
			mdl.IsHot.ColumnName().String():          input.IsHot,
			mdl.IsSlide.ColumnName().String():        input.IsSlide,
			mdl.Status.ColumnName().String():         input.Status,
			mdl.PublishTime.ColumnName().String():    input.PublishTime,
			mdl.Topic.ColumnName().String():          input.Topic,
			mdl.Template.ColumnName().String():       input.Template,
			mdl.UpdateTime.ColumnName().String():     input.UpdateTime,
		})
		if err != nil {
			logs.Error(err.Error())
		}
	}
	return err
}

/**
 * @description: ArticleSaveStatus 更新状态
 * @param {int64} articleId 文章ID
 * @param {int32} sortId 排序ID
 * @return {*}
 */
func (this *CmsArticle) ArticleSaveSortId(articleId int64, sortId int32) error {
	mdl, do := query.CmsArticleDo()
	_, err := do.Where(mdl.ArticleID.Eq(articleId)).UpdateColumns(
		map[string]interface{}{
			mdl.SortID.ColumnName().String():     sortId,
			mdl.UpdateTime.ColumnName().String(): time.Now(),
		},
	)
	return err
}

/**
 * @description: ArticleDestory 删除
 * @param {int64} articleId 文章ID
 * @return {*}
 */
func (this *CmsArticle) ArticleDestory(articleId int64) error {
	// // 删除附件
	// {
	// 	attachMdl, attachDo := query.CmsAttachDo()
	// 	if attachList, err := attachDo.Where(attachMdl.TableName_.Eq("article"), attachMdl.RecordID.Eq(articleId)).Find(); err != nil {
	// 		logs.Error(err.Error())
	// 	} else {
	// 		for _, attach := range attachList {
	// 			if attach.OriginalPath != "" && strings.HasPrefix(attach.OriginalPath, "/Uploads/") {
	// 				os.Remove(attach.OriginalPath[1:])
	// 			}
	// 		}
	// 	}
	// 	if _, err := attachDo.Where(attachMdl.TableName_.Eq("article"), attachMdl.RecordID.Eq(articleId)).Delete(); err != nil {
	// 		return err
	// 	}
	// }
	// // 删除相册
	// {
	// 	albumMdl, albumDo := query.CmsAlbumDo()
	// 	if albumList, err := albumDo.Where(albumMdl.TableName_.Eq("article"), albumMdl.RecordID.Eq(articleId)).Find(); err != nil {
	// 		logs.Error(err.Error())
	// 	} else {
	// 		for _, album := range albumList {
	// 			if album.OriginalPath != "" && strings.HasPrefix(album.OriginalPath, "/Uploads/") {
	// 				os.Remove(album.OriginalPath[1:])
	// 			}
	// 		}
	// 	}
	// 	if _, err := albumDo.Where(albumMdl.TableName_.Eq("article"), albumMdl.RecordID.Eq(articleId)).Delete(); err != nil {
	// 		return err
	// 	}
	// }
	mdl, do := query.CmsArticleDo()
	if _, err := do.Where(mdl.ArticleID.Eq(articleId)).Delete(); err != nil {
		return err
	}
	return nil
}

/**
 * @description: CategoryPaginate 分页
 * @param {*} page 页码
 * @param {int} limit 每页数量
 * @param {int64} channelId 频道ID
 * @param {*} title 标题
 * @param {string} callIndex 别名
 * @return {*}
 */
func (this *CmsArticle) CategoryPaginate(page, limit int, channelId int64, title, callIndex string) ([]*model.CmsArticleCategory, int64, error) {
	mdl, do := query.CmsArticleCategoryDo()
	if channelId > 0 {
		do = do.Where(mdl.ChannelID.Eq(channelId))
	}
	if title != "" {
		do = do.Where(mdl.Title.Like("%" + title + "%"))
	}
	if callIndex != "" {
		do = do.Where(mdl.CallIndex.Like("%" + callIndex + "%"))
	}
	return do.Order(mdl.SortID).FindByPage((page-1)*limit, limit)
}

/**
 * @description: CategoryFind 获取
 * @param {int64} categoryId 分类ID
 * @return {*}
 */
func (this *CmsArticle) CategoryFind(categoryId int64) (*model.CmsArticleCategory, error) {
	mdl, do := query.CmsArticleCategoryDo()
	return do.Where(mdl.CategoryID.Eq(categoryId)).First()
}

/**
 * @description: CategorySave 保存或更新
 * @param {*model.CmsArticleCategory} input 分类
 * @return {*}
 */
func (this *CmsArticle) CategorySave(input *model.CmsArticleCategory) error {
	mdl, do := query.CmsArticleCategoryDo()
	if input.CallIndex != "" {
		if count, _ := do.Where(mdl.CategoryID.Neq(input.CategoryID), mdl.CallIndex.Eq(input.CallIndex)).Count(); count > 0 {
			return errors.New("调用别名重复")
		}
	}
	{
		chnMdl, chnDo := query.CmsSiteChannelDo()
		var siteId int64
		chnDo.Where(chnMdl.ChannelID.Eq(input.ChannelID)).Pluck(chnMdl.SiteID, &siteId)
		input.SiteID = siteId
		var classLayer int32
		do.Where(mdl.CategoryID.Eq(input.ParentID)).Pluck(mdl.ClassLayer, &classLayer)
		classLayer++
		input.ClassLayer = classLayer
	}
	var err error
	input.UpdateTime = time.Now()
	input.ImgUrl2 = xgeneric.IFF(input.ImgUrl1 == "", "", input.ImgUrl2)
	if input.CategoryID <= 0 {
		input.CreateTime = time.Now()
		err = do.Create(input)
	} else {
		_, err = do.Where(mdl.CategoryID.Eq(input.CategoryID)).Updates(map[string]interface{}{
			mdl.ChannelID.ColumnName().String():      input.ChannelID,
			mdl.SiteID.ColumnName().String():         input.SiteID,
			mdl.Title.ColumnName().String():          input.Title,
			mdl.CallIndex.ColumnName().String():      input.CallIndex,
			mdl.ClassLayer.ColumnName().String():     input.ClassLayer,
			mdl.LinkURL.ColumnName().String():        input.LinkURL,
			mdl.ImgUrl1.ColumnName().String():        input.ImgUrl1,
			mdl.ImgUrl2.ColumnName().String():        input.ImgUrl2,
			mdl.SeoTitle.ColumnName().String():       input.SeoTitle,
			mdl.SeoKeyword.ColumnName().String():     input.SeoKeyword,
			mdl.SeoDescription.ColumnName().String(): input.SeoDescription,
			mdl.Content.ColumnName().String():        input.Content,
			mdl.SortID.ColumnName().String():         input.SortID,
			mdl.IsSearch.ColumnName().String():       input.IsSearch,
			mdl.IsShow.ColumnName().String():         input.IsShow,
			mdl.Status.ColumnName().String():         input.Status,
			mdl.Template.ColumnName().String():       input.Template,
			mdl.UpdateTime.ColumnName().String():     input.UpdateTime,
		})
	}
	return err
}

/**
 * @description: CategoryAutoUrl 自动生成Url
 * @param {int64} channelId 频道ID
 * @return {*}
 */
func (this *CmsArticle) CategoryAutoUrl(channelId int64) error {
	mdl, do := query.CmsArticleCategoryDo()
	categories, err := do.Where(mdl.ChannelID.Eq(channelId), mdl.LinkURL.Eq("")).Find()
	if err != nil {
		return err
	}
	channelMdl, channelDo := query.CmsSiteChannelDo()

	find, err := channelDo.Where(channelMdl.ChannelID.Eq(channelId)).Select(channelMdl.SiteID, channelMdl.Name).First()
	if err != nil {
		return err
	}
	siteFlag := ""
	siteMdl, siteDo := query.CmsSiteDo()
	if err := siteDo.Where(siteMdl.SiteID.Eq(find.SiteID)).Pluck(siteMdl.Flag, &siteFlag); err != nil {
		return err
	}
	var errOut error
	for _, category := range categories {
		linkUrl := fmt.Sprintf("/%s/%s/%s", siteFlag, find.Name, category.CallIndex)
		_, errOut = do.Where(mdl.CategoryID.Eq(category.CategoryID)).UpdateColumns(map[string]interface{}{
			mdl.LinkURL.ColumnName().String():    linkUrl,
			mdl.UpdateTime.ColumnName().String(): time.Now(),
		})
		if errOut != nil {
			return errOut
		}
	}
	return nil
}

/**
 * @description: CategorySaveSortId 保存排序
 * @param {int64} categoryId 分类ID
 * @param {int32} sortId 排序
 * @return {*}
 */
func (this *CmsArticle) CategorySaveSortId(categoryId int64, sortId int32) error {
	mdl, do := query.CmsArticleCategoryDo()
	_, err := do.Where(mdl.CategoryID.Eq(categoryId)).UpdateColumns(
		map[string]interface{}{
			mdl.SortID.ColumnName().String():     sortId,
			mdl.UpdateTime.ColumnName().String(): time.Now(),
		},
	)
	return err
}

/**
 * @description: CategoryDestory 删除
 * @param {int64} categoryId 分类ID
 * @return {*}
 */
func (this *CmsArticle) CategoryDestory(categoryId int64) error {
	mdl, do := query.CmsArticleCategoryDo()
	if count, _ := do.Where(mdl.ParentID.Eq(categoryId)).Count(); count > 0 {
		return errors.New("请先删除子分类")
	}
	artMdl, artDo := query.CmsArticleDo()
	if count, _ := artDo.Where(artMdl.CategoryID.Eq(categoryId)).Count(); count > 0 {
		return errors.New("请先删除分类下的文章")
	}
	if _, err := do.Where(mdl.CategoryID.Eq(categoryId)).Delete(); err != nil {
		return err
	}
	return nil
}

/**
 * @description: ArticleClone 克隆
 * @param {int64} articleId 文章ID
 * @return {*}
 */
func (this *CmsArticle) ArticleClone(articleId int64) (int64, error) {
	mdl, do := query.CmsArticleDo()
	art, err := do.Where(mdl.ArticleID.Eq(articleId)).First()
	if err != nil {
		return 0, err
	}
	art.ArticleID = 0
	art.CreateTime = time.Now()
	art.UpdateTime = time.Now()
	art.Status = 0
	err = do.Create(art)
	return art.ArticleID, err
}

/**
 * @description: ArticleChangeStatus 修改状态
 * @param {int64} articleId 文章ID
 * @param {int32} status 状态
 * @return {*}
 */
func (this *CmsArticle) ArticleChangeStatus(articleId int64, status int32) error {
	mdl, do := query.CmsArticleDo()
	_, err := do.Where(mdl.ArticleID.Eq(articleId)).UpdateColumns(
		map[string]interface{}{
			mdl.Status.ColumnName().String():     status,
			mdl.UpdateTime.ColumnName().String(): time.Now(),
		},
	)
	return err
}

/**
 * @description: CategoryTree 分类树
 * @param {*} channelId 频道ID
 * @param {int64} categoryId 分类ID
 * @return {*}
 */
func (this *CmsArticle) CategoryTree(channelId, categoryId int64) ([]*bizmodel.TreeNode, error) {
	out := make([]*bizmodel.TreeNode, 0)
	mdl, do := query.CmsArticleCategoryDo()
	list, err := do.Where(mdl.ChannelID.Eq(channelId), mdl.ParentID.Eq(0)).Order(mdl.SortID).Find()
	if err != nil {
		return out, err
	}
	for _, item := range list {
		child := &bizmodel.TreeNode{
			Id:       item.CategoryID,
			Name:     item.Title,
			Open:     true,
			Checked:  item.CategoryID == categoryId,
			Selected: item.CategoryID == categoryId,
			Children: nil,
		}
		children, _ := this.CategoryTreeByParentId(item.CategoryID, categoryId)
		if children != nil {
			child.Children = children
		}
		out = append(out, child)
	}
	return out, nil
}

/**
 * @description: CategoryTreeByParentId 根据父ID获取分类树
 * @param {*} parentId 父ID
 * @param {int64} categoryId 分类ID
 * @return {*}
 */
func (this *CmsArticle) CategoryTreeByParentId(parentId, categoryId int64) ([]*bizmodel.TreeNode, error) {
	out := make([]*bizmodel.TreeNode, 0)
	mdl, do := query.CmsArticleCategoryDo()
	list, err := do.Where(mdl.ParentID.Eq(parentId)).Order(mdl.SortID).Find()
	if err != nil {
		return out, err
	}
	for _, item := range list {
		child := &bizmodel.TreeNode{
			Id:       item.CategoryID,
			Name:     item.Title,
			Checked:  item.CategoryID == categoryId,
			Selected: item.CategoryID == categoryId,
			Children: nil,
		}
		children, _ := this.CategoryTreeByParentId(item.CategoryID, categoryId)
		if children != nil {
			child.Children = children
		}
		out = append(out, child)
	}
	return out, nil
}
