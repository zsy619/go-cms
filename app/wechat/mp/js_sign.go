package mp

import (
	"crypto/sha1"
	"fmt"
	"strings"
	"time"

	"haedu.gov.cn/cms/app/wechat/utils"

	"github.com/beego/beego/v2/core/logs"
	simplejson "github.com/bitly/go-simplejson"
)

// https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/JS-SDK.html#2
const (
	// APIURLPrefix 微信授权请求
	WxAPIURLPrefix = "https://api.weixin.qq.com/cgi-bin"
	// AuthURL get access_token route
	WxAuthURL = "/token?grant_type=client_credential&"
	// GetTicketURL get ticket route
	WxGetTicketURL = "/ticket/getticket?"
)

var (
	tokenExpire    int64     = 3600 // TokenExpire token缓存的时间
	ticketExpire   int64     = 3600 // TicketExpire tocket缓存时间
	jsToken        string    = ""
	jsTokenExpire  time.Time = time.Now()
	jsTicket       string    = ""
	jsTicketExpire time.Time = time.Now()
)

func getJsTicket() string {
	nowx := time.Now()
	seconds := nowx.Sub(jsTicketExpire).Seconds()
	if seconds >= float64(ticketExpire) { // 过期
		jsTicket = ""
	}
	return jsTicket
}

func setJsTicket(token string, expires int64) {
	ticketExpire = expires
	jsTicket = token
	jsTicketExpire = time.Now()
}

func getJsToken() string {
	nowx := time.Now()
	seconds := nowx.Sub(jsTokenExpire).Seconds()
	if seconds >= float64(tokenExpire) { // 过期
		jsToken = ""
	}
	return jsToken
}

func setJsToken(token string, expires int64) {
	tokenExpire = expires
	jsToken = token
	jsTokenExpire = time.Now()
}

type WxSign struct {
	Appid        string // Appid 公众号appid
	AppSecret    string // AppSecret 公众号秘钥
	TokenRdsKey  string // TokenRdsKey access_token缓存key
	TicketRdsKey string // TicketRdsKey ticket缓存key
}

// WxJsSign
type WxJsSign struct {
	Appid     string `json:"appid"`
	Noncestr  string `json:"noncestr"`
	Timestamp int64  `json:"timestamp"`
	Url       string `json:"url"`
	Signature string `json:"signature"`
}

// New 创建对象
func NewWxSign(appid, secret, tokenKey, ticketKey string) *WxSign {
	return &WxSign{
		Appid:        appid,
		AppSecret:    secret,
		TokenRdsKey:  tokenKey,
		TicketRdsKey: ticketKey,
	}
}

// GetJsSign GetJsSign
func (wSign *WxSign) GetJsSign(url string) (*WxJsSign, error) {
	jsTicket, err := wSign.GetTicket()
	if err != nil {
		return nil, err
	}
	// splite url
	urlSlice := strings.Split(url, "#")
	jsSign := &WxJsSign{
		Appid:    wSign.Appid,
		Noncestr: utils.GenerateNonceStr(),
		// Timestamp: strconv.FormatInt(time.Now().UTC().Unix(), 10),
		Timestamp: time.Now().UTC().Unix(),
		Url:       urlSlice[0],
	}
	jsSign.Signature = Signature(jsTicket, jsSign.Noncestr, jsSign.Timestamp, jsSign.Url)
	return jsSign, nil
}

// Signature 使用权限签名算法
// 参考：https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/JS-SDK.html#62
func Signature(jsTicket, noncestr string, timestamp int64, url string) string {
	h := sha1.New()
	notice := fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s", jsTicket, noncestr, timestamp, url)
	logs.Debug("Signature --> ", notice)
	h.Write([]byte(notice))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// GetAccessToken 获取普通api调用需要的access_token 因为有次数限制，需要缓存
func (wSign *WxSign) GetAccessToken() (accessToken string, err error) {
	// 首先从redis缓存中取出token
	accessToken = getJsToken()
	if len(accessToken) > 0 {
		return accessToken, nil
	}
	/*
		Response:
		{
			"access_token":"bxLdikRXVbTPdHSM05e5u5sUoXNKd8", // token
			"expires_in":7200 // 时效
		}
	*/
	url := fmt.Sprintf("%s%sappid=%s&secret=%s", WxAPIURLPrefix, WxAuthURL, wSign.Appid, wSign.AppSecret)
	logs.Debug("GetAccessToken: ", url)
	_, bs, e := utils.Get(url)
	if e != nil {
		err = fmt.Errorf("GetAccessToken: Request Failed, err-> %v", e)
		logs.Error(err.Error())
		return
	}
	var resp *simplejson.Json
	resp, e = simplejson.NewJson(bs)
	if e != nil {
		err = fmt.Errorf("GetAccessToken: Unmarshal Failed, err-> %v", err)
		logs.Error(err.Error())
		return
	}
	logs.Debug("GetAccessToken：", resp)
	if _, ok := resp.CheckGet("errcode"); ok {
		err = fmt.Errorf("GetAccessToken: get access token err-> %s", string(bs))
		logs.Error(err.Error())
		return
	}
	accessToken = resp.GetPath("access_token").MustString()
	expire := resp.GetPath("expires_in").MustInt64()
	if expire < 1 {
		expire = tokenExpire
	} else {
		expire = expire - 100
	}
	setJsToken(resp.GetPath("access_token").MustString(), expire)
	return
}

// GetTicket 获取JSAPI授权TICKET
func (wSign *WxSign) GetTicket() (ticket string, err error) {
	ticket = getJsTicket()
	if len(ticket) > 0 {
		return ticket, nil
	}
	// 重新获取ticket
	accessToken, e := wSign.GetAccessToken()
	if e != nil {
		err = e
		return
	}
	url := fmt.Sprintf("%s%saccess_token=%s&type=jsapi", WxAPIURLPrefix, WxGetTicketURL, accessToken)
	logs.Debug("GetJsTicket: url --> ", url, " accessToken --> ", accessToken)
	_, bs, e := utils.Get(url)
	if e != nil {
		err = fmt.Errorf("GetJsTicket: get ticket err-> %v", e)
		return
	}
	var resp *simplejson.Json
	resp, e = simplejson.NewJson(bs)
	if e != nil {
		err = fmt.Errorf("GetJsTicket: Unmarshal err-> %v", err)
		return
	}
	/*
		Response:
		{
			"errcode":0,
			"errmsg":"ok",
			"ticket":"bxLdikRXVbTPdHSM05e5u5sUoXNKd8",
			"expires_in":7200
		}
	*/
	logs.Debug("GetJsTicket: ", resp)
	if _, ok := resp.CheckGet("ticket"); !ok {
		err = fmt.Errorf("GetJsTicket: get ticket err-> %s", string(bs))
		return
	}
	ticket = resp.GetPath("ticket").MustString()
	expire := resp.GetPath("expires_in").MustInt64()
	if expire < 1 {
		expire = ticketExpire
	} else {
		expire = expire - 100
	}
	setJsTicket(resp.GetPath("ticket").MustString(), expire)
	// wSign.PushTicketByCache(resp.GetPath("ticket").MustString(), time.Duration(expire)*time.Second)
	return
}
