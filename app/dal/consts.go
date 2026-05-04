package dal

import (
	"fmt"
	"net/url"

	"github.com/beego/beego/v2/server/web"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	CmsDatabase  *DBExtension // gsg数据库
	AuthDatabase *DBExtension // 认证数据库
)

func init() {
	// cms数据库
	{
		dsn := GetDbConn("cms", "cms")
		dbConnect, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})
		if err != nil {
			fmt.Println(err.Error())
			panic(err)
		}

		CmsDatabase = NewDBWrapper(dbConnect)
	}
	// 认证数据库
	{
		dsn := GetDbConn("auth", "un2co_yunzhipin")
		dbConnect, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})
		if err != nil {
			fmt.Println(err.Error())
			// panic(err)
		} else {
			AuthDatabase = NewDBWrapper(dbConnect)
		}
	}
}

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
	fmt.Println(prefix, dsn)
	return dsn
}
