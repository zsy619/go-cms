package www

import "haedu.gov.cn/cms/app/cms/service"

type ApiCacheController struct{ BaseController }

// @router /api/cache/clear [get]
func (ctrl *ApiCacheController) Clear() {
	service.CleanCahe()
	DefaultSite = nil
	ctrl.JSONError("清除缓存成功")
}
