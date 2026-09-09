# 多租户 & 软删除 数据库迁移脚本

本目录包含**生产环境**使用的 SQL 迁移脚本，覆盖 4 类数据库方言。

> **设计参考**：芋道 yudao 框架的 `BaseDO` / `TenantBaseDO`。
> - 软删除使用单一字段 `deleted`（与 MyBatis-Plus `@TableLogic` 一致）
> - 不再使用 `is_deleted`，避免语义冲突

## 文件清单

| 文件 | 适用数据库 | 备注 |
|------|-----------|------|
| `01_mysql.sql` | MySQL / MariaDB | `tinyint(1)` 用 `0` 表示 false |
| `02_postgres.sql` | PostgreSQL / openGauss / KingbaseES (人大金仓) | `boolean` 类型 |
| `03_sqlserver.sql` | Microsoft SQL Server | `BIT` 类型，使用 `IF NOT EXISTS` 包裹 |
| `04_weixin_init.sql` | MySQL / MariaDB | **微信模块 6 张表初始建表**（`weixin_account/menu/mp_verify/request_content/request_rule/response_content`），与 domain 严格对齐 |

> 达梦 DM8 / Oracle 请参考 `02_postgres.sql` 自行调整（`boolean`→`NUMBER(1)`）。

## 字段说明（与芋道 yudao 一致）

每张业务表添加 2 列（按表结构中顺序）：

| 列名 | MySQL | PostgreSQL | SQL Server | 用途 |
|------|-------|------------|------------|------|
| `tenant_id` | `bigint NOT NULL DEFAULT 0` | `bigint NOT NULL DEFAULT 0` | `BIGINT NOT NULL DEFAULT 0` | **多租户隔离**（与 `cms_tenant.tenant_id` 关联） |
| `deleted`    | `tinyint(1) NOT NULL DEFAULT 0` | `boolean DEFAULT false` | `BIT NOT NULL DEFAULT 0` | **软删除标记**（单一字段） |

### 字段设计说明

| 字段 | 为什么这样设计 |
|------|---------------|
| `tenant_id` 放最后两列之前 | 业务字段在前，框架字段在后，便于阅读 |
| `deleted` 单一字段 | 与 MyBatis-Plus / GORM / 芋道保持一致，语义明确 |
| `deleted` MySQL 用 `tinyint(1)` 而非 `bit(1)` | `bit(1)` 在 Go MySQL 驱动中返回 `[]byte`，无法直接 Scan 到 `bool`；`tinyint(1)` 是 GORM 推荐的 bool 映射类型 |
| **不再使用 `is_deleted`** | 与芋道 yudao 对齐，避免双重删除标记造成的混乱 |

## 使用方法

### MySQL

```bash
mysqldump -u root -p cms > cms_backup.sql   # 1. 备份
mysql -u root -p cms < 01_mysql.sql          # 2. 迁移
```

### PostgreSQL

```bash
pg_dump -U postgres cms > cms_backup.sql     # 1. 备份
psql -U postgres -d cms -f 02_postgres.sql   # 2. 迁移
```

### SQL Server

```sql
BACKUP DATABASE [cms] TO DISK = 'D:\cms_backup.bak'  -- 1. 备份
-- 2. 在 SSMS 中打开 03_sqlserver.sql 并执行
```

## 脚本内容概要

1. **第 1 步**：创建 `cms_tenant` 租户主表
2. **第 2 步**：给 `cms_admin` 添加 `tenant_id` 和 `deleted`（同时迁移旧的 `is_deleted` 数据到 `deleted`，然后 DROP `is_deleted`）
3. **第 3 步**：给其他 37 张业务表批量添加 `tenant_id`、`deleted` 两列 + 索引
4. **第 4 步**：插入一条默认租户 `code='default'` 记录（保证系统启动后查询不报错）
5. **第 5 步**：验证 SQL（查询 `information_schema` 确认所有列已添加）

## 注意事项

1. **生产环境执行前务必先备份**
2. **大表 ALTER TABLE 可能锁表** —— 建议在低峰期执行，或使用 `pt-online-schema-change` / `gh-ost` 等工具
3. **SQL Server** 必须用 `GO` 分批执行，**不能**整段粘贴到 SSMS 查询窗口
4. **从 is_deleted 到 deleted 的数据迁移**：脚本已包含 `UPDATE cms_admin SET deleted = 1 WHERE is_deleted = 1;` 并 DROP COLUMN
5. **租户数据迁移**：现有数据行 `tenant_id=0`，需要根据业务手动分配：
   ```sql
   -- 把所有现有数据划归默认租户(tenant_id=1)
   UPDATE cms_article SET tenant_id = 1 WHERE tenant_id = 0;
   ```
6. **CMS 启动时自动迁移**（`cms.migrate=true`）也会创建新列，但如果生产库有大量数据，建议手动 SQL 迁移以避免 GORM AutoMigrate 锁表
