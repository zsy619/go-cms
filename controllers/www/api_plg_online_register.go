package www

import (
	"github.com/beego/beego/v2/core/logs"

	"haedu.gov.cn/cms/app/cms/domain"
	"haedu.gov.cn/cms/app/cms/service"
)

type PlgOnlineRegisterController struct{ BaseController }

// @router /api/plg/online/register/save [post]
func (ctrl *PlgOnlineRegisterController) Save() {
	mdl := domain.PlgOnlineRegister{}
	if err := ctrl.ParseForm(&mdl); err != nil {
		logs.Error("Save", err.Error())
		ctrl.JSONError(err.Error())
	}
	if err := service.NewPlgOnlineRegister().Save(&mdl); err != nil {
		logs.Error("Save", err.Error())
		ctrl.JSONError(err.Error())
		return
	}
	ctrl.JSONSuccess("保存成功", nil)
}
