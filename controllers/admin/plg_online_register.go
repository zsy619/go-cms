package admin

import (
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/app/lib"
	"haedu.gov.cn/cms/controllers/admin/vmodel"
	"haedu.gov.cn/tools/xjson"
)

type PlgOnlineRegisterController struct{ BaseController }

// @router /admin/plg/register [get]
func (this *PlgOnlineRegisterController) Index() {
	// fmt.Println(this.ControllerName, this.ActionName)
	this.display()
}

// @router /admin/plg/register/paginate [get]
func (this *PlgOnlineRegisterController) Paginate() {
	service := biz.NewPlgOnlineRegister()
	page, limit := this.GetPagingParameters()
	realName := this.GetString("realName")
	special := this.GetString("special")
	degree := this.GetString("degree")
	tags := this.GetString("tags")
	remark := this.GetString("remark")
	isRead, _ := this.GetInt32("isRead")
	list, total, err := service.Paginate(page, limit, realName, special, degree, tags, remark, isRead)
	if err != nil {
		logs.Error(err.Error())
	}
	this.JSONPage(lib.CodeSuccess, "", list, total)
}

// @router /admin/plg/register/edit [get]
func (this *PlgOnlineRegisterController) Edit() {
	registerId, _ := this.GetInt64("registerId")
	service := biz.NewPlgOnlineRegister()
	mdl, err := service.Find(registerId)
	if err != nil {
		logs.Error(err.Error())
		mdl = &model.PlgOnlineRegister{
			CreateTime: time.Now(),
		}
	}
	this.Data["mdl"] = mdl
	if registerId > 0 {
		service.ChangeRead(registerId, 1, int32(GlobalAdminId), GlobalAdminName)
	}
	this.display()
}

// @router /admin/plg/register/save [post]
func (this *PlgOnlineRegisterController) Save() {
	mdl := model.PlgOnlineRegister{}
	if err := this.ParseForm(&mdl); err != nil {
		logs.Error("Save", err.Error())
		this.JSONError(err.Error())
	}
	mdl.Tags = strings.ReplaceAll(mdl.Tags, "，", ",")
	mdl.CreateID = int32(GlobalAdminId)
	mdl.CreateName = GlobalAdminName
	mdl.UpdateID = int32(GlobalAdminId)
	mdl.UpdateName = GlobalAdminName
	if err := biz.NewPlgOnlineRegister().Save(&mdl); err != nil {
		logs.Error("Save", err.Error())
		this.JSONError(err.Error())
		return
	}
	this.JSONSuccess("保存成功", nil)
}

// @router /admin/plg/register/destory [post]
func (this *PlgOnlineRegisterController) Destory() {
	registerId, _ := this.GetInt64("registerId")
	if err := biz.NewPlgOnlineRegister().Destory(registerId); err != nil {
		logs.Error("Destory", err.Error())
		this.JSONError(err.Error())
		return
	}
	this.JSONSuccess("删除成功", nil)
}

// @router /admin/plg/register/read [post]
func (this *PlgOnlineRegisterController) ChangeRead() {
	var mdl vmodel.Register_ChangeReadModel
	if err := xjson.Unmarshal(this.Ctx.Input.RequestBody, &mdl); err != nil {
		logs.Error("ChangeRead", err.Error())
		this.JSONError(err.Error())
		return
	}
	for _, registerId := range mdl.RegisterIds {
		if err := biz.NewPlgOnlineRegister().ChangeRead(registerId, mdl.IsRead, int32(GlobalAdminId), GlobalAdminName); err != nil {
			logs.Error("ChangeRead", err.Error())
			this.JSONError(err.Error())
			return
		}
	}
	this.JSONSuccess("操作成功", nil)
}
