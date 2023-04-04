package biz

import (
	"os"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/app/dal/query"
)

// AlbumPaginate 获取
func (this *CmsArticle) AlbumPaginate(page, limit int, articleId int64) ([]*model.CmsArticleAlbum, int64, error) {
	mdl, do := query.CmsArticleAlbumDo()
	return do.Where(mdl.ArticleID.Eq(articleId)).Order(mdl.SortID).FindByPage((page-1)*limit, limit)
}

// AlbumSave 保存或更新
func (this *CmsArticle) AlbumSave(input *model.CmsArticleAlbum) error {
	mdl, do := query.CmsArticleAlbumDo()
	var err error
	input.UpdateTime = time.Now()
	if input.AlbumID <= 0 {
		input.CreateTime = time.Now()
		err = do.Create(input)
	} else {
		_, err = do.Where(mdl.AlbumID.Eq(input.AlbumID)).Updates(map[string]interface{}{
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
		})
	}
	return err
}

// AlbumDestory 删除
func (this *CmsArticle) AlbumDestory(albumId int64) error {
	mdl, do := query.CmsArticleAlbumDo()
	if finder, err := do.Where(mdl.AlbumID.Eq(albumId)).First(); err != nil {
		return err
	} else {
		if finder != nil && finder.OriginalPath != "" {
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
func (this *CmsArticle) AlbumSaveShow(albumId int64, show int32) error {
	mdl, do := query.CmsArticleAlbumDo()
	_, err := do.Where(mdl.AlbumID.Eq(albumId)).UpdateColumns(
		map[string]interface{}{
			mdl.IsShow.ColumnName().String():     show,
			mdl.UpdateTime.ColumnName().String(): time.Now(),
		},
	)
	return err
}

// AlbumSaveInfo 保存
func (this *CmsArticle) AlbumSaveInfo(albumId int64, title, linkUrl string, click, sortId int32, remark string) error {
	mdl, do := query.CmsArticleAlbumDo()
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
