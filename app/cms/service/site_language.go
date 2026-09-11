package service

import (
	"haedu.gov.cn/cms/app/cms/domain"
	"haedu.gov.cn/cms/app/db"
)

// SiteLanguageService 站点语言服务
type SiteLanguageService struct{}

// NewSiteLanguageService 创建站点语言服务实例
func NewSiteLanguageService() *SiteLanguageService {
	return &SiteLanguageService{}
}

// SiteLanguageList 查询站点语言关联列表
// 对应 SQL:
// SELECT a.site_id, b.language_id, a.language_code, b.name, b.description, b.icon
// FROM cms_site a
// LEFT JOIN cms_language b ON a.language_code = b.code
// WHERE b.language_id IS NOT NULL
// ORDER BY b.is_default DESC
func (svc *SiteLanguageService) SiteLanguageList() ([]*domain.SiteLanguage, int64, error) {
	var list []*domain.SiteLanguage

	err := db.CmsDatabase.Table("cms_site a").
		Select("a.site_id, b.language_id, a.language_code, b.name, b.description, b.icon").
		Joins("LEFT JOIN cms_language b ON a.language_code = b.code").
		Where("b.language_id IS NOT NULL").
		Where("a.deleted = ?", false).
		Where("b.deleted = ?", false).
		Order("b.is_default DESC, a.site_id ASC").
		Scan(&list).Error

	return list, int64(len(list)), err
}

// SiteLanguageListBySiteID 根据站点ID查询语言
func (svc *SiteLanguageService) SiteLanguageListBySiteID(siteID int64) ([]*domain.SiteLanguage, error) {
	var list []*domain.SiteLanguage

	err := db.CmsDatabase.Table("cms_site a").
		Select("a.site_id, b.language_id, a.language_code, b.name, b.description, b.icon").
		Joins("LEFT JOIN cms_language b ON a.language_code = b.code").
		Where("b.language_id IS NOT NULL").
		Where("a.deleted = ?", false).
		Where("b.deleted = ?", false).
		Where("a.site_id = ?", siteID).
		Order("b.is_default DESC").
		Scan(&list).Error

	return list, err
}

// SiteLanguageListByCode 根据语言代码查询站点
func (svc *SiteLanguageService) SiteLanguageListByCode(code string) ([]*domain.SiteLanguage, error) {
	var list []*domain.SiteLanguage

	err := db.CmsDatabase.Table("cms_site a").
		Select("a.site_id, b.language_id, a.language_code, b.name, b.description, b.icon").
		Joins("LEFT JOIN cms_language b ON a.language_code = b.code").
		Where("b.language_id IS NOT NULL").
		Where("a.deleted = ?", false).
		Where("b.deleted = ?", false).
		Where("a.language_code = ?", code).
		Order("b.is_default DESC").
		Scan(&list).Error

	return list, err
}
