package biz

import (
	"errors"
	"time"

	"haedu.gov.cn/tools/xgeneric"

	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/app/dal/query"
	"haedu.gov.cn/cms/global"
)

type CmsTag struct{}

func NewCmsTag() *CmsTag {
	return &CmsTag{}
}

// TagClone 克隆
func (svc *CmsTag) TagClone(tagId int64) (int64, error) {
	mdl, do := query.CmsTagDo()
	art, err := do.Where(mdl.TagID.Eq(tagId)).First()
	if err != nil {
		return 0, err
	}
	art.TagID = 0
	art.CreateTime = time.Now()
	art.UpdateTime = time.Now()
	art.Status = 0
	err = do.Create(art)
	return art.TagID, err
}

// TagChangeStatus 修改状态
func (svc *CmsTag) TagChangeStatus(tagId int64, status int32) error {
	mdl, do := query.CmsTagDo()
	_, err := do.Where(mdl.TagID.Eq(tagId)).UpdateColumns(
		map[string]interface{}{
			mdl.Status.ColumnName().String():     status,
			mdl.UpdateTime.ColumnName().String(): time.Now(),
		},
	)
	return err
}

// TagPaginate 分页查询
func (svc *CmsTag) TagPaginate(page, limit int, channelId int64, name, title string, status int32, siteId ...int64) ([]*model.CmsTag, int64, error) {
	mdl, do := query.CmsTagDo()
	if len(siteId) > 0 {
		do = do.Where(mdl.SiteID.In(siteId...))
	}
	if name != "" {
		do = do.Where(mdl.Name.Like("%" + name + "%"))
	}
	if title != "" {
		do = do.Where(mdl.Title.Like("%" + title + "%"))
	}
	if status >= 0 {
		do = do.Where(mdl.Status.Eq(status))
	}
	return do.Order(mdl.SortID).FindByPage((page-1)*limit, limit)
}

// TagFind 获取
func (svc *CmsTag) TagFind(tagId int64) (*model.CmsTag, error) {
	mdl, do := query.CmsTagDo()
	return do.Where(mdl.TagID.Eq(tagId)).First()
}

// TagSave 保存或更新
func (svc *CmsTag) TagSave(input *model.CmsTag) error {
	mdl, do := query.CmsTagDo()
	if input.Name != "" {
		if count, _ := do.Where(mdl.TagID.Neq(input.TagID), mdl.SiteID.Eq(input.SiteID), mdl.Name.Eq(input.Name)).Count(); count > 0 {
			return errors.New("同一站点下标签名称重复")
		}
	}
	var err error
	input.UpdateTime = time.Now()
	input.ImgUrl2 = xgeneric.IFF(input.ImgUrl1 == "", "", input.ImgUrl2)
	if input.TagID <= 0 {
		input.CreateTime = time.Now()
		err = do.Create(input)
	} else {
		_, err = do.Where(mdl.TagID.Eq(input.TagID)).Updates(map[string]interface{}{
			mdl.SiteID.ColumnName().String():         input.SiteID,
			mdl.Name.ColumnName().String():           input.Name,
			mdl.Title.ColumnName().String():          input.Title,
			mdl.SeoTitle.ColumnName().String():       input.SeoTitle,
			mdl.SeoKeyword.ColumnName().String():     input.SeoKeyword,
			mdl.SeoDescription.ColumnName().String(): input.SeoDescription,
			mdl.ImgUrl1.ColumnName().String():        input.ImgUrl1,
			mdl.ImgUrl2.ColumnName().String():        input.ImgUrl2,
			mdl.Remark.ColumnName().String():         input.Remark,
			mdl.SortID.ColumnName().String():         input.SortID,
			mdl.Status.ColumnName().String():         input.Status,
			mdl.Template.ColumnName().String():       input.Template,
			mdl.UpdateTime.ColumnName().String():     input.UpdateTime,
			mdl.UpdateID.ColumnName().String():       input.UpdateID,
			mdl.UpdateName.ColumnName().String():     input.UpdateName,
		})
	}
	return err
}

/**
 * @description: TagDestory 删除
 * @param {int64} tagId ID
 * @return {*}
 */
func (svc *CmsTag) TagDestory(tagId int64) error {
	mdl, do := query.CmsTagDo()
	if _, err := do.Where(mdl.TagID.Eq(tagId)).Delete(); err != nil {
		return err
	}
	return nil
}

/**
 * @description: TagSaveSortId 保存排序
 * @param {int64} tagId ID
 * @param {int32} sortId 排序
 * @return {*}
 */
func (svc *CmsTag) TagSaveSortId(tagId int64, sortId int32) error {
	mdl, do := query.CmsTagDo()
	_, err := do.Where(mdl.TagID.Eq(tagId)).UpdateColumns(
		map[string]interface{}{
			mdl.SortID.ColumnName().String():     sortId,
			mdl.UpdateTime.ColumnName().String(): time.Now(),
		},
	)
	return err
}

// SiteGet 获取站点
func (svc *CmsTag) SiteGet(roleId int64, roleType string) ([]*model.CmsSite, error) {
	if global.IsSuper(roleType) {
		siteList, _, _ := NewCmsSite().SitePaginate(1, 999999, "", "")
		return siteList, nil
	}
	siteIdList, _, _ := NewCmsAdmin().RoleSiteFind(roleId)
	siteList := []*model.CmsSite{}
	if len(siteIdList) > 0 {
		for i := 0; i < len(siteIdList); i++ {
			item, _ := NewCmsSite().SiteOne(siteIdList[i].SiteID)
			siteList = append(siteList, item)
		}
	}
	return siteList, nil
}
