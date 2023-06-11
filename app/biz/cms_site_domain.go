package biz

import (
	"fmt"
	"strconv"
	"strings"

	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/app/dal/query"
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
		id := idarr[i]
		id64, _ := strconv.ParseInt(id, 0, 64)
		if id != "" {
			if _, err := domainDo.Where(domain.SiteID.Eq(id64)).Delete(); err != nil {
				fmt.Println(err.Error())
			}
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
