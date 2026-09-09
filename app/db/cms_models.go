package db

import (
	"reflect"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"

	"haedu.gov.cn/cms/app/cms/domain"
)

// runAutoMigrate 执行多数据库自动迁移
// 该函数会在应用启动时自动创建或更新所有数据表，
// 同时自动同步表注释、列注释、索引、类型修正等元数据。
//
// 模型分组说明：
//   - CMS 模型组(cmsModels): 管理员、广告、文章、站点、相册、链接、标签、主题、专题、插件 → CmsDatabase
//   - Auth 模型组(authModels): 微信公众号相关 → AuthDatabase(若不可用则回退至 CmsDatabase)
//
// 各模型对应表注释：
//   - cms_ad_category_relation: 广告分类关联表
//   - cms_admin_log: 管理员操作日志表
//   - cms_admin_nav: 管理员导航菜单表
//   - cms_admin_notice: 管理员公告表
//   - cms_admin_role: 管理员角色表
//   - cms_admin_role_site: 角色站点关联表
//   - cms_admin_role_value: 角色权限值表
//   - cms_admin: 管理员用户表
//   - cms_ads: 广告表
//   - cms_ads_category: 广告分类表
//   - cms_ads_category_relation: 广告与分类关联表
//   - cms_album: 相册表
//   - cms_article: 文章表
//   - cms_article_category: 文章分类表
//   - cms_article_category_relation: 文章与分类关联表
//   - cms_article_comment: 文章评论表
//   - cms_article_label: 文章标签表
//   - cms_article_label_relation: 文章与标签关联表
//   - cms_article_property: 文章属性表
//   - cms_attach: 附件表
//   - cms_link: 快捷链接表
//   - cms_link_category: 链接分类表
//   - cms_link_category_relation: 链接与分类关联表
//   - cms_site: 站点表
//   - cms_site_channel: 站点频道表
//   - cms_site_channel_album: 频道相册表
//   - cms_site_channel_field: 频道字段表
//   - cms_site_domain: 站点域名表
//   - cms_tag: 标签表
//   - cms_tenant: 租户表(多租户主表)
//   - cms_theme: 主题表
//   - cms_topic: 专题表
//   - plg_online_register: 在线报名表
//   - weixin_account: 微信公众号表
//   - weixin_menu: 微信菜单表
//   - weixin_mp_verify: 微信认证验证表
//   - weixin_request_content: 微信请求内容表
//   - weixin_request_rule: 微信请求规则表
//   - weixin_response_content: 微信响应内容表
func runAutoMigrate() {
	migrateEnabled, _ := web.AppConfig.Bool("cms.migrate")
	if !migrateEnabled {
		logs.Info("数据库自动迁移已禁用")
		return
	}

	logs.Info("开始执行多数据库自动迁移...")

	// 表注释映射 (表名 → 注释)
	cmsTableComments := buildCMSTableComments()
	authTableComments := buildAuthTableComments()

	// CMS 模型注册表: model → 期望的表注释
	cmsModels := buildCMSModels()
	authModels := buildAuthModels()

	// 将 CMS 模型迁移至 CmsDatabase
	migrateWithFullMeta(CmsDatabase, cmsDialect, cmsModels, cmsTableComments)

	// 将 Auth/微信模型迁移至 AuthDatabase(若可用)，否则回退至 CmsDatabase
	if AuthDatabase != nil {
		migrateWithFullMeta(AuthDatabase, authDialect, authModels, authTableComments)
	} else {
		logs.Warn("AuthDatabase 不可用，微信模型将回退迁移至 CmsDatabase")
		migrateWithFullMeta(CmsDatabase, cmsDialect, authModels, authTableComments)
	}

	totalCount := len(cmsModels) + len(authModels)
	logs.Info("多数据库自动迁移完成, 共迁移 %d 个表", totalCount)
}

// buildCMSTableComments 构造 CMS 表注释映射
func buildCMSTableComments() map[string]string {
	return map[string]string{
		"cms_ad_category_relation":      "广告分类关联表",
		"cms_admin_log":                 "管理员操作日志表",
		"cms_admin_nav":                 "管理员导航菜单表",
		"cms_admin_notice":              "管理员公告表",
		"cms_admin_role":                "管理员角色表",
		"cms_admin_role_site":           "角色站点关联表",
		"cms_admin_role_value":          "角色权限值表",
		"cms_admin":                     "管理员用户表",
		"cms_ads":                       "广告表",
		"cms_ads_category":              "广告分类表",
		"cms_ads_category_relation":     "广告与分类关联表",
		"cms_album":                     "相册表",
		"cms_article":                   "文章表",
		"cms_article_category":          "文章分类表",
		"cms_article_category_relation": "文章与分类关联表",
		"cms_article_comment":           "文章评论表",
		"cms_article_label":             "文章标签表",
		"cms_article_label_relation":    "文章与标签关联表",
		"cms_article_property":          "文章属性表",
		"cms_attach":                    "附件表",
		"cms_link":                      "快捷链接表",
		"cms_link_category":             "链接分类表",
		"cms_link_category_relation":    "链接与分类关联表",
		"cms_site":                      "站点表",
		"cms_site_channel":              "站点频道表",
		"cms_site_channel_album":        "频道相册表",
		"cms_site_channel_field":        "频道字段表",
		"cms_site_domain":               "站点域名表",
		"cms_tag":                       "标签表",
		"cms_tenant":                    "租户表",
		"cms_theme":                     "主题表",
		"cms_topic":                     "专题表",
		"plg_online_register":           "在线报名表",
	}
}

