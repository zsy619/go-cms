package mp

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/beego/beego/v2/core/logs"
)

// 参考：https://github.com/slrem/wechat/blob/master/trader/common.go
type User struct {
	Request     Request
	AccessToken AccessToken
}

func NewUser(appId, appSecret string, refreshToken bool) *User {
	message := &User{
		Request:     Request{Token: ""},
		AccessToken: AccessToken{AppId: appId, AppSecret: appSecret},
	}
	if refreshToken {
		token, _ := message.AccessToken.Fresh()
		message.Request.Token = token
	}
	return message
}

func (this *User) CreateTag(tagname string) (tagid int, err error) {
	surl := UrlPrefix + "tags/create?access_token=" + this.Request.Token
	var p struct {
		Tag struct {
			Name string `json:"name"`
		} `json:"tag"`
	}
	p.Tag.Name = tagname
	str, err := json.Marshal(p)
	if err != nil {
		return
	}
	b, err := postjson(surl, string(str))
	if err != nil {
		return
	}
	var r struct {
		Tag struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
		} `json:"tag"`
	}
	err = json.Unmarshal(b, &r)
	if err != nil {
		return
	}
	if r.Tag.Id == 0 {
		err = errors.New(string(b))
	} else {
		tagid = r.Tag.Id
	}
	return
}

//获取公众号已创建的标签
func (this *User) GetTag() (tags []Tag, err error) {
	if err != nil {
		return
	}
	surl := UrlPrefix + "tags/get?access_token=" + this.Request.Token
	b, err := getbytes(surl)
	if err != nil {
		return
	}
	var r struct {
		Tags []Tag `json:"tags"`
	}
	err = json.Unmarshal(b, &r)
	if err != nil {
		return
	}
	if len(r.Tags) == 0 {
		err = errors.New(string(b))
	} else {
		tags = r.Tags
	}
	return
}

//编辑标签
// https://developers.weixin.qq.com/doc/offiaccount/User_Management/User_Tag_Management.html
func (this *User) UpdateTag(tagid int, tagname string) (err error) {
	surl := UrlPrefix + "tags/update?access_token=" + this.Request.Token
	var p struct {
		Tag struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
		} `json:"tag"`
	}
	p.Tag.Id, p.Tag.Name = tagid, tagname
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

