package www

import "haedu.gov.cn/cms/app/biz"

type ApiCacheController struct{ BaseController }

// @router /api/cache/clear [get]
func (this *ApiCacheController) Clear() {
	biz.CleanCahe()
	DefatulSite = nil
	this.JSONError("清除缓存成功")
}
