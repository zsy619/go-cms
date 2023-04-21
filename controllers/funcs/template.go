package funcs

import (
	"bytes"
	"fmt"
	"html/template"
	"path"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
)

/**
 * @description: 根据theme加载模板
 * @param {*} theme 主题名称，如defalut、h5、pc等，对应views/themes目录下的文件夹，区分大小写
 * @param {string} file 模板文件名称，如index.html、list.html等，对应views/themes/主题名称/views目录下的文件
 * @param {map[interface{}]interface{}} b 模板变量
 * @return {*}
 */
func TemplateTheme(theme, file string, b map[interface{}]interface{}) template.HTML {
	url := path.Join("./views/themes", theme, "views", file)
	fmt.Println("TemplateTheme url ---->", url)
	tmpl := template.New(file)
	InitFuncType(tmpl)
	tmplResult, err := tmpl.ParseFiles(url)
	if err != nil {
		fmt.Println("TemplateTheme---1>", err.Error())
		logs.Error(err)
		return ""
	}
	// tmpl.Execute(os.Stdout, b)
	buf := &bytes.Buffer{}
	err = tmplResult.Execute(buf, b)
	if err != nil {
		fmt.Println("TemplateTheme---2>", err.Error())
		logs.Error(err)
		return ""
	}
	return web.Str2html(buf.String())
}

/**
 * @description: 根据views目录加载模板
 * @param {string} views views目录路径，如themes/default/views
 * @param {string} file 模板文件名称，如index.html、list.html等，对应views/themes/主题名称/views目录下的文件
 * @param {map[interface{}]interface{}} b 模板变量
 * @return {*}
 */
func TemplateView(views, file string, b map[interface{}]interface{}) template.HTML {
	url := path.Join(views, "views", file)
	fmt.Println("TemplateView url ---->", url)
	tmpl := template.New(file)
	InitFuncType(tmpl)
	tmplResult, err := tmpl.ParseFiles(url)
	if err != nil {
		fmt.Println("TemplateView---1>", err.Error())
		logs.Error(err)
		return ""
	}
	// tmpl.Execute(os.Stdout, b)
	buf := &bytes.Buffer{}
	err = tmplResult.Execute(buf, b)
	if err != nil {
		fmt.Println("TemplateView---2>", err.Error())
		logs.Error(err)
		return ""
	}
	return web.Str2html(buf.String())
}

func UrlForView(theme, file string) string {
	return path.Join("themes", theme, "views", file)
}
