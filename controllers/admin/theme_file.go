package admin

import (
	"os"
	"path/filepath"
)

func (c *ThemeController) Files() {
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
	theme := c.GetString("theme")
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

	result.Views, _ = ListFiles(viewPath)
	result.Js, _ = ListFiles(jsPath)
	result.Css, _ = ListFiles(cssPath)
	result.Images, _ = ListFiles(imagePath)
	c.Data["json"] = result
	c.ServeJSON()
}

// ListFiles 获取指定目录下所有文件的路径列表
func ListFiles(dir string) ([]string, error) {
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