// buildAuthTableComments 构造 Auth 表注释映射
func buildAuthTableComments() map[string]string {
	return map[string]string{
		"weixin_account":          "微信公众号表",
		"weixin_menu":             "微信菜单表",
		"weixin_mp_verify":        "微信认证验证表",
		"weixin_request_content":  "微信请求内容表",
		"weixin_request_rule":     "微信请求规则表",
		"weixin_response_content": "微信响应内容表",
	}
}

// buildCMSModels 构造 CMS 模型列表
func buildCMSModels() []any {
	return []any{
		&domain.CmsAdCategoryRelation{},
		&domain.CmsAdminLog{},
		&domain.CmsAdminNav{},
		&domain.CmsAdminNotice{},
		&domain.CmsAdminRoleSite{},
		&domain.CmsAdminRoleValue{},
		&domain.CmsAdminRole{},
		&domain.CmsAdmin{},
		&domain.CmsAdsCategoryRelation{},
		&domain.CmsAdsCategory{},
		&domain.CmsAds{},
		&domain.CmsAlbum{},
		&domain.CmsArticleCategoryRelation{},
		&domain.CmsArticleCategory{},
		&domain.CmsArticleComment{},
		&domain.CmsArticleLabelRelation{},
		&domain.CmsArticleLabel{},
		&domain.CmsArticleProperty{},
		&domain.CmsArticle{},
		&domain.CmsAttach{},
		&domain.CmsLinkCategoryRelation{},
		&domain.CmsLinkCategory{},
		&domain.CmsLink{},
		&domain.CmsSiteChannelAlbum{},
		&domain.CmsSiteChannelField{},
		&domain.CmsSiteChannel{},
		&domain.CmsSiteDomain{},
		&domain.CmsSite{},
		&domain.CmsTag{},
		&domain.CmsTenant{},
		&domain.CmsTheme{},
		&domain.CmsTopic{},
		&domain.PlgOnlineRegister{},
	}
}

// buildAuthModels 构造 Auth 模型列表
func buildAuthModels() []any {
	return []any{
		&domain.WeixinAccount{},
		&domain.WeixinMenu{},
		&domain.WeixinMpVerify{},
		&domain.WeixinRequestContent{},
		&domain.WeixinRequestRule{},
		&domain.WeixinResponseContent{},
	}
}

// migrateWithFullMeta 对一组模型执行完整的元数据同步：
//   1. GORM AutoMigrate (CREATE TABLE IF NOT EXISTS, ADD COLUMN IF NOT EXISTS)
//   2. 列注释 (从 struct gorm tag comment: 提取)
//   3. 索引创建 (从 struct gorm tag index 提取)
//   4. 表注释 (从 tableComments 映射)
//   5. 类型修正 (bit(1) → tinyint(1) 等跨方言兼容)
func migrateWithFullMeta(db *DBExtension, dialect string, models []any, tableComments map[string]string) {
	if db == nil {
		logs.Warn("跳过迁移: 数据库实例为空")
		return
	}

	for _, model := range models {
		meta, err := extractTableMeta(model)
		if err != nil {
			logs.Warn("提取元数据失败[%T]: %v", model, err)
			continue
		}

		tableName := meta.TableName
		comment, hasComment := tableComments[tableName]

		// 1) GORM AutoMigrate (CREATE TABLE IF NOT EXISTS / ADD COLUMN)
		if err := db.AutoMigrate(model); err != nil {
			logs.Warn("GORM AutoMigrate 失败[%s]: %v", tableName, err)
			// 不 return,继续补充注释/索引
		} else {
			logs.Debug("GORM 迁移成功: %s", tableName)
		}

		// 2) 列注释 - 从 struct gorm tag 提取
		for _, col := range meta.AllFieldRefs {
			if !col.HasComment || col.DBType == "" {
				continue
			}
			if err := SetColumnComment(db, dialect, tableName, col.ColumnName, col.DBType, col.Comment); err != nil {
				logs.Warn("设置列注释失败 %s.%s: %v", tableName, col.ColumnName, err)
			}
		}

		// 3) 索引创建 - 从 struct gorm tag 提取
		for _, col := range meta.AllFieldRefs {
			if col.PrimaryKey || !col.IsIndex {
				continue
			}
			idxName := "idx_" + tableName + "_" + col.ColumnName
			if _, err := EnsureIndex(db, dialect, tableName, idxName, []string{col.ColumnName}, col.UniqueIndex); err != nil {
				logs.Warn("创建索引失败 %s: %v", idxName, err)
			}
		}

		// 4) 表注释
		if hasComment && comment != "" {
			if err := SetTableComment(db, dialect, tableName, comment); err != nil {
				logs.Warn("设置表注释失败 %s: %v", tableName, err)
			}
		}

		// 5) 类型修正 - bit(1) → tinyint(1) (MySQL 兼容性修复)
		if dialect == DialectMySQL {
			if err := FixBit1ToTinyInt1(db, dialect, tableName); err != nil {
				logs.Warn("类型修正失败 %s: %v", tableName, err)
			}
		}
	}
}

// ensureReflectImported 保留 reflect 引用以备扩展(结构体元数据扫描)
var _ = reflect.TypeOf
