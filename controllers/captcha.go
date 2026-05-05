package controllers

import (
	"encoding/json"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	captcha "github.com/mojocn/base64Captcha"
)

type CaptchaController struct {
	web.Controller
}

var (
	store    = captcha.DefaultMemStore
	verifyId = "captcha:yufei"
)

func NewDriver() *captcha.DriverString {
	driver := new(captcha.DriverString)
	driver.Height = 44
	driver.Width = 120
	driver.NoiseCount = 5
	// driver.ShowLineOptions = captcha.OptionShowSineLine | captcha.OptionShowSlimeLine | captcha.OptionShowHollowLine
	driver.ShowLineOptions = captcha.OptionShowHollowLine
	driver.Length = 6
	driver.Source = "123456789qwertyuipkjhgfdsazxcvbnm"
	driver.Fonts = []string{"wqy-microhei.ttc"}
	return driver
}

// 生成图形验证码
func (ctrl *CaptchaController) GenerateHandler() {
	driver := NewDriver().ConvertFonts()
	c := captcha.NewCaptcha(driver, store)
	id, content, answer := c.Driver.GenerateIdQuestionAnswer()
	logs.Debug(" ---> ", id, content, answer)
	item, _ := c.Driver.DrawCaptcha(content)
	_ = c.Store.Set(verifyId, answer)
	_, _ = item.WriteTo(ctrl.Ctx.ResponseWriter)
}

// 验证
func (ctrl *CaptchaController) VerifyHandle() {
	code := ctrl.Ctx.Request.FormValue("code")
	body := map[string]interface{}{"code": 1000, "msg": "failed"}
	if store.Verify(verifyId, code, true) {
		body = map[string]interface{}{"code": 1001, "msg": "ok"}
	}
	ctrl.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(ctrl.Ctx.ResponseWriter).Encode(body)
}

// VerifyCode 验证
func VerifyCode(code string) bool {
	return store.Verify(verifyId, code, true)
}
