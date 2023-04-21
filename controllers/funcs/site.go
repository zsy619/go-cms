package funcs

import (
	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/biz/bizmodel"
)

func SiteDefault() *bizmodel.ApiSiteModel {
	find, err := biz.NewApiSite().Default()
	if err != nil {
		return &bizmodel.ApiSiteModel{}
	}
	return find
}
