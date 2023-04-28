package admin

import (
	"time"

	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/controllers/admin/vmodel"
	"haedu.gov.cn/tools/xjson"
)

func (c *ArticleController) Comment() {
	channelId, _ := c.GetInt64("channelId")
	if channelId <= 0 {
		c.Abort("404")
		c.StopRun()
		return
	}
	c.Data["channelId"] = channelId
	c.display()
}

func (c *ArticleController) CommentPaginate() {
	page, limit := c.GetPagingParameters()
	channelId, _ := c.GetInt64("channelId")
	siteId, _ := c.GetInt64("siteId")
	status, _ := c.GetInt32("status")
	lock, _ := c.GetInt32("lock")
	reply, _ := c.GetInt32("reply")
	list, count, err := biz.NewCmsArticle().CommentPaginate(page, limit, siteId, channelId, status, lock, reply)
	if err != nil {
		logs.Error("CommentPaginate", err.Error())
	}
	c.JSONPageSuccess(list, count)
}

func (c *ArticleController) CommentEdit() {
	if c.IsPost() {
		mdl := model.CmsArticleComment{}
		if err := c.ParseForm(&mdl); err != nil {
			logs.Error("CommentEdit", err.Error())
			c.JSONError(err.Error())
		}
		if err := biz.NewCmsArticle().CommentEdit(&mdl); err != nil {
			logs.Error("CommentEdit", err.Error())
			c.JSONError(err.Error())
			return
		}
		c.JSONSuccess("保存成功", nil)
	}
	channelId, _ := c.GetInt64("channelId")
	if channelId <= 0 {
		c.Abort("404")
		c.StopRun()
		return
	}
	commentId, _ := c.GetInt64("commentId")
	clone, _ := c.GetInt("clone")
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
	c.Data["mdl"] = mdl
	c.display()
}

func (c *ArticleController) CommentChangeStatus() {
	var mdl vmodel.Comment_ChangeStatusModel
	if err := xjson.Unmarshal(c.Ctx.Input.RequestBody, &mdl); err != nil {
		logs.Error("ArticleChangeStatus", err.Error())
		c.JSONError(err.Error())
		return
	}
	for _, commentId := range mdl.CommentIds {
		if err := biz.NewCmsArticle().CommentChangeStatus(commentId, mdl.Status); err != nil {
			logs.Error("CommentChangeStatus", err.Error())
			c.JSONError(err.Error())
			return
		}
	}
	c.JSONSuccess("更改状态成功", nil)
}

func (c *ArticleController) CommentDestory() {
	commentId, _ := c.GetInt64("commentId")
	if err := biz.NewCmsArticle().CommentDestory(commentId); err != nil {
		logs.Error("CommentDestory", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("删除成功", nil)
}
