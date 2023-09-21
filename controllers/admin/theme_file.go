package admin

import (
	"io/ioutil"
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
func (c *ThemeController) LoadFiles() {
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
	theme := c.GetSafeString("theme")
	if theme == "" {
		result.Code = -1
		result.Msg = "主题名称不能为空"
		c.Data["json"] = result
		c.ServeJSON()
		return
	}
	themeDir := filepath.Join("./views/themes/", theme)
	viewPath := filepath.Join(themeDir, "views")
	jsPath := filepath.Join(themeDir, "static", "js")
	cssPath := filepath.Join(themeDir, "static", "css")
	imagePath := filepath.Join(themeDir, "static", "images")

	result.Views, _ = c.listFiles(viewPath)
	result.Js, _ = c.listFiles(jsPath)
	result.Css, _ = c.listFiles(cssPath)
	result.Images, _ = c.listFiles(imagePath)
	c.Data["json"] = result
	c.ServeJSON()
}

/**
 * @description: 根据类型获取主题文件列表
 * @return {*}
 */
func (c *ThemeController) LoadFilesType() {
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
	theme := c.GetSafeString("theme")
	typex := c.GetSafeString("type")
	if theme == "" {
		result.Code = -1
		result.Msg = "主题名称不能为空"
		c.Data["json"] = result
		c.ServeJSON()
		return
	}
	themeDir := filepath.Join("./views/themes/", theme)

	switch typex {
	case "View":
		viewPath := filepath.Join(themeDir, "views")
		result.Views, _ = c.listFiles(viewPath)
		break
	case "js":
		jsPath := filepath.Join(themeDir, "static", "js")
		result.Js, _ = c.listFiles(jsPath)
		break
	case "css":
		cssPath := filepath.Join(themeDir, "static", "css")
		result.Css, _ = c.listFiles(cssPath)
		break
	default:
		imagePath := filepath.Join(themeDir, "static", "images")
		result.Images, _ = c.listFiles(imagePath)
		break
	}
	c.Data["json"] = result
	c.ServeJSON()
}

/**
 * @description: 获取指定目录下所有文件的路径列表
 * @param {string} dir 目录路径
 * @return {*}
 */
func (c *ThemeController) listFiles(dir string) ([]string, error) {
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
func (c *ThemeController) FileContent() {
	theme := c.GetSafeString("theme")
	file := c.GetSafeString("file")
	typex := c.GetSafeString("type")
	if theme == "" || file == "" || typex == "" {
		c.JSONError("参数错误")
		return
	}
	if typex != "View" {
		typex = "static/" + typex
	} else {
		typex = "views"
	}
	themeDir := filepath.Join("./views/themes/", theme)
	filename := "./" + filepath.Join(themeDir, strings.ToLower(typex), file)
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		logs.Error(err.Error())
		c.JSONError("读取文件失败")
		return
	}
	c.JSONSuccess("", string(content))
}

/**
 * @description: 创建文件
 * @return {*}
 */
func (c *ThemeController) FileAdd() {
	theme := c.GetSafeString("theme")
	file := c.GetSafeString("file")
	typex := c.GetSafeString("type")
	if theme == "" || file == "" || typex == "" {
		c.JSONError("参数错误")
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
		c.JSONError("文件已存在")
		return
	}
	filex, err := os.Create(filename)
	if err != nil {
		logs.Error(err.Error())
		c.JSONError("创建文件失败")
		return
	}
	defer filex.Close()
	c.JSONSuccess("创建文件成功", nil)
}

/**
 * @description: 保存文件内容
 * @return {*}
 */
func (c *ThemeController) FileSave() {
	theme := c.GetSafeString("theme")
	file := c.GetSafeString("file")
	typex := c.GetSafeString("type")
	content := c.GetSafeString("content")
	if theme == "" || file == "" || typex == "" {
		c.JSONError("参数错误")
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
	}
	err := ioutil.WriteFile(filename, []byte(content), 0o644)
	if err != nil {
		logs.Error(err.Error())
		c.JSONError("保存文件失败")
		return
	}
	c.JSONSuccess("保存文件成功", nil)
}

/**
 * @description: 删除文件
 * @return {*}
 */
func (c *ThemeController) FileDelete() {
	theme := c.GetSafeString("theme")
	file := c.GetSafeString("file")
	typex := c.GetSafeString("type")
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
		c.JSONError("删除文件失败")
		return
	}
	c.JSONSuccess("删除文件成功", nil)
}
