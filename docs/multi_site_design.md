# CMS 多站点（Multi-Site）支持实施方案

> 基于现有 `cms_site` / `cms_site_domain` / `cms_site_channel` 基础设施，> 完善多租户隔离、域名解析、PC/移动端分流、主题继承、模板降级、灰度上线等关键能力，> 给出贴合本项目实际的端到端可执行方案。

---

## 目录

- [一、现状盘点](#一现状盘点)
- [二、目标与设计原则](#二目标与设计原则)
- [三、数据模型增强](#三数据模型增强)
- [四、迁移 SQL（多方言）](#四迁移-sql多方言)
- [五、域名解析器 SiteResolver](#五域名解析器-siteresolver)
- [六、Beego 中间件](#六beego-中间件)
- [七、SiteContext 上下文](#七sitecontext-上下文)
- [八、Theme 主题继承](#八theme-主题继承)
- [九、PC/移动端分流](#九pc移动端分流)
- [十、缓存策略](#十缓存策略)
- [十一、后台管理界面](#十一后台管理界面)
- [十二、运维与配置](#十二运维与配置)
- [十三、回滚与灰度策略](#十三回滚与灰度策略)
- [十四、风险与对策](#十四风险与对策)
- [十五、里程碑与工作量估算](#十五里程碑与工作量估算)
- [附录 A：典型场景示例](#附录-a典型场景示例)
- [附录 B：参考链接](#附录-b参考链接)

---

## 一、现状盘点

### 1.1 现有数据表（已具备基础）

| 表 | 字段数 | 说明 |
|------|------|------|
| `cms_site` | 30+ | 站点主表（名称、LOGO、模板、备案号、SEO 等完整配置）|
| `cms_site_domain` | 4 | 域名绑定（多域名可绑定同一站点）|
| `cms_site_channel` | 20+ | 站点频道（含 `site_id` 关联、模板路径）|
| `cms_site_channel_album` | - | 频道相册 |
| `cms_site_channel_field` | - | 频道扩展字段 |

### 1.2 现有代码能力

| 能力 | 位置 | 状态 |
|------|------|------|
| 根据 Host 找站点 | `app/cms/service/api_site.go::FindByHost` | ✅ 已实现 |
| 默认站点 fallback | `ApiSite::Default()` 查 `is_default=1` | ✅ 已实现 |
| 前台 BaseController 注入 | `controllers/www/b.go::Prepare()` | ✅ 已实现 |
| 站点配置缓存 | `app/tool/cache.go::SiteCache` | ⚠️ 60 秒过期，无刷新机制 |
| 后台站点管理 CRUD | `controllers/admin/site.go` | ✅ 已实现 |
| 多租户隔离 | `cms_site.tenant_id` | ✅ 已支持 |
| 主题模板路径 | `SiteTheme` 全局变量 + `cms_theme.is_default` | ⚠️ 无继承链 |
| PC/移动端分离 | `cms_site.is_mobile` | ⚠️ 仅布尔位，无模板分流 |
| 子域名 / 泛解析 | `cms_site_domain.domain` 精确匹配 | ⚠️ 不支持 `*.edu.cn` |
| 站点切换（后台）| 无 | ❌ 缺失 |
| 跨站点数据隔离 | 依赖 `site_id` 列 | ⚠️ 部分表缺 `site_id` |

### 1.3 现存问题（待解决）

1. **域名匹配仅精确**：`cms_site_domain.domain = "www.example.com"`，无法匹配 `*.example.com` 子域
2. **无主题继承**：所有同级站点必须独立指定模板，无法复用主站模板
3. **PC/移动端硬编码**：`is_mobile` 字段仅作标识，没有按 UA 自动分流到不同模板
4. **缓存单一**：站点信息变更后，最长 60s 才生效，无主动失效机制
5. **多租户维度缺失**：管理员切换租户后，站点列表是否只显示当前租户？目前未明确
6. **站点切换 UX**：后台管理员需在多个站点间切换内容，当前无统一入口
7. **域名泛解析**：教育部平台常需 `*.edu.cn` → 统一子站
8. **HTTPS / SNI**：多站点共享 IP 时需要 SNI，多证书管理

---

## 二、目标与设计原则

### 2.1 目标

1. **一个进程支撑 N 个独立站点**：每个站点拥有独立的域名、LOGO、模板、内容、SEO
2. **多租户隔离**：不同租户的站点不可互见、不可越权访问
3. **泛域名/子域名**：支持 `*.example.com` 模式路由
4. **主题继承**：子站可继承父站模板，仅按需覆盖
5. **PC / 移动端自动分流**：根据 UA 自动匹配 PC 模板 / 移动模板
6. **零侵入上线**：复用现有 `cms_site` 等表结构，新增少量字段即可
7. **缓存友好**：站点信息全程内存缓存 + 主动失效
8. **跨方言兼容**：MySQL / PG / MSSQL / Kingbase / GaussDB 全支持

### 2.2 设计原则

| 原则 | 落地 |
|------|------|
| **域名驱动** | 一切从 `Host` 头开始，无 Host 时用默认站 |
| **租户 + 站点二维** | `tenant_id` 一级隔离，`site_id` 二级隔离 |
| **继承而非重复** | 主题 / 配置可继承父站，未配置时 fallback |
| **PC 优先，移动端按需** | 默认 PC，启用移动端后按 UA 分流 |
| **缓存先行** | 站点信息进入 L1/L2 双层缓存 |
| **优雅失败** | 未知域名 → default site + 日志告警 |

---

## 三、数据模型增强

### 3.1 cms_site 增强字段

现有 30+ 字段已覆盖站点主体配置，仅需小幅扩展：

```sql
-- cms_site 新增字段
ALTER TABLE cms_site
  ADD COLUMN parent_site_id  BIGINT        NOT NULL DEFAULT 0      COMMENT "父站ID(用于主题继承)",
  ADD COLUMN theme_inherit   TINYINT(1)    NOT NULL DEFAULT 1      COMMENT "是否继承父站主题(0否1是)",
  ADD COLUMN mobile_template VARCHAR(128)  NOT NULL DEFAULT ""     COMMENT "移动端模板(留空则使用相同模板)",
  ADD COLUMN desktop_template VARCHAR(128) NOT NULL DEFAULT ""    COMMENT "PC端模板(默认=template)",
  ADD COLUMN site_kind       TINYINT       NOT NULL DEFAULT 0      COMMENT "站点类型:0主站1子站2微站3门户4学校",
  ADD COLUMN site_code       VARCHAR(64)   NOT NULL DEFAULT ""     COMMENT "站点代码(英文唯一标识,用于URL)",
  ADD COLUMN site_status     TINYINT       NOT NULL DEFAULT 1      COMMENT "状态:0停用1启用2维护中",
  ADD COLUMN enable_wildcard TINYINT(1)    NOT NULL DEFAULT 0      COMMENT "启用泛域名:0否1是",
  ADD COLUMN wildcard_pattern VARCHAR(128) NOT NULL DEFAULT ""    COMMENT "泛域名规则,如 *.edu.cn",
  ADD COLUMN site_settings   TEXT          NULL                     COMMENT "JSON扩展配置(任意KV)";

-- 索引
ALTER TABLE cms_site ADD UNIQUE INDEX uk_site_code (tenant_id, site_code);
ALTER TABLE cms_site ADD INDEX idx_parent_site (parent_site_id);
ALTER TABLE cms_site ADD INDEX idx_site_status (site_status);
```

### 3.2 cms_site_domain 增强字段

现有字段：domain_id, site_id, domain, remark，增强为支持泛解析：

```sql
-- cms_site_domain 新增字段
ALTER TABLE cms_site_domain
  ADD COLUMN domain_kind  TINYINT     NOT NULL DEFAULT 0      COMMENT "域名类型:0精确1泛域2正则",
  ADD COLUMN priority    INT         NOT NULL DEFAULT 100    COMMENT "匹配优先级(数值越小越高)",
  ADD COLUMN https_mode  TINYINT(1)  NOT NULL DEFAULT 1      COMMENT "HTTPS:0关闭1强制2自动",
  ADD COLUMN cert_id     BIGINT      NOT NULL DEFAULT 0      COMMENT "关联SSL证书ID(cms_ssl_cert)",
  ADD COLUMN is_primary  TINYINT(1)  NOT NULL DEFAULT 0      COMMENT "是否主域名",
  ADD COLUMN redirect_to VARCHAR(256) NOT NULL DEFAULT ""    COMMENT "强制跳转域名(空=不跳转)";

-- 索引
ALTER TABLE cms_site_domain ADD INDEX idx_priority (priority);
```

### 3.3 新增 cms_ssl_cert 表（HTTPS 多证书）

```sql
CREATE TABLE cms_ssl_cert (
  id              BIGINT       NOT NULL AUTO_INCREMENT,
  tenant_id       BIGINT       NOT NULL DEFAULT 0,
  cert_name       VARCHAR(64)  NOT NULL DEFAULT ""      COMMENT "证书别名",
  cert_domain     VARCHAR(256) NOT NULL DEFAULT ""      COMMENT "绑定域名(支持多个)",
  cert_type       TINYINT      NOT NULL DEFAULT 0       COMMENT "类型:0自签1LetsEncrypt2商业CA",
  cert_pem        LONGTEXT     NOT NULL                  COMMENT "证书PEM(完整链)",
  cert_key        LONGTEXT     NOT NULL                  COMMENT "私钥PEM",
  issuer          VARCHAR(128) NOT NULL DEFAULT ""      COMMENT "颁发者",
  not_before      DATETIME     DEFAULT NULL             COMMENT "生效时间",
  not_after       DATETIME     DEFAULT NULL             COMMENT "过期时间",
  auto_renew      TINYINT(1)   NOT NULL DEFAULT 0       COMMENT "自动续签",
  deleted         TINYINT(1)   NOT NULL DEFAULT 0,
  create_id       INT          NOT NULL DEFAULT 0,
  create_name     VARCHAR(64)  NOT NULL DEFAULT "",
  create_time     DATETIME     DEFAULT CURRENT_TIMESTAMP,
  update_time     DATETIME     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  INDEX idx_tenant (tenant_id),
  INDEX idx_cert_domain (cert_domain(128)),
  INDEX idx_not_after (not_after)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT="SSL证书表";
```

### 3.4 新增 cms_site_settings 表（站点扩展配置 JSON）

对于复杂站点配置（如 ICP 备案、统计代码、SEO 全局配置），使用 KV 表：

```sql
CREATE TABLE cms_site_settings (
  id              BIGINT       NOT NULL AUTO_INCREMENT,
  tenant_id       BIGINT       NOT NULL DEFAULT 0,
  site_id         BIGINT       NOT NULL DEFAULT 0,
  setting_key     VARCHAR(64)  NOT NULL DEFAULT "",
  setting_value   LONGTEXT     NOT NULL,
  setting_group   VARCHAR(32)  NOT NULL DEFAULT "general" COMMENT "配置分组",
  deleted         TINYINT(1)   NOT NULL DEFAULT 0,
  create_time     DATETIME     DEFAULT CURRENT_TIMESTAMP,
  update_time     DATETIME     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uk_site_key (tenant_id, site_id, setting_key),
  INDEX idx_site_group (site_id, setting_group)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT="站点扩展配置";
```

---

## 四、迁移 SQL（多方言）

所有迁移文件追加到现有 `migrations/05_multi_site_init.sql`（MySQL）。Postgres / MSSQL 同步生成。

### 4.1 MySQL（完整脚本）

```sql
-- migrations/05_multi_site_init.sql (MySQL)
SET NAMES utf8mb4;
START TRANSACTION;

-- 1. cms_site 增强
ALTER TABLE cms_site
  ADD COLUMN parent_site_id  BIGINT       NOT NULL DEFAULT 0 COMMENT "父站ID",
  ADD COLUMN theme_inherit   TINYINT(1)   NOT NULL DEFAULT 1 COMMENT "是否继承父站主题",
  ADD COLUMN mobile_template VARCHAR(128) NOT NULL DEFAULT "" COMMENT "移动端模板",
  ADD COLUMN desktop_template VARCHAR(128) NOT NULL DEFAULT "" COMMENT "PC端模板",
  ADD COLUMN site_kind       TINYINT      NOT NULL DEFAULT 0 COMMENT "站点类型",
  ADD COLUMN site_code       VARCHAR(64)  NOT NULL DEFAULT "" COMMENT "站点代码",
  ADD COLUMN site_status     TINYINT      NOT NULL DEFAULT 1 COMMENT "状态:0停1启2维护",
  ADD COLUMN enable_wildcard TINYINT(1)   NOT NULL DEFAULT 0 COMMENT "启用泛域名",
  ADD COLUMN wildcard_pattern VARCHAR(128) NOT NULL DEFAULT "" COMMENT "泛域名规则",
  ADD COLUMN site_settings   TEXT         NULL COMMENT "JSON扩展配置";

ALTER TABLE cms_site
  ADD UNIQUE INDEX uk_site_code (tenant_id, site_code),
  ADD INDEX idx_parent_site (parent_site_id),
  ADD INDEX idx_site_status (site_status);

-- 2. cms_site_domain 增强
ALTER TABLE cms_site_domain
  ADD COLUMN domain_kind  TINYINT      NOT NULL DEFAULT 0 COMMENT "域名类型",
  ADD COLUMN priority    INT          NOT NULL DEFAULT 100 COMMENT "优先级",
  ADD COLUMN https_mode  TINYINT(1)   NOT NULL DEFAULT 1 COMMENT "HTTPS模式",
  ADD COLUMN cert_id     BIGINT       NOT NULL DEFAULT 0 COMMENT "证书ID",
  ADD COLUMN is_primary  TINYINT(1)   NOT NULL DEFAULT 0 COMMENT "是否主域名",
  ADD COLUMN redirect_to VARCHAR(256) NOT NULL DEFAULT "" COMMENT "跳转域名";

ALTER TABLE cms_site_domain ADD INDEX idx_priority (priority);

-- 3. cms_ssl_cert 新表
CREATE TABLE IF NOT EXISTS cms_ssl_cert (
  id BIGINT NOT NULL AUTO_INCREMENT,
  tenant_id BIGINT NOT NULL DEFAULT 0,
  cert_name VARCHAR(64) NOT NULL DEFAULT "",
  cert_domain VARCHAR(256) NOT NULL DEFAULT "",
  cert_type TINYINT NOT NULL DEFAULT 0,
  cert_pem LONGTEXT NOT NULL,
  cert_key LONGTEXT NOT NULL,
  issuer VARCHAR(128) NOT NULL DEFAULT "",
  not_before DATETIME DEFAULT NULL,
  not_after DATETIME DEFAULT NULL,
  auto_renew TINYINT(1) NOT NULL DEFAULT 0,
  deleted TINYINT(1) NOT NULL DEFAULT 0,
  create_id INT NOT NULL DEFAULT 0,
  create_name VARCHAR(64) NOT NULL DEFAULT "",
  create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
  update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  INDEX idx_tenant (tenant_id),
  INDEX idx_cert_domain (cert_domain(128)),
  INDEX idx_not_after (not_after)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT="SSL证书表";

-- 4. cms_site_settings 新表
CREATE TABLE IF NOT EXISTS cms_site_settings (
  id BIGINT NOT NULL AUTO_INCREMENT,
  tenant_id BIGINT NOT NULL DEFAULT 0,
  site_id BIGINT NOT NULL DEFAULT 0,
  setting_key VARCHAR(64) NOT NULL DEFAULT "",
  setting_value LONGTEXT NOT NULL,
  setting_group VARCHAR(32) NOT NULL DEFAULT "general",
  deleted TINYINT(1) NOT NULL DEFAULT 0,
  create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
  update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uk_site_key (tenant_id, site_id, setting_key),
  INDEX idx_site_group (site_id, setting_group)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT="站点扩展配置";

COMMIT;
```

### 4.2 PostgreSQL（openGauss / KingbaseES 同结构）

```sql
-- migrations/06_multi_site_init.sql (PostgreSQL)
BEGIN;

ALTER TABLE cms_site ADD COLUMN IF NOT EXISTS parent_site_id BIGINT NOT NULL DEFAULT 0;
ALTER TABLE cms_site ADD COLUMN IF NOT EXISTS theme_inherit BOOLEAN NOT NULL DEFAULT TRUE;
ALTER TABLE cms_site ADD COLUMN IF NOT EXISTS mobile_template VARCHAR(128) NOT NULL DEFAULT "";
ALTER TABLE cms_site ADD COLUMN IF NOT EXISTS desktop_template VARCHAR(128) NOT NULL DEFAULT "";
ALTER TABLE cms_site ADD COLUMN IF NOT EXISTS site_kind SMALLINT NOT NULL DEFAULT 0;
ALTER TABLE cms_site ADD COLUMN IF NOT EXISTS site_code VARCHAR(64) NOT NULL DEFAULT "";
ALTER TABLE cms_site ADD COLUMN IF NOT EXISTS site_status SMALLINT NOT NULL DEFAULT 1;
ALTER TABLE cms_site ADD COLUMN IF NOT EXISTS enable_wildcard BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE cms_site ADD COLUMN IF NOT EXISTS wildcard_pattern VARCHAR(128) NOT NULL DEFAULT "";
ALTER TABLE cms_site ADD COLUMN IF NOT EXISTS site_settings TEXT;

CREATE UNIQUE INDEX IF NOT EXISTS uk_site_code ON cms_site(tenant_id, site_code);
CREATE INDEX IF NOT EXISTS idx_parent_site ON cms_site(parent_site_id);
CREATE INDEX IF NOT EXISTS idx_site_status ON cms_site(site_status);

CREATE TABLE IF NOT EXISTS cms_ssl_cert (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL DEFAULT 0,
  cert_name VARCHAR(64) NOT NULL DEFAULT "",
  cert_domain VARCHAR(256) NOT NULL DEFAULT "",
  cert_type SMALLINT NOT NULL DEFAULT 0,
  cert_pem TEXT NOT NULL,
  cert_key TEXT NOT NULL,
  issuer VARCHAR(128) NOT NULL DEFAULT "",
  not_before TIMESTAMP NULL,
  not_after TIMESTAMP NULL,
  auto_renew BOOLEAN NOT NULL DEFAULT FALSE,
  deleted BOOLEAN NOT NULL DEFAULT FALSE,
  create_id INT NOT NULL DEFAULT 0,
  create_name VARCHAR(64) NOT NULL DEFAULT "",
  create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_cert_tenant ON cms_ssl_cert(tenant_id);

CREATE TABLE IF NOT EXISTS cms_site_settings (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL DEFAULT 0,
  site_id BIGINT NOT NULL DEFAULT 0,
  setting_key VARCHAR(64) NOT NULL DEFAULT "",
  setting_value TEXT NOT NULL,
  setting_group VARCHAR(32) NOT NULL DEFAULT "general",
  deleted BOOLEAN NOT NULL DEFAULT FALSE,
  create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT uk_cms_site_settings UNIQUE (tenant_id, site_id, setting_key)
);

COMMIT;
```

### 4.3 SQL Server

```sql
-- migrations/07_multi_site_init.sql (SQL Server)
BEGIN TRY
BEGIN TRAN;

IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID("cms_site") AND name = "parent_site_id")
  ALTER TABLE cms_site ADD parent_site_id BIGINT NOT NULL DEFAULT 0;
IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID("cms_site") AND name = "theme_inherit")
  ALTER TABLE cms_site ADD theme_inherit BIT NOT NULL DEFAULT 1;
-- ... 其他字段类似 ...

IF NOT EXISTS (SELECT * FROM sys.tables WHERE name = "cms_ssl_cert")
BEGIN
CREATE TABLE cms_ssl_cert (
  id BIGINT IDENTITY(1,1) PRIMARY KEY,
  tenant_id BIGINT NOT NULL DEFAULT 0,
  cert_name NVARCHAR(64) NOT NULL DEFAULT "",
  cert_domain NVARCHAR(256) NOT NULL DEFAULT "",
  cert_type TINYINT NOT NULL DEFAULT 0,
  cert_pem NVARCHAR(MAX) NOT NULL,
  cert_key NVARCHAR(MAX) NOT NULL,
  issuer NVARCHAR(128) NOT NULL DEFAULT "",
  not_before DATETIME2 NULL,
  not_after DATETIME2 NULL,
  auto_renew BIT NOT NULL DEFAULT 0,
  deleted BIT NOT NULL DEFAULT 0,
  create_id INT NOT NULL DEFAULT 0,
  create_name NVARCHAR(64) NOT NULL DEFAULT "",
  create_time DATETIME2 DEFAULT GETDATE(),
  update_time DATETIME2 DEFAULT GETDATE()
);
END;

COMMIT;
END TRY
BEGIN CATCH
ROLLBACK;
THROW;
END CATCH;
```

### 4.4 自动迁移注册

在 `app/db/cms_models.go::buildCMSModels()` 中追加：

```go
// 多站点支持
&domain.CmsSslCert{},
&domain.CmsSiteSettings{},
```

并在 `buildCMSTableComments()` 增加：

```go
"cms_ssl_cert":     "SSL证书表",
"cms_site_settings": "站点扩展配置",
```

---

## 五、域名解析器 SiteResolver

新增独立模块 `app/site/`，与现有 `app/cms/` `app/db/` 解耦。

### 5.1 模块结构

```
app/site/
├── resolver.go            # 域名 → 站点解析（核心）
├── context.go             # SiteContext（goroutine 安全的当前站点上下文）
├── middleware.go          # Beego 中间件（注入 SiteContext）
├── theme.go               # 主题继承 / PC/移动端分流
├── cache.go               # L1 内存缓存 + L2 Redis（可选）
├── cert.go                # HTTPS / SNI / 证书加载
├── service.go             # 业务 API：Get/Set/Clear
├── filter.go              # Beego Filter
├── wildcard.go            # 泛域名匹配算法
├── redirect.go            # HTTPS / 主域名跳转
└── README.md
```

### 5.2 域名解析算法

```go
// app/site/resolver.go
package site

import (
  "context"
  "strings"
  "github.com/beego/beego/v2/core/logs"
  mapper "haedu.gov.cn/cms/app/cms/mapper"
  domain "haedu.gov.cn/cms/app/cms/domain"
)

// Resolver 域名解析器
type Resolver struct {
  cache *Cache
}

// NewResolver 构造
func NewResolver() *Resolver {
  return &Resolver{cache: NewCache(50000, 3600)}
}

// Resolve 根据 host 解析站点
// 返回: 站点 / 域名绑定 / 是否默认 / error
func (r *Resolver) Resolve(ctx context.Context, host string) (*SiteContext, error) {
  host = normalizeHost(host)
  if host == "" {
    return nil, errors.New("empty host")
  }

  // 1. 精确匹配（缓存优先）
  if sc := r.cache.Get(host); sc != nil {
    return sc, nil
  }

  // 2. 精确匹配 → cms_site_domain
  mdl, do := mapper.CmsSiteDomainDo()
  domainRow, err := do.Where(
    mdl.Domain.Eq(host),
    mdl.Deleted.Is(false),
  ).First()
  if err == nil && domainRow != nil {
    sc := r.loadSite(domainRow.SiteID)
    r.cache.Set(host, sc)
    return sc, nil
  }

  // 3. 泛域名匹配（*.example.com）
  if sc := r.matchWildcard(host); sc != nil {
    r.cache.Set(host, sc)
    return sc, nil
  }

  // 4. 降级到默认站
  sc := r.loadDefault()
  if sc != nil {
    r.cache.Set(host, sc)
  }
  return sc, nil
}

// normalizeHost 规范化 host
func normalizeHost(host string) string {
  // 移除端口
  if idx := strings.Index(host, ":"); idx > 0 {
    host = host[:idx]
  }
  return strings.ToLower(strings.TrimSpace(host))
}

// matchWildcard 泛域名匹配
func (r *Resolver) matchWildcard(host string) *SiteContext {
  // 从 cms_site 加载所有 enable_wildcard=1 的站点
  // 按 wildcard_pattern 优先级匹配
  mdl, do := mapper.CmsSiteDo()
  rows, _ := do.Where(
    mdl.EnableWildcard.Is(true),
    mdl.Deleted.Is(false),
    mdl.SiteStatus.Eq(1),
  ).Find()

  for _, siteRow := range rows {
    if matchWildcardPattern(host, siteRow.WildcardPattern) {
      return r.buildContext(siteRow, "wildcard")
    }
  }
  return nil
}

// matchWildcardPattern 匹配通配符
func matchWildcardPattern(host, pattern string) bool {
  if pattern == "" || !strings.HasPrefix(pattern, "*.") {
    return false
  }
  // *.example.com 匹配 sub.example.com / a.b.example.com
  suffix := pattern[1:] // .example.com
  return strings.HasSuffix(host, suffix)
}
```

---

## 六、Beego 中间件

### 6.1 启动注册（main.go）

```go
// main.go
import (
  "haedu.gov.cn/cms/app/site"
)

func main() {
  // ... 现有初始化 ...
  // 多站点支持：在所有业务 Filter 之前
  site.Init()
  web.InsertFilter("/*", web.BeforeRouter, site.Middleware())
  web.InsertFilter("/*", web.BeforeRouter, site.HTTPSRedirectFilter())
  web.InsertFilter("/*", web.BeforeStatic, site.MobileDetectFilter())
  // ... 启动服务 ...
}
```

### 6.2 域名解析中间件

```go
// app/site/middleware.go
package site

import (
  "github.com/beego/beego/v2/server/web"
  "github.com/beego/beego/v2/server/web/context"
)

// Middleware 域名解析中间件
func Middleware() web.FilterFunc {
  return func(c *context.Context) {
    host := c.Request.Host
    sc, err := defaultResolver.Resolve(c.Request.Context(), host)
    if err != nil || sc == nil {
      logs.Warn("无法解析站点: host=%s, err=%v", host, err)
      c.Abort(404, "site not found")
      return
    }

    // 注入到 context
    c.Input.SetData("site_context", sc)

    // 设置租户上下文（兼容现有 tenant 模块）
    if sc.TenantID > 0 {
      tenant.SetCurrentTenant(c.Request.Context(), sc.TenantID)
    }

    // 设置到 SiteContext 工具方法可访问
    SetCurrentSite(c.Request.Context(), sc)
  }
}
```

### 6.3 HTTPS 重定向中间件

```go
// app/site/middleware.go

// HTTPSRedirectFilter HTTPS 强制跳转
func HTTPSRedirectFilter() web.FilterFunc {
  return func(c *context.Context) {
    sc := GetCurrentSite(c.Request.Context())
    if sc == nil {
      return
    }

    // 找出当前 host 的域名绑定
    binding := findBindingForHost(sc, c.Request.Host)
    if binding == nil {
      return
    }

    switch binding.HTTPSMode {
    case 1: // 强制 HTTPS
      if c.Request.TLS == nil {
        redirectToHTTPS(c)
        return
      }
    case 0: // 强制 HTTP
      if c.Request.TLS != nil {
        redirectToHTTP(c)
        return
      }
    }

    // 跳转主域名
    if binding.RedirectTo != "" && !strings.HasSuffix(c.Request.Host, binding.RedirectTo) {
      c.Redirect(301, "https://"+binding.RedirectTo+c.Request.URL.Path)
      return
    }

    // HSTS
    if c.Request.TLS != nil {
      c.ResponseWriter.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
    }
  }
}
```

---

## 七、SiteContext 上下文

### 7.1 数据结构

```go
// app/site/context.go
package site

import (
  "context"
  "sync"
  "haedu.gov.cn/cms/app/cms/domain"
)

// SiteContext 当前请求的站点上下文
type SiteContext struct {
  Site          *domain.CmsSite       // 站点主表
  DomainBinding *domain.CmsSiteDomain // 域名绑定
  TenantID      int64                // 租户ID
  SiteID        int64                // 站点ID
  Template      string               // 当前使用的模板（已解析继承）
  IsMobile      bool                 // 是否移动端（按UA判断）
  Device        string               // 设备类型 pc/mobile/tablet
  Host          string               // 当前 host
  IsWildcard    bool                 // 是否泛域匹配
  IsDefault     bool                 // 是否默认站 fallback
  Settings      map[string]string    // 站点 KV 配置
  InheritChain  []int64              // 主题继承链 site_id 列表
}

// 工具方法
func (sc *SiteContext) GetSetting(key, def string) string {
  if v, ok := sc.Settings[key]; ok { return v }
  return def
}

func (sc *SiteContext) Logo() string {
  if sc.Site.Logo1 != "" { return sc.Site.Logo1 }
  return sc.GetSetting("default_logo", "")
}

func (sc *SiteContext) Copyright() string {
  if sc.Site.Copyright != "" { return sc.Site.Copyright }
  return sc.GetSetting("default_copyright", "© 2026")
}
```

### 7.2 goroutine 安全访问

```go
// app/site/context.go

// 通过 context.Context 携带 SiteContext
type ctxKey int

const siteCtxKey ctxKey = 1

// WithSiteContext 注入
func WithSiteContext(ctx context.Context, sc *SiteContext) context.Context {
  return context.WithValue(ctx, siteCtxKey, sc)
}

// FromContext 取出
func FromContext(ctx context.Context) *SiteContext {
  if sc, ok := ctx.Value(siteCtxKey).(*SiteContext); ok {
    return sc
  }
  return nil
}

// GetCurrentSite 兼容旧 API（带 fallback 到全局）
func GetCurrentSite(ctx context.Context) *SiteContext {
  if sc := FromContext(ctx); sc != nil {
    return sc
  }
  return GlobalDefaultSite // 兜底
}

var (
  GlobalDefaultSite *SiteContext
  siteContextMu     sync.RWMutex
)

// SetCurrentSite 兼容旧 API
func SetCurrentSite(ctx context.Context, sc *SiteContext) {
  siteContextMu.Lock()
  GlobalDefaultSite = sc
  siteContextMu.Unlock()
}
```

---

## 八、Theme 主题继承

### 8.1 继承链算法

```
主题继承示意（按 parent_site_id 反向追溯）：

       国家级主站 (id=1, template=gov)
           ↑
       省级门户站 (id=10, theme_inherit=true, template="")
           ↑
       市级子站 (id=100, theme_inherit=true, template="city")  ← 仅覆盖 PC 模板
           ↑
       县区微站 (id=1000, theme_inherit=false, template="county")
```

### 8.2 解析逻辑

```go
// app/site/theme.go
package site

// ResolveTemplate 解析模板（含继承）
func ResolveTemplate(site *domain.CmsSite, isMobile bool) (template string, chain []int64) {
  template = site.Template
  if isMobile && site.MobileTemplate != "" {
    template = site.MobileTemplate
  } else if !isMobile && site.DesktopTemplate != "" {
    template = site.DesktopTemplate
  }

  chain = append(chain, site.SiteID)

  // 沿 parent_site_id 向上查找（直到找到 template 非空且 theme_inherit=true）
  current := site
  for current.ParentSiteID > 0 && template == "" && current.ThemeInherit {
    parent := loadSite(current.ParentSiteID)
    if parent == nil { break }
    template = parent.Template
    chain = append(chain, parent.SiteID)
    current = parent
  }

  if template == "" {
    template = loadDefaultTheme() // cms_theme.is_default = 1
  }
  return template, chain
}
```

### 8.3 模板路径解析

```go
// 模板查找顺序
func ResolveViewPath(siteID int64, viewName string) string {
  chain := getInheritChain(siteID) // 上节算法结果
  for _, sid := range chain {
    template := getSiteTemplate(sid)
    // 尝试 views/themes/<template>/<viewName>.html
    path := "themes/" + template + "/" + viewName + ".html"
    if fileExists(path) {
      return path
    }
  }
  // 兜底：默认主题
  return "themes/default/" + viewName + ".html"
}
```

### 8.4 主题切换器（后台）

后台管理员可在站点编辑页面快速切换主站点的 `template`，保存后通过缓存失效立即看到效果。

---

## 九、PC/移动端分流

### 9.1 UA 检测

```go
// app/site/middleware.go
// MobileDetectFilter 移动端检测
func MobileDetectFilter() web.FilterFunc {
  return func(c *context.Context) {
    sc := GetCurrentSite(c.Request.Context())
    if sc == nil { return }

    ua := c.Request.UserAgent()
    isMobile := isMobileUA(ua)
    sc.IsMobile = isMobile
    sc.Device = detectDevice(ua) // pc/mobile/tablet

    // 重新解析模板
    template, _ := ResolveTemplate(sc.Site, isMobile)
    sc.Template = template
  }
}

func isMobileUA(ua string) bool {
  keywords := []string{"iPhone", "iPad", "Android", "Mobile", "MQQBrowser", "IEMobile"}
  for _, k := range keywords {
    if strings.Contains(ua, k) { return true }
  }
  return false
}
```

### 9.2 移动端子域（可选）

除了 UA 检测，可选地启用移动子域 `m.example.com`：

```
cms_site_domain
├── www.example.com → 主站 PC 模板
├── m.example.com   → 主站 mobile_template
└── *.example.com   → 泛解析到子站
```

由域名直接判定（无需 UA），更适合跨设备分享。

### 9.3 自适应 vs 独立移动模板

| 方案 | 适用 |
|------|------|
| **响应式（同一模板自适应）** | 简单站、内容站 |
| **独立移动模板** | 复杂业务、电商、视频 |
| **混合**（首页响应式 + 详情独立移动模板）| 主流 CMS |

本项目默认 `mobile_template = ""` 时走自适应，配置后走独立模板。

---

## 十、缓存策略

### 10.1 多级缓存

```
L1 进程内存 cache（host → SiteContext）
    ↓ 过期 / 主动失效
L2 Redis cache（跨进程共享）
    ↓ 过期
DB cms_site + cms_site_domain
```

### 10.2 实现

```go
// app/site/cache.go
package site

import (
  "github.com/beego/beego/v2/core/logs"
  cache "github.com/zsy619/tools/xcache"
)

type Cache struct {
  l1 *cache.CacheItemModel  // 进程内存
  l2 cache.Cache            // Redis（可选）
  ttl int
}

func NewCache(size, ttl int) *Cache {
  return &Cache{
    l1: cache.NewCacheItemModel("SiteCache", "站点内存缓存", ttl),
    l2: nil, // 默认关闭，按需开启
    ttl: ttl,
  }
}

func (c *Cache) Get(host string) *SiteContext {
  if v, ok := c.l1.Get(host); ok {
    return v.(*SiteContext)
  }
  if c.l2 != nil {
    // L2 查询

  }
  return nil
}

func (c *Cache) Set(host string, sc *SiteContext) {
  c.l1.Set(host, sc)
  if c.l2 != nil {
    // L2 写入

  }
}

// InvalidateSite 站点变更时主动失效
func (c *Cache) InvalidateSite(siteID int64) {
  // 找出该站点绑定的所有 host，逐一失效
  hosts := listHostsBySiteID(siteID)
  for _, h := range hosts {
    c.l1.Delete(h)
    if c.l2 != nil { c.l2.Delete(h) }
  }
  logs.Info("失效站点缓存: siteID=%d, hosts=%v", siteID, hosts)
}
```

### 10.3 主动失效触发点

| 触发点 | 代码 |
|------|------|
| 后台保存站点配置 | `service.CmsSite::SaveOne` → `cache.InvalidateSite(siteID)` |
| 后台保存域名绑定 | `service.CmsSiteDomain::Save` → 同上 |
| 后台删除站点 | `service.CmsSite::Delete` → 同上 |
| 上传新模板 | `service.CmsTheme::Save` → 失效所有站点（按 template 名字匹配）|
| 手动清空 | 后台「缓存管理 → 清空站点缓存」按钮 |

---

## 十一、后台管理界面

### 11.1 站点列表页（views/admin/Site/Index.html）增强

- 表格新增列：站点代码、状态、维护中标识、泛域启用、移动模板、继承父站
- 顶部按钮：新建站点 / 复制站点（继承现有配置）
- 行操作：编辑 / 域名管理 / 模板设置 / 缓存刷新 / 启用禁用

### 11.2 站点编辑页（SiteEdit.html）增强

新增 Tab 页：

- 基础信息（保留）
- 模板设置（PC / 移动 / 继承）
- 域名绑定（精确 / 泛域 / HTTPS / 主域名跳转）
- SEO 配置（META / 统计代码 / robots）
- 高级设置（JSON KV 配置）

### 11.3 域名管理组件

```html
<table class="layui-table" id="domain-table">
  <thead>
    <tr>
      <th>域名</th><th>类型</th><th>优先级</th><th>HTTPS</th><th>主域名</th><th>证书</th><th>操作</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><input class="layui-input" name="domain" value="www.example.com"></td>
      <td>
        <select name="domain_kind">
          <option value="0">精确</option>
          <option value="1">泛域（*.example.com）</option>
        </select>
      </td>
      <td><input class="layui-input" name="priority" value="100"></td>
      <td>
        <select name="https_mode">
          <option value="2">自动</option>
          <option value="1">强制 HTTPS</option>
          <option value="0">强制 HTTP</option>
        </select>
      </td>
      <td><input type="checkbox" lay-skin="switch"></td>
      <td>
        <select name="cert_id">
          <option value="0">不绑定</option>
          <option value="1">主站证书 (Let's Encrypt)</option>
        </select>
      </td>
      <td><button class="layui-btn layui-btn-xs layui-btn-danger">删除</button></td>
    </tr>
  </tbody>
</table>
<button class="layui-btn" id="add-domain">+ 添加域名</button>
```

---

## 十二、运维与配置

### 12.1 conf/app.conf 新增配置

```ini
# ===== 多站点配置 =====
# 启用多站点模式
site.enabled = true
# 默认租户ID（未匹配时使用）
site.default_tenant_id = 0
# 默认站点ID（Host 匹配失败时使用）
site.default_site_id = 1
# 泛域名匹配（按需开启）
site.enable_wildcard = true
# 移动端自动检测（PC 优先，移动模板）
site.enable_mobile_detect = true
# HTTPS 强制跳转
site.https_redirect = false
# 站点信息缓存 L1（秒）
site.cache_l1_ttl = 3600
# 站点信息缓存 L2 Redis（秒，0 表示关闭）
site.cache_l2_ttl = 7200
# 多租户模式：是否严格按租户隔离（true=严格 false=共享）
site.strict_tenant_isolation = true
# 未知 host 行为：fallback / 404
site.unknown_host_action = fallback
# 日志：站点解析耗时警告阈值（毫秒）
site.warn_resolve_ms = 50
```

### 12.2 Nginx 配置（多域名 → 单后端）

```nginx
# /etc/nginx/conf.d/cms.conf
upstream cms_backend {
  server 127.0.0.1:8125;
  keepalive 64;
}

# HTTP → HTTPS 重定向
server {
  listen 80;
  listen [::]:80;
  server_name ~^(www\|m)\.(.+)$ ~^([^.]+\.)?(.+\.edu\.cn)$;
  return 301 https://$host$request_uri;
}

# HTTPS 主配置（SNI 多证书）
server {
  listen 443 ssl http2;
  server_name ~^(www\|m)\.(.+)$ ~^([^.]+\.)?(.+\.edu\.cn)$;

  ssl_protocols TLSv1.2 TLSv1.3;
  ssl_ciphers HIGH:!aNULL:!MD5;
  ssl_prefer_server_ciphers on;

  # SNI 多证书（按 host 动态加载）
  ssl_certificate     /etc/nginx/ssl/$ssl_server_name.crt;
  ssl_certificate_key /etc/nginx/ssl/$ssl_server_name.key;

  # 通用证书 fallback（acme.sh 通配符证书）
  ssl_certificate     /etc/nginx/ssl/_wildcard_.edu.cn.crt;
  ssl_certificate_key /etc/nginx/ssl/_wildcard_.edu.cn.key;

  location / {
    proxy_pass http://cms_backend;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
    proxy_set_header X-Forwarded-Host $host;
  }

  location /static/ {
    alias /opt/cms/static/;
    expires 30d;
  }
}
```

### 12.3 监控指标

```go
// 站点解析命中
siteResolutionTotal = promauto.NewCounterVec(prometheus.CounterOpts{
  Name: "cms_site_resolution_total",
  Help: "Total site resolutions",
}, ["source"]) // cache/db/wildcard/default

siteResolutionDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
  Name: "cms_site_resolution_duration_seconds",
  Help: "Site resolution latency",
  Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5},
}, ["source"])

siteCacheInvalidation = promauto.NewCounterVec(prometheus.CounterOpts{
  Name: "cms_site_cache_invalidation_total",
  Help: "Cache invalidation events",
}, ["reason"]) // save/delete/cert/theme/manual
```

---

## 十三、回滚与灰度策略

### 13.1 灰度上线

| 阶段 | 内容 | 范围 | 验证指标 |
|------|------|------|----------|
| **第 1 阶段** | 新增 site 模块基础 + 缓存 | 全量 | 启动 0 报错 |
| **第 2 阶段** | 域名精确匹配（仅 1 个租户）| 租户 ID=1 | 解析命中率 100% |
| **第 3 阶段** | 泛域名匹配 + 主题继承 | 租户 ID=1,2 | 命中延迟 < 5ms |
| **第 4 阶段** | PC/移动端分流 + HTTPS | 全量租户 | UA 分流准确率 |
| **第 5 阶段** | SSL 多证书管理 | 全量 | 证书过期告警 |
| **第 6 阶段** | 后台管理 UI 增强 | 全量管理员 | 操作可用性 |

### 13.2 Feature Flag

```ini
# 仅对租户 ID 1 / 2 启用新多站点逻辑
site.tenant_whitelist = 1,2
# 仅对启用 enable_wildcard=1 的站点启用新匹配
site.wildcard_rollout_pct = 0   # 0% = 不启用，100% = 全量
```

### 13.3 回滚方案

1. **配置回滚**：`site.enabled=false` → 走旧 `SiteByHost(host)` 单查逻辑
2. **代码回滚**：site 模块可独立 Revert
3. **数据库回滚**：新增字段 `enable_wildcard / site_code` 等可保留（不影响）
4. **新表回滚**：`cms_ssl_cert / cms_site_settings` 可直接 DROP
5. **紧急开关**：Site Resolver 顶部加 `if !site.Enabled { return legacy }` 一键回滚

---

## 十四、风险与对策

| # | 风险 | 影响 | 对策 |
|---|------|------|------|
| 1 | 泛域名匹配性能 | 启动慢 / 每请求 O(N) 扫描 | 启动时预加载到内存 + LRU 缓存 |
| 2 | 多证书管理复杂 | SNI 错配 | acme.sh 自动续签 + 监控告警（提前 30 天）|
| 3 | 主题继承深度过大 | 模板查找慢 | 限制继承深度 ≤ 3，超过则切断 |
| 4 | PC/移动端 UA 误判 | 用户体验差 | 提供手动切换链接 / URL 中显式 ?device=mobile |
| 5 | 多租户数据泄露 | 安全事故 | 严格 tenant_id 过滤 + SQL 注入防护 + 审计日志 |
| 6 | 不同方言 SQL 差异 | 跨库迁移失败 | migrations 多方言 + GORM AutoMigrate 兜底 |
| 7 | 站点数量爆炸 | 管理难 | 站点分组 + 标签 + 检索 |
| 8 | 缓存不一致 | 用户看到旧配置 | 主动失效 + 多级缓存 + TTL 兜底 |
| 9 | 域名备案合规 | 法律风险 | 后台必填备案号字段 + 验证逻辑 |
| 10 | 旧域名停用 | 用户访问失败 | 301 跳转到新域名 + 保留期配置 |

### 14.1 主题继承深度限制

```go
const MaxThemeInheritDepth = 3

func ResolveTemplate(site *domain.CmsSite, isMobile bool) (string, []int64) {
  depth := 0
  current := site
  for depth < MaxThemeInheritDepth {
    if current.ParentSiteID == 0 || !current.ThemeInherit { break }
    parent := loadSite(current.ParentSiteID)
    if parent == nil { break }
    current = parent
    depth++
  }
  // ...
}
```

---

## 十五、里程碑与工作量估算

### 15.1 里程碑

| 里程碑 | 内容 | 周期 | 验收 |
|------|------|------|------|
| **M1** | 数据模型扩展（cms_site + cms_site_domain + 2 张新表）+ 多方言迁移脚本 | 1 周 | GORM AutoMigrate 无报错 |
| **M2** | `app/site/` 模块基础（resolver / cache / context） | 1 周 | 单租户域名解析命中 < 5ms |
| **M3** | 域名精确匹配 + 中间件接入 | 0.5 周 | 多域名正确路由 |
| **M4** | 泛域名匹配 + 主题继承 | 1 周 | *.edu.cn 解析正确 |
| **M5** | PC/移动端 UA 检测 + 模板分流 | 0.5 周 | 移动端走 mobile_template |
| **M6** | HTTPS 强制跳转 + SNI 多证书 | 1 周 | 自动续签证书 |
| **M7** | 后台管理 UI 增强（域名 / 模板 / 缓存） | 1 周 | 后台可完成全流程 |
| **M8** | 监控指标 + 灰度开关 + 文档 | 0.5 周 | Prometheus 指标 + 灰度上线 |

**总计约 6.5 周**（2 名 Go 后端 + 1 名前端 + 0.5 名运维）。

### 15.2 优先级

1. **P0**：M1 + M2 + M3（数据 + 解析 + 中间件）
2. **P1**：M4 + M5（泛域 + 移动端）
3. **P2**：M6 + M7 + M8（HTTPS + UI + 监控）

### 15.3 验收清单

- [ ] 多租户下不同租户的站点互不可见
- [ ] 多域名（精确）正确解析到对应站点
- [ ] 泛域名（*.example.com）正确匹配子域
- [ ] 主题继承链 ≤ 3 级，子站可继承父站模板
- [ ] 移动端 UA 自动识别并使用 mobile_template
- [ ] HTTPS 强制跳转正常
- [ ] SNI 多证书配置生效，证书过期前 30 天告警
- [ ] 后台保存站点配置后立即生效（缓存失效）
- [ ] 未知 host fallback 到默认站
- [ ] Prometheus 指标可观测
- [ ] 多方言数据库迁移无报错（MySQL / PG / MSSQL / Kingbase / GaussDB）
- [ ] 灰度开关可控制范围
- [ ] 紧急回滚一键关闭

---

## 附录 A：典型场景示例

### A.1 河南省智慧就业平台（多级子站）

```
主站：haedu.gov.cn (cms_site.id=1)
├── /xm/  厦门市子站：xm.haedu.gov.cn (cms_site.id=10, parent=1)
├── /zz/  郑州市子站：zz.haedu.gov.cn (cms_site.id=11, parent=1)
├── /xx/  学校门户：*.edu.cn (cms_site.id=100, parent=1, enable_wildcard=1, wildcard_pattern=*.edu.cn)
└── /m/   移动端：m.haedu.gov.cn (cms_site.id=1, mobile_template=gov-mobile)
```

### A.2 单租户多品牌

```
一家教育集团，多个品牌站：
主品牌：www.brand1.com (cms_site.id=1, template=brand1)
子品牌：www.brand2.com (cms_site.id=2, parent=1, theme_inherit=true, template=brand2)
└─ 共用基础模板，仅覆盖差异化部分
```

### A.3 多租户 SaaS 模式

```
SaaS 平台，每个租户独立站点：
租户 A：{tenant}.example-a.com (tenant_id=100, site_code=a)
租户 B：{tenant}.example-b.com (tenant_id=200, site_code=b)
└─ 通过 tenant_id 严格隔离数据
```

---

## 附录 B：参考链接

- Beego 中间件：https://beego.vip/docs/web/filter/
- LetsEncrypt 通配符证书：https://letsencrypt.org/
- Nginx SNI 多证书：https://nginx.org/en/docs/http/configuring_https_servers.html
- GORM 多方言：https://gorm.io/docs/
- 项目现有站点模块：`app/cms/service/api_site.go`
- 项目现有站点缓存：`app/tool/cache.go::SiteCache`
- 项目现有域名解析：`controllers/www/b.go::SiteByHost`
- 多语言方案：`docs/i18n_design.md`（站点切换 + 语言切换可叠加）
