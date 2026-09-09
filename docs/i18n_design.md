# CMS 内容多语言支持（i18n）实施方案

> 基于 Beego v2.3.8 + GORM 多方言 + 多租户的 CMS 内容管理系统，结合 https://chat.deepseek.com/share/lukvbu834bxzfk0t05 中的多语言设计思路，给出贴合本项目实际的端到端可执行方案。

---

## 目录

- [一、目标与设计原则](#一目标与设计原则)
- [二、总体架构](#二总体架构)
- [三、数据模型设计](#三数据模型设计)
- [四、迁移 SQL（多方言）](#四迁移-sql多方言)
- [五、Go 代码结构](#五go-代码结构)
- [六、Beego i18n 集成](#六beego-i18n-集成)
- [七、Domain / Mapper / Service 三层落地](#七domain--mapper--service-三层落地)
- [八、REST API 设计](#八rest-api-设计)
- [九、模板侧实现](#九模板侧实现)
- [十、前端（LayUI）集成](#十前端layui集成)
- [十一、运维与配置](#十一运维与配置)
- [十二、回滚与灰度策略](#十二回滚与灰度策略)
- [十三、风险与对策](#十三风险与对策)
- [十四、里程碑与工作量估算](#十四里程碑与工作量估算)

---

## 一、目标与设计原则

### 1.1 目标

1. **UI 多语言**：后台 / 前台模板文案按用户语言切换（中、英、繁、其它按需扩展）
2. **内容多语言**：业务数据（文章、栏目、专题、广告、标签等）按语言存储与展示
3. **租户维度**：每个租户可独立配置默认语言和支持语言集
4. **优雅降级**：翻译缺失时自动 fallback 到主语言，URL 无 lang 参数时按租户默认 / 浏览器语言
5. **零侵入上线**：复用现有 domain/mapper 三层，新增翻译表即可，旧业务平滑迁移
6. **跨方言兼容**：MySQL / PostgreSQL / SQL Server / Kingbase / GaussDB 全支持

### 1.2 设计原则

| 原则 | 落地 |
|------|------|
| **约定优于配置** | 默认 zh-CN，URL/Header/Query 任意一层带 lang 即可切换 |
| **翻译与数据分离** | 新增翻译表 `cms_xxx_i18n`，主表字段保留作默认语言值 |
| **多租户感知** | `cms_tenant` 增加 `default_lang` 与 `support_langs` 字段 |
| **不强侵入既有代码** | 翻译读取走装饰层 service，mapper 不变 |
| **翻译缓存** | 启用 Beego 内置 i18n 缓存 + 项目级多语种内存缓存 |
| **性能优先** | 翻译读取走索引 `(source_type, source_id, lang)` 复合唯一索引 |

---

## 二、总体架构

```
                     ┌─────────────────────────────────────────┐
                     │              Browser / Client            │
                     │  URL: /:lang/article/list?lang=en-US      │
                     │  Header: Accept-Language: en-US          │
                     └────────────────────┬────────────────────┘
                                          │
                                          ▼
┌──────────────────────────────────────────────────────────────────────┐
│                       Beego HTTP / Filter                            │
│  ┌──────────────────┐   ┌──────────────────┐   ┌─────────────────┐  │
│  │ LanguageFilter   │──▶│ LocaleResolver   │──▶│ beego i18n Set  │  │
│  │ (中间件)         │   │ (URL/Header)     │   │ (zh-CN/en-US..) │  │
│  └──────────────────┘   └──────────────────┘   └─────────────────┘  │
└──────────────────────────────────┬───────────────────────────────────┘
                                   │
        ┌──────────────────────────┼─────────────────────────────┐
        ▼                          ▼                             ▼
┌────────────────┐       ┌──────────────────┐          ┌────────────────┐
│ Templates      │       │ Controllers      │          │ Service 层     │
│ {{T "key"}}    │       │ ctx.Input.Su   │          │ i18n.Get()    │
│ ({{lang}})    │       │ 数据 lang 注入   │          │ + 翻译表查询  │
└────────────────┘       └──────────────────┘          └────────────────┘
                                                              │
                                                              ▼
                                              ┌────────────────────────────────┐
│                                              │   MySQL / PG / MSSQL / ..     │
│                                              │   cms_xxx_i18n                │
│                                              │   (source_type, source_id,    │
│                                              │    lang, field, value)        │
└──────────────────────────────────────────────────────────────────────┘
```

---

## 三、数据模型设计

### 3.1 总体策略：主表 + 翻译表分离

| 方案 | 优点 | 缺点 | 推荐度 |
|------|------|------|--------|
| **A. 单表多列 (title_zh, title_en ...)** | 单表查询、无需 JOIN | 列爆炸、零翻译占 NULL、索引膨胀 | ⭐⭐ |
| **B. 主表 + EAV 翻译表** | 灵活、字段任意扩展 | 需要聚合查询、多行翻译 | ⭐⭐⭐⭐ |
| **C. 主表 JSON 字段** | 无需新表 | 数据库 JSON 函数支持差异、检索差 | ⭐⭐ |
| **D. 主表 + 业务行级翻译表 (article_i18n)** | 业务清晰、查询简单、缓存友好 | 新增表、迁移多 | ⭐⭐⭐⭐⭐ |

**推荐方案 D**：每张可翻译业务表配套一张 `cms_xxx_i18n`，按业务行存储翻译字段；主表保留**默认语言**字段作主表 fallback。

### 3.2 通用翻译表（轻量场景）

```sql
cms_translation
├── id           BIGINT PK
├── tenant_id    BIGINT NOT NULL DEFAULT 0
├── source_type  VARCHAR(64)   NOT NULL   -- 如 article / category / tag / link / ads
├── source_id    BIGINT        NOT NULL   -- 业务主键
├── lang         VARCHAR(16)   NOT NULL   -- zh-CN / en-US / zh-HK
├── field        VARCHAR(64)   NOT NULL   -- 字段名如 title / summary / seo_keywords
├── value        TEXT          NOT NULL
├── deleted      TINYINT(1)    NOT NULL DEFAULT 0
├── create_id    INT           DEFAULT 0
├── create_name  VARCHAR(64)   DEFAULT ""
├── update_id    INT           DEFAULT 0
├── update_name  VARCHAR(64)   DEFAULT ""
├── create_time  DATETIME
├── update_time  DATETIME
├── UNIQUE INDEX uk_stsl (source_type, source_id, lang, field)
├── INDEX idx_lang (lang)
├── INDEX idx_tenant (tenant_id)
```

适用场景：广告位 / 链接 / 标签 / 简短字段，**轻量 + 单表覆盖**。

### 3.3 业务行级翻译表（重量级场景，推荐）

以文章为例 `cms_article_i18n`：

```sql
cms_article_i18n
├── id               BIGINT PK
├── tenant_id        BIGINT NOT NULL DEFAULT 0
├── article_id       BIGINT NOT NULL       -- 关联 cms_article.id
├── lang             VARCHAR(16) NOT NULL  -- zh-CN / en-US
├── title            VARCHAR(64) NOT NULL
├── sub_title        VARCHAR(128) DEFAULT ""
├── summary          VARCHAR(512) DEFAULT ""
├── content          LONGTEXT
├── seo_title        VARCHAR(128) DEFAULT ""
├── seo_keywords     VARCHAR(255) DEFAULT ""
├── seo_description  VARCHAR(512) DEFAULT ""
├── author           VARCHAR(64) DEFAULT ""
├── source           VARCHAR(64) DEFAULT ""
├── deleted          TINYINT(1) NOT NULL DEFAULT 0
├── create_id        INT DEFAULT 0
├── create_name      VARCHAR(64) DEFAULT ""
├── update_id        INT DEFAULT 0
├── update_name      VARCHAR(64) DEFAULT ""
├── create_time      DATETIME
├── update_time      DATETIME
├── UNIQUE INDEX uk_atl (article_id, lang)
├── INDEX idx_lang_tenant (lang, tenant_id)
```

同样模式可推广到：

- `cms_article_category_i18n`（栏目名、SEO、模板说明）
- `cms_topic_i18n`（专题标题、描述）
- `cms_tag_i18n`（标签名）
- `cms_ads_i18n`（广告标题、副标题、链接文案）
- `cms_link_i18n`（链接名、说明）
- `cms_theme_i18n`（主题名、说明）
- `cms_site_i18n`（站点名、描述）

### 3.4 租户语言配置

`cms_tenant` 增加字段：

```sql
ALTER TABLE cms_tenant ADD COLUMN default_lang VARCHAR(16) NOT NULL DEFAULT "zh-CN";
ALTER TABLE cms_tenant ADD COLUMN support_langs VARCHAR(255) NOT NULL DEFAULT "zh-CN,en-US";
-- support_langs 为逗号分隔的语言代码列表
```

### 3.5 用户语言偏好

`cms_admin` 增加字段：

```sql
ALTER TABLE cms_admin ADD COLUMN language VARCHAR(16) NOT NULL DEFAULT "";
-- 空表示跟随租户默认
```

### 3.6 domain struct 落地（示例 cms_article_i18n）

```go
// app/cms/domain/cms_article_i18n.gen.go
type CmsArticleI18n struct {
    ID              int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement;comment:主键ID"`
    TenantID        int64     `gorm:"column:tenant_id;type:bigint;not null;default:0;index;comment:租户ID"`
    ArticleID       int64     `gorm:"column:article_id;type:bigint;not null;index;comment:文章ID"`
    Lang            string    `gorm:"column:lang;type:varchar(16);not null;index;comment:语言代码 zh-CN/en-US"`
    Title           string    `gorm:"column:title;type:varchar(64);not null;comment:标题"`
    SubTitle        string    `gorm:"column:sub_title;type:varchar(128);comment:副标题"`
    Summary         string    `gorm:"column:summary;type:varchar(512);comment:摘要"`
    Content         string    `gorm:"column:content;type:longtext;comment:正文"`
    SEOTitle        string    `gorm:"column:seo_title;type:varchar(128);comment:SEO标题"`
    SEOKeywords     string    `gorm:"column:seo_keywords;type:varchar(255);comment:SEO关键词"`
    SEODescription  string    `gorm:"column:seo_description;type:varchar(512);comment:SEO描述"`
    Author          string    `gorm:"column:author;type:varchar(64);comment:作者"`
    Source          string    `gorm:"column:source;type:varchar(64);comment:来源"`
    Deleted         bool      `gorm:"column:deleted;type:tinyint(1);not null;default:0;index;comment:逻辑删除"`
    CreateID        int32     `gorm:"column:create_id;type:int;default:0;comment:创建人ID"`
    CreateName      string    `gorm:"column:create_name;type:varchar(64);default:0;comment:创建人姓名"`
    UpdateID        int32     `gorm:"column:update_id;type:int;default:0;comment:更新人ID"`
    UpdateName      string    `gorm:"column:update_name;type:varchar(64);default:0;comment:更新人姓名"`
    CreateTime      time.Time `gorm:"column:create_time;type:datetime;comment:创建时间"`
    UpdateTime      time.Time `gorm:"column:update_time;type:datetime;comment:更新时间"`
}

func (CmsArticleI18n) TableName() string { return "cms_article_i18n" }
```

---

## 四、迁移 SQL（多方言）

所有迁移文件追加到现有 `migrations/05_i18n_init.sql`（MySQL）。Postgres / MSSQL 同步生成。

### 4.1 MySQL（含 utf8mb4）

```sql
-- migrations/05_i18n_init.sql (MySQL)
SET NAMES utf8mb4;

-- 通用翻译表（轻量场景）
CREATE TABLE IF NOT EXISTS `cms_translation` (
  `id`           BIGINT       NOT NULL AUTO_INCREMENT,
  `tenant_id`    BIGINT       NOT NULL DEFAULT 0,
  `source_type`  VARCHAR(64)  NOT NULL DEFAULT "",
  `source_id`    BIGINT       NOT NULL DEFAULT 0,
  `lang`         VARCHAR(16)  NOT NULL DEFAULT "zh-CN",
  `field`        VARCHAR(64)  NOT NULL DEFAULT "",
  `value`        TEXT         NOT NULL,
  `deleted`      TINYINT(1)   NOT NULL DEFAULT 0,
  `create_id`    INT          NOT NULL DEFAULT 0,
  `create_name`  VARCHAR(64)  NOT NULL DEFAULT "",
  `update_id`    INT          NOT NULL DEFAULT 0,
  `update_name`  VARCHAR(64)  NOT NULL DEFAULT "",
  `create_time`  DATETIME     DEFAULT CURRENT_TIMESTAMP,
  `update_time`  DATETIME     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_stsl` (`source_type`,`source_id`,`lang`,`field`),
  KEY `idx_lang` (`lang`),
  KEY `idx_tenant` (`tenant_id`),
  KEY `idx_deleted` (`deleted`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT="通用翻译表";

-- 文章翻译表
CREATE TABLE IF NOT EXISTS `cms_article_i18n` (
  `id`              BIGINT       NOT NULL AUTO_INCREMENT,
  `tenant_id`       BIGINT       NOT NULL DEFAULT 0,
  `article_id`      BIGINT       NOT NULL DEFAULT 0,
  `lang`            VARCHAR(16)  NOT NULL DEFAULT "zh-CN",
  `title`           VARCHAR(64)  NOT NULL DEFAULT "",
  `sub_title`       VARCHAR(128) NOT NULL DEFAULT "",
  `summary`         VARCHAR(512) NOT NULL DEFAULT "",
  `content`         LONGTEXT     NOT NULL,
  `seo_title`       VARCHAR(128) NOT NULL DEFAULT "",
  `seo_keywords`    VARCHAR(255) NOT NULL DEFAULT "",
  `seo_description` VARCHAR(512) NOT NULL DEFAULT "",
  `author`          VARCHAR(64)  NOT NULL DEFAULT "",
  `source`          VARCHAR(64)  NOT NULL DEFAULT "",
  `deleted`         TINYINT(1)   NOT NULL DEFAULT 0,
  `create_id`       INT          NOT NULL DEFAULT 0,
  `create_name`     VARCHAR(64)  NOT NULL DEFAULT "",
  `update_id`       INT          NOT NULL DEFAULT 0,
  `update_name`     VARCHAR(64)  NOT NULL DEFAULT "",
  `create_time`     DATETIME     DEFAULT CURRENT_TIMESTAMP,
  `update_time`     DATETIME     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_atl` (`article_id`,`lang`),
  KEY `idx_lang_tenant` (`lang`,`tenant_id`),
  KEY `idx_deleted` (`deleted`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT="文章翻译表";

-- 栏目 / 专题 / 标签 / 广告 / 链接 / 主题 / 站点翻译表结构类似
-- 省略重复样板代码

-- 租户语言配置
ALTER TABLE `cms_tenant` ADD COLUMN `default_lang` VARCHAR(16) NOT NULL DEFAULT "zh-CN" COMMENT "默认语言";
ALTER TABLE `cms_tenant` ADD COLUMN `support_langs` VARCHAR(255) NOT NULL DEFAULT "zh-CN,en-US" COMMENT "支持语言列表";

-- 管理员语言偏好
ALTER TABLE `cms_admin` ADD COLUMN `language` VARCHAR(16) NOT NULL DEFAULT "" COMMENT "用户偏好语言";
```

### 4.2 PostgreSQL（openGauss / KingbaseES 同结构）

```sql
-- migrations/06_i18n_init.sql (PostgreSQL)
CREATE TABLE IF NOT EXISTS cms_translation (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NOT NULL DEFAULT 0,
  source_type VARCHAR(64) NOT NULL DEFAULT "",
  source_id BIGINT NOT NULL DEFAULT 0,
  lang VARCHAR(16) NOT NULL DEFAULT "zh-CN",
  field VARCHAR(64) NOT NULL DEFAULT "",
  value TEXT NOT NULL,
  deleted BOOLEAN NOT NULL DEFAULT FALSE,
  create_id INT NOT NULL DEFAULT 0,
  create_name VARCHAR(64) NOT NULL DEFAULT "",
  update_id INT NOT NULL DEFAULT 0,
  update_name VARCHAR(64) NOT NULL DEFAULT "",
  create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT uk_cms_translation UNIQUE (source_type, source_id, lang, field)
);
CREATE INDEX idx_cms_translation_lang ON cms_translation(lang);
CREATE INDEX idx_cms_translation_tenant ON cms_translation(tenant_id);

-- 字段注释（PG 特有）
COMMENT ON TABLE cms_translation IS "通用翻译表";
COMMENT ON COLUMN cms_translation.lang IS "语言代码";
-- ... 其他字段
```

### 4.3 SQL Server

```sql
-- migrations/07_i18n_init.sql (SQL Server)
IF NOT EXISTS (SELECT * FROM sys.tables WHERE name = "cms_translation")
BEGIN
CREATE TABLE cms_translation (
  id BIGINT IDENTITY(1,1) PRIMARY KEY,
  tenant_id BIGINT NOT NULL DEFAULT 0,
  source_type NVARCHAR(64) NOT NULL DEFAULT "",
  source_id BIGINT NOT NULL DEFAULT 0,
  lang NVARCHAR(16) NOT NULL DEFAULT "zh-CN",
  field NVARCHAR(64) NOT NULL DEFAULT "",
  value NVARCHAR(MAX) NOT NULL,
  deleted BIT NOT NULL DEFAULT 0,
  create_id INT NOT NULL DEFAULT 0,
  create_name NVARCHAR(64) NOT NULL DEFAULT "",
  update_id INT NOT NULL DEFAULT 0,
  update_name NVARCHAR(64) NOT NULL DEFAULT "",
  create_time DATETIME2 DEFAULT GETDATE(),
  update_time DATETIME2 DEFAULT GETDATE()
);
CREATE UNIQUE INDEX uk_stsl ON cms_translation(source_type, source_id, lang, field);
END
```

### 4.4 通过 cms.migrate=true 自动建表

在 `app/db/cms_models.go` 的 `buildCMSModels()` 中追加：

```go
// i18n 翻译表
&domain.CmsTranslation{},
&domain.CmsArticleI18n{},
&domain.CmsArticleCategoryI18n{},
&domain.CmsTopicI18n{},
&domain.CmsTagI18n{},
&domain.CmsAdsI18n{},
&domain.CmsLinkI18n{},
&domain.CmsThemeI18n{},
&domain.CmsSiteI18n{},
```

并在 `buildCMSTableComments()` 增加对应注释，自动迁移即可生效。

---

## 五、Go 代码结构

新增独立模块 `app/i18n/`，与现有 `app/cms/`、`app/db/` 解耦。

```
app/i18n/
├── locale.go              # 语言上下文 / 解析
├── resolver.go            # 多级解析（URL > Header > Cookie > 租户默认 > 默认）
├── filter.go              # Beego Filter 中间件
├── manager.go             # 翻译管理器（内存缓存 + DB 回源）
├── service.go             # 业务用：GetTranslation / SaveTranslation
├── message.go             # Beego i18n 初始化（conf/locale_*.ini）
├── conf/
│   ├── locale_zh-CN.ini   # 内置中文文案
│   ├── locale_en-US.ini   # 内置英文文案
│   ├── locale_zh-HK.ini   # 内置繁体
│   └── locale_ja-JP.ini   # 日文（可选）
└── README.md
```

### 5.1 conf/locale_zh-CN.ini 示例

```ini
[Common]
app_name = 内容管理平台
welcome = 欢迎回来
login = 登录
logout = 退出
save = 保存
cancel = 取消
confirm = 确认
delete = 删除
edit = 编辑
add = 新增
search = 搜索
reset = 重置
submit = 提交

[Article]
title = 标题
category = 所属栏目
status = 状态
status_draft = 草稿
status_published = 已发布
status_auditing = 审核中
i18n_missing = 当前语言暂无翻译，已显示默认语言版本

[Error]
not_found = 资源不存在
forbidden = 权限不足
server_error = 服务器繁忙，请稍后重试
```

### 5.2 conf/locale_en-US.ini 示例

```ini
[Common]
app_name = Content Management Platform
welcome = Welcome back
login = Sign In
logout = Sign Out
save = Save
cancel = Cancel
confirm = Confirm
delete = Delete
edit = Edit
add = Add
search = Search
reset = Reset
submit = Submit

[Article]
title = Title
category = Category
status = Status
status_draft = Draft
status_published = Published
status_auditing = Auditing
i18n_missing = Translation not available, fallback to default

[Error]
not_found = Not Found
forbidden = Forbidden
server_error = Server busy, please retry later
```

---

## 六、Beego i18n 集成

### 6.1 初始化（在 main.go 启动时）

```go
// main.go
import (
  "haedu.gov.cn/cms/app/i18n"
  _ "github.com/beego/beego/v2/server/web"
)

func main() {
  // ... 现有初始化 ...
  i18n.InitI18n("./app/i18n/conf")
  web.InsertFilter("/*", web.BeforeRouter, i18n.LocaleFilter())
  // ... 启动服务 ...
}
```

### 6.2 LocaleFilter 实现

```go
// app/i18n/filter.go
package i18n

import (
  "github.com/beego/beego/v2/server/web"
  "github.com/beego/beego/v2/server/web/context"
  ctx "context"
)

var supportedLangs = []string{"zh-CN", "en-US", "zh-HK", "ja-JP"}

func LocaleFilter() web.FilterFunc {
  return func(c *context.Context) {
  lang := resolveLang(c)
  c.Input.SetData("lang", lang)
  // 设置 beego i18n
  web.I18n.SetLang(c, lang)
  }
}

func resolveLang(c *context.Context) string {
  // 1. URL 路径前缀 /:lang/...
  if lang := c.Input.Param(":lang"); lang != "" && IsSupported(lang) {
    return lang
  }
  // 2. Query 参数 ?lang=en-US
  if lang := c.Input.Query("lang"); lang != "" && IsSupported(lang) {
    return lang
  }
  // 3. Cookie
  if ck, _ := c.Request.Cookie("lang"); ck != nil && IsSupported(ck.Value) {
    return ck.Value
  }
  // 4. HTTP Header Accept-Language (协商)
  if al := c.Request.Header.Get("Accept-Language"); al != "" {
    if lang := negotiateAccept(al); IsSupported(lang) {
      return lang
    }
  }
  // 5. 租户默认 / 系统默认
  if tenantID, ok := c.Input.GetData("tenant_id").(int64); ok && tenantID > 0 {
    if lang := getTenantDefaultLang(tenantID); lang != "" {
      return lang
    }
  }
  return "zh-CN"
}
```

### 6.3 语言协商

```go
// app/i18n/resolver.go
func negotiateAccept(header string) string {
  // header 形如: en-US,en;q=0.9,zh-CN;q=0.8
  parts := strings.Split(header, ",")
  best := struct{ lang string; q float64 }{"", 0}
  for _, p := range parts {
    segs := strings.Split(strings.TrimSpace(p), ";")
    if len(segs) == 0 {
      continue
    }
    q := 1.0
    if len(segs) > 1 {
      fmt.Sscanf(segs[1], "q=%f", &q)
    }
    if q > best.q {
      best.lang, best.q = segs[0], q
    }
  }
  return best.lang
}

func IsSupported(lang string) bool {
  for _, l := range supportedLangs {
    if strings.EqualFold(l, lang) {
      return true
    }
  }
  return false
}
```

### 6.4 message.go（Beego i18n 初始化）

```go
// app/i18n/message.go
package i18n

import (
  "embed"
  "github.com/beego/beego/v2/server/web"
)

//go:embed conf/*.ini
var localeFS embed.FS

func InitI18n(dir string) {
  // 方式1: 使用 embed.FS（推荐，无需依赖运行时路径）
  for _, lang := range supportedLangs {
    data, _ := localeFS.ReadFile("conf/locale_" + lang + ".ini")
    if len(data) > 0 {
      web.I18n.AddTranslation(lang, parseIni(data))
    }
  }
  // 启用多语言缓存
  web.I18n.SetMessage(&Messages{mem: make(map[string]map[string]string)})
}
```

---

## 七、Domain / Mapper / Service 三层落地

### 7.1 域名 / 映射 与现有体系一致

新增 domain `cms_article_i18n.gen.go`、`cms_translation.gen.go` 等，在 `app/cms/mapper/a.go` 注册 `CmsArticleI18nDo()` / `CmsTranslationDo()`。

### 7.2 service 层：装饰器模式

保持现有 `cms_article.go` 调用模式，向后兼容；新增 `cms_article_i18n.go` 提供多语言读写：

```go
// app/cms/service/cms_article_i18n.go
package service

import (
  "context"
  "time"
  "haedu.gov.cn/cms/app/cms/domain"
  "haedu.gov.cn/cms/app/cms/mapper"
  "haedu.gov.cn/cms/app/i18n"
)

// ArticleI18nSave 保存文章翻译
func ArticleI18nSave(articleID int64, lang string, data *ArticleI18nDTO) error {
  mdl, do := mapper.CmsArticleI18nDo()
  now := time.Now()
  return do.Where(mdl.ArticleID.Eq(articleID), mdl.Lang.Eq(lang)).Save(
    &domain.CmsArticleI18n{
      ArticleID:       articleID,
      Lang:            lang,
      Title:           data.Title,
      SubTitle:        data.SubTitle,
      Summary:         data.Summary,
      Content:         data.Content,
      SEOTitle:        data.SEOTitle,
      SEOKeywords:     data.SEOKeywords,
      SEODescription:  data.SEODescription,
      Author:          data.Author,
      Source:          data.Source,
      UpdateTime:      now,
    },
  )
}

// ArticleGetWithLang 取文章 + 当前语言翻译（缺翻译 fallback 主表）
func ArticleGetWithLang(ctx context.Context, articleID int64) (*ArticleVO, error) {
  lang := i18n.FromContext(ctx) // 获取当前请求语言
  mdl, do := mapper.CmsArticleI18nDo()

  // 主表行
  article, err := ArticleGet(articleID)
  if err != nil {
    return nil, err
  }

  // 翻译行（不存在则用主表默认语言值）
  i18nRow, _ := do.Where(mdl.ArticleID.Eq(articleID), mdl.Lang.Eq(lang), mdl.Deleted.Is(false)).First()
  if i18nRow != nil {
    article.Title = orDefault(i18nRow.Title, article.Title)
    article.SubTitle = orDefault(i18nRow.SubTitle, article.SubTitle)
    article.Summary = orDefault(i18nRow.Summary, article.Summary)
    article.Content = orDefault(i18nRow.Content, article.Content)
    article.Author = orDefault(i18nRow.Author, article.Author)
    // SEO 字段同上
  }
  article.Lang = lang
  return article, nil
}

// ArticleListWithLang 取列表 + 当前语言翻译
func ArticleListWithLang(ctx context.Context, cond ArticleListCond) ([]*ArticleVO, int64, error) {
  lang := i18n.FromContext(ctx)
  // 1. 查主表（按 tenant / site / category / status ...）
  list, total, err := ArticleList(cond)
  if err != nil {
    return nil, 0, err
  }
  if len(list) == 0 {
    return list, total, nil
  }
  // 2. 批量查翻译（一次 IN 查询）
  ids := make([]int64, len(list))
  for i, a := range list { ids[i] = a.ArticleID }
  mdl, do := mapper.CmsArticleI18nDo()
  rows, _ := do.Where(mdl.ArticleID.In(ids), mdl.Lang.Eq(lang), mdl.Deleted.Is(false)).Find()
  // 3. 内存合并
  trans := map[int64]*domain.CmsArticleI18n{}
  for _, r := range rows { trans[r.ArticleID] = r }
  for _, a := range list {
    if r, ok := trans[a.ArticleID]; ok {
      mergeArticleLang(a, r)
    }
  }
  return list, total, nil
}
```

### 7.3 cms_translation 通用版（轻量场景）

```go
// app/cms/service/cms_translation.go
func TranslationGet(sourceType string, sourceID int64, lang string) (map[string]string, error) {
  mdl, do := mapper.CmsTranslationDo()
  rows, err := do.Where(
    mdl.SourceType.Eq(sourceType),
    mdl.SourceID.Eq(sourceID),
    mdl.Lang.Eq(lang),
    mdl.Deleted.Is(false),
  ).Find()
  if err != nil {
    return nil, err
  }
  out := make(map[string]string, len(rows))
  for _, r := range rows {
    out[r.Field] = r.Value
  }
  return out, nil
}

func TranslationBatchGet(sourceType string, sourceIDs []int64, lang string) (map[int64]map[string]string, error) {
  // 一次 IN 查询，避免 N+1
  mdl, do := mapper.CmsTranslationDo()
  rows, err := do.Where(
    mdl.SourceType.Eq(sourceType),
    mdl.SourceID.In(sourceIDs),
    mdl.Lang.Eq(lang),
    mdl.Deleted.Is(false),
  ).Find()
  // 聚合 ...
  return nil, nil
}
```

---

## 八、REST API 设计

> 在 `controllers/www/` 下新增 `/api/v2/i18n/*` 与 `/api/v2/article/*` 多语言版本。

### 8.1 通用规范

所有 API 支持以下三种语言切换方式（优先级从高到低）：

1. URL 路径前缀：`/zh-CN/api/v2/article/list` 或 `/en-US/api/v2/article/list`
2. 查询参数：`?lang=en-US`
3. HTTP Header：`Accept-Language: en-US,zh-CN;q=0.8`

响应中始终返回 `Lang` 字段，便于前端识别：

```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "lang": "en-US",
    "fallback": false,
    "items": [...]
  }
}
```

`fallback=true` 表示当前语言无翻译，已返回主表默认语言。

### 8.2 接口清单

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v2/i18n/lang/list` | 获取当前租户支持的语言列表 |
| GET | `/api/v2/i18n/lang/current` | 获取当前请求解析到的语言 |
| POST | `/api/v2/i18n/article/save` | 后台保存文章翻译 |
| GET | `/api/v2/i18n/article/get` | 后台获取文章翻译（按 lang） |
| POST | `/api/v2/i18n/category/save` | 保存栏目翻译 |
| POST | `/api/v2/i18n/tag/save` | 保存标签翻译 |
| POST | `/api/v2/i18n/ads/save` | 保存广告翻译 |
| POST | `/api/v2/i18n/link/save` | 保存链接翻译 |
| POST | `/api/v2/i18n/topic/save` | 保存专题翻译 |
| POST | `/api/v2/i18n/site/save` | 保存站点翻译 |
| POST | `/api/v2/i18n/theme/save` | 保存主题翻译 |
| GET | `/api/v2/article/list` | 前台获取文章列表（多语言版） |
| GET | `/api/v2/article/get` | 前台获取文章详情（多语言版） |
| GET | `/api/v2/category/list` | 前台栏目列表（多语言版） |
| GET | `/api/v2/tag/list` | 前台标签列表（多语言版） |
| GET | `/api/v2/ads/list` | 前台广告列表（多语言版） |
| GET | `/api/v2/link/list` | 前台链接列表（多语言版） |

### 8.3 示例：保存文章翻译

```http
POST /api/v2/i18n/article/save
Content-Type: application/json

{
  "article_id": 1024,
  "lang": "en-US",
  "title": "Smart Job Fair 2023",
  "sub_title": "Annual Recruitment Event",
  "summary": "Summary in English...",
  "content": "<p>Full HTML content</p>",
  "seo_title": "Job Fair 2023 | Smart Employment Platform",
  "seo_keywords": "job fair,recruitment,employment",
  "seo_description": "Annual recruitment event..."
}
```

响应：

```json
{
  "code": 0,
  "msg": "ok",
  "data": { "id": 5678 }
}
```

### 8.4 示例：前台获取文章详情

```http
GET /api/v2/article/get?article_id=1024&lang=en-US
```

响应：

```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "lang": "en-US",
    "fallback": false,
    "article_id": 1024,
    "title": "Smart Job Fair 2023",
    "sub_title": "Annual Recruitment Event",
    "summary": "Summary in English...",
    "content": "<p>Full HTML content</p>",
    "create_time": "2026-09-09T16:30:00+08:00"
  }
}
```

---

## 九、模板侧实现

### 9.1 模板中输出文案

Beego 模板支持 i18n 转义函数 `{{T "key"}}`：

```html
<!-- views/admin/Article/list.html -->
<!DOCTYPE html>
<html lang="{{.Lang}}">
<head>
  <title>{{T "Article.title"}} - {{T "Common.app_name"}}</title>
</head>
<body>
  <h1>{{T "Article.title"}}</h1>
  <button class="layui-btn">{{T "Common.add"}}</button>
  <button class="layui-btn layui-btn-primary">{{T "Common.search"}}</button>
  <table>
    <thead>
      <tr>
        <th>{{T "Article.title"}}</th>
        <th>{{T "Article.category"}}</th>
        <th>{{T "Article.status"}}</th>
      </tr>
    </thead>
  </table>
</body>
</html>
```

### 9.2 语言切换器组件

```html
<!-- views/admin/Common/lang_switch.html -->
<div class="layui-form lang-switch">
  <select lay-filter="lang">
    <option value="zh-CN" {{if eq .Lang "zh-CN"}}selected{{end}}>简体中文</option>
    <option value="zh-HK" {{if eq .Lang "zh-HK"}}selected{{end}}>繁體中文</option>
    <option value="en-US" {{if eq .Lang "en-US"}}selected{{end}}>English</option>
    <option value="ja-JP" {{if eq .Lang "ja-JP"}}selected{{end}}>日本語</option>
  </select>
</div>
<script>
layui.form.on("select(lang)", function(data){
  document.cookie = "lang=" + data.value + ";path=/;max-age=31536000";
  location.reload();
});
</script>
```

### 9.3 控制器注入 Lang

在基础控制器 `controllers/a.go` 中统一注入：

```go
// controllers/a.go
type BaseController struct {
  web.Controller
}

func (c *BaseController) Prepare() {
  // 从 context 取 Lang（由 LocaleFilter 写入）
  c.Data.Lang = c.GetString("lang")
  if c.Data.Lang == "" {
    c.Data.Lang = "zh-CN"
  }
  // 注入支持的语言列表（前端语言切换器使用）
  c.Data.SupportedLangs = []map[string]string{
    {"code": "zh-CN", "name": "简体中文"},
    {"code": "zh-HK", "name": "繁體中文"},
    {"code": "en-US", "name": "English"},
  }
}
```

---

## 十、前端（LayUI）集成

### 10.1 动态加载 JS 文案

通过全局注入 `window.Lang` 对象：

```js
// static/admin/js/i18n.js
window.Lang = {
  save: "{{T "Common.save"}}",
  cancel: "{{T "Common.cancel"}}",
  confirm: "{{T "Common.confirm"}}",
  delete: "{{T "Common.delete"}}",
  // ...
};
```

### 10.2 LayUI 组件国际化

LayUI 表格 / 表单 / 分页等组件文案统一通过 `window.Lang` 注入：

```js
layui.use(["table", "form", "laypage"], function(){
  var table = layui.table;
  table.render({
    elem: "#article-list",
    cols: [[
      {field: "title", title: window.Lang.title},
      {field: "category", title: window.Lang.category},
      {field: "status_text", title: window.Lang.status},
      {fixed: "right", title: window.Lang.action, toolbar: "#toolbar"}
    ]],
    page: {
      layout: [""prev", "page", "next", "skip", "count"],
      prev: window.Lang.prev_page,
      next: window.Lang.next_page
    }
  });
});
```

### 10.3 KindEditor / UEditor 富文本

富文本编辑器无需 i18n（编辑器内置），但加载的字体可按语言优化：

- zh-CN: 默认中文
- en-US: Arial / sans-serif
- ja-JP: MS Mincho
- ar-SA: RTL 布局（未来支持）

---

## 十一、运维与配置

### 11.1 conf/app.conf 新增配置

```ini
# ===== i18n 多语言配置 =====
# 默认语言（zh-CN / en-US / zh-HK / ja-JP）
i18n.default_lang = zh-CN
# 支持的语言（逗号分隔）
i18n.support_langs = zh-CN,zh-HK,en-US,ja-JP
# 内容缺失时是否 fallback 默认语言（true 推荐）
i18n.fallback_default = true
# 翻译缓存大小（条目数）
i18n.cache_size = 50000
# 翻译缓存过期（秒）
i18n.cache_ttl = 3600
# 启用浏览器语言协商（HTTP Accept-Language）
i18n.accept_language = true
# 启用 Cookie 持久化（30天）
i18n.cookie_ttl = 2592000
```

### 11.2 应用启动顺序（main.go 修改）

```go
func main() {
  // 1. 加载配置
  web.BConfig.RunMode = web.AppConfig.DefaultString("runmode", "dev")

  // 2. 初始化 i18n（早于路由注册）
  i18n.InitI18n("./app/i18n/conf")
  web.InsertFilter("/*", web.BeforeRouter, i18n.LocaleFilter())

  // 3. 数据库迁移
  if enabled, _ := web.AppConfig.Bool("cms.migrate"); enabled {
    db.RunAutoMigrate()
  }

  // 4. 启动 banner
  banner.Print()

  // 5. 启动 HTTP 服务
  web.Run()
}
```

### 11.3 监控指标

在 `controllers/metrics.go` 中扩展 Prometheus 指标：

```go
// 翻译命中率
i18nTranslationHits = promauto.NewCounterVec(prometheus.CounterOpts{
    Name: "cms_i18n_translation_hits_total",
    Help: "Total translation cache hits",
}, ["lang", "source_type"])

i18nTranslationMisses = promauto.NewCounterVec(prometheus.CounterOpts{
    Name: "cms_i18n_translation_misses_total",
    Help: "Total translation cache misses (DB fallback)",
}, ["lang", "source_type"])

i18nFallback = promauto.NewCounterVec(prometheus.CounterOpts{
    Name: "cms_i18n_fallback_total",
    Help: "Total times fallback to default lang",
}, ["lang"])
```

---

## 十二、回滚与灰度策略

### 12.1 灰度上线路径

| 阶段 | 内容 | 范围 | 验证指标 |
|------|------|------|----------|
| **第 1 阶段** | 仅 UI 文案多语言（conf/locale_*.ini） | 全量 | 模板渲染正确率 100% |
| **第 2 阶段** | + 通用翻译表 cms_translation | 内部租户 1 | DB 命中率、API 性能 |
| **第 3 阶段** | + 文章翻译表 cms_article_i18n | 内部租户 1+2 | 内容完整性、缓存命中 |
| **第 4 阶段** | + 全部业务翻译表 | 全量租户 | 内存占用、慢查询 |
| **第 5 阶段** | + 前台 v2 API 多语言版 | 前台灰度 | 浏览器兼容、移动端 |

### 12.2 Feature Flag

通过配置开关控制灰度：

```ini
# 仅对租户 ID 1 / 2 启用 i18n 业务翻译
i18n.tenant_whitelist = 1,2
# 仅对 URL 含 lang= 的请求启用多语言业务翻译
i18n.requires_lang_param = true
```

### 12.3 回滚方案

1. **配置回滚**：`i18n.tenant_whitelist=` 置空 → 所有租户回到无 i18n 状态
2. **数据回滚**：翻译表数据不回滚（独立表，不影响主表）
3. **代码回滚**：i18n 中间件代码可独立 Revert
4. **紧急开关**：在 LocaleFilter 顶部加 `if !i18n.Enabled { return }` 一键关闭

---

## 十三、风险与对策

| # | 风险 | 影响 | 对策 |
|---|------|------|------|
| 1 | 翻译表与主表数据不一致 | 内容缺失显示 | 兜底 fallback 策略 + 对账脚本 |
| 2 | 字符集不全（如繁体缺字） | 显示乱码 | 全库 utf8mb4 + InnoDB utf8mb4_unicode_ci |
| 3 | URL 中 lang 引发 SEO 重复内容 | 搜索权重分散 | 配置 `<link rel="alternate" hreflang>` |
| 4 | 大量翻译导致内存暴涨 | OOM | LRU 缓存（默认 50000 条 / 1小时 TTL）|
| 5 | 多语言内容审核复杂 | 内容安全 | 审核模块增加按 lang 维度过滤 |
| 6 | 管理员手工维护成本高 | 录入繁琐 | 集成机器翻译 API（百度 / DeepL / 腾讯）|
| 7 | 富文本编辑器语言 | 编辑器界面 | KindEditor 替换为 TinyMCE 多语言版 |
| 8 | DB 不同方言语法差异 | 跨方言脚本差异 | 通过 GORM AutoMigrate + migrations 多方言 |
| 9 | Excel 导入/导出乱码 | 业务中断 | 统一 CSV UTF-8 BOM |
| 10 | 时间日期格式差异 | 显示不一致 | 国际化 time.Time 格式化（Jan 1, 2026 / 2026年1月1日）|

### 13.1 自动机器翻译集成（可选）

`app/i18n/translate.go`：

```go
type Translator interface {
  Translate(text string, from, to string) (string, error)
}

type BaiduTranslator struct {
  AppID string
  Key   string
}

func (t *BaiduTranslator) Translate(text string, from, to string) (string, error) {
  // 调用百度翻译开放平台 API
  // 返回翻译结果
}
```

后台翻译编辑器可挂此模块，编辑人员输入原文后一键翻译为其他语言版本。

---

## 十四、里程碑与工作量估算

### 14.1 里程碑

| 里程碑 | 内容 | 周期 | 验收 |
|------|------|------|------|
| **M1** | UI 文案国际化（locale ini + Beego i18n） | 1 周 | 后台 / 前台所有页面可切换语言 |
| **M2** | cms_translation 通用表 + 通用 service | 0.5 周 | 任意字段可翻译 |
| **M3** | cms_article_i18n + Article 服务层多语言 | 1 周 | 文章读写翻译正常 |
| **M4** | 栏目 / 标签 / 专题 / 广告 / 链接翻译表 | 1 周 | 各业务 CRUD 支持翻译 |
| **M5** | 站点 / 主题 / 租户配置翻译 | 0.5 周 | 多租户独立语言 |
| **M6** | v2 API 多语言版 + 前台灰度 | 1 周 | API 兼容旧版 + 新版多语言 |
| **M7** | 集成机器翻译 + 后台翻译辅助 UI | 1 周 | 一键翻译 |
| **M8** | 监控 / 告警 / 灰度开关 / 文档 | 0.5 周 | 全量上线 |

**总计约 6.5 周**（3 名 Go 后端 + 1 名前端 + 1 名 QA）。

### 14.2 优先级

1. **P0**：M1 + M3 + M4（后台编辑人员能用，文章可翻译）
2. **P1**：M2 + M5（通用翻译 + 多租户语言）
3. **P2**：M6 + M7 + M8（前台 API + 机器翻译 + 监控）

### 14.3 验收清单

- [ ] 后台切换中英文，UI 文案正确
- [ ] 后台保存文章翻译，前台 en-US 取到英文内容
- [ ] 前台 zh-CN 取到中文（即使没有翻译）
- [ ] URL `/en-US/article/list` 路径前缀可用
- [ ] Accept-Language 协商工作
- [ ] Cookie 持久化语言（30 天）
- [ ] 翻译缺失自动 fallback 默认语言，响应中 `fallback=true`
- [ ] Prometheus 指标可观测
- [ ] 多方言数据库迁移无报错（MySQL / PG / MSSQL / Kingbase / GaussDB）
- [ ] 灰度开关可控制上线范围
- [ ] 紧急回滚一键关闭

---

## 附录 A：与芋道 yudao / MyBatis-Plus 多语言方案对照

| 维度 | 芋道 yudao | 本 CMS 方案 |
|------|------|------|
| 框架 | Spring Boot | Beego v2.3.8 |
| ORM | MyBatis-Plus | GORM + gorm.io/gen |
| 翻译存储 | 表行级 + 通用 I18n 表 | 同上（方案对齐）|
| 多租户 | yudao.base.TenantBaseDO | `tenant_id` 列 + `app/cms/tenant/` Scope |
| 软删除 | yudao.base.BaseDO（deleted BIT）| `deleted tinyint(1)`（已对齐）|
| UI 国际化 | Spring MessageSource | Beego i18n + ini 文件 |
| 内容国际化 | DictTranslService | `cms_translation` + `cms_xxx_i18n` |

## 附录 B：参考链接

- DeepSeek 讨论：https://chat.deepseek.com/share/lukvbu834bxzfk0t05
- Beego i18n 文档：https://beego.vip/docs/mvc/controller/
- gorm.io/gen：https://gorm.io/gen/
- 项目迁移脚本：`cms/migrations/README.md`
- 自动迁移机制：`cms/app/db/migrate.go`

