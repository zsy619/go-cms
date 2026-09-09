-- ============================================================
-- 多租户 + 软删除 字段迁移脚本
-- 数据库: Microsoft SQL Server
-- 设计参考: 芋道 yudao BaseDO
-- 生成时间: 2026-09-09 10:28:57
-- ============================================================
USE [cms];
GO

-- ========== 1. 租户主表 cms_tenant ==========
IF NOT EXISTS (SELECT * FROM sys.objects WHERE object_id = OBJECT_ID(N'[dbo].[cms_tenant]') AND type in (N'U'))
BEGIN
CREATE TABLE [dbo].[cms_tenant] (
    [tenant_id]    BIGINT IDENTITY(1,1) NOT NULL,
    [code]         NVARCHAR(64) NOT NULL,
    [name]         NVARCHAR(128) NOT NULL,
    [short_name]   NVARCHAR(64) NULL,
    [contact_name] NVARCHAR(64) NULL,
    [contact_tel]  NVARCHAR(32) NULL,
    [contact_mail] NVARCHAR(128) NULL,
    [description]  NVARCHAR(512) NULL,
    [logo]         NVARCHAR(512) NULL,
    [domain]       NVARCHAR(256) NULL,
    [expire_time]  DATETIME NULL,
    [max_users]    INT NULL DEFAULT 0,
    [status]       TINYINT NULL DEFAULT 1,
    [create_id]    INT NULL,
    [create_name]  NVARCHAR(64) NULL,
    [create_time]  DATETIME NULL DEFAULT GETDATE(),
    [update_id]    INT NULL,
    [update_name]  NVARCHAR(64) NULL,
    [update_time]  DATETIME NULL,
    [deleted]      BIT NOT NULL DEFAULT 0,
    CONSTRAINT [PK_cms_tenant] PRIMARY KEY CLUSTERED ([tenant_id]),
    CONSTRAINT [UK_cms_tenant_code] UNIQUE ([code])
);
EXEC sp_addextendedproperty 'MS_Description', N'租户表(多租户主表)', 'SCHEMA', N'dbo', 'TABLE', N'cms_tenant';
CREATE INDEX [idx_cms_tenant_status] ON [dbo].[cms_tenant]([status]);
CREATE INDEX [idx_cms_tenant_deleted] ON [dbo].[cms_tenant]([deleted]);
END
GO

-- ========== 2. 给 cms_admin 添加 tenant_id 和 deleted ==========
IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID(N'[dbo].[cms_admin]') AND name = 'tenant_id')
BEGIN
    ALTER TABLE [dbo].[cms_admin] ADD [tenant_id] BIGINT NOT NULL DEFAULT 0;
    CREATE INDEX [idx_cms_admin_tenant] ON [dbo].[cms_admin]([tenant_id]);
END
GO

IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID(N'[dbo].[cms_admin]') AND name = 'deleted')
BEGIN
    ALTER TABLE [dbo].[cms_admin] ADD [deleted] BIT NOT NULL DEFAULT 0;
    CREATE INDEX [idx_cms_admin_deleted] ON [dbo].[cms_admin]([deleted]);
END
GO

-- ========== 3. 给所有业务表添加列 ==========
IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID(N'[dbo].[cms_ad_category_relation]') AND name = 'tenant_id')
BEGIN
    ALTER TABLE [dbo].[cms_ad_category_relation] ADD [tenant_id] BIGINT NOT NULL DEFAULT 0;
    ALTER TABLE [dbo].[cms_ad_category_relation] ADD [deleted] BIT NOT NULL DEFAULT 0;
    CREATE INDEX [idx_cms_ad_category_relation_tenant] ON [dbo].[cms_ad_category_relation]([tenant_id]);
    CREATE INDEX [idx_cms_ad_category_relation_deleted] ON [dbo].[cms_ad_category_relation]([deleted]);
END
GO

IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID(N'[dbo].[cms_admin_log]') AND name = 'tenant_id')
BEGIN
    ALTER TABLE [dbo].[cms_admin_log] ADD [tenant_id] BIGINT NOT NULL DEFAULT 0;
    ALTER TABLE [dbo].[cms_admin_log] ADD [deleted] BIT NOT NULL DEFAULT 0;
    CREATE INDEX [idx_cms_admin_log_tenant] ON [dbo].[cms_admin_log]([tenant_id]);
    CREATE INDEX [idx_cms_admin_log_deleted] ON [dbo].[cms_admin_log]([deleted]);
