package www

import (
	"haedu.gov.cn/cms/app/biz/bizmodel"
	"haedu.gov.cn/cms/app/dal/model"
)

// JgszController 机构设置
type JgszController struct{ CategoryBaseController }

func (this *JgszController) Prepare() {
	this.CategoryBaseController.Prepare()
	article, err := this.ArticleArticle(this.ActionName, 0)
	if err != nil {
		article = &bizmodel.ApiArticleOneModel{}
	}
	this.Data["article"] = article
}

// Xzbgs 校长办公室
func (this *JgszController) Xzbgs() {
	this.displayCategory()
}

// Dzbgs 党政办公室
func (this *JgszController) Dzbgs() {
	this.displayCategory()
}

// Xwbgs 校务办公室
func (this *JgszController) Xwbgs() {
	this.displayCategory()
}

// Rsc 人事处
func (this *JgszController) Rsc() {
	album, err := this.ArticleAlbum(this.ActionName, 0, 0)
	if err != nil {
		album = []*model.CmsAlbum{}
	}
	this.Data["album"] = album
	this.displayCategory()
}

// Jwc 教务处
func (this *JgszController) Jwc() {
	album, err := this.ArticleAlbum(this.ActionName, 0, 0)
	if err != nil {
		album = []*model.CmsAlbum{}
	}
	this.Data["album"] = album
	this.displayCategory()
}

// Xsc 学生处
func (this *JgszController) Xsc() {
	album, err := this.ArticleAlbum(this.ActionName, 0, 0)
	if err != nil {
		album = []*model.CmsAlbum{}
	}
	this.Data["album"] = album
	this.displayCategory()
}

// Tw 团委
func (this *JgszController) Tw() {
	album, err := this.ArticleAlbum(this.ActionName, 0, 0)
	if err != nil {
		album = []*model.CmsAlbum{}
	}
	this.Data["album"] = album
	this.displayCategory()
}

// Cwc 财务处
func (this *JgszController) Cwc() {
	this.displayCategory()
}

// Zsjybgs 招生就业办公室
func (this *JgszController) Zsjybgs() {
	this.displayCategory()
}

// Zwc 总务处
func (this *JgszController) Zwc() {
	this.displayCategory()
}

// Daglzx 档案管理中心
func (this *JgszController) Daglzx() {
	this.displayCategory()
}
