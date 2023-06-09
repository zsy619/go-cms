package biz

import (
	"haedu.gov.cn/cms/global"
	"os"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/app/dal/query"
)

type CmsAttach struct{}

func NewCmsAttach() *CmsAttach {
	return &CmsAttach{}
}

// AttachPaginate 获取
func (this *CmsAttach) AttachPaginate(page, limit int, tableName string, recordId int64, typeId int32, adminId int64, roleType string) ([]*model.CmsAttach, int64, error) {
	mdl, do := query.CmsAttachDo()
	if global.IsSuper(roleType) {
		return do.Where(mdl.TableName_.Eq(tableName), mdl.RecordID.Eq(recordId), mdl.TypeID.Eq(typeId)).Order(mdl.SortID).FindByPage((page-1)*limit, limit)
	} else {
		return do.Where(mdl.TableName_.Eq(tableName), mdl.RecordID.Eq(recordId), mdl.TypeID.Eq(typeId), mdl.CreateID.Eq(int32(adminId))).Order(mdl.SortID).FindByPage((page-1)*limit, limit)
	}
}

// AttachSave 保存或更新
func (this *CmsAttach) AttachSave(input *model.CmsAttach) error {
	mdl, do := query.CmsAttachDo()
	var err error
	input.UpdateTime = time.Now()
	if input.AttachID <= 0 {
		input.CreateTime = time.Now()
		err = do.Create(input)
	} else {
		_, err = do.Where(mdl.AttachID.Eq(input.AttachID)).Updates(map[string]interface{}{
			mdl.TableName_.ColumnName().String(): input.TableName_,
			mdl.RecordID.ColumnName().String():   input.RecordID,
			mdl.TypeID.ColumnName().String():     input.TypeID,
			mdl.IsShow.ColumnName().String():     input.IsShow,
			mdl.Remark.ColumnName().String():     input.Remark,
			mdl.SortID.ColumnName().String():     input.SortID,
			mdl.CreateID.ColumnName().String():   input.CreateID,
			mdl.UpdateTime.ColumnName().String(): input.UpdateTime,
		})
	}
	return err
}

// AttachDestory 删除
func (this *CmsAttach) AttachDestory(attachId int64) error {
	mdl, do := query.CmsAttachDo()
	if finder, err := do.Where(mdl.AttachID.Eq(attachId)).First(); err != nil {
		return err
	} else {
		if finder != nil && finder.OriginalPath != "" && strings.HasPrefix(finder.OriginalPath, "/Uploads/") {
			err := os.Remove(finder.OriginalPath[1:])
			if err != nil {
				logs.Error("AttachDestory", err.Error())
			}
		}
	}
	_, err := do.Where(mdl.AttachID.Eq(attachId)).Delete()
	return err
}

// AttachSaveShow 保存排序
func (this *CmsAttach) AttachSaveShow(attachId int64, show int32) error {
	mdl, do := query.CmsAttachDo()
	_, err := do.Where(mdl.AttachID.Eq(attachId)).UpdateColumns(
		map[string]interface{}{
			mdl.IsShow.ColumnName().String():     show,
			mdl.UpdateTime.ColumnName().String(): time.Now(),
		},
	)
	return err
}

// AttachSaveInfo 保存
func (this *CmsAttach) AttachSaveInfo(attachId int64, title string, point, click, sortId int32, remark string) error {
	mdl, do := query.CmsAttachDo()
	_, err := do.Where(mdl.AttachID.Eq(attachId)).UpdateColumns(
		map[string]interface{}{
			mdl.Title.ColumnName().String():      title,
			mdl.Point.ColumnName().String():      point,
			mdl.Click.ColumnName().String():      click,
			mdl.SortID.ColumnName().String():     sortId,
			mdl.Remark.ColumnName().String():     remark,
			mdl.UpdateTime.ColumnName().String(): time.Now(),
		},
	)
	return err
}
