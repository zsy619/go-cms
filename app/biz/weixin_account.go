package biz

import (
	"errors"
	"time"

	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/app/dal/query"
)

type WeixinAccount struct{}

func NewWeixinAccount() *WeixinAccount {
	return &WeixinAccount{}
}

// AccountPaginate 分页
func (this *WeixinAccount) AccountPaginate(page, limit int, name string, status int32) ([]*model.WeixinAccount, int64, error) {
	mdl, do := query.WeixinAccountDo()
	if name != "" {
		do = do.Where(mdl.Name.Like("%" + name + "%"))
	}
	if status >= 0 {
		do = do.Where(mdl.Status.Eq(status))
	}
	return do.Order(mdl.SortID).FindByPage((page-1)*limit, limit)
}

// AccountFind 获取
func (this *WeixinAccount) AccountFind(AccountId int64) (*model.WeixinAccount, error) {
	mdl, do := query.WeixinAccountDo()
	return do.Where(mdl.AccountID.Eq(AccountId)).First()
}

// AccountSave 保存或更新
func (this *WeixinAccount) AccountSave(input *model.WeixinAccount) error {
	mdl, do := query.WeixinAccountDo()
	if input.Name != "" {
		if count, _ := do.Where(mdl.AccountID.Neq(input.AccountID), mdl.Name.Eq(input.Name)).Count(); count > 0 {
			return errors.New("公众号名称重复")
		}
	}
	if input.OriginalID != "" {
		if count, _ := do.Where(mdl.AccountID.Neq(input.AccountID), mdl.OriginalID.Eq(input.OriginalID)).Count(); count > 0 {
			return errors.New("公众号原始ID重复")
		}
	}
	if input.WxCode != "" {
		if count, _ := do.Where(mdl.AccountID.Neq(input.AccountID), mdl.WxCode.Eq(input.WxCode)).Count(); count > 0 {
			return errors.New("公众平台微信号重复")
		}
	}
	var err error
	input.UpdateTime = time.Now()
	if input.AccountID <= 0 {
		input.CreateTime = time.Now()
		err = do.Create(input)
	} else {
		_, err = do.Where(mdl.AccountID.Eq(input.AccountID)).Updates(map[string]interface{}{
			mdl.Name.ColumnName().String():       input.Name,
			mdl.OriginalID.ColumnName().String(): input.OriginalID,
			mdl.WxCode.ColumnName().String():     input.WxCode,
			mdl.Token.ColumnName().String():      input.Token,
			mdl.AppID.ColumnName().String():      input.AppID,
			mdl.AppSecret.ColumnName().String():  input.AppSecret,
			mdl.IsPush.ColumnName().String():     input.IsPush,
			mdl.SortID.ColumnName().String():     input.SortID,
			mdl.Status.ColumnName().String():     input.Status,
			mdl.UpdateTime.ColumnName().String(): input.UpdateTime,
		})
	}
	return err
}

// AccountSaveSortId 修改排序
func (this *WeixinAccount) AccountSaveSortId(accountId int64, sortId int32) error {
	mdl, do := query.WeixinAccountDo()
	_, err := do.Where(mdl.AccountID.Eq(accountId)).UpdateColumns(
		map[string]interface{}{
			mdl.SortID.ColumnName().String():     sortId,
			mdl.UpdateTime.ColumnName().String(): time.Now(),
		},
	)
	return err
}

// AccountDestory 删除
func (this *WeixinAccount) AccountDestory(accountId int64) error {
	mdl, do := query.WeixinAccountDo()
	if _, err := do.Where(mdl.AccountID.Eq(accountId)).Delete(); err != nil {
		return err
	}
	return nil
}

// AccountChangeStatus 修改状态
func (this *WeixinAccount) AccountChangeStatus(accountId int64, status int32) error {
	mdl, do := query.WeixinAccountDo()
	_, err := do.Where(mdl.AccountID.Eq(accountId)).UpdateColumns(
		map[string]interface{}{
			mdl.Status.ColumnName().String():     status,
			mdl.UpdateTime.ColumnName().String(): time.Now(),
		},
	)
	return err
}
