package admin

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/app/lib"
)

var UEditorConfig = map[string]interface{}{
	// 上传图片配置项
	"imageActionName":     "image",                                           /* 执行上传图片的action名称 */
	"imageFieldName":      "file",                                            /* 提交的图片表单名称 */
	"imageMaxSize":        1024 * 1024 * 10,                                  /* 上传大小限制，单位B */
	"imageAllowFiles":     []string{".png", ".jpg", ".jpeg", ".gif", ".bmp"}, /* 上传图片格式显示 */
	"imageCompressEnable": true,                                              /* 是否压缩图片,默认是true */
	"imageCompressBorder": 5000,                                              /* 图片压缩最长边限制 */
	"imageInsertAlign":    "none",                                            /* 插入的图片浮动方式 */
	"imageUrlPrefix":      "",                                                /* 图片访问路径前缀 */
	"imagePathFormat":     "/Uploads/UEditor/image/",                         /* 上传保存路径,可以自定义保存路径和文件名格式 */

	/* 涂鸦图片上传配置项 */
	"scrawlActionName":  "crawl",                   /* 执行上传涂鸦的action名称 */
	"scrawlFieldName":   "file",                    /* 提交的图片表单名称 */
	"scrawlPathFormat":  "/Uploads/UEditor/crawl/", /* 上传保存路径,可以自定义保存路径和文件名格式 */
	"scrawlMaxSize":     1024 * 1024 * 10,          /* 上传大小限制，单位B */
	"scrawlUrlPrefix":   "",                        /* 图片访问路径前缀 */
	"scrawlInsertAlign": "none",

	/* 截图工具上传 */
	"snapscreenActionName":  "snap",                   /* 执行上传截图的action名称 */
	"snapscreenPathFormat":  "/Uploads/UEditor/snap/", /* 上传保存路径,可以自定义保存路径和文件名格式 */
	"snapscreenUrlPrefix":   "",                       /* 图片访问路径前缀 */
	"snapscreenInsertAlign": "none",                   /* 插入的图片浮动方式 */

	/* 抓取远程图片配置 */
	"catcherLocalDomain": []string{
		"127.0.0.1",
		"localhost",
	},
	"catcherActionName": "catch",                                           /* 执行抓取远程图片的action名称 */
	"catcherFieldName":  "source",                                          /* 提交的图片列表表单名称 */
	"catcherPathFormat": "/Uploads/UEditor/catch/",                         /* 上传保存路径,可以自定义保存路径和文件名格式 */
	"catcherUrlPrefix":  "",                                                /* 图片访问路径前缀 */
	"catcherMaxSize":    1024 * 1024 * 10,                                  /* 上传大小限制，单位B */
	"catcherAllowFiles": []string{".png", ".jpg", ".jpeg", ".gif", ".bmp"}, /* 抓取图片格式显示 */

	/* 上传视频配置 */
	"videoActionName": "video",                                                                                                                                            /* 执行上传视频的action名称 */
	"videoFieldName":  "file",                                                                                                                                             /* 提交的视频表单名称 */
	"videoPathFormat": "/Uploads/UEditor/video/",                                                                                                                          /* 上传保存路径,可以自定义保存路径和文件名格式 */
	"videoUrlPrefix":  "",                                                                                                                                                 /* 视频访问路径前缀 */
	"videoMaxSize":    1024 * 1024 * 100,                                                                                                                                  /* 上传大小限制，单位B，默认100MB */
	"videoAllowFiles": []string{".flv", ".swf", ".mkv", ".avi", ".rm", ".rmvb", ".mpeg", ".mpg", ".ogg", ".ogv", ".mov", ".wmv", ".mp4", ".webm", ".mp3", ".wav", ".mid"}, /* 上传视频格式显示 */

	/* 上传文件配置 */
	"fileActionName": "file",                                                                      /* controller里,执行上传视频的action名称 */
	"fileFieldName":  "file",                                                                      /* 提交的文件表单名称 */
	"filePathFormat": "/Uploads/UEditor/file/",                                                    /* 上传保存路径,可以自定义保存路径和文件名格式 */
	"fileUrlPrefix":  "",                                                                          /* 文件访问路径前缀 */
	"fileMaxSize":    1024 * 1024 * 100,                                                           /* 上传大小限制，单位B，默认50MB */
	"fileAllowFiles": []string{".zip", ".pdf", ".doc", ".docx", ".ppt", ".pptx", ".xls", ".xlsx"}, /* 上传文件格式显示 */

	/* 列出图片 */
	"imageManagerActionName":  "listImage",                                       /* 执行图片管理的action名称 */
	"imageManagerListPath":    "/Uploads/UEditor/image/",                         /* 指定要列出图片的目录 */
	"imageManagerListSize":    20,                                                /* 每次列出文件数量 */
	"imageManagerUrlPrefix":   "",                                                /* 图片访问路径前缀 */
	"imageManagerInsertAlign": "none",                                            /* 插入的图片浮动方式 */
	"imageManagerAllowFiles":  []string{".png", ".jpg", ".jpeg", ".gif", ".bmp"}, /* 列出的文件类型 */

	/* 列出指定目录下的文件 */
	"fileManagerActionName": "listFile",                                                                                                                                                                                                                                                                                                                                  /* 执行文件管理的action名称 */
	"fileManagerListPath":   "/Uploads/UEditor/file/",                                                                                                                                                                                                                                                                                                                    /* 指定要列出文件的目录 */
	"fileManagerUrlPrefix":  "",                                                                                                                                                                                                                                                                                                                                          /* 文件访问路径前缀 */
	"fileManagerListSize":   20,                                                                                                                                                                                                                                                                                                                                          /* 每次列出文件数量 */
	"fileManagerAllowFiles": []string{".png", ".jpg", ".jpeg", ".gif", ".bmp", ".flv", ".swf", ".mkv", ".avi", ".rm", ".rmvb", ".mpeg", ".mpg", ".ogg", ".ogv", ".mov", ".wmv", ".mp4", ".webm", ".mp3", ".wav", ".mid", ".rar", ".zip", ".tar", ".gz", ".7z", ".bz2", ".cab", ".iso", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".pdf", ".txt", ".md", ".xml"}, /* 列出的文件类型 */

	/* 公式渲染 */
	"formulaConfig": map[string]interface{}{
		"imageUrlTemplate": "https://latex.codecogs.com/svg.image?{}",
	},
}

