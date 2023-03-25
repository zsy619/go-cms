package biz

import (
	"errors"
	"fmt"
	"haedu.gov.cn/cms/app/dal"
	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/app/dal/query"
	"strconv"
	"strings"
)

type CmsSite struct{}

func NewCmsSiteModel() *CmsSite {
	return &CmsSite{}
}

func (this *CmsSite) List(name string, title string, domainStr string, page int, limit int) ([]map[string]interface{}, int64, error) {
	site, siteDo := query.CmsSiteDo()
	domain, _ := query.CmsSiteDomainDo()
	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 15
	}
	db := dal.CmsDatabase.Table("cms_site").Joins("left join cms_site_domain on(cms_site.site_id=cms_site_domain.site_id and cms_site.is_deleted=0)").Select("cms_site.site_id,cms_site.name,cms_site.dir_path,cms_site.title,cms_site_domain.domain,cms_site.sort_id,cms_site.mobile,cms_site.is_default")
	if name != "" {
		db = db.Where(site.Name.Like("%" + name + "%"))
	}
	if title != "" {
		db = db.Where(site.Title.Like("%" + title + "%"))
	}
	if domainStr != "" {
		db = db.Where(domain.Domain.Like("%" + domainStr + "%"))
	}
	var list []map[string]interface{}
	var err error
	if err := db.Offset((page - 1) * limit).Limit(limit).Find(&list).Error; err != nil {
		fmt.Println(err.Error())
	}
	count, err := siteDo.Where(site.IsDeleted.Is(false)).Count()
	return list, count, err
}

func (this *CmsSite) Delete(ids string) {
	if ids == "" {
		return
	}
	site, siteDo := query.CmsSiteDo()
	idarr := strings.Split(ids, ",")
	for i := 0; i < len(idarr); i++ {
		var id = idarr[i]
		id64, _ := strconv.ParseInt(id, 0, 64)
		if id != "" {
			siteDo.Where(site.SiteID.Eq(id64)).Update(site.IsDeleted, true)
		}
	}
}

func (this *CmsSite) One(id int64) *model.CmsSite {
	site, siteDo := query.CmsSiteDo()
	list, err := siteDo.Where(site.SiteID.Eq(id)).Find()
	if err != nil {
		return nil
	}
	return list[0]
}

func (this *CmsSite) Save(mdl *model.CmsSite, domains []string, remarks []string) error {
	if mdl.Title == "" {
		return errors.New("站点名称不能为空")
	}
	if mdl.DirPath == "" {
		return errors.New("生成目录名不能为空")
	}

	site, siteDo := query.CmsSiteDo()
	count, err := siteDo.Where(site.Name.Eq(mdl.Name), site.IsDeleted.Is(false), site.SiteID.Neq(mdl.SiteID)).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("网站名称已存在,请修改后重试")
	}
	err = siteDo.Save(mdl)
	if err != nil {
		return err
	}
	domain, domainDo := query.CmsSiteDomainDo()
	_, err = domainDo.Where(domain.SiteID.Eq(mdl.SiteID)).Delete()
	if domains != nil {
		domainLen := len(domains)
		for i := 0; i < domainLen; i++ {
			if domains[i] == "" {
				continue
			}
			domain := &model.CmsSiteDomain{}
			domain.Domain = domains[i]
			domain.Remark = remarks[i]
			domain.SiteID = mdl.SiteID
			domainDo.Save(domain)
		}
	}
	return err
}
