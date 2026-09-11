// Package db 提供多数据库自动迁移能力,涵盖 GORM 模型注册、表/列注释同步、索引创建、
// 跨方言类型修正等"启动期元数据对齐"工作。
//
// 关键约定:
//   - CMS 模型 → CmsDatabase(主业务库)
//   - Auth/微信模型 → AuthDatabase(若不可用则回退至 CmsDatabase)
//   - 表注释中文文案统一由 buildCMSTableComments / buildAuthTableComments 提供
//   - 模型注册列表分别由 buildCMSModels / buildAuthModels 提供,与 domain 包严格对应
//   - 任何一步失败仅记告警日志,不会中断整体迁移流程,保证"局部错误不影响全局"
//
// 本文件不持有连接,仅描述"应当迁移哪些表/模型",实际执行由 migrateWithFullMeta 完成。
package db

import (
	"reflect"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"

	"haedu.gov.cn/cms/app/cms/domain"
)

// runAutoMigrate 执行多数据库自动迁移
//
// 流程概述:
//  1. 读取配置文件 cms.migrate 开关,关闭则直接返回(常用于生产环境禁用)
//  2. 分别构造 CMS / Auth 两组模型的"表注释"映射与"模型注册表"
//  3. 调用 migrateWithFullMeta 将 CMS 模型迁移到 CmsDatabase
//  4. 根据 AuthDatabase 是否可用,决定 Auth 模型的迁移目标
//  5. 输出迁移完成的统计日志
//
// 各模型对应表注释:
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
//   - cms_language: 语言表
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
	// 读取迁移开关;生产环境通常设为 false,避免 AutoMigrate 在高峰期执行
	migrateEnabled, _ := web.AppConfig.Bool("cms.migrate")
	if !migrateEnabled {
		logs.Info("数据库自动迁移已禁用")
		return
	}
	logs.Info("开始执行多数据库自动迁移...")
	// 表注释映射 (表名 → 注释),在 AutoMigrate 之后用于补齐 COMMENT
	cmsTableComments := buildCMSTableComments()
	authTableComments := buildAuthTableComments()
	// CMS 模型注册表: model 列表,AutoMigrate 与 extractTableMeta 都依赖此列表
	cmsModels := buildCMSModels()
	authModels := buildAuthModels()
	// 将 CMS 模型迁移至 CmsDatabase(主业务库,涵盖站点/文章/广告/标签/链接/语言/插件 等)
	migrateWithFullMeta(CmsDatabase, cmsDialect, cmsModels, cmsTableComments)
	// 将 Auth/微信模型迁移至 AuthDatabase(若可用),否则回退至 CmsDatabase
	if AuthDatabase != nil {
		migrateWithFullMeta(AuthDatabase, authDialect, authModels, authTableComments)
	} else {
		logs.Warn("AuthDatabase 不可用，微信模型将回退迁移至 CmsDatabase")
		migrateWithFullMeta(CmsDatabase, cmsDialect, authModels, authTableComments)
	}
	// 输出迁移汇总,便于运维快速确认迁移规模
	totalCount := len(cmsModels) + len(authModels)
	logs.Info("多数据库自动迁移完成, 共迁移 %d 个表", totalCount)
}

// buildCMSTableComments 构造 CMS 表注释映射
//
// 返回值: 表名 → 人类可读的中文注释
//
// 该映射用于在 AutoMigrate 之后,对已存在的表执行 ALTER TABLE ... COMMENT = '...',
// 让 DBA 在 Navicat/数据字典等工具中能直观看到表的中文业务含义。
//
// 注意: 一旦新增 CMS 表,必须同时在 domain 包定义结构体 + 在本 map 注册表注释 +
// 在 buildCMSModels 中加入模型指针,三者缺一不可。
func buildCMSTableComments() map[string]string {
	return map[string]string{
		"cms_ad_category_relation":      "广告分类关联表",  // 多对多: cms_ads ↔ cms_ads_category
		"cms_admin_log":                 "管理员操作日志表", // 记录登录/增删改 等审计事件
		"cms_admin_nav":                 "管理员导航菜单表", // 后台左侧菜单的层级结构
		"cms_admin_notice":              "管理员公告表",   // 系统级公告/通知
		"cms_admin_role":                "管理员角色表",   // RBAC 中的角色定义
		"cms_admin_role_site":           "角色站点关联表",  // 角色与可管理站点的多对多映射
		"cms_admin_role_value":          "角色权限值表",   // 角色对应的细粒度权限点
		"cms_admin":                     "管理员用户表",   // 后台账号
		"cms_ads":                       "广告表",      // 广告位内容
		"cms_ads_category":              "广告分类表",    // 广告位分类
		"cms_ads_category_relation":     "广告与分类关联表", // 广告 ↔ 分类 多对多
		"cms_album":                     "相册表",      // 图片集合,关联到频道
		"cms_article":                   "文章表",      // 内容主体,关联 channel/category/标签
		"cms_article_category":          "文章分类表",    // 文章所属分类
		"cms_article_category_relation": "文章与分类关联表", // 文章 ↔ 分类 多对多
		"cms_article_comment":           "文章评论表",    // 用户评论,支持回复(自引用)
		"cms_article_label":             "文章标签表",    // 文章打标
		"cms_article_label_relation":    "文章与标签关联表", // 文章 ↔ 标签 多对多
		"cms_article_property":          "文章属性表",    // 文章扩展属性(推荐/置顶 等)
		"cms_attach":                    "附件表",      // 通用文件上传,记录 OSS/本地路径
		"cms_language":                  "语言表",      // 多语言站点配置(代码/图标/默认)
		"cms_link":                      "快捷链接表",    // 友情链接/底部导航
		"cms_link_category":             "链接分类表",    // 链接分组
		"cms_link_category_relation":    "链接与分类关联表", // 链接 ↔ 分类 多对多
		"cms_site":                      "站点表",      // 多站点主表,按 tenant_id 隔离
		"cms_site_channel":              "站点频道表",    // 站点下的栏目/频道
		"cms_site_channel_album":        "频道相册表",    // 频道 ↔ 相册 的多对多绑定
		"cms_site_channel_field":        "频道字段表",    // 频道自定义字段(扩展模型)
		"cms_site_domain":               "站点域名表",    // 多域名绑定到同一站点
		"cms_tag":                       "标签表",      // 全局标签,可在文章/产品 等场景复用
		"cms_tenant":                    "租户表",      // 多租户隔离主表
		"cms_theme":                     "主题表",      // 站点主题/皮肤
		"cms_topic":                     "专题表",      // 内容聚合专题
		"plg_online_register":           "在线报名表",    // 表单插件,记录用户提交
	}
}

