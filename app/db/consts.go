package db

import (
	"fmt"
	"net/url"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// 数据库连接池配置常量
const (
	MaxIdleConns    = 10   // 最大空闲连接数
	MaxOpenConns    = 100  // 最大打开连接数
	ConnMaxLifetime = 3600 // 连接最大生命周期(秒)
	ConnMaxIdleTime = 600  // 连接最大空闲时间(秒)
)

// 数据库方言常量 — 支持主流及国产数据库
const (
	DialectMySQL     = "mysql"     // MySQL / MariaDB
	DialectPostgres  = "postgres"  // PostgreSQL
	DialectSQLServer = "sqlserver" // Microsoft SQL Server
	DialectOracle    = "oracle"    // Oracle
	DialectDameng    = "dameng"    // 达梦数据库(DM8)
	DialectKingbase  = "kingbase"  // 人大金仓(KingbaseES, PG兼容)
	DialectGaussDB   = "gaussdb"   // 高斯数据库(openGauss, PG兼容)
)

var (
	CmsDatabase  *DBExtension // CMS主数据库
	AuthDatabase *DBExtension // 认证数据库
	cmsDialect   string       // CMS数据库方言
	authDialect  string       // 认证数据库方言
)

// dsnOptions DSN构建所需参数
type dsnOptions struct {
	host     string // 主机地址
	port     string // 端口
	user     string // 用户名
	password string // 密码
	dbname   string // 数据库名
	timezone string // 时区
	sslmode  string // SSL模式(PG/Oracle/国产库)
	prefix   string // 配置前缀,用于读取额外参数
}

// dbDialectConfig 每种数据库方言的配置信息
// 通过 RegisterDialect 函数注册,支持运行时扩展
type dbDialectConfig struct {
	opener       func(string) gorm.Dialector        // GORM驱动开启器(dsn→Dialector)
	buildDSN     func(dsnOptions) string            // DSN构建函数
	tableComment func(table, comment string) string // 表注释SQL语句生成函数
	defaultPort  string                             // 默认端口
}

// dialectRegistry 全局方言注册表
var dialectRegistry = map[string]dbDialectConfig{}

// RegisterDialect 注册一个数据库方言
// 在 init() 中调用,将驱动工厂函数、DSN构建器、表注释生成器关联到方言名
func RegisterDialect(name string, cfg dbDialectConfig) {
	dialectRegistry[name] = cfg
}

// init 注册内置驱动并初始化数据库连接
func init() {
	// 注册 MySQL 驱动
	RegisterDialect(DialectMySQL, dbDialectConfig{
		opener: mysql.Open,
		buildDSN: func(o dsnOptions) string {
			dsn := o.user + ":" + o.password + "@tcp(" + o.host + ":" + o.port + ")/" + o.dbname + "?charset=utf8&parseTime=true"
			if o.timezone != "" {
				dsn = dsn + "&loc=" + url.QueryEscape(o.timezone)
			}
			return dsn
		},
		tableComment: func(table, comment string) string {
			return fmt.Sprintf("ALTER TABLE %s COMMENT = '%s'", table, comment)
		},
		defaultPort: "3306",
	})

	// 注册 PostgreSQL 驱动
	RegisterDialect(DialectPostgres, dbDialectConfig{
		opener: postgres.Open,
		buildDSN: func(o dsnOptions) string {
			if o.sslmode == "" {
				o.sslmode = "disable"
			}
			return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
				o.host, o.user, o.password, o.dbname, o.port, o.sslmode, o.timezone)
		},
		tableComment: func(table, comment string) string {
			return fmt.Sprintf("COMMENT ON TABLE %s IS '%s'", table, comment)
		},
		defaultPort: "5432",
	})

	// 注册 SQL Server 驱动
	RegisterDialect(DialectSQLServer, dbDialectConfig{
		opener: sqlserver.Open,
		buildDSN: func(o dsnOptions) string {
			// 格式: sqlserver://user:password@host:port?database=dbname&encrypt=disable
			dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s",
				url.QueryEscape(o.user), url.QueryEscape(o.password),
				o.host, o.port, o.dbname)
			encrypt, _ := web.AppConfig.String(o.prefix + ".encrypt")
			if encrypt == "" {
				encrypt = "disable"
			}
			dsn = dsn + "&encrypt=" + encrypt
			return dsn
		},
		tableComment: func(table, comment string) string {
			// SQL Server 使用扩展属性存储表注释,dbo是默认schema
			return fmt.Sprintf(`EXEC sys.sp_addextendedproperty @name=N'MS_Description', @value=N'%s', @level0type=N'SCHEMA', @level0name=N'dbo', @level1type=N'TABLE', @level1name=N'%s'`,
				comment, table)
		},
		defaultPort: "1433",
	})

	// 注册 人大金仓 — PG兼容,复用PostgreSQL驱动
	RegisterDialect(DialectKingbase, dbDialectConfig{
		opener: postgres.Open,
		buildDSN: func(o dsnOptions) string {
			if o.sslmode == "" {
				o.sslmode = "disable"
			}
			return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
				o.host, o.user, o.password, o.dbname, o.port, o.sslmode, o.timezone)
		},
		tableComment: func(table, comment string) string {
			return fmt.Sprintf("COMMENT ON TABLE %s IS '%s'", table, comment)
		},
		defaultPort: "54321", // Kingbase默认端口
	})

	// 注册 高斯数据库 — PG兼容,复用PostgreSQL驱动
	RegisterDialect(DialectGaussDB, dbDialectConfig{
		opener: postgres.Open,
		buildDSN: func(o dsnOptions) string {
			if o.sslmode == "" {
				o.sslmode = "disable"
			}
			return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
				o.host, o.user, o.password, o.dbname, o.port, o.sslmode, o.timezone)
		},
		tableComment: func(table, comment string) string {
			return fmt.Sprintf("COMMENT ON TABLE %s IS '%s'", table, comment)
		},
		defaultPort: "5432",
	})

	// CMS数据库初始化
	{
		dsn := GetDbConn("cms", "cms")
		cmsDialect = getDialect("cms")
		dbConnect, err := openDatabase(cmsDialect, dsn)
		if err != nil {
			logs.Error("CMS数据库连接失败[%s]: %v", cmsDialect, err)
			panic(err)
		}

		configureConnectionPool(dbConnect, "CMS")
		CmsDatabase = NewDBWrapper(dbConnect)
		CmsDatabase.dialect = cmsDialect
		logs.Info("CMS数据库[%s]连接成功, 连接池: MaxIdle=%d, MaxOpen=%d", cmsDialect, MaxIdleConns, MaxOpenConns)
	}

	// 认证数据库初始化
	{
		dsn := GetDbConn("auth", "un2co_yunzhipin")
		authDialect = getDialect("auth")
		dbConnect, err := openDatabase(authDialect, dsn)
		if err != nil {
			logs.Warning("认证数据库连接失败[%s]: %v", authDialect, err)
		} else {
			configureConnectionPool(dbConnect, "认证")
			AuthDatabase = NewDBWrapper(dbConnect)
			AuthDatabase.dialect = authDialect
			logs.Info("认证数据库[%s]连接成功", authDialect)
		}
	}

	// 自动迁移数据库表结构
	runAutoMigrate()
}

