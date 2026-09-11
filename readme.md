# CMS 内容管理系统

基于 **Go + Beego + GORM** 的多租户 CMS 内容管理系统，提供文章/广告/站点/链接/标签/微信等核心模块，支持 MySQL/PostgreSQL/SQL Server/人大金仓/高斯等多种数据库方言。

---

## 目录

- [项目简介](#项目简介)
- [技术栈](#技术栈)
- [核心特性](#核心特性)
- [目录结构](#目录结构)
- [三层架构（domain / mapper / service）](#三层架构domain--mapper--service)
- [快速开始](#快速开始)
- [数据库配置（多方言）](#数据库配置多方言)
- [GORM 自动迁移与元数据同步](#gorm-自动迁移与元数据同步)
- [核心功能模块](#核心功能模块)
- [微信公众号集成](#微信公众号集成)
- [多租户与软删除](#多租户与软删除)
- [视图与模板](#视图与模板)
- [部署说明](#部署说明)
- [常见问题](#常见问题)

---

## 项目简介

本项目 是一个通用的 CMS 内容管理系统，基于 Beego + GORM 构建，提供：

- **多租户隔离**：通过 `tenant_id` 字段实现业务数据隔离
- **软删除**：通过 `deleted` 字段（`tinyint(1)`）实现数据软删除，与 MyBatis-Plus / 芋道 yudao 设计对齐
- **多数据库方言**：同一套代码可运行于 MySQL / PostgreSQL / SQL Server / 人大金仓 / 高斯数据库
- **自动元数据同步**：应用启动时自动建表、补齐列注释、表注释、索引，并修复历史部署中的 `bit(1)` 兼容性问题
- **微信公众号集成**：支持公众号配置、自定义菜单、消息规则、被动回复、网页授权等完整能力
- **就业业务插件**：空中宣讲会、招聘会、岗位、企业、线上报名等业务模块

---

## 技术栈

| 层级       | 技术                                        | 版本                     | 用途                             |
| ---------- | ------------------------------------------- | ------------------------ | -------------------------------- |
| 语言       | Go                                          | 1.27                     | 程序语言                         |
| Web 框架   | Beego                                       | v2.3.8                   | MVC 框架 / 路由 / 配置 / Session |
| ORM        | GORM                                        | v1.30.0                  | 数据库访问                       |
| 代码生成   | gorm.io/gen                                 | v0.3.27                  | 生成 domain/mapper               |
| 数据库驱动 | GORM drivers                                | mysql/postgres/sqlserver | 多方言数据库支持                 |
| 监控       | Prometheus                                  | v1.22.0                  | 指标采集（`/metrics`）           |
| 模板引擎   | Beego Template                              | -                        | HTML 视图渲染                    |
| 前端 UI    | LayUI + KindEditor + ECharts + UEditor Plus | -                        | 后台界面 / 富文本 / 图表         |
| 热重载     | Air                                         | -                        | 开发模式热部署                   |

---

## 核心特性

- **严格三层架构**：`domain`（领域实体）+ `mapper`（数据访问）+ `service`（业务编排）一一对应，39 张业务表 / 722 个字段 0 偏差
- **零配置建表**：应用启动时自动 `CREATE TABLE IF NOT EXISTS`，无需手动执行 DDL
- **列注释自动同步**：从 struct tag `comment:xxx` 自动同步到所有方言（MySQL `ALTER ... COMMENT`、PG `COMMENT ON COLUMN`、MSSQL `sp_addextendedproperty`）
- **索引自动创建**：从 struct tag `index/uniqueIndex` 自动建索引，幂等执行
- **bit(1) 自动修复**：检测到历史部署的 `deleted bit(1)` 列时自动转为 `tinyint(1)`，修复 Go 驱动 Scan 失败
- **多租户支持**：`tenant_id` 自动过滤，支持 GORM Scope 用法
- **CAS 单点登录**：可配置启用 CAS 2.0 协议
- **双角色登录**：管理员 / 学校双角色统一入口
- **微信公众平台**：完整的公众号接入能力
- **服务管理**：支持 Windows service / Linux systemd 注册

---

## 目录结构

```
cms/
├── app/                          # 业务应用层
│   ├── banner/                   # 启动 Banner 渲染（Spring Boot 风格）
│   ├── cms/                      # CMS 核心业务
│   │   ├── domain/               # 领域实体（*.gen.go 自动生成）
│   │   ├── mapper/               # 数据访问层（*.gen.go 自动生成 + a.go 访问器）
│   │   ├── service/              # 服务层（手工编写）
│   │   └── tenant/               # 多租户实现（Scope/Context/Reflect）
│   ├── db/                       # 数据库基础设施
│   │   ├── consts.go             # 数据库方言注册 / DSN 构建 / SQL 钩子
│   │   ├── db_extension.go       # GORM DB 扩展（GetList/SaveOne/Count 等）
│   │   ├── db_buildwhere.go      # 动态 WHERE 构建器
│   │   ├── cms_models.go         # 39 个 domain 实体的 AutoMigrate 注册
│   │   ├── migrate.go            # 自动迁移编排器（元数据提取/列注释/索引/类型修复）
│   │   └── logger.go             # SQL 日志
│   ├── tool/                     # 工具层（config/cache/const）
│   ├── wechat/                   # 微信集成（mp/pay/models/utils/simulate）
│   └── models/                   # 公共数据模型
├── controllers/                  # 控制器层
│   ├── admin/                    # 后台管理（25+ 个控制器）
│   │   ├── login.go              # 登录认证（CAS + 账号密码）
│   │   ├── article.go            # 文章管理
│   │   ├── ads.go / link.go / site.go / theme.go / tag.go
│   │   ├── notice.go / menu.go / nav.go / index.go / admin.go
│   │   ├── file.go / cache.go / common*.go / weixin*.go
│   │   └── vmodel/               # 视图模型（请求/响应 DTO）
│   ├── www/                      # 前台展示（API 接口）
│   │   ├── api_*.go              # 30+ 个 REST API
│   │   ├── article.go / channel.go / tag.go / topic.go
│   │   ├── sse.go                # Server-Sent Events
│   │   └── wechat_*.go           # 微信回调
│   ├── plugin/                   # 业务插件
│   │   ├── airKeynote.go         # 空中宣讲会
│   │   ├── jobFair.go            # 招聘会
│   │   ├── job.go                # 就业岗位
│   │   ├── xsbm.go               # 线上报名
│   │   └── company.go            # 企业管理
│   ├── funcs/                    # 模板函数（format/time/string/url）
│   ├── a.go, b.go                # 控制器基类
│   ├── captcha.go                # 验证码
│   ├── filter_admin.go                # 过滤器
│   ├── metrics.go                # Prometheus 指标
│   └── error.go                  # 错误处理
├── views/                        # Beego 视图模板
│   └── admin/                    # 后台 HTML 模板（19+ 模块）
│       ├── layout/               # 布局
│       ├── login/                # 登录页（含动态背景动画）
│       ├── Article/ / Ads/ / Site/ / Tag/ / Theme/ / Topic/
│       └── Weixin/ / ...
├── conf/                          # 配置
│   ├── app.conf                  # 应用主配置
│   └── app_remote.conf           # 远程配置
├── migrations/                   # 数据库迁移脚本（多方言 DDL）
│   ├── 01_mysql.sql              # MySQL/MariaDB（39 表 + bit(1)→tinyint(1) 修复段）
│   ├── 02_postgres.sql           # PostgreSQL / openGauss / KingbaseES
│   ├── 03_sqlserver.sql          # SQL Server
│   ├── 04_weixin_init.sql        # 微信模块 6 表建表脚本
│   └── README.md                 # 迁移说明文档
├── static/                       # 静态资源（CSS/JS/Images/Plugins）
├── main.go                       # 程序入口
├── go.mod / go.sum               # Go 模块依赖
├── .air.toml                     # Air 热重载配置
├── build.sh / run.sh / update.sh # 构建/启动/更新脚本
└── readme.md                     # 本文档
```

---

## 三层架构（domain / mapper / service）

### 39 个业务实体一览

```
cms_ad_category_relation, cms_admin, cms_admin_log, cms_admin_nav,
cms_admin_notice, cms_admin_role, cms_admin_role_site, cms_admin_role_value,
cms_ads, cms_ads_category, cms_ads_category_relation, cms_album,
cms_article, cms_article_category, cms_article_category_relation,
cms_article_comment, cms_article_label, cms_article_label_relation,
cms_article_property, cms_attach, cms_link, cms_link_category,
cms_link_category_relation, cms_site, cms_site_channel,
cms_site_channel_album, cms_site_channel_field, cms_site_domain,
cms_tag, cms_tenant, cms_theme, cms_topic,
plg_online_register,
weixin_account, weixin_menu, weixin_mp_verify,
weixin_request_content, weixin_request_rule, weixin_response_content
```

### 典型调用示例

```go
// service 层调用模式
mdl, do := mapper.CmsArticleDo()

// 多条件查询
articles, total, err := do.Where(
    mdl.Deleted.Is(false),
    mdl.Status.Eq(2),
    mdl.SiteID.Eq(siteId),
).Order(mdl.SortID.Desc()).FindByPage((page-1)*limit, limit)

// 软删除（标记 deleted=true）
_, err = do.Where(mdl.ArticleID.Eq(articleId)).UpdateColumns(
    map[string]interface{}{
        mdl.Deleted.ColumnName().String(): true,
        mdl.UpdateTime.ColumnName().String(): time.Now(),
    },
)
```

### 关键设计

| 字段        | 类型         | NOT NULL | DEFAULT | 用途       |
| ----------- | ------------ | -------- | ------- | ---------- |
| `tenant_id` | `bigint`     | ✅        | `0`     | 多租户隔离 |
| `deleted`   | `tinyint(1)` | ✅        | `0`     | 软删除标记 |

> **重要**：`deleted` 列必须为 `tinyint(1)` 而非 `bit(1)`，否则 Go MySQL 驱动会返回 `[]byte`，无法 Scan 到 `bool`，导致 `Scan error: couldn't convert "\\x00" into type bool` 错误。

---

## 快速开始

### 环境要求

- Go **1.27+**
- MySQL 5.7+ / PostgreSQL 12+ / SQL Server 2017+（任一）
- Git

### 安装步骤

```bash
# 1. 进入项目目录
cd cms

# 2. 安装依赖
go mod tidy

# 3. 配置数据库（编辑 conf/app.conf）
#    默认端口 8125

# 4. 启动服务
go run main.go

# 或使用 air 热重载（开发推荐）
air

# 5. 浏览器访问
#    前台首页: http://localhost:8125
#    后台登录: http://localhost:8125/cms/admin/login
```

### 默认登录

后台路径：/cms/admin/login

默认管理员账号通过 SQL 注入或 cms_tenant 默认租户初始化时创建。

---

## 数据库配置（多方言）

```ini
# conf/app.conf

# ===== CMS 主数据库 =====
cms.dialect = mysql        # mysql | postgres | sqlserver | kingbase | gaussdb
cms.host = 127.0.0.1
cms.user = root
cms.password = 123456
cms.port = 3306
cms.db = cms
cms.prefix = cms_           # 表名前缀（实际表名为 cms_xxx）
cms.timezone = Asia/Shanghai
cms.migrate = true          # 启动时自动迁移（建表+注释+索引+类型修复）

# ===== Auth 认证数据库（微信公众号等独立库） =====
auth.dialect = mysql
auth.host = 127.0.0.1
auth.user = root
auth.password = 123456
auth.port = 3306
auth.db = un2co_yunzhipin
auth.prefix = un2co_
auth.timezone = Asia/Shanghai

# ===== Web 服务 =====
httpport = 8125
runmode = dev                # dev | prod
sessionprovider = "file"
sessionproviderconfig = "./logs/session"
LOCAL_DOMAIN = http://localhost:8125
```

### 支持的方言

| 方言            | 常量      | 默认端口 | 说明                 |
| --------------- | --------- | -------- | -------------------- |
| MySQL / MariaDB | mysql     | 3306     | 默认方言             |
| PostgreSQL      | postgres  | 5432     | 含 openGauss         |
| SQL Server      | sqlserver | 1433     | TDS 加密默认 disable |
| 人大金仓        | kingbase  | 54321    | KingbaseES，PG 协议  |
| 高斯数据库      | gaussdb   | 5432     | openGauss，PG 协议   |

### 双数据库支持

cms.* 用于业务主库，auth.* 用于认证/微信库。两库可为不同方言甚至不同主机，迁移时自动按分组应用。

---

## GORM 自动迁移与元数据同步

应用启动时（cms.migrate=true）自动执行：

### 1. GORM AutoMigrate

CREATE TABLE IF NOT EXISTS + ADD COLUMN IF NOT EXISTS — 39 张业务表全部覆盖。

### 2. 列注释同步

从 struct gorm tag comment:xxx 提取注释，按方言生成 SQL：

| 方言       | SQL                                                |
| ---------- | -------------------------------------------------- |
| MySQL      | ALTER TABLE `t` MODIFY COLUMN `c` TYPE COMMENT xxx |
| PostgreSQL | COMMENT ON COLUMN t.c IS xxx                       |
| SQL Server | EXEC sys.sp_addextendedproperty ...                |

### 3. 索引自动创建

从 struct gorm tag index / uniqueIndex 提取，命名规则：idx_<table>_<column>。

### 4. 表注释

通过 buildCMSTableComments() / buildAuthTableComments() 集中维护，例如：

```go
"cms_article":    "文章表",
"weixin_account": "微信公众号表",
```

### 5. 类型修复

FixBit1ToTinyInt1() 检测 deleted 列若为 bit(1) 自动转为 tinyint(1)，解决 Go 驱动 Scan 失败。

### 关键 API

```go
// 公共 SQL 构建器（跨方言）
GetTableCommentSQL(dialect, table, comment)
GetColumnCommentSQL(dialect, table, column, ctype, comment)
GetAddIndexSQL(dialect, table, idxName, cols, unique)
GetAddColumnSQL(dialect, table, column, ctype, notNull, def)
GetModifyColumnSQL(dialect, table, column, newType)
GetCheckColumnSQL(dialect, table, column)

// 业务操作（容错）
EnsureColumn(db, dialect, table, column, ctype, notNull, def, comment) (added bool, err)
EnsureIndex(db, dialect, table, idxName, cols, unique) (added bool, err)
SetTableComment(db, dialect, table, comment) error
SetColumnComment(db, dialect, table, column, ctype, comment) error
FixBit1ToTinyInt1(db, dialect, table) error
```

---

## 核心功能模块

### CMS 基础模块

| 模块     | 控制器          | 视图                 | 说明                             |
| -------- | --------------- | -------------------- | -------------------------------- |
| 文章管理 | article.go      | views/admin/Article/ | 文章 CRUD / 审核 / 回收站 / 评论 |
| 栏目分类 | article.go      | 同上                 | 树形栏目结构 + 模板配置          |
| 专题管理 | article.go      | -                    | 专题聚合                         |
| 标签管理 | tag.go          | views/admin/Tag/     | 标签 CRUD                        |
| 主题管理 | theme.go        | views/admin/Theme/   | 模板主题                         |
| 链接管理 | link.go         | views/admin/Link/    | 快捷链接                         |
| 广告管理 | ads.go          | views/admin/Ads/     | 广告位 / 分类 / 内容             |
| 站点管理 | site.go         | views/admin/Site/    | 多站点 / 频道 / 域名             |
| 相册管理 | common_album.go | -                    | 图片附件                         |
| 公告管理 | notice.go       | views/admin/Notice/  | 后台公告                         |
| 菜单管理 | menu.go         | -                    | 后台导航菜单                     |
| 缓存管理 | cache.go        | views/admin/Cache/   | Redis / 内存缓存                 |
| 文件管理 | file.go         | views/admin/File/    | 资源上传 / 浏览                  |

### 业务插件模块

| 模块       | 控制器               | 说明           |
| ---------- | -------------------- | -------------- |
| 空中宣讲会 | plugin/airKeynote.go | 视频宣讲直播   |
| 招聘会     | plugin/jobFair.go    | 线下招聘会管理 |
| 就业岗位   | plugin/job.go        | 岗位发布       |
| 线上报名   | plugin/xsbm.go       | 用户报名       |
| 企业管理   | plugin/company.go    | 企业入驻       |

### 前台 API

| 模块     | 控制器                         | 说明                   |
| -------- | ------------------------------ | ---------------------- |
| 文章 API | www/api_article.go             | 文章列表 / 详情 / 搜索 |
| 广告 API | www/api_ads.go                 | 广告位渲染             |
| 链接 API | www/api_link.go                | 友情链接               |
| 标签 API | www/api_tag.go                 | 标签云                 |
| 专题 API | www/api_topic.go               | 专题页                 |
| 站点 API | www/api_site.go                | 站点配置               |
| 报名 API | www/api_plg_online_register.go | 线上报名               |
| 实时推送 | www/sse.go                     | Server-Sent Events     |

---

## 微信公众号集成

### 模块组成

```
app/wechat/
├── mp/            # 公众号（授权/JSSDK/消息加解密/被动回复）
├── pay/           # 微信支付 V2/V3
├── models/        # 数据模型
├── simulate/      # 接口模拟测试
└── utils/         # 工具函数
```

### 6 张微信业务表

| 表名                    | 用途                                        |
| ----------------------- | ------------------------------------------- |
| weixin_account          | 公众号配置（AppID/AppSecret/Token/AES Key） |
| weixin_menu             | 自定义菜单（含 parent_id/account_id 外键）  |
| weixin_mp_verify        | 公众号验证文件（MP_verify_xxx.txt）         |
| weixin_request_content  | 关键字触发的回复内容                        |
| weixin_request_rule     | 关键字规则（默认/模糊/响应类型）            |
| weixin_response_content | 消息记录日志（请求/响应/xml）               |

### 控制器

- controllers/admin/weixin_account.go - 公众号管理后台
- controllers/admin/weixin_content.go - 规则/回复内容
- controllers/admin/weixin_menu.go - 自定义菜单
- controllers/admin/weixin_mp_verify.go - 验证文件
- controllers/admin/weixin_request.go - 请求规则
- controllers/www/wechat_mp.go - 公众号消息入口
- controllers/www/wechat_mp_web_auth.go - 网页授权
- controllers/www/wechat_mp_verify.go - 公众号验证回调

### 建表初始化

通过 migrations/04_weixin_init.sql 或 GORM AutoMigrate 自动创建（应用启动时）。

---

## 多租户与软删除

### 设计参考

设计参考芋道 yudao 的 TenantBaseDO：

- 软删除：单一字段 deleted（与 MyBatis-Plus @TableLogic 一致）
- 多租户：单一字段 tenant_id
- 不再使用 is_deleted

### 字段规范

每张业务表自动包含：

| 列名      | MySQL                         | PostgreSQL                | SQL Server                |
| --------- | ----------------------------- | ------------------------- | ------------------------- |
| tenant_id | bigint NOT NULL DEFAULT 0     | bigint NOT NULL DEFAULT 0 | BIGINT NOT NULL DEFAULT 0 |
| deleted   | tinyint(1) NOT NULL DEFAULT 0 | boolean DEFAULT false     | BIT NOT NULL DEFAULT 0    |

### 多租户实现

app/cms/tenant/ 提供 GORM Scope 用法，详见 app/cms/tenant/README.md。

### 历史数据迁移

migrations/01_mysql.sql 中包含 is_deleted → deleted 数据迁移段：

```sql
-- 数据迁移: 把 is_deleted 的值复制到 deleted
UPDATE `cms_admin` SET `deleted` = 1 WHERE `is_deleted` = 1;

-- 删除旧的 is_deleted 字段
ALTER TABLE `cms_admin` DROP COLUMN `is_deleted`;
```

末尾还有 bit(1) → tinyint(1) 修复段，自动执行 39 张表的列类型变更。

---

## 视图与模板

### Beego 模板语法

- 模板文件：*.html
- 模板标签：{{...}}、{{< partial >}}、{{< block >}} 等
- 自定义函数：controllers/funcs/ 提供 format / time / string / url 等

### 登录页背景动画

views/admin/login/login.html 内置三层动画：

1. **6 色动画渐变**（60fps 流畅过渡）
2. **鼠标光晕跟随**（CSS 变量 --mx/--my + radial-gradient）
3. **5 个浮动圆球** + **80 个粒子**（4 种形状）+ 连接线 + 点击波纹

### 模板布局

views/admin/layout/ 提供主布局：

- <aside> 侧边栏（菜单/导航）
- <nav> 顶部栏（用户信息/工具栏）
- <section> 内容区

---

## 部署说明

### 开发环境

```bash
# 方式1: 直接运行
go run main.go

# 方式2: Air 热重载（推荐）
air
```

### 生产部署

#### 1. 编译二进制

```bash
CGO_ENABLED=0 go build -ldflags="-s -w" -o cms main.go
```

#### 2. 上传到服务器

```bash
scp cms user@server:/opt/cms/
```

#### 3. 服务注册（Linux systemd）

/etc/systemd/system/cms.service：

```ini
[Unit]
Description=CMS Service
After=network.target mysql.service

[Service]
Type=simple
User=cms
WorkingDirectory=/opt/cms
ExecStart=/opt/cms/cms
Restart=always
RestartSec=5
LimitNOFILE=65535

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl daemon-reload
sudo systemctl enable cms
sudo systemctl start cms
sudo systemctl status cms
```

#### 4. 内置服务管理

```bash
# Windows
cms install
cms uninstall

# Linux
sudo ./cms install
sudo ./cms uninstall
```

#### 5. Nginx 反向代理

```nginx
upstream cms {
    server 127.0.0.1:8125;
}

server {
    listen 80;
    server_name cms.example.com;

    location / {
        proxy_pass http://cms;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_read_timeout 300s;
    }

    location /static/ {
        alias /opt/cms/static/;
        expires 30d;
    }

    location /Uploads/ {
        alias /opt/cms/Uploads/;
        expires 7d;
    }
}
```

#### 6. 监控指标

Prometheus 通过 /metrics 端点采集：

```bash
curl http://localhost:8125/metrics
```

---

## 数据库迁移

详见 migrations/README.md。

### 核心流程

1. 启用 cms.migrate=true，应用启动时自动迁移（推荐开发/测试环境）
2. 生产环境执行 migrations/01_mysql.sql（含 tenant_id/deleted 添加 + is_deleted 迁移 + bit(1) 修复段）
3. 跨方言切换：修改 cms.dialect 即可，无需改业务代码

### 关键迁移脚本

| 文件               | 适用                            |
| ------------------ | ------------------------------- |
| 01_mysql.sql       | MySQL/MariaDB                   |
| 02_postgres.sql    | PostgreSQL/openGauss/KingbaseES |
| 03_sqlserver.sql   | SQL Server                      |
| 04_weixin_init.sql | 微信模块 6 表                   |

---

## 开发提示

### 代码生成

domain/mapper 由 gorm.io/gen 自动生成，配置文件见 app/db/cms_models.go。如需重新生成：

```bash
# 在项目根目录执行 gen 生成命令
go run gen.go
```

### 添加新表的步骤

1. 在数据库创建表（带 tenant_id + deleted）
2. 在 app/cms/domain/ 添加 struct 定义（带 gorm tag）
3. 在 app/cms/mapper/ 添加 mapper（参考其他 *.gen.go）
4. 在 app/db/cms_models.go 的 buildCMSModels() 添加 &domain.Xxx{}
5. 在 app/db/cms_models.go 的 buildCMSTableComments() 添加注释
6. 在 app/cms/service/ 编写业务方法
7. 重启应用，AutoMigrate 自动建表 + 同步注释

### 字段命名规范

| 字段                    | 类型          | 用途                                        |
| ----------------------- | ------------- | ------------------------------------------- |
| xxx_id                  | bigint        | 主键（auto_increment）                      |
| create_id / create_name | int / varchar | 创建人                                      |
| create_time             | datetime      | 创建时间（默认 CURRENT_TIMESTAMP）          |
| update_id / update_name | int / varchar | 更新人                                      |
| update_time             | datetime      | 修改时间                                    |
| sort_id                 | int           | 排序字段                                    |
| status                  | tinyint       | 状态（0草稿 1提交 2审核通过 3未通过 4驳回） |
| tenant_id               | bigint        | 多租户隔离                                  |
| deleted                 | tinyint(1)    | 软删除标记                                  |

---

## 常见问题

### Q: 登录报 Scan error on column ... name "deleted": couldn't convert "\x00" into type bool？

A: 数据库 deleted 列是 bit(1) 类型，Go MySQL 驱动无法 Scan 为 bool。

解决方案：
1. 手动执行 migrations/01_mysql.sql 末尾的修复段（39 张表 ALTER TABLE）
2. 或启用 cms.migrate=true，应用启动时 FixBit1ToTinyInt1() 自动修复
3. 长期方案：domain gorm tag 统一为 type:tinyint(1)（本项目已修复）

### Q: 如何启用 CAS 单点登录？

A: 在 conf/app.conf 中配置：

```ini
cas.enabled = true
cas.url = http://cas-server:10240
```

然后调用 /cms/admin/login 即可跳转 CAS 登录。

### Q: 多租户如何开启？

A: 应用已默认启用 tenant_id 字段。需要在请求处理时通过 app/cms/tenant/tenant_context.go 设置当前租户 ID。

### Q: 如何新增方言（如达梦 DM8 / Oracle）？

A: 在 app/db/consts.go 的 init() 中调用 RegisterDialect(name, dbDialectConfig{...})，提供 opener/buildDSN/tableComment/columnComment/addIndex/addColumn/modifyColumn/checkColumn 七个 SQL 钩子。

### Q: Session 存储在哪里？

A: 默认 ./logs/session，可在 conf/app.conf 修改 sessionproviderconfig。

### Q: 如何关闭自动建表（避免 GORM 锁表）？

A: 设置 cms.migrate = false，手动执行 migrations/01_mysql.sql。

### Q: 微信菜单修改后未生效？

A: 自定义菜单需通过 weixin_menu 表保存后，调用微信公众号接口同步（控制器：controllers/admin/weixin_menu.go）。

### Q: 应用启动很慢？

A: 检查：
1. cms.migrate=true 时首次启动会建 39 张表，耗时 5-10 秒
2. 数据库索引是否存在缺失
3. auth.* 配置错误会导致回退到 CMS 库，日志中会有 AuthDatabase 不可用 警告

### Q: 表注释没有同步？

A: 检查：
1. app/db/cms_models.go 中 buildCMSTableComments() / buildAuthTableComments() 是否包含目标表
2. 数据库用户是否有 ALTER TABLE 权限
3. 方言是否支持（PG/MS 需要 superuser 或 schema owner 权限）

---

## 项目里程碑

| 日期       | 提交    | 说明                                                                     |
| ---------- | ------- | ------------------------------------------------------------------------ |
| 2026-09-09 | d7f1410 | 一键修订 domain/mapper/service 三层映射关系（39 实体 / 722 字段 0 偏差） |
| 2026-09-09 | 295c77f | 实现 GORM 自动建表 + 列注释 + 索引 + 类型修复                            |
| 2026-09-09 | 2437fb0 | 添加多租户 + 软删除 + weixin 模块迁移脚本                                |
| 2026-09-09 | 115bd6f | 登录页背景动画 + 居中布局                                                |
| 2026-09-09 | a524bcb | 修订 Deleted 字段引用 + bit(1) 修复                                      |
| 2026-09-09 | 5282a7c | Spring Boot 风格启动 Banner                                              |
| 2026-09-09 | 27a93e7 | 升级 go 1.27 + 新增 sqlserver/postgres 驱动                              |

---

## 许可证

本项目可用于商业用途。

---

## 联系方式

如有技术问题，请联系项目维护团队。
