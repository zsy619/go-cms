package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/beego/beego/v2/core/logs"

	"haedu.gov.cn/cms/app/cms/domain"
	"haedu.gov.cn/cms/app/cms/mapper"
	lib "haedu.gov.cn/cms/app/tool"
)

type WeixinAccount struct{}

func NewWeixinAccount() *WeixinAccount {
	return &WeixinAccount{}
}

// AccountPaginate 分页
func (svc *WeixinAccount) AccountPaginate(page, limit int, name string, status int32) ([]*domain.WeixinAccount, int64, error) {
	mdl, do := mapper.WeixinAccountDo()
	if name != "" {
		do = do.Where(mdl.Name.Like("%" + name + "%"))
	}
	if status >= 0 {
		do = do.Where(mdl.Status.Eq(status))
	}
	return do.Order(mdl.SortID).FindByPage((page-1)*limit, limit)
}

// AccountFind 获取
func (svc *WeixinAccount) AccountFind(accountId int64) (*domain.WeixinAccount, error) {
	mdl, do := mapper.WeixinAccountDo()
	return do.Where(mdl.AccountID.Eq(accountId)).First()
}

func (svc *WeixinAccount) AccountFindCache(accountId int64) (*domain.WeixinAccount, error) {
	cacheKey := fmt.Sprintf("AccountFindCache_%d", accountId)
	/*if ok, v := WeiXinCache.Get(cacheKey); ok {
		return v.(*model.WeixinAccount), nil
	}*/
	if found, item := lib.WechatAccountCache.Get(cacheKey); found {
		mdl := item.(*domain.WeixinAccount)
		logs.Debug("WechatAccountList[Cache]::", "cacheKey", cacheKey, "WechatAccountList", mdl)
		return mdl, nil
	}
	mdl, do := mapper.WeixinAccountDo()
	find, err := do.Where(mdl.AccountID.Eq(accountId)).First()
	if err == nil {
		// WeiXinCache.Set(cacheKey, find, 7100)
		lib.WechatAccountCache.Set(cacheKey, find)
		return find, nil
	}
	return nil, err
}

// AccountSave 保存或更新
func (svc *WeixinAccount) AccountSave(input *domain.WeixinAccount) error {
	mdl, do := mapper.WeixinAccountDo()
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
			mdl.AppAesKey.ColumnName().String():  input.AppAesKey,
			mdl.AppID.ColumnName().String():      input.AppID,
			mdl.AppSecret.ColumnName().String():  input.AppSecret,
			mdl.IsPush.ColumnName().String():     input.IsPush,
			mdl.SortID.ColumnName().String():     input.SortID,
			mdl.Status.ColumnName().String():     input.Status,
			mdl.UpdateTime.ColumnName().String(): input.UpdateTime,
			mdl.UpdateID.ColumnName().String():   input.UpdateID,
			mdl.UpdateName.ColumnName().String(): input.UpdateName,
		})
	}
	return err
}

// AccountSaveSortId 修改排序
func (svc *WeixinAccount) AccountSaveSortId(accountId int64, sortId int32) error {
	mdl, do := mapper.WeixinAccountDo()
	_, err := do.Where(mdl.AccountID.Eq(accountId)).UpdateColumns(
		map[string]interface{}{
			mdl.SortID.ColumnName().String():     sortId,
			mdl.UpdateTime.ColumnName().String(): time.Now(),
		},
	)
	return err
}

// AccountDestory 删除
func (svc *WeixinAccount) AccountDestory(accountId int64) error {
	mdl, do := mapper.WeixinAccountDo()
	if _, err := do.Where(mdl.AccountID.Eq(accountId)).Delete(); err != nil {
		return err
	}
	return nil
}

// AccountChangeStatus 修改状态
func (svc *WeixinAccount) AccountChangeStatus(accountId int64, status int32) error {
	mdl, do := mapper.WeixinAccountDo()
	_, err := do.Where(mdl.AccountID.Eq(accountId)).UpdateColumns(
		map[string]interface{}{
			mdl.Status.ColumnName().String():     status,
			mdl.UpdateTime.ColumnName().String(): time.Now(),
		},
	)
	return err
}
