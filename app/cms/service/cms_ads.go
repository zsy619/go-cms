package service

import (
	"errors"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/zsy619/tools/xgeneric"

	"haedu.gov.cn/cms/app/cms/domain"
	"haedu.gov.cn/cms/app/cms/mapper"
	"haedu.gov.cn/cms/global"
)

// CmsAds 广告管理服务
type CmsAds struct{}

// NewCmsAds 创建广告服务实例
func NewCmsAds() *CmsAds {
	return &CmsAds{}
}

// CategoryPaginate 分页查询广告分类
// @param page 页码
// @param limit 每页数量
// @param channelId 频道ID
// @param title 标题(模糊搜索)
// @param callIndex 调用别名
// @param siteId 站点ID列表
// @return []*domain.CmsAdsCategory 广告分类列表,总数,错误信息
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

// CategoryFind 根据ID获取广告分类
// @param categoryId 分类ID
// @return *domain.CmsAdsCategory 分类实体,错误信息
func (svc *CmsAds) CategoryFind(categoryId int64) (*domain.CmsAdsCategory, error) {
	mdl, do := mapper.CmsAdsCategoryDo()
	return do.Where(mdl.CategoryID.Eq(categoryId)).First()
}

// CategorySave 保存或更新广告分类
// @param input 广告分类实体
// @return error 错误信息
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
			logs.Error("更新广告分类失败: adsId=%d, error=%v", input.CategoryID, err)
		}
		adsMdl, adsDo := mapper.CmsAdsDo()
		_, err = adsDo.Where(adsMdl.CategoryID.Eq(input.CategoryID)).UpdateColumns(map[string]interface{}{
			adsMdl.SiteID.ColumnName().String():     input.SiteID,
			adsMdl.UpdateTime.ColumnName().String(): input.UpdateTime,
		})
		if err != nil {
			logs.Error("更新广告分类下的广告站点ID失败: categoryId=%d, error=%v", input.CategoryID, err)
		}
	}
	return err
}

// CategorySaveSortId 保存广告分类排序
// @param categoryId 分类ID
// @param sortId 排序值
// @return error 错误信息
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

// AdsClone 克隆广告
// @param adsId 原广告ID
// @return int64 新广告ID,错误信息
func (svc *CmsAds) AdsClone(adsId int64) (int64, error) {
	mdl, do := mapper.CmsAdsDo()
	art, err := do.Where(mdl.AdsID.Eq(adsId)).First()
	if err != nil {
		logs.Error("克隆广告查询失败: adsId=%d, error=%v", adsId, err)
		return 0, err
	}
	art.AdsID = 0
	art.CreateTime = time.Now()
	art.UpdateTime = time.Now()
	art.Status = 0
	err = do.Create(art)
	if err != nil {
		logs.Error("克隆广告创建失败: adsId=%d, error=%v", adsId, err)
	}
	return art.AdsID, err
}

// AdsChangeStatus 修改广告状态
// @param adsId 广告ID
// @param status 状态值
// @return error 错误信息
func (svc *CmsAds) AdsChangeStatus(adsId int64, status int32) error {
	mdl, do := mapper.CmsAdsDo()
	_, err := do.Where(mdl.AdsID.Eq(adsId)).UpdateColumns(
		map[string]interface{}{
			mdl.Status.ColumnName().String():     status,
			mdl.UpdateTime.ColumnName().String(): time.Now(),
		},
	)
	if err != nil {
		logs.Error("修改广告状态失败: adsId=%d, status=%d, error=%v", adsId, status, err)
	}
	return err
}

// CategoryDestory 删除广告分类
// @param categoryId 分类ID
// @return error 错误信息
func (svc *CmsAds) CategoryDestory(categoryId int64) error {
	mdl, do := mapper.CmsAdsCategoryDo()
	if _, err := do.Where(mdl.CategoryID.Eq(categoryId)).Delete(); err != nil {
		logs.Error("删除广告分类失败: categoryId=%d, error=%v", categoryId, err)
		return err
	}
	return svc.AdsDestroyByCategoryId(categoryId)
}

// AdsDestroyByCategoryId 根据分类ID删除广告
// @param categoryId 分类ID
// @return error 错误信息
func (svc *CmsAds) AdsDestroyByCategoryId(categoryId int64) error {
	mdl, do := mapper.CmsAdsDo()
	if _, err := do.Where(mdl.CategoryID.Eq(categoryId)).Delete(); err != nil {
		logs.Error("删除广告失败: categoryId=%d, error=%v", categoryId, err)
		return err
	}
	return nil
}

// AdsPaginate 分页查询广告
// @param page 页码
// @param limit 每页数量
// @param channelId 频道ID
// @param categoryId 分类ID
// @param title 标题
// @param callIndex 调用别名
// @param status 状态
// @param siteId 站点ID列表
// @return []*domain.CmsAds 广告列表,总数,错误信息
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

// AdsFind 根据ID获取广告
// @param adsId 广告ID
// @return *domain.CmsAds 广告实体,错误信息
func (svc *CmsAds) AdsFind(adsId int64) (*domain.CmsAds, error) {
	mdl, do := mapper.CmsAdsDo()
	return do.Where(mdl.AdsID.Eq(adsId)).First()
}

// AdsSave 保存或更新广告
// @param input 广告实体
// @return error 错误信息
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
			logs.Error("获取站点ID失败: categoryId=%d, error=%v", input.CategoryID, err)
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
		if err != nil {
			logs.Error("更新广告失败: adsId=%d, error=%v", input.AdsID, err)
		}
	}
	return err
}

// AdsDestory 删除广告
// @param adsId 广告ID
// @return error 错误信息
func (svc *CmsAds) AdsDestory(adsId int64) error {
	mdl, do := mapper.CmsAdsDo()
	if _, err := do.Where(mdl.AdsID.Eq(adsId)).Delete(); err != nil {
		logs.Error("删除广告失败: adsId=%d, error=%v", adsId, err)
		return err
	}
	return nil
}

// AdsSaveSortId 保存广告排序
// @param adsId 广告ID
// @param sortId 排序值
// @return error 错误信息
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

// SiteCategoryGet 获取站点与分类列表(根据角色权限)
// @param roleId 角色ID
// @param roleType 角色类型
// @return []*domain.CmsSite 站点列表, []*domain.CmsAdsCategory 分类列表, error
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

// SiteIdsGet 根据站点筛选条件和角色获取站点ID集合
// @param siteId 站点ID(>0时直接使用)
// @param roleType 角色类型
// @param roleId 角色ID
// @return []int64 站点ID列表
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

// FindByDate 根据日期范围查询广告创建数量(用于统计报表)
// @param selectTime 日期列表
// @return []int64 每天的广告数量列表
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