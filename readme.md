# 河南省大中专学生智慧就业平台 - CMS内容管理系统

## 项目简介

本项目是**河南省大中专学生智慧就业平台**的CMS内容管理系统，负责对接教育部API，提供就业信息管理、招聘会管理、空中宣讲会等功能模块。系统支持多站点管理、内容发布、审核流程、微信公众平台集成等核心功能。

### 项目背景

为贯彻落实国家关于做好大中专毕业生就业工作的决策部署，河南省建立了智慧就业平台，通过信息化手段提升就业服务质量。本CMS系统作为平台的核心组成部分，承担着信息发布、内容管理、数据对接等关键职责。

---

## 技术栈

### 后端技术

| 技术 | 版本 | 说明 |
|------|------|------|
| **Go** | 1.25+ | 程序语言 |
| **Beego** | v2.3.8 | Web框架 |
| **GORM** | v1.30.0 | ORM框架 |
| **GORM Gen** | v0.3.27 | 代码生成工具 |
| **MySQL** | 5.7+ | 数据库 |
| **Prometheus** | v1.22.0 | 监控指标 |

### 前端技术

| 技术 | 版本 | 说明 |
|------|------|------|
| **LayUI** | 2.5.x~2.9.x | 前端UI框架 |
| **ECharts** | - | 数据可视化图表 |
| **KindEditor** | v4.1.11 | 富文本编辑器 |
| **UEditor Plus** | v3.0.0 | 百度编辑器增强版 |

---

## 目录结构

```
cms/
├── app/                        # 业务应用层
│   ├── cms/                   # CMS核心业务
│   │   ├── domain/           # 领域模型(自动生成)
│   │   ├── mapper/           # 数据访问层(自动生成)
│   │   └── service/          # 服务层
│   ├── dal/                  # 数据访问层
│   │   ├── db_extension.go   # 数据库扩展
│   │   ├── db_buildwhere.go  # 查询构建器
│   │   └── consts.go        # 常量定义
│   ├── tool/                 # 工具层
│   │   ├── config.go        # 配置管理
│   │   ├── cache.go         # 缓存管理
│   │   └── const.go         # 公共常量
│   └── wechat/              # 微信集成模块
│       ├── mp/             # 公众号相关(授权、JSSDK、消息)
│       ├── pay/            # 微信支付
│       ├── models/         # 数据模型
│       └── utils/          # 工具函数
├── controllers/              # 控制器层
│   ├── admin/              # 后台管理控制器
│   │   ├── vmodel/        # 视图模型
│   │   ├── ads.go         # 广告管理
│   │   ├── article.go     # 文章管理
│   │   ├── link.go        # 链接管理
│   │   ├── weixin*.go     # 微信管理系列
│   │   └── login.go       # 登录认证
│   ├── www/               # 前台展示控制器
│   │   ├── api_*.go       # API接口
│   │   ├── article.go     # 文章展示
│   │   └── wechat_*.go    # 微信相关
│   ├── plugin/           # 插件控制器(就业业务)
│   │   ├── airKeynote.go # 空中宣讲会
│   │   ├── jobFair.go    # 招聘会
│   │   ├── job.go        # 就业岗位
│   │   ├── xsbm.go       # 线上报名
│   │   └── company.go    # 企业管理
│   └── funcs/            # 模板函数
├── static/                # 静态资源
│   ├── admin/           # 后台静态资源
│   │   ├── css/        # 样式文件
│   │   ├── js/         # 脚本文件
│   │   ├── lib/        # 第三方库(LayUI、KindEditor等)
│   │   └── images/     # 图片资源
│   ├── theme/          # 前台主题
│   └── plugins/        # 前端插件
├── views/                 # 视图模板
│   └── admin/          # 后台HTML模板
│       ├── Index/      # 首页模板
│       ├── Article/    # 文章模板
│       ├── Ads/        # 广告模板
│       └── ...
├── conf/                 # 配置文件
│   ├── app.conf        # 应用配置
│   └── app_remote.conf # 远程配置
├── main.go              # 程序入口
├── go.mod              # Go模块依赖
└── go.sum              # 依赖校验
```

---

## 快速开始

### 环境要求

- Go 1.16+
- MySQL 5.7+
- Git

### 安装步骤

1. **克隆项目**
   ```bash
   cd /Volumes/E/JYW/河南省大中专学生智慧就业平台/教育部API对接
   ```

2. **安装依赖**
   ```bash
   cd cms
   go mod tidy
   ```

3. **配置数据库**
   编辑 `conf/app.conf` 文件，修改数据库连接信息：
   ```ini
   cms.host = 127.0.0.1
   cms.port = 3306
   cms.user = root
   cms.password = 123456
   cms.db = cms
   ```

4. **启动服务**
   ```bash
   # 开发模式
   go run main.go
   
   # 或使用air热重载
   air
   ```

5. **访问系统**
   打开浏览器访问：`http://localhost:8125`
   后台管理：`http://localhost:8125/cms/admin/login`

### 服务管理

```bash
# 安装为系统服务(Windows)
go run main.go install

# 卸载服务
go run main.go uninstall

# Linux系统服务
sudo ./cms install
```

---

## 核心功能模块

### 1. 内容管理

- **文章管理** - 内容发布、编辑、审核、回收站
- **栏目管理** - 树形栏目结构、栏目模板配置
- **专题管理** - 内容专题聚合
- **标签管理** - 内容标签系统

### 2. 广告管理

- **广告位管理** - 创建广告位、设置展示位置
- **广告投放** - 广告创建、排序、时段控制
- **数据统计** - 点击统计、曝光统计

