package admin

import (
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/controllers/admin/vmodel"
	"haedu.gov.cn/tools/xjson"
	"time"
)

type NoticeController struct {
	BaseController
}

// Index 系统公告管理
// @router /admin/notice/index [get]
func (c *NoticeController) Index() {
	// 获取角色权限
	roleMap := c.RolePowerGet("notice")
	c.Data["roleMap"] = roleMap
	c.display()
}

// NoticeEdit 系统公告编辑
// @router /admin/notice/NoticeEdit [get]
func (c *NoticeController) NoticeEdit() {
	noticeId, _ := c.GetInt64("noticeId")
	mdl, err := biz.NewCmsNotice().NoticeFind(noticeId)
	if err != nil {
		mdl = &model.AdminNotice{
			SortID:      99,
			CreateName:  GlobalAdminName,
			PublishTime: time.Now(),
		}
	}
	c.Data["mdl"] = mdl
	// 获取角色权限
	roleMap := c.RolePowerGet("notice")
	c.Data["roleMap"] = roleMap
	c.display()
}

// NoticeHomePaginate 首页系统公告列表
func (c *NoticeController) NoticeHomePaginate() {
	noticeList, err := biz.NewCmsNotice().NoticeShow()
	if err != nil {
		logs.Error("NoticeHomePaginate", err.Error())
	}
	c.JSONSuccess("", noticeList)
}

// NoticePaginate 获取系统公告列表
func (c *NoticeController) NoticePaginate() {
	page, limit := c.GetPagingParameters()
	title := c.GetString("title")
	status, _ := c.GetInt32("status")
	noticeList, count, err := biz.NewCmsNotice().NoticePaginate(page, limit, title, status, GlobalAdminId, GlobalRoleType)
	if err != nil {
		logs.Error("NoticePaginate", err.Error())
	}
	c.JSONPageSuccess(noticeList, count)
}

// NoticeSaveSortId 保存排序
// @router /admin/notice/NoticeSaveSortId [post]
func (c *NoticeController) NoticeSaveSortId() {
	mdls := []vmodel.Notice_SaveSortIdModel{}
	data := c.Ctx.Input.RequestBody
	fmt.Println("NoticeSaveSortId", string(data))
	if err := xjson.Unmarshal(c.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("NoticeSaveSortId", err.Error())
		c.JSONError(err.Error())
	}
	service := biz.NewCmsNotice()
	for _, mdl := range mdls {
		if err := service.NoticeSaveSortId(mdl.NoticeId, int32(GlobalAdminId), GlobalAdminName, int32(mdl.SortId)); err != nil {
			logs.Error("NoticeSaveSortId", err.Error())
			c.JSONError(err.Error())
			return
		}
	}
	c.JSONSuccess("保存成功", nil)
}

// NoticeChangeStatus 更改状态
// @router /admin/notice/NoticeChangeStatus [post]
func (c *NoticeController) NoticeChangeStatus() {
	var mdl vmodel.Notice_ChangeStatusModel
	if err := xjson.Unmarshal(c.Ctx.Input.RequestBody, &mdl); err != nil {
		logs.Error("NoticeChangeStatus", err.Error())
		c.JSONError(err.Error())
		return
	}
	for _, noticeId := range mdl.NoticeIds {
		if err := biz.NewCmsNotice().NoticeChangeStatus(noticeId, int32(GlobalAdminId), GlobalAdminName, mdl.Status); err != nil {
			logs.Error("NoticeChangeStatus", err.Error())
			c.JSONError(err.Error())
			return
		}
	}
	c.JSONSuccess("更改状态成功", nil)
}

// NoticeDestroy 删除
// @router /admin/notice/NoticeDestroy [post]
func (c *NoticeController) NoticeDestroy() {
	noticeId, _ := c.GetInt64("noticeId")
	if err := biz.NewCmsNotice().NoticeDestroy(noticeId); err != nil {
		logs.Error("NoticeDestroy", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("删除成功", nil)
}

// NoticeSave 保存
// @router /admin/notice/NoticeSave [post]
func (c *NoticeController) NoticeSave() {
	mdl := model.AdminNotice{}
	if err := c.ParseForm(&mdl); err != nil {
		logs.Error("NoticeSave", err.Error())
		c.JSONError(err.Error())
	}
	if mdl.NoticeID > 0 {
		mdl.UpdateID = int32(GlobalAdminId)
		mdl.UpdateName = GlobalAdminName
	} else {
		mdl.CreateID = int32(GlobalAdminId)
		mdl.CreateName = GlobalAdminName
	}
	if err := biz.NewCmsNotice().NoticeSave(&mdl); err != nil {
		logs.Error("NoticeSave", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("保存成功", nil)
}
