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

func (this *CmsArticle) ArticlePaginate(page, limit int, channelId, categoryId int64, title string) ([]*model.CmsArticle, int64, error) {
	mdl, do := query.CmsArticleDo()
	if channelId > 0 {
		do = do.Where(mdl.ChannelID.Eq(channelId))
	}
	if channelId > 0 {
		do = do.Where(mdl.ChannelID.Eq(channelId))
	}
	if title != "" {
		do = do.Where(mdl.Title.Like("%" + title + "%"))
	}
	return do.FindByPage((page-1)*limit, limit)
}

func (this *CmsArticle) CategoryPagination(page, limit int, channelId int64, title, callIndex string) ([]*model.CmsArticleCategory, int64, error) {
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
