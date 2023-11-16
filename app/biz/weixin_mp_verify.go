package biz

import (
	"errors"
	"fmt"
	"time"

	"github.com/beego/beego/v2/core/logs"

	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/app/dal/query"
	"haedu.gov.cn/cms/app/lib"
)

type WeixinMpVerify struct{}

func NewWeixinMpVerify() *WeixinMpVerify {
	return &WeixinMpVerify{}
}

func (svc *WeixinMpVerify) Get(accountId int64) ([]*model.WeixinMpVerify, error) {
	mdl, do := query.WeixinMpVerifyDo()
	return do.Where(mdl.AccountID.Eq(accountId)).Order(mdl.SortID).Find()
}

/**
 * @description: 获取缓存
 * @return {*}
 */
func (svc *WeixinMpVerify) GetCache() ([]*model.WeixinMpVerify, error) {
	cacheKey := fmt.Sprintf("%s_%d", "Weixin_Mp_Verify", 0)
	/*if found, item := ApiCache.Get(cacheKey); found {
		list := item.([]*model.WeixinMpVerify)
		logs.Debug("GetCache[Cache]::", "cacheKey", cacheKey, "WeixinMpVerify", list)
		return list, nil
	}*/
	if found, item := lib.WechatVerifyCache.Get(cacheKey); found {
		list := item.([]*model.WeixinMpVerify)
		logs.Debug("WechatVerifyList[Cache]::", "cacheKey", cacheKey, "WechatVerifyList", list)
		return list, nil
	}
	mdl, do := query.WeixinMpVerifyDo()
	list, err := do.Order(mdl.SortID).Find()
	if len(list) > 0 {
		// ApiCache.Set(cacheKey, list, 60*60*24)
		lib.WechatVerifyCache.Set(cacheKey, list)
	}
	return list, err
}

/**
 * @description: 刷新缓存
 * @return {*}
 */
func (svc *WeixinMpVerify) RefeshCache() {
	cacheKey := fmt.Sprintf("%s_%d", "Weixin_Mp_Verify", 0)
	ApiCache.Delete(cacheKey)
	_, _ = svc.GetCache()
}

func (svc *WeixinMpVerify) Find(verifyId int64) (*model.WeixinMpVerify, error) {
	mdl, do := query.WeixinMpVerifyDo()
	return do.Where(mdl.VerifyID.Eq(verifyId)).First()
}

func (svc *WeixinMpVerify) Save(m *model.WeixinMpVerify) error {
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

func (svc *WeixinMpVerify) Destory(verifyId int64) error {
	mdl, do := query.WeixinMpVerifyDo()
	_, err := do.Where(mdl.VerifyID.Eq(verifyId)).Delete()
	return err
}

func (svc *WeixinMpVerify) SaveSortId(verifyId int64, sortId int32) error {
	mdl, do := query.WeixinMpVerifyDo()
	_, err := do.Where(mdl.VerifyID.Eq(verifyId)).UpdateColumns(
		map[string]interface{}{
			mdl.SortID.ColumnName().String():     sortId,
			mdl.UpdateTime.ColumnName().String(): time.Now(),
		},
	)
	return err
}