//删除标签
func (this *User) DelTag(tagid int) (err error) {
	surl := UrlPrefix + "tags/delete?access_token=" + this.Request.Token
	var p struct {
		Tag struct {
			Id int `json:"id"`
		} `json:"tag"`
	}
	p.Tag.Id = tagid
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

//获取标签下粉丝列表
/*
tagid 标签id, nextopenid 第一个拉取的OPENID，不填默认从头开始拉取
*/
func (this *User) GetUserByTag(tagid int, nextopenid string) (useropenid []string, lastopenid string, err error) {
	surl := "https://api.weixin.qq.com/cgi-bin/user/tag/get?access_token=" + this.Request.Token
	logs.Debug(surl)
	var p struct {
		TagId      int    `json:"tagid"`
		NextOpenId string `json:"next_openid"`
	}
	p.TagId, p.NextOpenId = tagid, nextopenid
	str, err := json.Marshal(p)
	if err != nil {
		logs.Debug("GetUserByTag ", " -- 001")
		return
	}
	b, err := postjson(surl, string(str))
	if err != nil {
		logs.Debug("GetUserByTag ", " -- 002")
		return
	}
	if strings.Contains(string(b), "errcode") {
		err = errors.New(string(b))
		logs.Debug("GetUserByTag ", " -- 003")
		return
	}
	var r struct {
		Count int `json:"count"`
		Data  struct {
			OpenId []string `json:"openid"`
		} `json:"data"`
		NextOpenId string `json:"next_openid"`
	}
	err = json.Unmarshal(b, &r)
	if err != nil {
		logs.Debug("GetUserByTag ", " -- 004")
		return
	}
	useropenid, lastopenid = r.Data.OpenId, r.NextOpenId
	return
}

//批量为用户打标签
func (this *User) BatchTagToUsers(useropenids []string, tagid int) (err error) {
	surl := "https://api.weixin.qq.com/cgi-bin/tags/members/batchtagging?access_token=" + this.Request.Token
	var p struct {
		OpenIds []string `json:"openid_list"`
		TagId   int      `json:"tagid"`
	}
	p.OpenIds, p.TagId = useropenids, tagid
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

//批量为用户取消标签
func (this *User) BatchCancelTag(useropenid []string, tagid int) (err error) {
	surl := "https://api.weixin.qq.com/cgi-bin/tags/members/batchuntagging?access_token=" + this.Request.Token
	var p struct {
		OpenIds []string `json:"openid_list"`
		TagId   int      `json:"tagid"`
	}
	p.OpenIds, p.TagId = useropenid, tagid
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

//获取用户身上的标签列表
func (this *User) GetTagsByUser(useropenid string) (tagids []int, err error) {
	surl := "https://api.weixin.qq.com/cgi-bin/tags/getidlist?access_token=" + this.Request.Token
	var p struct {
		OpenId string `json:"openid"`
	}
	p.OpenId = useropenid
	str, err := json.Marshal(p)
	if err != nil {
		return
	}
	b, err := postjson(surl, string(str))
	if err != nil {
		return
	}
	if strings.Contains(string(b), "errcode") {
		err = errors.New(string(b))
		return
	}
	var r struct {
		TagIds []int `json:"tagid_list"`
	}
	err = json.Unmarshal(b, &r)
	if err != nil {
		return
	}
	tagids = r.TagIds
	return
}

//设置用户备注名
func (this *User) SetRemark(useropenid string, remark string) (err error) {
	surl := "https://api.weixin.qq.com/cgi-bin/user/info/updateremark?access_token=" + this.Request.Token
	var p struct {
		OpenId string `json:"openid"`
		Remark string `json:"remark"`
	}
	p.OpenId, p.Remark = useropenid, remark
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

//获取用户基本信息（包括UnionID机制）
// 参考：https://developers.weixin.qq.com/doc/offiaccount/User_Management/Get_users_basic_information_UnionID.html#UinonId
func (this *User) GetUserInfo(openid string) (user UserInfo, err error) {
	surl := "https://api.weixin.qq.com/cgi-bin/user/info?access_token=" + this.Request.Token + "&openid=" + openid + "&lang=zh_CN "
	b, err := getbytes(surl)
	if err != nil {
		return
	}
	err = json.Unmarshal(b, &user)
	if err != nil {
		return
	}
	if user.Openid == "" {
		err = errors.New(string(b))
	}
	return
}

// 批量获取用户基本信息
// 开发者可通过该接口来批量获取用户基本信息。最多支持一次拉取100条。
func (this *User) BatchGetUserInfo(input UserListRequest) (users UserListReponse, err error) {
	surl := "https://api.weixin.qq.com/cgi-bin/user/info/batchget?access_token=" + this.Request.Token
	str, err := json.Marshal(input)
	if err != nil {
		return
	}
	b, err := postjson(surl, string(str))
	if err != nil {
		return
	}
	err = json.Unmarshal(b, &users)
	if err != nil {
		return
	}
	return
}

//获取用户列表
/*
nextopenid 第一个拉取的OPENID，填空默认从头开始拉取
附：关注者数量超过10000时
当公众号关注者数量超过10000时，可通过填写next_openid的值，从而多次拉取列表的方式来满足需求。
参考：https://developers.weixin.qq.com/doc/offiaccount/User_Management/Getting_a_User_List.html
*/
func (this *User) GetFans(nextopenid string) (fans Fans, err error) {
	surl := "https://api.weixin.qq.com/cgi-bin/user/get?access_token=" + this.Request.Token
	if len(nextopenid) > 0 {
		surl += "&next_openid=" + nextopenid
	}
	logs.Debug(surl)
	b, err := getbytes(surl)
	if err != nil {
		logs.Debug("GetFans", " 001 ---> ", err.Error())
		return
	}
	if strings.Contains(string(b), "errcode") {
		err = errors.New(string(b))
		logs.Debug("GetFans", " 002 ---> ", string(b))
		return
	}

	err = json.Unmarshal(b, &fans)
	if err != nil {
		logs.Debug("GetFans", " 003 ---> ", err.Error())
	}
	return
}
