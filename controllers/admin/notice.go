package admin

import (
	"fmt"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/tools/xjson"

	"haedu.gov.cn/cms/app/cms/domain"
	"haedu.gov.cn/cms/app/cms/service"
	"haedu.gov.cn/cms/controllers/admin/vmodel"
)

type NoticeController struct{ BaseController }

// Index 系统公告管理
// @router /admin/notice/index [get]
func (ctrl *NoticeController) Index() {
	// 获取角色权限
	roleMap := ctrl.RolePowerGet("notice")
	ctrl.Data["roleMap"] = roleMap
	ctrl.display()
}

// NoticeEdit 系统公告编辑
// @router /admin/notice/NoticeEdit [get]
func (ctrl *NoticeController) NoticeEdit() {
	noticeId, _ := ctrl.GetInt64("noticeId")
	mdl, err := service.NewCmsAdminNotice().NoticeFind(noticeId)
	if err != nil {
		mdl = &domain.CmsAdminNotice{
			SortID:      99,
			CreateName:  GlobalAdminName,
			PublishTime: time.Now(),
		}
	}
	ctrl.Data["mdl"] = mdl
	// 获取角色权限
	roleMap := ctrl.RolePowerGet("notice")
	ctrl.Data["roleMap"] = roleMap
	ctrl.display()
}

// NoticeHomePaginate 首页系统公告列表
func (ctrl *NoticeController) NoticeHomePaginate() {
	noticeList, err := service.NewCmsAdminNotice().NoticeShow()
	if err != nil {
		logs.Error("NoticeHomePaginate", err.Error())
	}
	ctrl.JSONSuccess("", noticeList)
}

// NoticePaginate 获取系统公告列表
func (ctrl *NoticeController) NoticePaginate() {
	page, limit := ctrl.GetPagingParameters()
	title := ctrl.GetSafeString("title")
	status, _ := ctrl.GetInt32("status")
	noticeList, count, err := service.NewCmsAdminNotice().NoticePaginate(page, limit, title, status, GlobalAdminId, GlobalRoleType)
	if err != nil {
		logs.Error("NoticePaginate", err.Error())
	}
	ctrl.JSONPageSuccess(noticeList, count)
}

// NoticeSaveSortId 保存排序
// @router /admin/notice/NoticeSaveSortId [post]
func (ctrl *NoticeController) NoticeSaveSortId() {
	mdls := []vmodel.Notice_SaveSortIdModel{}
	data := ctrl.Ctx.Input.RequestBody
	fmt.Println("NoticeSaveSortId", string(data))
	if err := xjson.Unmarshal(ctrl.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("NoticeSaveSortId", err.Error())
		ctrl.JSONError(err.Error())
	}
	service := service.NewCmsAdminNotice()
	for _, mdl := range mdls {
		if err := service.NoticeSaveSortId(mdl.NoticeId, int32(GlobalAdminId), GlobalAdminName, int32(mdl.SortId)); err != nil {
			logs.Error("NoticeSaveSortId", err.Error())
			ctrl.JSONError(err.Error())
			return
		}
	}
	ctrl.JSONSuccess("保存成功", nil)
}

// NoticeChangeStatus 更改状态
// @router /admin/notice/NoticeChangeStatus [post]
func (ctrl *NoticeController) NoticeChangeStatus() {
	var mdl vmodel.Notice_ChangeStatusModel
	if err := xjson.Unmarshal(ctrl.Ctx.Input.RequestBody, &mdl); err != nil {
		logs.Error("NoticeChangeStatus", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	for _, noticeId := range mdl.NoticeIds {
		if err := service.NewCmsAdminNotice().NoticeChangeStatus(noticeId, int32(GlobalAdminId), GlobalAdminName, mdl.Status); err != nil {
			logs.Error("NoticeChangeStatus", err.Error())
			ctrl.JSONError(err.Error())
			return
		}
	}
	ctrl.JSONSuccess("更改状态成功", nil)
}

// NoticeDestroy 删除
// @router /admin/notice/NoticeDestroy [post]
func (ctrl *NoticeController) NoticeDestroy() {
	noticeId, _ := ctrl.GetInt64("noticeId")
	if err := service.NewCmsAdminNotice().NoticeDestroy(noticeId); err != nil {
		logs.Error("NoticeDestroy", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("删除成功", nil)
}

// NoticeSave 保存
// @router /admin/notice/NoticeSave [post]
func (ctrl *NoticeController) NoticeSave() {
	mdl := domain.CmsAdminNotice{}
	if err := ctrl.ParseForm(&mdl); err != nil {
		logs.Error("NoticeSave", err.Error())
		ctrl.JSONError(err.Error())
	}
	if mdl.NoticeID > 0 {
		mdl.UpdateID = int32(GlobalAdminId)
		mdl.UpdateName = GlobalAdminName
	} else {
		mdl.CreateID = int32(GlobalAdminId)
		mdl.CreateName = GlobalAdminName
	}
	if err := service.NewCmsAdminNotice().NoticeSave(&mdl); err != nil {
		logs.Error("NoticeSave", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("保存成功", nil)
}
