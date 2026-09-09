-- ============================================================
-- 多租户 + 软删除 字段迁移脚本
-- 数据库: PostgreSQL / openGauss / KingbaseES
-- 设计参考: 芋道 yudao BaseDO
-- 生成时间: 2026-09-09 10:28:57
-- ============================================================

\connect cms

-- ========== 1. 租户主表 cms_tenant ==========
CREATE TABLE IF NOT EXISTS cms_tenant (
  tenant_id bigserial PRIMARY KEY,
  code varchar(64) NOT NULL,
  name varchar(128) NOT NULL,
  short_name varchar(64),
  contact_name varchar(64),
  contact_tel varchar(32),
  contact_mail varchar(128),
  description varchar(512),
  logo varchar(512),
  domain varchar(256),
  expire_time timestamp,
  max_users int DEFAULT 0,
  status smallint DEFAULT 1,
  create_id int,
  create_name varchar(64),
  create_time timestamp DEFAULT CURRENT_TIMESTAMP,
  update_id int,
  update_name varchar(64),
  update_time timestamp,
  deleted boolean DEFAULT false,
  CONSTRAINT uk_cms_tenant_code UNIQUE (code)
);
COMMENT ON TABLE cms_tenant IS '租户表(多租户主表)';
CREATE INDEX IF NOT EXISTS idx_cms_tenant_status ON cms_tenant(status);
CREATE INDEX IF NOT EXISTS idx_cms_tenant_deleted ON cms_tenant(deleted);

-- ========== 2. 给 cms_admin 添加 tenant_id 和 deleted ==========
ALTER TABLE cms_admin ADD COLUMN IF NOT EXISTS tenant_id bigint NOT NULL DEFAULT 0;
CREATE INDEX IF NOT EXISTS idx_cms_admin_tenant ON cms_admin(tenant_id);
ALTER TABLE cms_admin ADD COLUMN IF NOT EXISTS deleted boolean DEFAULT false;
CREATE INDEX IF NOT EXISTS idx_cms_admin_deleted ON cms_admin(deleted);

-- ========== 3. 给所有业务表添加列 ==========
ALTER TABLE cms_ad_category_relation
  ADD COLUMN IF NOT EXISTS tenant_id bigint NOT NULL DEFAULT 0;
ALTER TABLE cms_ad_category_relation
  ADD COLUMN IF NOT EXISTS deleted boolean DEFAULT false;
CREATE INDEX IF NOT EXISTS idx_cms_ad_category_relation_tenant ON cms_ad_category_relation(tenant_id);
CREATE INDEX IF NOT EXISTS idx_cms_ad_category_relation_deleted ON cms_ad_category_relation(deleted);

ALTER TABLE cms_admin_log
  ADD COLUMN IF NOT EXISTS tenant_id bigint NOT NULL DEFAULT 0;
ALTER TABLE cms_admin_log
  ADD COLUMN IF NOT EXISTS deleted boolean DEFAULT false;
CREATE INDEX IF NOT EXISTS idx_cms_admin_log_tenant ON cms_admin_log(tenant_id);
CREATE INDEX IF NOT EXISTS idx_cms_admin_log_deleted ON cms_admin_log(deleted);

ALTER TABLE cms_admin_nav
  ADD COLUMN IF NOT EXISTS tenant_id bigint NOT NULL DEFAULT 0;
ALTER TABLE cms_admin_nav
  ADD COLUMN IF NOT EXISTS deleted boolean DEFAULT false;
CREATE INDEX IF NOT EXISTS idx_cms_admin_nav_tenant ON cms_admin_nav(tenant_id);
CREATE INDEX IF NOT EXISTS idx_cms_admin_nav_deleted ON cms_admin_nav(deleted);

ALTER TABLE cms_admin_notice
  ADD COLUMN IF NOT EXISTS tenant_id bigint NOT NULL DEFAULT 0;
ALTER TABLE cms_admin_notice
  ADD COLUMN IF NOT EXISTS deleted boolean DEFAULT false;
CREATE INDEX IF NOT EXISTS idx_cms_admin_notice_tenant ON cms_admin_notice(tenant_id);
CREATE INDEX IF NOT EXISTS idx_cms_admin_notice_deleted ON cms_admin_notice(deleted);

