package service

import (
	"fmt"
	"strconv"
	"strings"

	"haedu.gov.cn/cms/app/cms/domain"
	"haedu.gov.cn/cms/app/cms/mapper"
)

type CmsSiteDomain struct{}

func NewCmsSiteDomainModel() *CmsSiteDomain {
	return &CmsSiteDomain{}
}

func (svc *CmsSiteDomain) Delete(ids string) {
	if ids == "" {
		return
	}
	domain, domainDo := mapper.CmsSiteDomainDo()
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

func (svc *CmsSiteDomain) List(siteID int64) []*domain.CmsSiteDomain {
	domain, domainDo := mapper.CmsSiteDomainDo()
	list, err := domainDo.Where(domain.SiteID.Eq(siteID)).Find()
	if err != nil {
		return nil
	}
	return list
}

func (svc *CmsSiteDomain) One(domainUrl string) *domain.CmsSiteDomain {
	domain, domainDo := mapper.CmsSiteDomainDo()
	mdl, err := domainDo.Where(domain.Domain.Eq(domainUrl)).First()
	if err != nil {
		return nil
	}
	return mdl
}