type UEditorController struct{ BaseController }

func (ue *UEditorController) Handle() {
	op := ue.Ctx.Request.Form.Get("action")
	if op == "" {
		op = ue.Ctx.Input.Query("action")
	}
	fmt.Println("upload-------->", op)
	switch op {
	case "showPost":
		body := ue.Ctx.Input.RequestBody
		ue.Data["json"] = string(body)
		ue.ServeJSON()
	// case "config":
	// 	// 返回配置信息
	// 	file, err := os.Open("./static/admin/lib/ueditor-plus-v3.0.0/ueditor.config.js")
	// 	if err != nil {
	// 		fmt.Println(err.Error())
	// 		fmt.Fprintf(ue.Ctx.ResponseWriter, "打开文件错误 : %v", err)
	// 		return
	// 	}
	// 	defer file.Close()
	// 	fd, err := ioutil.ReadAll(file)
	// 	if err != nil {
	// 		fmt.Println(err.Error())
	// 		fmt.Fprintf(ue.Ctx.ResponseWriter, "读取文件失败 : %v", err)
	// 		return
	// 	}

	// 	src := string(fd)
	// 	re, _ := regexp.Compile(`\/\*[\S\s]+?\*\/`) // 匹配里面的注释
	// 	src = re.ReplaceAllString(src, "")
	// 	tt := []byte(src)
	// 	var r interface{}
	// 	err = json.Unmarshal(tt, &r) // 这个byte要解码
	// 	if err != nil {
	// 		fmt.Println(err.Error())
	// 		fmt.Fprintf(ue.Ctx.ResponseWriter, "json decode failed %v", err)
	// 		return
	// 	}

	// 	tt, err = json.Marshal(r)
	// 	if err != nil {
	// 		fmt.Println(err.Error())
	// 		fmt.Fprintf(ue.Ctx.ResponseWriter, "json encode failed %v", err)
	// 		return
	// 	}
	// 	fmt.Fprint(ue.Ctx.ResponseWriter, string(tt))
	case "image":
		ue.UploadFileX(op, UEditorConfig["imageFieldName"].(string), UEditorConfig["imageAllowFiles"].([]string), "不支持的图片格式", UEditorConfig["imagePathFormat"].(string)[1:])
	case "listImage":
		// 列出图片
		ue.ListFileX("."+UEditorConfig["imagePathFormat"].(string), UEditorConfig["imageAllowFiles"].([]string))
	case "video":
		ue.UploadFileX(op, UEditorConfig["videoFieldName"].(string), UEditorConfig["videoAllowFiles"].([]string), "不支持的视频格式", UEditorConfig["videoPathFormat"].(string)[1:])
	case "file":
		ue.UploadFileX(op, UEditorConfig["fileFieldName"].(string), UEditorConfig["fileAllowFiles"].([]string), "不支持的文件格式", UEditorConfig["filePathFormat"].(string)[1:])
	case "listFile":
		// 列出文件
		ue.ListFileX("."+UEditorConfig["filePathFormat"].(string), UEditorConfig["fileAllowFiles"].([]string))
	case "crawl":
		// 涂鸦上传
	case "catch":
		// 抓取远程图片
	default:
		ue.Data["json"] = UEditorConfig
		ue.ServeJSON()
	}
}

