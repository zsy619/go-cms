-- 语言设置表初始化脚本
-- 执行前请先在 MySQL 中创建数据库和表

-- 创建表
CREATE TABLE IF NOT EXISTS cms_language (
  language_id bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  name varchar(64) NOT NULL COMMENT '语言名称',
  code varchar(16) NOT NULL COMMENT '语言代码',
  icon varchar(64) DEFAULT NULL COMMENT '图标/国旗',
  sort_id int DEFAULT 99 COMMENT '排序',
  description varchar(256) DEFAULT NULL COMMENT '描述',
  status tinyint DEFAULT 1 COMMENT '状态: 0禁用 1启用',
  is_default tinyint(1) DEFAULT FALSE COMMENT '是否默认',
  create_id int DEFAULT NULL COMMENT '创建人ID',
  create_name varchar(64) DEFAULT NULL COMMENT '创建人姓名',
  create_time datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  update_id int DEFAULT NULL COMMENT '更新人ID',
  update_name varchar(64) DEFAULT NULL COMMENT '更新人姓名',
  update_time datetime DEFAULT NULL COMMENT '修改时间',
  tenant_id bigint NOT NULL DEFAULT 0 COMMENT '租户ID',
  deleted tinyint(1) NOT NULL DEFAULT 0 COMMENT '逻辑删除',
  PRIMARY KEY (language_id),
  KEY idx_tenant_id (tenant_id),
  KEY idx_deleted (deleted)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='语言设置表';

-- 初始化语言数据
INSERT INTO cms_language (name, code, icon, sort_id, description, status, is_default) VALUES
('简体中文', 'zh-CN', '🇨🇳', 1, '简体中文（Simplified Chinese）', 1, TRUE),
('English', 'en', '🇺🇸', 2, '英语（English）', 1, FALSE),
('繁体中文', 'zh-TW', '🇨🇳', 3, '繁体中文（Traditional Chinese）', 1, FALSE),
('日本語', 'ja', '🇯🇵', 4, '日语（Japanese）', 1, FALSE),
('한국어', 'ko', '🇰🇷', 5, '韩语（Korean）', 1, FALSE),
('Français', 'fr', '🇫🇷', 6, '法语（French）', 1, FALSE),
('Deutsch', 'de', '🇩🇪', 7, '德语（German）', 1, FALSE),
('Español', 'es', '🇪🇸', 8, '西班牙语（Spanish）', 1, FALSE),
('Português', 'pt', '🇵🇹', 9, '葡萄牙语（Portuguese）', 1, FALSE),
('Русский', 'ru', '🇷🇺', 10, '俄语（Russian）', 1, FALSE);
