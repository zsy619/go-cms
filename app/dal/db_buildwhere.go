package dal

import (
	"errors"
	"reflect"
	"strings"

	"github.com/zsy619/tools/xdatabase"
	"gorm.io/gorm"
)

// https://github.com/qicmsg/go_vcard/blob/master/app/models/entity/Gorm.go

///region 业务上扩展方法

func BuildWhere(db *gorm.DB, where interface{}) (*gorm.DB, error) {
	var err error
	t := reflect.TypeOf(where).Kind()
	switch t {
	case reflect.Struct, reflect.Map:
		db = db.Where(where)
	case reflect.Slice:
		for _, item := range where.([]interface{}) {
			item := item.([]interface{})
			column := item[0]
			if reflect.TypeOf(column).Kind() == reflect.String {
				count := len(item)
				if count == 1 {
					return nil, errors.New("切片长度不能小于2")
				}
				columnstr := column.(string)
				// 拼接参数形式
				if strings.Contains(columnstr, "?") {
					db = db.Where(column, item[1:]...)
				} else {
					cond := "and" // cond
					opt := "="
					_opt := " = "
					var val interface{}
					if count == 2 {
						opt = "="
						val = item[1]
					} else {
						opt = strings.ToLower(item[1].(string))
						_opt = " " + strings.ReplaceAll(opt, " ", "") + " "
						val = item[2]
					}

					if count == 4 {
						cond = strings.ToLower(strings.ReplaceAll(item[3].(string), " ", ""))
					}

					/*
					   '=', '<', '>', '<=', '>=', '<>', '!=', '<=>',
					   'like', 'like binary', 'not like', 'ilike',
					   '&', '|', '^', '<<', '>>',
					   'rlike', 'regexp', 'not regexp',
					   '~', '~*', '!~', '!~*', 'similar to',
					   'not similar to', 'not ilike', '~~*', '!~~*',
					*/

					if strings.Contains(" in notin ", _opt) {
						// val 是数组类型
						column = columnstr + " " + opt + " (?)"
					} else if strings.Contains(" = < > <= >= <> != <=> like likebinary notlike ilike rlike regexp notregexp", _opt) {
						column = columnstr + " " + opt + " ?"
					}

					if cond == "and" {
						db = db.Where(column, val)
					} else {
						db = db.Or(column, val)
					}
				}
			} else if t == reflect.Map /*Map*/ {
				db = db.Where(item)
			} else {
				/*
					// 解决and 与 or 混合查询，但这种写法有问题，会抛出 invalid query condition
					db = db.Where(func(db *gorm.DB) *gorm.DB {
						db, err = BuildWhere(db, item)
						if err != nil {
							panic(err)
						}
						return db
					})*/

				db, err = BuildWhere(db, item)
				if err != nil {
					return nil, err
				}
			}
		}
	default:
		return nil, errors.New("参数有误")
	}
	return db, nil
}

func BuildQueryList(db *gorm.DB, wheres interface{}, columns interface{}, orderBy interface{}, page, rows int) (*gorm.DB, error) {
	var err error
	db, err = BuildWhere(db, wheres)
	if err != nil {
		return nil, err
	}
	db = db.Select(columns)
	if orderBy != nil && orderBy != "" {
		db = db.Order(orderBy)
	}
	if page > 0 && rows > 0 {
		db = db.Limit(rows).Offset((page - 1) * rows)
	}
	return db, err
}

///endregion

// 分页封装
// https://gobea.cn/blog/detail/Xg5pqmoa.html
func Paginate(page int, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}
		offset := xdatabase.Offset(page, pageSize)
		return db.Offset(offset).Limit(pageSize)
	}
}
