package pay

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/cms/app/wechat/models"
	"haedu.gov.cn/cms/app/wechat/utils"
)

type JsApiPay struct {
	TotalFee           int                 // 商品金额，用于统一下单
	OpenId             string              // 用于调用统一下单接口
	AccessToken        string              // 用于获取收货地址js函数入口参数
	AppId              string              //
	Key                string              // 支付appkey
	AppSecret          string              //
	UnifiedOrderResult models.WxPayData    //
	w                  http.ResponseWriter //
	r                  *http.Request       // 页面请求
}

func NewJsApiPay(totalFee int, appId, appSecret string, w http.ResponseWriter, r *http.Request) *JsApiPay {
	jsApiPay := new(JsApiPay)
	jsApiPay.TotalFee = totalFee
	jsApiPay.AppId = appId
	jsApiPay.AppSecret = appSecret
	jsApiPay.w = w
	jsApiPay.r = r
	jsApiPay.UnifiedOrderResult = models.NewWxPayData() // 初始化
	return jsApiPay
}

/**
* 网页授权获取用户基本信息的全部过程
* 详情请参看网页授权获取用户基本信息：http://mp.weixin.qq.com/wiki/17/c0f37d5704f0b64713d5d2c37b468d75.html
* 第一步：利用url跳转获取code
* 第二步：利用code去获取openid和access_token
**/
// 网页授权获取用户基本信息的全部过程
func (this *JsApiPay) GetOpenidAndAccessToken() {
	code := this.r.URL.Query().Get("code")
	if len(code) > 0 {
		this.GetOpenidAndAccessTokenFromCode(code)
	} else {
		// 构造网页授权获取code的URL
		host := this.r.URL.Host
		path := this.r.URL.Path
		redirect_uri := url.QueryEscape("http://" + host + path)
		fmt.Println(redirect_uri)
		// data := models.NewWxPayData()
		// data.SetValue("appid", this.AppId)
		// data.SetValue("redirect_uri", redirect_uri)
		// data.SetValue("response_type", "code")
		// data.SetValue("scope", "snsapi_base")
		// data.SetValue("state", "STATE"+"#wechat_redirect")
		// url, err := data.ToUrl()
		// if err != nil {
		// 	fmt.Println(err.Error())
		// }
		// url = "https://open.weixin.qq.com/connect/oauth2/authorize?" + url
		// http.Redirect(this.w, this.r, url, 302)
	}
}

/**
* 通过code换取网页授权access_token和openid的返回数据，正确时返回的JSON数据包如下：
* {
*  "access_token":"ACCESS_TOKEN",
*  "expires_in":7200,
*  "refresh_token":"REFRESH_TOKEN",
*  "openid":"OPENID",
*  "scope":"SCOPE",
*  "unionid": "o6_bmasdasdsad6_2sgVt7hMZOPfL"
* }
* 其中access_token可用于获取共享收货地址
* openid是微信支付jsapi支付接口统一下单时必须的参数
* 更详细的说明请参考网页授权获取用户基本信息：http://mp.weixin.qq.com/wiki/17/c0f37d5704f0b64713d5d2c37b468d75.html
* @失败时抛异常WxPayException
**/
// 通过code换取网页授权access_token和openid的返回数据
// <param name="code"></param>
// <exception cref="WxPayException"></exception>
func (this *JsApiPay) GetOpenidAndAccessTokenFromCode(code string) error {
	// 构造获取openid及access_token的url
	data := models.NewWxPayData()
	data.SetValue("appid", this.AppId)
	data.SetValue("secret", this.AppSecret)
	data.SetValue("code", code)
	data.SetValue("grant_type", "authorization_code")
	url, err := data.ToUrl()
	if err != nil {
		return err
	}
	url = "https://api.weixin.qq.com/sns/oauth2/access_token?" + url
	datas, err := utils.HTTPGet(url)
	if err != nil {
		return err
	}
	// result := string(datas)
	// fmt.Println(result)
	result := models.ResAccessToken{}
	err = json.Unmarshal(datas, &result)
	if err != nil {
		return err
	}
	this.OpenId = result.OpenId
	this.AccessToken = result.AccessToken
	return nil
}

/**
 * 调用统一下单，获得下单结果
 * @return 统一下单结果
 * @失败时抛异常WxPayException
 **/

