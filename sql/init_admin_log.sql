-- 操作日志表 - cms_admin_log
-- 用于记录后台管理操作日志

CREATE TABLE IF NOT EXISTS cms_admin_log (
    log_id        BIGINT       PRIMARY KEY AUTO_INCREMENT COMMENT '主键',
    user_id       BIGINT       DEFAULT 0 COMMENT '用户ID',
    user_name     VARCHAR(128) DEFAULT '' COMMENT '账号',
    method        VARCHAR(32)  DEFAULT '' COMMENT '请求方法',
    path          VARCHAR(128) DEFAULT '' COMMENT '请求路径',
    query         VARCHAR(256) DEFAULT '' COMMENT '请求参数(JSON)',
    status_code   VARCHAR(64)  DEFAULT '' COMMENT '响应状态码',
    ip            VARCHAR(64)  DEFAULT '' COMMENT 'IP地址',
    create_time   DATETIME     DEFAULT CURRENT_TIMESTAMP COMMENT '记录时间',
    tenant_id     BIGINT       NOT NULL DEFAULT 0 COMMENT '租户ID',
    deleted       TINYINT(1)   NOT NULL DEFAULT 0 COMMENT '逻辑删除标识',
    INDEX idx_user_id (user_id),
    INDEX idx_create_time (create_time),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_deleted (deleted),
    INDEX idx_path (path),
    INDEX idx_method (method)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='操作日志表';
