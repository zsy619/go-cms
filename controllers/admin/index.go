package admin

import (
	"fmt"

	"haedu.gov.cn/cms/app/biz"
)

type IndexController struct {
	BaseController
}

func (c *IndexController) Index() {
	c.Data["userName"] = GlobalRealName
	c.displayNoLayout()
}

func (c *IndexController) Welcome() {
	c.display()
}

func (c *IndexController) UserPassword() {
	c.display()
}

// UserPasswordSave 修改密码
// @router /admin/index/UserPasswordSave [post]
func (c *IndexController) UserPasswordSave() {
	oldPassword := c.GetString("old_password")
	newPassword := c.GetString("new_password")
	confirmPassword := c.GetString("again_password")
	fmt.Println(oldPassword, newPassword, confirmPassword)

	if newPassword != confirmPassword {
		c.JSONError("两次输入的密码不一致")
	}
	err := biz.NewCmsAdmin().ModifyPassword(GlobalAdminId, oldPassword, newPassword)
	if err != nil {
		c.JSONError(err.Error())
	}

	c.JSONSuccess("修改成功", nil)
}

func (c *IndexController) UserSetting() {
	c.display()
}
