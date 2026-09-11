# 多语言功能变更日志

## [1.0.0] - 2024

### 新增功能

#### 后台管理
- 语言设置页面（站点管理 → 语言设置）
- 支持 197 种语言（含 emoji 国旗图标）
- 语言列表（layui table，不分页）
- 支持搜索（按名称、代码）
- 语言新增、编辑、删除
- 排序、状态切换
- 默认语言设置
- 一键初始化语言数据

#### 数据模型
- 新增 cms_language 表
- cms_site 添加 language_code 字段
- cms_site_channel 添加 language_code 字段

#### API 接口

##### 后台管理 API
- GET /admin/site/language - 列表页面
- GET /admin/site/language/data - 列表数据
- GET /admin/site/language/edit - 编辑页面
- POST /admin/site/language/save - 保存
- GET /admin/site/language/delete - 删除
- GET /admin/site/language/sort - 排序
- GET /admin/site/language/status - 状态
- GET /admin/site/language/init - 初始化

##### 前台 API
- GET /api/site/language/list - 站点语言关联列表
- GET /api/site/language/find - 按站点ID查询
- GET /api/site/language/by-code - 按语言代码查询

#### 服务层
- CmsLanguage 服务（CRUD、搜索、排序）
- SiteLanguageService 服务（关联查询）

#### 数据库
- utf8mb4 字符集支持
- GORM AutoMigrate 自动迁移

#### 文档
- 后台使用文档
- API 接口文档
- 字段说明文档
- 快速开始指南
- 前端集成指南

### 技术栈

- Go + Beego v2
- GORM v2
- MySQL 5.7+ / 8.0+
- Layui v2.9.21

### 文件变更

#### 新增文件
1. app/cms/domain/cms_language.gen.go
2. app/cms/domain/site_language.go
3. app/cms/service/cms_language.go
4. app/cms/service/site_language.go
5. controllers/admin/language.go
6. views/admin/Site/Language.html
7. views/admin/Site/LanguageEdit.html
8. sql/init_language.sql
9. views/help/dev/lang_admin.html
10. views/help/dev/api_site_language.html
11. views/help/dev/site_language_field.html
12. views/help/dev/lang_quickstart.html
13. views/help/dev/lang_frontend.html
14. docs/LANGUAGE.md
15. docs/CHANGELOG_LANGUAGE.md

#### 修改文件
1. app/db/cms_models.go - 注册 CmsLanguage 模型
2. app/db/consts.go - 字符集改为 utf8mb4
3. app/cms/domain/cms_site.gen.go - 添加 LanguageCode 字段
4. app/cms/domain/cms_site_channel.gen.go - 添加 LanguageCode 字段
5. controllers/admin/aa.go - 注册 language 路由
6. controllers/admin/site.go - SiteEdit 加载语言列表
7. controllers/www/api_site.go - 添加 API 接口
8. controllers/www/a.go - 添加 API 路由
9. views/admin/Site/SiteEdit.html - 添加语言选择
10. views/admin/Site/ChannelEdit.html - 添加语言选择
