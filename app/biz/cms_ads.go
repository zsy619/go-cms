package biz

import (
	"errors"
	"time"

	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/app/dal/query"
	"haedu.gov.cn/tools/xgeneric"
)

type CmsAds struct{}

func NewCmsAds() *CmsAds {
	return &CmsAds{}
}

// CategoryPaginate 分页查询
func (this *CmsAds) CategoryPaginate(page, limit int, siteId, channelId int64, title, callIndex string) ([]*model.CmsAdsCategory, int64, error) {
	mdl, do := query.CmsAdsCategoryDo()
	if siteId > 0 {
		do = do.Where(mdl.SiteID.Eq(siteId))
	}
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
func (this *CmsAds) CategoryFind(categoryId int64) (*model.CmsAdsCategory, error) {
	mdl, do := query.CmsAdsCategoryDo()
	return do.Where(mdl.CategoryID.Eq(categoryId)).First()
}

// CategorySave 保存或更新
func (this *CmsAds) CategorySave(input *model.CmsAdsCategory) error {
	mdl, do := query.CmsAdsCategoryDo()
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
			mdl.SiteID.ColumnName().String():     input.SiteID,
			mdl.ChannelID.ColumnName().String():  input.ChannelID,
			mdl.Title.ColumnName().String():      input.Title,
			mdl.CallIndex.ColumnName().String():  input.CallIndex,
			mdl.Remark.ColumnName().String():     input.Remark,
			mdl.SortID.ColumnName().String():     input.SortID,
			mdl.Template.ColumnName().String():   input.Template,
			mdl.UpdateTime.ColumnName().String(): input.UpdateTime,
		})
		adsMdl, adsDo := query.CmsAdsDo()
		adsDo.Where(adsMdl.CategoryID.Eq(input.CategoryID)).UpdateColumns(map[string]interface{}{
			adsMdl.SiteID.ColumnName().String():     input.SiteID,
			adsMdl.UpdateTime.ColumnName().String(): input.UpdateTime,
		})
	}
	return err
}

func (this *CmsAds) CategorySaveSortId(categoryId int64, sortId int32) error {
	mdl, do := query.CmsAdsCategoryDo()
	_, err := do.Where(mdl.CategoryID.Eq(categoryId)).UpdateColumns(
		map[string]interface{}{
			mdl.SortID.ColumnName().String():     sortId,
			mdl.UpdateTime.ColumnName().String(): time.Now(),
		},
	)
	return err
}

// AdsClone 克隆
func (this *CmsAds) AdsClone(adsId int64) (int64, error) {
	mdl, do := query.CmsAdsDo()
	art, err := do.Where(mdl.AdsID.Eq(adsId)).First()
	if err != nil {
		return 0, err
	}
	art.AdsID = 0
	art.CreateTime = time.Now()
	art.UpdateTime = time.Now()
	art.Status = 0
	err = do.Create(art)
	return art.AdsID, err
}

// AdsChangeStatus 修改状态
func (this *CmsAds) AdsChangeStatus(adsId int64, status int32) error {
	mdl, do := query.CmsAdsDo()
	_, err := do.Where(mdl.AdsID.Eq(adsId)).UpdateColumns(
		map[string]interface{}{
			mdl.Status.ColumnName().String():     status,
			mdl.UpdateTime.ColumnName().String(): time.Now(),
		},
	)
	return err
}

// CategoryDestory 删除
func (this *CmsAds) CategoryDestory(categoryId int64) error {
	mdl, do := query.CmsAdsCategoryDo()
	if _, err := do.Where(mdl.CategoryID.Eq(categoryId)).Delete(); err != nil {
		return err
	}
	return this.AdsDestroyByCategoryId(categoryId)
}

// AdsDestroyByCategoryId 删除
func (this *CmsAds) AdsDestroyByCategoryId(categoryId int64) error {
	mdl, do := query.CmsAdsDo()
	if _, err := do.Where(mdl.CategoryID.Eq(categoryId)).Delete(); err != nil {
		return err
	}
	return nil
}

