package models

import (
	"crypto/md5"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/beego/beego/v2/core/logs"
)

// WxPayData 微信支付协议接口数据类
// 所有的API接口通信都依赖这个数据结构，
// 在调用接口之前先填充各个字段的值，然后进行接口通信，
// 这样设计的好处是可扩展性强，用户可随意对协议进行更改而不用重新设计数据结构，
// 还可以随意组合出不同的协议数据包，不用为每个协议设计一个数据包结构
type WxPayData struct {
	m_values map[string]any // 采用排序的Dictionary的好处是方便对数据包进行签名，不用再签名之前再做一次排序
}

// NewWxPayData 创建微信支付数据实例
// @return WxPayData 微信支付数据对象
func NewWxPayData() WxPayData {
	result := WxPayData{}
	result.m_values = make(map[string]any)
	return result
}

// IsEmpty 判断数据是否为空
// @return bool 是否为空
func (wpd *WxPayData) IsEmpty() bool {
	if wpd.m_values == nil || len(wpd.m_values) == 0 {
		return true
	}
	return false
}

// SetValue 设置某个字段的值
// @param key 字段名
// @param value 字段值
func (wpd *WxPayData) SetValue(key string, value any) {
	wpd.m_values[key] = value
}

// GetValue 根据字段名获取某个字段的值
// @param key 字段名
// @return any 字段值
func (wpd *WxPayData) GetValue(key string) any {
	if value, ok := wpd.m_values[key]; ok {
		return value
	}
	return nil
}

// IsSet 判断某个字段是否已设置
// @param key 字段名
// @return bool 是否已设置
func (wpd *WxPayData) IsSet(key string) bool {
	_, ok := wpd.m_values[key]
	return ok
}

// ToXml 将Dictionary转成xml
// @return string xml字符串, error 错误信息
func (wpd *WxPayData) ToXml() (string, error) {
	if len(wpd.m_values) == 0 {
		return "", errors.New("wxPayData数据为空")
	}

	xmlStr := "<xml>"
	for key, value := range wpd.m_values {
		if value == nil {
			return "", errors.New("wxPayData内部含有值为null的字段")
		}
		vx, bl := wpd.interfaceToString(value)
		if len(vx) == 0 {
			return "", errors.New("wxPayData字段数据类型错误")
		}
		if bl {
			xmlStr += "<" + key + ">" + vx + "</" + key + ">"
		} else {
			xmlStr += "<" + key + ">" + "<![CDATA[" + vx + "]]></" + key + ">"
		}
		logs.Debug("ToUrl: key=%s, value=%s", key, vx)
	}
	xmlStr += "</xml>"
	return xmlStr, nil
}

// ToUrl Dictionary格式转化成url参数格式
// @return string url格式串(不包含sign字段值), error 错误信息
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
		logs.Debug("ToUrl: key=%s, value=%s", key, vx)
		if key == "sign" || vx == "" {
			continue
		}
		buff += key + "=" + vx + "&"
	}
	buff = buff[0 : len(buff)-1]
	return buff, nil
}

// interfaceToString 将interface转换为字符串
// @param value 任意类型的值
// @return string 转换后的字符串, bool 是否为基本类型
func (wpd *WxPayData) interfaceToString(value any) (string, bool) {
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

// ToJson Dictionary格式化成Json
// @return string json字符串, error 错误信息
func (wpd *WxPayData) ToJson() (string, error) {
	datas, err := json.Marshal(wpd.m_values)
	if err != nil {
		return "", err
	}
	return string(datas), nil
}

// ToPrintStr 格式化成能在Web页面上显示的结果
// @return string 格式化字符串, error 错误信息
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

// MakeSign 生成签名
// @param key API密钥
// @return string 签名字符串, error 错误信息
func (wpd *WxPayData) MakeSign(key string) (string, error) {
	str, err := wpd.ToUrl()
	if err != nil {
		return "", err
	}
	str += "&key=" + key
	logs.Debug("MakeSign: str=%s", str)
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return strings.ToUpper(md5str), nil
}

// CheckSign 检测签名是否正确
// @param key API密钥
// @return bool 是否正确, error 错误信息
func (wpd *WxPayData) CheckSign(key string) (bool, error) {
	if !wpd.IsSet("sign") {
		return false, errors.New("wxPayData签名存在但不合法")
	}
	value := wpd.GetValue("sign")
	if value == nil || value.(string) == "" {
		return false, errors.New("wxPayData签名存在但不合法")
	}
	calSign, err := wpd.MakeSign(key)
	if err != nil {
		return false, err
	}
	if calSign == value {
		return true, nil
	}
	return false, errors.New("wxPayData签名验证错误")
}

// GetValues 获取Dictionary
// @return map[string]any 字典数据
func (wpd *WxPayData) GetValues() map[string]any {
	return wpd.m_values
}

// FromXml 将xml转为WxPayData对象
// @param contentXml 待转换的xml字符串
// @param key API密钥
// @return map[string]any 转换后的字典, error 错误信息
func (wpd *WxPayData) FromXml(contentXml, key string) (map[string]any, error) {
	if len(contentXml) == 0 {
		return nil, NewWxPayException("将空的xml串转换为WxPayData不合法!")
	}

	dec := xml.NewDecoder(strings.NewReader(contentXml))
	ele, val := "", ""

	for t, err := dec.Token(); err == nil; t, err = dec.Token() {
		switch token := t.(type) {
		case xml.StartElement:
			ele = token.Name.Local
			if strings.ToLower(ele) == "xml" {
				continue
			}
		case xml.EndElement:
			name := token.Name.Local
			if strings.ToLower(name) == "xml" {
				break
			}
			if ele == name && ele != "" {
				wpd.SetValue(ele, val)
				ele = ""
				val = ""
			}
		case xml.CharData:
			val = string(token)
		default:
			logs.Debug("FromXml token: %v", token)
		}
	}

	if wpd.GetValue("return_code") != "SUCCESS" {
		return wpd.m_values, nil
	}
	_, err := wpd.CheckSign(key)
	if err != nil {
		return nil, err
	}

	return wpd.m_values, nil
}