// <summary>
// 调用统一下单，获得下单结果
// </summary>
// <returns></returns>
// <exception cref="WxPayException"></exception>
func (this *JsApiPay) GetUnifiedOrderResult(orderNo string) (models.WxPayData, error) {
	// 统一下单
	data := models.NewWxPayData()
	data.SetValue("body", orderNo)                                   // 商品描述 商品简单描述，该字段请按照规范传递
	data.SetValue("attach", orderNo)                                 // 附加数据 附加数据，在查询API和支付通知中原样返回，可作为自定义参数使用。
	data.SetValue("out_trade_no", orderNo)                           // 商户订单号 商户系统内部订单号，要求32个字符内，只能是数字、大小写字母_-|* 且在同一个商户号下唯一。详见商户订单号
	data.SetValue("total_fee", this.TotalFee)                        // 标价金额 订单总金额，单位为分
	data.SetValue("time_start", time.Now().Format("20060102150405")) // 交易起始时间 订单生成时间，格式为yyyyMMddHHmmss，如2009年12月25日9点10分10秒表示为20091225091010。
	// 	订单失效时间，格式为yyyyMMddHHmmss，如2009年12月27日9点10分10秒表示为20091227091010。订单失效时间是针对订单号而言的，由于在请求支付的时候有一个必传参数prepay_id只有两小时的有效期，所以在重入时间超过2小时的时候需要重新请求下单接口获取新的prepay_id。
	// 	time_expire只能第一次下单传值，不允许二次修改，二次修改系统将报错。如用户支付失败后，需再次支付，需更换原订单号重新下单。
	// 	建议：最短失效时间间隔大于1分钟
	data.SetValue("time_expire", time.Now().Add(120*time.Minute).Format("20060102150405")) // 交易结束时间
	// data.SetValue("goods_tag", "奋进新孔店@支付") //订单优惠标记，使用代金券或立减优惠功能时需要的参数，说明详见代金券或立减优惠

	// 	JSAPI--JSAPI支付（或小程序支付）、NATIVE--Native支付、APP--app支付，MWEB--H5支付，不同trade_type决定了调起支付的方式，请根据支付产品正确上传
	// 	MICROPAY--付款码支付，付款码支付有单独的支付接口，所以接口不需要上传，该字段在对账单中会出现
	data.SetValue("trade_type", "JSAPI") // 交易类型

	// trade_type=JSAPI时（即JSAPI支付），此参数必传，此参数为微信用户在商户对应appid下的唯一标识。openid如何获取，可参考【获取openid】。
	// 企业号请使用【企业号OAuth2.0接口】获取企业号内成员userid，再调用【企业号userid转openid接口】进行转换
	data.SetValue("openid", this.OpenId) // 用户标识

	xml, _ := data.ToXml()
	logs.Debug("GetUnifiedOrderResult ", xml)

	result, err := UnifiedOrder(data, 6)
	if err != nil {
		logs.Error("GetUnifiedOrderResult ", err.Error())
		return models.NewWxPayData(), err
	}
	logs.Info("GetUnifiedOrderResult ", result)
	if !result.IsSet("appid") || !result.IsSet("prepay_id") || result.GetValue("prepay_id").(string) == "" {
		return data, errors.New("UnifiedOrder response error!")
	}

	this.UnifiedOrderResult = *result
	return *result, nil
}

/**
 *
 * 从统一下单成功返回的数据中获取微信浏览器调起jsapi支付所需的参数，
 * 微信浏览器调起JSAPI时的输入参数格式如下：
 * {
 *   "appId" : "wx2421b1c4370ec43b",     //公众号名称，由商户传入
 *   "timeStamp":" 1395712654",         //时间戳，自1970年以来的秒数
 *   "nonceStr" : "e61463f8efa94090b1f366cccfbbb444", //随机串
 *   "package" : "prepay_id=u802345jgfjsdfgsdg888",
 *   "signType" : "MD5",         //微信签名方式:
 *   "paySign" : "70EA570631E4BB79628FBCA90534C63FF7FADD89" //微信签名
 * }
 * @return string 微信浏览器调起JSAPI时的输入参数，json格式可以直接做参数用
 * 更详细的说明请参考网页端调起支付API：http://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=7_7
 *
 **/
// <summary>
// 从统一下单成功返回的数据中获取微信浏览器调起jsapi支付所需的参数
// </summary>
// <returns></returns>
func (this *JsApiPay) GetJsApiParameters() (string, models.WxPayData) {
	logs.Debug("JsApiPay::GetJsApiParam is processing...")

	jsApiParam := models.NewWxPayData()
	jsApiParam.SetValue("appId", this.UnifiedOrderResult.GetValue("appid"))
	jsApiParam.SetValue("timeStamp", utils.CurrentTimeStampString())
	jsApiParam.SetValue("nonceStr", utils.GenerateNonceStr())
	jsApiParam.SetValue("package", "prepay_id="+this.UnifiedOrderResult.GetValue("prepay_id").(string))
	jsApiParam.SetValue("signType", "MD5")
	paySign, err := jsApiParam.MakeSign(this.Key)
	if err != nil {
		logs.Error(err.Error())
	}
	jsApiParam.SetValue("paySign", paySign)

	parameters, err := jsApiParam.ToJson()
	if err != nil {
		logs.Error(err.Error())
	}
	logs.Debug("Get jsApiParam : " + parameters)
	return parameters, jsApiParam
}
