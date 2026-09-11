package mapper

import (
	"context"

	"haedu.gov.cn/cms/app/db"
)

var defaultContext = context.Background()

// Do() helpers for all 39 entities - one CmsXxxDo() function per entity
// CmsAdCategoryRelationDo returns the entity model and its Data Object query builder
func CmsAdCategoryRelationDo() (cmsAdCategoryRelation, *cmsAdCategoryRelationDo) {
	u := Use(db.CmsDatabase.DB).CmsAdCategoryRelation
	return u, u.WithContext(defaultContext).Debug()
}

// CmsAdminDo returns the entity model and its Data Object query builder
func CmsAdminDo() (cmsAdmin, *cmsAdminDo) {
	u := Use(db.CmsDatabase.DB).CmsAdmin
	return u, u.WithContext(defaultContext).Debug()
}

// CmsAdminLogDo returns the entity model and its Data Object query builder
func CmsAdminLogDo() (cmsAdminLog, *cmsAdminLogDo) {
	u := Use(db.CmsDatabase.DB).CmsAdminLog
	return u, u.WithContext(defaultContext).Debug()
}

// CmsAdminNavDo returns the entity model and its Data Object query builder
func CmsAdminNavDo() (cmsAdminNav, *cmsAdminNavDo) {
	u := Use(db.CmsDatabase.DB).CmsAdminNav
	return u, u.WithContext(defaultContext).Debug()
}

// CmsAdminNoticeDo returns the entity model and its Data Object query builder
func CmsAdminNoticeDo() (cmsAdminNotice, *cmsAdminNoticeDo) {
	u := Use(db.CmsDatabase.DB).CmsAdminNotice
	return u, u.WithContext(defaultContext).Debug()
}

// CmsAdminRoleDo returns the entity model and its Data Object query builder
func CmsAdminRoleDo() (cmsAdminRole, *cmsAdminRoleDo) {
	u := Use(db.CmsDatabase.DB).CmsAdminRole
	return u, u.WithContext(defaultContext).Debug()
}

// CmsAdminRoleSiteDo returns the entity model and its Data Object query builder
func CmsAdminRoleSiteDo() (cmsAdminRoleSite, *cmsAdminRoleSiteDo) {
	u := Use(db.CmsDatabase.DB).CmsAdminRoleSite
	return u, u.WithContext(defaultContext).Debug()
}

// CmsAdminRoleValueDo returns the entity model and its Data Object query builder
func CmsAdminRoleValueDo() (cmsAdminRoleValue, *cmsAdminRoleValueDo) {
	u := Use(db.CmsDatabase.DB).CmsAdminRoleValue
	return u, u.WithContext(defaultContext).Debug()
}

// CmsAdsDo returns the entity model and its Data Object query builder
func CmsAdsDo() (cmsAds, *cmsAdsDo) {
	u := Use(db.CmsDatabase.DB).CmsAds
	return u, u.WithContext(defaultContext).Debug()
}

// CmsAdsCategoryDo returns the entity model and its Data Object query builder
func CmsAdsCategoryDo() (cmsAdsCategory, *cmsAdsCategoryDo) {
	u := Use(db.CmsDatabase.DB).CmsAdsCategory
	return u, u.WithContext(defaultContext).Debug()
}

// CmsAdsCategoryRelationDo returns the entity model and its Data Object query builder
func CmsAdsCategoryRelationDo() (cmsAdsCategoryRelation, *cmsAdsCategoryRelationDo) {
	u := Use(db.CmsDatabase.DB).CmsAdsCategoryRelation
	return u, u.WithContext(defaultContext).Debug()
}

// CmsAlbumDo returns the entity model and its Data Object query builder
func CmsAlbumDo() (cmsAlbum, *cmsAlbumDo) {
	u := Use(db.CmsDatabase.DB).CmsAlbum
	return u, u.WithContext(defaultContext).Debug()
}

// CmsArticleDo returns the entity model and its Data Object query builder
func CmsArticleDo() (cmsArticle, *cmsArticleDo) {
	u := Use(db.CmsDatabase.DB).CmsArticle
	return u, u.WithContext(defaultContext).Debug()
}

