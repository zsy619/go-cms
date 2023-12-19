package admin

import (
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/tools/xjson"

	"haedu.gov.cn/cms/app/cms/domain"
	"haedu.gov.cn/cms/app/cms/service"
	"haedu.gov.cn/cms/app/lib"
	"haedu.gov.cn/cms/controllers/admin/vmodel"
)

type PlgOnlineRegisterController struct{ BaseController }

// @router /admin/plg/register [get]
func (ctrl *PlgOnlineRegisterController) Index() {
	// fmt.Println(ctrl.ControllerName, ctrl.ActionName)
	ctrl.display()
}

// @router /admin/plg/register/paginate [get]
func (ctrl *PlgOnlineRegisterController) Paginate() {
	service := service.NewPlgOnlineRegister()
	page, limit := ctrl.GetPagingParameters()
	realName := ctrl.GetSafeString("realName")
	special := ctrl.GetSafeString("special")
	degree := ctrl.GetSafeString("degree")
	tags := ctrl.GetSafeString("tags")
	remark := ctrl.GetSafeString("remark")
	isRead, _ := ctrl.GetInt32("isRead")
	list, total, err := service.Paginate(page, limit, realName, special, degree, tags, remark, isRead)
	if err != nil {
		logs.Error(err.Error())
	}
	ctrl.JSONPage(lib.CodeSuccess, "", list, total)
}

// @router /admin/plg/register/edit [get]
func (ctrl *PlgOnlineRegisterController) Edit() {
	registerId, _ := ctrl.GetInt64("registerId")
	service := service.NewPlgOnlineRegister()
	mdl, err := service.Find(registerId)
	if err != nil {
		logs.Error(err.Error())
		mdl = &domain.PlgOnlineRegister{
			CreateTime: time.Now(),
			IP:         ctrl.Ctx.Input.IP(),
		}
	}
	ctrl.Data["mdl"] = mdl
	if registerId > 0 {
		service.ChangeRead(registerId, 1, int32(GlobalAdminId), GlobalAdminName)
	}
	ctrl.display()
}

// @router /admin/plg/register/save [post]
func (ctrl *PlgOnlineRegisterController) Save() {
	mdl := domain.PlgOnlineRegister{}
	if err := ctrl.ParseForm(&mdl); err != nil {
		logs.Error("Save", err.Error())
		ctrl.JSONError(err.Error())
	}
	mdl.Tags = strings.ReplaceAll(mdl.Tags, "，", ",")
	mdl.CreateID = int32(GlobalAdminId)
	mdl.CreateName = GlobalAdminName
	mdl.UpdateID = int32(GlobalAdminId)
	mdl.UpdateName = GlobalAdminName
	if err := service.NewPlgOnlineRegister().Save(&mdl); err != nil {
		logs.Error("Save", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("保存成功", nil)
}

// @router /admin/plg/register/destory [post]
func (ctrl *PlgOnlineRegisterController) Destory() {
	registerId, _ := ctrl.GetInt64("registerId")
	if err := service.NewPlgOnlineRegister().Destory(registerId); err != nil {
		logs.Error("Destory", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("删除成功", nil)
}

// @router /admin/plg/register/read [post]
func (ctrl *PlgOnlineRegisterController) ChangeRead() {
	var mdl vmodel.Register_ChangeReadModel
	if err := xjson.Unmarshal(ctrl.Ctx.Input.RequestBody, &mdl); err != nil {
		logs.Error("ChangeRead", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	for _, registerId := range mdl.RegisterIds {
		if err := service.NewPlgOnlineRegister().ChangeRead(registerId, mdl.IsRead, int32(GlobalAdminId), GlobalAdminName); err != nil {
			logs.Error("ChangeRead", err.Error())
			ctrl.JSONError(err.Error())
			return
		}
	}
	ctrl.JSONSuccess("操作成功", nil)
}
