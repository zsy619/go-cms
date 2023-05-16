function Json2Param(data) {
    const params = new URLSearchParams();
    for (const key in data) {
        params.append(key, data[key]);
    }
    const encodedData = params.toString();
    return encodedData;
}

/**
 * @description: ajax请求
 * @param {*} cache 是否缓存
 * @param {*} url 请求地址
 * @param {*} type 请求类型
 * @param {*} data 请求参数，json格式
 * @param {*} beforeSend 请求前执行
 * @param {*} success 请求成功执行
 * @param {*} complete 请求完成执行
 * @return {*}
 */
function AjaxRequest(cache, url, type, data, beforeSend, success, complete) {
    $.ajax({
        cache: cache,
        url: url,
        dataType: "json",
        type: type,
        beforeSend: function () {
            if (beforeSend) {
                beforeSend();
            }
        },
        success: function (res) {
            if (success) {
                success(res);
            }
        },
        error: function () {
        },
        complete: function () {
            if (complete) {
                complete();
            }
        }
    });
}

/**
 * @description: 站点定义
 * @return {*}
 */
function Site() {
    this.cache = true;
}

/**
 * @description: 获取站点信息
 * @param {*} data 参数，json格式
 * @param {*} beforeSend 请求前执行
 * @param {*} success 请求成功执行
 * @param {*} complete 请求完成执行
 * @return {*}
 */
Site.prototype.default = function (data, beforeSend, success, complete) {
    const params = Json2Param(data);
    var url = "/api/site/default?" + params;
    AjaxRequest(true, url, "get", data, beforeSend, success, complete);
}

/**
 * @description: 获取站点信息
 * @param {*} data 参数，json格式
 * @param {*} beforeSend 请求前执行
 * @param {*} success 请求成功执行
 * @param {*} complete 请求完成执行
 * @return {*}
 */
Site.prototype.find = function (data, beforeSend, success, complete) {
    const params = Json2Param(data);
    var url = "/api/site/find?" + params;
    AjaxRequest(true, url, "get", data, beforeSend, success, complete);
}

/**
 * @description: 获取站点频道
 * @param {*} data 参数，json格式
 * @param {*} beforeSend 请求前执行
 * @param {*} success 请求成功执行
 * @param {*} complete 请求完成执行
 * @return {*}
 */
Site.prototype.channelGet = function (data, beforeSend, success, complete) {
    const params = Json2Param(data);
    var url = "/api/channel/get?" + params;
    AjaxRequest(true, url, "get", data, beforeSend, success, complete);
}

/**
 * @description: 获取站点菜单
 * @param {*} data 参数，json格式
 * @param {*} beforeSend 请求前执行
 * @param {*} success 请求成功执行
 * @param {*} complete 请求完成执行
 * @return {*}
 */
Site.prototype.menu = function (data, beforeSend, success, complete) {
    const params = Json2Param(data);
    var url = "/api/site/menu?" + params;
    AjaxRequest(true, url, "get", data, beforeSend, success, complete);
}

/**
 * @description: 获取站点菜单
 * @param {*} data 参数，json格式
 * @param {*} beforeSend 请求前执行
 * @param {*} success 请求成功执行
 * @param {*} complete 请求完成执行
 * @return {*}
 */
Site.prototype.menuFlag = function (data, beforeSend, success, complete) {
    const params = Json2Param(data);
    var url = "/api/site/menu/flag?" + params;
    AjaxRequest(true, url, "get", data, beforeSend, success, complete);
}

/**
 * @description: 广告定义
 * @return {*}
 */
function Ads() {
    this.url = "/api/ads/get";
    this.cache = true;
}

/**
 * @description: 获取广告，按is_top升序排序
 * @param {*} data 参数，json格式
 * @param {*} beforeSend 请求前执行
 * @param {*} success 请求成功执行
 * @param {*} complete 请求完成执行
 * @return {*}
 */