// CmsArticleCategoryDo returns the entity model and its Data Object query builder
func CmsArticleCategoryDo() (cmsArticleCategory, *cmsArticleCategoryDo) {
	u := Use(db.CmsDatabase.DB).CmsArticleCategory
	return u, u.WithContext(defaultContext).Debug()
}

// CmsArticleCategoryRelationDo returns the entity model and its Data Object query builder
func CmsArticleCategoryRelationDo() (cmsArticleCategoryRelation, *cmsArticleCategoryRelationDo) {
	u := Use(db.CmsDatabase.DB).CmsArticleCategoryRelation
	return u, u.WithContext(defaultContext).Debug()
}

// CmsArticleCommentDo returns the entity model and its Data Object query builder
func CmsArticleCommentDo() (cmsArticleComment, *cmsArticleCommentDo) {
	u := Use(db.CmsDatabase.DB).CmsArticleComment
	return u, u.WithContext(defaultContext).Debug()
}

// CmsArticleLabelDo returns the entity model and its Data Object query builder
func CmsArticleLabelDo() (cmsArticleLabel, *cmsArticleLabelDo) {
	u := Use(db.CmsDatabase.DB).CmsArticleLabel
	return u, u.WithContext(defaultContext).Debug()
}

// CmsArticleLabelRelationDo returns the entity model and its Data Object query builder
func CmsArticleLabelRelationDo() (cmsArticleLabelRelation, *cmsArticleLabelRelationDo) {
	u := Use(db.CmsDatabase.DB).CmsArticleLabelRelation
	return u, u.WithContext(defaultContext).Debug()
}

// CmsArticlePropertyDo returns the entity model and its Data Object query builder
func CmsArticlePropertyDo() (cmsArticleProperty, *cmsArticlePropertyDo) {
	u := Use(db.CmsDatabase.DB).CmsArticleProperty
	return u, u.WithContext(defaultContext).Debug()
}

// CmsAttachDo returns the entity model and its Data Object query builder
func CmsAttachDo() (cmsAttach, *cmsAttachDo) {
	u := Use(db.CmsDatabase.DB).CmsAttach
	return u, u.WithContext(defaultContext).Debug()
}

// CmsLanguageDo returns the entity model and its Data Object query builder
func CmsLanguageDo() (cmsLanguage, *cmsLanguageDo) {
	u := Use(db.CmsDatabase.DB).CmsLanguage
	return u, u.WithContext(defaultContext).Debug()
}

// CmsLinkDo returns the entity model and its Data Object query builder
func CmsLinkDo() (cmsLink, *cmsLinkDo) {
	u := Use(db.CmsDatabase.DB).CmsLink
	return u, u.WithContext(defaultContext).Debug()
}

// CmsLinkCategoryDo returns the entity model and its Data Object query builder
func CmsLinkCategoryDo() (cmsLinkCategory, *cmsLinkCategoryDo) {
	u := Use(db.CmsDatabase.DB).CmsLinkCategory
	return u, u.WithContext(defaultContext).Debug()
}

// CmsLinkCategoryRelationDo returns the entity model and its Data Object query builder
func CmsLinkCategoryRelationDo() (cmsLinkCategoryRelation, *cmsLinkCategoryRelationDo) {
	u := Use(db.CmsDatabase.DB).CmsLinkCategoryRelation
	return u, u.WithContext(defaultContext).Debug()
}

// CmsSiteDo returns the entity model and its Data Object query builder
func CmsSiteDo() (cmsSite, *cmsSiteDo) {
	u := Use(db.CmsDatabase.DB).CmsSite
	return u, u.WithContext(defaultContext).Debug()
}

// CmsSiteChannelDo returns the entity model and its Data Object query builder
func CmsSiteChannelDo() (cmsSiteChannel, *cmsSiteChannelDo) {
	u := Use(db.CmsDatabase.DB).CmsSiteChannel
	return u, u.WithContext(defaultContext).Debug()
}

