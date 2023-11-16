package www

import "haedu.gov.cn/cms/app/biz"

type ApiCacheController struct{ BaseController }

// @router /api/cache/clear [get]
func (ctrl *ApiCacheController) Clear() {
	biz.CleanCahe()
	DefaultSite = nil
	ctrl.JSONError("清除缓存成功")
}