### 3. 就业业务插件

| 模块 | 说明 | 访问路径 |
|------|------|----------|
| **空中宣讲会** | 企业视频宣讲、直播回放 | /plugin/airKeynote/* |
| **招聘会** | 双选会、专场招聘会管理 | /plugin/jobFair/* |
| **就业岗位** | 职位发布、简历投递 | /plugin/job/* |
| **线上报名** | 活动在线报名 | /plugin/xsbm/* |
| **企业管理** | 企业入驻、资质审核 | /plugin/company/* |

### 4. 微信管理

- **公众号配置** - AppID、AppSecret配置
- **自定义菜单** - 菜单创建、个性化配置
- **自动回复** - 关键词回复、关注回复
- **网页授权** - OAuth2用户授权
- **模板消息** - 消息推送

### 5. 用户权限

- **管理员管理** - 超级管理员、普通管理员
- **学校账户** - 高校用户管理
- **角色权限** - 基于角色的权限控制(RBAC)

---

## API接口

### 文章接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/article/get` | 获取文章列表 |
| GET | `/api/article/get/new` | 获取最新文章 |
| GET | `/api/article/paginate` | 分页获取文章 |
| GET | `/api/article/find` | 获取文章详情 |
| GET | `/api/article/prev_next` | 获取上一篇/下一篇 |
| GET | `/api/article/click` | 文章点击+1 |
| GET | `/api/article/like` | 文章点赞+1 |

### 栏目接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/category/nav` | 获取栏目导航 |
| GET | `/api/category/get` | 获取栏目列表 |
| GET | `/api/category/find` | 获取栏目详情 |

### 广告接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/ads/get` | 获取广告列表 |
| GET | `/api/ads/find` | 获取广告详情 |

### 链接接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/link/get` | 获取链接列表 |
| GET | `/api/link/find` | 获取链接详情 |

---

## 配置说明

### 应用配置 (conf/app.conf)

```ini
# 应用基本信息
appname = cms
httpport = 8125
runmode = dev

# 会话配置
sessionon = true
sessionprovider = "file"
sessionproviderconfig = "./logs/session"

# CMS数据库配置
cms.host = 127.0.0.1
cms.port = 3306
cms.user = root
cms.password = 123456
cms.db = cms
cms.prefix = cms_
cms.timezone = Asia/Shanghai

# 认证数据库配置
auth.host = 127.0.0.1
auth.port = 3306
auth.db = un2co_yunzhipin

# CAS单点登录(可选)
cas.enabled = false
cas.url = http://localhost:10240

# 本地域名
LOCAL_DOMAIN = http://localhost:8125
```

### 静态资源路径

系统会自动映射以下路径：

| 路径 | 映射目录 | 说明 |
|------|----------|------|
| `/views` | views/ | 模板目录 |
| `/images` | static/images | 图片目录 |
| `/css` | static/css | 样式目录 |
| `/js` | static/js | 脚本目录 |
| `/Upload` | Upload/ | 上传文件 |
| `/Public` | Public/ | 公共资源 |

---

## 数据库设计

### 数据表前缀

- `cms_ads_*` - 广告相关表
- `cms_article_*` - 文章内容表
- `cms_link_*` - 链接管理表
- `cms_notice_*` - 通知公告表
- `cms_channel_*` - 频道/栏目表
- `cms_admin_*` - 管理员表
- `cms_site_*` - 站点配置表
- `plg_*` - 插件相关表

### 状态值说明

| 状态值 | 说明 |
|--------|------|
| 0 | 草稿/正常 |
| 1 | 已提交 |
| 2 | 审核通过 |
| 3 | 审核未通过 |
| 4 | 驳回 |

---

## 开发指南

### 代码生成

使用GORM Gen自动生成代码：

```bash
cd cmd
go run generate.go
```

### 新增功能模块

1. 在 `app/cms/domain/` 创建领域模型
2. 在 `app/cms/mapper/` 生成数据访问层
3. 在 `app/cms/service/` 编写业务逻辑
4. 在 `controllers/` 创建控制器
5. 在 `views/admin/` 创建模板文件

### 模板语法

```html
<!-- 模板继承 -->
{{template "layout.html" .}}

<!-- 模板渲染 -->
{{.VariableName}}

<!-- 循环遍历 -->
{{range .List}}
    <li>{{.Name}}</li>
{{end}}
```

---

## 部署说明

### 开发环境

```bash
go run main.go
```

### 生产环境

1. 编译二进制文件：
   ```bash
   go build -o cms main.go
   ```

2. 配置Nginx反向代理（可选）：
   ```nginx
   location / {
       proxy_pass http://127.0.0.1:8125;
   }
   ```

3. 使用systemd管理服务（Linux）：
   ```ini
   [Unit]
   Description=CMS Service
   
   [Service]
   ExecStart=/path/to/cms
   Restart=always
   
   [Install]
   WantedBy=multi-user.target
   ```

---

## 常见问题

### Q: 如何开启CAS单点登录？
A: 在 `conf/app.conf` 中设置 `cas.enabled = true` 并配置 `cas.url`。

### Q: 如何修改数据库连接？
A: 编辑 `conf/app.conf` 中的 `cms.*` 配置项。

### Q: 如何启用微信支付？
A: 在 `app/tool/config.go` 中配置微信支付相关参数。

### Q: Session存储在哪里？
A: 默认存储在 `./logs/session` 目录，支持修改为Redis存储。

---

## 许可证

本项目仅供河南省大中专学生智慧就业平台内部使用。

---

## 联系方式

如有技术问题，请联系项目维护团队。