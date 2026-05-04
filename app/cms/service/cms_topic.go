package service

import (
	"errors"
	"time"

	"haedu.gov.cn/tools/xgeneric"

	"haedu.gov.cn/cms/app/cms/domain"
	"haedu.gov.cn/cms/app/cms/mapper"
	"haedu.gov.cn/cms/global"
)

type CmsTopic struct{}

func NewCmsTopic() *CmsTopic {
	return &CmsTopic{}
}

// TopicClone 克隆
func (svc *CmsTopic) TopicClone(topicId int64) (int64, error) {
	mdl, do := mapper.CmsTopicDo()
	art, err := do.Where(mdl.TopicID.Eq(topicId)).First()
	if err != nil {
		return 0, err
	}
	art.TopicID = 0
	art.CreateTime = time.Now()
	art.UpdateTime = time.Now()
	art.Status = 0
	err = do.Create(art)
	return art.TopicID, err
}

// TopicChangeStatus 修改状态
func (svc *CmsTopic) TopicChangeStatus(topicId int64, status int32) error {
	mdl, do := mapper.CmsTopicDo()
	_, err := do.Where(mdl.TopicID.Eq(topicId)).UpdateColumns(
		map[string]interface{}{
			mdl.Status.ColumnName().String():     status,
			mdl.UpdateTime.ColumnName().String(): time.Now(),
		},
	)
	return err
}

// TopicPaginate 分页查询
func (svc *CmsTopic) TopicPaginate(page, limit int, channelId int64, name, title string, status int32, siteId ...int64) ([]*domain.CmsTopic, int64, error) {
	mdl, do := mapper.CmsTopicDo()
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

// TopicFind 获取
func (svc *CmsTopic) TopicFind(topicId int64) (*domain.CmsTopic, error) {
	mdl, do := mapper.CmsTopicDo()
	return do.Where(mdl.TopicID.Eq(topicId)).First()
}

// TopicSave 保存或更新
func (svc *CmsTopic) TopicSave(input *domain.CmsTopic) error {
	mdl, do := mapper.CmsTopicDo()
	if input.Name != "" {
		if count, _ := do.Where(mdl.TopicID.Neq(input.TopicID), mdl.SiteID.Eq(input.SiteID), mdl.Name.Eq(input.Name)).Count(); count > 0 {
			return errors.New("专题名称重复")
		}
	}
	var err error
	input.UpdateTime = time.Now()
	input.ImgUrl2 = xgeneric.IFF(input.ImgUrl1 == "", "", input.ImgUrl2)
	if input.TopicID <= 0 {
		input.CreateTime = time.Now()
		err = do.Create(input)
	} else {
		_, err = do.Where(mdl.TopicID.Eq(input.TopicID)).Updates(map[string]interface{}{
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
 * @description: TopicDestory 删除
 * @param {int64} topicId ID
 * @return {*}
 */
func (svc *CmsTopic) TopicDestory(topicId int64) error {
	mdl, do := mapper.CmsTopicDo()
	if _, err := do.Where(mdl.TopicID.Eq(topicId)).Delete(); err != nil {
		return err
	}
	return nil
}

/**
 * @description: TopicSaveSortId 保存排序
 * @param {int64} topicId ID
 * @param {int32} sortId 排序
 * @return {*}
 */
func (svc *CmsTopic) TopicSaveSortId(topicId int64, sortId int32) error {
	mdl, do := mapper.CmsTopicDo()
	_, err := do.Where(mdl.TopicID.Eq(topicId)).UpdateColumns(
		map[string]interface{}{
			mdl.SortID.ColumnName().String():     sortId,
			mdl.UpdateTime.ColumnName().String(): time.Now(),
		},
	)
	return err
}

// SiteGet 获取站点
func (svc *CmsTopic) SiteGet(roleId int64, roleType string) ([]*domain.CmsSite, error) {
	if global.IsSuper(roleType) {
		siteList, _, _ := NewCmsSite().SitePaginate(1, 999999, "", "")
		return siteList, nil
	}
	siteIdList, _, _ := NewCmsAdmin().RoleSiteFind(roleId)
	siteList := []*domain.CmsSite{}
	if len(siteIdList) > 0 {
		for i := 0; i < len(siteIdList); i++ {
			item, _ := NewCmsSite().SiteOne(siteIdList[i].SiteID)
			siteList = append(siteList, item)
		}
	}
	return siteList, nil
}
