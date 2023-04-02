package www

// LxwmController 联系我们
type LxwmController struct{ CategoryBaseController }

// Bgsdh 办公室电话
func (this *LxwmController) Bgsdh() {
	this.displayCategory()
}

// Xzyx 校长邮箱
func (this *LxwmController) Xzyx() {
	this.displayCategory()
}

// Dlwz 地理位置
func (this *LxwmController) Dlwz() {
	this.displayCategory()
}
