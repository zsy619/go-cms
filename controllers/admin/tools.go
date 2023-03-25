package admin

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/cms/app/lib"
	"haedu.gov.cn/tools/xgeneric"
)

type UploadResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	File struct {
		Name string `json:"name"` // 保存文件名称
		Ext  string `json:"ext"`  // 文件后缀
		Url1 string `json:"url1"` // 文件相对路径
		Url2 string `json:"url2"` // 文件绝对路径
	} `json:"file"`
}

type WebUploadResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Url []string `json:"url"`
	} `json:"data"`
}

// ToolsController 工具类接口
type ToolsController struct{ BaseController }

// @router admin/tools/ImageUpload [post]
func (c *ToolsController) ImageUpload() {
	result := UploadResult{
		Code: 1,
		Msg:  "上传失败",
		File: struct {
			Name string "json:\"name\""
			Ext  string "json:\"ext\""
			Url1 string "json:\"url1\""
			Url2 string "json:\"url2\""
		}{},
	}
	file, head, err := c.GetFile("file")
	if err != nil {
		result.Code = 1
		result.Msg = "文件上传失败"
		c.Data["json"] = result
		c.ServeJSON()
		return
	}
	defer file.Close()

	ext := strings.ToLower(path.Ext(head.Filename))
	fmt.Println(ext)
	if !strings.Contains(".jpg,.jpeg,.png,.gif", ext) {
		result.Code = 1
		result.Msg = "不支持的文件类型"
		c.Data["json"] = result
		c.ServeJSON()
		return
	}

	// 生成唯一的文件名
	filename := generateFilename(ext)

	// 保存上传的文件到指定目录
	uploadDir := "Uploads/images/" + time.Now().Format("2006/01/")
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		fmt.Println("err", err.Error())
		result.Code = 1
		result.Msg = err.Error()
		c.Data["json"] = result
		c.ServeJSON()
		return
	}
	if err := c.SaveToFile("file", uploadDir+filename); err != nil {
		result.Code = 1
		result.Msg = "文件上传失败"
		c.Data["json"] = result
		c.ServeJSON()
		return
	}

	// 返回上传成功的文件 URL 地址
	result.Code = 0
	result.Msg = "succcess"
	result.File.Ext = ext
	result.File.Name = filename
	result.File.Url1 = "/" + uploadDir + filename
	result.File.Url2 = xgeneric.IFF(strings.HasSuffix(lib.C_LOCAL_DOMAIN(), "/"), lib.C_LOCAL_DOMAIN()+uploadDir+filename, lib.C_LOCAL_DOMAIN()+"/"+uploadDir+filename)
	c.Data["json"] = result
	c.ServeJSON()
}

// @router admin/tools/Upload [post]
func (c *ToolsController) Upload() {
	result := UploadResult{
		Code: 1,
		Msg:  "上传失败",
		File: struct {
			Name string "json:\"name\""
			Ext  string "json:\"ext\""
			Url1 string "json:\"url1\""
			Url2 string "json:\"url2\""
		}{},
	}
	file, head, err := c.GetFile("file")
	if err != nil {
		result.Code = 1
		result.Msg = "文件上传失败"
		c.Data["json"] = result
		c.ServeJSON()
		return
	}
	defer file.Close()

	ext := strings.ToLower(path.Ext(head.Filename))
	// if !strings.Contains("jpg,jpeg,png,gif", ext) {
	// 	result.Code = 1
	// 	result.Msg = "不支持的文件类型"
	// 	c.Data["json"] = result
	// 	c.ServeJSON()
	// 	return
	// }

	// 生成唯一的文件名
	filename := generateFilename(ext)

	// 保存上传的文件到指定目录
	uploadDir := "Uploads/files/" + time.Now().Format("2006/01/")
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		fmt.Println("err", err.Error())
		result.Code = 1
		result.Msg = err.Error()
		c.Data["json"] = result
		c.ServeJSON()
		return
	}
	if err := c.SaveToFile("file", uploadDir+filename); err != nil {
		result.Code = 1
		result.Msg = "文件上传失败"
		c.Data["json"] = result
		c.ServeJSON()
		return
	}

	// 返回上传成功的文件 URL 地址
	result.Code = 0
	result.Msg = "succcess"
	result.File.Ext = ext
	result.File.Name = filename
	result.File.Url1 = "/" + uploadDir + filename
	result.File.Url2 = xgeneric.IFF(strings.HasSuffix(lib.C_LOCAL_DOMAIN(), "/"), lib.C_LOCAL_DOMAIN()+uploadDir+filename, lib.C_LOCAL_DOMAIN()+"/"+uploadDir+filename)
	c.Data["json"] = result
	c.ServeJSON()
}

// kindeditor 图片上传
// @router admin/tools/KindEditorUpload [post]
func (c *ToolsController) KindEditorUpload() {
	// 获取上传文件
	file, header, err := c.GetFile("imgFile")
	if err != nil {
		logs.Error(err)
		fmt.Fprintln(c.Ctx.ResponseWriter, "{\"error\": 1, \"message\": \""+err.Error()+"\"}")
		return
	}
	defer file.Close()

	// 创建目标文件夹
	targetDir := "Uploads/images/" + time.Now().Format("2006/01/")
	if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
		logs.Error(err)
		fmt.Fprintln(c.Ctx.ResponseWriter, "{\"error\": 1, \"message\": \""+err.Error()+"\"}")
		return
	}
	// 创建目标文件
	ext := filepath.Ext(header.Filename)
	filename := strconv.FormatInt(time.Now().UnixNano(), 10) + ext
	targetPath := targetDir + filename
	targetFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Println(err)
		fmt.Fprintln(c.Ctx.ResponseWriter, "{\"error\": 1, \"message\": \""+err.Error()+"\"}")
		return
	}
	defer targetFile.Close()

	// 将上传文件拷贝到目标文件
	_, err = io.Copy(targetFile, file)
	if err != nil {
		log.Println(err)
		fmt.Fprintln(c.Ctx.ResponseWriter, "{\"error\": 1, \"message\": \""+err.Error()+"\"}")
		return
	}
	outPath := xgeneric.IFF(strings.HasSuffix(lib.C_LOCAL_DOMAIN(), "/"), lib.C_LOCAL_DOMAIN()+strings.ReplaceAll(targetPath, "\\", "/"), lib.C_LOCAL_DOMAIN()+"/"+strings.ReplaceAll(targetPath, "\\", "/"))
	// 返回上传结果
	result := "{\"error\": 0, \"url\": \"" + outPath + "\"}"
	fmt.Fprintln(c.Ctx.ResponseWriter, result)
}

func generateFilename(ext string) string {
	return fmt.Sprintf("%s%s", randomString(10), ext)
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var result strings.Builder
	for i := 0; i < length; i++ {
		randomIndex := rand.Intn(len(charset))
		result.WriteByte(charset[randomIndex])
	}
	return result.String()
}
