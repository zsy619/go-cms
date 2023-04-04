package biz

import (
	"time"

	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/app/dal/query"
)

// AttachPaginate 获取
func (this *CmsArticle) AttachPaginate(page, limit int, articleId int64) ([]*model.CmsArticleAttach, int64, error) {
	mdl, do := query.CmsArticleAttachDo()
	return do.Where(mdl.ArticleID.Eq(articleId)).Order(mdl.SortID).FindByPage((page-1)*limit, limit)
}

// AttachSave 保存或更新
func (this *CmsArticle) AttachSave(input *model.CmsArticleAttach) error {
	mdl, do := query.CmsArticleAttachDo()
	var err error
	input.UpdateTime = time.Now()
	if input.AttachID <= 0 {
		input.CreateTime = time.Now()
		err = do.Create(input)
	} else {
		_, err = do.Where(mdl.AttachID.Eq(input.AttachID)).Updates(map[string]interface{}{
			mdl.IsShow.ColumnName().String():     input.IsShow,
			mdl.Remark.ColumnName().String():     input.Remark,
			mdl.SortID.ColumnName().String():     input.SortID,
			mdl.UpdateTime.ColumnName().String(): input.UpdateTime,
		})
	}
	return err
}

// AttachDestory 删除
func (this *CmsArticle) AttachDestory(attachId int64) error {
	mdl, do := query.CmsArticleAttachDo()
	_, err := do.Where(mdl.AttachID.Eq(attachId)).Delete()
	return err
}

// AttachSaveShow 保存排序
func (this *CmsArticle) AttachSaveShow(attachId int64, show int32) error {
	mdl, do := query.CmsArticleAttachDo()
	_, err := do.Where(mdl.AttachID.Eq(attachId)).UpdateColumns(
		map[string]interface{}{
			mdl.IsShow.ColumnName().String():     show,
			mdl.UpdateTime.ColumnName().String(): time.Now(),
		},
	)
	return err
}

// AttachSaveInfo 保存
func (this *CmsArticle) AttachSaveInfo(attachId int64, title string, point, click, sortId int32, remark string) error {
	mdl, do := query.CmsArticleAttachDo()
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
