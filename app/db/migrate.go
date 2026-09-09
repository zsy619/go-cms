package db

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/beego/beego/v2/core/logs"
)

// ============================================================
// GORM 自动迁移:表结构 + 表注释 + 列注释 + 索引 + 类型修正
// ============================================================

// ColumnMeta 通过反射提取的列元数据
type ColumnMeta struct {
	GoName       string // Go 字段名(如 UserID)
	ColumnName   string // 数据库列名(如 user_id)
	GoType       string // Go 类型(如 int64)
	DBType       string // 数据库类型(从 gorm tag type:xxx 提取,如 bigint/varchar(64))
	Tag          string // 完整 gorm tag
	PrimaryKey   bool
	AutoIncr     bool
	NotNull      bool
	IsIndex      bool
	UniqueIndex  bool
	DefaultValue string // DEFAULT 子句
	Comment      string
	HasComment   bool
}

// TableMeta 单张表的元数据(从 struct 提取)
type TableMeta struct {
	StructName   string       // Go 类型名(如 CmsAdmin)
	TableName    string       // 数据库表名(如 cms_admin)
	Columns      []ColumnMeta // 业务字段
	AllFieldRefs []ColumnMeta // 所有字段(含 tenant_id/deleted)
	Indexes      []IndexMeta  // 索引定义
}

// IndexMeta 索引元数据
type IndexMeta struct {
	Name    string
	Columns []string
	Unique  bool
	Comment string // 来源字段的 comment(从 struct field 推导)
}

// extractTableMeta 从 domain struct 提取表元数据
// 使用反射解析 gorm tag,提取表名、字段、列注释、索引
func extractTableMeta(model any) (TableMeta, error) {
	var meta TableMeta
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return meta, errors.New("model must be a struct or pointer to struct")
	}

	meta.StructName = t.Name()

	// 提取 TableName(调用 TableName() 方法或使用类型名)
	if tn, ok := model.(TableNameAble); ok {
		meta.TableName = tn.TableName()
	} else {
		meta.TableName = toSnakeCase(t.Name())
	}

	// 遍历字段
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}

		// 跳过嵌入字段(暂不支持关联)
		if f.Anonymous {
			continue
		}

		tag := f.Tag.Get("gorm")
		if tag == "" {
			continue
		}

		cm := ColumnMeta{
			GoName: f.Name,
			GoType: f.Type.String(),
			Tag:    tag,
		}

		// 解析 column:xxx
		cm.ColumnName = parseTagValue(tag, "column")
		if cm.ColumnName == "" {
			cm.ColumnName = toSnakeCase(f.Name)
		}

		// 解析 type:xxx
		cm.DBType = parseTagValue(tag, "type")

		// 解析 primaryKey / autoIncrement
		cm.PrimaryKey = tagContains(tag, "primaryKey")
		cm.AutoIncr = tagContains(tag, "autoIncrement")
		cm.NotNull = tagContains(tag, "not null")

		// 解析 index (注意:index:col_name 是命名索引,单独的 index 是匿名)
		idxVal := parseTagValue(tag, "index")
		if idxVal != "" {
			cm.IsIndex = true
		} else if tagContains(tag, "index") {
			cm.IsIndex = true
		}

		// 解析 uniqueIndex / unique
		if tagContains(tag, "uniqueIndex") || tagContains(tag, "unique") {
			cm.UniqueIndex = true
		}

		// 解析 default:xxx
		cm.DefaultValue = parseTagValue(tag, "default")

		// 解析 comment:xxx
		cm.Comment = parseTagValue(tag, "comment")
		cm.HasComment = cm.Comment != ""

		meta.AllFieldRefs = append(meta.AllFieldRefs, cm)
	}

	return meta, nil
}

// parseTagValue 解析 gorm tag 中的 key:value 部分
// 找到 key: 后读取直到 ; 或字符串末尾
func parseTagValue(tag, key string) string {
	key = key + ":"
	idx := strings.Index(tag, key)
	if idx < 0 {
		return ""
	}
	start := idx + len(key)
	end := strings.Index(tag[start:], ";")
	if end < 0 {
		return strings.TrimSpace(tag[start:])
	}
	return strings.TrimSpace(tag[start : start+end])
}

