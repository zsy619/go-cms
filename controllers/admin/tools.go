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
	"haedu.gov.cn/tools/xio"

	"haedu.gov.cn/cms/app/cms/domain"
	"haedu.gov.cn/cms/app/cms/service"
	lib "haedu.gov.cn/cms/app/tool"
)

type UploadResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	File struct {
		Name1 string `json:"name1"` // 原始文件名称
		Name2 string `json:"name2"` // 保存文件名称
		Ext   string `json:"ext"`   // 文件后缀
		Size  int64  `json:"size"`  // 文件大小
		Url1  string `json:"url1"`  // 文件相对路径
		Url2  string `json:"url2"`  // 文件绝对路径
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
func (ctrl *ToolsController) ImageUpload() {
	result := UploadResult{
		Code: 1,
		Msg:  "上传失败",
		File: struct {
			Name1 string `json:"name1"` // 原始文件名称
			Name2 string `json:"name2"` // 保存文件名称
			Ext   string `json:"ext"`   // 文件后缀
			Size  int64  `json:"size"`  // 文件大小
			Url1  string `json:"url1"`  // 文件相对路径
			Url2  string `json:"url2"`  // 文件绝对路径
		}{},
	}
	file, head, err := ctrl.GetFile("file")
	if err != nil {
		result.Code = 1
		result.Msg = "文件上传失败"
		ctrl.Data["json"] = result
		ctrl.ServeJSON()
		return
	}
	defer file.Close()

	result.File.Name1 = head.Filename
	result.File.Size = head.Size

	ext := strings.ToLower(path.Ext(head.Filename))
	// fmt.Println(ext)
	if !strings.Contains(".jpg,.jpeg,.png,.gif,.bmp,.ico", ext) {
		result.Code = 1
		result.Msg = "不支持的文件类型"
		ctrl.Data["json"] = result
		ctrl.ServeJSON()
		return
	}

	// 生成唯一的文件名
	filename := generateFilename(ext)
	pth := ctrl.GetSafeString("path")
	table := ctrl.GetSafeString("table")
	if pth == "" {
		pth = "images"
	}
	if table == "" {
		table = pth
	}
	// 保存上传的文件到指定目录
	uploadDir := "Uploads/" + pth + "/" + time.Now().Format("2006/01/")
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		fmt.Println("err", err.Error())
		result.Code = 1
		result.Msg = err.Error()
		ctrl.Data["json"] = result
		ctrl.ServeJSON()
		return
	}
	if err := ctrl.SaveToFile("file", uploadDir+filename); err != nil {
		result.Code = 1
		result.Msg = "文件上传失败"
		ctrl.Data["json"] = result
		ctrl.ServeJSON()
		return
	}

	// 返回上传成功的文件 URL 地址
	result.Code = 0
	result.Msg = "succcess"
	result.File.Ext = ext
	result.File.Name2 = filename
	result.File.Url1 = "/" + uploadDir + filename
	result.File.Url2 = lib.C_LOCAL_DOMAIN_Backslash() + uploadDir + filename

	// 保存到数据库
	{
		mdl := domain.CmsAlbum{
			TableName_:   table,
			RecordID:     0,
			TypeID:       0,
			Title:        filename,
			ThumbPath:    result.File.Url2,
			OriginalPath: result.File.Url1,
			FilePath:     result.File.Url2,
			FileSize:     result.File.Size,
			FileExt:      result.File.Ext,
			LinkURL:      "",
			Click:        0,
			IsShow:       1,
			CreateID:     int32(GlobalAdminId),
		}
		if err := service.NewCmsAlbum().AlbumSave(&mdl); err != nil {
			logs.Error("ImageUpload AlbumSave", err.Error())
		}
	}

	// 写入日志
	service.NewCmsAdmin().LoginLog(GlobalAdminId, GlobalAdminName, "ImageUpload", result.File.Url1, "", "OK", ctrl.GetClientIp())

	ctrl.Data["json"] = result
	ctrl.ServeJSON()
}

