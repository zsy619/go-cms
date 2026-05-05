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

// CmsLink 链接管理服务
type CmsLink struct{}

// NewCmsLink 创建链接服务实例
func NewCmsLink() *CmsLink {
	return &CmsLink{}
}

// CategoryPaginate 分页查询链接分类
// @param page 页码
// @param limit 每页数量
// @param channelId 频道ID
// @param title 标题(模糊搜索)
// @param callIndex 调用别名
// @param siteId 站点ID列表
// @return []*domain.CmsLinkCategory 分类列表,总数,错误信息
func (svc *CmsLink) CategoryPaginate(page, limit int, channelId int64, title, callIndex string, siteId ...int64) ([]*domain.CmsLinkCategory, int64, error) {
	mdl, do := mapper.CmsLinkCategoryDo()
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

// CategoryFind 根据ID获取链接分类
// @param categoryId 分类ID
// @return *domain.CmsLinkCategory 分类实体,错误信息
func (svc *CmsLink) CategoryFind(categoryId int64) (*domain.CmsLinkCategory, error) {
	mdl, do := mapper.CmsLinkCategoryDo()
	return do.Where(mdl.CategoryID.Eq(categoryId)).First()
}

// CategorySave 保存或更新链接分类
// @param input 链接分类实体
// @return error 错误信息
func (svc *CmsLink) CategorySave(input *domain.CmsLinkCategory) error {
	mdl, do := mapper.CmsLinkCategoryDo()
	if input.CallIndex != "" {
		if count, _ := do.Where(mdl.CategoryID.Neq(input.CategoryID), mdl.CallIndex.Eq(input.CallIndex)).Count(); count > 0 {
			return errors.New("调用别名重复")
		}
	}
	var err error
	input.UpdateTime = time.Now()
	input.ImgUrl2 = xgeneric.IFF(input.ImgUrl1 == "", "", input.ImgUrl2)
	if input.CategoryID <= 0 {
		input.CreateTime = time.Now()
		err = do.Create(input)
	} else {
		_, err = do.Where(mdl.CategoryID.Eq(input.CategoryID)).Updates(map[string]interface{}{
			mdl.SiteID.ColumnName().String():           input.SiteID,
			mdl.ChannelID.ColumnName().String():        input.ChannelID,
			mdl.Title.ColumnName().String():            input.Title,
			mdl.CallIndex.ColumnName().String():        input.CallIndex,
			mdl.LinkURL.ColumnName().String():          input.LinkURL,
			mdl.ImgUrl1.ColumnName().String():         input.ImgUrl1,
			mdl.ImgUrl2.ColumnName().String():         input.ImgUrl2,
			mdl.SeoTitle.ColumnName().String():         input.SeoTitle,
			mdl.SeoKeyword.ColumnName().String():       input.SeoKeyword,
			mdl.SeoDescription.ColumnName().String():   input.SeoDescription,
			mdl.Content.ColumnName().String():          input.Content,
			mdl.SortID.ColumnName().String():          input.SortID,
			mdl.Template.ColumnName().String():         input.Template,
			mdl.UpdateTime.ColumnName().String():      input.UpdateTime,
		})
		if err != nil {
			logs.Error("更新链接分类失败: categoryId=%d, error=%v", input.CategoryID, err)
		}
		linkMdl, linkDo := mapper.CmsLinkDo()
		_, err = linkDo.Where(linkMdl.CategoryID.Eq(input.CategoryID)).UpdateColumns(map[string]interface{}{
			linkMdl.SiteID.ColumnName().String():     input.SiteID,
			linkMdl.UpdateTime.ColumnName().String(): input.UpdateTime,
		})
		if err != nil {
			logs.Error("更新链接分类下的链接站点ID失败: categoryId=%d, error=%v", input.CategoryID, err)
		}
	}
	return err
}

// CategorySaveSortId 保存链接分类排序
// @param categoryId 分类ID
// @param sortId 排序值
// @return error 错误信息
func (svc *CmsLink) CategorySaveSortId(categoryId int64, sortId int32) error {
	mdl, do := mapper.CmsLinkCategoryDo()
	_, err := do.Where(mdl.CategoryID.Eq(categoryId)).UpdateColumns(
		map[string]interface{}{
			mdl.SortID.ColumnName().String():     sortId,
			mdl.UpdateTime.ColumnName().String(): time.Now(),
		},
	)
	return err
}

// LinkClone 克隆链接
// @param linkId 原链接ID
// @return int64 新链接ID,错误信息
func (svc *CmsLink) LinkClone(linkId int64) (int64, error) {
	mdl, do := mapper.CmsLinkDo()
	art, err := do.Where(mdl.LinkID.Eq(linkId)).First()
	if err != nil {
		logs.Error("克隆链接查询失败: linkId=%d, error=%v", linkId, err)
		return 0, err
	}
	art.LinkID = 0
	art.CreateTime = time.Now()
	art.UpdateTime = time.Now()
	art.Status = 0
	err = do.Create(art)
	if err != nil {
		logs.Error("克隆链接创建失败: linkId=%d, error=%v", linkId, err)
	}
	return art.LinkID, err
}

// LinkChangeStatus 修改链接状态
// @param linkId 链接ID
// @param status 状态值
// @return error 错误信息
func (svc *CmsLink) LinkChangeStatus(linkId int64, status int32) error {
	mdl, do := mapper.CmsLinkDo()
	_, err := do.Where(mdl.LinkID.Eq(linkId)).UpdateColumns(
		map[string]interface{}{
			mdl.Status.ColumnName().String():     status,
			mdl.UpdateTime.ColumnName().String(): time.Now(),
		},
	)
	if err != nil {
		logs.Error("修改链接状态失败: linkId=%d, status=%d, error=%v", linkId, status, err)
	}
	return err
}

// CategoryDestory 删除链接分类
// @param categoryId 分类ID
// @return error 错误信息
func (svc *CmsLink) CategoryDestory(categoryId int64) error {
	mdl, do := mapper.CmsLinkCategoryDo()
	if _, err := do.Where(mdl.CategoryID.Eq(categoryId)).Delete(); err != nil {
		logs.Error("删除链接分类失败: categoryId=%d, error=%v", categoryId, err)
		return err
	}
	return svc.LinkDestroyByCategoryId(categoryId)
}

// LinkDestroyByCategoryId 根据分类ID删除链接
// @param categoryId 分类ID
// @return error 错误信息
func (svc *CmsLink) LinkDestroyByCategoryId(categoryId int64) error {
	mdl, do := mapper.CmsLinkDo()
	if _, err := do.Where(mdl.CategoryID.Eq(categoryId)).Delete(); err != nil {
		logs.Error("删除链接失败: categoryId=%d, error=%v", categoryId, err)
		return err
	}
	return nil
}

// LinkPaginate 分页查询链接
// @param page 页码
// @param limit 每页数量
// @param channelId 频道ID
// @param categoryId 分类ID
// @param title 标题
// @param callIndex 调用别名
// @param status 状态
// @param siteId 站点ID列表
// @return []*domain.CmsLink 链接列表,总数,错误信息
func (svc *CmsLink) LinkPaginate(page, limit int, channelId, categoryId int64, title, callIndex string, status int32, siteId ...int64) ([]*domain.CmsLink, int64, error) {
	mdl, do := mapper.CmsLinkDo()
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

// LinkFind 根据ID获取链接
// @param linkId 链接ID
// @return *domain.CmsLink 链接实体,错误信息
func (svc *CmsLink) LinkFind(linkId int64) (*domain.CmsLink, error) {
	mdl, do := mapper.CmsLinkDo()
	return do.Where(mdl.LinkID.Eq(linkId)).First()
}

// LinkSave 保存或更新链接
// @param input 链接实体
// @return error 错误信息
func (svc *CmsLink) LinkSave(input *domain.CmsLink) error {
	mdl, do := mapper.CmsLinkDo()
	if input.CallIndex != "" {
		if count, _ := do.Where(mdl.LinkID.Neq(input.LinkID), mdl.CallIndex.Eq(input.CallIndex)).Count(); count > 0 {
			return errors.New("调用别名重复")
		}
	}
	var err error
	input.UpdateTime = time.Now()
	input.ImgUrl2 = xgeneric.IFF(input.ImgUrl1 == "", "", input.ImgUrl2)
	if input.LinkID <= 0 {
		input.CreateTime = time.Now()
		err = do.Create(input)
	} else {
		siteId := int64(0)
		catMdl, catDo := mapper.CmsLinkCategoryDo()
		if err := catDo.Where(catMdl.CategoryID.Eq(input.CategoryID)).Pluck(catMdl.SiteID, &siteId); err != nil {
			logs.Error("获取链接站点ID失败: categoryId=%d, error=%v", input.CategoryID, err)
		}
		_, err = do.Where(mdl.LinkID.Eq(input.LinkID)).Updates(map[string]interface{}{
			mdl.SiteID.ColumnName().String():       siteId,
			mdl.CategoryID.ColumnName().String():  input.CategoryID,
			mdl.Title.ColumnName().String():       input.Title,
			mdl.CallIndex.ColumnName().String():   input.CallIndex,
			mdl.LinkURL.ColumnName().String():     input.LinkURL,
			mdl.Target.ColumnName().String():      input.Target,
			mdl.ImgUrl1.ColumnName().String():     input.ImgUrl1,
			mdl.ImgUrl2.ColumnName().String():     input.ImgUrl2,
			mdl.Remark.ColumnName().String():      input.Remark,
			mdl.SortID.ColumnName().String():      input.SortID,
			mdl.Status.ColumnName().String():      input.Status,
			mdl.IsLock.ColumnName().String():      input.IsLock,
			mdl.IsTop.ColumnName().String():       input.IsTop,
			mdl.IsRed.ColumnName().String():       input.IsRed,
			mdl.IsHot.ColumnName().String():       input.IsHot,
			mdl.IsSlide.ColumnName().String():     input.IsSlide,
			mdl.UpdateTime.ColumnName().String():  input.UpdateTime,
			mdl.UpdateID.ColumnName().String():    input.UpdateID,
			mdl.UpdateName.ColumnName().String():  input.UpdateName,
		})
		if err != nil {
			logs.Error("更新链接失败: linkId=%d, error=%v", input.LinkID, err)
		}
	}
	return err
}

// LinkDestory 删除链接
// @param linkId 链接ID
// @return error 错误信息
func (svc *CmsLink) LinkDestory(linkId int64) error {
	mdl, do := mapper.CmsLinkDo()
	if _, err := do.Where(mdl.LinkID.Eq(linkId)).Delete(); err != nil {
		logs.Error("删除链接失败: linkId=%d, error=%v", linkId, err)
		return err
	}
	return nil
}

// LinkSaveSortId 保存链接排序
// @param linkId 链接ID
// @param sortId 排序值
// @return error 错误信息
func (svc *CmsLink) LinkSaveSortId(linkId int64, sortId int32) error {
	mdl, do := mapper.CmsLinkDo()
	_, err := do.Where(mdl.LinkID.Eq(linkId)).UpdateColumns(
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
// @return []*domain.CmsSite 站点列表, []*domain.CmsLinkCategory 分类列表, error
func (svc *CmsLink) SiteCategoryGet(roleId int64, roleType string) ([]*domain.CmsSite, []*domain.CmsLinkCategory, error) {
	if global.IsSuper(roleType) {
		list, _, _ := svc.CategoryPaginate(1, 99999, -1, "", "")
		siteList, _, _ := NewCmsSite().SitePaginate(1, 999999, "", "")
		return siteList, list, nil
	}
	siteIdList, _, _ := NewCmsAdmin().RoleSiteFind(roleId)
	siteList := []*domain.CmsSite{}
	categoryList := []*domain.CmsLinkCategory{}
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
func (svc *CmsLink) SiteIdsGet(siteId int64, roleType string, roleId int64) []int64 {
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

// FindByDate 根据日期范围查询链接创建数量(用于统计报表)
// @param selectTime 日期列表
// @return []int64 每天的链接数量列表
func (svc *CmsLink) FindByDate(selectTime ...time.Time) []int64 {
	mdl, do := mapper.CmsLinkDo()
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