END
GO

IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID(N'[dbo].[cms_admin_nav]') AND name = 'tenant_id')
BEGIN
    ALTER TABLE [dbo].[cms_admin_nav] ADD [tenant_id] BIGINT NOT NULL DEFAULT 0;
    ALTER TABLE [dbo].[cms_admin_nav] ADD [deleted] BIT NOT NULL DEFAULT 0;
    CREATE INDEX [idx_cms_admin_nav_tenant] ON [dbo].[cms_admin_nav]([tenant_id]);
    CREATE INDEX [idx_cms_admin_nav_deleted] ON [dbo].[cms_admin_nav]([deleted]);
END
GO

IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID(N'[dbo].[cms_admin_notice]') AND name = 'tenant_id')
BEGIN
    ALTER TABLE [dbo].[cms_admin_notice] ADD [tenant_id] BIGINT NOT NULL DEFAULT 0;
    ALTER TABLE [dbo].[cms_admin_notice] ADD [deleted] BIT NOT NULL DEFAULT 0;
    CREATE INDEX [idx_cms_admin_notice_tenant] ON [dbo].[cms_admin_notice]([tenant_id]);
    CREATE INDEX [idx_cms_admin_notice_deleted] ON [dbo].[cms_admin_notice]([deleted]);
END
GO

IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID(N'[dbo].[cms_admin_role]') AND name = 'tenant_id')
BEGIN
    ALTER TABLE [dbo].[cms_admin_role] ADD [tenant_id] BIGINT NOT NULL DEFAULT 0;
    ALTER TABLE [dbo].[cms_admin_role] ADD [deleted] BIT NOT NULL DEFAULT 0;
    CREATE INDEX [idx_cms_admin_role_tenant] ON [dbo].[cms_admin_role]([tenant_id]);
    CREATE INDEX [idx_cms_admin_role_deleted] ON [dbo].[cms_admin_role]([deleted]);
END
GO

IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID(N'[dbo].[cms_admin_role_site]') AND name = 'tenant_id')
BEGIN
    ALTER TABLE [dbo].[cms_admin_role_site] ADD [tenant_id] BIGINT NOT NULL DEFAULT 0;
    ALTER TABLE [dbo].[cms_admin_role_site] ADD [deleted] BIT NOT NULL DEFAULT 0;
    CREATE INDEX [idx_cms_admin_role_site_tenant] ON [dbo].[cms_admin_role_site]([tenant_id]);
    CREATE INDEX [idx_cms_admin_role_site_deleted] ON [dbo].[cms_admin_role_site]([deleted]);
END
GO

IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID(N'[dbo].[cms_admin_role_value]') AND name = 'tenant_id')
BEGIN
    ALTER TABLE [dbo].[cms_admin_role_value] ADD [tenant_id] BIGINT NOT NULL DEFAULT 0;
    ALTER TABLE [dbo].[cms_admin_role_value] ADD [deleted] BIT NOT NULL DEFAULT 0;
    CREATE INDEX [idx_cms_admin_role_value_tenant] ON [dbo].[cms_admin_role_value]([tenant_id]);
    CREATE INDEX [idx_cms_admin_role_value_deleted] ON [dbo].[cms_admin_role_value]([deleted]);
END
GO

IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID(N'[dbo].[cms_ads]') AND name = 'tenant_id')
BEGIN
    ALTER TABLE [dbo].[cms_ads] ADD [tenant_id] BIGINT NOT NULL DEFAULT 0;
    ALTER TABLE [dbo].[cms_ads] ADD [deleted] BIT NOT NULL DEFAULT 0;
    CREATE INDEX [idx_cms_ads_tenant] ON [dbo].[cms_ads]([tenant_id]);
    CREATE INDEX [idx_cms_ads_deleted] ON [dbo].[cms_ads]([deleted]);
END
GO

IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID(N'[dbo].[cms_ads_category]') AND name = 'tenant_id')
BEGIN
    ALTER TABLE [dbo].[cms_ads_category] ADD [tenant_id] BIGINT NOT NULL DEFAULT 0;
    ALTER TABLE [dbo].[cms_ads_category] ADD [deleted] BIT NOT NULL DEFAULT 0;
    CREATE INDEX [idx_cms_ads_category_tenant] ON [dbo].[cms_ads_category]([tenant_id]);
    CREATE INDEX [idx_cms_ads_category_deleted] ON [dbo].[cms_ads_category]([deleted]);