// buildAuthTableComments 构造 Auth 表注释映射
//
// 返回值: 表名 → 人类可读的中文注释
//
// 与 buildCMSTableComments 同源设计,但只服务于 Auth 库(微信公众号相关业务)。
// 注意: AuthDatabase 在初始化失败时,本组表会回退到 CmsDatabase,因此注释仍由本函数统一定义,
// 避免注释分裂。
func buildAuthTableComments() map[string]string {
	return map[string]string{
		"weixin_account":          "微信公众号表",  // 公众号接入凭证(app_id/secret/aes_key 等)
		"weixin_menu":             "微信菜单表",   // 自定义菜单的 JSON 配置
		"weixin_mp_verify":        "微信认证验证表", // 服务器 URL 校验 token
		"weixin_request_content":  "微信请求内容表", // 用户消息/事件 关键字回复内容
		"weixin_request_rule":     "微信请求规则表", // 关键字匹配规则
		"weixin_response_content": "微信响应内容表", // 关键字命中后返回的回复(文本/图文 等)
	}
}

// buildCMSModels 构造 CMS 模型列表
//
// 返回值: GORM 模型指针切片(传给 AutoMigrate 时使用指针,触发 GORM 反射元数据)
//
// 该列表会被传入 migrateWithFullMeta,顺序执行以下工作:
//  1. AutoMigrate 自动建表/补列
//  2. 同步列注释、索引
//  3. 同步表注释
//  4. 跨方言类型修正
//
// 列表顺序遵循"先依赖后被依赖"原则,例如:
//   - *_relation 表放在对应主表之后,确保先有主表再建关联
//   - cms_article 放在 cms_article_category 之后,因为文章会引用分类 id
func buildCMSModels() []any {
	return []any{
		&domain.CmsAdCategoryRelation{},      // 广告 ↔ 广告分类 关联
		&domain.CmsAdminLog{},                // 管理员操作日志
		&domain.CmsAdminNav{},                // 后台导航菜单
		&domain.CmsAdminNotice{},             // 系统公告
		&domain.CmsAdminRoleSite{},           // 角色 ↔ 站点 关联
		&domain.CmsAdminRoleValue{},          // 角色权限点
		&domain.CmsAdminRole{},               // 管理员角色
		&domain.CmsAdmin{},                   // 管理员用户
		&domain.CmsAdsCategoryRelation{},     // 广告 ↔ 广告分类 关联(独立命名空间,与上面 ad_category_relation 区分)
		&domain.CmsAdsCategory{},             // 广告分类
		&domain.CmsAds{},                     // 广告位内容
		&domain.CmsAlbum{},                   // 相册
		&domain.CmsArticleCategoryRelation{}, // 文章 ↔ 文章分类 关联
		&domain.CmsArticleCategory{},         // 文章分类
		&domain.CmsArticleComment{},          // 文章评论(自引用支持回复)
		&domain.CmsArticleLabelRelation{},    // 文章 ↔ 标签 关联
		&domain.CmsArticleLabel{},            // 文章标签
		&domain.CmsArticleProperty{},         // 文章扩展属性
		&domain.CmsArticle{},                 // 文章主表
		&domain.CmsAttach{},                  // 通用附件
		&domain.CmsLanguage{},                // 多语言配置(代码/图标/默认)
		&domain.CmsLinkCategoryRelation{},    // 链接 ↔ 链接分类 关联
		&domain.CmsLinkCategory{},            // 链接分类
		&domain.CmsLink{},                    // 友情链接/底部导航
		&domain.CmsSiteChannelAlbum{},        // 频道 ↔ 相册 关联
		&domain.CmsSiteChannelField{},        // 频道自定义字段
		&domain.CmsSiteChannel{},             // 站点频道
		&domain.CmsSiteDomain{},              // 站点多域名绑定
		&domain.CmsSite{},                    // 站点主表
		&domain.CmsTag{},                     // 全局标签
		&domain.CmsTenant{},                  // 多租户主表
		&domain.CmsTheme{},                   // 站点主题
		&domain.CmsTopic{},                   // 内容专题
		&domain.PlgOnlineRegister{},          // 在线报名表单
	}
}

