package biz

import (
	"errors"
	"time"

	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/app/dal/query"
)

type CmsAd struct{}

func NewCmsAd() *CmsAd {
	return &CmsAd{}
}

// CategoryPaginate 分页查询
func (this *CmsAd) CategoryPaginate(page, limit int, siteId, channelId int64, title, callIndex string) ([]*model.CmsAdCategory, int64, error) {
	mdl, do := query.CmsAdCategoryDo()
	if title != "" {
		do = do.Where(mdl.Title.Like("%" + title + "%"))
	}
	if callIndex != "" {
		do = do.Where(mdl.CallIndex.Like("%" + callIndex + "%"))
	}
	return do.Order(mdl.SortID).FindByPage((page-1)*limit, limit)
}

// CategoryFind 获取
func (this *CmsAd) CategoryFind(categoryId int64) (*model.CmsAdCategory, error) {
	mdl, do := query.CmsAdCategoryDo()
	return do.Where(mdl.CategoryID.Eq(categoryId)).First()
}

// CategorySave 保存或更新
func (this *CmsAd) CategorySave(input *model.CmsAdCategory) error {
	mdl, do := query.CmsAdCategoryDo()
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
			mdl.Title.ColumnName().String():      input.Title,
			mdl.CallIndex.ColumnName().String():  input.CallIndex,
			mdl.Remark.ColumnName().String():     input.Remark,
			mdl.SortID.ColumnName().String():     input.SortID,
			mdl.UpdateTime.ColumnName().String(): input.UpdateTime,
		})
	}
	return err
}

func (this *CmsAd) CategorySaveSortId(categoryId int64, sortId int32) error {
	mdl, do := query.CmsAdCategoryDo()
	_, err := do.Where(mdl.CategoryID.Eq(categoryId)).UpdateColumns(
		map[string]interface{}{
			mdl.SortID.ColumnName().String():     sortId,
			mdl.UpdateTime.ColumnName().String(): time.Now(),
		},
	)
	return err
}

// AdClone 克隆
func (this *CmsAd) AdClone(adId int64) (int64, error) {
	mdl, do := query.CmsAdDo()
	art, err := do.Where(mdl.AdID.Eq(adId)).First()
	if err != nil {
		return 0, err
	}
	art.AdID = 0
	art.CreateTime = time.Now()
	art.UpdateTime = time.Now()
	art.Status = 0
	err = do.Create(art)
	return art.AdID, err
}

// AdChangeStatus 修改状态
func (this *CmsAd) AdChangeStatus(adId int64, status int32) error {
	mdl, do := query.CmsAdDo()
	_, err := do.Where(mdl.AdID.Eq(adId)).UpdateColumns(
		map[string]interface{}{
			mdl.Status.ColumnName().String():     status,
			mdl.UpdateTime.ColumnName().String(): time.Now(),
		},
	)
	return err
}

// CategoryDestory 删除
func (this *CmsAd) CategoryDestory(categoryId int64) error {
	mdl, do := query.CmsAdCategoryDo()
	if _, err := do.Where(mdl.CategoryID.Eq(categoryId)).Delete(); err != nil {
		return err
	}
	return this.AdDestroyByCategoryId(categoryId)
}

// AdDestroyByCategoryId 删除
func (this *CmsAd) AdDestroyByCategoryId(categoryId int64) error {
	mdl, do := query.CmsAdDo()
	if _, err := do.Where(mdl.CategoryID.Eq(categoryId)).Delete(); err != nil {
		return err
	}
	return nil
}

// AdPaginate 分页查询
func (this *CmsAd) AdPaginate(page, limit int, siteId, channelId, categoryId int64, title, callIndex string, status int32) ([]*model.CmsAd, int64, error) {
	mdl, do := query.CmsAdDo()
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

// AdFind 获取
func (this *CmsAd) AdFind(adId int64) (*model.CmsAd, error) {
	mdl, do := query.CmsAdDo()
	return do.Where(mdl.AdID.Eq(adId)).First()
}

// AdSave 保存或更新
func (this *CmsAd) AdSave(input *model.CmsAd) error {
	mdl, do := query.CmsAdDo()
	if input.CallIndex != "" {
		if count, _ := do.Where(mdl.AdID.Neq(input.AdID), mdl.CallIndex.Eq(input.CallIndex)).Count(); count > 0 {
			return errors.New("调用别名重复")
		}
	}
	var err error
	input.UpdateTime = time.Now()
	if input.AdID <= 0 {
		input.CreateTime = time.Now()
		err = do.Create(input)
	} else {
		_, err = do.Where(mdl.AdID.Eq(input.AdID)).Updates(map[string]interface{}{
			mdl.CategoryID.ColumnName().String(): input.CategoryID,
			mdl.Title.ColumnName().String():      input.Title,
			mdl.CallIndex.ColumnName().String():  input.CallIndex,
			mdl.LinkURL.ColumnName().String():    input.LinkURL,
			mdl.Target.ColumnName().String():     input.Target,
			mdl.ImgUrl1.ColumnName().String():    input.ImgUrl1,
			mdl.ImgUrl2.ColumnName().String():    input.ImgUrl2,
			mdl.Remark.ColumnName().String():     input.Remark,
			mdl.SortID.ColumnName().String():     input.SortID,
			mdl.Status.ColumnName().String():     input.Status,
			mdl.IsLock.ColumnName().String():     input.IsLock,
			mdl.IsTop.ColumnName().String():      input.IsTop,
			mdl.IsRed.ColumnName().String():      input.IsRed,
			mdl.IsHot.ColumnName().String():      input.IsHot,
			mdl.IsSlide.ColumnName().String():    input.IsSlide,
			mdl.BeginTime.ColumnName().String():  input.BeginTime,
			mdl.EndTime.ColumnName().String():    input.EndTime,
			mdl.UpdateTime.ColumnName().String(): input.UpdateTime,
		})
	}
	return err
}

// AdDestory 删除
func (this *CmsAd) AdDestory(adId int64) error {
	mdl, do := query.CmsAdDo()
	if _, err := do.Where(mdl.AdID.Eq(adId)).Delete(); err != nil {
		return err
	}
	return nil
}

func (this *CmsAd) AdSaveSortId(adId int64, sortId int32) error {
	mdl, do := query.CmsAdDo()
	_, err := do.Where(mdl.AdID.Eq(adId)).UpdateColumns(
		map[string]interface{}{
			mdl.SortID.ColumnName().String():     sortId,
			mdl.UpdateTime.ColumnName().String(): time.Now(),
		},
	)
	return err
}