END
GO

IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID(N'[dbo].[cms_ads_category_relation]') AND name = 'tenant_id')
BEGIN
    ALTER TABLE [dbo].[cms_ads_category_relation] ADD [tenant_id] BIGINT NOT NULL DEFAULT 0;
    ALTER TABLE [dbo].[cms_ads_category_relation] ADD [deleted] BIT NOT NULL DEFAULT 0;
    CREATE INDEX [idx_cms_ads_category_relation_tenant] ON [dbo].[cms_ads_category_relation]([tenant_id]);
    CREATE INDEX [idx_cms_ads_category_relation_deleted] ON [dbo].[cms_ads_category_relation]([deleted]);
END
GO

IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID(N'[dbo].[cms_album]') AND name = 'tenant_id')
BEGIN
    ALTER TABLE [dbo].[cms_album] ADD [tenant_id] BIGINT NOT NULL DEFAULT 0;
    ALTER TABLE [dbo].[cms_album] ADD [deleted] BIT NOT NULL DEFAULT 0;
    CREATE INDEX [idx_cms_album_tenant] ON [dbo].[cms_album]([tenant_id]);
    CREATE INDEX [idx_cms_album_deleted] ON [dbo].[cms_album]([deleted]);
END
GO

IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID(N'[dbo].[cms_article]') AND name = 'tenant_id')
BEGIN
    ALTER TABLE [dbo].[cms_article] ADD [tenant_id] BIGINT NOT NULL DEFAULT 0;
    ALTER TABLE [dbo].[cms_article] ADD [deleted] BIT NOT NULL DEFAULT 0;
    CREATE INDEX [idx_cms_article_tenant] ON [dbo].[cms_article]([tenant_id]);
    CREATE INDEX [idx_cms_article_deleted] ON [dbo].[cms_article]([deleted]);
END
GO

IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID(N'[dbo].[cms_article_category]') AND name = 'tenant_id')
BEGIN
    ALTER TABLE [dbo].[cms_article_category] ADD [tenant_id] BIGINT NOT NULL DEFAULT 0;
    ALTER TABLE [dbo].[cms_article_category] ADD [deleted] BIT NOT NULL DEFAULT 0;
    CREATE INDEX [idx_cms_article_category_tenant] ON [dbo].[cms_article_category]([tenant_id]);
    CREATE INDEX [idx_cms_article_category_deleted] ON [dbo].[cms_article_category]([deleted]);
END
GO

IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID(N'[dbo].[cms_article_category_relation]') AND name = 'tenant_id')
BEGIN
    ALTER TABLE [dbo].[cms_article_category_relation] ADD [tenant_id] BIGINT NOT NULL DEFAULT 0;
    ALTER TABLE [dbo].[cms_article_category_relation] ADD [deleted] BIT NOT NULL DEFAULT 0;
    CREATE INDEX [idx_cms_article_category_relation_tenant] ON [dbo].[cms_article_category_relation]([tenant_id]);
    CREATE INDEX [idx_cms_article_category_relation_deleted] ON [dbo].[cms_article_category_relation]([deleted]);
END
GO

IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID(N'[dbo].[cms_article_comment]') AND name = 'tenant_id')
BEGIN
    ALTER TABLE [dbo].[cms_article_comment] ADD [tenant_id] BIGINT NOT NULL DEFAULT 0;
    ALTER TABLE [dbo].[cms_article_comment] ADD [deleted] BIT NOT NULL DEFAULT 0;
    CREATE INDEX [idx_cms_article_comment_tenant] ON [dbo].[cms_article_comment]([tenant_id]);
    CREATE INDEX [idx_cms_article_comment_deleted] ON [dbo].[cms_article_comment]([deleted]);
END
GO

IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID(N'[dbo].[cms_article_label]') AND name = 'tenant_id')
BEGIN
    ALTER TABLE [dbo].[cms_article_label] ADD [tenant_id] BIGINT NOT NULL DEFAULT 0;
    ALTER TABLE [dbo].[cms_article_label] ADD [deleted] BIT NOT NULL DEFAULT 0;
    CREATE INDEX [idx_cms_article_label_tenant] ON [dbo].[cms_article_label]([tenant_id]);
    CREATE INDEX [idx_cms_article_label_deleted] ON [dbo].[cms_article_label]([deleted]);
END
GO

IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID(N'[dbo].[cms_article_label_relation]') AND name = 'tenant_id')
BEGIN
    ALTER TABLE [dbo].[cms_article_label_relation] ADD [tenant_id] BIGINT NOT NULL DEFAULT 0;
    ALTER TABLE [dbo].[cms_article_label_relation] ADD [deleted] BIT NOT NULL DEFAULT 0;
    CREATE INDEX [idx_cms_article_label_relation_tenant] ON [dbo].[cms_article_label_relation]([tenant_id]);
    CREATE INDEX [idx_cms_article_label_relation_deleted] ON [dbo].[cms_article_label_relation]([deleted]);
END
GO

IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID(N'[dbo].[cms_article_property]') AND name = 'tenant_id')
BEGIN
    ALTER TABLE [dbo].[cms_article_property] ADD [tenant_id] BIGINT NOT NULL DEFAULT 0;
    ALTER TABLE [dbo].[cms_article_property] ADD [deleted] BIT NOT NULL DEFAULT 0;
    CREATE INDEX [idx_cms_article_property_tenant] ON [dbo].[cms_article_property]([tenant_id]);
    CREATE INDEX [idx_cms_article_property_deleted] ON [dbo].[cms_article_property]([deleted]);
END
GO

IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID(N'[dbo].[cms_attach]') AND name = 'tenant_id')
BEGIN
    ALTER TABLE [dbo].[cms_attach] ADD [tenant_id] BIGINT NOT NULL DEFAULT 0;
    ALTER TABLE [dbo].[cms_attach] ADD [deleted] BIT NOT NULL DEFAULT 0;
    CREATE INDEX [idx_cms_attach_tenant] ON [dbo].[cms_attach]([tenant_id]);
    CREATE INDEX [idx_cms_attach_deleted] ON [dbo].[cms_attach]([deleted]);
END
GO

IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID(N'[dbo].[cms_link]') AND name = 'tenant_id')
BEGIN
    ALTER TABLE [dbo].[cms_link] ADD [tenant_id] BIGINT NOT NULL DEFAULT 0;
    ALTER TABLE [dbo].[cms_link] ADD [deleted] BIT NOT NULL DEFAULT 0;
    CREATE INDEX [idx_cms_link_tenant] ON [dbo].[cms_link]([tenant_id]);
    CREATE INDEX [idx_cms_link_deleted] ON [dbo].[cms_link]([deleted]);
END
GO

IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID(N'[dbo].[cms_link_category]') AND name = 'tenant_id')
BEGIN
    ALTER TABLE [dbo].[cms_link_category] ADD [tenant_id] BIGINT NOT NULL DEFAULT 0;
    ALTER TABLE [dbo].[cms_link_category] ADD [deleted] BIT NOT NULL DEFAULT 0;
    CREATE INDEX [idx_cms_link_category_tenant] ON [dbo].[cms_link_category]([tenant_id]);
    CREATE INDEX [idx_cms_link_category_deleted] ON [dbo].[cms_link_category]([deleted]);
END
GO

IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID(N'[dbo].[cms_link_category_relation]') AND name = 'tenant_id')
BEGIN
    ALTER TABLE [dbo].[cms_link_category_relation] ADD [tenant_id] BIGINT NOT NULL DEFAULT 0;
    ALTER TABLE [dbo].[cms_link_category_relation] ADD [deleted] BIT NOT NULL DEFAULT 0;
    CREATE INDEX [idx_cms_link_category_relation_tenant] ON [dbo].[cms_link_category_relation]([tenant_id]);
    CREATE INDEX [idx_cms_link_category_relation_deleted] ON [dbo].[cms_link_category_relation]([deleted]);
END
GO

IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID(N'[dbo].[cms_site]') AND name = 'tenant_id')
BEGIN
    ALTER TABLE [dbo].[cms_site] ADD [tenant_id] BIGINT NOT NULL DEFAULT 0;
    ALTER TABLE [dbo].[cms_site] ADD [deleted] BIT NOT NULL DEFAULT 0;
    CREATE INDEX [idx_cms_site_tenant] ON [dbo].[cms_site]([tenant_id]);
    CREATE INDEX [idx_cms_site_deleted] ON [dbo].[cms_site]([deleted]);
END
GO

IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID(N'[dbo].[cms_site_channel]') AND name = 'tenant_id')
BEGIN
    ALTER TABLE [dbo].[cms_site_channel] ADD [tenant_id] BIGINT NOT NULL DEFAULT 0;
    ALTER TABLE [dbo].[cms_site_channel] ADD [deleted] BIT NOT NULL DEFAULT 0;
    CREATE INDEX [idx_cms_site_channel_tenant] ON [dbo].[cms_site_channel]([tenant_id]);
    CREATE INDEX [idx_cms_site_channel_deleted] ON [dbo].[cms_site_channel]([deleted]);
END
GO

IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID(N'[dbo].[cms_site_channel_album]') AND name = 'tenant_id')
BEGIN
    ALTER TABLE [dbo].[cms_site_channel_album] ADD [tenant_id] BIGINT NOT NULL DEFAULT 0;
    ALTER TABLE [dbo].[cms_site_channel_album] ADD [deleted] BIT NOT NULL DEFAULT 0;
    CREATE INDEX [idx_cms_site_channel_album_tenant] ON [dbo].[cms_site_channel_album]([tenant_id]);
    CREATE INDEX [idx_cms_site_channel_album_deleted] ON [dbo].[cms_site_channel_album]([deleted]);
END
GO

IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID(N'[dbo].[cms_site_channel_field]') AND name = 'tenant_id')
BEGIN
    ALTER TABLE [dbo].[cms_site_channel_field] ADD [tenant_id] BIGINT NOT NULL DEFAULT 0;
    ALTER TABLE [dbo].[cms_site_channel_field] ADD [deleted] BIT NOT NULL DEFAULT 0;
    CREATE INDEX [idx_cms_site_channel_field_tenant] ON [dbo].[cms_site_channel_field]([tenant_id]);
    CREATE INDEX [idx_cms_site_channel_field_deleted] ON [dbo].[cms_site_channel_field]([deleted]);
END
GO

IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID(N'[dbo].[cms_site_domain]') AND name = 'tenant_id')
BEGIN
    ALTER TABLE [dbo].[cms_site_domain] ADD [tenant_id] BIGINT NOT NULL DEFAULT 0;
    ALTER TABLE [dbo].[cms_site_domain] ADD [deleted] BIT NOT NULL DEFAULT 0;
    CREATE INDEX [idx_cms_site_domain_tenant] ON [dbo].[cms_site_domain]([tenant_id]);
    CREATE INDEX [idx_cms_site_domain_deleted] ON [dbo].[cms_site_domain]([deleted]);
END
GO

IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID(N'[dbo].[cms_tag]') AND name = 'tenant_id')
BEGIN
    ALTER TABLE [dbo].[cms_tag] ADD [tenant_id] BIGINT NOT NULL DEFAULT 0;
    ALTER TABLE [dbo].[cms_tag] ADD [deleted] BIT NOT NULL DEFAULT 0;
    CREATE INDEX [idx_cms_tag_tenant] ON [dbo].[cms_tag]([tenant_id]);
    CREATE INDEX [idx_cms_tag_deleted] ON [dbo].[cms_tag]([deleted]);
END
GO

IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID(N'[dbo].[cms_theme]') AND name = 'tenant_id')
BEGIN
    ALTER TABLE [dbo].[cms_theme] ADD [tenant_id] BIGINT NOT NULL DEFAULT 0;
    ALTER TABLE [dbo].[cms_theme] ADD [deleted] BIT NOT NULL DEFAULT 0;
    CREATE INDEX [idx_cms_theme_tenant] ON [dbo].[cms_theme]([tenant_id]);
    CREATE INDEX [idx_cms_theme_deleted] ON [dbo].[cms_theme]([deleted]);
END
GO

IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID(N'[dbo].[cms_topic]') AND name = 'tenant_id')
BEGIN
    ALTER TABLE [dbo].[cms_topic] ADD [tenant_id] BIGINT NOT NULL DEFAULT 0;
    ALTER TABLE [dbo].[cms_topic] ADD [deleted] BIT NOT NULL DEFAULT 0;
    CREATE INDEX [idx_cms_topic_tenant] ON [dbo].[cms_topic]([tenant_id]);
    CREATE INDEX [idx_cms_topic_deleted] ON [dbo].[cms_topic]([deleted]);
END
GO

IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID(N'[dbo].[plg_online_register]') AND name = 'tenant_id')
BEGIN
    ALTER TABLE [dbo].[plg_online_register] ADD [tenant_id] BIGINT NOT NULL DEFAULT 0;
    ALTER TABLE [dbo].[plg_online_register] ADD [deleted] BIT NOT NULL DEFAULT 0;
    CREATE INDEX [idx_plg_online_register_tenant] ON [dbo].[plg_online_register]([tenant_id]);
    CREATE INDEX [idx_plg_online_register_deleted] ON [dbo].[plg_online_register]([deleted]);