// buildAuthModels 构造 Auth 模型列表
//
// 返回值: GORM 模型指针切片
//
// 与 buildCMSModels 配对使用;仅包含微信生态相关的 6 张表,均在 AuthDatabase 中。
// 若 AuthDatabase 不可用,这些模型会被回退迁移至 CmsDatabase,见 runAutoMigrate。
func buildAuthModels() []any {
	return []any{
		&domain.WeixinAccount{},         // 微信公众号接入凭证
		&domain.WeixinMenu{},            // 自定义菜单 JSON
		&domain.WeixinMpVerify{},        // 公众号服务器 URL 校验 token
		&domain.WeixinRequestContent{},  // 关键字回复内容
		&domain.WeixinRequestRule{},     // 关键字匹配规则
		&domain.WeixinResponseContent{}, // 命中后的回复内容(文本/图文)
	}
}

// migrateWithFullMeta 对一组模型执行完整的元数据同步
//
// 工作步骤:
//  1. GORM AutoMigrate  → CREATE TABLE IF NOT EXISTS / ADD COLUMN IF NOT EXISTS
//  2. 列注释             → 从 struct gorm tag comment: 提取,通过 SetColumnComment 写入
//  3. 索引创建           → 从 struct gorm tag index 提取,通过 EnsureIndex 创建
//  4. 表注释             → 从 tableComments 映射,通过 SetTableComment 写入
//  5. 类型修正           → bit(1) → tinyint(1) 等跨方言兼容,仅 MySQL 执行
//
// 设计要点:
//   - 任一步骤失败不会中断整体流程,只会日志告警并继续处理下一张表
//   - 这样保证即使个别表的迁移报错,其他表仍能完成同步,避免"一张失败全表回滚"
//   - db 为 nil 时直接返回(例如 AuthDatabase 未配置的场景)
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
		// 1) GORM AutoMigrate: CREATE TABLE IF NOT EXISTS / ADD COLUMN
		if err := db.AutoMigrate(model); err != nil {
			logs.Warn("GORM AutoMigrate 失败[%s]: %v", tableName, err)
			// 不 return,继续补充注释/索引,即便建表失败也尽量把已存在表的元数据补齐
		} else {
			logs.Debug("GORM 迁移成功: %s", tableName)
		}
		// 2) 列注释 - 从 struct gorm tag comment 提取
		for _, col := range meta.AllFieldRefs {
			if !col.HasComment || col.DBType == "" {
				continue
			}
			if err := SetColumnComment(db, dialect, tableName, col.ColumnName, col.DBType, col.Comment); err != nil {
				logs.Warn("设置列注释失败 %s.%s: %v", tableName, col.ColumnName, err)
			}
			// 3) 索引创建 - 从 struct gorm tag index 提取
			// 排除主键(主键索引由数据库自动创建) 与 未声明索引需求的列
			if col.PrimaryKey || !col.IsIndex {
				continue
			}
			idxName := "idx_" + tableName + "_" + col.ColumnName
			if _, err := EnsureIndex(db, dialect, tableName, idxName, []string{col.ColumnName}, col.UniqueIndex); err != nil {
				logs.Warn("创建索引失败 %s: %v", idxName, err)
			}
		}
		// 4) 表注释 - 来自 buildCMSTableComments / buildAuthTableComments
		if hasComment && comment != "" {
			if err := SetTableComment(db, dialect, tableName, comment); err != nil {
				logs.Warn("设置表注释失败 %s: %v", tableName, err)
			}
		}
		// 5) 类型修正 - bit(1) → tinyint(1) (MySQL 兼容性修复)
		// 历史原因: 早期 MySQL 默认 bool 映射为 bit(1),部分 ORM 客户端读取不友好,
		// 这里统一改成 tinyint(1),跨方言语义更清晰
		if dialect == DialectMySQL {
			if err := FixBit1ToTinyInt1(db, dialect, tableName); err != nil {
				logs.Warn("类型修正失败 %s: %v", tableName, err)
			}
		}
	}
}

// ensureReflectImported 保留 reflect 引用以备扩展(结构体元数据扫描)
//
// 说明: extractTableMeta 内部依赖 reflect 反射读取 struct tag,本变量作为占位,
// 保证 reflect 包不被 goimports 误删,也方便后续新增反射工具时复用 import。
var _ = reflect.TypeOf
