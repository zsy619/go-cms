package admin

import (
	"fmt"
	"haedu.gov.cn/cms/controllers/admin/vmodel"
	"time"

	"haedu.gov.cn/cms/app/biz"
	"haedu.gov.cn/cms/app/dal/query"
)

type IndexController struct {
	BaseController
}

func (c *IndexController) Index() {
	c.Data["userName"] = GlobalRealName
	c.displayNoLayout()
}

func (c *IndexController) Welcome() {
	c.Data["roleId"] = GlobalRoleId
	c.Data["roleType"] = GlobalRoleType
	noticeList, _ := biz.NewCmsAdminNotice().NoticeShow()
	c.Data["noticeList"] = noticeList
	c.display()
}

func (c *IndexController) Count() {
	result := struct {
		SiteCount     int64 `json:"site_count"`
		ChannelCount  int64 `json:"channel_count"`
		CategoryCount int64 `json:"category_count"`
		ArticleCount  int64 `json:"article_count"`
	}{}

	_, siteDo := query.CmsSiteDo()
	siteCount, _ := siteDo.Count()

	_, channelDo := query.CmsSiteChannelDo()
	channelCount, _ := channelDo.Count()

	categoryMdl, categoryDo := query.CmsArticleCategoryDo()
	categoryCount, _ := categoryDo.Where(categoryMdl.Status.Eq(2)).Count()

	articleMdl, articleDo := query.CmsArticleDo()
	articleCount, _ := articleDo.Where(articleMdl.Status.Eq(2)).Count()

	result.SiteCount = siteCount
	result.ChannelCount = channelCount
	result.CategoryCount = categoryCount
	result.ArticleCount = articleCount

	c.JSONSuccess("success", result)
}

func (c *IndexController) UserPassword() {
	c.display()
}

// UserPasswordSave 修改密码
// @router /admin/index/UserPasswordSave [post]
func (c *IndexController) UserPasswordSave() {
	oldPassword := c.GetString("old_password")
	newPassword := c.GetString("new_password")
	confirmPassword := c.GetString("again_password")
	fmt.Println(oldPassword, newPassword, confirmPassword)

	if newPassword != confirmPassword {
		c.JSONError("两次输入的密码不一致")
	}
	// 检查密码是否符合规则
	if psErr := CheckPasswordRole(newPassword); psErr != nil {
		c.JSONError(psErr.Error())
	}
	err := biz.NewCmsAdmin().ModifyPassword(GlobalAdminId, oldPassword, newPassword)
	if err != nil {
		c.JSONError(err.Error())
	}

	c.JSONSuccess("修改成功", nil)
}

func (c *IndexController) UserSetting() {
	c.display()
}

// ReportFormsGet 获取首页报表数据
func (c *IndexController) ReportFormsGet() {
	currentDate := time.Now()
	var times []time.Time
	var showTimes []string
	for i := 0; i < 7; i++ {
		duration, _ := time.ParseDuration(fmt.Sprintf("%dh", -(7-i)*24))
		yesTime := currentDate.Add(duration)
		yesTime = time.Date(yesTime.Year(), yesTime.Month(), yesTime.Day(), 0, 0, 0, 0, time.Local)
		times = append(times, yesTime)
		showTimes = append(showTimes, yesTime.Format("01-02"))
	}
	linkCounts := biz.NewCmsLink().FindByDate(times...)
	adsCounts := biz.NewCmsAds().FindByDate(times...)
	articleCounts := biz.NewCmsArticle().FindByDate(times...)

	mdl := vmodel.ReportFormsModel{}
	mdl.Dates = showTimes
	item := vmodel.ReportFormsItemModel{}
	item.Name = "链接"
	item.Data = linkCounts
	item.Type = "line"
	mdl.Items = append(mdl.Items, item)

	item = vmodel.ReportFormsItemModel{}
	item.Name = "广告"
	item.Data = adsCounts
	item.Type = "line"
	mdl.Items = append(mdl.Items, item)

	item = vmodel.ReportFormsItemModel{}
	item.Name = "内容"
	item.Data = articleCounts
	item.Type = "line"
	mdl.Items = append(mdl.Items, item)

	c.JSONSuccess("", mdl)
}
