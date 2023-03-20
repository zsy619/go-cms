package main

import (
	"gorm.io/gen"
	"haedu.gov.cn/cms/app/dal"
)

// 参考：https://blog.csdn.net/Jeffid/article/details/126898000
func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath:      "./app/dal/query",
		ModelPkgPath: "./app/dal/model", // 默认情况下会跟随OutPath参数，在同目录下生成model目录
		Mode:         gen.WithDefaultQuery,
		WithUnitTest: false,

		// 表字段可为 null 值时, 对应结体字段使用指针类型
		FieldNullable: false, // generate pointer when field is nullable

		// 表字段默认值与模型结构体字段零值不一致的字段, 在插入数据时需要赋值该字段值为零值的, 结构体字段须是指针类型才能成功, 即`FieldCoverable:true`配置下生成的结构体字段.
		// 因为在插入时遇到字段为零值的会被GORM赋予默认值. 如字段`age`表默认值为10, 即使你显式设置为0最后也会被GORM设为10提交.
		// 如果该字段没有上面提到的插入时赋零值的特殊需要, 则字段为非指针类型使用起来会比较方便.
		FieldCoverable: false, // generate pointer when field has default value, to fix problem zero value cannot be assign: https://gorm.io/docs/create.html#Default-Values

		// 模型结构体字段的数字类型的符号表示是否与表字段的一致, `false`指示都用有符号类型
		FieldSignable: false, // detect integer field's unsigned type, adjust generated data type

		// 生成 gorm 标签的字段索引属性
		FieldWithIndexTag: true, // generate with gorm index tag

		// 生成 gorm 标签的字段类型属性
		FieldWithTypeTag: true, // generate with gorm column type tag
	})

	g.UseDB(dal.CmsDatabase.DB)

	// 将特定字段名的 json 标签加上`string`属性,即 MarshalJSON 时该字段由数字类型转成字符串类型
	jsonField := gen.FieldJSONTagWithNS(func(columnName string) (tagContent string) {
		// toStringField := `balance, `
		// if strings.Contains(toStringField, columnName) {
		// 	return columnName + ",string"
		// }
		return columnName
	})
	// 将非默认字段名的字段定义为自动时间戳和软删除字段;
	// 自动时间戳默认字段名为:`updated_at`、`created_at, 表字段数据类型为: INT 或 DATETIME
	// 软删除默认字段名为:`deleted_at`, 表字段数据类型为: DATETIME
	autoUpdateTimeField := gen.FieldGORMTag("update_time", "column:update_time;type:int unsigned;autoUpdateTime")
	autoCreateTimeField := gen.FieldGORMTag("create_time", "column:create_time;type:int unsigned;autoCreateTime")
	softDeleteField := gen.FieldType("delete_time", "soft_delete.DeletedAt")
	formField := gen.FieldNewTagWithNS("form", func(columnName string) (tagContent string) {
		// toStringField := `balance, `
		// if strings.Contains(toStringField, columnName) {
		// 	return columnName + ",string"
		// }
		return columnName
	})
	// 模型自定义选项组
	fieldOpts := []gen.ModelOpt{formField, jsonField, autoCreateTimeField, autoUpdateTimeField, softDeleteField}

	// generate all table from database
	g.ApplyBasic(g.GenerateAllTable(fieldOpts...)...)
	// g.ApplyBasic(g.GenerateModel("gsg_bys_xsqy"))
	// g.ApplyBasic(g.GenerateAllTableWithSharding(`sso_{.}`)...)
	// g.GenerateAllTable()

	// g.ApplyBasic(g.GenerateModelAs("sso_micro_svc", "SsoMicroSvc"))

	g.Execute()
}

// dataMap mapping relationship
var dataMap = map[string]func(detailType string) (dataType string){
	// // int mapping
	// "int": func(detailType string) (dataType string) { return "int32" },

	// // bool mapping
	// "tinyint": func(detailType string) (dataType string) {
	// 	if strings.HasPrefix(detailType, "tinyint(1)") {
	// 		return "bool"
	// 	}
	// 	return "byte"
	// },
}
