package funcs

import (
	"bytes"
	"html/template"
	"path"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
)

// TemplateTheme 根据theme加载模板
// @param theme 主题名称，如defalut、h5、pc等，对应views/themes目录下的文件夹，区分大小写
// @param file 模板文件名称，如index.html、list.html等，对应views/themes/主题名称/views目录下的文件
// @param b 模板变量
// @return template.HTML 渲染后的HTML
func TemplateTheme(theme, file string, b map[any]any) template.HTML {
	url := path.Join("./views/themes", theme, "views", file)
	logs.Debug("TemplateTheme url: %s", url)
	tmpl := template.New(file)
	InitFuncs(tmpl)
	tmplResult, err := tmpl.ParseFiles(url)
	if err != nil {
		logs.Error("TemplateTheme解析失败: url=%s, error=%v", url, err)
		return ""
	}
	buf := &bytes.Buffer{}
	err = tmplResult.Execute(buf, b)
	if err != nil {
		logs.Error("TemplateTheme执行失败: url=%s, error=%v", url, err)
		return ""
	}
	return web.Str2html(buf.String())
}

// TemplateView 根据views目录加载模板
// @param views views目录路径，如themes/default/views
// @param file 模板文件名称，如index.html、list.html等，对应views/themes/主题名称/views目录下的文件
// @param b 模板变量
// @return template.HTML 渲染后的HTML
func TemplateView(views, file string, b map[any]any) template.HTML {
	url := path.Join(views, "views", file)
	logs.Debug("TemplateView url: %s", url)
	tmpl := template.New(file)
	InitFuncs(tmpl)
	tmplResult, err := tmpl.ParseFiles(url)
	if err != nil {
		logs.Error("TemplateView解析失败: url=%s, error=%v", url, err)
		return ""
	}
	buf := &bytes.Buffer{}
	err = tmplResult.Execute(buf, b)
	if err != nil {
		logs.Error("TemplateView执行失败: url=%s, error=%v", url, err)
		return ""
	}
	return web.Str2html(buf.String())
}

// UrlForView 获取模板路径
// @param theme 主题名称，如defalut、h5、pc等，对应views/themes目录下的文件夹，区分大小写
// @param file 视图名称 如index.html、list.html等，对应views/themes/主题名称/views目录下的文件
// @return string 模板路径
func UrlForView(theme, file string) string {
	return path.Join("themes", theme, "views", file)
}