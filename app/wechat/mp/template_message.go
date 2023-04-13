package mp

import (
	"encoding/json"
	"errors"
	"fmt"
)

// 参考 https://github.com/slrem/wechat/blob/master/trader/templatemsg.go
type TemplateMessage struct {
	Request     Request
	AccessToken AccessToken
}

func NewTemplateMessage(appId, appSecret string, refreshToken bool) *TemplateMessage {
	message := &TemplateMessage{
		Request:     Request{Token: ""},
		AccessToken: AccessToken{AppId: appId, AppSecret: appSecret},
	}
	//判定是否刷新token
	if refreshToken {
		token, _ := message.AccessToken.Fresh()
		message.Request.Token = token
	}
	return message
}

//设置所属行业
func (this *TemplateMessage) SetIndustry(industryId1, industryId2 int) (err error) {
	surl := TemplateURL + "api_set_industry?access_token=" + this.Request.Token
	var p struct {
		IndustryId1 string `json:"industry_id1"`
		IndustryId2 string `json:"industry_id2"`
	}
	p.IndustryId1 = fmt.Sprint(industryId1)
	p.IndustryId2 = fmt.Sprint(industryId2)
	str, err := json.Marshal(p)
	if err != nil {
		return
	}
	b, err := postjson(surl, string(str))
	if err != nil {
		return
	}
	var r res
	err = json.Unmarshal(b, &r)
	if err != nil {
		return
	}
	if r.ErrCode != 0 {
		err = errors.New(string(b))
	}
	return
}

//获取设置的行业信息
func (this *TemplateMessage) GetIndustry() (jsonstr string, err error) {
	surl := TemplateURL + "get_industry?access_token=" + this.Request.Token
	b, err := getbytes(surl)
	return string(b), err
}

//获得模板ID
func (this *TemplateMessage) GetTemplateId(templateIdShort string) (templateId string, err error) {
	surl := TemplateURL + "api_add_template?access_token=" + this.Request.Token
	var p struct {
		TemplateIdShort string `json:"template_id_short"`
	}
	p.TemplateIdShort = templateIdShort
	str, err := json.Marshal(p)
	if err != nil {
		return
	}
	b, err := postjson(surl, string(str))
	if err != nil {
		return
	}
	var r struct {
		ErrCode    int    `json:"errcode"`
		ErrMsg     string `json:"errmsg"`
		TemplateId string `json:"template_id"`
	}

	err = json.Unmarshal(b, &r)
	if err != nil {
		return
	}
	if r.ErrCode != 0 {
		err = errors.New(string(b))
	} else {
		templateId = r.TemplateId
	}
	return
}

//获取模板列表 return josn字符串
func (this *TemplateMessage) GetALLTemplate() (jsonstr string, err error) {
	surl := TemplateURL + "get_all_private_template?access_token=" + this.Request.Token
	b, err := getbytes(surl)
	return string(b), err
}

//删除模板
func (this *TemplateMessage) DelTemplate(templateId string) (err error) {
	surl := TemplateURL + "del_private_template?access_token=" + this.Request.Token
	var p struct {
		TemplateId string `json:"template_id"`
	}
	p.TemplateId = templateId
	str, err := json.Marshal(p)
	if err != nil {
		return
	}
	b, err := postjson(surl, string(str))
	if err != nil {
		return
	}
	var r res
	err = json.Unmarshal(b, &r)
	if err != nil {
		return
	}
	if r.ErrCode != 0 {
		err = errors.New(string(b))
	}
	return
}

//发送模板消息
func (this *TemplateMessage) SendTemplateMsg(jsonContext string) (msgid int, err error) {
	surl := "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=" + this.Request.Token
	b, err := postjson(surl, jsonContext)
	if err != nil {
		return
	}
	var r struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
		MsgId   int    `json:"msgid"`
	}
	err = json.Unmarshal(b, &r)
	if err != nil {
		return
	}
	if r.ErrCode != 0 {
		err = errors.New(string(b))
	} else {
		msgid = r.MsgId
	}
	return
}
