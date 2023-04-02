package www

// XyzxController 校园资讯
type XyzxController struct{ CategoryBaseController }

// Xyxw 校园新闻
func (this *XyzxController) Xyxw() {
	this.displayCategory()
}

// Tzgg 通知公告
func (this *XyzxController) Tzgg() {
	this.displayCategory()
}
