package www

// JyjxController 教育教学
type JyjxController struct{ CategoryBaseController }

// Zyjs 专业介绍
func (this *JyjxController) Zyjs() {
	this.displayCategory()
}

// Tsbx 特色班型
func (this *JyjxController) Tsbx() {
	this.displayCategory()
}

// Jxhd 教学活动
func (this *JyjxController) Jxhd() {
	this.displayCategory()
}

// Jsfc 教师风采
func (this *JyjxController) Jsfc() {
	this.displayCategory()
}
