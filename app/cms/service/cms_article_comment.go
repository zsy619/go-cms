package service

import (
	"time"

	"gorm.io/gen"

	"haedu.gov.cn/cms/app/cms/domain"
	"haedu.gov.cn/cms/app/cms/mapper"
	service_model "haedu.gov.cn/cms/app/cms/service/model"
)

/**
 * @description: 评论分页
 * @param {*} page 页码
 * @param {int} limit 每页数量
 * @param {*} siteId 站点ID
 * @param {int64} channelId 频道ID
 * @param {int32} status 状态
 * @param {int32} lock 锁定
 * @param {int32} reply 回复
 * @return {*}
 */
func (svc *CmsArticle) CommentPaginate(page, limit int, siteId, channelId int64, status, lock, reply int32) ([]*service_model.ApiArticleCommentModel, int64, error) {
	cmtMdl, cmtDo := mapper.CmsArticleCommentDo()
	artMdl, _ := mapper.CmsArticleDo()
	conds := []gen.Condition{}
	if channelId > 0 {
		conds = append(conds, cmtMdl.ChannelID.Eq(channelId))
	}
	if siteId > 0 {
		conds = append(conds, cmtMdl.SiteID.Eq(siteId))
	}
	if status >= 0 {
		conds = append(conds, cmtMdl.Status.Eq(status))
	}
	if lock == 0 {
		conds = append(conds, cmtMdl.IsLock.Is(false))
	}
	if lock == 1 {
		conds = append(conds, cmtMdl.IsLock.Is(true))
	}
	if reply == 0 {
		conds = append(conds, cmtMdl.IsReply.Is(false))
	}
	if reply == 1 {
		conds = append(conds, cmtMdl.IsReply.Is(true))
	}
	out := make([]*service_model.ApiArticleCommentModel, 0)
	count, err := cmtDo.LeftJoin(artMdl, cmtMdl.ArticleID.EqCol(artMdl.ArticleID)).Select(cmtMdl.ALL, artMdl.Title).Where(conds...).Order(cmtMdl.AddTime.Desc()).ScanByPage(&out, (page-1)*limit, limit)
	return out, count, err
}

/**
 * @description: 删除评论
 * @param {int64} commentId 评论ID
 * @return {*}
 */
func (svc *CmsArticle) CommentDestory(commentId int64) error {
	cmtMdl, cmtDo := mapper.CmsArticleCommentDo()
	if _, err := cmtDo.Where(cmtMdl.CommentID.Eq(commentId)).Delete(); err != nil {
		return err
	}
	return nil
}

/**
 * @description: 编辑评论
 * @param {*model.CmsArticleComment} input 评论信息
 * @return {*}
 */
func (svc *CmsArticle) CommentEdit(input *domain.CmsArticleComment) error {
	mdl, do := mapper.CmsArticleCommentDo()
	var err error
	input.ReplyTime = time.Now()
	_, err = do.Where(mdl.CommentID.Eq(input.CommentID)).Updates(map[string]interface{}{
		mdl.IsLock.ColumnName().String():       input.IsLock,
		mdl.IsReply.ColumnName().String():      input.IsReply,
		mdl.ReplyTime.ColumnName().String():    input.ReplyTime,
		mdl.ReplyUser.ColumnName().String():    input.ReplyUser,
		mdl.ReplyContent.ColumnName().String(): input.ReplyContent,
		mdl.Status.ColumnName().String():       input.Status,
	})
	return err
}

/**
 * @description: 获取评论详情
 * @param {int64} commentId 评论ID
 * @return {*}
 */
func (svc *CmsArticle) CommentFind(commentId int64) (*domain.CmsArticleComment, error) {
	mdl, do := mapper.CmsArticleCommentDo()
	return do.Where(mdl.CommentID.Eq(commentId)).First()
}

/**
 * @description: 更改评论状态
 * @param {int64} commentId 评论ID
 * @param {int32} status 状态
 * @return {*}
 */
func (svc *CmsArticle) CommentChangeStatus(commentId int64, status int32) error {
	mdl, do := mapper.CmsArticleCommentDo()
	_, err := do.Where(mdl.CommentID.Eq(commentId)).UpdateColumns(
		map[string]interface{}{
			mdl.Status.ColumnName().String(): status,
		},
	)
	return err
}