// tagContains 检查 gorm tag 是否包含某个独立标记
// 通过 ; 分割后精确匹配
func tagContains(tag, key string) bool {
	for _, part := range strings.Split(tag, ";") {
		if strings.TrimSpace(part) == key {
			return true
		}
	}
	return false
}

// toSnakeCase 将 CamelCase 转为 snake_case
func toSnakeCase(s string) string {
	var b strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			b.WriteByte('_')
		}
		if r >= 'A' && r <= 'Z' {
			b.WriteByte(byte(r) + 32)
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// ============================================================
// 检查列是否存在 (跨方言)
// ============================================================

// columnExists 检查表的某列是否存在(通过执行方言特定的检查 SQL)
func columnExists(db *DBExtension, dialect, tableName, columnName string) (bool, error) {
	sql := GetCheckColumnSQL(dialect, tableName, columnName)
	var count int64
	if err := db.Raw(sql).Scan(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// tableExists 检查表是否存在
func tableExists(db *DBExtension, tableName string) (bool, error) {
	var count int64
	if err := db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_name = ?", tableName).Scan(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// ============================================================
// 列存在性查询(MySQL) — 给外部包使用
// ============================================================

// HasColumn 检查某列是否存在(基于当前数据库方言)
func HasColumn(db *DBExtension, dialect, tableName, columnName string) bool {
	ok, _ := columnExists(db, dialect, tableName, columnName)
	return ok
}

// ============================================================
// 自动迁移:列存在性检查 + 添加列
// ============================================================

// EnsureColumn 检查并添加缺失的列
// 返回: true 表示已添加(或重新同步),false 表示已存在且无需变更
func EnsureColumn(db *DBExtension, dialect, tableName, columnName, dbType string, notNull bool, defValue, comment string) (bool, error) {
	exists, err := columnExists(db, dialect, tableName, columnName)
	if err != nil {
		return false, err
	}
	if exists {
		return false, nil
	}

	// 列不存在,添加
	sql := GetAddColumnSQL(dialect, tableName, columnName, dbType, notNull, defValue)
	if err := db.Exec(sql).Error; err != nil {
		return false, fmt.Errorf("添加列 %s.%s 失败: %w", tableName, columnName, err)
	}
	logs.Info("已添加列: %s.%s (%s)", tableName, columnName, dbType)
	return true, nil
}

// EnsureIndex 检查并创建缺失的索引
func EnsureIndex(db *DBExtension, dialect, tableName, indexName string, columns []string, unique bool) (bool, error) {
	// 简化:尝试创建,若已存在则忽略错误
	sql := GetAddIndexSQL(dialect, tableName, indexName, columns, unique)
	if err := db.Exec(sql).Error; err != nil {
		// 索引已存在错误可忽略(各方言错误码不同)
		if isIndexExistsError(err) {
			return false, nil
		}
		return false, fmt.Errorf("创建索引 %s 失败: %w", indexName, err)
	}
	logs.Info("已创建索引: %s ON %s (%v)", indexName, tableName, columns)
	return true, nil
}

// isIndexExistsError 索引已存在错误的容错判断
func isIndexExistsError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	// MySQL: 1061 (Duplicate key name), 1826 (Duplicate foreign key)
	// PostgreSQL: 42710 (duplicate object)
	// SQL Server: 2714 (There is already an object named)
	if strings.Contains(msg, "duplicate key name") ||
		strings.Contains(msg, "already exists") ||
		strings.Contains(msg, "duplicate object") ||
		strings.Contains(msg, "there is already") {
		return true
	}
	return false
}

// SetTableComment 设置表注释(覆盖式)
func SetTableComment(db *DBExtension, dialect, tableName, comment string) error {
	if comment == "" {
		return nil
	}
	sql := GetTableCommentSQL(dialect, tableName, comment)
	if err := db.Exec(sql).Error; err != nil {
		return fmt.Errorf("设置表 %s 注释失败: %w", tableName, err)
	}
	return nil
}

// SetColumnComment 设置列注释(必须保留类型定义,否则会改类型)
func SetColumnComment(db *DBExtension, dialect, tableName, columnName, dbType, comment string) error {
	if comment == "" {
		return nil
	}
	sql := GetColumnCommentSQL(dialect, tableName, columnName, dbType, comment)
	if err := db.Exec(sql).Error; err != nil {
		return fmt.Errorf("设置列 %s.%s 注释失败: %w", tableName, columnName, err)
	}
	return nil
}

// ============================================================
// 完整迁移:从 domain struct 自动同步到数据库
// ============================================================

// AutoMigrateWithMeta 完整的自动迁移:
//   1. CREATE TABLE IF NOT EXISTS(由 GORM 处理)
//   2. 添加缺失的列
//   3. 创建缺失的索引
//   4. 设置表注释
//   5. 设置所有列注释
func AutoMigrateWithMeta(db *DBExtension, dialect string, model any, tableComment string) error {
	if db == nil {
		return errors.New("db is nil")
	}

	// 1. GORM AutoMigrate 处理表创建 + 列添加
	if err := db.AutoMigrate(model); err != nil {
		logs.Warn("GORM AutoMigrate 失败[%T]: %v", model, err)
		// 不 return,继续做后续手工补充
	}

	meta, err := extractTableMeta(model)
	if err != nil {
		return fmt.Errorf("提取元数据失败: %w", err)
	}

	// 2. 遍历每个字段,确保列存在 + 索引存在 + 注释正确
	for _, col := range meta.AllFieldRefs {
		// 跳过非业务字段(嵌入的 TenantID/Deleted 会经过表名,需要保留)
		// 已由 GORM 处理,这里主要补充注释

		// 设置列注释
		if col.HasComment && col.DBType != "" {
			if err := SetColumnComment(db, dialect, meta.TableName, col.ColumnName, col.DBType, col.Comment); err != nil {
				logs.Warn("设置列注释失败 %s.%s: %v", meta.TableName, col.ColumnName, err)
			}
		}

		// 索引处理
		if col.IsIndex && !col.PrimaryKey {
			idxName := "idx_" + meta.TableName + "_" + col.ColumnName
			if _, err := EnsureIndex(db, dialect, meta.TableName, idxName, []string{col.ColumnName}, col.UniqueIndex); err != nil {
				logs.Warn("创建索引失败 %s: %v", idxName, err)
			}
		}
	}

	// 3. 设置表注释
	if tableComment != "" {
		if err := SetTableComment(db, dialect, meta.TableName, tableComment); err != nil {
			logs.Warn("设置表注释失败 %s: %v", meta.TableName, err)
		}
	}

	return nil
}

// ============================================================
// 跨方言兼容的类型修复(bit(1) → tinyint(1))
// ============================================================

// FixBit1ToTinyInt1 将指定表的所有 deleted 列从 bit(1) 转为 tinyint(1)
// 用于修复历史部署中 bit(1) 无法被 GORM 扫描为 bool 的问题
func FixBit1ToTinyInt1(db *DBExtension, dialect, tableName string) error {
	switch dialect {
	case DialectMySQL:
		// 检查当前类型
		var colType string
		err := db.Raw("SELECT COLUMN_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = 'deleted'", tableName).Scan(&colType).Error
		if err != nil {
			return err
		}
		if strings.Contains(strings.ToLower(colType), "bit") {
			sql := GetModifyColumnSQL(dialect, tableName, "deleted", "tinyint(1) NOT NULL DEFAULT 0")
			if err := db.Exec(sql).Error; err != nil {
				return fmt.Errorf("修复 %s.deleted 类型失败: %w", tableName, err)
			}
			logs.Info("已修复列类型 bit(1)→tinyint(1): %s.deleted", tableName)
		}
	case DialectPostgres, DialectKingbase, DialectGaussDB:
		// PG 系不存在 bit(1) 问题,默认 boolean 类型
		return nil
	case DialectSQLServer:
		// MSSQL 用 BIT 类型,无问题
		return nil
	}
	return nil
}

// ============================================================
// 数据库支持函数(扩展 GORM 的 DB.Exec 错误处理)
// ============================================================

// IsSQLNoSuchTable 检查是否为"表不存在"错误
func IsSQLNoSuchTable(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "no such table") ||
		strings.Contains(msg, "doesn't exist") ||
		strings.Contains(msg, "1146") // MySQL error code
}

// ============================================================
// 检查 sql 包的 sql.NullString 等辅助
// ============================================================

var _ sql.NullString // 引用 database/sql 包,保证导入
