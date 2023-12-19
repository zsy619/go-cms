package service

import (
	"errors"
	"time"

	"github.com/beego/beego/v2/core/logs"

	"haedu.gov.cn/cms/app/cms/domain"
	"haedu.gov.cn/cms/app/cms/mapper"
)

func (svc *CmsArticle) PropertyPaginate(page, limit int, parentId, articleId int64, status int32, callIndex, title string) ([]*domain.CmsArticleProperty, int64, error) {
	mdl, do := mapper.CmsArticlePropertyDo()
	if parentId > 0 {
		do = do.Where(mdl.ParentID.Eq(parentId))
	}
	if articleId >= 0 {
		do = do.Where(mdl.ArticleID.Eq(articleId))
	}
	if status >= 0 {
		do = do.Where(mdl.Status.Eq(status))
	}
	if callIndex != "" {
		do = do.Where(mdl.CallIndex.Like("%" + callIndex + "%"))
	}
	if title != "" {
		do = do.Where(mdl.Title.Like("%" + title + "%"))
	}
	return do.Where(mdl.IsDeleted.Is(false)).Order(mdl.SortID).FindByPage((page-1)*limit, limit)
}

func (svc *CmsArticle) PropertyFind(propertyId int64) (*domain.CmsArticleProperty, error) {
	mdl, do := mapper.CmsArticlePropertyDo()
	return do.Where(mdl.IsDeleted.Is(false)).Where(mdl.PropertyID.Eq(propertyId)).First()
}

func (svc *CmsArticle) PropertySave(input *domain.CmsArticleProperty) error {
	mdl, do := mapper.CmsArticlePropertyDo()
	if input.CallIndex != "" {
		if count, _ := do.Where(mdl.PropertyID.Neq(input.PropertyID), mdl.CallIndex.Eq(input.CallIndex)).Count(); count > 0 {
			return errors.New("调用别名重复")
		}
	}

	var err error
	if input.PropertyID <= 0 {
		input.CreateTime = time.Now()
		err = do.Create(input)
	} else {
		input.UpdateTime = time.Now()
		_, err = do.Where(mdl.PropertyID.Eq(input.PropertyID)).Updates(map[string]interface{}{
			mdl.ParentID.ColumnName().String():   input.ParentID,
			mdl.Title.ColumnName().String():      input.Title,
			mdl.CallIndex.ColumnName().String():  input.CallIndex,
			mdl.Value.ColumnName().String():      input.Value,
			mdl.SortID.ColumnName().String():     input.SortID,
			mdl.Status.ColumnName().String():     input.Status,
			mdl.BelongTo.ColumnName().String():   input.BelongTo,
			mdl.UpdateID.ColumnName().String():   input.UpdateID,
			mdl.UpdateName.ColumnName().String(): input.UpdateName,
			mdl.UpdateTime.ColumnName().String(): input.UpdateTime,
		})
		if err != nil {
			logs.Error(err.Error())
		}
	}
	return err
}

func (svc *CmsArticle) PropertySaveSortId(PropertyId int64, sortId int32) error {
	mdl, do := mapper.CmsArticlePropertyDo()
	_, err := do.Where(mdl.PropertyID.Eq(PropertyId)).UpdateColumns(
		map[string]interface{}{
			mdl.SortID.ColumnName().String():     sortId,
			mdl.UpdateTime.ColumnName().String(): time.Now(),
		},
	)
	return err
}

func (svc *CmsArticle) PropertyDestroy(PropertyId int64) error {
	mdl, do := mapper.CmsArticlePropertyDo()
	if _, err := do.Where(mdl.PropertyID.Eq(PropertyId)).Delete(); err != nil {
		return err
	}
	return nil
}

func (svc *CmsArticle) PropertyChangeStatus(PropertyId int64, status int32) error {
	mdl, do := mapper.CmsArticlePropertyDo()
	_, err := do.Where(mdl.PropertyID.Eq(PropertyId)).UpdateColumns(
		map[string]interface{}{
			mdl.Status.ColumnName().String():     status,
			mdl.UpdateTime.ColumnName().String(): time.Now(),
		},
	)
	return err
}
