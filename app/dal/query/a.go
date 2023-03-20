package query

import (
	"context"

	"haedu.gov.cn/cms/app/dal"
)

var defaultContext = context.Background()

// CmsAdminDo
func CmsAdminDo() (cmsAdmin, *cmsAdminDo) {
	u := Use(dal.CmsDatabase.DB).CmsAdmin
	return u, u.WithContext(defaultContext).Debug()
}

// CmsAdminLogDo
func CmsAdminLogDo() (cmsAdminLog, *cmsAdminLogDo) {
	u := Use(dal.CmsDatabase.DB).CmsAdminLog
	return u, u.WithContext(defaultContext).Debug()
}

// CmsSiteDo
func CmsSiteDo() (cmsSite, *cmsSiteDo) {
	u := Use(dal.CmsDatabase.DB).CmsSite
	return u, u.WithContext(defaultContext).Debug()
}

// CmsSiteDomainDo
func CmsSiteDomainDo() (cmsSiteDomain, *cmsSiteDomainDo) {
	u := Use(dal.CmsDatabase.DB).CmsSiteDomain
	return u, u.WithContext(defaultContext).Debug()
}

// CmsSiteChannelDo
func CmsSiteChannelDo() (cmsSiteChannel, *cmsSiteChannelDo) {
	u := Use(dal.CmsDatabase.DB).CmsSiteChannel
	return u, u.WithContext(defaultContext).Debug()
}

// CmsArticleCategoryDo
func CmsArticleCategoryDo() (cmsArticleCategory, *cmsArticleCategoryDo) {
	u := Use(dal.CmsDatabase.DB).CmsArticleCategory
	return u, u.WithContext(defaultContext).Debug()
}

// CmsArticleCategoryRelationDo
func CmsArticleCategoryRelationDo() (cmsArticleCategoryRelation, *cmsArticleCategoryRelationDo) {
	u := Use(dal.CmsDatabase.DB).CmsArticleCategoryRelation
	return u, u.WithContext(defaultContext).Debug()
}

// CmsArticleDo
func CmsArticleDo() (cmsArticle, *cmsArticleDo) {
	u := Use(dal.CmsDatabase.DB).CmsArticle
	return u, u.WithContext(defaultContext).Debug()
}

// CmsArticleAlbumDo
func CmsArticleAlbumDo() (cmsArticleAlbum, *cmsArticleAlbumDo) {
	u := Use(dal.CmsDatabase.DB).CmsArticleAlbum
	return u, u.WithContext(defaultContext).Debug()
}

// CmsArticleAttachDo
func CmsArticleAttachDo() (cmsArticleAttach, *cmsArticleAttachDo) {
	u := Use(dal.CmsDatabase.DB).CmsArticleAttach
	return u, u.WithContext(defaultContext).Debug()
}

// CmsArticleLabelDo
func CmsArticleLabelDo() (cmsArticleLabel, *cmsArticleLabelDo) {
	u := Use(dal.CmsDatabase.DB).CmsArticleLabel
	return u, u.WithContext(defaultContext).Debug()
}

// CmsArticleLabelRelationDo
func CmsArticleLabelRelationDo() (cmsArticleLabelRelation, *cmsArticleLabelRelationDo) {
	u := Use(dal.CmsDatabase.DB).CmsArticleLabelRelation
	return u, u.WithContext(defaultContext).Debug()
}

// CmsLinkCategoryDo
func CmsLinkCategoryDo() (cmsLinkCategory, *cmsLinkCategoryDo) {
	u := Use(dal.CmsDatabase.DB).CmsLinkCategory
	return u, u.WithContext(defaultContext).Debug()
}

// CmsLinkCategoryRelationDo
func CmsLinkCategoryRelationDo() (cmsLinkCategoryRelation, *cmsLinkCategoryRelationDo) {
	u := Use(dal.CmsDatabase.DB).CmsLinkCategoryRelation
	return u, u.WithContext(defaultContext).Debug()
}

// CmsLinkDo
func CmsLinkDo() (cmsLink, *cmsLinkDo) {
	u := Use(dal.CmsDatabase.DB).CmsLink
	return u, u.WithContext(defaultContext).Debug()
}
