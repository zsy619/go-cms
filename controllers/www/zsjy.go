package www

// 招生就业
type ZsjyController struct{ BaseController }

func (this *ZsjyController) Prepare() {
	this.BaseController.Prepare()
	category, _ := this.CategoryOne(0, this.ActionName)
	this.Data["category"] = category
}

// 招生简章
// @router zsjy/zsjz [get]
func (this *ZsjyController) Zsjz() {
	this.displayCategory()
}

// 线上报名
// @router zsjy/xsbm [get]
func (this *ZsjyController) Xsbm() {
	this.displayCategory()
}

// 线上报名
// @router zsjy/jyfw [get]
func (this *ZsjyController) Jyfw() {
	this.displayCategory()
}