Ads.prototype.get = function (data, beforeSend, success, complete) {
    const params = Json2Param(data);
    var url = this.url + "?" + params;
    AjaxRequest(this.cache, url, "get", data, beforeSend, success, complete);
}

/**
 * @description: 获取最新广告，按ads_id倒序排序
 * @param {*} data 参数，json格式
 * @param {*} beforeSend 请求前执行
 * @param {*} success 请求成功执行
 * @param {*} complete 请求完成执行
 * @return {*}
 */
Ads.prototype.getNew = function (data, beforeSend, success, complete) {
    const params = Json2Param(data);
    var url = this.url + "/new?" + params;
    AjaxRequest(this.cache, url, "get", data, beforeSend, success, complete);
}

/**
 * @description: 分页取广告，按is_top升序排序
 * @param {*} data 参数，json格式
 * @param {*} beforeSend 请求前执行
 * @param {*} success 请求成功执行
 * @param {*} complete 请求完成执行
 * @return {*}
 */
Ads.prototype.paginate = function (data, beforeSend, success, complete) {
    const params = Json2Param(data);
    var url = "/api/ads/paginate?" + params;
    AjaxRequest(this.cache, url, "get", data, beforeSend, success, complete);
}

/**
 * @description: 文章定义
 * @return {*}
 */
function Article() {
    this.url = "/api/article/get";
    this.urlNew = "/api/article/get/new";
    this.cache = true;
}

/**
 * @description: 获取文章
 * @param {*} data 参数，json格式
 * @param {*} beforeSend 请求前执行
 * @param {*} success 请求成功执行
 * @param {*} complete 请求完成执行
 * @return {*}
 */
Article.prototype.get = function (data, beforeSend, success, complete) {
    const params = Json2Param(data);
    var url = this.url + "?" + params;
    AjaxRequest(this.cache, url, "get", data, beforeSend, success, complete);
}

/**
 * @description: 获取最新文章
 * @param {*} data 参数，json格式
 * @param {*} beforeSend 请求前执行
 * @param {*} success 请求成功执行
 * @param {*} complete 请求完成执行
 * @return {*}
 */
Article.prototype.getNew = function (data, beforeSend, success, complete) {
    const params = Json2Param(data);
    var url = this.urlNew + "?" + params;
    AjaxRequest(this.cache, url, "get", data, beforeSend, success, complete);
}

/**
 * @description: 获取相册
 * @param {*} data 参数，json格式
 * @param {*} beforeSend 请求前执行
 * @param {*} success 请求成功执行
 * @param {*} complete 请求完成执行
 * @return {*}
 */
Article.prototype.album = function (data, beforeSend, success, complete) {
    const params = Json2Param(data);
    var url = "/api/article/album?" + params;
    AjaxRequest(this.cache, url, "get", data, beforeSend, success, complete);
}

/**
 * @description: 获取附件
 * @param {*} data 参数，json格式
 * @param {*} beforeSend 请求前执行
 * @param {*} success 请求成功执行
 * @param {*} complete 请求完成执行
 * @return {*}
 */
Article.prototype.attach = function (data, beforeSend, success, complete) {
    const params = Json2Param(data);
    var url = "/api/article/attach?" + params;
    AjaxRequest(this.cache, url, "get", data, beforeSend, success, complete);
}

/**
 * @description: 获取文章明细
 * @param {*} data 参数，json格式
 * @param {*} beforeSend 请求前执行
 * @param {*} success 请求成功执行
 * @param {*} complete 请求完成执行
 * @return {*}
 */
Article.prototype.article = function (data, beforeSend, success, complete) {
    const params = Json2Param(data);
    var url = "/api/article/article?" + params;
    AjaxRequest(this.cache, url, "get", data, beforeSend, success, complete);
}

/**
 * @description: 获取文章明细
 * @param {*} data 参数，json格式
 * @param {*} beforeSend 请求前执行
 * @param {*} success 请求成功执行
 * @param {*} complete 请求完成执行
 * @return {*}
 */