ALTER TABLE cms_admin_role
  ADD COLUMN IF NOT EXISTS tenant_id bigint NOT NULL DEFAULT 0;
ALTER TABLE cms_admin_role
  ADD COLUMN IF NOT EXISTS deleted boolean DEFAULT false;
CREATE INDEX IF NOT EXISTS idx_cms_admin_role_tenant ON cms_admin_role(tenant_id);
CREATE INDEX IF NOT EXISTS idx_cms_admin_role_deleted ON cms_admin_role(deleted);

ALTER TABLE cms_admin_role_site
  ADD COLUMN IF NOT EXISTS tenant_id bigint NOT NULL DEFAULT 0;
ALTER TABLE cms_admin_role_site
  ADD COLUMN IF NOT EXISTS deleted boolean DEFAULT false;
CREATE INDEX IF NOT EXISTS idx_cms_admin_role_site_tenant ON cms_admin_role_site(tenant_id);
CREATE INDEX IF NOT EXISTS idx_cms_admin_role_site_deleted ON cms_admin_role_site(deleted);

ALTER TABLE cms_admin_role_value
  ADD COLUMN IF NOT EXISTS tenant_id bigint NOT NULL DEFAULT 0;
ALTER TABLE cms_admin_role_value
  ADD COLUMN IF NOT EXISTS deleted boolean DEFAULT false;
CREATE INDEX IF NOT EXISTS idx_cms_admin_role_value_tenant ON cms_admin_role_value(tenant_id);
CREATE INDEX IF NOT EXISTS idx_cms_admin_role_value_deleted ON cms_admin_role_value(deleted);

ALTER TABLE cms_ads
  ADD COLUMN IF NOT EXISTS tenant_id bigint NOT NULL DEFAULT 0;
ALTER TABLE cms_ads
  ADD COLUMN IF NOT EXISTS deleted boolean DEFAULT false;
CREATE INDEX IF NOT EXISTS idx_cms_ads_tenant ON cms_ads(tenant_id);
CREATE INDEX IF NOT EXISTS idx_cms_ads_deleted ON cms_ads(deleted);

ALTER TABLE cms_ads_category
  ADD COLUMN IF NOT EXISTS tenant_id bigint NOT NULL DEFAULT 0;
ALTER TABLE cms_ads_category
  ADD COLUMN IF NOT EXISTS deleted boolean DEFAULT false;
CREATE INDEX IF NOT EXISTS idx_cms_ads_category_tenant ON cms_ads_category(tenant_id);
CREATE INDEX IF NOT EXISTS idx_cms_ads_category_deleted ON cms_ads_category(deleted);

ALTER TABLE cms_ads_category_relation
  ADD COLUMN IF NOT EXISTS tenant_id bigint NOT NULL DEFAULT 0;
ALTER TABLE cms_ads_category_relation
  ADD COLUMN IF NOT EXISTS deleted boolean DEFAULT false;
CREATE INDEX IF NOT EXISTS idx_cms_ads_category_relation_tenant ON cms_ads_category_relation(tenant_id);
CREATE INDEX IF NOT EXISTS idx_cms_ads_category_relation_deleted ON cms_ads_category_relation(deleted);

ALTER TABLE cms_album
  ADD COLUMN IF NOT EXISTS tenant_id bigint NOT NULL DEFAULT 0;
ALTER TABLE cms_album
  ADD COLUMN IF NOT EXISTS deleted boolean DEFAULT false;
CREATE INDEX IF NOT EXISTS idx_cms_album_tenant ON cms_album(tenant_id);
CREATE INDEX IF NOT EXISTS idx_cms_album_deleted ON cms_album(deleted);

ALTER TABLE cms_article
  ADD COLUMN IF NOT EXISTS tenant_id bigint NOT NULL DEFAULT 0;
ALTER TABLE cms_article
  ADD COLUMN IF NOT EXISTS deleted boolean DEFAULT false;
CREATE INDEX IF NOT EXISTS idx_cms_article_tenant ON cms_article(tenant_id);
CREATE INDEX IF NOT EXISTS idx_cms_article_deleted ON cms_article(deleted);

ALTER TABLE cms_article_category
  ADD COLUMN IF NOT EXISTS tenant_id bigint NOT NULL DEFAULT 0;
