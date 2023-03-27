package lib

import (
	"strings"

	"github.com/beego/beego/v2/server/web"
	"haedu.gov.cn/tools/xcas"
	"haedu.gov.cn/tools/xgeneric"
)

const (
	Url_Student_Login = "/student/login"
	Url_Company_Login = "/company/login"
	Url_School_Login  = "/school/login"
	Url_Admin_Login   = "/admin/login"

	Url_Student_Index = "/stu/home/index"
	Url_Company_Index = "/cpy/self/index"
	Url_School_Index  = "/admin/home/index"
	Url_Admin_Index   = "/admin/home/index"

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
	if strings.HasSuffix(CasHost, "/") {
		CasHost = CasHost[:len(CasHost)-1]
	}
	CasLoginPath = CasHost + xcas.CASLoginURI
	CasLogoutPath = CasHost + xcas.CASLogoutURI
	CasValidatePath = CasHost + xcas.CASValidateURI
	CaseServiceValidatePath = CasHost + xcas.CASVersion2ServiceValidateURI

	LoginPathOfStudent = C_LOCAL_DOMAIN() + Url_Student_Login
	LoginPathOfCompany = C_LOCAL_DOMAIN() + Url_Company_Login
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
