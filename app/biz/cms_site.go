package biz

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/app/dal/query"
)

type CmsSite struct{}

func NewCmsSite() *CmsSite {
	return &CmsSite{}
}

func (this *CmsSite) SiteSaveSortId(siteId int64, sortId int32) error {
	site, siteDo := query.CmsSiteDo()
	_, err := siteDo.Where(site.SiteID.Eq(siteId)).UpdateColumns(
		map[string]interface{}{
			site.SortID.ColumnName().String():     sortId,
			site.UpdateTime.ColumnName().String(): time.Now(),
		},
	)
	return err
}

func (this *CmsSite) SitePaginate(page, limit int, name string, title string) ([]*model.CmsSite, int64, error) {
	site, siteDo := query.CmsSiteDo()
	siteDo = siteDo.Where(site.IsDeleted.Is(false))
	if name != "" {
		siteDo = siteDo.Where(site.Name.Like("%" + name + "%"))
	}
	if title != "" {
		siteDo = siteDo.Where(site.Title.Like("%" + title + "%"))
	}
	return siteDo.FindByPage((page-1)*limit, limit)
}

func (this *CmsSite) SiteDelete(ids string) {
	if ids == "" {
		return
	}
	site, siteDo := query.CmsSiteDo()
	idarr := strings.Split(ids, ",")
	for i := 0; i < len(idarr); i++ {
		id := idarr[i]
		id64, _ := strconv.ParseInt(id, 0, 64)
		if id != "" {
			siteDo.Where(site.SiteID.Eq(id64)).Update(site.IsDeleted, true)
		}
	}
}

func (this *CmsSite) SiteOne(id int64) *model.CmsSite {
	site, siteDo := query.CmsSiteDo()
	list, err := siteDo.Where(site.SiteID.Eq(id)).Find()
	if err != nil {
		return nil
	}
	return list[0]
}

func (this *CmsSite) SiteSave(mdl *model.CmsSite, domains []string, remarks []string) error {
	if mdl.Title == "" {
		return errors.New("站点名称不能为空")
	}
	if mdl.DirPath == "" {
		return errors.New("生成目录名不能为空")
	}

	site, siteDo := query.CmsSiteDo()
	if mdl.IsDefault {
		if count, _ := siteDo.Where(site.IsDefault.Is(mdl.IsDefault), site.IsDeleted.Is(false), site.SiteID.Neq(mdl.SiteID)).Count(); count > 0 {
			return errors.New("默认站点只能有一个，请修改后重试")
		}
	}
	if count, _ := siteDo.Where(site.Name.Eq(mdl.Name), site.IsDeleted.Is(false), site.SiteID.Neq(mdl.SiteID)).Count(); count > 0 {
		return errors.New("网站名称已存在，请修改后重试")
	}
	if count, _ := siteDo.Where(site.DirPath.Eq(mdl.DirPath), site.IsDeleted.Is(false), site.SiteID.Neq(mdl.SiteID)).Count(); count > 0 {
		return errors.New("生成目录名已存在，请修改后重试")
	}
	if mdl.SiteID <= 0 {
		if err := siteDo.Save(mdl); err != nil {
			return err
		}
	} else {
		if _, err := siteDo.Where(site.SiteID.Eq(mdl.SiteID)).UpdateColumns(map[string]interface{}{
			site.Name.ColumnName().String():            mdl.Name,
			site.Title.ColumnName().String():           mdl.Title,
			site.DirPath.ColumnName().String():         mdl.DirPath,
			site.IsDefault.ColumnName().String():       mdl.IsDefault,
			site.IsMobile.ColumnName().String():        mdl.IsMobile,
			site.Company.ColumnName().String():         mdl.Company,
			site.Address.ColumnName().String():         mdl.Address,
			site.Telphone.ColumnName().String():        mdl.Telphone,
			site.Fax.ColumnName().String():             mdl.Fax,
			site.Email.ColumnName().String():           mdl.Email,
			site.Crod.ColumnName().String():            mdl.Crod,
			site.HomeTitle.ColumnName().String():       mdl.HomeTitle,
			site.Copyright.ColumnName().String():       mdl.Copyright,
			site.MetaKeyword.ColumnName().String():     mdl.MetaKeyword,
			site.MetaDescription.ColumnName().String(): mdl.MetaDescription,
		}); err != nil {
			return err
		}
	}
	domain, domainDo := query.CmsSiteDomainDo()
	domainDo.Where(domain.SiteID.Eq(mdl.SiteID)).Delete()
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
	return nil
}
