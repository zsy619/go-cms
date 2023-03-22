
var CodeSuccess = 0; // 成功
var CodeError = 1;// 失败
var CodeInvalid = 2;// 无效
var CodeNoAuth = 3;// 无权限
var CodeParamError = 4;// 参数错误
var CodeNoLogin = 5;// 未登录
var CodeNoData = 6;// 无数据
var CodeFatal = 7;// 严重错误

var Code400 = 400;
var Code401 = 401;
var Code402 = 402;
var Code403 = 403;
var Code404 = 404;
var Code500 = 500;
var Code501 = 501;
var Code502 = 502;
var Code503 = 503;
var Code504 = 504;

var Lcjd_Submit = "SUBMITTED"; // 提交
var Lcjd_Company = "COMPANY";   // 公司审批中
var Lcjd_School = "SCHOOL";    // 学校审批中

var Shzt_Refine = "REFINE";   // 完善信息
var Shzt_Await = "AWAIT";    // 等待中
var Shzt_Audit = "AUDIT";      // 审核中，已进入到审核流程
var Shzt_Approval = "APPROVAL"; // 审核通过
var Shzt_Reject = "REJECT";   // 审核拒绝
var Shzt_Recall = "RECALL";   // 撤回

function PrintLayuiTable(tablelayid) {
    //自定义打印table
    let h = window.open("Print_window", "_blank");
    if (h) {
        var v = document.createElement("div");
        var f = ["<head>", "<style>", "body{font-size: 12px; color: #666;}", "table{width: 100%; border-collapse: collapse; border-spacing: 0;}", "th,td{line-height: 20px; padding: 9px 15px; border: 1px solid #ccc; text-align: left; font-size: 12px; color: #666;}", "a{color: #666; text-decoration:none;}", "*.layui-hide{display: none}", ".picture{ }", "</style>", "</head>"].join("");
        $(v).append($("[lay-id=\"" + tablelayid + "\"]").find(".layui-table-box").find(".layui-table-header").html());
        $(v).find("tr").after($("[lay-id=\"" + tablelayid + "\"] .layui-table-body.layui-table-main table").html());
        $(v).find("th.layui-table-patch").remove();
        $(v).find(".layui-table-col-special").remove();
        h.document.write(f + $(v).prop("outerHTML"));
        h.document.close();
        h.print();
        h.close();
    } else {
        layer.msg("打印窗口被阻塞~");
    }
}