package tool

import (
	"strings"

	"github.com/beego/beego/v2/server/web"
	"github.com/zsy619/tools/xcas"
	"github.com/zsy619/tools/xgeneric"
)

const (
	Url_School_Login = "/cas/school/login"
	Url_Admin_Login  = "/cas/admin/login"

	Url_School_Index = "/admin/home/index"
	Url_Admin_Index  = "/admin/home/index"

	Url_File_Upload = "/file/upload"
)

var (
	CasEnabled              bool
	CasHost                 string
	CasLoginPath            string
	CasLogoutPath           string
	CasValidatePath         string
	CaseServiceValidatePath string

	LoginPathOfStudent string
	LoginPathOfCompany string
	LoginPathOfSchool  string
	LoginPathOfAdmin   string

	LoginCasPathOfStudent string
	LoginCasPathOfCompany string
	LoginCasPathOfSchool  string
	LoginCasPathOfAdmin   string

	Smn_Url      string
	Smn_Username string
	Smn_Password string
	Smn_Auth     string
	Smn_Salt     string
	Smn_Alias    string
	Smn_Subject  string

	local_domain string
)

func init() {
	CasEnabled, _ = web.AppConfig.Bool("cas.enabled")
	CasHost, _ = web.AppConfig.String("cas.url")
	CasHost = strings.TrimSuffix(CasHost, "/")
	CasLoginPath = CasHost + xcas.CASLoginURI
	CasLogoutPath = CasHost + xcas.CASLogoutURI
	CasValidatePath = CasHost + xcas.CASValidateURI
	CaseServiceValidatePath = CasHost + xcas.CASVersion2ServiceValidateURI

	LoginPathOfSchool = C_LOCAL_DOMAIN() + Url_School_Login
	LoginPathOfAdmin = C_LOCAL_DOMAIN() + Url_Admin_Login

	LoginCasPathOfStudent = CasLoginPath + "?service=" + LoginPathOfStudent
	LoginCasPathOfCompany = CasLoginPath + "?service=" + LoginPathOfCompany
	LoginCasPathOfSchool = CasLoginPath + "?service=" + LoginPathOfSchool
	LoginCasPathOfAdmin = CasLoginPath + "?admin=1&service=" + LoginPathOfAdmin

	Smn_Url = C("smn.url", "http://smn.hnzhjypt.com")
	Smn_Username = C("smn.username", "")
	Smn_Password = C("smn.password", "")
	Smn_Auth = C("smn.auth", "")
	Smn_Salt = C("smn.salt", "")
	Smn_Alias = C("smn.alias", "")
	Smn_Subject = C("smn.subject", "")

	// 初始化 Google OAuth 配置
	InitGoogleOAuth()
}

// GoogleOAuth Google OAuth 配置(包级 init 时调用)
var (
	GoogleEnabled      bool
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURI  string
	GoogleAuthURL      string
	GoogleTokenURL     string
	GoogleUserInfoURL  string
)

func InitGoogleOAuth() {
	GoogleEnabled, _ = web.AppConfig.Bool("google.enabled")
	GoogleClientID = C("google.clientID", "")
	GoogleClientSecret = C("google.clientSecret", "")
	GoogleRedirectURI = C_LOCAL_DOMAIN() + "/cms/admin/google/callback"
	GoogleAuthURL = "https://accounts.google.com/o/oauth2/v2/auth"
	GoogleTokenURL = "https://oauth2.googleapis.com/token"
	GoogleUserInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"
}

func C(name, value string) string {
	rt, err := web.AppConfig.String(name)
	if err != nil {
		return value
	}
	return rt
}

func C_SITE_NAME() string {
	return C("site.name", "")
}

func C_LOCAL_DOMAIN() string {
	if local_domain == "" {
		return C("LOCAL_DOMAIN", "")
	}
	return local_domain
}

func C_LOCAL_DOMAIN_Backslash() string {
	local_domain := C_LOCAL_DOMAIN()
	return xgeneric.IFF(strings.HasSuffix(local_domain, "/"), local_domain, local_domain+"/")
}

func C_LOCAL_DOMAIN_SET(domain string) {
	local_domain = domain
}

func Easemob_OrgName() string {
	return C("easemob.orgname", "1112200414042300")
}

func Easemob_AppName() string {
	return C("easemob.appname", "yunshipin")
}

func Easemob_Url() string {
	return C("easemob.url", "https://a1.easemob.com")
}

func Easemob_ClientId() string {
	return C("easemob.clientid", "YXA6BzydySCpRqGXlSgMqAo8_Q")
}

func Easemob_ClientSecret() string {
	return C("easemob.clientsecret", "YXA6WwTPDaYDu7oGKlO3uPD6WqRmz0c")
}
