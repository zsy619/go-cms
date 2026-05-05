package admin

import (
	"github.com/beego/beego/v2/core/logs"
	"github.com/zsy619/tools/xjson"

	"haedu.gov.cn/cms/app/cms/domain"
	"haedu.gov.cn/cms/app/cms/service"
	lib "haedu.gov.cn/cms/app/tool"
	"haedu.gov.cn/cms/controllers/admin/vmodel"
	"haedu.gov.cn/cms/global"
)

type TopicController struct{ BaseController }

// Index 标签管理
// @router /admin/Topic/index [get]
func (ctrl *TopicController) Index() {
	siteList, _ := service.NewCmsTopic().SiteGet(GlobalRoleId, GlobalRoleType)
	ctrl.Data["site"] = siteList
	// 获取角色权限
	roleMap := ctrl.RolePowerGet("topic")
	ctrl.Data["roleMap"] = roleMap
	ctrl.display()
}

// TopicEdit 编辑
// @router /admin/Topic/TopicEdit [get]
func (ctrl *TopicController) TopicEdit() {
	topicId, _ := ctrl.GetInt64("topicId")
	clone, _ := ctrl.GetInt("clone")
	mdl, err := service.NewCmsTopic().TopicFind(topicId)
	if err != nil {
		mdl = &domain.CmsTopic{
			SortID: 99,
		}
	}
	// 是否克隆
	if clone == 1 {
		mdl.TopicID = 0
	}
	ctrl.Data["mdl"] = mdl
	siteList, _ := service.NewCmsTopic().SiteGet(GlobalRoleId, GlobalRoleType)
	ctrl.Data["site"] = siteList
	// 获取角色权限
	roleMap := ctrl.RolePowerGet("topic")
	ctrl.Data["roleMap"] = roleMap
	ctrl.display()
}

// TopicSave 保存
// @router /admin/Topic/TopicSave [post]
func (ctrl *TopicController) TopicSave() {
	mdl := domain.CmsTopic{}
	if err := ctrl.ParseForm(&mdl); err != nil {
		logs.Error("TopicSave", err.Error())
		ctrl.JSONError(err.Error())
	}
	if mdl.TopicID <= 0 {
		mdl.CreateID = int32(GlobalAdminId)
		mdl.CreateName = GlobalAdminName
	} else {
		mdl.UpdateID = int32(GlobalAdminId)
		mdl.UpdateName = GlobalAdminName
	}
	if err := service.NewCmsTopic().TopicSave(&mdl); err != nil {
		logs.Error("TopicSave", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("保存成功", nil)
}

// TopicSaveSortId 保存排序
// @router /admin/Topic/TopicSaveSortId [post]
func (ctrl *TopicController) TopicSaveSortId() {
	mdls := []vmodel.Topic_SaveSortIdModel{}
	data := ctrl.Ctx.Input.RequestBody
	logs.Debug("TopicSaveSortId", string(data))
	if err := xjson.Unmarshal(ctrl.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("TopicSaveSortId", err.Error())
		ctrl.JSONError(err.Error())
	}
	service := service.NewCmsTopic()
	for _, mdl := range mdls {
		if err := service.TopicSaveSortId(mdl.TopicId, int32(mdl.SortId)); err != nil {
			logs.Error("TopicSaveSortId", err.Error())
			ctrl.JSONError(err.Error())
			return
		}
	}
	ctrl.JSONSuccess("保存成功", nil)
}

// TopicDestory 删除
// @router /admin/Topic/TopicDestory [post]
func (ctrl *TopicController) TopicDestory() {
	topicId, _ := ctrl.GetInt64("topicId")
	if err := service.NewCmsTopic().TopicDestory(topicId); err != nil {
		logs.Error("TopicDestory", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("删除成功", nil)
}

// TopicChangeStatus 更改状态
// @router /admin/Topic/TopicChangeStatus [post]
func (ctrl *TopicController) TopicChangeStatus() {
	var mdl vmodel.Topic_ChangeStatusModel
	if err := xjson.Unmarshal(ctrl.Ctx.Input.RequestBody, &mdl); err != nil {
		logs.Error("TopicChangeStatus", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	for _, topicId := range mdl.TopicIds {
		if err := service.NewCmsTopic().TopicChangeStatus(topicId, mdl.Status); err != nil {
			logs.Error("TopicChangeStatus", err.Error())
			ctrl.JSONError(err.Error())
			return
		}
	}
	ctrl.JSONSuccess("更改状态成功", nil)
}

// TopicPaginate 列表
// @router /admin/Topic/TopicPaginate [get]
func (ctrl *TopicController) TopicPaginate() {
	page, limit := ctrl.GetPagingParameters()
	siteId, _ := ctrl.GetInt64("siteId")
	status, _ := ctrl.GetInt32("status")
	title := ctrl.GetSafeString("title")
	name := ctrl.GetSafeString("name")
	var siteIds []int64
	if siteId > 0 {
		siteIds = append(siteIds, siteId)
	} else {
		if !global.IsSuper(GlobalRoleType) {
			siteIdList, _, _ := service.NewCmsAdmin().RoleSiteFind(GlobalRoleId)
			for _, item := range siteIdList {
				siteIds = append(siteIds, item.SiteID)
			}
		}
	}
	list, count, _ := service.NewCmsTopic().TopicPaginate(page, limit, -1, name, title, status, siteIds...)
	ctrl.JSONPage(lib.CodeSuccess, "", list, count)
}
