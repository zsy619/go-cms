# 多语言功能文档

CMS 系统支持完整的多语言功能，包括语言管理、站点语言关联、API 查询等。

## 文档目录

### 后台管理
- 语言设置 (lang_admin.html) - 语言管理后台使用文档
- 站点语言字段 (site_language_field.html) - cms_site/cms_site_channel 字段说明
- 快速开始 (lang_quickstart.html) - 5 分钟快速集成

### API 接口
- 站点语言 API (api_site_language.html) - 多语言查询接口

### 前端集成
- 前端集成 (lang_frontend.html) - 前端多语言切换实现

## 核心特性

- 197 种语言 - 内置全球主流语言（含 emoji 国旗）
- 多维度搜索 - 按名称、代码筛选
- CRUD 完整 - 支持新增、编辑、删除、排序、状态切换
- 一键初始化 - 自动导入全部语言数据
- 三级关联 - cms_language → cms_site → cms_site_channel
- 完整 API - 列表、按站点、按语言代码查询
- emoji 支持 - 完整 utf8mb4 字符集支持

## 数据模型

### cms_language（语言设置表）

| 字段 | 类型 | 说明 |
|------|------|------|
| language_id | bigint | 主键 |
| name | varchar(64) | 语言名称（本民族语言） |
| code | varchar(16) | 语言代码（ISO 639-1） |
| icon | varchar(64) | emoji 图标 |
| sort_id | int | 排序 |
| description | varchar(256) | 中文描述 |
| status | tinyint | 状态（0禁用/1启用） |
| is_default | tinyint(1) | 是否默认 |

### cms_site（站点表新增字段）

| 字段 | 类型 | 说明 |
|------|------|------|
| language_code | varchar(16) | 语言代码 |

### cms_site_channel（频道表新增字段）

| 字段 | 类型 | 说明 |
|------|------|------|
| language_code | varchar(16) | 语言代码 |

## API 接口

### 后台管理接口

| URL | 方法 | 说明 |
|-----|------|------|
| /admin/site/language | GET | 语言列表页面 |
| /admin/site/language/data | GET | 获取语言列表数据 |
| /admin/site/language/edit | GET | 编辑页面 |
| /admin/site/language/save | POST | 保存语言 |
| /admin/site/language/delete | GET | 删除语言 |
| /admin/site/language/sort | GET | 更新排序 |
| /admin/site/language/status | GET | 切换状态 |
| /admin/site/language/init | GET | 初始化数据 |

### 前台 API 接口

| URL | 方法 | 说明 |
|-----|------|------|
| /api/site/language/list | GET | 站点语言关联列表 |
| /api/site/language/find | GET | 根据站点ID查询 |
| /api/site/language/by-code | GET | 根据语言代码查询 |

## 文件结构

```
cms/
├── app/
│   └── cms/
│       ├── domain/
│       │   ├── cms_language.gen.go    # 语言表模型
│       │   ├── cms_site.gen.go        # 站点表（含language_code）
│       │   ├── cms_site_channel.gen.go # 频道表（含language_code）
│       │   └── site_language.go        # 站点语言视图模型
│       └── service/
│           ├── cms_language.go         # 语言服务
│           └── site_language.go        # 站点语言服务
├── controllers/
│   ├── admin/
│   │   └── language.go                 # 后台语言控制器
│   └── www/
│       └── api_site.go                 # 前台 API
├── views/
│   ├── admin/Site/
│   │   ├── Language.html               # 语言列表
│   │   ├── LanguageEdit.html           # 语言编辑
│   │   ├── SiteEdit.html               # 站点编辑（含语言选择）
│   │   └── ChannelEdit.html            # 频道编辑（含语言选择）
│   └── help/dev/
│       ├── lang_admin.html             # 后台文档
│       ├── api_site_language.html      # API 文档
│       ├── site_language_field.html    # 字段文档
│       ├── lang_quickstart.html        # 快速开始
│       └── lang_frontend.html          # 前端集成
└── sql/
    └── init_language.sql               # 初始化 SQL
```

## 快速开始

### 1. 初始化数据

```bash
# 方式一：API 初始化
curl http://localhost:8125/admin/site/language/init

# 方式二：SQL 初始化
mysql -u root -p cms < sql/init_language.sql

# 方式三：重启服务（GORM AutoMigrate）
```

### 2. 设置站点语言

访问 站点管理 → 站点列表 → 编辑，在"站点语言"下拉框中选择语言。

### 3. 设置频道语言

访问 站点管理 → 频道列表 → 编辑，在"频道语言"下拉框中选择语言。

### 4. 调用 API

```bash
curl http://localhost:8125/api/site/language/list
```

## 注意事项

- 数据库必须使用 utf8mb4 字符集以支持 emoji
- 修改字符集后需要重启服务
- GORM AutoMigrate 会自动添加 language_code 字段
- 旧数据库需要手动执行 ALTER TABLE

## 技术支持

- GitHub: https://github.com/zsy619/go-cms
- 文档目录: views/help/dev/
