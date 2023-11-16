package admin

import (
	"strconv"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/tools/xjson"

	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/controllers/admin/vmodel"
)

func (ctrl *ArticleController) Comment() {
	channelId, _ := ctrl.GetInt64("channelId")
	if channelId <= 0 {
		ctrl.Abort("404")
		ctrl.StopRun()
		return
	}
	ctrl.Data["channelId"] = channelId
	// 获取角色权限
	roleMap := ctrl.RolePowerGet("channel_" + strconv.FormatInt(channelId, 10) + "_comment")
	ctrl.Data["roleMap"] = roleMap
	ctrl.display()
}

func (ctrl *ArticleController) CommentPaginate() {
	page, limit := ctrl.GetPagingParameters()
	channelId, _ := ctrl.GetInt64("channelId")
	siteId, _ := ctrl.GetInt64("siteId")
	status, _ := ctrl.GetInt32("status")
	lock, _ := ctrl.GetInt32("lock")
	reply, _ := ctrl.GetInt32("reply")
	list, count, err := biz.NewCmsArticle().CommentPaginate(page, limit, siteId, channelId, status, lock, reply)
	if err != nil {
		logs.Error("CommentPaginate", err.Error())
	}
	ctrl.JSONPageSuccess(list, count)
}

func (ctrl *ArticleController) CommentEdit() {
	if ctrl.IsPost() {
		mdl := model.CmsArticleComment{}
		if err := ctrl.ParseForm(&mdl); err != nil {
			logs.Error("CommentEdit", err.Error())
			ctrl.JSONError(err.Error())
		}
		if err := biz.NewCmsArticle().CommentEdit(&mdl); err != nil {
			logs.Error("CommentEdit", err.Error())
			ctrl.JSONError(err.Error())
			return
		}
		ctrl.JSONSuccess("保存成功", nil)
	}
	channelId, _ := ctrl.GetInt64("channelId")
	if channelId <= 0 {
		ctrl.Abort("404")
		ctrl.StopRun()
		return
	}
	commentId, _ := ctrl.GetInt64("commentId")
	clone, _ := ctrl.GetInt("clone")
	mdl, err := biz.NewCmsArticle().CommentFind(commentId)
	if err != nil {
		mdl = &model.CmsArticleComment{
			ChannelID: channelId,
			ReplyUser: GlobalAdminName,
			ReplyTime: time.Now(),
		}
	}
	// 是否克隆
	if clone == 1 {
		mdl.CommentID = 0
	}
	ctrl.Data["mdl"] = mdl
	// 获取角色权限
	roleMap := ctrl.RolePowerGet("channel_" + strconv.FormatInt(channelId, 10) + "_comment")
	ctrl.Data["roleMap"] = roleMap
	ctrl.display()
}

func (ctrl *ArticleController) CommentChangeStatus() {
	var mdl vmodel.Comment_ChangeStatusModel
	if err := xjson.Unmarshal(ctrl.Ctx.Input.RequestBody, &mdl); err != nil {
		logs.Error("ArticleChangeStatus", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	for _, commentId := range mdl.CommentIds {
		if err := biz.NewCmsArticle().CommentChangeStatus(commentId, mdl.Status); err != nil {
			logs.Error("CommentChangeStatus", err.Error())
			ctrl.JSONError(err.Error())
			return
		}
	}
	ctrl.JSONSuccess("更改状态成功", nil)
}

func (ctrl *ArticleController) CommentDestory() {
	commentId, _ := ctrl.GetInt64("commentId")
	if err := biz.NewCmsArticle().CommentDestory(commentId); err != nil {
		logs.Error("CommentDestory", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("删除成功", nil)
}
