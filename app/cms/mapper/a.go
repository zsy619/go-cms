package mapper

import (
	"context"

	"haedu.gov.cn/cms/app/db"
)

var defaultContext = context.Background()

// AlbumDo
func CmsAlbumDo() (cmsAlbum, *cmsAlbumDo) {
	u := Use(db.CmsDatabase.DB).CmsAlbum
	return u, u.WithContext(defaultContext).Debug()
}

// AttachDo
func CmsAttachDo() (cmsAttach, *cmsAttachDo) {
	u := Use(db.CmsDatabase.DB).CmsAttach
	return u, u.WithContext(defaultContext).Debug()
}

// CmsAdminRoleDo
func CmsAdminRoleDo() (cmsAdminRole, *cmsAdminRoleDo) {
	u := Use(db.CmsDatabase.DB).CmsAdminRole
	return u, u.WithContext(defaultContext).Debug()
}

// CmsAdminRoleValueDo
func CmsAdminRoleValueDo() (cmsAdminRoleValue, *cmsAdminRoleValueDo) {
	u := Use(db.CmsDatabase.DB).CmsAdminRoleValue
	return u, u.WithContext(defaultContext).Debug()
}

// CmsAdminRoleSiteDo
func CmsAdminRoleSiteDo() (cmsAdminRoleSite, *cmsAdminRoleSiteDo) {
	u := Use(db.CmsDatabase.DB).CmsAdminRoleSite
	return u, u.WithContext(defaultContext).Debug()
}

// CmsAdminNavDo
func CmsAdminNavDo() (cmsAdminNav, *cmsAdminNavDo) {
	u := Use(db.CmsDatabase.DB).CmsAdminNav
	return u, u.WithContext(defaultContext).Debug()
}

// CmsAdminDo
func CmsAdminDo() (cmsAdmin, *cmsAdminDo) {
	u := Use(db.CmsDatabase.DB).CmsAdmin
	return u, u.WithContext(defaultContext).Debug()
}

// CmsAdminLogDo
func CmsAdminLogDo() (cmsAdminLog, *cmsAdminLogDo) {
	u := Use(db.CmsDatabase.DB).CmsAdminLog
	return u, u.WithContext(defaultContext).Debug()
}

// CmsThemeDo
func CmsThemeDo() (cmsTheme, *cmsThemeDo) {
	u := Use(db.CmsDatabase.DB).CmsTheme
	return u, u.WithContext(defaultContext).Debug()
}

// CmsSiteDo
func CmsSiteDo() (cmsSite, *cmsSiteDo) {
	u := Use(db.CmsDatabase.DB).CmsSite
	return u, u.WithContext(defaultContext).Debug()
}

// CmsSiteDomainDo
func CmsSiteDomainDo() (cmsSiteDomain, *cmsSiteDomainDo) {
	u := Use(db.CmsDatabase.DB).CmsSiteDomain
	return u, u.WithContext(defaultContext).Debug()
}

// CmsSiteChannelDo
func CmsSiteChannelDo() (cmsSiteChannel, *cmsSiteChannelDo) {
	u := Use(db.CmsDatabase.DB).CmsSiteChannel
	return u, u.WithContext(defaultContext).Debug()
}

// CmsArticleCategoryDo
func CmsArticleCategoryDo() (cmsArticleCategory, *cmsArticleCategoryDo) {
	u := Use(db.CmsDatabase.DB).CmsArticleCategory
	return u, u.WithContext(defaultContext).Debug()
}

// CmsArticleCategoryRelationDo
func CmsArticleCategoryRelationDo() (cmsArticleCategoryRelation, *cmsArticleCategoryRelationDo) {
	u := Use(db.CmsDatabase.DB).CmsArticleCategoryRelation
	return u, u.WithContext(defaultContext).Debug()
}

// CmsArticleDo
func CmsArticleDo() (cmsArticle, *cmsArticleDo) {
	u := Use(db.CmsDatabase.DB).CmsArticle
	return u, u.WithContext(defaultContext).Debug()
}

// CmsArticleLabelDo
func CmsArticleLabelDo() (cmsArticleLabel, *cmsArticleLabelDo) {
	u := Use(db.CmsDatabase.DB).CmsArticleLabel
	return u, u.WithContext(defaultContext).Debug()
}

// CmsArticleLabelRelationDo
func CmsArticleLabelRelationDo() (cmsArticleLabelRelation, *cmsArticleLabelRelationDo) {
	u := Use(db.CmsDatabase.DB).CmsArticleLabelRelation
	return u, u.WithContext(defaultContext).Debug()
}

