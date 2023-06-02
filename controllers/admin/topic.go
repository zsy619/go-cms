package admin

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/app/lib"
	"haedu.gov.cn/cms/controllers/admin/vmodel"
	"haedu.gov.cn/tools/xjson"
)

type TopicController struct{ BaseController }

// Index 标签管理
// @router /admin/Topic/index [get]
func (c *TopicController) Index() {
	siteList, _ := biz.NewCmsTopic().SiteGet(GlobalRoleId, GlobalRoleType)
	c.Data["site"] = siteList
	c.display()
}

// TopicEdit 编辑
// @router /admin/Topic/TopicEdit [get]
func (c *TopicController) TopicEdit() {
	topicId, _ := c.GetInt64("topicId")
	clone, _ := c.GetInt("clone")
	mdl, err := biz.NewCmsTopic().TopicFind(topicId)
	if err != nil {
		mdl = &model.CmsTopic{
			SortID: 99,
		}
	}
	// 是否克隆
	if clone == 1 {
		mdl.TopicID = 0
	}
	c.Data["mdl"] = mdl
	siteList, _ := biz.NewCmsTopic().SiteGet(GlobalRoleId, GlobalRoleType)
	c.Data["site"] = siteList
	c.display()
}

// TopicSave 保存
// @router /admin/Topic/TopicSave [post]
func (c *TopicController) TopicSave() {
	mdl := model.CmsTopic{}
	if err := c.ParseForm(&mdl); err != nil {
		logs.Error("TopicSave", err.Error())
		c.JSONError(err.Error())
	}
	if mdl.TopicID <= 0 {
		mdl.CreateID = int32(GlobalAdminId)
		mdl.CreateName = GlobalAdminName
	} else {
		mdl.UpdateID = int32(GlobalAdminId)
		mdl.UpdateName = GlobalAdminName
	}
	if err := biz.NewCmsTopic().TopicSave(&mdl); err != nil {
		logs.Error("TopicSave", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("保存成功", nil)
}

// TopicSaveSortId 保存排序
// @router /admin/Topic/TopicSaveSortId [post]
func (c *TopicController) TopicSaveSortId() {
	mdls := []vmodel.Topic_SaveSortIdModel{}
	data := c.Ctx.Input.RequestBody
	fmt.Println("TopicSaveSortId", string(data))
	if err := xjson.Unmarshal(c.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("TopicSaveSortId", err.Error())
		c.JSONError(err.Error())
	}
	service := biz.NewCmsTopic()
	for _, mdl := range mdls {
		if err := service.TopicSaveSortId(mdl.TopicId, int32(mdl.SortId)); err != nil {
			logs.Error("TopicSaveSortId", err.Error())
			c.JSONError(err.Error())
			return
		}
	}
	c.JSONSuccess("保存成功", nil)
}

// TopicDestory 删除
// @router /admin/Topic/TopicDestory [post]
func (c *TopicController) TopicDestory() {
	topicId, _ := c.GetInt64("topicId")
	if err := biz.NewCmsTopic().TopicDestory(topicId); err != nil {
		logs.Error("TopicDestory", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("删除成功", nil)
}

// TopicChangeStatus 更改状态
// @router /admin/Topic/TopicChangeStatus [post]
func (c *TopicController) TopicChangeStatus() {
	var mdl vmodel.Topic_ChangeStatusModel
	if err := xjson.Unmarshal(c.Ctx.Input.RequestBody, &mdl); err != nil {
		logs.Error("TopicChangeStatus", err.Error())
		c.JSONError(err.Error())
		return
	}
	for _, topicId := range mdl.TopicIds {
		if err := biz.NewCmsTopic().TopicChangeStatus(topicId, mdl.Status); err != nil {
			logs.Error("TopicChangeStatus", err.Error())
			c.JSONError(err.Error())
			return
		}
	}
	c.JSONSuccess("更改状态成功", nil)
}

// TopicPaginate 列表
// @router /admin/Topic/TopicPaginate [get]
func (c *TopicController) TopicPaginate() {
	page, limit := c.GetPagingParameters()
	siteId, _ := c.GetInt64("siteId")
	status, _ := c.GetInt32("status")
	title := c.GetString("title")
	name := c.GetString("name")
	if GlobalRoleType == "super" || siteId != 0 {
		list, count, _ := biz.NewCmsTopic().TopicPaginate(page, limit, siteId, -1, name, title, status)
		c.JSONPage(lib.CodeSuccess, "", list, count)
	} else {
		siteIdList, _, _ := biz.NewCmsAdmin().RoleSiteFind(GlobalRoleId)
		list := make([]interface{}, 0)
		var totalCount int64 = 0
		if len(siteIdList) > 0 {
			for i := 0; i < len(siteIdList); i++ {
				topicList, count, _ := biz.NewCmsTopic().TopicPaginate(page, limit, siteIdList[i].SiteID, -1, name, title, status)
				totalCount += count
				if count > 0 {
					for j := 0; j < len(topicList); j++ {
						list = append(list, topicList[j])
					}
				}
			}
		}
		c.JSONPage(lib.CodeSuccess, "", list, totalCount)
	}
}
