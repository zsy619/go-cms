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
	if !find.Enabled {
		return nil, errors.New("用户已禁用")
	}
	// 密码加密格式 0不加密 1默认加密 2MD5类型
	switch find.PasswordFormat {
	case 1: // 默认加密

		break
	case 2: // MD5加密
		password = xcrypto.GetMD5Hash(password + find.PasswordSalt)
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

func (m *CmsAdmin) AdminPaginate(page, limit int, roleId int64, realName, userName string) ([]*model.CmsAdmin, int64, error) {
	mdl, do := query.CmsAdminDo()
	if roleId > 0 {
		do = do.Where(mdl.RoleID.Eq(roleId))
	}
	if realName != "" {
		do = do.Where(mdl.RealName.Like("%" + realName + "%"))
	}
	if userName != "" {
		do = do.Where(mdl.UserName.Like("%" + userName + "%"))
	}
	return do.Order(mdl.SortID).FindByPage((page-1)*limit, limit)
}

// AdminDestory 删除
func (this *CmsAdmin) AdminDestory(userId int64) error {
	mdl, do := query.CmsAdminDo()
	if _, err := do.Where(mdl.UserID.Eq(userId)).Delete(); err != nil {
		return err
	}
	return nil
}

// AdminSave 保存或更新
func (this *CmsAdmin) AdminSave(input *model.CmsAdmin) error {
	mdl, do := query.CmsAdminDo()
	if input.NickName != "" {
		if count, _ := do.Where(mdl.UserID.Neq(input.UserID), mdl.NickName.Eq(input.NickName)).Count(); count > 0 {
			return errors.New("昵称重复")
		}
	}
	if input.UserName != "" {
		if count, _ := do.Where(mdl.UserID.Neq(input.UserID), mdl.UserName.Eq(input.UserName)).Count(); count > 0 {
			return errors.New("登录名重复")
		}
	}
	if input.Email != "" {
		if count, _ := do.Where(mdl.UserID.Neq(input.UserID), mdl.Email.Eq(input.Email)).Count(); count > 0 {
			return errors.New("邮箱重复")
		}
	}
	if input.Mobile != "" {
		if count, _ := do.Where(mdl.UserID.Neq(input.UserID), mdl.Mobile.Eq(input.Mobile)).Count(); count > 0 {
			return errors.New("手机号重复")
		}
	}
	if input.IDCard != "" {
		if count, _ := do.Where(mdl.UserID.Neq(input.UserID), mdl.IDCard.Eq(input.IDCard)).Count(); count > 0 {
			return errors.New("身份证号重复")
		}
	}
	if input.UserNumber != "" {
		if count, _ := do.Where(mdl.UserID.Neq(input.UserID), mdl.UserNumber.Eq(input.UserNumber)).Count(); count > 0 {
			return errors.New("特征标识重复")
		}
	}
	roleMdl, roleDo := query.CmsAdminRoleDo()
	var roleType string
	if err := roleDo.Where(roleMdl.RoleID.Eq(input.RoleID)).Pluck(roleMdl.Type, &roleType); err != nil {
		return err
	}
	input.RoleType = roleType
	var err error
	input.UpdateTime = time.Now()
	if input.UserID <= 0 {
		input.CreateTime = time.Now()
		err = do.Create(input)
	} else {
		_, err = do.Where(mdl.UserID.Eq(input.UserID)).Updates(map[string]interface{}{
			mdl.RoleID.ColumnName().String():     input.RoleID,
			mdl.RoleType.ColumnName().String():   input.RoleType,
			mdl.RealName.ColumnName().String():   input.RealName,
			mdl.NickName.ColumnName().String():   input.NickName,
			mdl.IDCard.ColumnName().String():     input.IDCard,
			mdl.Sex.ColumnName().String():        input.Sex,
			mdl.UserName.ColumnName().String():   input.UserName,
			mdl.UserNumber.ColumnName().String(): input.UserNumber,
			mdl.Email.ColumnName().String():      input.Email,
			mdl.Telphone.ColumnName().String():   input.Telphone,
			mdl.Mobile.ColumnName().String():     input.Mobile,
			mdl.MobilePin.ColumnName().String():  input.MobilePin,
			mdl.UserType.ColumnName().String():   input.UserType,
			mdl.Enabled.ColumnName().String():    input.Enabled,
			mdl.Remark.ColumnName().String():     input.Remark,
			mdl.Avatar.ColumnName().String():     input.Avatar,
			mdl.IsAudit.ColumnName().String():    input.IsAudit,
			mdl.SortID.ColumnName().String():     input.SortID,
			mdl.UpdateTime.ColumnName().String(): input.UpdateTime,
		})
		if input.Password != "" {
			_, err = do.Where(mdl.RoleID.Eq(input.RoleID)).Updates(map[string]interface{}{
				mdl.Password.ColumnName().String():       input.Password,
				mdl.PasswordSalt.ColumnName().String():   input.PasswordSalt,
				mdl.PasswordFormat.ColumnName().String(): input.PasswordFormat,
			})
		}
	}
	return err
}

func (this *CmsAdmin) AdminSaveSortId(userId int64, sortId int32) error {
	mdl, do := query.CmsAdminDo()
	_, err := do.Where(mdl.UserID.Eq(userId)).UpdateColumns(
		map[string]interface{}{
			mdl.SortID.ColumnName().String():     sortId,
			mdl.UpdateTime.ColumnName().String(): time.Now(),
		},
	)
	return err
}

func (m *CmsAdmin) AdminFind(userId int64) (*model.CmsAdmin, error) {
	mdl, do := query.CmsAdminDo()
	return do.Where(mdl.UserID.Eq(userId)).First()
}

func (m *CmsAdmin) LogPaginate(page, limit int, userId int64, userName string) ([]*model.CmsAdminLog, int64, error) {
	mdl, do := query.CmsAdminLogDo()
	if userId > 0 {
		do = do.Where(mdl.UserID.Eq(userId))
	}
	return do.Where(mdl.UserName.Like("%"+userName+"%")).FindByPage((page-1)*limit, limit)
}

func (m *CmsAdmin) RolePaginate(page, limit int, name string) ([]*model.CmsAdminRole, int64, error) {
	mdl, do := query.CmsAdminRoleDo()
	return do.Where(mdl.Name.Like("%"+name+"%")).Order(mdl.SortID).FindByPage((page-1)*limit, limit)
}

func (m *CmsAdmin) RoleFind(roleId int64) (*model.CmsAdminRole, error) {
	mdl, do := query.CmsAdminRoleDo()
	return do.Where(mdl.RoleID.Eq(roleId)).First()
}

// RoleSave 保存或更新
func (this *CmsAdmin) RoleSave(input *model.CmsAdminRole) error {
	mdl, do := query.CmsAdminRoleDo()
	if input.Name != "" {
		if count, _ := do.Where(mdl.RoleID.Neq(input.RoleID), mdl.Name.Eq(input.Name)).Count(); count > 0 {
			return errors.New("调用别名重复")
		}
	}
	var err error
	input.UpdateTime = time.Now()
	if input.RoleID <= 0 {
		input.CreateTime = time.Now()
		err = do.Create(input)
	} else {
		_, err = do.Where(mdl.RoleID.Eq(input.RoleID)).Updates(map[string]interface{}{
			mdl.Name.ColumnName().String():       input.Name,
			mdl.Type.ColumnName().String():       input.Type,
			mdl.Remark.ColumnName().String():     input.Remark,
			mdl.SortID.ColumnName().String():     input.SortID,
			mdl.UpdateTime.ColumnName().String(): input.UpdateTime,
		})
	}
	return err
}

func (this *CmsAdmin) RoleValSave(roleId int64, values map[string]string) error {
	mdl, do := query.CmsAdminRoleValueDo()
	_, err := do.Where(mdl.RoleID.Eq(roleId)).Delete()
	for k, v := range values {
		err = do.Save(&model.CmsAdminRoleValue{
			RoleID:  roleId,
			NavName: k,
			Action:  v,
		})
	}
	return err
}

// RoleDestory 删除
func (this *CmsAdmin) RoleDestory(roleId int64) error {
	adminMdl, adminDo := query.CmsAdminDo()
	if count, _ := adminDo.Where(adminMdl.RoleID.Eq(roleId)).Count(); count > 0 {
		return errors.New("该角色下有用户，无法删除")
	}
	valueMdl, valueDo := query.CmsAdminRoleValueDo()
	if _, err := valueDo.Where(valueMdl.RoleID.Eq(roleId)).Delete(); err != nil {
		return err
	}
	mdl, do := query.CmsAdminRoleDo()
	if _, err := do.Where(mdl.RoleID.Eq(roleId)).Delete(); err != nil {
		return err
	}
	return nil
}

func (this *CmsAdmin) RoleSaveSortId(roleId int64, sortId int32) error {
	mdl, do := query.CmsAdminRoleDo()
	_, err := do.Where(mdl.RoleID.Eq(roleId)).UpdateColumns(
		map[string]interface{}{
			mdl.SortID.ColumnName().String():     sortId,
			mdl.UpdateTime.ColumnName().String(): time.Now(),
		},
	)
	return err
}

func (this *CmsAdmin) RoleValueFind(roleId int64) ([]*model.CmsAdminRoleValue, int64, error) {
	mdl, do := query.CmsAdminRoleValueDo()
	return do.Where(mdl.RoleID.Eq(roleId)).FindByPage(0, 99999)
}

func (this *CmsAdmin) NavFind(roleId int64) ([]*model.CmsAdminNav, int64, error) {
	mdl, do := query.CmsAdminNavDo()
	return do.Order(mdl.SortID).FindByPage(0, 99999)
}

func (this *CmsAdmin) OneByUserId(userId int64) *model.CmsAdmin {
	mdl, do := query.CmsAdminDo()
	first, err := do.Where(mdl.UserID.Eq(userId)).First()
	if err != nil {
		return nil
	}
	return first
}