// AdsPaginate 分页查询
func (this *CmsAds) AdsPaginate(page, limit int, siteId, channelId, categoryId int64, title, callIndex string, status int32) ([]*model.CmsAds, int64, error) {
	mdl, do := query.CmsAdsDo()
	if siteId > 0 {
		do = do.Where(mdl.SiteID.Eq(siteId))
	}
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
	return do.Order(mdl.IsTop.Desc(), mdl.SortID).FindByPage((page-1)*limit, limit)
}

// AdsFind 获取
func (this *CmsAds) AdsFind(adsId int64) (*model.CmsAds, error) {
	mdl, do := query.CmsAdsDo()
	return do.Where(mdl.AdsID.Eq(adsId)).First()
}

// AdsSave 保存或更新
func (this *CmsAds) AdsSave(input *model.CmsAds) error {
	mdl, do := query.CmsAdsDo()
	if input.CallIndex != "" {
		if count, _ := do.Where(mdl.AdsID.Neq(input.AdsID), mdl.CallIndex.Eq(input.CallIndex)).Count(); count > 0 {
			return errors.New("调用别名重复")
		}
	}
	var err error
	input.UpdateTime = time.Now()
	input.ImgUrl2 = xgeneric.IFF(input.ImgUrl1 == "", "", input.ImgUrl2)
	if input.AdsID <= 0 {
		input.CreateTime = time.Now()
		err = do.Create(input)
	} else {
		siteId := int64(0)
		catMdl, catDo := query.CmsLinkCategoryDo()
		if err := catDo.Where(catMdl.CategoryID.Eq(input.CategoryID)).Pluck(catMdl.SiteID, &siteId); err != nil {
		}
		_, err = do.Where(mdl.AdsID.Eq(input.AdsID)).Updates(map[string]interface{}{
			mdl.SiteID.ColumnName().String():     siteId,
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
			mdl.UpdateID.ColumnName().String():   input.UpdateID,
			mdl.UpdateName.ColumnName().String(): input.UpdateName,
		})
	}
	return err
}

/**
 * @description: AdsDestory 删除
 * @param {int64} adsId 广告ID
 * @return {*}
 */
func (this *CmsAds) AdsDestory(adsId int64) error {
	mdl, do := query.CmsAdsDo()
	if _, err := do.Where(mdl.AdsID.Eq(adsId)).Delete(); err != nil {
		return err
	}
	return nil
}

/**
 * @description: AdsSaveSortId 保存排序
 * @param {int64} adsId 广告ID
 * @param {int32} sortId 排序
 * @return {*}
 */
func (this *CmsAds) AdsSaveSortId(adsId int64, sortId int32) error {
	mdl, do := query.CmsAdsDo()
	_, err := do.Where(mdl.AdsID.Eq(adsId)).UpdateColumns(
		map[string]interface{}{
			mdl.SortID.ColumnName().String():     sortId,
			mdl.UpdateTime.ColumnName().String(): time.Now(),
		},
	)
	return err
}

// SiteCategoryGet 获取站点与分类
func (this *CmsAds) SiteCategoryGet(roleId int64, roleType string) ([]*model.CmsSite, []*model.CmsAdsCategory, error) {
	if roleType == "super" {
		list, _, _ := this.CategoryPaginate(1, 99999, -1, -1, "", "")
		siteList, _, _ := NewCmsSite().SitePaginate(1, 999999, "", "")
		return siteList, list, nil
	}
	siteIdList, _, _ := NewCmsAdmin().RoleSiteFind(roleId)
	siteList := []*model.CmsSite{}
	categoryList := []*model.CmsAdsCategory{}
	if len(siteIdList) > 0 {
		for i := 0; i < len(siteIdList); i++ {
			list, _, _ := this.CategoryPaginate(1, 99999, siteIdList[i].SiteID, -1, "", "")
			if len(list) > 0 {
				for j := 0; j < len(list); j++ {
					categoryList = append(categoryList, list[j])
				}
			}
			item, _ := NewCmsSite().SiteOne(siteIdList[i].SiteID)
			siteList = append(siteList, item)
		}
	}
	return siteList, categoryList, nil
}
