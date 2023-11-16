package admin

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/tools/xio"
)

/**
 * @description: 获取主题文件列表
 * @return {*}
 */
func (ctrl *ThemeController) LoadFiles() {
	result := struct {
		Code   int      `json:"code"`
		Msg    string   `json:"msg"`
		Views  []string `json:"views"`
		Js     []string `json:"js"`
		Css    []string `json:"css"`
		Images []string `json:"images"`
	}{
		Code: 0,
	}
	theme := ctrl.GetSafeString("theme")
	if theme == "" {
		result.Code = -1
		result.Msg = "主题名称不能为空"
		ctrl.Data["json"] = result
		ctrl.ServeJSON()
		return
	}
	themeDir := filepath.Join("./views/themes/", theme)
	viewPath := filepath.Join(themeDir, "views")
	jsPath := filepath.Join(themeDir, "static", "js")
	cssPath := filepath.Join(themeDir, "static", "css")
	imagePath := filepath.Join(themeDir, "static", "images")

	result.Views, _ = ctrl.listFiles(viewPath)
	result.Js, _ = ctrl.listFiles(jsPath)
	result.Css, _ = ctrl.listFiles(cssPath)
	result.Images, _ = ctrl.listFiles(imagePath)
	ctrl.Data["json"] = result
	ctrl.ServeJSON()
}

/**
 * @description: 根据类型获取主题文件列表
 * @return {*}
 */
func (ctrl *ThemeController) LoadFilesType() {
	result := struct {
		Code   int      `json:"code"`
		Msg    string   `json:"msg"`
		Views  []string `json:"views"`
		Js     []string `json:"js"`
		Css    []string `json:"css"`
		Images []string `json:"images"`
	}{
		Code: 0,
	}
	theme := ctrl.GetSafeString("theme")
	typex := ctrl.GetSafeString("type")
	if theme == "" {
		result.Code = -1
		result.Msg = "主题名称不能为空"
		ctrl.Data["json"] = result
		ctrl.ServeJSON()
		return
	}
	themeDir := filepath.Join("./views/themes/", theme)

	switch typex {
	case "View":
		viewPath := filepath.Join(themeDir, "views")
		result.Views, _ = ctrl.listFiles(viewPath)
	case "js":
		jsPath := filepath.Join(themeDir, "static", "js")
		result.Js, _ = ctrl.listFiles(jsPath)
	case "css":
		cssPath := filepath.Join(themeDir, "static", "css")
		result.Css, _ = ctrl.listFiles(cssPath)
	default:
		imagePath := filepath.Join(themeDir, "static", "images")
		result.Images, _ = ctrl.listFiles(imagePath)
	}
	ctrl.Data["json"] = result
	ctrl.ServeJSON()
}

/**
 * @description: 获取指定目录下所有文件的路径列表
 * @param {string} dir 目录路径
 * @return {*}
 */
func (ctrl *ThemeController) listFiles(dir string) ([]string, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, filepath.Base(path))
		}
		return nil
	})
	if err != nil {
		return files, err
	}

	return files, nil
}

/**
 * @description: 获取文件内容
 * @return {*}
 */
func (ctrl *ThemeController) FileContent() {
	theme := ctrl.GetSafeString("theme")
	file := ctrl.GetSafeString("file")
	typex := ctrl.GetSafeString("type")
	if theme == "" || file == "" || typex == "" {
		ctrl.JSONError("参数错误")
		return
	}
	if typex != "View" {
		typex = "static/" + typex
	} else {
		typex = "views"
	}
	themeDir := filepath.Join("./views/themes/", theme)
	filename := "./" + filepath.Join(themeDir, strings.ToLower(typex), file)
	content, err := os.ReadFile(filename)
	if err != nil {
		logs.Error(err.Error())
		ctrl.JSONError("读取文件失败")
		return
	}
	ctrl.JSONSuccess("", string(content))
}

/**
 * @description: 创建文件
 * @return {*}
 */
func (ctrl *ThemeController) FileAdd() {
	theme := ctrl.GetSafeString("theme")
	file := ctrl.GetSafeString("file")
	typex := ctrl.GetSafeString("type")
	if theme == "" || file == "" || typex == "" {
		ctrl.JSONError("参数错误")
		return
	}
	if typex != "View" {
		typex = "static/" + typex
	} else {
		typex = "views"
	}
	themeDir := filepath.Join("./views/themes/", theme)
	filename := "./" + filepath.Join(themeDir, strings.ToLower(typex), file)
	if xio.FileIsExisted(filename) {
		ctrl.JSONError("文件已存在")
		return
	}
	filex, err := os.Create(filename)
	if err != nil {
		logs.Error(err.Error())
		ctrl.JSONError("创建文件失败")
		return
	}
	defer filex.Close()
	ctrl.JSONSuccess("创建文件成功", nil)
}

/**
 * @description: 保存文件内容
 * @return {*}
 */
func (ctrl *ThemeController) FileSave() {
	theme := ctrl.GetSafeString("theme")
	file := ctrl.GetSafeString("file")
	typex := ctrl.GetSafeString("type")
	content := ctrl.GetSafeString("content")
	if theme == "" || file == "" || typex == "" {
		ctrl.JSONError("参数错误")
		return
	}
	if typex != "View" {
		typex = "static/" + typex
	} else {
		typex = "views"
	}
	themeDir := filepath.Join("./views/themes/", theme)
	filename := "./" + filepath.Join(themeDir, strings.ToLower(typex), file)
	if !xio.IsFileExist(filename) {
		fmt.Println("file not exists")
	}
	err := os.WriteFile(filename, []byte(content), 0o644)
	if err != nil {
		logs.Error(err.Error())
		ctrl.JSONError("保存文件失败")
		return
	}
	ctrl.JSONSuccess("保存文件成功", nil)
}

/**
 * @description: 删除文件
 * @return {*}
 */
func (ctrl *ThemeController) FileDelete() {
	theme := ctrl.GetSafeString("theme")
	file := ctrl.GetSafeString("file")
	typex := ctrl.GetSafeString("type")
	if typex != "View" {
		typex = "static/" + typex
	} else {
		typex = "views"
	}
	themeDir := filepath.Join("./views/themes/", theme)
	filename := "./" + filepath.Join(themeDir, strings.ToLower(typex), file)
	err := os.Remove(filename)
	if err != nil {
		logs.Error(err.Error())
		ctrl.JSONError("删除文件失败")
		return
	}
	ctrl.JSONSuccess("删除文件成功", nil)
}
