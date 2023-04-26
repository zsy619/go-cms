function Json2Param(data) {
    const params = new URLSearchParams();
    for (const key in data) {
        params.append(key, data[key]);
    }
    const encodedData = params.toString();
    return encodedData;
}

/**
 * @description: 站点定义
 * @return {*}
 */
function Site() {
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
 * @description: 获取广告
 * @param {*} data 参数，json格式
 * @param {*} beforeSend 请求前执行
 * @param {*} success 请求成功执行
 * @param {*} complete 请求完成执行
 * @return {*}
 */
Ads.prototype.get = function (data, beforeSend, success, complete) {
    const params = Json2Param(data);
    var url = this.url + "?" + params;
    $.ajax({
        cache: this.cache,
        url: url,
        dataType: "json",
        type: "get",
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
    $.ajax({
        cache: this.cache,
        url: url,
        dataType: "json",
        type: "get",
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