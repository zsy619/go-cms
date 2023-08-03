package www

import (
	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/dal/model"
)

type PlgOnlineRegisterController struct{ BaseController }

// @router /api/plg/online/register/save [post]
func (this *PlgOnlineRegisterController) Save() {
	mdl := model.PlgOnlineRegister{}
	if err := this.ParseForm(&mdl); err != nil {
		logs.Error("Save", err.Error())
		this.JSONError(err.Error())
	}
	if err := biz.NewPlgOnlineRegister().Save(&mdl); err != nil {
		logs.Error("Save", err.Error())
		this.JSONError(err.Error())
		return
	}
	this.JSONSuccess("保存成功", nil)
}
