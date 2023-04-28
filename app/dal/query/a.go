package query

import (
	"context"

	"haedu.gov.cn/cms/app/dal"
)

var defaultContext = context.Background()

// AlbumDo
func CmsAlbumDo() (cmsAlbum, *cmsAlbumDo) {
	u := Use(dal.CmsDatabase.DB).CmsAlbum
	return u, u.WithContext(defaultContext).Debug()
}

// AttachDo
func CmsAttachDo() (cmsAttach, *cmsAttachDo) {
	u := Use(dal.CmsDatabase.DB).CmsAttach
	return u, u.WithContext(defaultContext).Debug()
}

// CmsAdminRoleDo
func CmsAdminRoleDo() (cmsAdminRole, *cmsAdminRoleDo) {
	u := Use(dal.CmsDatabase.DB).CmsAdminRole
	return u, u.WithContext(defaultContext).Debug()
}

// CmsAdminRoleValueDo
func CmsAdminRoleValueDo() (cmsAdminRoleValue, *cmsAdminRoleValueDo) {
	u := Use(dal.CmsDatabase.DB).CmsAdminRoleValue
	return u, u.WithContext(defaultContext).Debug()
}

// CmsAdminNavDo
func CmsAdminNavDo() (cmsAdminNav, *cmsAdminNavDo) {
	u := Use(dal.CmsDatabase.DB).CmsAdminNav
	return u, u.WithContext(defaultContext).Debug()
}

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

// CmsThemeDo
func CmsThemeDo() (cmsTheme, *cmsThemeDo) {
	u := Use(dal.CmsDatabase.DB).CmsTheme
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

// CmsArticleCommentDo
func CmsArticleCommentDo() (cmsArticleComment, *cmsArticleCommentDo) {
	u := Use(dal.CmsDatabase.DB).CmsArticleComment
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

// WeixinAccountDo
func WeixinAccountDo() (weixinAccount, *weixinAccountDo) {
	u := Use(dal.CmsDatabase.DB).WeixinAccount
	return u, u.WithContext(defaultContext).Debug()
}

// WeixinMenuDo
func WeixinMenuDo() (weixinMenu, *weixinMenuDo) {
	u := Use(dal.CmsDatabase.DB).WeixinMenu
	return u, u.WithContext(defaultContext).Debug()
}

// WeixinRequestRuleDo
func WeixinRequestRuleDo() (weixinRequestRule, *weixinRequestRuleDo) {
	u := Use(dal.CmsDatabase.DB).WeixinRequestRule
	return u, u.WithContext(defaultContext).Debug()
}

// WeixinRequestContentDo
func WeixinRequestContentDo() (weixinRequestContent, *weixinRequestContentDo) {
	u := Use(dal.CmsDatabase.DB).WeixinRequestContent
	return u, u.WithContext(defaultContext).Debug()
}

// WeixinResponseContentDo
func WeixinResponseContentDo() (weixinResponseContent, *weixinResponseContentDo) {
	u := Use(dal.CmsDatabase.DB).WeixinResponseContent
	return u, u.WithContext(defaultContext).Debug()
}

// PlgOnlineRegisterDo
func PlgOnlineRegisterDo() (plgOnlineRegister, *plgOnlineRegisterDo) {
	u := Use(dal.CmsDatabase.DB).PlgOnlineRegister
	return u, u.WithContext(defaultContext).Debug()
}

// CmsAdsCategoryDo
func CmsAdsCategoryDo() (cmsAdsCategory, *cmsAdsCategoryDo) {
	u := Use(dal.CmsDatabase.DB).CmsAdsCategory
	return u, u.WithContext(defaultContext).Debug()
}

// CmsAdsCategoryRelationDo
func CmsAdsCategoryRelationDo() (cmsAdsCategoryRelation, *cmsAdsCategoryRelationDo) {
	u := Use(dal.CmsDatabase.DB).CmsAdsCategoryRelation
	return u, u.WithContext(defaultContext).Debug()
}

// CmsAdsDo
func CmsAdsDo() (cmsAds, *cmsAdsDo) {
	u := Use(dal.CmsDatabase.DB).CmsAds
	return u, u.WithContext(defaultContext).Debug()
}

// CmsTagDo
func CmsTagDo() (cmsTag, *cmsTagDo) {
	u := Use(dal.CmsDatabase.DB).CmsTag
	return u, u.WithContext(defaultContext).Debug()
}

// CmsTopicDo
func CmsTopicDo() (cmsTopic, *cmsTopicDo) {
	u := Use(dal.CmsDatabase.DB).CmsTopic
	return u, u.WithContext(defaultContext).Debug()
}
