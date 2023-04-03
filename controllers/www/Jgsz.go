package www

import (
	"haedu.gov.cn/cms/app/biz/bmodel"
	"haedu.gov.cn/cms/app/dal/model"
)

// JgszController 机构设置
type JgszController struct{ CategoryBaseController }

func (c *JgszController) Prepare() {
	c.CategoryBaseController.Prepare()
	article, err := c.ArticleArticle(c.ActionName, 0)
	if err != nil {
		article = &bmodel.ApiArticleOneModel{}
	}
	c.Data["article"] = article
}

// Xzbgs 校长办公室
func (c *JgszController) Xzbgs() {
	c.displayCategory()
}

// Dzbgs 党政办公室
func (c *JgszController) Dzbgs() {
	c.displayCategory()
}

// Xwbgs 校务办公室
func (c *JgszController) Xwbgs() {
	c.displayCategory()
}

// Rsc 人事处
func (c *JgszController) Rsc() {
	album, err := c.ArticleAlbum(c.ActionName, 0)
	if err != nil {
		album = []*model.CmsArticleAlbum{}
	}
	c.Data["album"] = album
	c.displayCategory()
}

// Jwc 教务处
func (c *JgszController) Jwc() {
	album, err := c.ArticleAlbum(c.ActionName, 0)
	if err != nil {
		album = []*model.CmsArticleAlbum{}
	}
	c.Data["album"] = album
	c.displayCategory()
}

// Xsc 学生处
func (c *JgszController) Xsc() {
	album, err := c.ArticleAlbum(c.ActionName, 0)
	if err != nil {
		album = []*model.CmsArticleAlbum{}
	}
	c.Data["album"] = album
	c.displayCategory()
}

// Tw 团委
func (c *JgszController) Tw() {
	album, err := c.ArticleAlbum(c.ActionName, 0)
	if err != nil {
		album = []*model.CmsArticleAlbum{}
	}
	c.Data["album"] = album
	c.displayCategory()
}

// Cwc 财务处
func (c *JgszController) Cwc() {
	c.displayCategory()
}

// Zsjybgs 招生就业办公室
func (c *JgszController) Zsjybgs() {
	c.displayCategory()
}

// Zwc 总务处
func (c *JgszController) Zwc() {
	c.displayCategory()
}

// Daglzx 档案管理中心
func (c *JgszController) Daglzx() {
	c.displayCategory()
}
