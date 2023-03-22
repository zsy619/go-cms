package admin

import (
	"fmt"
	"math/rand"
	"os"
	"path"
	"strings"
	"time"

	"haedu.gov.cn/cms/app/lib"
	"haedu.gov.cn/tools/xgeneric"
)

type UploadResult struct {
	Code     int    `json:"code"`
	Msg      string `json:"msg"`
	FileExt  string `json:"fileExt"`  // 文件后缀
	FileUrl1 string `json:"fileUrl1"` // 文件相对路径
	FileUrl2 string `json:"fileUrl2"` // 文件绝对路径
	FileName string `json:"fileName"` // 保存文件名称
}

// ToolController 工具类接口
type ToolController struct{ BaseController }

// WeUpload 上传图片
// @router admin/tool/WeUpload [post]
func (c *ToolController) WeUpload() {
	result := UploadResult{}
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
	if !strings.Contains("jpg,jpeg,png,gif", ext) {
		result.Code = 1
		result.Msg = "不支持的文件类型"
		c.Data["json"] = result
		c.ServeJSON()
		return
	}

	// 生成唯一的文件名
	filename := generateFilename(ext)

	// 保存上传的文件到指定目录
	uploadDir := "Uploads/images" + time.Now().Format("2006/01/")
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
	result.FileExt = ext
	result.FileName = filename
	result.FileUrl1 = "/" + uploadDir + filename
	result.FileUrl2 = xgeneric.IFF(strings.HasSuffix(lib.C_LOCAL_DOMAIN(), "/"), lib.C_LOCAL_DOMAIN()+uploadDir+filename, lib.C_LOCAL_DOMAIN()+"/"+uploadDir+filename)
	c.Data["json"] = result
	c.ServeJSON()
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