// @router admin/tools/Upload [post]
func (ctrl *ToolsController) Upload() {
	result := UploadResult{
		Code: 1,
		Msg:  "上传失败",
		File: struct {
			Name1 string `json:"name1"` // 原始文件名称
			Name2 string `json:"name2"` // 保存文件名称
			Ext   string `json:"ext"`   // 文件后缀
			Size  int64  `json:"size"`  // 文件大小
			Url1  string `json:"url1"`  // 文件相对路径
			Url2  string `json:"url2"`  // 文件绝对路径
		}{},
	}
	file, head, err := ctrl.GetFile("file")
	if err != nil {
		result.Code = 1
		result.Msg = "文件上传失败"
		ctrl.Data["json"] = result
		ctrl.ServeJSON()
		return
	}
	defer file.Close()
	if isAllow := xio.IsAllowFile(head.Filename); !isAllow {
		result.Code = 1
		result.Msg = "不支持的文件类型"
		ctrl.Data["json"] = result
		ctrl.ServeJSON()
		return
	}

	result.File.Name1 = head.Filename
	result.File.Size = head.Size

	ext := strings.ToLower(path.Ext(head.Filename))
	// if !strings.Contains("jpg,jpeg,png,gif", ext) {
	// 	result.Code = 1
	// 	result.Msg = "不支持的文件类型"
	// 	ctrl.Data["json"] = result
	// 	ctrl.ServeJSON()
	// 	return
	// }

	// 生成唯一的文件名
	filename := generateFilename(ext)
	pth := ctrl.GetSafeString("path")
	table := ctrl.GetSafeString("table")
	if pth == "" {
		pth = "files"
	}
	if table == "" {
		table = pth
	}
	// 保存上传的文件到指定目录
	uploadDir := "Uploads/" + pth + "/" + time.Now().Format("2006/01/")
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		fmt.Println("err", err.Error())
		result.Code = 1
		result.Msg = err.Error()
		ctrl.Data["json"] = result
		ctrl.ServeJSON()
		return
	}
	if err := ctrl.SaveToFile("file", uploadDir+filename); err != nil {
		result.Code = 1
		result.Msg = "文件上传失败"
		ctrl.Data["json"] = result
		ctrl.ServeJSON()
		return
	}

	// 返回上传成功的文件 URL 地址
	result.Code = 0
	result.Msg = "succcess"
	result.File.Ext = ext
	result.File.Name2 = filename
	result.File.Url1 = "/" + uploadDir + filename
	result.File.Url2 = lib.C_LOCAL_DOMAIN_Backslash() + uploadDir + filename

	// 保存到数据库
	{
		mdl := domain.CmsAttach{
			TableName_:   table,
			RecordID:     0,
			TypeID:       0,
			Title:        filename,
			FilePath:     result.File.Url2,
			OriginalPath: result.File.Url1,
			FileSize:     result.File.Size,
			FileExt:      result.File.Ext,
			Click:        0,
			IsShow:       1,
			CreateID:     int32(GlobalAdminId),
		}
		if err := service.NewCmsAttach().AttachSave(&mdl); err != nil {
			logs.Error("Upload AttachSave", err.Error())
		}
	}
	// 写入日志
	service.NewCmsAdmin().LoginLog(GlobalAdminId, GlobalAdminName, "UploadFile", result.File.Url1, "", "OK", ctrl.GetClientIp())

	ctrl.Data["json"] = result
	ctrl.ServeJSON()
}

// kindeditor 图片上传
// @router admin/tools/KindEditorUpload [post]
func (ctrl *ToolsController) KindEditorUpload() {
	// 获取上传文件
	file, header, err := ctrl.GetFile("imgFile")
	if err != nil {
		logs.Error(err)
		fmt.Fprintln(ctrl.Ctx.ResponseWriter, "{\"error\": 1, \"message\": \""+err.Error()+"\"}")
		return
	}
	defer file.Close()

	// 创建目标文件夹
	targetDir := "Uploads/images/" + time.Now().Format("2006/01/")
	if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
		logs.Error(err)
		fmt.Fprintln(ctrl.Ctx.ResponseWriter, "{\"error\": 1, \"message\": \""+err.Error()+"\"}")
		return
	}
	// 创建目标文件
	ext := filepath.Ext(header.Filename)
	filename := strconv.FormatInt(time.Now().UnixNano(), 10) + ext
	targetPath := targetDir + filename
	targetFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Println(err)
		fmt.Fprintln(ctrl.Ctx.ResponseWriter, "{\"error\": 1, \"message\": \""+err.Error()+"\"}")
		return
	}
	defer targetFile.Close()

	// 将上传文件拷贝到目标文件
	_, err = io.Copy(targetFile, file)
	if err != nil {
		log.Println(err)
		fmt.Fprintln(ctrl.Ctx.ResponseWriter, "{\"error\": 1, \"message\": \""+err.Error()+"\"}")
		return
	}
	outPath := lib.C_LOCAL_DOMAIN_Backslash() + strings.ReplaceAll(targetPath, "\\", "/")
	table := ctrl.GetSafeString("table")
	if table == "" {
		table = "images"
	}
	// 保存到数据库
	{
		mdl := domain.CmsAlbum{
			TableName_:   table,
			RecordID:     0,
			TypeID:       0,
			Title:        filename,
			ThumbPath:    outPath,
			OriginalPath: targetPath,
			FilePath:     outPath,
			FileSize:     header.Size,
			FileExt:      ext,
			LinkURL:      "",
			Click:        0,
			IsShow:       1,
			CreateID:     int32(GlobalAdminId),
		}
		if err := service.NewCmsAlbum().AlbumSave(&mdl); err != nil {
			logs.Error("KindEditorUpload AlbumSave", err.Error())
		}
	}
	// 返回上传结果
	result := "{\"error\": 0, \"url\": \"" + outPath + "\"}"
	fmt.Fprintln(ctrl.Ctx.ResponseWriter, result)
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