// CmsArticleCommentDo
func CmsArticleCommentDo() (cmsArticleComment, *cmsArticleCommentDo) {
	u := Use(db.CmsDatabase.DB).CmsArticleComment
	return u, u.WithContext(defaultContext).Debug()
}

// CmsArticlePropertyDo
func CmsArticlePropertyDo() (cmsArticleProperty, *cmsArticlePropertyDo) {
	u := Use(db.CmsDatabase.DB).CmsArticleProperty
	return u, u.WithContext(defaultContext).Debug()
}

// CmsLinkCategoryDo
func CmsLinkCategoryDo() (cmsLinkCategory, *cmsLinkCategoryDo) {
	u := Use(db.CmsDatabase.DB).CmsLinkCategory
	return u, u.WithContext(defaultContext).Debug()
}

// CmsLinkCategoryRelationDo
func CmsLinkCategoryRelationDo() (cmsLinkCategoryRelation, *cmsLinkCategoryRelationDo) {
	u := Use(db.CmsDatabase.DB).CmsLinkCategoryRelation
	return u, u.WithContext(defaultContext).Debug()
}

// CmsLinkDo
func CmsLinkDo() (cmsLink, *cmsLinkDo) {
	u := Use(db.CmsDatabase.DB).CmsLink
	return u, u.WithContext(defaultContext).Debug()
}

// WeixinAccountDo
func WeixinAccountDo() (weixinAccount, *weixinAccountDo) {
	u := Use(db.CmsDatabase.DB).WeixinAccount
	return u, u.WithContext(defaultContext).Debug()
}

// WeixinMenuDo
func WeixinMenuDo() (weixinMenu, *weixinMenuDo) {
	u := Use(db.CmsDatabase.DB).WeixinMenu
	return u, u.WithContext(defaultContext).Debug()
}

// WeixinMpVerifyDo
func WeixinMpVerifyDo() (weixinMpVerify, *weixinMpVerifyDo) {
	u := Use(db.CmsDatabase.DB).WeixinMpVerify
	return u, u.WithContext(defaultContext).Debug()
}

// WeixinRequestRuleDo
func WeixinRequestRuleDo() (weixinRequestRule, *weixinRequestRuleDo) {
	u := Use(db.CmsDatabase.DB).WeixinRequestRule
	return u, u.WithContext(defaultContext).Debug()
}

// WeixinRequestContentDo
func WeixinRequestContentDo() (weixinRequestContent, *weixinRequestContentDo) {
	u := Use(db.CmsDatabase.DB).WeixinRequestContent
	return u, u.WithContext(defaultContext).Debug()
}

// WeixinResponseContentDo
func WeixinResponseContentDo() (weixinResponseContent, *weixinResponseContentDo) {
	u := Use(db.CmsDatabase.DB).WeixinResponseContent
	return u, u.WithContext(defaultContext).Debug()
}

// PlgOnlineRegisterDo
func PlgOnlineRegisterDo() (plgOnlineRegister, *plgOnlineRegisterDo) {
	u := Use(db.CmsDatabase.DB).PlgOnlineRegister
	return u, u.WithContext(defaultContext).Debug()
}

// CmsAdsCategoryDo
func CmsAdsCategoryDo() (cmsAdsCategory, *cmsAdsCategoryDo) {
	u := Use(db.CmsDatabase.DB).CmsAdsCategory
	return u, u.WithContext(defaultContext).Debug()
}

// CmsAdsCategoryRelationDo
func CmsAdsCategoryRelationDo() (cmsAdsCategoryRelation, *cmsAdsCategoryRelationDo) {
	u := Use(db.CmsDatabase.DB).CmsAdsCategoryRelation
	return u, u.WithContext(defaultContext).Debug()
}

// CmsAdsDo
func CmsAdsDo() (cmsAds, *cmsAdsDo) {
	u := Use(db.CmsDatabase.DB).CmsAds
	return u, u.WithContext(defaultContext).Debug()
}

// CmsTagDo
func CmsTagDo() (cmsTag, *cmsTagDo) {
	u := Use(db.CmsDatabase.DB).CmsTag
	return u, u.WithContext(defaultContext).Debug()
}

// CmsTopicDo
func CmsTopicDo() (cmsTopic, *cmsTopicDo) {
	u := Use(db.CmsDatabase.DB).CmsTopic
	return u, u.WithContext(defaultContext).Debug()
}

// CmsAdminNoticeDo
func CmsAdminNoticeDo() (cmsAdminNotice, *cmsAdminNoticeDo) {
	u := Use(db.CmsDatabase.DB).CmsAdminNotice
	return u, u.WithContext(defaultContext).Debug()
}