ALTER TABLE cms_article_category
  ADD COLUMN IF NOT EXISTS deleted boolean DEFAULT false;
CREATE INDEX IF NOT EXISTS idx_cms_article_category_tenant ON cms_article_category(tenant_id);
CREATE INDEX IF NOT EXISTS idx_cms_article_category_deleted ON cms_article_category(deleted);

ALTER TABLE cms_article_category_relation
  ADD COLUMN IF NOT EXISTS tenant_id bigint NOT NULL DEFAULT 0;
ALTER TABLE cms_article_category_relation
  ADD COLUMN IF NOT EXISTS deleted boolean DEFAULT false;
CREATE INDEX IF NOT EXISTS idx_cms_article_category_relation_tenant ON cms_article_category_relation(tenant_id);
CREATE INDEX IF NOT EXISTS idx_cms_article_category_relation_deleted ON cms_article_category_relation(deleted);

ALTER TABLE cms_article_comment
  ADD COLUMN IF NOT EXISTS tenant_id bigint NOT NULL DEFAULT 0;
ALTER TABLE cms_article_comment
  ADD COLUMN IF NOT EXISTS deleted boolean DEFAULT false;
CREATE INDEX IF NOT EXISTS idx_cms_article_comment_tenant ON cms_article_comment(tenant_id);
CREATE INDEX IF NOT EXISTS idx_cms_article_comment_deleted ON cms_article_comment(deleted);

ALTER TABLE cms_article_label
  ADD COLUMN IF NOT EXISTS tenant_id bigint NOT NULL DEFAULT 0;
ALTER TABLE cms_article_label
  ADD COLUMN IF NOT EXISTS deleted boolean DEFAULT false;
CREATE INDEX IF NOT EXISTS idx_cms_article_label_tenant ON cms_article_label(tenant_id);
CREATE INDEX IF NOT EXISTS idx_cms_article_label_deleted ON cms_article_label(deleted);

ALTER TABLE cms_article_label_relation
  ADD COLUMN IF NOT EXISTS tenant_id bigint NOT NULL DEFAULT 0;
ALTER TABLE cms_article_label_relation
  ADD COLUMN IF NOT EXISTS deleted boolean DEFAULT false;
CREATE INDEX IF NOT EXISTS idx_cms_article_label_relation_tenant ON cms_article_label_relation(tenant_id);
CREATE INDEX IF NOT EXISTS idx_cms_article_label_relation_deleted ON cms_article_label_relation(deleted);

ALTER TABLE cms_article_property
  ADD COLUMN IF NOT EXISTS tenant_id bigint NOT NULL DEFAULT 0;
ALTER TABLE cms_article_property
  ADD COLUMN IF NOT EXISTS deleted boolean DEFAULT false;
CREATE INDEX IF NOT EXISTS idx_cms_article_property_tenant ON cms_article_property(tenant_id);
CREATE INDEX IF NOT EXISTS idx_cms_article_property_deleted ON cms_article_property(deleted);

ALTER TABLE cms_attach
  ADD COLUMN IF NOT EXISTS tenant_id bigint NOT NULL DEFAULT 0;
ALTER TABLE cms_attach
  ADD COLUMN IF NOT EXISTS deleted boolean DEFAULT false;
CREATE INDEX IF NOT EXISTS idx_cms_attach_tenant ON cms_attach(tenant_id);
CREATE INDEX IF NOT EXISTS idx_cms_attach_deleted ON cms_attach(deleted);

ALTER TABLE cms_link
  ADD COLUMN IF NOT EXISTS tenant_id bigint NOT NULL DEFAULT 0;
ALTER TABLE cms_link
  ADD COLUMN IF NOT EXISTS deleted boolean DEFAULT false;
CREATE INDEX IF NOT EXISTS idx_cms_link_tenant ON cms_link(tenant_id);
CREATE INDEX IF NOT EXISTS idx_cms_link_deleted ON cms_link(deleted);

ALTER TABLE cms_link_category
  ADD COLUMN IF NOT EXISTS tenant_id bigint NOT NULL DEFAULT 0;
ALTER TABLE cms_link_category
  ADD COLUMN IF NOT EXISTS deleted boolean DEFAULT false;
CREATE INDEX IF NOT EXISTS idx_cms_link_category_tenant ON cms_link_category(tenant_id);
CREATE INDEX IF NOT EXISTS idx_cms_link_category_deleted ON cms_link_category(deleted);

