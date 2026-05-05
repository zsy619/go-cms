package admin

import (
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
)

type FileController struct{ BaseController }

// Index 文件管理首页
// @router /admin/file/index [get]
func (ctrl *FileController) Index() {
	ctrl.display()
}

func (ctrl *FileController) CreateRandomStr(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func (ctrl *FileController) GetFileDataVm() {
	runCommand, err := os.Executable()
	if err != nil {
		logs.Debug(err)
	}
	runPath := filepath.Dir(runCommand)
	logs.Debug("runPath---->", runPath)
	path := ctrl.GetSafeString("path")
	if path == "" {
		path = "./Uploads"
	}

	rs := struct {
		Images []map[string]interface{} `json:"images"`
		Count  int                      `json:"count"`
	}{}

	infos, _ := os.ReadDir(path)
	for _, info := range infos {
		logs.Debug("Name:%-30s IsDir:%v", info.Name(), info.IsDir())
		fs, err := filepath.Abs(info.Name())
		if err != nil {
			logs.Debug(err)
			continue
		}
		logs.Debug("fs---->", fs)
		if info.IsDir() {
			rs.Images = append(rs.Images, map[string]interface{}{
				"thumb": "",
				"name":  info.Name(),
				"type":  "dir",
				"path":  path[1:] + strings.Replace(fs, runPath, "", -1),
			})
		} else {
			rs.Images = append(rs.Images, map[string]interface{}{
				"thumb": info.Name(),
				"name":  info.Name(),
				"type":  filepath.Ext(info.Name())[1:],
				"path":  path[1:] + strings.Replace(fs, runPath, "", -1),
			})
		}
	}
	// rand.Shuffle(len(rs.Images), func(i, j int) {
	// 	rs.Images[i], rs.Images[j] = rs.Images[j], rs.Images[i]
	// })

	rs.Count = len(rs.Images)
	ctrl.Data["json"] = rs
	ctrl.ServeJSON()
}
