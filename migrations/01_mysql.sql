-- ============================================================
-- 多租户 + 软删除 字段迁移脚本
-- 数据库: MySQL / MariaDB
-- 设计参考: 芋道 yudao BaseDO
-- 生成时间: 2026-09-09 10:28:57
-- 影响表数: 39
-- ============================================================

-- 字段说明:
--   tenant_id  bigint NOT NULL DEFAULT 0 - 多租户ID(参考芋道 TenantBaseDO)
--   deleted tinyint(1) NOT NULL DEFAULT 0 - 软删除标记(参考芋道 BaseDO.deleted)

-- 注意事项:
--   1) 不再使用 is_deleted 字段,与芋道 yudao 保持一致
--   2) deleted 字段为 tinyint(1) 类型,0=未删除,1=已删除
--   3) 建议在执行前备份数据库

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ========== 1. 租户主表 cms_tenant ==========
CREATE TABLE IF NOT EXISTS `cms_tenant` (
  `tenant_id` bigint NOT NULL AUTO_INCREMENT COMMENT '租户ID',
  `code` varchar(64) NOT NULL COMMENT '租户编码(全局唯一)',
  `name` varchar(128) NOT NULL COMMENT '租户名称',
  `short_name` varchar(64) DEFAULT NULL COMMENT '租户简称',
  `contact_name` varchar(64) DEFAULT NULL COMMENT '联系人姓名',
  `contact_tel` varchar(32) DEFAULT NULL COMMENT '联系电话',
  `contact_mail` varchar(128) DEFAULT NULL COMMENT '联系邮箱',
  `description` varchar(512) DEFAULT NULL COMMENT '租户描述',
  `logo` varchar(512) DEFAULT NULL COMMENT '租户LOGO',
  `domain` varchar(256) DEFAULT NULL COMMENT '独立域名',
  `expire_time` datetime DEFAULT NULL COMMENT '到期时间',
  `max_users` int DEFAULT '0' COMMENT '最大用户数(0=不限)',
  `status` tinyint DEFAULT '1' COMMENT '状态0禁用1启用',
  `create_id` int DEFAULT NULL COMMENT '创建人ID',
  `create_name` varchar(64) DEFAULT NULL COMMENT '创建人姓名',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_id` int DEFAULT NULL COMMENT '更新人ID',
  `update_name` varchar(64) DEFAULT NULL COMMENT '更新人姓名',
  `update_time` datetime DEFAULT NULL COMMENT '修改时间',
  `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)',
  PRIMARY KEY (`tenant_id`),
  UNIQUE KEY `uk_code` (`code`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted` (`deleted`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='租户表(多租户主表)';

-- ========== 2. 给 cms_admin 添加 tenant_id 和 deleted ==========
-- (cms_admin 原本有 is_deleted 字段,需要先删除或迁移到 deleted)
ALTER TABLE `cms_admin`
  ADD COLUMN `tenant_id` bigint NOT NULL DEFAULT 0 COMMENT '租户ID' AFTER `update_time`,
  ADD INDEX `idx_tenant` (`tenant_id`);

-- 如果原表已有 is_deleted 字段,迁移数据到 deleted 后删除 is_deleted
ALTER TABLE `cms_admin`
  ADD COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)' AFTER `tenant_id`,
  ADD INDEX `idx_deleted` (`deleted`);

-- 数据迁移: 把 is_deleted 的值复制到 deleted(若 is_deleted 存在)
UPDATE `cms_admin` SET `deleted` = 1 WHERE `is_deleted` = 1;

-- 删除旧的 is_deleted 字段
ALTER TABLE `cms_admin` DROP COLUMN `is_deleted`;

-- ========== 3. 给所有业务表添加 tenant_id 和 deleted 列 ==========
-- cms_ad_category_relation
ALTER TABLE `cms_ad_category_relation`
  ADD COLUMN `tenant_id` bigint NOT NULL DEFAULT 0 COMMENT '租户ID(多租户隔离)',
  ADD COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)',
  ADD INDEX `idx_tenant` (`tenant_id`),
  ADD INDEX `idx_deleted` (`deleted`);

-- cms_admin_log
ALTER TABLE `cms_admin_log`
  ADD COLUMN `tenant_id` bigint NOT NULL DEFAULT 0 COMMENT '租户ID(多租户隔离)',
  ADD COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)',
  ADD INDEX `idx_tenant` (`tenant_id`),
  ADD INDEX `idx_deleted` (`deleted`);

-- cms_admin_nav
ALTER TABLE `cms_admin_nav`
  ADD COLUMN `tenant_id` bigint NOT NULL DEFAULT 0 COMMENT '租户ID(多租户隔离)',
  ADD COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)',
  ADD INDEX `idx_tenant` (`tenant_id`),
  ADD INDEX `idx_deleted` (`deleted`);

-- cms_admin_notice
ALTER TABLE `cms_admin_notice`
  ADD COLUMN `tenant_id` bigint NOT NULL DEFAULT 0 COMMENT '租户ID(多租户隔离)',
  ADD COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)',
  ADD INDEX `idx_tenant` (`tenant_id`),
  ADD INDEX `idx_deleted` (`deleted`);

-- cms_admin_role
ALTER TABLE `cms_admin_role`
  ADD COLUMN `tenant_id` bigint NOT NULL DEFAULT 0 COMMENT '租户ID(多租户隔离)',
  ADD COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)',
  ADD INDEX `idx_tenant` (`tenant_id`),
  ADD INDEX `idx_deleted` (`deleted`);

-- cms_admin_role_site
ALTER TABLE `cms_admin_role_site`
  ADD COLUMN `tenant_id` bigint NOT NULL DEFAULT 0 COMMENT '租户ID(多租户隔离)',
  ADD COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)',
  ADD INDEX `idx_tenant` (`tenant_id`),
  ADD INDEX `idx_deleted` (`deleted`);

-- cms_admin_role_value
ALTER TABLE `cms_admin_role_value`
  ADD COLUMN `tenant_id` bigint NOT NULL DEFAULT 0 COMMENT '租户ID(多租户隔离)',
  ADD COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)',
  ADD INDEX `idx_tenant` (`tenant_id`),
  ADD INDEX `idx_deleted` (`deleted`);

-- cms_ads
ALTER TABLE `cms_ads`
  ADD COLUMN `tenant_id` bigint NOT NULL DEFAULT 0 COMMENT '租户ID(多租户隔离)',
  ADD COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)',
  ADD INDEX `idx_tenant` (`tenant_id`),
  ADD INDEX `idx_deleted` (`deleted`);

-- cms_ads_category
ALTER TABLE `cms_ads_category`
  ADD COLUMN `tenant_id` bigint NOT NULL DEFAULT 0 COMMENT '租户ID(多租户隔离)',
  ADD COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)',
  ADD INDEX `idx_tenant` (`tenant_id`),
  ADD INDEX `idx_deleted` (`deleted`);

-- cms_ads_category_relation
ALTER TABLE `cms_ads_category_relation`
  ADD COLUMN `tenant_id` bigint NOT NULL DEFAULT 0 COMMENT '租户ID(多租户隔离)',
  ADD COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)',
  ADD INDEX `idx_tenant` (`tenant_id`),
  ADD INDEX `idx_deleted` (`deleted`);

-- cms_album
ALTER TABLE `cms_album`
  ADD COLUMN `tenant_id` bigint NOT NULL DEFAULT 0 COMMENT '租户ID(多租户隔离)',
  ADD COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)',
  ADD INDEX `idx_tenant` (`tenant_id`),
  ADD INDEX `idx_deleted` (`deleted`);

-- cms_article
ALTER TABLE `cms_article`
  ADD COLUMN `tenant_id` bigint NOT NULL DEFAULT 0 COMMENT '租户ID(多租户隔离)',
  ADD COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)',
  ADD INDEX `idx_tenant` (`tenant_id`),
  ADD INDEX `idx_deleted` (`deleted`);

-- cms_article_category
ALTER TABLE `cms_article_category`
  ADD COLUMN `tenant_id` bigint NOT NULL DEFAULT 0 COMMENT '租户ID(多租户隔离)',
  ADD COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)',
  ADD INDEX `idx_tenant` (`tenant_id`),
  ADD INDEX `idx_deleted` (`deleted`);

-- cms_article_category_relation
ALTER TABLE `cms_article_category_relation`
  ADD COLUMN `tenant_id` bigint NOT NULL DEFAULT 0 COMMENT '租户ID(多租户隔离)',
  ADD COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)',
  ADD INDEX `idx_tenant` (`tenant_id`),
  ADD INDEX `idx_deleted` (`deleted`);

-- cms_article_comment
ALTER TABLE `cms_article_comment`
  ADD COLUMN `tenant_id` bigint NOT NULL DEFAULT 0 COMMENT '租户ID(多租户隔离)',
  ADD COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)',
  ADD INDEX `idx_tenant` (`tenant_id`),
  ADD INDEX `idx_deleted` (`deleted`);

-- cms_article_label
ALTER TABLE `cms_article_label`
  ADD COLUMN `tenant_id` bigint NOT NULL DEFAULT 0 COMMENT '租户ID(多租户隔离)',
  ADD COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)',
  ADD INDEX `idx_tenant` (`tenant_id`),
  ADD INDEX `idx_deleted` (`deleted`);

-- cms_article_label_relation
ALTER TABLE `cms_article_label_relation`
  ADD COLUMN `tenant_id` bigint NOT NULL DEFAULT 0 COMMENT '租户ID(多租户隔离)',
  ADD COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)',
  ADD INDEX `idx_tenant` (`tenant_id`),
  ADD INDEX `idx_deleted` (`deleted`);

-- cms_article_property
ALTER TABLE `cms_article_property`
  ADD COLUMN `tenant_id` bigint NOT NULL DEFAULT 0 COMMENT '租户ID(多租户隔离)',
  ADD COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)',
  ADD INDEX `idx_tenant` (`tenant_id`),
  ADD INDEX `idx_deleted` (`deleted`);

-- cms_attach
ALTER TABLE `cms_attach`
  ADD COLUMN `tenant_id` bigint NOT NULL DEFAULT 0 COMMENT '租户ID(多租户隔离)',
  ADD COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)',
  ADD INDEX `idx_tenant` (`tenant_id`),
  ADD INDEX `idx_deleted` (`deleted`);

-- cms_link
ALTER TABLE `cms_link`
  ADD COLUMN `tenant_id` bigint NOT NULL DEFAULT 0 COMMENT '租户ID(多租户隔离)',
  ADD COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)',
  ADD INDEX `idx_tenant` (`tenant_id`),
  ADD INDEX `idx_deleted` (`deleted`);

-- cms_link_category
ALTER TABLE `cms_link_category`
  ADD COLUMN `tenant_id` bigint NOT NULL DEFAULT 0 COMMENT '租户ID(多租户隔离)',
  ADD COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)',
  ADD INDEX `idx_tenant` (`tenant_id`),
  ADD INDEX `idx_deleted` (`deleted`);

-- cms_link_category_relation
ALTER TABLE `cms_link_category_relation`
  ADD COLUMN `tenant_id` bigint NOT NULL DEFAULT 0 COMMENT '租户ID(多租户隔离)',
  ADD COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)',
  ADD INDEX `idx_tenant` (`tenant_id`),
  ADD INDEX `idx_deleted` (`deleted`);

-- cms_site
ALTER TABLE `cms_site`
  ADD COLUMN `tenant_id` bigint NOT NULL DEFAULT 0 COMMENT '租户ID(多租户隔离)',
  ADD COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)',
  ADD INDEX `idx_tenant` (`tenant_id`),
  ADD INDEX `idx_deleted` (`deleted`);

-- cms_site_channel
ALTER TABLE `cms_site_channel`
  ADD COLUMN `tenant_id` bigint NOT NULL DEFAULT 0 COMMENT '租户ID(多租户隔离)',
  ADD COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)',
  ADD INDEX `idx_tenant` (`tenant_id`),
  ADD INDEX `idx_deleted` (`deleted`);

-- cms_site_channel_album
ALTER TABLE `cms_site_channel_album`
  ADD COLUMN `tenant_id` bigint NOT NULL DEFAULT 0 COMMENT '租户ID(多租户隔离)',
  ADD COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)',
  ADD INDEX `idx_tenant` (`tenant_id`),
  ADD INDEX `idx_deleted` (`deleted`);

-- cms_site_channel_field
ALTER TABLE `cms_site_channel_field`
  ADD COLUMN `tenant_id` bigint NOT NULL DEFAULT 0 COMMENT '租户ID(多租户隔离)',
  ADD COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)',
  ADD INDEX `idx_tenant` (`tenant_id`),
  ADD INDEX `idx_deleted` (`deleted`);

-- cms_site_domain
ALTER TABLE `cms_site_domain`
  ADD COLUMN `tenant_id` bigint NOT NULL DEFAULT 0 COMMENT '租户ID(多租户隔离)',
  ADD COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)',
  ADD INDEX `idx_tenant` (`tenant_id`),
  ADD INDEX `idx_deleted` (`deleted`);

-- cms_tag
ALTER TABLE `cms_tag`
  ADD COLUMN `tenant_id` bigint NOT NULL DEFAULT 0 COMMENT '租户ID(多租户隔离)',
  ADD COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)',
  ADD INDEX `idx_tenant` (`tenant_id`),
  ADD INDEX `idx_deleted` (`deleted`);

-- cms_theme
ALTER TABLE `cms_theme`
  ADD COLUMN `tenant_id` bigint NOT NULL DEFAULT 0 COMMENT '租户ID(多租户隔离)',
  ADD COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)',
  ADD INDEX `idx_tenant` (`tenant_id`),
  ADD INDEX `idx_deleted` (`deleted`);

-- cms_topic
ALTER TABLE `cms_topic`
  ADD COLUMN `tenant_id` bigint NOT NULL DEFAULT 0 COMMENT '租户ID(多租户隔离)',
  ADD COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)',
  ADD INDEX `idx_tenant` (`tenant_id`),
  ADD INDEX `idx_deleted` (`deleted`);

-- plg_online_register
ALTER TABLE `plg_online_register`
  ADD COLUMN `tenant_id` bigint NOT NULL DEFAULT 0 COMMENT '租户ID(多租户隔离)',
  ADD COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)',
  ADD INDEX `idx_tenant` (`tenant_id`),
  ADD INDEX `idx_deleted` (`deleted`);

-- weixin_account
ALTER TABLE `weixin_account`
  ADD COLUMN `tenant_id` bigint NOT NULL DEFAULT 0 COMMENT '租户ID(多租户隔离)',
  ADD COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)',
  ADD INDEX `idx_tenant` (`tenant_id`),
  ADD INDEX `idx_deleted` (`deleted`);

-- weixin_menu
ALTER TABLE `weixin_menu`
  ADD COLUMN `tenant_id` bigint NOT NULL DEFAULT 0 COMMENT '租户ID(多租户隔离)',
  ADD COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)',
  ADD INDEX `idx_tenant` (`tenant_id`),
  ADD INDEX `idx_deleted` (`deleted`);

-- weixin_mp_verify
ALTER TABLE `weixin_mp_verify`
  ADD COLUMN `tenant_id` bigint NOT NULL DEFAULT 0 COMMENT '租户ID(多租户隔离)',
  ADD COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)',
  ADD INDEX `idx_tenant` (`tenant_id`),
  ADD INDEX `idx_deleted` (`deleted`);

-- weixin_request_content
ALTER TABLE `weixin_request_content`
  ADD COLUMN `tenant_id` bigint NOT NULL DEFAULT 0 COMMENT '租户ID(多租户隔离)',
  ADD COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)',
  ADD INDEX `idx_tenant` (`tenant_id`),
  ADD INDEX `idx_deleted` (`deleted`);

-- weixin_request_rule
ALTER TABLE `weixin_request_rule`
  ADD COLUMN `tenant_id` bigint NOT NULL DEFAULT 0 COMMENT '租户ID(多租户隔离)',
  ADD COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)',
  ADD INDEX `idx_tenant` (`tenant_id`),
  ADD INDEX `idx_deleted` (`deleted`);

-- weixin_response_content
ALTER TABLE `weixin_response_content`
  ADD COLUMN `tenant_id` bigint NOT NULL DEFAULT 0 COMMENT '租户ID(多租户隔离)',
  ADD COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)',
  ADD INDEX `idx_tenant` (`tenant_id`),
  ADD INDEX `idx_deleted` (`deleted`);

-- ========== 4. 初始化一个默认租户(可选) ==========
INSERT INTO `cms_tenant` (`code`, `name`, `short_name`, `status`, `create_name`)
VALUES ('default', '默认租户', '默认', 1, 'system')
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`);