END
GO

IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID(N'[dbo].[weixin_account]') AND name = 'tenant_id')
BEGIN
    ALTER TABLE [dbo].[weixin_account] ADD [tenant_id] BIGINT NOT NULL DEFAULT 0;
    ALTER TABLE [dbo].[weixin_account] ADD [deleted] BIT NOT NULL DEFAULT 0;
    CREATE INDEX [idx_weixin_account_tenant] ON [dbo].[weixin_account]([tenant_id]);
    CREATE INDEX [idx_weixin_account_deleted] ON [dbo].[weixin_account]([deleted]);
END
GO

IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID(N'[dbo].[weixin_menu]') AND name = 'tenant_id')
BEGIN
    ALTER TABLE [dbo].[weixin_menu] ADD [tenant_id] BIGINT NOT NULL DEFAULT 0;
    ALTER TABLE [dbo].[weixin_menu] ADD [deleted] BIT NOT NULL DEFAULT 0;
    CREATE INDEX [idx_weixin_menu_tenant] ON [dbo].[weixin_menu]([tenant_id]);
    CREATE INDEX [idx_weixin_menu_deleted] ON [dbo].[weixin_menu]([deleted]);
END
GO

IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID(N'[dbo].[weixin_mp_verify]') AND name = 'tenant_id')
BEGIN
    ALTER TABLE [dbo].[weixin_mp_verify] ADD [tenant_id] BIGINT NOT NULL DEFAULT 0;
    ALTER TABLE [dbo].[weixin_mp_verify] ADD [deleted] BIT NOT NULL DEFAULT 0;
    CREATE INDEX [idx_weixin_mp_verify_tenant] ON [dbo].[weixin_mp_verify]([tenant_id]);
    CREATE INDEX [idx_weixin_mp_verify_deleted] ON [dbo].[weixin_mp_verify]([deleted]);
END
GO

IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID(N'[dbo].[weixin_request_content]') AND name = 'tenant_id')
BEGIN
    ALTER TABLE [dbo].[weixin_request_content] ADD [tenant_id] BIGINT NOT NULL DEFAULT 0;
    ALTER TABLE [dbo].[weixin_request_content] ADD [deleted] BIT NOT NULL DEFAULT 0;
    CREATE INDEX [idx_weixin_request_content_tenant] ON [dbo].[weixin_request_content]([tenant_id]);
    CREATE INDEX [idx_weixin_request_content_deleted] ON [dbo].[weixin_request_content]([deleted]);
END
GO

IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID(N'[dbo].[weixin_request_rule]') AND name = 'tenant_id')
BEGIN
    ALTER TABLE [dbo].[weixin_request_rule] ADD [tenant_id] BIGINT NOT NULL DEFAULT 0;
    ALTER TABLE [dbo].[weixin_request_rule] ADD [deleted] BIT NOT NULL DEFAULT 0;
    CREATE INDEX [idx_weixin_request_rule_tenant] ON [dbo].[weixin_request_rule]([tenant_id]);
    CREATE INDEX [idx_weixin_request_rule_deleted] ON [dbo].[weixin_request_rule]([deleted]);
END
GO

IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID(N'[dbo].[weixin_response_content]') AND name = 'tenant_id')
BEGIN
    ALTER TABLE [dbo].[weixin_response_content] ADD [tenant_id] BIGINT NOT NULL DEFAULT 0;
    ALTER TABLE [dbo].[weixin_response_content] ADD [deleted] BIT NOT NULL DEFAULT 0;
    CREATE INDEX [idx_weixin_response_content_tenant] ON [dbo].[weixin_response_content]([tenant_id]);
    CREATE INDEX [idx_weixin_response_content_deleted] ON [dbo].[weixin_response_content]([deleted]);
END
GO

-- ========== 4. 初始化默认租户 ==========
IF NOT EXISTS (SELECT 1 FROM [dbo].[cms_tenant] WHERE [code] = 'default')
BEGIN
    INSERT INTO [dbo].[cms_tenant] ([code], [name], [short_name], [status], [create_name])
    VALUES (N'default', N'默认租户', N'默认', 1, N'system');
END
GO