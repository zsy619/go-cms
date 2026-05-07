package dal

import (
	"fmt"
	"net/url"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// 数据库连接池配置常量
const (
	MaxIdleConns    = 10                        // 最大空闲连接数
	MaxOpenConns    = 100                       // 最大打开连接数
	ConnMaxLifetime = 3600                      // 连接最大生命周期(秒)
	ConnMaxIdleTime = 600                       // 连接最大空闲时间(秒)
	DialectMySQL    = "mysql"                   // MySQL 数据库方言
	DialectPostgres = "postgres"                // PostgreSQL 数据库方言
)

var (
	CmsDatabase  *DBExtension // CMS主数据库
	AuthDatabase *DBExtension // 认证数据库
	cmsDialect   string       // CMS数据库方言(mysql/postgres)
	authDialect  string       // 认证数据库方言(mysql/postgres)
)

// init 初始化数据库连接
// 根据配置中的cms.dialect和auth.dialect自动选择MySQL或PostgreSQL驱动
func init() {
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
// @return 数据库方言(mysql/postgres),默认返回mysql
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

// openDatabase 根据数据库方言打开对应数据库连接
// @param dialect 数据库方言(mysql/postgres)
// @param dsn 数据库连接字符串
// @return *gorm.DB 数据库连接实例
func openDatabase(dialect string, dsn string) (*gorm.DB, error) {
	config := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
	}

	switch dialect {
	case DialectPostgres:
		return gorm.Open(postgres.Open(dsn), config)
	case DialectMySQL:
		return gorm.Open(mysql.Open(dsn), config)
	default:
		return nil, fmt.Errorf("不支持的数据库方言: %s, 仅支持 mysql 或 postgres", dialect)
	}
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
//   - MySQL: user:password@tcp(host:port)/dbname?charset=utf8&parseTime=true&loc=...
//   - PostgreSQL: host=host user=user password=password dbname=dbname port=port sslmode=disable TimeZone=...
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
	dbport, _ := web.AppConfig.String(prefix + ".port")
	if dbport == "" {
		dbport = "3306"
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
	dialect := getDialect(prefix)

	var dsn string
	switch dialect {
	case DialectPostgres:
		sslmode, _ := web.AppConfig.String(prefix + ".sslmode")
		if sslmode == "" {
			sslmode = "disable"
		}
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
			dbhost, dbuser, dbpassword, dbname, dbport, sslmode, timezone)
	case DialectMySQL:
		dsn = dbuser + ":" + dbpassword + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname + "?charset=utf8&parseTime=true"
		if timezone != "" {
			dsn = dsn + "&loc=" + url.QueryEscape(timezone)
		}
	}
	logs.Debug("数据库连接配置[%s/%s]: %s", prefix, dialect, dsn)
	return dsn
}
