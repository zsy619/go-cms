package service

import (
	"errors"
	"fmt"
	"time"

	"haedu.gov.cn/tools/xgeneric"

	"haedu.gov.cn/cms/app/cms/domain"
	"haedu.gov.cn/cms/app/cms/mapper"
	"haedu.gov.cn/cms/global"
)

type CmsAds struct{}

func NewCmsAds() *CmsAds {
	return &CmsAds{}
}

// CategoryPaginate 分页查询
func (svc *CmsAds) CategoryPaginate(page, limit int, channelId int64, title, callIndex string, siteId ...int64) ([]*domain.CmsAdsCategory, int64, error) {
	mdl, do := mapper.CmsAdsCategoryDo()
	if len(siteId) > 0 {
		do = do.Where(mdl.SiteID.In(siteId...))
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
func (svc *CmsAds) CategoryFind(categoryId int64) (*domain.CmsAdsCategory, error) {
	mdl, do := mapper.CmsAdsCategoryDo()
	return do.Where(mdl.CategoryID.Eq(categoryId)).First()
}

// CategorySave 保存或更新
func (svc *CmsAds) CategorySave(input *domain.CmsAdsCategory) error {
	mdl, do := mapper.CmsAdsCategoryDo()
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
		if err != nil {
			fmt.Println(err.Error())
		}
		adsMdl, adsDo := mapper.CmsAdsDo()
		_, err = adsDo.Where(adsMdl.CategoryID.Eq(input.CategoryID)).UpdateColumns(map[string]interface{}{
			adsMdl.SiteID.ColumnName().String():     input.SiteID,
			adsMdl.UpdateTime.ColumnName().String(): input.UpdateTime,
		})
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	return err
}

func (svc *CmsAds) CategorySaveSortId(categoryId int64, sortId int32) error {
	mdl, do := mapper.CmsAdsCategoryDo()
	_, err := do.Where(mdl.CategoryID.Eq(categoryId)).UpdateColumns(
		map[string]interface{}{
			mdl.SortID.ColumnName().String():     sortId,
			mdl.UpdateTime.ColumnName().String(): time.Now(),
		},
	)
	return err
}

// AdsClone 克隆
func (svc *CmsAds) AdsClone(adsId int64) (int64, error) {
	mdl, do := mapper.CmsAdsDo()
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
func (svc *CmsAds) AdsChangeStatus(adsId int64, status int32) error {
	mdl, do := mapper.CmsAdsDo()
	_, err := do.Where(mdl.AdsID.Eq(adsId)).UpdateColumns(
		map[string]interface{}{
			mdl.Status.ColumnName().String():     status,
			mdl.UpdateTime.ColumnName().String(): time.Now(),
		},
	)
	return err
}

// CategoryDestory 删除
func (svc *CmsAds) CategoryDestory(categoryId int64) error {
	mdl, do := mapper.CmsAdsCategoryDo()
	if _, err := do.Where(mdl.CategoryID.Eq(categoryId)).Delete(); err != nil {
		return err
	}
	return svc.AdsDestroyByCategoryId(categoryId)
}

// AdsDestroyByCategoryId 删除
func (svc *CmsAds) AdsDestroyByCategoryId(categoryId int64) error {
	mdl, do := mapper.CmsAdsDo()
	if _, err := do.Where(mdl.CategoryID.Eq(categoryId)).Delete(); err != nil {
		return err
	}
	return nil
}

// AdsPaginate 分页查询
func (svc *CmsAds) AdsPaginate(page, limit int, channelId, categoryId int64, title, callIndex string, status int32, siteId ...int64) ([]*domain.CmsAds, int64, error) {
	mdl, do := mapper.CmsAdsDo()
	if len(siteId) > 0 {
		do = do.Where(mdl.SiteID.In(siteId...))
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
func (svc *CmsAds) AdsFind(adsId int64) (*domain.CmsAds, error) {
	mdl, do := mapper.CmsAdsDo()
	return do.Where(mdl.AdsID.Eq(adsId)).First()
}

// AdsSave 保存或更新
func (svc *CmsAds) AdsSave(input *domain.CmsAds) error {
	mdl, do := mapper.CmsAdsDo()
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
		catMdl, catDo := mapper.CmsLinkCategoryDo()
		if err := catDo.Where(catMdl.CategoryID.Eq(input.CategoryID)).Pluck(catMdl.SiteID, &siteId); err != nil {
			fmt.Println(err.Error())
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
func (svc *CmsAds) AdsDestory(adsId int64) error {
	mdl, do := mapper.CmsAdsDo()
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
func (svc *CmsAds) AdsSaveSortId(adsId int64, sortId int32) error {
	mdl, do := mapper.CmsAdsDo()
	_, err := do.Where(mdl.AdsID.Eq(adsId)).UpdateColumns(
		map[string]interface{}{
			mdl.SortID.ColumnName().String():     sortId,
			mdl.UpdateTime.ColumnName().String(): time.Now(),
		},
	)
	return err
}

// SiteCategoryGet 获取站点与分类
func (svc *CmsAds) SiteCategoryGet(roleId int64, roleType string) ([]*domain.CmsSite, []*domain.CmsAdsCategory, error) {
	if global.IsSuper(roleType) {
		list, _, _ := svc.CategoryPaginate(1, 99999, -1, "", "")
		siteList, _, _ := NewCmsSite().SitePaginate(1, 999999, "", "")
		return siteList, list, nil
	}
	siteIdList, _, _ := NewCmsAdmin().RoleSiteFind(roleId)
	siteList := []*domain.CmsSite{}
	categoryList := []*domain.CmsAdsCategory{}
	if len(siteIdList) > 0 {
		var siteIds []int64
		for i := 0; i < len(siteIdList); i++ {
			siteIds = append(siteIds, siteIdList[i].SiteID)
			item, _ := NewCmsSite().SiteOne(siteIdList[i].SiteID)
			siteList = append(siteList, item)
		}
		categoryList, _, _ = svc.CategoryPaginate(1, 99999, -1, "", "", siteIds...)
	}
	return siteList, categoryList, nil
}

// SiteIdsGet 根据传入的站点筛选条件、角色类型、角色ID获取站点ID集合
func (svc *CmsAds) SiteIdsGet(siteId int64, roleType string, roleId int64) []int64 {
	var siteIds []int64
	if siteId > 0 {
		siteIds = append(siteIds, siteId)
	} else {
		if !global.IsSuper(roleType) {
			siteIdList, _, _ := NewCmsAdmin().RoleSiteFind(roleId)
			for _, item := range siteIdList {
				siteIds = append(siteIds, item.SiteID)
			}
		}
	}
	return siteIds
}

func (svc *CmsAds) FindByDate(selectTime ...time.Time) []int64 {
	mdl, do := mapper.CmsAdsDo()
	duration, _ := time.ParseDuration("24h")
	var counts []int64
	if len(selectTime) > 0 {
		for i := 0; i < len(selectTime); i++ {
			count, _ := do.Where(mdl.Status.Eq(2), mdl.CreateTime.Between(selectTime[i], selectTime[i].Add(duration))).Count()
			counts = append(counts, count)
		}
	}
	return counts
}
