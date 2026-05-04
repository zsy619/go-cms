package models

import (
	"crypto/md5"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	_ "github.com/beego/beego/v2/server/web"
)

// 微信支付协议接口数据类，所有的API接口通信都依赖这个数据结构，
// 在调用接口之前先填充各个字段的值，然后进行接口通信，
// 这样设计的好处是可扩展性强，用户可随意对协议进行更改而不用重新设计数据结构，
// 还可以随意组合出不同的协议数据包，不用为每个协议设计一个数据包结构
type WxPayData struct {
	m_values map[string]interface{} // 采用排序的Dictionary的好处是方便对数据包进行签名，不用再签名之前再做一次排序
}

func NewWxPayData() WxPayData {
	result := WxPayData{}
	result.m_values = make(map[string]interface{})
	return result
}

func (wpd *WxPayData) IsEmpty() bool {
	if wpd.m_values == nil || len(wpd.m_values) == 0 {
		return true
	}
	return false
}

/**
* 设置某个字段的值
* @param key 字段名
* @param value 字段值
 */
func (wpd *WxPayData) SetValue(key string, value interface{}) {
	wpd.m_values[key] = value
}

/**
* 根据字段名获取某个字段的值
* @param key 字段名
* @return key对应的字段值
 */
func (wpd *WxPayData) GetValue(key string) interface{} {
	if value, ok := wpd.m_values[key]; ok {
		return value
	}
	return nil
}

/**
* 判断某个字段是否已设置
* @param key 字段名
* @return 若字段key已被设置，则返回true，否则返回false
 */
func (wpd *WxPayData) IsSet(key string) bool {
	_, ok := wpd.m_values[key]
	return ok
}

/**
* @将Dictionary转成xml
* @return 经转换得到的xml串
* @throws WxPayException
**/
func (wpd *WxPayData) ToXml() (string, error) {
	// 数据为空时不能转化为xml格式
	if len(wpd.m_values) == 0 {
		return "", errors.New("wxPayData数据为空")
	}

	xml := "<xml>"
	for key, value := range wpd.m_values {
		if value == nil {
			return "", errors.New("wxPayData内部含有值为null的字段")
		}
		vx, bl := wpd.interfaceToString(value)
		if len(vx) == 0 {
			return "", errors.New("wxPayData字段数据类型错误")
		}
		if bl {
			xml += "<" + key + ">" + vx + "</" + key + ">"
		} else {
			xml += "<" + key + ">" + "<![CDATA[" + vx + "]]></" + key + ">"
		}
		logs.Debug("ToUrl", key, vx)
	}
	xml += "</xml>"
	return xml, nil
}

/**
* @Dictionary格式转化成url参数格式
* @ return url格式串, 该串不包含sign字段值
 */
func (wpd *WxPayData) ToUrl() (string, error) {
	buff := ""
	var keys []string
	for k := range wpd.m_values {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		value := wpd.m_values[key]
		if value == nil {
			return "", errors.New("wxPayData内部含有值为null的字段")
		}
		vx, _ := wpd.interfaceToString(value)
		logs.Debug("ToUrl", key, vx)
		if key == "sign" || vx == "" {
			continue
		}
		buff += key + "=" + vx + "&"
	}
	buff = buff[0 : len(buff)-1]
	return buff, nil
}

func (wpd *WxPayData) interfaceToString(value interface{}) (string, bool) {
	switch value := value.(type) {
	case string:
		return value, false
	case int:
		return strconv.Itoa(value), true
	case int16:
		return strconv.FormatInt(int64(value), 10), true
	case int32:
		return strconv.FormatInt(int64(value), 10), true
	case int64:
		return strconv.FormatInt(value, 10), true
	default:
		return "", false
	}
}

/**
* @Dictionary格式化成Json
* @return json串数据
 */
func (wpd *WxPayData) ToJson() (string, error) {
	datas, err := json.Marshal(wpd.m_values)
	if err != nil {
		return "", err
	}
	return string(datas), nil
}

/**
* @values格式化成能在Web页面上显示的结果（因为web页面上不能直接输出xml格式的字符串）
 */
func (wpd *WxPayData) ToPrintStr() (string, error) {
	str := ""
	for key, value := range wpd.m_values {
		if value == nil {
			return "", errors.New("wxPayData内部含有值为null的字段")
		}
		if key == "sign" || value.(string) == "" {
			continue
		}
		str += fmt.Sprintf("%s=%s<br>", key, value)
	}
	return str, nil
}

/**
* @生成签名，详见签名生成算法
* @return 签名, sign字段不参加签名
* https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=4_3
 */
func (wpd *WxPayData) MakeSign(key string) (string, error) {
	// 转url格式
	str, err := wpd.ToUrl()
	if err != nil {
		return "", err
	}
	// 在string后加入API KEY
	str += "&key=" + key
	logs.Debug("MakeSign ", str)
	// MD5加密
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has) // 将[]byte转成16进制
	// 所有字符转为大写
	return strings.ToUpper(md5str), nil
}

/**
* 检测签名是否正确
* 正确返回true，错误抛异常
 */
func (wpd *WxPayData) CheckSign(key string) (bool, error) {
	// 如果没有设置签名，则跳过检测
	if !wpd.IsSet("sign") {
		return false, errors.New("wxPayData签名存在但不合法")
	}
	value := wpd.GetValue("sign")
	if value == nil || value.(string) == "" { // 如果设置了签名但是签名为空，则抛异常
		return false, errors.New("wxPayData签名存在但不合法")
	}
	// 在本地计算新的签名
	cal_sign, err := wpd.MakeSign(key)
	if err != nil {
		return false, err
	}
	if cal_sign == value {
		return true, nil
	}
	return false, errors.New("wxPayData签名验证错误")
}

/**
* @获取Dictionary
**/
func (wpd *WxPayData) GetValues() map[string]interface{} {
	return wpd.m_values
}

/**
 * @将xml转为WxPayData对象并返回对象内部的数据
 * @param string 待转换的xml串
 * @return 经转换得到的Dictionary
 * @throws WxPayException
 **/
func (wpd *WxPayData) FromXml(contentXml, key string) (map[string]interface{}, error) {
	if len(contentXml) == 0 {
		return nil, NewWxPayException("将空的xml串转换为WxPayData不合法!")
	}

	dec := xml.NewDecoder(strings.NewReader(contentXml))
	ele, val := "", ""

	for t, err := dec.Token(); err == nil; t, err = dec.Token() {
		switch token := t.(type) {
		case xml.StartElement: // 处理元素开始（标签）
			ele = token.Name.Local
			// fmt.Printf("wpd is the sta: %s\n", ele)
			if strings.ToLower(ele) == "xml" {
				// xmlFlag = true
				continue
			}
		case xml.EndElement: // 处理元素结束（标签）
			name := token.Name.Local
			// fmt.Printf("wpd is the end: %s\n", name)
			if strings.ToLower(name) == "xml" {
				break
			}
			if ele == name && ele != "" {
				wpd.SetValue(ele, val)
				ele = ""
				val = ""
			}
		case xml.CharData: // 处理字符数据（这里就是元素的文本）
			// content := string(token)
			// fmt.Printf("wpd is the content: %v\n", content)
			val = string(token)
		default: // 异常处理(Log输出）
			log.Println(token)
		}
	}

	// 2015-06-29 错误是没有签名
	if wpd.GetValue("return_code") != "SUCCESS" {
		return wpd.m_values, nil
	}
	_, err := wpd.CheckSign(key) // 验证签名,不通过会抛异常
	if err != nil {
		return nil, err
	}

	return wpd.m_values, nil
}
