package admin

import (
	"fmt"
	"time"

	"haedu.gov.cn/cms/app/cms/mapper"
	"haedu.gov.cn/cms/app/cms/service"
	"haedu.gov.cn/cms/controllers/admin/vmodel"
)

type IndexController struct{ BaseController }

func (ctrl *IndexController) Index() {
	ctrl.Data["userName"] = GlobalRealName
	ctrl.displayNoLayout()
}

func (ctrl *IndexController) Welcome() {
	ctrl.Data["roleId"] = GlobalRoleId
	ctrl.Data["roleType"] = GlobalRoleType
	noticeList, _ := service.NewCmsAdminNotice().NoticeShow()
	ctrl.Data["noticeList"] = noticeList
	ctrl.display()
}

func (ctrl *IndexController) Count() {
	result := struct {
		SiteCount     int64 `json:"site_count"`
		ChannelCount  int64 `json:"channel_count"`
		CategoryCount int64 `json:"category_count"`
		ArticleCount  int64 `json:"article_count"`
	}{}

	_, siteDo := mapper.CmsSiteDo()
	siteCount, _ := siteDo.Count()

	_, channelDo := mapper.CmsSiteChannelDo()
	channelCount, _ := channelDo.Count()

	categoryMdl, categoryDo := mapper.CmsArticleCategoryDo()
	categoryCount, _ := categoryDo.Where(categoryMdl.Status.Eq(2)).Count()

	articleMdl, articleDo := mapper.CmsArticleDo()
	articleCount, _ := articleDo.Where(articleMdl.Status.Eq(2)).Count()

	result.SiteCount = siteCount
	result.ChannelCount = channelCount
	result.CategoryCount = categoryCount
	result.ArticleCount = articleCount

	ctrl.JSONSuccess("success", result)
}

func (ctrl *IndexController) UserPassword() {
	ctrl.display()
}

// UserPasswordSave 修改密码
// @router /admin/index/UserPasswordSave [post]
func (ctrl *IndexController) UserPasswordSave() {
	oldPassword := ctrl.GetSafeString("old_password")
	newPassword := ctrl.GetSafeString("new_password")
	confirmPassword := ctrl.GetSafeString("again_password")
	fmt.Println(oldPassword, newPassword, confirmPassword)

	if newPassword != confirmPassword {
		ctrl.JSONError("两次输入的密码不一致")
	}
	// 检查密码是否符合规则
	if psErr := CheckPasswordRole(newPassword); psErr != nil {
		ctrl.JSONError(psErr.Error())
	}
	err := service.NewCmsAdmin().ModifyPassword(GlobalAdminId, oldPassword, newPassword)
	if err != nil {
		ctrl.JSONError(err.Error())
	}

	ctrl.JSONSuccess("修改成功", nil)
}

func (ctrl *IndexController) UserSetting() {
	ctrl.display()
}

// ReportFormsGet 获取首页报表数据
func (ctrl *IndexController) ReportFormsGet() {
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
	linkCounts := service.NewCmsLink().FindByDate(times...)
	adsCounts := service.NewCmsAds().FindByDate(times...)
	articleCounts := service.NewCmsArticle().FindByDate(times...)

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

	ctrl.JSONSuccess("", mdl)
}