// getDialect 获取数据库方言配置
// @param prefix 配置前缀,如"cms"将读取cms.dialect配置
// @return 数据库方言,默认返回mysql
func getDialect(prefix string) string {
	if prefix == "" {
		prefix = "db"
	}
	dialect, _ := web.AppConfig.String(prefix + ".dialect")
	if dialect == "" {
		dialect = DialectMySQL
	}
	return dialect
}

// getDialectConfig 获取方言配置,支持PostgreSQL兼容别名自动映射
// Kingbase/GaussDB 返回 postgres 的配置(它们共用PG驱动)
func getDialectConfig(dialect string) (dbDialectConfig, bool) {
	cfg, ok := dialectRegistry[dialect]
	return cfg, ok
}

// openDatabase 根据数据库方言打开对应数据库连接
// @param dialect 数据库方言
// @param dsn 数据库连接字符串
// @return *gorm.DB 数据库连接实例
func openDatabase(dialect string, dsn string) (*gorm.DB, error) {
	config := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
	}

	cfg, ok := getDialectConfig(dialect)
	if !ok {
		supportList := registeredDialects()
		return nil, fmt.Errorf("不支持的数据库方言: %s, 当前支持: %v", dialect, supportList)
	}

	return gorm.Open(cfg.opener(dsn), config)
}

