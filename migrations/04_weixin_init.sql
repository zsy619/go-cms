-- ============================================================
-- weixin_* 微信模块 - 初始建表脚本
-- 数据库: MySQL / MariaDB
-- 适用范围: 微信公众平台管理相关 6 张表
-- 表列表:
--   weixin_account           - 公众号配置表
--   weixin_menu              - 自定义菜单表
--   weixin_mp_verify         - 公众号验证文件表
--   weixin_request_content   - 请求回复内容表(规则关联)
--   weixin_request_rule      - 请求规则表(关键字触发)
--   weixin_response_content  - 响应内容日志表(消息记录)
-- 设计: 与 cms_* 系列表保持一致风格(双字段 tenant_id + deleted)
-- ============================================================


-- WeixinAccount mapped from table <weixin_account>
CREATE TABLE IF NOT EXISTS `weixin_account` (
  `account_id` bigint NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT '主键',
  `name` varchar(64) DEFAULT NULL COMMENT '公众号名称',
  `original_id` varchar(64) DEFAULT NULL COMMENT '公众号原始ID',
  `wx_code` varchar(64) DEFAULT NULL COMMENT '公众平台微信号',
  `token` varchar(512) DEFAULT NULL COMMENT '令牌ToKen',
  `app_id` varchar(128) DEFAULT NULL COMMENT '开发者IDAppId',
  `app_secret` varchar(128) DEFAULT NULL COMMENT '开发者密码AppSecret',
  `app_aes_key` varchar(128) DEFAULT NULL COMMENT '消息加解密密钥',
  `is_push` tinyint(1) DEFAULT NULL COMMENT '内容推送',
  `sort_id` int DEFAULT NULL COMMENT '排序',
  `status` tinyint DEFAULT NULL COMMENT '状态0草稿1提交2审核通过3审核未通过4驳回',
  `belong_to` varchar(64) DEFAULT NULL COMMENT '归属',
  `create_id` int DEFAULT NULL COMMENT '创建人ID',
  `create_name` varchar(64) DEFAULT NULL COMMENT '创建人姓名',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_id` int DEFAULT NULL COMMENT '更新人ID',
  `update_name` varchar(64) DEFAULT NULL COMMENT '更新人姓名',
  `update_time` datetime DEFAULT NULL COMMENT '修改时间',
  `tenant_id` bigint NOT NULL DEFAULT 0 COMMENT '租户ID(多租户隔离)',
  `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '逻辑删除标识(0未删除1已删除)'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='WeixinAccount';


-- WeixinMenu mapped from table <weixin_menu>
CREATE TABLE IF NOT EXISTS `weixin_menu` (
  `menu_id` bigint NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT '主键',
  `parent_id` bigint NOT NULL COMMENT '父级ID',
  `account_id` bigint NOT NULL COMMENT '归属公众号',
  `name` varchar(40) DEFAULT NULL COMMENT '菜单标题，不超过16个字节，子菜单不超过40个字节',
  `type` varchar(32) DEFAULT NULL COMMENT '菜单的响应动作类型，view表示网页类型，click表示点击类型，miniprogram表示小程序类型',
  `key` varchar(128) DEFAULT NULL COMMENT '菜单KEY值，用于消息接口推送，不超过128字节',
  `url` varchar(1024) DEFAULT NULL COMMENT '网页链接，用户点击菜单可打开链接，不超过1024字节。当type为miniprogram时，不支持小程序的老版本客户端将打开本url',
  `app_id` varchar(128) DEFAULT NULL COMMENT '小程序appid',
  `page_path` varchar(128) DEFAULT NULL COMMENT '小程序页面路径',
  `media_id` varchar(256) DEFAULT NULL COMMENT 'media_id类型和view_limited类型必须',
  `article_id` varchar(256) DEFAULT NULL COMMENT 'article_id类型和article_view_limited类型必须',
  `sort_id` int DEFAULT NULL COMMENT '排序',
  `create_id` int DEFAULT NULL COMMENT '创建人ID',
  `create_name` varchar(64) DEFAULT NULL COMMENT '创建人姓名',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_id` int DEFAULT NULL COMMENT '更新人ID',
  `update_name` varchar(64) DEFAULT NULL COMMENT '更新人姓名',
  `update_time` datetime DEFAULT NULL COMMENT '修改时间',
  `tenant_id` bigint NOT NULL DEFAULT 0 COMMENT '租户ID(多租户隔离)',
  `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '逻辑删除标识(0未删除1已删除)'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='WeixinMenu';


-- WeixinMpVerify mapped from table <weixin_mp_verify>
CREATE TABLE IF NOT EXISTS `weixin_mp_verify` (
  `verify_id` bigint NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT '主键',
  `account_id` bigint NOT NULL COMMENT '归属公众号',
  `path` varchar(128) DEFAULT NULL COMMENT '访问路径，默认/MP_verify_公众号原始ID.txt',
  `file_path` varchar(128) DEFAULT NULL COMMENT '文件路径',
  `file_name` varchar(128) DEFAULT NULL COMMENT '文件名称',
  `sort_id` int DEFAULT NULL COMMENT '排序',
  `status` tinyint DEFAULT NULL COMMENT '状态0草稿1提交2审核通过3审核未通过4驳回',
  `create_id` int DEFAULT NULL COMMENT '创建人ID',
  `create_name` varchar(64) DEFAULT NULL COMMENT '创建人姓名',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_id` int DEFAULT NULL COMMENT '更新人ID',
  `update_name` varchar(64) DEFAULT NULL COMMENT '更新人姓名',
  `update_time` datetime DEFAULT NULL COMMENT '修改时间',
  `tenant_id` bigint NOT NULL DEFAULT 0 COMMENT '租户ID(多租户隔离)',
  `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '逻辑删除标识(0未删除1已删除)'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='WeixinMpVerify';


-- WeixinRequestContent mapped from table <weixin_request_content>
CREATE TABLE IF NOT EXISTS `weixin_request_content` (
  `content_id` bigint NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT '主键',
  `rule_id` bigint NOT NULL COMMENT '规则ID',
  `title` varchar(512) DEFAULT NULL COMMENT '回复标题',
  `content` text DEFAULT NULL COMMENT '回复内容',
  `link_url` varchar(512) DEFAULT NULL COMMENT '详情链接地址',
  `img_url` varchar(512) DEFAULT NULL COMMENT '图片地址',
  `media_url` varchar(512) DEFAULT NULL COMMENT '语音或视频地址',
  `media_hd_url` varchar(512) DEFAULT NULL COMMENT '高清语音或者视频地址',
  `sort_id` int DEFAULT NULL COMMENT '排序',
  `create_id` int DEFAULT NULL COMMENT '创建人ID',
  `create_name` varchar(64) DEFAULT NULL COMMENT '创建人姓名',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_id` int DEFAULT NULL COMMENT '更新人ID',
  `update_name` varchar(64) DEFAULT NULL COMMENT '更新人姓名',
  `update_time` datetime DEFAULT NULL COMMENT '修改时间',
  `tenant_id` bigint NOT NULL DEFAULT 0 COMMENT '租户ID(多租户隔离)',
  `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '逻辑删除标识(0未删除1已删除)'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='WeixinRequestContent';


-- WeixinRequestRule mapped from table <weixin_request_rule>
CREATE TABLE IF NOT EXISTS `weixin_request_rule` (
  `rule_id` bigint NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT '主键',
  `account_id` bigint NOT NULL COMMENT '归属公众号',
  `name` varchar(256) DEFAULT NULL COMMENT '规则名称',
  `keywords` varchar(2048) DEFAULT NULL COMMENT '请求关键词,逗号分隔',
  `request_type` int DEFAULT NULL COMMENT '请求类型(0默认回复1文字2图片3语音4链接5地理位置6关注7取消关注8扫描带参数二维码事件9上报地理位置事件10自定义菜单事件）',
  `response_type` int DEFAULT NULL COMMENT '回复类型(1文本2图文3语音4视频5第三方接口)',
  `is_like_query` tinyint DEFAULT NULL COMMENT '是否模糊查询',
  `is_default` tinyint DEFAULT NULL COMMENT '是否默认回复',
  `sort_id` int DEFAULT NULL COMMENT '排序',
  `create_id` int DEFAULT NULL COMMENT '创建人ID',
  `create_name` varchar(64) DEFAULT NULL COMMENT '创建人姓名',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_id` int DEFAULT NULL COMMENT '更新人ID',
  `update_name` varchar(64) DEFAULT NULL COMMENT '更新人姓名',
  `update_time` datetime DEFAULT NULL COMMENT '修改时间',
  `tenant_id` bigint NOT NULL DEFAULT 0 COMMENT '租户ID(多租户隔离)',
  `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '逻辑删除标识(0未删除1已删除)'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='WeixinRequestRule';


-- WeixinResponseContent mapped from table <weixin_response_content>
CREATE TABLE IF NOT EXISTS `weixin_response_content` (
  `content_id` bigint NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT '主键',
  `account_id` bigint NOT NULL COMMENT '归属公众号',
  `open_id` varchar(128) DEFAULT NULL COMMENT '用户微信ID',
  `request_type` varchar(32) DEFAULT NULL COMMENT '数据类型 文本消息：text 图片消息:image 地理位置消息:location 链接消息:link 事件:event',
  `request_content` varchar(2048) DEFAULT NULL COMMENT '数据内容',
  `response_type` varchar(32) DEFAULT NULL COMMENT '回复的类型 文本消息：text 图片消息:image 地理位置消息:location 链接消息:link',
  `reponse_content` varchar(2048) DEFAULT NULL COMMENT '系统回复的内容',
  `meida_hd_url` varchar(512) DEFAULT NULL COMMENT '高清语音或者视频地址',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `xml_content` varchar(2048) DEFAULT NULL COMMENT 'xml原始内容',
  `add_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '录入系统的时间',
  `tenant_id` bigint NOT NULL DEFAULT 0 COMMENT '租户ID(多租户隔离)',
  `deleted` tinyint(1) NOT NULL DEFAULT 0 COMMENT '逻辑删除标识(0未删除1已删除)'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='WeixinResponseContent';

-- ============================================================
-- 多租户 + 软删除索引 (应用于所有 weixin_* 表)
-- ============================================================

-- weixin_account
CREATE INDEX idx_weixin_account_tenant ON weixin_account(tenant_id);
CREATE INDEX idx_weixin_account_deleted ON weixin_account(deleted);
CREATE INDEX idx_weixin_account_account ON weixin_account(account_id);

-- weixin_menu
CREATE INDEX idx_weixin_menu_tenant ON weixin_menu(tenant_id);
CREATE INDEX idx_weixin_menu_deleted ON weixin_menu(deleted);
CREATE INDEX idx_weixin_menu_account ON weixin_menu(account_id);

-- weixin_mp_verify
CREATE INDEX idx_weixin_mp_verify_tenant ON weixin_mp_verify(tenant_id);
CREATE INDEX idx_weixin_mp_verify_deleted ON weixin_mp_verify(deleted);
CREATE INDEX idx_weixin_mp_verify_account ON weixin_mp_verify(account_id);

-- weixin_request_content
CREATE INDEX idx_weixin_request_content_tenant ON weixin_request_content(tenant_id);
CREATE INDEX idx_weixin_request_content_deleted ON weixin_request_content(deleted);
CREATE INDEX idx_weixin_request_content_account ON weixin_request_content(account_id);

-- weixin_request_rule
CREATE INDEX idx_weixin_request_rule_tenant ON weixin_request_rule(tenant_id);
CREATE INDEX idx_weixin_request_rule_deleted ON weixin_request_rule(deleted);
CREATE INDEX idx_weixin_request_rule_account ON weixin_request_rule(account_id);

-- weixin_response_content
CREATE INDEX idx_weixin_response_content_tenant ON weixin_response_content(tenant_id);
CREATE INDEX idx_weixin_response_content_deleted ON weixin_response_content(deleted);
CREATE INDEX idx_weixin_response_content_account ON weixin_response_content(account_id);

