package dal

import (
	"net/url"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// 数据库连接池配置常量
const (
	MaxIdleConns    = 10  // 最大空闲连接数
	MaxOpenConns    = 100 // 最大打开连接数
	ConnMaxLifetime  = 3600 // 连接最大生命周期(秒)
	ConnMaxIdleTime = 600  // 连接最大空闲时间(秒)
)

var (
	CmsDatabase  *DBExtension // CMS主数据库
	AuthDatabase *DBExtension // 认证数据库
)

// init 初始化数据库连接
// 说明: 分别连接CMS主数据库和认证数据库,配置连接池参数以优化性能和资源利用
func init() {
	// CMS数据库初始化
	{
		dsn := GetDbConn("cms", "cms")
		dbConnect, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true, // 使用单数表名
			},
		})
		if err != nil {
			logs.Error("CMS数据库连接失败: %v", err)
			panic(err)
		}

		// 配置数据库连接池
		sqlDB, err := dbConnect.DB()
		if err != nil {
			logs.Error("获取CMS数据库连接池失败: %v", err)
			panic(err)
		}
		sqlDB.SetMaxIdleConns(MaxIdleConns)
		sqlDB.SetMaxOpenConns(MaxOpenConns)
		sqlDB.SetConnMaxLifetime(time.Duration(ConnMaxLifetime) * time.Second)
		sqlDB.SetConnMaxIdleTime(time.Duration(ConnMaxIdleTime) * time.Second)

		CmsDatabase = NewDBWrapper(dbConnect)
		logs.Info("CMS数据库连接成功, 连接池配置: MaxIdle=%d, MaxOpen=%d", MaxIdleConns, MaxOpenConns)
	}

	// 认证数据库初始化
	{
		dsn := GetDbConn("auth", "un2co_yunzhipin")
		dbConnect, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})
		if err != nil {
			logs.Warning("认证数据库连接失败: %v", err)
		} else {
			// 配置数据库连接池
			sqlDB, err := dbConnect.DB()
			if err != nil {
				logs.Warning("获取认证数据库连接池失败: %v", err)
			} else {
				sqlDB.SetMaxIdleConns(MaxIdleConns)
				sqlDB.SetMaxOpenConns(MaxOpenConns)
				sqlDB.SetConnMaxLifetime(time.Duration(ConnMaxLifetime) * time.Second)
				sqlDB.SetConnMaxIdleTime(time.Duration(ConnMaxIdleTime) * time.Second)
			}

			AuthDatabase = NewDBWrapper(dbConnect)
			logs.Info("认证数据库连接成功")
		}
	}

	// 自动迁移数据库表结构
	runAutoMigrate()
}

// GetDbConn 构建数据库连接字符串(DSN)
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
	dsn := dbuser + ":" + dbpassword + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname + "?charset=utf8&parseTime=true"
	if timezone != "" {
		dsn = dsn + "&loc=" + url.QueryEscape(timezone)
	}
	logs.Debug("数据库连接配置[%s]: %s", prefix, dsn)
	return dsn
}