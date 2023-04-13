package pay

import (
	"fmt"
	"smart_sso/controllers/wechat/models"
	"smart_sso/controllers/wechat/utils"
	"smart_sso/global"
	"time"

	"github.com/beego/beego/v2/core/logs"
)

/**
* 统一下单
* @param WxPaydata inputObj 提交给统一下单API的参数
* @param int timeOut 超时时间
* @throws WxPayException
* @return 成功时返回，其他抛异常
**/
func UnifiedOrder(inputObj models.WxPayData, timeOut int) (*models.WxPayData, error) {
	url := "https://api.mch.weixin.qq.com/pay/unifiedorder"
	//检测必填参数
	if !inputObj.IsSet("out_trade_no") {
		return nil, models.NewWxPayException("缺少统一支付接口必填参数out_trade_no！")
	}
	if !inputObj.IsSet("body") {
		return nil, models.NewWxPayException("缺少统一支付接口必填参数body！")
	}
	if !inputObj.IsSet("total_fee") {
		return nil, models.NewWxPayException("缺少统一支付接口必填参数total_fee！")
	}
	if !inputObj.IsSet("trade_type") {
		return nil, models.NewWxPayException("缺少统一支付接口必填参数trade_type！")
	}

	//关联参数
	if inputObj.GetValue("trade_type").(string) == "JSAPI" && !inputObj.IsSet("openid") {
		return nil, models.NewWxPayException("统一支付接口中，缺少必填参数openid！trade_type为JSAPI时，openid为必填参数！")
	}
	if inputObj.GetValue("trade_type").(string) == "NATIVE" && !inputObj.IsSet("product_id") {
		return nil, models.NewWxPayException("统一支付接口中，缺少必填参数product_id！trade_type为JSAPI时，product_id为必填参数！")
	}

	//异步通知url未设置，则使用配置文件中的url
	if !inputObj.IsSet("notify_url") {
		inputObj.SetValue("notify_url", global.WebSite+"wechat/paynotify") //异步通知url CurrentUrl + "alipay/wx_notify_url.aspx";
	}

	inputObj.SetValue("appid", global.WxMpConfig.AppId)      //公众账号ID 微信支付分配的公众账ID（企业号corpid即为此appId）
	inputObj.SetValue("mch_id", global.WxMchConfig.MchId)    //商户号 微信支付分配的商户号
	inputObj.SetValue("spbill_create_ip", "127.0.0.1")       //终端ip 支持IPV4和IPV6两种格式的IP地址。用户的客户端IP
	inputObj.SetValue("sign_type", "MD5")                    //签名类型，默认为MD5，支持HMAC-SHA256和MD5。
	inputObj.SetValue("nonce_str", utils.GenerateNonceStr()) //随机字符串 随机字符串，长度要求在32位以内。推荐随机数生成算法

	//签名
	sign, err := inputObj.MakeSign(global.WxMchConfig.MchApiKey) //签名 通过签名算法计算得出的签名值，详见签名生成算法
	if err != nil {
		logs.Error("UnifiedOrder ", err.Error())
		return nil, err
	}
	inputObj.SetValue("sign", sign)

	xml, err := inputObj.ToXml()
	if err != nil {
		logs.Error("UnifiedOrder ", err.Error())
		return nil, err
	}
	logs.Debug("UnifiedOrder ", xml)

	start := time.Now()

	response, err := utils.HttpServicePost(xml, url, false, timeOut)

	end := time.Now()
	timeCost := end.Sub(start).Microseconds()

	result := models.NewWxPayData()
	fmt.Println(result.IsEmpty())
	result.FromXml(string(response), global.WxMchConfig.MchApiKey)

	logs.Debug("UnifiedOrder ", string(response))

	ReportCostTime(url, timeCost, result) //测速上报

	return &result, nil
}

/**
* 测速上报
* @param string interface_url 接口URL
* @param int timeCost 接口耗时
* @param WxPayData inputObj参数数组
**/
func ReportCostTime(interface_url string, timeCost int64, inputObj models.WxPayData) {

}
