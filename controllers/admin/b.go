package admin

import (
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/controllers/admin/vmodel"
	"haedu.gov.cn/cms/global"
	"html/template"
	"regexp"
	"strings"

	"haedu.gov.cn/cms/app/dal/model"

	"haedu.gov.cn/cms/controllers"
)

type BaseController struct {
	controllers.BaseController
}

func (c *BaseController) Prepare() {
	fmt.Println("Admin BaseController Prepare")
	c.BaseController.Prepare()
	c.EnableXSRF = true
	c.XSRFExpire = 3600
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.Data["xsrf_token"] = c.XSRFToken()

	if GlobalAdminId == 0 {
		user := c.GetSession("user").(*model.CmsAdmin)
		GlobalAdminId = user.UserID
		GlobalUserType = int(user.UserType) // 1:管理员 2:学校
		GlobalAdminName = user.UserName
		GlobalRealName = user.RealName
		GlobalRoleId = user.RoleID
		GlobalRoleType = user.RoleType
	}
	c.Data["roleId"] = GlobalRoleId
}

func (c *BaseController) Finish() {
	fmt.Println("Admin BaseController Finish")
}

// 渲染模版
func (this *BaseController) display(tpl ...string) {
	var tplname string
	if len(tpl) > 0 {
		tplname = tpl[0] + ".html"
	} else {
		tplname = "admin/" + this.ControllerName + "/" + this.ActionName + ".html"
	}
	this.Layout = "admin/layout/layout.html"
	this.TplName = tplname
}

// 渲染模版
func (this *BaseController) displayNoLayout(tpl ...string) {
	var tplname string
	if len(tpl) > 0 {
		tplname = tpl[0] + ".html"
	} else {
		tplname = "admin/" + this.ControllerName + "/" + this.ActionName + ".html"
	}
	this.TplName = tplname
}

// 登录人ID
func (this *BaseController) IsLogin() int64 {
	id := this.GetSession(`adminId`)
	if id == nil {
		return 0
	} else {
		switch id.(type) {
		case int64:
			rt := id.(int64)
			return rt
		default:
			return 0
		}
	}
}

func (this *BaseController) RolePowerGet(navName string) vmodel.RoleAction {
	roleAction := vmodel.RoleAction{}
	if global.IsSuper(GlobalRoleType) {
		roleAction.IsSuccess = true
		roleAction.IsSuper = true
		roleAction.IsHasAdd = true
		roleAction.IsHasAudit = true
		roleAction.IsHasEdit = true
		roleAction.IsHasDelete = true
		roleAction.IsHasView = true
		roleAction.IsHasAttach = true
		roleAction.IsHasAlbum = true
	} else {
		mdl, err := biz.NewCmsAdmin().RolePower(GlobalRoleId, navName)
		if err != nil {
			logs.Error("RoleValueFind", err.Error())
		}
		if len(mdl.Action) > 0 {
			action := strings.Split(mdl.Action, ",")
			if len(action) > 0 {
				roleAction.IsSuccess = true
				for i := 0; i < len(action); i++ {
					if action[i] == "Add" { // 添加、拷贝
						roleAction.IsHasAdd = true
					}
					if action[i] == "Audit" { // 审核
						roleAction.IsHasAudit = true
					}
					if action[i] == "Edit" { // 编辑、保存
						roleAction.IsHasEdit = true
					}
					if action[i] == "Delete" { // 删除
						roleAction.IsHasDelete = true
					}
					if action[i] == "View" { // 查看
						roleAction.IsHasView = true
					}
					if action[i] == "Attach" { // 附件
						roleAction.IsHasAttach = true
					}
					if action[i] == "Album" { // 相册
						roleAction.IsHasAlbum = true
					}
				}
			}
		}
	}
	return roleAction
}

// CheckPasswordRole 检查密码规则
func CheckPasswordRole(ps string) error {
	if len(ps) < 8 {
		return fmt.Errorf("密码长度不得小于8位")
	}
	num := `[0-9]{1}`
	a_z := `[a-z]{1}`
	A_Z := `[A-Z]{1}`
	symbol := `[!@#~$%^&*()+|_]{1}`
	if b, err := regexp.MatchString(num, ps); !b || err != nil {
		return fmt.Errorf("密码中必须包含数字")
	}
	if b, err := regexp.MatchString(a_z, ps); !b || err != nil {
		return fmt.Errorf("密码中必须包含小写字母")
	}
	if b, err := regexp.MatchString(A_Z, ps); !b || err != nil {
		return fmt.Errorf("密码中必须包含大写字母")
	}
	if b, err := regexp.MatchString(symbol, ps); !b || err != nil {
		return fmt.Errorf("密码中必须包含特殊字符")
	}
	return nil
}
