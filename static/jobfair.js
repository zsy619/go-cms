function Json2Param(data) {
    const params = new URLSearchParams();
    for (const key in data) {
        params.append(key, data[key]);
    }
    const encodedData = params.toString();
    return encodedData;
}

/**
 * @description: 获取就业站点url
 * @return {*}
 */
function getJobfairSite() {
    // return "http://192.168.80.238:8888/";
    return "http://192.168.90.117:8888/";
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

function Job() {
    this.urlJobHotList = "/subsite/JobInterface/GetHotJobList";
    this.urlJobList = "/subsite/JobInterface/GetJobList";
    this.urlJobDetail = "/subsite/JobInterface/GetJobDetail";
    this.cache = true;
}

Job.prototype.getJobHotList = function (data, beforeSend, success, complete){
    const params = Json2Param(data);
    var url = getJobfairSite() + this.urlJobHotList + "?" + params;
    AjaxRequest(this.cache, url, "post", data, beforeSend, success, complete);
}

Job.prototype.getJobList = function (data, beforeSend, success, complete){
    const params = Json2Param(data);
    var url = getJobfairSite() + this.urlJobList + "?" + params;
    AjaxRequest(this.cache, url, "post", data, beforeSend, success, complete);
}

Job.prototype.getJobDetail = function (data, beforeSend, success, complete){
    const params = Json2Param(data);
    var url = getJobfairSite() + this.urlJobDetail + "?" + params;
    AjaxRequest(this.cache, url, "post", data, beforeSend, success, complete);
}

function Company() {
    this.urlHotList = "/subsite/CompanyInterface/GetHotCompanyList";
    this.urlList = "/subsite/CompanyInterface/GetCompanyList";
    this.urlDetail = "/subsite/CompanyInterface/GetCompanyDetail";
    this.cache = true;
}

Company.prototype.getCompanyHotList = function (data, beforeSend, success, complete){
    const params = Json2Param(data);
    var url = getJobfairSite() + this.urlHotList + "?" + params;
    AjaxRequest(this.cache, url, "post", data, beforeSend, success, complete);
}

Company.prototype.getCompanyList = function (data, beforeSend, success, complete){
    const params = Json2Param(data);
    var url = getJobfairSite() + this.urlList + "?" + params;
    AjaxRequest(this.cache, url, "post", data, beforeSend, success, complete);
}

Company.prototype.getCompanyDetail = function (data, beforeSend, success, complete){
    const params = Json2Param(data);
    var url = getJobfairSite() + this.urlDetail + "?" + params;
    AjaxRequest(this.cache, url, "post", data, beforeSend, success, complete);
}

function AirKeynote() {
    this.urlList = "/subsite/AirKeynoteInterface/GetAirKeynoteList";
    this.cache = true;
}

AirKeynote.prototype.getAirKeynoteList = function (data, beforeSend, success, complete){
    const params = Json2Param(data);
    var url = getJobfairSite() + this.urlList + "?" + params;
    AjaxRequest(this.cache, url, "post", data, beforeSend, success, complete);
}

function JobFair() {
    this.urlList = "/subsite/JobfairInterface/GetList";
    this.urlPageList = "/subsite/JobfairInterface/GetJobfairList";
    this.urlDetail = "/subsite/JobfairInterface/GetZphCompanyAndJob";
    this.urlDetailCount = "/subsite/JobfairInterface/GetZphDetailCount";
    this.urlCompanyJobList = "/subsite/JobfairInterface/GetCompanyAllJob";
    this.cache = true;
}

JobFair.prototype.getHomeJobFairList = function (data, beforeSend, success, complete){
    const params = Json2Param(data);
    var url = getJobfairSite() + this.urlList + "?" + params;
    AjaxRequest(this.cache, url, "post", data, beforeSend, success, complete);
}

JobFair.prototype.getJobFairList = function (data, beforeSend, success, complete){
    const params = Json2Param(data);
    var url = getJobfairSite() + this.urlPageList + "?" + params;
    AjaxRequest(this.cache, url, "post", data, beforeSend, success, complete);
}

JobFair.prototype.getJobFairDetailCount = function (data, beforeSend, success, complete){
    const params = Json2Param(data);
    var url = getJobfairSite() + this.urlDetailCount + "?" + params;
    AjaxRequest(this.cache, url, "post", data, beforeSend, success, complete);
}

JobFair.prototype.getJobFairDetail = function (data, beforeSend, success, complete){
    const params = Json2Param(data);
    var url = getJobfairSite() + this.urlDetail + "?" + params;
    AjaxRequest(this.cache, url, "post", data, beforeSend, success, complete);
}

JobFair.prototype.getJobFairCompanyJobList = function (data, beforeSend, success, complete){
    const params = Json2Param(data);
    var url = getJobfairSite() + this.urlCompanyJobList + "?" + params;
    AjaxRequest(this.cache, url, "post", data, beforeSend, success, complete);
}