ALTER TABLE cms_link_category_relation
  ADD COLUMN IF NOT EXISTS tenant_id bigint NOT NULL DEFAULT 0;
ALTER TABLE cms_link_category_relation
  ADD COLUMN IF NOT EXISTS deleted boolean DEFAULT false;
CREATE INDEX IF NOT EXISTS idx_cms_link_category_relation_tenant ON cms_link_category_relation(tenant_id);
CREATE INDEX IF NOT EXISTS idx_cms_link_category_relation_deleted ON cms_link_category_relation(deleted);

ALTER TABLE cms_site
  ADD COLUMN IF NOT EXISTS tenant_id bigint NOT NULL DEFAULT 0;
ALTER TABLE cms_site
  ADD COLUMN IF NOT EXISTS deleted boolean DEFAULT false;
CREATE INDEX IF NOT EXISTS idx_cms_site_tenant ON cms_site(tenant_id);
CREATE INDEX IF NOT EXISTS idx_cms_site_deleted ON cms_site(deleted);

ALTER TABLE cms_site_channel
  ADD COLUMN IF NOT EXISTS tenant_id bigint NOT NULL DEFAULT 0;
ALTER TABLE cms_site_channel
  ADD COLUMN IF NOT EXISTS deleted boolean DEFAULT false;
CREATE INDEX IF NOT EXISTS idx_cms_site_channel_tenant ON cms_site_channel(tenant_id);
CREATE INDEX IF NOT EXISTS idx_cms_site_channel_deleted ON cms_site_channel(deleted);

ALTER TABLE cms_site_channel_album
  ADD COLUMN IF NOT EXISTS tenant_id bigint NOT NULL DEFAULT 0;
ALTER TABLE cms_site_channel_album
  ADD COLUMN IF NOT EXISTS deleted boolean DEFAULT false;
CREATE INDEX IF NOT EXISTS idx_cms_site_channel_album_tenant ON cms_site_channel_album(tenant_id);
CREATE INDEX IF NOT EXISTS idx_cms_site_channel_album_deleted ON cms_site_channel_album(deleted);

ALTER TABLE cms_site_channel_field
  ADD COLUMN IF NOT EXISTS tenant_id bigint NOT NULL DEFAULT 0;
ALTER TABLE cms_site_channel_field
  ADD COLUMN IF NOT EXISTS deleted boolean DEFAULT false;
CREATE INDEX IF NOT EXISTS idx_cms_site_channel_field_tenant ON cms_site_channel_field(tenant_id);
CREATE INDEX IF NOT EXISTS idx_cms_site_channel_field_deleted ON cms_site_channel_field(deleted);

ALTER TABLE cms_site_domain
  ADD COLUMN IF NOT EXISTS tenant_id bigint NOT NULL DEFAULT 0;
ALTER TABLE cms_site_domain
  ADD COLUMN IF NOT EXISTS deleted boolean DEFAULT false;
CREATE INDEX IF NOT EXISTS idx_cms_site_domain_tenant ON cms_site_domain(tenant_id);
CREATE INDEX IF NOT EXISTS idx_cms_site_domain_deleted ON cms_site_domain(deleted);

ALTER TABLE cms_tag
  ADD COLUMN IF NOT EXISTS tenant_id bigint NOT NULL DEFAULT 0;
ALTER TABLE cms_tag
  ADD COLUMN IF NOT EXISTS deleted boolean DEFAULT false;
CREATE INDEX IF NOT EXISTS idx_cms_tag_tenant ON cms_tag(tenant_id);
CREATE INDEX IF NOT EXISTS idx_cms_tag_deleted ON cms_tag(deleted);

ALTER TABLE cms_theme
  ADD COLUMN IF NOT EXISTS tenant_id bigint NOT NULL DEFAULT 0;
ALTER TABLE cms_theme
  ADD COLUMN IF NOT EXISTS deleted boolean DEFAULT false;
CREATE INDEX IF NOT EXISTS idx_cms_theme_tenant ON cms_theme(tenant_id);
CREATE INDEX IF NOT EXISTS idx_cms_theme_deleted ON cms_theme(deleted);

