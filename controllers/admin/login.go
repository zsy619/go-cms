package admin

import (
	"net/http"

	"github.com/beego/beego/v2/core/logs"
	"github.com/zsy619/tools/xcache"
	"github.com/zsy619/tools/xcas"
	"github.com/zsy619/tools/xgeneric"

	"haedu.gov.cn/cms/app/cms/domain"
	"haedu.gov.cn/cms/app/cms/service"
	lib "haedu.gov.cn/cms/app/tool"
	"haedu.gov.cn/cms/controllers"
	"haedu.gov.cn/cms/controllers/admin/vmodel"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type LoginController struct{ controllers.BaseController }

// AdminLogin 管理员登录页面渲染
// @router cms/admin/login [get]
func (ctrl *LoginController) AdminLogin() {
	ctrl.Data["captcha"] = "/captcha"
	ctrl.TplName = "admin/login/login.html"
}

// SavaAdminState 保存管理员登录状态到Session和全局变量
// @param user *domain.CmsAdmin 管理员用户实体
func (ctrl *LoginController) SavaAdminState(user *domain.CmsAdmin) {
	ctrl.SetSession("adminId", user.UserID)
	ctrl.SetSession("adminAccount", user.UserName)
	ctrl.SetSession("adminName", user.UserName)
	ctrl.SetSession("realName", user.RealName)
	ctrl.SetSession("adminState", 1)
	ctrl.SetSession("adminLevel", 1)
	ctrl.SetSession("user", user)

	GlobalAdminId = user.UserID
	GlobalUserType = int(user.UserType)
	GlobalAdminName = user.UserName
	GlobalRealName = user.RealName
	GlobalIsAudit = user.IsAudit
	GlobalRoleId = user.RoleID
	GlobalRoleType = user.RoleType
}

// AdminLoginVerify 管理员登录验证处理
// @router cms/admin/login/verify [post]
func (ctrl *LoginController) AdminLoginVerify() {
	result := vmodel.LoginResult{
		Code:    0,
		Message: "登陆完成，载入中...",
		Url:     "",
	}

	captcha := ctrl.GetSafeString("captcha")
	if !controllers.VerifyCode(captcha) {
		result.Code = 1
		result.Message = "验证码错误"
		ctrl.Data["json"] = &result
		ctrl.ServeJSON()
		return
	}

	username := ctrl.GetSafeString("username")
	password := ctrl.GetSafeString("password")
	logs.Debug("登录尝试: username=%s, captcha=%s", username, captcha)

	adminDo := service.NewCmsAdmin()
	user, err := adminDo.Login(username, password, 0, service.LoginAll)
	if err != nil {
		logs.Error("登录失败: %v", err)
		result.Code = 2
		result.Message = "账号密码错误"
		ctrl.Data["json"] = &result
		ctrl.ServeJSON()
		return
	}
	ctrl.SavaAdminState(user)
	adminDo.LoginLog(user.UserID, username, "AdminLoginVerify", "", "", "OK", ctrl.GetClientIp())
	result.Url = "/admin/inde

// GoogleLogin Google OAuth 登录入口
// @router cms/admin/google [get]
func (ctrl *LoginController) GoogleLogin() {
	if !lib.GoogleEnabled {
		ctrl.Abort("404")
		return
	}
	params := url.Values{}
	params.Set("client_id", lib.GoogleClientID)
	params.Set("redirect_uri", lib.GoogleRedirectURI)
	params.Set("response_type", "code")
	params.Set("scope", "openid email profile")
	params.Set("access_type", "online")
	params.Set("prompt", "select_account")
	authURL := lib.GoogleAuthURL + "?" + params.Encode()
	ctrl.Redirect(authURL, 302)
}

// GoogleCallback Google OAuth 回调处理
// @router cms/admin/google/callback [get]
func (ctrl *LoginController) GoogleCallback() {
	if !lib.GoogleEnabled {
		ctrl.Abort("404")
		return
	}
	code := ctrl.GetString("code")
	if code == "" {
		ctrl.Abort("400")
		return
	}
	tokenParams := url.Values{}
	tokenParams.Set("code", code)
	tokenParams.Set("client_id", lib.GoogleClientID)
	tokenParams.Set("client_secret", lib.GoogleClientSecret)
	tokenParams.Set("redirect_uri", lib.GoogleRedirectURI)
	tokenParams.Set("grant_type", "authorization_code")
	resp, err := http.PostForm(lib.GoogleTokenURL, tokenParams)
	if err != nil || resp.StatusCode != 200 {
		logs.Error("Google Token exchange failed: %v status:%d", err, resp.StatusCode)
		ctrl.Abort("500")
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var tokenResp struct {
		AccessToken string `json:"access_token"`
		IDToken     string `json:"id_token"`
		ExpiresIn  int    `json:"expires_in"`
	}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		logs.Error("Parse Google token failed: %v", err)
		ctrl.Abort("500")
		return
	}
	userInfoResp, err := http.Get(lib.GoogleUserInfoURL + "?access_token=" + tokenResp.AccessToken)
	if err != nil || userInfoResp.StatusCode != 200 {
		logs.Error("Google UserInfo failed: %v", err)
		ctrl.Abort("500")
		return
	}
	defer userInfoResp.Body.Close()
	userInfoBody, _ := io.ReadAll(userInfoResp.Body)
	var userInfo struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		VerifiedEmail bool   `json:"verified_email"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
	}
	if err := json.Unmarshal(userInfoBody, &userInfo); err != nil {
		logs.Error("Parse Google userinfo failed: %v", err)
		ctrl.Abort("500")
		return
	}
	if userInfo.Email == "" || !userInfo.VerifiedEmail {
		logs.Warning("Google login: invalid email:%s verified:%v", userInfo.Email, userInfo.VerifiedEmail)
		ctrl.Abort("403")
		return
	}
	adminDo := service.NewCmsAdmin()
	user, err := adminDo.FindByAccount(userInfo.Email)
	if err != nil || user == nil || user.UserName == "" {
		logs.Info("Google login: creating user:%s", userInfo.Email)
		user = &domain.CmsAdmin{
			UserName:   userInfo.Email,
			NickName:   userInfo.Name,
			RealName:   userInfo.Name,
			Email:      userInfo.Email,
			UserNumber: userInfo.ID,
			Mobile:     userInfo.ID,
			UserType:   1,
			Enabled:    true,
			IsAudit:    1,
			Remark:     "Google OAuth",
		}
		if err := adminDo.AdminSave(user); err != nil {
			logs.Error("Google login: create user failed: %v", err)
			ctrl.Abort("500")
			return
		}
		user, err = adminDo.FindByAccount(userInfo.Email)
		if err != nil || user == nil {
			logs.Error("Google login: find user after create failed: %v", err)
			ctrl.Abort("500")
			return
		}
	}
	ctrl.SavaAdminState(user)
	adminDo.LoginLog(user.UserID, userInfo.Email, "GoogleOAuth", userInfo.Name, userInfo.Picture, "OK", ctrl.GetClientIp())
	ctrl.Redirect("/admin/index", 302)
}
x"
	ctrl.Data["json"] = &result
	ctrl.ServeJSON()
}

// Logout 管理员退出登录
// @router cms/admin/logout [get]
func (ctrl *LoginController) Logout() {
	ctrl.DestroySession()
	ctrl.Redirect("/admin/login", 302)
}

// School 学校用户登录入口
func (ctrl *LoginController) School() {
	ctrl.login(lib.LoginCasPathOfSchool, lib.LoginPathOfSchool, "school")
}

// Admin 管理员登录入口
func (ctrl *LoginController) Admin() {
	ctrl.login(lib.LoginCasPathOfAdmin, lib.LoginPathOfAdmin, "admin")
}

// login 处理CAS单点登录的核心逻辑
// @param loginCasPath CAS登录重定向路径
// @param loginPath 服务端回调路径
// @param kind 登录类型(school/admin)
func (ctrl *LoginController) login(loginCasPath, loginPath string, kind string) {
	xcache.SetDiskvString("cas_login", kind)
	GlobalAuthFlag = kind
	ticket := ctrl.GetSafeString("ticket")
	logs.Debug("CAS ticket: %s", ticket)

	if ticket == "" {
		ctrl.Redirect(loginCasPath, http.StatusFound)
		return
	}

	url := lib.CaseServiceValidatePath + "?service=" + loginPath + "&ticket=" + ticket
	logs.Debug("CAS验证URL: %s", url)

	serviceResponse, err := xcas.CasVersion2ServiceValidateAction(url)
	if err != nil {
		logs.Error("CAS验证请求失败: %v", err)
		ctrl.Redirect(loginCasPath, http.StatusFound)
		ctrl.StopRun()
		return
	}

	if serviceResponse.Failure != nil {
		logs.Warning("CAS验证失败: %s", serviceResponse.Failure.Message)
		ctrl.Ctx.ResponseWriter.Write([]byte(serviceResponse.Failure.Message))
		ctrl.StopRun()
	}

	if serviceResponse.Success != nil {
		kindx := ""
		for _, attribute := range serviceResponse.Success.Attributes.UserAttributes.Attributes {
			if attribute.Name == xcas.UserTypeField {
				kindx = attribute.Value
				break
			}
		}
		if kindx != kind {
			ctrl.Ctx.ResponseWriter.Write([]byte("invalid user type"))
			logs.Warning("用户类型不匹配: expected=%s, got=%s", kind, kindx)
			ctrl.StopRun()
		}

		account := serviceResponse.Success.User
		userDo := service.NewCmsAdmin()
		user, err := userDo.FindByAccount(account)
		if user == nil || user.UserName == "" || err != nil {
			logs.Warning("用户不存在，将创建新用户: account=%s, error=%v", account, err)
			user = &domain.CmsAdmin{
				UserName:   account,
				NickName:   account,
				RealName:   account,
				Email:      account + "@hnzhjypt.com",
				UserNumber: account,
				Mobile:     account,
				UserType:   xgeneric.IFF[int32](kind == "school", 2, 1),
				Enabled:    true,
				IsAudit:    1,
				Remark:     "CAS登录创建",
			}
			err := userDo.AdminSave(user)
			if err != nil {
				logs.Error("创建用户失败: %v", err)
				ctrl.Ctx.ResponseWriter.Write([]byte("创建用户失败"))
				ctrl.StopRun()
			}
		}

		user, err = userDo.FindByAccount(account)
		if err != nil {
			logs.Error("获取用户信息失败: %v", err)
			ctrl.Ctx.ResponseWriter.Write([]byte("未能获取用户信息"))
			ctrl.StopRun()
		}

		logs.Info("用户登录成功: userId=%d, account=%s", user.UserID, account)
		ctrl.SavaAdminState(user)
		userDo.LoginLog(user.UserID, account, "AdminLoginCase", "", "", "OK", ctrl.GetClientIp())
		url = "/admin/index"
		ctrl.Redirect(url, http.StatusFound)
		ctrl.StopRun()
	}
}
