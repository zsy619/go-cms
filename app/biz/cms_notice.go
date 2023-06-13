package biz

import (
	"errors"
	"fmt"
	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/app/dal/query"
	"haedu.gov.cn/cms/global"
	"time"
)

type CmsNotice struct {
}

func NewCmsNotice() *CmsNotice {
	return &CmsNotice{}
}

// NoticePaginate 获取系统公告列表
func (this *CmsNotice) NoticePaginate(page, limit int, title string, status int32, adminId int64, roleType string) ([]*model.AdminNotice, int64, error) {
	mdl, do := query.AdminNoticeDo()
	if title != "" {
		do = do.Where(mdl.Title.Like("%" + title + "%"))
	}
	if status >= 0 {
		do = do.Where(mdl.Status.Eq(status))
	}
	if global.IsSuper(roleType) {
		do = do.Order(mdl.SortID).Order(mdl.CreateTime.Desc())
	} else {
		do = do.Where(mdl.CreateID.Eq(int32(adminId))).Order(mdl.SortID).Order(mdl.CreateTime.Desc())
	}
	return do.FindByPage((page-1)*limit, limit)
}

// NoticeShow 首页展示系统公告列表
func (this *CmsNotice) NoticeShow() ([]*model.AdminNotice, error) {
	mdl, do := query.AdminNoticeDo()
	return do.Where(mdl.Status.Eq(2)).Order(mdl.IsTop.Desc(), mdl.CreateTime.Desc()).Limit(6).Find()
}

// NoticeFind 通过notice_id获取详情
func (this *CmsNotice) NoticeFind(noticeId int64) (*model.AdminNotice, error) {
	mdl, do := query.AdminNoticeDo()
	return do.Where(mdl.NoticeID.Eq(noticeId)).First()
}

// NoticeDestroy 根据notice_id删除
func (this *CmsNotice) NoticeDestroy(noticeId int64) error {
	mdl, do := query.AdminNoticeDo()
	if _, err := do.Where(mdl.NoticeID.Eq(noticeId)).Delete(); err != nil {
		return err
	}
	return nil
}

// NoticeSave 保存
func (this *CmsNotice) NoticeSave(input *model.AdminNotice) error {
	if input.Title == "" {
		return errors.New("标题不能为空")
	}
	if input.Content == "" {
		return errors.New("内容不能为空")
	}
	mdl, do := query.AdminNoticeDo()
	var err error
	input.UpdateTime = time.Now()
	if input.NoticeID <= 0 {
		input.CreateTime = time.Now()
		err = do.Create(input)
	} else {
		_, err = do.Where(mdl.NoticeID.Eq(input.NoticeID)).Updates(map[string]interface{}{
			mdl.Title.ColumnName().String():       input.Title,
			mdl.SubTitle.ColumnName().String():    input.SubTitle,
			mdl.Content.ColumnName().String():     input.Content,
			mdl.SortID.ColumnName().String():      input.SortID,
			mdl.IsTop.ColumnName().String():       input.IsTop,
			mdl.Status.ColumnName().String():      input.Status,
			mdl.PublishTime.ColumnName().String(): input.PublishTime,
			mdl.UpdateID.ColumnName().String():    input.UpdateID,
			mdl.UpdateName.ColumnName().String():  input.UpdateName,
			mdl.UpdateTime.ColumnName().String():  input.UpdateTime,
		})
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	return err
}

// NoticeSaveSortId 更新排序
func (this *CmsNotice) NoticeSaveSortId(noticeId int64, updateId int32, updateName string, sortId int32) error {
	mdl, do := query.AdminNoticeDo()
	_, err := do.Where(mdl.NoticeID.Eq(noticeId)).UpdateColumns(
		map[string]interface{}{
			mdl.SortID.ColumnName().String():     sortId,
			mdl.UpdateTime.ColumnName().String(): time.Now(),
			mdl.UpdateID.ColumnName().String():   updateId,
			mdl.UpdateName.ColumnName().String(): updateName,
		},
	)
	return err
}

// NoticeChangeStatus 更新审核状态
func (this *CmsNotice) NoticeChangeStatus(noticeId int64, updateId int32, updateName string, status int32) error {
	mdl, do := query.AdminNoticeDo()
	_, err := do.Where(mdl.NoticeID.Eq(noticeId)).UpdateColumns(
		map[string]interface{}{
			mdl.Status.ColumnName().String():     status,
			mdl.UpdateTime.ColumnName().String(): time.Now(),
			mdl.UpdateID.ColumnName().String():   updateId,
			mdl.UpdateName.ColumnName().String(): updateName,
		},
	)
	return err
}
