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

type CmsAlbum struct{}

func NewCmsAlbum() *CmsAlbum {
	return &CmsAlbum{}
}

// AlbumPaginate 获取
func (this *CmsAlbum) AlbumPaginate(page, limit int, tableName string, recordId int64, typeId int32, adminId int64, roleType string) ([]*model.CmsAlbum, int64, error) {
	mdl, do := query.CmsAlbumDo()
	if global.IsSuper(roleType) {
		return do.Where(mdl.TableName_.Eq(tableName), mdl.RecordID.Eq(recordId), mdl.TypeID.Eq(typeId)).Order(mdl.SortID).FindByPage((page-1)*limit, limit)
	} else {
		return do.Where(mdl.TableName_.Eq(tableName), mdl.RecordID.Eq(recordId), mdl.TypeID.Eq(typeId), mdl.CreateID.Eq(int32(adminId))).Order(mdl.SortID).FindByPage((page-1)*limit, limit)
	}
}

func (this *CmsAlbum) AlbumSearch(page, limit int, tableName, title, ext string, adminId int64, roleType string) ([]*model.CmsAlbum, int64, error) {
	mdl, do := query.CmsAlbumDo()
	if tableName != "" {
		do = do.Where(mdl.TableName_.Eq(tableName))
	}
	if title != "" {
		do = do.Where(mdl.Title.Like("%" + title + "%"))
	}
	if ext != "" {
		do = do.Where(mdl.FileExt.Like("%" + ext + "%"))
	}
	if global.IsSuper(roleType) {
		return do.Order(mdl.SortID).FindByPage((page-1)*limit, limit)
	} else {
		do = do.Where(mdl.CreateID.Eq(int32(adminId)))
		return do.Order(mdl.SortID).FindByPage((page-1)*limit, limit)
	}
}

// AlbumSave 保存或更新
func (this *CmsAlbum) AlbumSave(input *model.CmsAlbum) error {
	mdl, do := query.CmsAlbumDo()
	var err error
	input.UpdateTime = time.Now()
	if input.AlbumID <= 0 {
		input.CreateTime = time.Now()
		err = do.Create(input)
	} else {
		_, err = do.Where(mdl.AlbumID.Eq(input.AlbumID)).Updates(map[string]interface{}{
			mdl.TableName_.ColumnName().String():   input.TableName_,
			mdl.RecordID.ColumnName().String():     input.RecordID,
			mdl.TypeID.ColumnName().String():       input.TypeID,
			mdl.Title.ColumnName().String():        input.Title,
			mdl.ThumbPath.ColumnName().String():    input.ThumbPath,
			mdl.OriginalPath.ColumnName().String(): input.OriginalPath,
			mdl.FileSize.ColumnName().String():     input.FileSize,
			mdl.FileExt.ColumnName().String():      input.FileExt,
			mdl.LinkURL.ColumnName().String():      input.LinkURL,
			mdl.Click.ColumnName().String():        input.Click,
			mdl.IsShow.ColumnName().String():       input.IsShow,
			mdl.Remark.ColumnName().String():       input.Remark,
			mdl.SortID.ColumnName().String():       input.SortID,
			mdl.UpdateTime.ColumnName().String():   input.UpdateTime,
			mdl.CreateID.ColumnName().String():     input.CreateID,
		})
	}
	return err
}

// AlbumDestory 删除
func (this *CmsAlbum) AlbumDestory(albumId int64) error {
	mdl, do := query.CmsAlbumDo()
	if finder, err := do.Where(mdl.AlbumID.Eq(albumId)).First(); err != nil {
		return err
	} else {
		if finder != nil && finder.OriginalPath != "" && strings.HasPrefix(finder.OriginalPath, "/Uploads/") {
			err := os.Remove(finder.OriginalPath[1:])
			if err != nil {
				logs.Error("AlbumDestory", err.Error())
			}
		}
	}
	_, err := do.Where(mdl.AlbumID.Eq(albumId)).Delete()
	return err
}

// AlbumSaveShow 保存排序
func (this *CmsAlbum) AlbumSaveShow(albumId int64, show int32) error {
	mdl, do := query.CmsAlbumDo()
	_, err := do.Where(mdl.AlbumID.Eq(albumId)).UpdateColumns(
		map[string]interface{}{
			mdl.IsShow.ColumnName().String():     show,
			mdl.UpdateTime.ColumnName().String(): time.Now(),
		},
	)
	return err
}

// AlbumSaveInfo 保存
func (this *CmsAlbum) AlbumSaveInfo(albumId int64, title, linkUrl string, click, sortId int32, remark string) error {
	mdl, do := query.CmsAlbumDo()
	_, err := do.Where(mdl.AlbumID.Eq(albumId)).UpdateColumns(
		map[string]interface{}{
			mdl.Title.ColumnName().String():      title,
			mdl.LinkURL.ColumnName().String():    linkUrl,
			mdl.Click.ColumnName().String():      click,
			mdl.SortID.ColumnName().String():     sortId,
			mdl.Remark.ColumnName().String():     remark,
			mdl.UpdateTime.ColumnName().String(): time.Now(),
		},
	)
	return err
}
