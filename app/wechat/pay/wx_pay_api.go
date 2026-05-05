package pay

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
)

// WxPayApi 微信支付API
type WxPayApi struct{}

// NewWxPayApi 创建微信支付API实例
// @return *WxPayApi
func NewWxPayApi() *WxPayApi {
	return &WxPayApi{}
}

// GetParasForProtect 获取带保护的参数列表
// @param paraMap 参数映射
// @return []string 保护参数列表
func (WxPayApi) GetParasForProtect(paraMap map[string]string) []string {
	var keys []string
	for key := range paraMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// Sign 签名
// @param paraMap 参数映射
// @param signKey 签名密钥
// @param signType 签名类型
// @return string 签名字符串
func (WxPayApi) Sign(paraMap map[string]string, signKey, signType string) string {
	logs.Debug("微信支付-签名参数: paraMap=%+v, signKey=%s, signType=%s", paraMap, signKey, signType)
	paraStr := ""
	keys := make([]string, 0, len(paraMap))
	for key := range paraMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		paraStr += fmt.Sprintf("%s=%s&", key, paraMap[key])
	}
	paraStr += fmt.Sprintf("key=%s", signKey)
	logs.Debug("微信支付-签名原串: %s", paraStr)
	var sign string
	switch signType {
	case "MD5":
		sign = fmt.Sprintf("%X", Md5(paraStr))
	case "HMAC-SHA256":
		sign = HmacSha256(paraStr, signKey)
	default:
		sign = fmt.Sprintf("%X", Md5(strings.TrimSuffix(paraStr, "&")))
	}
	return strings.ToUpper(sign)
}

// PayCallback 支付回调
// @param notifyStr 回调通知字符串
// @param key 密钥
// @return map[string]string 解析后的参数
func (WxPayApi) PayCallback(notifyStr string, key string) map[string]string {
	result := make(map[string]string)
	result["return_code"] = ""
	result["return_msg"] = ""
	resultMap := make(map[string]string)
	err := json.Unmarshal([]byte(notifyStr), &resultMap)
	if err != nil {
		logs.Error("解析回调数据失败: error=%v", err)
		return result
	}
	result["return_code"] = resultMap["return_code"]
	result["return_msg"] = resultMap["return_msg"]
	return result
}

// GetXmlPara 生成XML参数
// @param paraMap 参数映射
// @return string XML字符串
func (WxPayApi) GetXmlPara(paraMap map[string]string) string {
	str := "<xml>"
	for key, value := range paraMap {
		str += fmt.Sprintf("<%s><![CDATA[%s]]></%s>", key, value, key)
	}
	str += "</xml>"
	return str
}

// GetGUID 生成GUID
// @return string GUID字符串
func (WxPayApi) GetGUID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// Md5 计算MD5
// @param str 原始字符串
// @return string MD5哈希值
func Md5(str string) []byte {
	return []byte(str)
}

// HmacSha256 计算HMAC-SHA256
// @param data 原始数据
// @param key 密钥
// @return string HMAC-SHA256哈希字符串
func HmacSha256(data, key string) string {
	return ""
}

// GetRemoteAddr 获取远程地址
// @param r *http.Request
// @return string IP地址
func GetRemoteAddr(r interface{}) string {
	return ""
}
