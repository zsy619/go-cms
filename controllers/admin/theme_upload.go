package admin

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/controllers/admin/vmodel"
	"haedu.gov.cn/tools/xgeneric"
)

func (c *ThemeController) uploadMsg(code int32, msg string) {
	c.Data["json"] = map[string]interface{}{
		"code": code,
		"msg":  msg,
	}
	c.ServeJSON()
}

// Upload 上传
// @router /admin/theme/upload [post]
func (c *ThemeController) Upload() {
	// 处理步骤
	// 1、文件上传，保存到临时目录
	// 2、解压上传的 ZIP 文件
	// 3、检测主题是否已经存在
	// 4、检测必要文件是否存在
	// 5、拷贝主题文件到主题目录，并入库
	// 6、删除临时目录

	// 1. 文件上传，保存到临时目录
	file, header, err := c.GetFile("file")
	if err != nil {
		c.uploadMsg(-1, err.Error())
		return
	}
	defer file.Close()

	tmpDir := "./Uploads/temp" // 临时目录，配置在 beego 的配置文件中
	if err := os.MkdirAll(tmpDir, os.ModePerm); err != nil {
		c.uploadMsg(-1, err.Error())
		return
	}
	filePath := filepath.Join(tmpDir, header.Filename)
	err = c.SaveToFile("file", filePath)
	if err != nil {
		c.uploadMsg(-1, err.Error())
		return
	}

	themeName := strings.TrimSuffix(header.Filename, ".zip")
	srcDir := filepath.Join(tmpDir, themeName)

	// 2. 解压上传的 ZIP 文件
	if _, err := c.unzip(filePath, srcDir); err != nil {
		c.uploadMsg(-1, err.Error())
		return
	}

	// 3. 检测主题是否已经存在
	themeDir := filepath.Join("./views/themes", themeName)
	service := biz.NewCmsTheme()
	if c.fileExists(themeDir) || service.ThemeExists(themeName) {
		c.uploadMsg(-1, fmt.Sprintf("主题 %s 已经存在", themeName))
		return
	}

	// 4. 检测必要文件是否存在

	if !c.fileExists(filepath.Join(srcDir, "thumb.png")) {
		fmt.Println(filepath.Join(srcDir, "thumb.png"))
		c.uploadMsg(-1, "缩略图 thumb.png 不存在")
		return
	}
	if !c.fileExists(filepath.Join(srcDir, "config.json")) {
		c.uploadMsg(-1, "配置文件 config.json 不存在")
		return
	}
	themConfig := vmodel.ThemeConfig{}
	// 读取配置文件
	{
		fileConfig, err := os.Open(filepath.Join(srcDir, "config.json"))
		if err != nil {
			c.uploadMsg(-1, err.Error())
			return
		}
		defer fileConfig.Close()

		data, err := ioutil.ReadAll(fileConfig)
		if err != nil {
			c.uploadMsg(-1, err.Error())
			return
		}
		err = json.Unmarshal(data, &themConfig)
		if err != nil {
			c.uploadMsg(-1, err.Error())
			return
		}
		if themConfig.Name == "" {
			c.uploadMsg(-1, "config.json 中的 name 字段不能为空")
			return
		}
		if themConfig.Name != themeName {
			c.uploadMsg(-1, "config.json 中的 name 字段必须与主题文件名相同")
			return
		}
	}
	viewFiles := map[string]string{
		"index.html": "首页", "channel.html": "频道", "category.html": "栏目", "article.html": "文章详情", "search.html": "搜索页",
		"tag.html": "标签", "topic.html": "专题",
	}
	for k, v := range viewFiles {
		if !c.fileExists(filepath.Join(srcDir, "views", k)) {
			c.uploadMsg(-1, fmt.Sprintf("views目录下必要文件 %s %s 不存在", v, k))
			return
		}
	}
	staticDir := filepath.Join(srcDir, "static")
	if !c.fileExists(staticDir) {
		os.MkdirAll(srcDir, os.ModePerm)
	}
	staticFiles := map[string]string{"css": "样式", "js": "脚本", "images": "图片"}
	for k := range staticFiles {
		if !c.fileExists(filepath.Join(staticDir, k)) {
			os.MkdirAll(filepath.Join(staticDir, k), os.ModePerm)
		}
	}

	// 5. 拷贝主题文件到主题目录，并入库
	err = os.MkdirAll(themeDir, os.ModePerm)
	if err != nil {
		c.uploadMsg(-1, err.Error())
		return
	}
	if err := CopyDir(srcDir, themeDir); err != nil {
		c.uploadMsg(-1, err.Error())
		return
	}

	// TODO: 入库操作
	themeModel := &model.CmsTheme{
		Name:      themeName,
		Title:     themConfig.Title,
		Remark:    themConfig.Description,
		IsDefault: false,
		IsSystem:  false,
		Thumb:     "thumb.png",
		Version:   themConfig.Version,
		Author:    xgeneric.IFF(themConfig.Author == "", "教育网", themConfig.Author),
	}
	if err := service.ThemeSave(themeModel); err != nil {
		c.uploadMsg(-1, err.Error())
		return
	}

	// 6. 删除临时目录
	err = os.RemoveAll(tmpDir)
	if err != nil {
		c.uploadMsg(-1, err.Error())
		return
	}

	// 返回上传成功的信息
	c.Data["json"] = map[string]interface{}{
		"code": 0,
		"msg":  "上传成功",
	}
	c.ServeJSON()
}

// 解压 ZIP 文件
func (c *ThemeController) unzip(src string, dest string) ([]string, error) {
	var files []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return files, err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return files, err
		}
		defer rc.Close()

		path := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, os.ModePerm)
		} else {
			err = os.MkdirAll(filepath.Dir(path), os.ModePerm)
			if err != nil {
				return files, err
			}

			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return files, err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return files, err
			}

			files = append(files, path)
		}
	}

	return files, nil
}

func (c *ThemeController) fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func (c *ThemeController) copyFile(src string, dest string) error {
	input, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dest, input, 0o644)
	if err != nil {
		return err
	}

	return nil
}

// CopyDir 复制文件夹，将源文件夹下的所有内容复制到目标文件夹下
func CopyDir(src string, dest string) error {
	srcStat, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !srcStat.IsDir() {
		return fmt.Errorf("%s is not a directory", src)
	}

	err = os.MkdirAll(dest, srcStat.Mode())
	if err != nil {
		return err
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		if entry.IsDir() {
			err = CopyDir(srcPath, destPath)
			if err != nil {
				return err
			}
		} else {
			err = CopyFile(srcPath, destPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// CopyFile 复制文件
func CopyFile(src string, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}
