package admin

import (
	"fmt"
	"time"

	"github.com/beego/beego/v2/core/logs"

	"haedu.gov.cn/cms/app/cms/mapper"
	"haedu.gov.cn/cms/app/cms/service"
	"haedu.gov.cn/cms/controllers/admin/vmodel"
)

// IndexController 后台首页控制器
type IndexController struct{ BaseController }

// Index 后台首页
func (ctrl *IndexController) Index() {
	ctrl.Data["userName"] = GlobalRealName
	ctrl.displayNoLayout()
}

// Welcome 欢迎页面
func (ctrl *IndexController) Welcome() {
	ctrl.Data["roleId"] = GlobalRoleId
	ctrl.Data["roleType"] = GlobalRoleType
	noticeList, _ := service.NewCmsAdminNotice().NoticeShow()
	ctrl.Data["noticeList"] = noticeList
	ctrl.display()
}

// Count 获取统计数据(站点、频道、栏目、文章数量)
// @return JSON响应包含四个统计数据
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

// UserPassword 修改密码页面
func (ctrl *IndexController) UserPassword() {
	ctrl.display()
}

// UserPasswordSave 修改密码处理
// @router /admin/index/UserPasswordSave [post]
func (ctrl *IndexController) UserPasswordSave() {
	oldPassword := ctrl.GetSafeString("old_password")
	newPassword := ctrl.GetSafeString("new_password")
	confirmPassword := ctrl.GetSafeString("again_password")
	logs.Debug("修改密码请求: oldPassword=%s, newPassword=%s, confirmPassword=%s", oldPassword, newPassword, confirmPassword)

	if newPassword != confirmPassword {
		ctrl.JSONError("两次输入的密码不一致")
	}
	if psErr := CheckPasswordRole(newPassword); psErr != nil {
		ctrl.JSONError(psErr.Error())
	}
	err := service.NewCmsAdmin().ModifyPassword(GlobalAdminId, oldPassword, newPassword)
	if err != nil {
		logs.Error("修改密码失败: adminId=%d, error=%v", GlobalAdminId, err)
		ctrl.JSONError(err.Error())
	}

	logs.Info("用户修改密码成功: adminId=%d", GlobalAdminId)
	ctrl.JSONSuccess("修改成功", nil)
}

// UserSetting 用户设置页面
func (ctrl *IndexController) UserSetting() {
	ctrl.display()
}

// ReportFormsGet 获取首页报表数据(最近7天的链接、广告、文章创建数量统计)
// @return JSON响应包含日期和各类数据统计
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