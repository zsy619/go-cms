package biz

import (
	"errors"
	"time"

	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/app/dal/query"
)

type CmsArticle struct{}

func NewCmsArticle() *CmsArticle {
	return &CmsArticle{}
}

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

// ArticleFind 获取
func (this *CmsArticle) ArticleFind(articleId int64) (*model.CmsArticle, error) {
	mdl, do := query.CmsArticleDo()
	return do.Where(mdl.ArticleID.Eq(articleId)).First()
}

// ArticleSave 保存或更新
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
			mdl.CallIndex.ColumnName().String():      input.CallIndex,
			mdl.Source.ColumnName().String():         input.Source,
			mdl.Author.ColumnName().String():         input.Author,
			mdl.LinkURL.ColumnName().String():        input.LinkURL,
			mdl.ImgURL.ColumnName().String():         input.ImgURL,
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
			mdl.UpdateTime.ColumnName().String():     input.UpdateTime,
		})
	}
	return err
}

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

// ArticleDestory 删除
func (this *CmsArticle) ArticleDestory(articleId int64) error {
	attachMdl, attachDo := query.CmsArticleAttachDo()
	if _, err := attachDo.Where(attachMdl.ArticleID.Eq(articleId)).Delete(); err != nil {
		return err
	}
	albumMdl, albumDo := query.CmsArticleAlbumDo()
	if _, err := albumDo.Where(albumMdl.ArticleID.Eq(articleId)).Delete(); err != nil {
		return err
	}
	mdl, do := query.CmsArticleDo()
	if _, err := do.Where(mdl.ArticleID.Eq(articleId)).Delete(); err != nil {
		return err
	}
	return nil
}

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

// CategoryFind 获取
func (this *CmsArticle) CategoryFind(categoryId int64) (*model.CmsArticleCategory, error) {
	mdl, do := query.CmsArticleCategoryDo()
	return do.Where(mdl.CategoryID.Eq(categoryId)).First()
}

// CategorySave 保存或更新
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
			mdl.ImgURL.ColumnName().String():         input.ImgURL,
			mdl.SeoTitle.ColumnName().String():       input.SeoTitle,
			mdl.SeoKeyword.ColumnName().String():     input.SeoKeyword,
			mdl.SeoDescription.ColumnName().String(): input.SeoDescription,
			mdl.Content.ColumnName().String():        input.Content,
			mdl.SortID.ColumnName().String():         input.SortID,
			mdl.Status.ColumnName().String():         input.Status,
			mdl.UpdateTime.ColumnName().String():     input.UpdateTime,
		})
	}
	return err
}

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

// CategoryDestory 删除
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

// ArticleClone 克隆
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

// ArticleChangeStatus 修改状态
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

func (this *CmsArticle) CategoryTree(channelId, categoryId int64) ([]*TreeNode, error) {
	out := make([]*TreeNode, 0)
	mdl, do := query.CmsArticleCategoryDo()
	list, err := do.Where(mdl.ChannelID.Eq(channelId), mdl.ParentID.Eq(0)).Order(mdl.SortID).Find()
	if err != nil {
		return out, err
	}
	for _, item := range list {
		child := &TreeNode{
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

func (this *CmsArticle) CategoryTreeByParentId(parentId, categoryId int64) ([]*TreeNode, error) {
	out := make([]*TreeNode, 0)
	mdl, do := query.CmsArticleCategoryDo()
	list, err := do.Where(mdl.ParentID.Eq(parentId)).Order(mdl.SortID).Find()
	if err != nil {
		return out, err
	}
	for _, item := range list {
		child := &TreeNode{
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
