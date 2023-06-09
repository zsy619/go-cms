package admin

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type FileController struct{ BaseController }

// Index 文件管理首页
// @router /admin/file/index [get]
func (this *FileController) Index() {
	this.display()
}

func (this *FileController) CreateRandomStr(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func (this *FileController) GetFileDataVm() {
	runCommand, err := os.Executable()
	if err != nil {
		fmt.Println(err)
	}
	runPath := filepath.Dir(runCommand)
	fmt.Println("runPath---->", runPath)
	path := this.GetString("path")
	if path == "" {
		path = "./Uploads"
	}

	rs := struct {
		Images []map[string]interface{} `json:"images"`
		Count  int                      `json:"count"`
	}{}

	infos, _ := ioutil.ReadDir(path)
	for _, info := range infos {
		fmt.Printf("Name:%-30s    Size:%-10d字节    Modtime:%s    IsDir:%v\n", info.Name(), info.Size(), info.ModTime().Format("2006-01-02 03:04:05"), info.IsDir())
		fs, err := filepath.Abs(info.Name())
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("fs---->", fs)
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
	this.Data["json"] = rs
	this.ServeJSON()
}
