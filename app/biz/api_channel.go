package biz

import (
	"errors"
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/cms/app/biz/bizmodel"
	"haedu.gov.cn/cms/app/dal/query"
)

type ApiChannel struct{}

func NewApiChannel() *ApiChannel {
	return &ApiChannel{}
}

/**
 * @description: Find 频道获取
 * @param {string} name 频道名称
 * @param {int64} channel_id 频道ID
 * @return {*}
 */
func (this *ApiChannel) Find(name string, channel_id int64) (*bizmodel.ApiChannelFindModel, error) {
	if name == "" && channel_id <= 0 {
		return nil, errors.New("参数错误")
	}
	cacheKey := fmt.Sprintf("ApiChannel_Find_%s_%d", name, channel_id)
	if found, item := ApiCache.Get(cacheKey); found {
		logs.Debug("ApiCache")
		return item.(*bizmodel.ApiChannelFindModel), nil
	}
	mdl, do := query.CmsSiteChannelDo()
	if name != "" {
		do = do.Where(mdl.Name.Eq(name))
	} else if channel_id > 0 {
		do = do.Where(mdl.ChannelID.Eq(channel_id))
	}
	find := &bizmodel.ApiChannelFindModel{}
	err := do.Select(mdl.ChannelID, mdl.ParentID, mdl.Title, mdl.Name, mdl.Kind, mdl.ClassLayer, mdl.ImgUrl1, mdl.ImgUrl2, mdl.SortID, mdl.IsAlbum, mdl.IsAttach, mdl.IsSpec, mdl.Template).Order(mdl.SortID).Scan(&find)
	if err != nil {
		return nil, err
	}
	if find == nil || find.Name == "" {
		return nil, errors.New("频道不存在")
	}
	ApiCache.Set(cacheKey, find, 3600)
	return find, nil
}