// CmsSiteChannelAlbumDo returns the entity model and its Data Object query builder
func CmsSiteChannelAlbumDo() (cmsSiteChannelAlbum, *cmsSiteChannelAlbumDo) {
	u := Use(db.CmsDatabase.DB).CmsSiteChannelAlbum
	return u, u.WithContext(defaultContext).Debug()
}

// CmsSiteChannelFieldDo returns the entity model and its Data Object query builder
func CmsSiteChannelFieldDo() (cmsSiteChannelField, *cmsSiteChannelFieldDo) {
	u := Use(db.CmsDatabase.DB).CmsSiteChannelField
	return u, u.WithContext(defaultContext).Debug()
}

// CmsSiteDomainDo returns the entity model and its Data Object query builder
func CmsSiteDomainDo() (cmsSiteDomain, *cmsSiteDomainDo) {
	u := Use(db.CmsDatabase.DB).CmsSiteDomain
	return u, u.WithContext(defaultContext).Debug()
}

// CmsTagDo returns the entity model and its Data Object query builder
func CmsTagDo() (cmsTag, *cmsTagDo) {
	u := Use(db.CmsDatabase.DB).CmsTag
	return u, u.WithContext(defaultContext).Debug()
}

// CmsTenantDo returns the entity model and its Data Object query builder
func CmsTenantDo() (cmsTenant, *cmsTenantDo) {
	u := Use(db.CmsDatabase.DB).CmsTenant
	return u, u.WithContext(defaultContext).Debug()
}

// CmsThemeDo returns the entity model and its Data Object query builder
func CmsThemeDo() (cmsTheme, *cmsThemeDo) {
	u := Use(db.CmsDatabase.DB).CmsTheme
	return u, u.WithContext(defaultContext).Debug()
}

// CmsTopicDo returns the entity model and its Data Object query builder
func CmsTopicDo() (cmsTopic, *cmsTopicDo) {
	u := Use(db.CmsDatabase.DB).CmsTopic
	return u, u.WithContext(defaultContext).Debug()
}

// PlgOnlineRegisterDo returns the entity model and its Data Object query builder
func PlgOnlineRegisterDo() (plgOnlineRegister, *plgOnlineRegisterDo) {
	u := Use(db.CmsDatabase.DB).PlgOnlineRegister
	return u, u.WithContext(defaultContext).Debug()
}

// WeixinAccountDo returns the entity model and its Data Object query builder
func WeixinAccountDo() (weixinAccount, *weixinAccountDo) {
	u := Use(db.CmsDatabase.DB).WeixinAccount
	return u, u.WithContext(defaultContext).Debug()
}

// WeixinMenuDo returns the entity model and its Data Object query builder
func WeixinMenuDo() (weixinMenu, *weixinMenuDo) {
	u := Use(db.CmsDatabase.DB).WeixinMenu
	return u, u.WithContext(defaultContext).Debug()
}

// WeixinMpVerifyDo returns the entity model and its Data Object query builder
func WeixinMpVerifyDo() (weixinMpVerify, *weixinMpVerifyDo) {
	u := Use(db.CmsDatabase.DB).WeixinMpVerify
	return u, u.WithContext(defaultContext).Debug()
}

// WeixinRequestContentDo returns the entity model and its Data Object query builder
func WeixinRequestContentDo() (weixinRequestContent, *weixinRequestContentDo) {
	u := Use(db.CmsDatabase.DB).WeixinRequestContent
	return u, u.WithContext(defaultContext).Debug()
}

// WeixinRequestRuleDo returns the entity model and its Data Object query builder
func WeixinRequestRuleDo() (weixinRequestRule, *weixinRequestRuleDo) {
	u := Use(db.CmsDatabase.DB).WeixinRequestRule
	return u, u.WithContext(defaultContext).Debug()
}

// WeixinResponseContentDo returns the entity model and its Data Object query builder
func WeixinResponseContentDo() (weixinResponseContent, *weixinResponseContentDo) {
	u := Use(db.CmsDatabase.DB).WeixinResponseContent
	return u, u.WithContext(defaultContext).Debug()
}

