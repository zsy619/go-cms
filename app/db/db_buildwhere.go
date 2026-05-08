package db

import (
	"errors"
	"reflect"
	"strings"

	"github.com/zsy619/tools/xdatabase"
	"gorm.io/gorm"
)

// BuildWhere 构建动态查询条件
// 支持结构体、Map和Slice三种形式的查询条件
func BuildWhere(db *gorm.DB, where any) (*gorm.DB, error) {
	t := reflect.TypeOf(where).Kind()
	switch t {
	case reflect.Struct, reflect.Map:
		db = db.Where(where)
	case reflect.Slice:
		for _, item := range where.([]any) {
			item := item.([]any)
			column := item[0]
			if reflect.TypeOf(column).Kind() == reflect.String {
				count := len(item)
				if count == 1 {
					return nil, errors.New("切片长度不能小于2")
				}
				columnstr := column.(string)
				if strings.Contains(columnstr, "?") {
					db = db.Where(column, item[1:]...)
				} else {
					cond := "and"
					var val any
					var opt string
					if count == 2 {
						opt = "="
						val = item[1]
					} else {
						opt = strings.ToLower(item[1].(string))
						val = item[2]
						if count == 4 {
							cond = strings.ToLower(strings.ReplaceAll(item[3].(string), " ", ""))
						}
					}
					column = columnstr + " " + opt + " ?"

					if strings.Contains(" in notin ", opt) {
						column = columnstr + " " + opt + " (?)"
					} else if strings.Contains(" = < > <= >= <> != <=> like likebinary notlike ilike rlike regexp notregexp", opt) {
						column = columnstr + " " + opt + " ?"
					}

					if cond == "and" {
						db = db.Where(column, val)
					} else {
						db = db.Or(column, val)
					}
				}
			} else if t == reflect.Map {
				db = db.Where(item)
			} else {
				var err error
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

// BuildQueryList 构建分页查询
// wheres: 查询条件, columns: 选择列, orderBy: 排序, page: 页码, rows: 每页数量
func BuildQueryList(db *gorm.DB, wheres any, columns any, orderBy any, page, rows int) (*gorm.DB, error) {
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
	return db, nil
}

// Paginate 分页函数
// 返回一个gorm回调函数用于分页
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
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
