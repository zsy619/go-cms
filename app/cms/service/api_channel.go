package service

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/beego/beego/v2/core/logs"
	"github.com/zsy619/tools/xgeneric"

	"haedu.gov.cn/cms/app/cms/mapper"
	service_model "haedu.gov.cn/cms/app/cms/service/model"
	lib "haedu.gov.cn/cms/app/tool"
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
func (svc *ApiChannel) Find(name string, channel_id int64) (*service_model.ApiChannelModel, error) {
	if name == "" && channel_id <= 0 {
		return nil, errors.New("参数错误")
	}
	cacheKey := fmt.Sprintf("ApiChannel_Find_%s_%d", name, channel_id)
	/*if found, item := ApiCache.Get(cacheKey); found {
		logs.Debug("ApiCache")
		return item.(*bizmodel.ApiChannelModel), nil
	}*/
	if found, item := lib.ArticleCache.Get(cacheKey); found {
		mdl := item.(*service_model.ApiChannelModel)
		logs.Debug("ArticleList[Cache]::", "cacheKey", cacheKey, "ArticleList", mdl)
		return mdl, nil
	}
	_, do := mapper.CmsSiteChannelDo()
	find := &service_model.ApiChannelModel{}
	field := `a.channel_id,a.parent_id,a.title,a.name,a.kind,a.class_layer,a.link_url,a.img_url1,a.img_url2,a.sort_id,a.is_album,a.is_attach,a.is_spec,a.tmpl_chnl,a.tmpl_cat,a.tmpl_lst,a.tmpl_dtl` +
		`,b.flag as site_flag`
	sql := `SELECT ` + field + ` FROM cms_site_channel a` +
		` LEFT JOIN cms_site b ON b.site_id=a.site_id` +
		` WHERE 1=1 ` +
		xgeneric.IFF(channel_id > 0, ` AND a.channel_id=`+strconv.FormatInt(channel_id, 10), ``) +
		xgeneric.IFF(name != "", ` AND a.name='`+name+`'`, ``)
	err := do.Debug().UnderlyingDB().Raw(sql).Scan(&find).Error
	if err != nil {
		return nil, err
	}
	if find == nil || find.Name == "" {
		return nil, errors.New("频道不存在")
	}
	// ApiCache.Set(cacheKey, find, 3600)
	lib.ChannelCache.Set(cacheKey, find)
	return find, nil
}