ALTER TABLE cms_topic
  ADD COLUMN IF NOT EXISTS tenant_id bigint NOT NULL DEFAULT 0;
ALTER TABLE cms_topic
  ADD COLUMN IF NOT EXISTS deleted boolean DEFAULT false;
CREATE INDEX IF NOT EXISTS idx_cms_topic_tenant ON cms_topic(tenant_id);
CREATE INDEX IF NOT EXISTS idx_cms_topic_deleted ON cms_topic(deleted);

ALTER TABLE plg_online_register
  ADD COLUMN IF NOT EXISTS tenant_id bigint NOT NULL DEFAULT 0;
ALTER TABLE plg_online_register
  ADD COLUMN IF NOT EXISTS deleted boolean DEFAULT false;
CREATE INDEX IF NOT EXISTS idx_plg_online_register_tenant ON plg_online_register(tenant_id);
CREATE INDEX IF NOT EXISTS idx_plg_online_register_deleted ON plg_online_register(deleted);

ALTER TABLE weixin_account
  ADD COLUMN IF NOT EXISTS tenant_id bigint NOT NULL DEFAULT 0;
ALTER TABLE weixin_account
  ADD COLUMN IF NOT EXISTS deleted boolean DEFAULT false;
CREATE INDEX IF NOT EXISTS idx_weixin_account_tenant ON weixin_account(tenant_id);
CREATE INDEX IF NOT EXISTS idx_weixin_account_deleted ON weixin_account(deleted);

ALTER TABLE weixin_menu
  ADD COLUMN IF NOT EXISTS tenant_id bigint NOT NULL DEFAULT 0;
ALTER TABLE weixin_menu
  ADD COLUMN IF NOT EXISTS deleted boolean DEFAULT false;
CREATE INDEX IF NOT EXISTS idx_weixin_menu_tenant ON weixin_menu(tenant_id);
CREATE INDEX IF NOT EXISTS idx_weixin_menu_deleted ON weixin_menu(deleted);

ALTER TABLE weixin_mp_verify
  ADD COLUMN IF NOT EXISTS tenant_id bigint NOT NULL DEFAULT 0;
ALTER TABLE weixin_mp_verify
  ADD COLUMN IF NOT EXISTS deleted boolean DEFAULT false;
CREATE INDEX IF NOT EXISTS idx_weixin_mp_verify_tenant ON weixin_mp_verify(tenant_id);
CREATE INDEX IF NOT EXISTS idx_weixin_mp_verify_deleted ON weixin_mp_verify(deleted);

ALTER TABLE weixin_request_content
  ADD COLUMN IF NOT EXISTS tenant_id bigint NOT NULL DEFAULT 0;
ALTER TABLE weixin_request_content
  ADD COLUMN IF NOT EXISTS deleted boolean DEFAULT false;
CREATE INDEX IF NOT EXISTS idx_weixin_request_content_tenant ON weixin_request_content(tenant_id);
CREATE INDEX IF NOT EXISTS idx_weixin_request_content_deleted ON weixin_request_content(deleted);

ALTER TABLE weixin_request_rule
  ADD COLUMN IF NOT EXISTS tenant_id bigint NOT NULL DEFAULT 0;
ALTER TABLE weixin_request_rule
  ADD COLUMN IF NOT EXISTS deleted boolean DEFAULT false;
CREATE INDEX IF NOT EXISTS idx_weixin_request_rule_tenant ON weixin_request_rule(tenant_id);
CREATE INDEX IF NOT EXISTS idx_weixin_request_rule_deleted ON weixin_request_rule(deleted);

ALTER TABLE weixin_response_content
  ADD COLUMN IF NOT EXISTS tenant_id bigint NOT NULL DEFAULT 0;
ALTER TABLE weixin_response_content
  ADD COLUMN IF NOT EXISTS deleted boolean DEFAULT false;
CREATE INDEX IF NOT EXISTS idx_weixin_response_content_tenant ON weixin_response_content(tenant_id);
CREATE INDEX IF NOT EXISTS idx_weixin_response_content_deleted ON weixin_response_content(deleted);

-- ========== 4. 初始化默认租户 ==========
INSERT INTO cms_tenant (code, name, short_name, status, create_name)
VALUES ('default', '默认租户', '默认', 1, 'system')
ON CONFLICT (code) DO NOTHING;