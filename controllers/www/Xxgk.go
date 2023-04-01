package www

// XxgkController 学校概况
type XxgkController struct{ CategoryBaseController }

// Xxjj 学校简介
// @router /xxgk/xxjj [get]
func (this *XxgkController) Xxjj() {
	this.displayCategory()
}

// Xxry 学校荣誉
// @router /xxgk/xxry [get]
func (this *XxgkController) Xxry() {
	this.displayCategory()
}

// Xyfg 校园风光
// @router /xxgk/xyfg [get]
func (this *XxgkController) Xyfg() {
	this.displayCategory()
}
