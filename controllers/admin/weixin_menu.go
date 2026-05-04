package admin

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	"github.com/zsy619/tools/xjson"

	"haedu.gov.cn/cms/app/cms/domain"
	"haedu.gov.cn/cms/app/cms/service"
	lib "haedu.gov.cn/cms/app/tool"
	"haedu.gov.cn/cms/app/wechat/mp"
	"haedu.gov.cn/cms/controllers/admin/vmodel"
)

func (ctrl *WeixinController) MenuFind() {
	accountId, _ := ctrl.GetInt64("accountId")
	list, count, err := service.NewWeixinMenu().MenuPaginate(1, 99999, accountId)
	if err != nil {
		logs.Error("MenuFind", err.Error())
	}
	ctrl.JSONPage(lib.CodeSuccess, "", list, count)
}

func (ctrl *WeixinController) MenuEdit() {
	accountId, _ := ctrl.GetInt64("accountId")
	if accountId <= 0 {
		ctrl.Abort("404")
		ctrl.StopRun()
		return
	}
	ctrl.Data["accountId"] = accountId
	menuId, _ := ctrl.GetInt64("menuId")
	parentId, _ := ctrl.GetInt64("parentId")
	mdl, err := service.NewWeixinMenu().MenuFind(menuId)
	if err != nil {
		mdl = &domain.WeixinMenu{
			AccountID: accountId,
			ParentID:  parentId,
			SortID:    99,
			Type:      "view",
		}
	}
	ctrl.Data["mdl"] = mdl
	ctrl.display()
}

func (ctrl *WeixinController) MenuSave() {
	mdl := domain.WeixinMenu{}
	if err := ctrl.ParseForm(&mdl); err != nil {
		logs.Error("MenuSave", err.Error())
		ctrl.JSONError(err.Error())
	}
	if mdl.MenuID == 0 {
		mdl.CreateID = int32(GlobalAdminId)
		mdl.CreateName = GlobalAdminName
	} else {
		mdl.UpdateID = int32(GlobalAdminId)
		mdl.UpdateName = GlobalAdminName
	}
	if err := service.NewWeixinMenu().MenuSave(&mdl); err != nil {
		logs.Error("MenuSave", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("保存成功", nil)
}

func (ctrl *WeixinController) MenuSaveSortId() {
	mdls := []vmodel.Menu_SaveSortIdModel{}
	data := ctrl.Ctx.Input.RequestBody
	fmt.Println("MenuSaveSortId", string(data))
	if err := xjson.Unmarshal(ctrl.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("MenuSaveSortId", err.Error())
		ctrl.JSONError(err.Error())
	}
	for _, mdl := range mdls {
		if err := service.NewWeixinMenu().MenuSaveSortId(mdl.MenuId, int32(mdl.SortId)); err != nil {
			logs.Error("MenuSaveSortId", err.Error())
			ctrl.JSONError(err.Error())
			return
		}
	}
	ctrl.JSONSuccess("保存成功", nil)
}

func (ctrl *WeixinController) MenuDestory() {
	menuId, _ := ctrl.GetInt64("menuId")
	if err := service.NewWeixinMenu().MenuDestory(menuId); err != nil {
		logs.Error("MenuDestory", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("删除成功", nil)
}

// MenuSync 同步菜单
// @router /admin/weixin/menusync [post]
func (ctrl *WeixinController) MenuSync() {
	accountId, _ := ctrl.GetInt64("accountId")
	finder, err := service.NewWeixinAccount().AccountFind(accountId)
	if err != nil {
		logs.Error("MenuSync", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	fmt.Println("MenuSync", finder)
	message := mp.NewMessage(finder.AppID, finder.AppSecret, true)
	pbuttons := []mp.Button{}
	service := service.NewWeixinMenu()
	parents, err := service.MenuFindByParentId(accountId, 0)
	if err == nil {
		for _, p := range parents {
			pbutton := mp.Button{
				Name:      p.Name,
				Key:       p.Key,
				Url:       p.URL,
				Type:      p.Type,
				SubButton: []mp.Button{},
			}
			children, err := service.MenuFindByParentId(accountId, p.MenuID)
			if err == nil && len(children) > 0 {
				for _, child := range children {
					cbutton := mp.Button{
						Name: child.Name,
						Key:  child.Key,
						Type: child.Type,
						Url:  child.URL,
					}
					pbutton.SubButton = append(pbutton.SubButton, cbutton)
				}
			}
			pbuttons = append(pbuttons, pbutton)
		}
	}

	err = message.CreateCustomMenu(&pbuttons)
	if err != nil {
		logs.Error("MenuSync", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("同步成功", nil)
}
