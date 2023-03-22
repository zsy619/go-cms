package biz

import (
	"context"
	"errors"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/app/dal/query"
	"haedu.gov.cn/tools/xcrypto"
)

type LoginType int

const (
	LoginName LoginType = iota + 1
	LoginMobile
	LoginEmail
	LoginNameMobile
	LoginNameEmail
	LoginMobileEmail
	LoginAll
)

type CmsAdmin struct{}

func NewCmsAdmin() *CmsAdmin {
	return &CmsAdmin{}
}

// Login 登录
func (m *CmsAdmin) Login(login_key, password string, user_type int, login_type LoginType) (*model.CmsAdmin, error) {
	mdl, do := query.CmsAdminDo()
	ctx := context.Background()
	switch login_type {
	case LoginName:
		do = do.Where(mdl.IsDeleted.Is(false), mdl.UserName.Eq(login_key))
	case LoginMobile:
		do = do.Where(mdl.IsDeleted.Is(false), mdl.Mobile.Eq(login_key))
	case LoginEmail:
		do = do.Where(mdl.IsDeleted.Is(false), mdl.Email.Eq(login_key))
	case LoginNameMobile:
		do = do.Where(
			do.WithContext(ctx).Where(mdl.IsDeleted.Is(false)),
		).Where(
			do.Or(mdl.UserName.Eq(login_key)).Or(mdl.Mobile.Eq(login_key)),
		)
	case LoginNameEmail:
		do = do.Where(
			do.WithContext(ctx).Where(mdl.IsDeleted.Is(false)),
		).Where(
			do.Or(mdl.UserName.Eq(login_key)).Or(mdl.Email.Eq(login_key)),
		)
	case LoginMobileEmail:
		do = do.Where(
			do.WithContext(ctx).Where(mdl.IsDeleted.Is(false)),
		).Where(
			do.Or(mdl.Mobile.Eq(login_key)).Or(mdl.Email.Eq(login_key)),
		)
	case LoginAll:
		do = do.Where(
			do.WithContext(ctx).Where(mdl.IsDeleted.Is(false)),
		).Where(
			do.Or(mdl.UserName.Eq(login_key)).Or(mdl.Mobile.Eq(login_key)).Or(mdl.Email.Eq(login_key)),
		)
	}
	find, err := do.First()
	if err != nil {
		return nil, err
	}
	if find == nil {
		return nil, errors.New("用户不存在")
	}
	// 密码加密格式 0不加密 1默认加密 2MD5类型
	switch find.PasswordFormat {
	case 1: // 默认加密

		break
	case 2: // MD5加密
		password = xcrypto.GetMD5Hash(password)
		break
	default:
		break
	}
	logs.Debug("Login --> ", find.Password, password)
	if find.Password == password {
		return find, nil
	}
	return nil, errors.New("账号密码错误")
}

// LoginLog 写入登录日志
func (m *CmsAdmin) LoginLog(userId int64, userName string, method, path, queryx, statusCode, ip string) {
	_, logDo := query.CmsAdminLogDo()
	log := &model.CmsAdminLog{
		UserID:     userId,
		UserName:   userName,
		Method:     method,
		Path:       path,
		Query:      queryx,
		StatusCode: statusCode,
		IP:         ip,
		CreateTime: time.Now(),
	}
	logDo.Create(log)
	mdl, do := query.CmsAdminDo()
	do.Where(mdl.UserID.Eq(userId)).UpdateColumns(
		map[string]interface{}{
			mdl.LastIP.ColumnName().String():   ip,
			mdl.LastTime.ColumnName().String(): time.Now().Format("2006-01-02 15:04:05"),
		},
	)
}