Article.prototype.find = function (data, beforeSend, success, complete) {
    const params = Json2Param(data);
    var url = "/api/article/find?" + params;
    AjaxRequest(this.cache, url, "get", data, beforeSend, success, complete);
}

/**
 * @description: 分页获取文章列表
 * @param {*} data 参数，json格式
 * @param {*} beforeSend 请求前执行
 * @param {*} success 请求成功执行
 * @param {*} complete 请求完成执行
 * @return {*}
 */
Article.prototype.paginate = function (data, beforeSend, success, complete) {
    const params = Json2Param(data);
    var url = "/api/article/paginate?" + params;
    AjaxRequest(this.cache, url, "get", data, beforeSend, success, complete);
}

/**
 * @description: 栏目定义
 * @return {*}
 */
function Category() {
    this.url = "/api/category/get";
    this.cache = true;
}

/**
 * @description: 获取频道下的栏目
 * @param {*} data 参数，json格式
 * @param {*} beforeSend 请求前执行
 * @param {*} success 请求成功执行
 * @param {*} complete 请求完成执行
 * @return {*}
 */
Category.prototype.get = function (data, beforeSend, success, complete) {
    const params = Json2Param(data);
    var url = "/api/category/get?" + params;
    AjaxRequest(this.cache, url, "get", data, beforeSend, success, complete);
}

/**
 * @description: 栏目路径
 * @param {*} data 参数，json格式
 * @param {*} beforeSend 请求前执行
 * @param {*} success 请求成功执行
 * @param {*} complete 请求完成执行
 * @return {*}
 */
Category.prototype.nav = function (data, beforeSend, success, complete) {
    const params = Json2Param(data);
    var url = "/api/category/nav?" + params;
    AjaxRequest(this.cache, url, "get", data, beforeSend, success, complete);
}

/**
 * @description: 获取栏目详情
 * @param {*} data 参数，json格式
 * @param {*} beforeSend 请求前执行
 * @param {*} success 请求成功执行
 * @param {*} complete 请求完成执行
 * @return {*}
 */
Category.prototype.find = function (data, beforeSend, success, complete) {
    const params = Json2Param(data);
    var url = "/api/category/find?" + params;
    AjaxRequest(this.cache, url, "get", data, beforeSend, success, complete);
}

function Link() {
    this.cache = true;
}

/**
 * @description: 获取链接，按is_top升序排序
 * @param {*} data 参数，json格式
 * @param {*} beforeSend 请求前执行
 * @param {*} success 请求成功执行
 * @param {*} complete 请求完成执行
 * @return {*}
 */
Link.prototype.get = function (data, beforeSend, success, complete) {
    const params = Json2Param(data);
    var url = "/api/link/get?" + params;
    AjaxRequest(this.cache, url, "get", data, beforeSend, success, complete);
}

/**
 * @description: 获取最新链接，按link_id倒序排序
 * @param {*} data 参数，json格式
 * @param {*} beforeSend 请求前执行
 * @param {*} success 请求成功执行
 * @param {*} complete 请求完成执行
 * @return {*}
 */
Link.prototype.getNew = function (data, beforeSend, success, complete) {
    const params = Json2Param(data);
    var url = "/api/link/get/new?" + params;
    AjaxRequest(this.cache, url, "get", data, beforeSend, success, complete);
}

/**
 * @description: 分页取链接，按is_top升序排序
 * @param {*} data 参数，json格式
 * @param {*} beforeSend 请求前执行
 * @param {*} success 请求成功执行
 * @param {*} complete 请求完成执行
 * @return {*}
 */
Link.prototype.paginate = function (data, beforeSend, success, complete) {
    const params = Json2Param(data);
    var url = "/api/link/paginate?" + params;
    AjaxRequest(this.cache, url, "get", data, beforeSend, success, complete);
}
