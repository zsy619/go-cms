package admin

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/app/lib"
	"haedu.gov.cn/cms/controllers/admin/vmodel"
	"haedu.gov.cn/cms/controllers/www"
	"haedu.gov.cn/tools/xio"
	"haedu.gov.cn/tools/xjson"
)

type WeixinMpVerifyController struct{ BaseController }

// @router /admin/weixin/verify [get]
func (this *WeixinMpVerifyController) Index() {
	list, _, _ := biz.NewWeixinAccount().AccountPaginate(1, 9999, "", -1)
	this.Data["accountList"] = list
	this.display()
}

// @router /admin/weixin/verify/list [get]
func (c *WeixinMpVerifyController) List() {
	accountId, _ := c.GetInt64("account_id")
	list, err := biz.NewWeixinMpVerify().Get(accountId)
	if err != nil {
		logs.Error("List", err.Error())
	}
	c.JSONPageSuccess(list, int64(len(list)))
}

// @router /admin/weixin/verify/edit [get]
func (this *WeixinMpVerifyController) Edit() {
	if this.IsPost() {
		mdl := model.WeixinMpVerify{}
		if err := this.ParseForm(&mdl); err != nil {
			logs.Error("Edit", err.Error())
			this.JSONError(err.Error())
		}
		if mdl.VerifyID <= 0 {
			mdl.CreateID = int32(GlobalAdminId)
		} else {
			mdl.UpdateID = int32(GlobalAdminId)
		}
		if err := biz.NewWeixinMpVerify().Save(&mdl); err != nil {
			logs.Error("Edit", err.Error())
			this.JSONError(err.Error())
			return
		}
		this.JSONSuccess("保存成功", nil)
	}
	{
		list, _, _ := biz.NewWeixinAccount().AccountPaginate(1, 9999, "", -1)
		this.Data["accountList"] = list
	}
	{
		verifyId, _ := this.GetInt64("verify_id")
		this.Data["verify_id"] = verifyId
		finder, err := biz.NewWeixinMpVerify().Find(verifyId)
		if finder == nil || err != nil {
			finder = &model.WeixinMpVerify{
				SortID: 99,
				Path:   "/",
			}
		}
		this.Data["mdl"] = finder
	}
	this.display()
}

// @router /admin/weixin/verify/savesortid [post]
func (this *WeixinMpVerifyController) SaveSortId() {
	mdls := []vmodel.Verify_SaveSortIdModel{}
	data := this.Ctx.Input.RequestBody
	fmt.Println("SaveSortId", string(data))
	if err := xjson.Unmarshal(this.Ctx.Input.RequestBody, &mdls); err != nil {
		logs.Error("SaveSortId", err.Error())
		this.JSONError(err.Error())
		return
	}
	for _, mdl := range mdls {
		if err := biz.NewWeixinMpVerify().SaveSortId(mdl.VerifyId, int32(mdl.SortId)); err != nil {
			logs.Error("SaveSortId", err.Error())
			this.JSONError(err.Error())
			return
		}
	}
	this.JSONSuccess("保存成功", nil)
}

// @router /admin/weixin/verify/destory [post]
func (c *WeixinMpVerifyController) Destory() {
	verifyId, _ := c.GetInt64("verify_id")
	if err := biz.NewWeixinMpVerify().Destory(verifyId); err != nil {
		logs.Error("Destory", err.Error())
		c.JSONError(err.Error())
		return
	}
	c.JSONSuccess("删除成功", nil)
}

// @router admin/weixin/verify/Upload [post]
func (c *WeixinMpVerifyController) Upload() {
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
	file, head, err := c.GetFile("file")
	if err != nil {
		result.Code = 1
		result.Msg = "文件上传失败"
		c.Data["json"] = result
		c.ServeJSON()
		return
	}
	defer file.Close()
	if isAllow := xio.IsAllowFile(head.Filename); !isAllow {
		result.Code = 1
		result.Msg = "不支持的文件类型"
		c.Data["json"] = result
		c.ServeJSON()
		return
	}

	result.File.Name1 = head.Filename
	result.File.Size = head.Size

	ext := strings.ToLower(path.Ext(head.Filename))

	// 生成唯一的文件名
	filename := head.Filename
	pth := c.GetSafeString("path")
	table := c.GetSafeString("table")
	if pth == "" {
		pth = "weixin/verify"
	}
	if table == "" {
		table = "mp_verify"
	}
	// 保存上传的文件到指定目录
	uploadDir := "Uploads/" + pth + "/"
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
	result.File.Name2 = filename
	result.File.Url1 = "/" + uploadDir + filename
	result.File.Url2 = lib.C_LOCAL_DOMAIN_Backslash() + uploadDir + filename

	// 保存到数据库
	{
		mdl := model.CmsAttach{
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
		if err := biz.NewCmsAttach().AttachSave(&mdl); err != nil {
			logs.Error("Upload AttachSave", err.Error())
		}
	}

	c.Data["json"] = result
	c.ServeJSON()
}

/**
 * @description: 刷新缓存
 * @return {*}
 */
// @router /admin/weixin/verify/refrshcache [post]
func (c *WeixinMpVerifyController) RefrshCache() {
	biz.NewWeixinMpVerify().RefeshCache()
	www.InitWechatMpVerifyRouter()
	c.JSONSuccess("重置成功", nil)
}
