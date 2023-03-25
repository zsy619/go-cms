package biz

import (
	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/app/dal/query"
	"strconv"
	"strings"
)

type CmsSiteDomain struct{}

func NewCmsSiteDomainModel() *CmsSiteDomain {
	return &CmsSiteDomain{}
}

func (this *CmsSiteDomain) Delete(ids string) {
	if ids == "" {
		return
	}
	domain, domainDo := query.CmsSiteDomainDo()
	idarr := strings.Split(ids, ",")
	for i := 0; i < len(idarr); i++ {
		var id = idarr[i]
		id64, _ := strconv.ParseInt(id, 0, 64)
		if id != "" {
			domainDo.Where(domain.SiteID.Eq(id64)).Delete()
		}
	}
}

func (this *CmsSiteDomain) List(siteID int64) []*model.CmsSiteDomain {
	domain, domainDo := query.CmsSiteDomainDo()
	list, err := domainDo.Where(domain.SiteID.Eq(siteID)).Find()
	if err != nil {
		return nil
	}
	return list
}