SET FOREIGN_KEY_CHECKS = 1;

-- 验证: 确认所有表都有 tenant_id 和 deleted 列
SELECT TABLE_NAME, COLUMN_NAME
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA = DATABASE()
  AND COLUMN_NAME IN ('tenant_id', 'deleted')
  AND (TABLE_NAME LIKE 'cms_%' OR TABLE_NAME LIKE 'plg_%' OR TABLE_NAME LIKE 'weixin_%')
ORDER BY TABLE_NAME, COLUMN_NAME;


-- ============================================================
-- 修复: deleted 列类型 bit(1) -> tinyint(1)
-- 问题: MySQL bit(1) 驱动返回 []byte,无法 Scan 到 Go bool 字段
--       导致 'sql: Scan error on column ... name "deleted": couldn't convert "\x00" into type bool'
-- 影响: 所有 cms_* / plg_* / weixin_* 表(39 张)
-- 时间: 2026/09/09
-- ============================================================

-- 步骤1: ALTER 所有表, 将 bit(1) 改为 tinyint(1)
-- tinyint(1) 是 GORM 推荐的 bool 映射类型 (0=false, 1=true)
-- 转换无数据丢失: bit(1) 内部就是 0/1, tinyint(1) 也是 0/1

-- 1) cms_*
ALTER TABLE `cms_ad_category_relation`    MODIFY COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)';
ALTER TABLE `cms_admin`                   MODIFY COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)';
ALTER TABLE `cms_admin_log`               MODIFY COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)';
ALTER TABLE `cms_admin_nav`               MODIFY COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)';
ALTER TABLE `cms_admin_notice`            MODIFY COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)';
ALTER TABLE `cms_admin_role`              MODIFY COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)';
ALTER TABLE `cms_admin_role_site`         MODIFY COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)';
ALTER TABLE `cms_admin_role_value`        MODIFY COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)';
ALTER TABLE `cms_ads`                     MODIFY COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)';
ALTER TABLE `cms_ads_category`            MODIFY COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)';
ALTER TABLE `cms_ads_category_relation`   MODIFY COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)';
ALTER TABLE `cms_album`                   MODIFY COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)';
ALTER TABLE `cms_article`                 MODIFY COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)';
ALTER TABLE `cms_article_category`        MODIFY COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)';
ALTER TABLE `cms_article_category_relation` MODIFY COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)';
ALTER TABLE `cms_article_comment`         MODIFY COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)';
ALTER TABLE `cms_article_label`           MODIFY COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)';
ALTER TABLE `cms_article_label_relation`  MODIFY COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)';
ALTER TABLE `cms_article_property`        MODIFY COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)';
ALTER TABLE `cms_attach`                  MODIFY COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)';
ALTER TABLE `cms_link`                    MODIFY COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)';
ALTER TABLE `cms_link_category`           MODIFY COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)';
ALTER TABLE `cms_link_category_relation`  MODIFY COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)';
ALTER TABLE `cms_site`                    MODIFY COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)';
ALTER TABLE `cms_site_channel`            MODIFY COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)';
ALTER TABLE `cms_site_channel_album`      MODIFY COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)';
ALTER TABLE `cms_site_channel_field`      MODIFY COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)';
ALTER TABLE `cms_site_domain`             MODIFY COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)';
ALTER TABLE `cms_tag`                     MODIFY COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)';
ALTER TABLE `cms_tenant`                  MODIFY COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)';
ALTER TABLE `cms_theme`                   MODIFY COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)';
ALTER TABLE `cms_topic`                   MODIFY COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)';

-- 2) plg_*
ALTER TABLE `plg_online_register`            MODIFY COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)';

-- 3) weixin_*
ALTER TABLE `weixin_account`                  MODIFY COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)';
ALTER TABLE `weixin_menu`                     MODIFY COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)';
ALTER TABLE `weixin_mp_verify`                MODIFY COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)';
ALTER TABLE `weixin_request_content`          MODIFY COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)';
ALTER TABLE `weixin_request_rule`             MODIFY COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)';
ALTER TABLE `weixin_response_content`         MODIFY COLUMN `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '软删除标记(0未删除1已删除)';

-- 验证: 确认所有 deleted 列已转为 tinyint(1)
SELECT TABLE_NAME, COLUMN_NAME, COLUMN_TYPE
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA = DATABASE()
  AND COLUMN_NAME = 'deleted'
  AND (TABLE_NAME LIKE 'cms_%' OR TABLE_NAME LIKE 'plg_%' OR TABLE_NAME LIKE 'weixin_%')
ORDER BY TABLE_NAME;