// registeredDialects 返回所有已注册的方言列表
func registeredDialects() []string {
	list := make([]string, 0, len(dialectRegistry))
	for k := range dialectRegistry {
		list = append(list, k)
	}
	return list
}

// GetTableCommentSQL 根据数据库方言生成表注释SQL语句
func GetTableCommentSQL(dialect, tableName, comment string) string {
	cfg, ok := getDialectConfig(dialect)
	if !ok {
		cfg, ok = getDialectConfig(DialectMySQL)
		if !ok {
			return fmt.Sprintf("ALTER TABLE %s COMMENT = '%s'", tableName, comment)
		}
	}
	return cfg.tableComment(tableName, comment)
}

// configureConnectionPool 配置数据库连接池参数
func configureConnectionPool(db *gorm.DB, dbName string) {
	sqlDB, err := db.DB()
	if err != nil {
		logs.Error("%s数据库连接池配置失败: %v", dbName, err)
		panic(err)
	}
	sqlDB.SetMaxIdleConns(MaxIdleConns)
	sqlDB.SetMaxOpenConns(MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(ConnMaxLifetime) * time.Second)
	sqlDB.SetConnMaxIdleTime(time.Duration(ConnMaxIdleTime) * time.Second)
}

// GetDbConn 构建数据库连接字符串(DSN)
// 根据数据库中言自动生成对应格式的连接串:
//   - MySQL:       user:password@tcp(host:port)/dbname?charset=utf8&parseTime=true&loc=...
//   - PostgreSQL:  host=host user=user password=pwd dbname=db port=port sslmode=disable
//   - SQL Server:  sqlserver://user:password@host:port?database=db
//   - Kingbase:    host=host user=user ... (与PostgreSQL相同)
//   - GaussDB:     host=host user=user ... (与PG相同,默认端口5432)
//
// Oracle 和 达梦 需分别引入对应驱动包,DSN格式:
//   - Oracle:      oracle://user:password@host:port/service_name
//   - 达梦:        dm://user:password@host:port?schema=SCHEMA_NAME
//
// @param prefix 配置前缀,如"cms"将读取cms.host,cms.port等配置
// @param defaultDb 默认数据库名称
// @return string 数据库DSN连接字符串
func GetDbConn(prefix string, defaultDb string) string {
	if prefix == "" {
		prefix = "db"
	}
	dbhost, _ := web.AppConfig.String(prefix + ".host")
	if dbhost == "" {
		dbhost = "127.0.0.1"
	}
	dbuser, _ := web.AppConfig.String(prefix + ".user")
	if dbuser == "" {
		dbuser = "root"
	}
	dbpassword, _ := web.AppConfig.String(prefix + ".password")
	if dbpassword == "" {
		dbpassword = "123456"
	}
	dbname, _ := web.AppConfig.String(prefix + ".db")
	if dbname == "" {
		dbname = defaultDb
	}
	timezone, _ := web.AppConfig.String(prefix + ".timezone")
	if timezone == "" {
		timezone = "Asia/Shanghai"
	}
	sslmode, _ := web.AppConfig.String(prefix + ".sslmode")

	dialect := getDialect(prefix)
	cfg, ok := getDialectConfig(dialect)

	// 如果方言未注册,尝试使用默认解析
	dbport, _ := web.AppConfig.String(prefix + ".port")
	if !ok {
		logs.Warn("方言 %s 未注册驱动, 使用配置中的连接参数", dialect)
		if dbport == "" {
			dbport = "3306"
		}
		dsn := dbuser + ":" + dbpassword + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname
		logs.Debug("数据库连接配置[%s/%s]: %s", prefix, dialect, dsn)
		return dsn
	}

	// 使用方言默认端口
	if dbport == "" {
		dbport = cfg.defaultPort
	}

	opts := dsnOptions{
		host:     dbhost,
		port:     dbport,
		user:     dbuser,
		password: dbpassword,
		dbname:   dbname,
		timezone: timezone,
		sslmode:  sslmode,
		prefix:   prefix,
	}

	dsn := cfg.buildDSN(opts)
	logs.Debug("数据库连接配置[%s/%s/%s]: %s", prefix, dialect, dbport, dsn)
	return dsn
}

// GetDialectDefaultPort 获取方言的默认端口
func GetDialectDefaultPort(dialect string) string {
	cfg, ok := getDialectConfig(dialect)
	if ok {
		return cfg.defaultPort
	}
	return "3306"
}