func (ue *UEditorController) UploadFileX(op string, fieldName string, exts []string, extsMsg, filePath string) {
	file, h, err := ue.Ctx.Request.FormFile("file")
	if err != nil {
		fmt.Println(err.Error())
		ue.JSONError(err.Error())
	}
	defer file.Close()
	if !ue.isExtFile(h.Filename, exts) {
		ue.JSONError(extsMsg)
	}
	ext := strings.ToLower(path.Ext(h.Filename))
	filename := generateFilename(ext)
	// 文件路径
	err = os.MkdirAll(filePath, 0o777)
	if err != nil {
		ue.JSONError(err.Error())
	}
	// 文件名
	fileName := filePath + filename
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o666)
	if err != nil {
		ue.JSONError(err.Error())
	}
	io.Copy(f, file)
	ret_json := map[string]interface{}{
		"state":    "SUCCESS",
		"url":      lib.C_LOCAL_DOMAIN_Backslash() + fileName,
		"title":    h.Filename,
		"original": h.Filename,
		"type":     h.Header.Get("Content-Type"),
	}

	// 保存到数据库
	{
		if op == "images" {
			mdl := model.CmsAlbum{
				TableName_:   op,
				RecordID:     0,
				TypeID:       0,
				Title:        h.Filename,
				ThumbPath:    lib.C_LOCAL_DOMAIN_Backslash() + fileName,
				OriginalPath: "/" + fileName,
				FilePath:     lib.C_LOCAL_DOMAIN_Backslash() + fileName,
				FileSize:     h.Size,
				FileExt:      ext,
				LinkURL:      "",
				Click:        0,
				IsShow:       1,
				CreateID:     int32(GlobalAdminId),
			}
			if err := biz.NewCmsAlbum().AlbumSave(&mdl); err != nil {
				logs.Error("UploadFileX AlbumSave", err.Error())
			}
		} else {
			mdl := model.CmsAttach{
				TableName_:   op,
				RecordID:     0,
				TypeID:       0,
				Title:        filename,
				FilePath:     lib.C_LOCAL_DOMAIN_Backslash() + fileName,
				OriginalPath: "/" + fileName,
				FileSize:     h.Size,
				FileExt:      ext,
				Click:        0,
				IsShow:       1,
				CreateID:     int32(GlobalAdminId),
			}
			if err := biz.NewCmsAttach().AttachSave(&mdl); err != nil {
				logs.Error("UploadFileX AttachSave", err.Error())
			}
		}
	}

	json, _ := json.Marshal(ret_json)
	fmt.Fprintf(ue.Ctx.ResponseWriter, string(json))
}

func (ue *UEditorController) ListFileX(dir string, exts []string) {
	result := struct {
		State string `json:"state"`
		List  []struct {
			Url   string `json:"url"`
			Mtime int64  `json:"mtime"`
		} `json:"list"`
		Start int `json:"start"`
		Total int `json:"total"`
	}{}
	pageSize, _ := ue.GetInt("siez", 20) // 每页大小

	startIndex, _ := ue.GetInt("start", 0)
	endIndex := startIndex + pageSize
	count := 0
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && ue.isExtFile(path, exts) {
			if count >= startIndex && count < endIndex {
				fmt.Println(path)
				result.List = append(result.List, struct {
					Url   string `json:"url"`
					Mtime int64  `json:"mtime"`
				}{
					Url:   "/" + path,
					Mtime: info.ModTime().Unix(),
				})
			}
			count++
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	result.State = "SUCCESS"
	result.Start = startIndex
	result.Total = count
	ue.Data["json"] = result
	ue.ServeJSON()
}

func (ue *UEditorController) isExtFile(path string, exts []string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	for _, e := range exts {
		if ext == e {
			return true
		}
	}
	return false
}
