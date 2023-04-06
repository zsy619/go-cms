package admin

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
)

type UEditorController struct{ BaseController }

func (ue *UEditorController) Upload() {
	op := ue.Ctx.Request.Form.Get("action")

	switch op {
	case "showPost":

	case "config":
		file, err := os.Open("./static/admin/lib/ueditor-plus-v3.0.0/ueditor.config.js")
		if err != nil {
			fmt.Println(err.Error())
			fmt.Fprintf(ue.Ctx.ResponseWriter, "打开文件错误 : %v", err)
			return
		}
		defer file.Close()
		fd, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Fprintf(ue.Ctx.ResponseWriter, "读取文件失败 : %v", err)
			return
		}

		src := string(fd)
		re, _ := regexp.Compile(`\/\*[\S\s]+?\*\/`) // 匹配里面的注释
		src = re.ReplaceAllString(src, "")
		tt := []byte(src)
		var r interface{}
		err = json.Unmarshal(tt, &r) // 这个byte要解码
		if err != nil {
			fmt.Println(err.Error())
			fmt.Fprintf(ue.Ctx.ResponseWriter, "json decode failed %v", err)
			return
		}

		tt, err = json.Marshal(r)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Fprintf(ue.Ctx.ResponseWriter, "json encode failed %v", err)
			return
		}
		fmt.Fprint(ue.Ctx.ResponseWriter, string(tt))
	// 上传图片的功能
	case "uploadimage":
		err := ue.Ctx.Request.ParseForm()
		if err != nil {
			fmt.Println(err.Error())
			fmt.Fprintf(ue.Ctx.ResponseWriter, "uploadimage parseform fail : %v", err)
			return
		}
		// fmt.Println("打印所有的请求 : ", ue.Ctx.Request.PostForm)
		// fmt.Println("打印 upfile :", ue.Ctx.Request.PostForm.Get("upfile"))
		// 开始上传
		// 文件路径
		file, h, err := ue.Ctx.Request.FormFile("upfile")
		defer file.Close()
		if err != nil {
			fmt.Println(err.Error())
			ue.JSONError(err.Error())
		}
		ext := strings.ToLower(path.Ext(h.Filename))
		filename := generateFilename(ext)
		fmt.Println("file header: ", h, "\n detail-filename: ", h.Filename, " ;\n fileHeader:", h.Header)
		// 文件路径
		filePath := "Uploads/UEditor/" + time.Now().Format("2006/01/")
		err = os.MkdirAll(filePath, 0o777)
		if err != nil {
			ue.JSONError(err.Error())
		}
		// 文件名
		fileName := filePath + filename
		f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o666)
		ue.JSONError(err.Error())
		defer f.Close()
		io.Copy(f, file)
		ret_json := map[string]interface{}{
			"state":    "SUCCESS",
			"url":      "/" + fileName,
			"title":    h.Filename,
			"original": h.Filename,
			"type":     h.Header.Get("Content-Type"),
		}
		json, _ := json.Marshal(ret_json)
		fmt.Fprintf(ue.Ctx.ResponseWriter, string(json))

	default:
		fmt.Fprint(ue.Ctx.ResponseWriter, `{"msg":"请求地址错误"}`)
	}
}
