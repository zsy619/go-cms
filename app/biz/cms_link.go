package biz

import (
	"errors"
	"time"

	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/app/dal/query"
)

type CmsLink struct{}

func NewCmsLink() *CmsLink {
	return &CmsLink{}
}

// CategoryPaginate 分页查询
func (this *CmsLink) CategoryPaginate(page, limit int, siteId, channelId int64, title, callIndex string) ([]*model.CmsLinkCategory, int64, error) {
	mdl, do := query.CmsLinkCategoryDo()
	if title != "" {
		do = do.Where(mdl.Title.Like("%" + title + "%"))
	}
	if callIndex != "" {
		do = do.Where(mdl.CallIndex.Like("%" + callIndex + "%"))
	}
	return do.Order(mdl.SortID).FindByPage((page-1)*limit, limit)
}

// CategoryFind 获取
func (this *CmsLink) CategoryFind(categoryId int64) (*model.CmsLinkCategory, error) {
	mdl, do := query.CmsLinkCategoryDo()
	return do.Where(mdl.CategoryID.Eq(categoryId)).First()
}

// CategorySave 保存或更新
func (this *CmsLink) CategorySave(input *model.CmsLinkCategory) error {
	mdl, do := query.CmsLinkCategoryDo()
	if input.CallIndex != "" {
		if count, _ := do.Where(mdl.CategoryID.Neq(input.CategoryID), mdl.CallIndex.Eq(input.CallIndex)).Count(); count > 0 {
			return errors.New("调用别名重复")
		}
	}
	var err error
	input.UpdateTime = time.Now()
	if input.CategoryID <= 0 {
		input.CreateTime = time.Now()
		err = do.Create(input)
	} else {
		_, err = do.Where(mdl.CategoryID.Eq(input.CategoryID)).Updates(map[string]interface{}{
			mdl.Title.ColumnName().String():          input.Title,
			mdl.CallIndex.ColumnName().String():      input.CallIndex,
			mdl.LinkURL.ColumnName().String():        input.LinkURL,
			mdl.SeoTitle.ColumnName().String():       input.SeoTitle,
			mdl.SeoKeyword.ColumnName().String():     input.SeoKeyword,
			mdl.SeoDescription.ColumnName().String(): input.SeoDescription,
			mdl.Content.ColumnName().String():        input.Content,
			mdl.SortID.ColumnName().String():         input.SortID,
			mdl.UpdateTime.ColumnName().String():     input.UpdateTime,
		})
	}
	return err
}

func (this *CmsLink) CategorySaveSortId(categoryId int64, sortId int32) error {
	mdl, do := query.CmsLinkCategoryDo()
	_, err := do.Where(mdl.CategoryID.Eq(categoryId)).UpdateColumns(
		map[string]interface{}{
			mdl.SortID.ColumnName().String():     sortId,
			mdl.UpdateTime.ColumnName().String(): time.Now(),
		},
	)
	return err
}

// LinkClone 克隆
func (this *CmsLink) LinkClone(linkId int64) (int64, error) {
	mdl, do := query.CmsLinkDo()
	art, err := do.Where(mdl.LinkID.Eq(linkId)).First()
	if err != nil {
		return 0, err
	}
	art.LinkID = 0
	art.CreateTime = time.Now()
	art.UpdateTime = time.Now()
	art.Status = 0
	err = do.Create(art)
	return art.LinkID, err
}

// LinkChangeStatus 修改状态
func (this *CmsLink) LinkChangeStatus(linkId int64, status int32) error {
	mdl, do := query.CmsLinkDo()
	_, err := do.Where(mdl.LinkID.Eq(linkId)).UpdateColumns(
		map[string]interface{}{
			mdl.Status.ColumnName().String():     status,
			mdl.UpdateTime.ColumnName().String(): time.Now(),
		},
	)
	return err
}

// CategoryDestory 删除
func (this *CmsLink) CategoryDestory(categoryId int64) error {
	mdl, do := query.CmsLinkCategoryDo()
	if _, err := do.Where(mdl.CategoryID.Eq(categoryId)).Delete(); err != nil {
		return err
	}
	return this.LinkDestroyByCategoryId(categoryId)
}

// LinkDestroyByCategoryId 删除
func (this *CmsLink) LinkDestroyByCategoryId(categoryId int64) error {
	mdl, do := query.CmsLinkDo()
	if _, err := do.Where(mdl.CategoryID.Eq(categoryId)).Delete(); err != nil {
		return err
	}
	return nil
}

// LinkPaginate 分页查询
func (this *CmsLink) LinkPaginate(page, limit int, siteId, channelId, categoryId int64, title, callIndex string, status int32) ([]*model.CmsLink, int64, error) {
	mdl, do := query.CmsLinkDo()
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
	return do.Order(mdl.IsTop.Desc(), mdl.SortID).FindByPage((page-1)*limit, limit)
}

// LinkFind 获取
func (this *CmsLink) LinkFind(linkId int64) (*model.CmsLink, error) {
	mdl, do := query.CmsLinkDo()
	return do.Where(mdl.LinkID.Eq(linkId)).First()
}

// LinkSave 保存或更新
func (this *CmsLink) LinkSave(input *model.CmsLink) error {
	mdl, do := query.CmsLinkDo()
	if input.CallIndex != "" {
		if count, _ := do.Where(mdl.LinkID.Neq(input.LinkID), mdl.CallIndex.Eq(input.CallIndex)).Count(); count > 0 {
			return errors.New("调用别名重复")
		}
	}
	var err error
	input.UpdateTime = time.Now()
	if input.LinkID <= 0 {
		input.CreateTime = time.Now()
		err = do.Create(input)
	} else {
		_, err = do.Where(mdl.LinkID.Eq(input.LinkID)).Updates(map[string]interface{}{
			mdl.CategoryID.ColumnName().String(): input.CategoryID,
			mdl.Title.ColumnName().String():      input.Title,
			mdl.CallIndex.ColumnName().String():  input.CallIndex,
			mdl.LinkURL.ColumnName().String():    input.LinkURL,
			mdl.Target.ColumnName().String():     input.Target,
			mdl.ImgURL.ColumnName().String():     input.ImgURL,
			mdl.Remark.ColumnName().String():     input.Remark,
			mdl.SortID.ColumnName().String():     input.SortID,
			mdl.Status.ColumnName().String():     input.Status,
			mdl.IsLock.ColumnName().String():     input.IsLock,
			mdl.IsTop.ColumnName().String():      input.IsTop,
			mdl.IsRed.ColumnName().String():      input.IsRed,
			mdl.IsHot.ColumnName().String():      input.IsHot,
			mdl.IsSlide.ColumnName().String():    input.IsSlide,
			mdl.UpdateTime.ColumnName().String(): input.UpdateTime,
		})
		if err == nil {
			// NewApiLink().InitCache()
		}
	}
	return err
}

// LinkDestory 删除
func (this *CmsLink) LinkDestory(linkId int64) error {
	mdl, do := query.CmsLinkDo()
	if _, err := do.Where(mdl.LinkID.Eq(linkId)).Delete(); err != nil {
		return err
	}
	// NewApiLink().InitCache()
	return nil
}

func (this *CmsLink) LinkSaveSortId(linkId int64, sortId int32) error {
	mdl, do := query.CmsLinkDo()
	_, err := do.Where(mdl.LinkID.Eq(linkId)).UpdateColumns(
		map[string]interface{}{
			mdl.SortID.ColumnName().String():     sortId,
			mdl.UpdateTime.ColumnName().String(): time.Now(),
		},
	)
	return err
}
