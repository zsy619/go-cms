package biz

import (
	"errors"
	"fmt"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/app/dal/query"
)

type WeixinMpVerify struct{}

func NewWeixinMpVerify() *WeixinMpVerify {
	return &WeixinMpVerify{}
}

func (this *WeixinMpVerify) Get(accountId int64) ([]*model.WeixinMpVerify, error) {
	mdl, do := query.WeixinMpVerifyDo()
	return do.Where(mdl.AccountID.Eq(accountId)).Order(mdl.SortID).Find()
}

func (this *WeixinMpVerify) GetStatus(status int32) ([]*model.WeixinMpVerify, error) {
	cacheKey := fmt.Sprintf("%s_%d", "Mp_Verify", status)
	if found, item := ApiCache.Get(cacheKey); found {
		list := item.([]*model.WeixinMpVerify)
		logs.Debug("GetStatus[Cache]::", "cacheKey", cacheKey, "MpVerify", list)
		return list, nil
	}
	mdl, do := query.WeixinMpVerifyDo()
	list, err := do.Where(mdl.Status.Eq(status)).Order(mdl.SortID).Find()
	if list != nil && len(list) > 0 {
		ApiCache.Set(cacheKey, list, 60*60*24)
	}
	return list, err
}

func (this *WeixinMpVerify) Find(verifyId int64) (*model.WeixinMpVerify, error) {
	mdl, do := query.WeixinMpVerifyDo()
	return do.Where(mdl.VerifyID.Eq(verifyId)).First()
}

func (this *WeixinMpVerify) Save(m *model.WeixinMpVerify) error {
	if m == nil || m.AccountID <= 0 {
		return errors.New("参数错误")
	}
	mdl, do := query.WeixinMpVerifyDo()
	m.UpdateTime = time.Now()
	if m.VerifyID <= 0 {
		m.CreateTime = time.Now()
		return do.Create(m)
	} else {
		_, err := do.Where(mdl.VerifyID.Eq(m.VerifyID)).UpdateColumns(map[string]interface{}{
			mdl.AccountID.ColumnName().String():  m.AccountID,
			mdl.Path.ColumnName().String():       m.Path,
			mdl.FilePath.ColumnName().String():   m.FilePath,
			mdl.FileName.ColumnName().String():   m.FileName,
			mdl.SortID.ColumnName().String():     m.SortID,
			mdl.Status.ColumnName().String():     m.Status,
			mdl.UpdateID.ColumnName().String():   m.UpdateID,
			mdl.UpdateTime.ColumnName().String(): m.UpdateTime,
		})
		return err
	}
}

func (this *WeixinMpVerify) Destory(verifyId int64) error {
	mdl, do := query.WeixinMpVerifyDo()
	_, err := do.Where(mdl.VerifyID.Eq(verifyId)).Delete()
	return err
}

func (this *WeixinMpVerify) SaveSortId(verifyId int64, sortId int32) error {
	mdl, do := query.WeixinMpVerifyDo()
	_, err := do.Where(mdl.VerifyID.Eq(verifyId)).UpdateColumns(
		map[string]interface{}{
			mdl.SortID.ColumnName().String():     sortId,
			mdl.UpdateTime.ColumnName().String(): time.Now(),
		},
	)
	return err
}
