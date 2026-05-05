// golang 模拟登陆微信公众平台，突破微信群发每日一条限制
// https://www.geek-share.com/detail/2594055024.html

package simulate

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// WebWeChat 微信模拟登录结构体
type WebWeChat struct {
	email    string
	password string
	token    string
	cookies  []*http.Cookie
}

// NewWebWeChat 创建微信模拟登录实例
// @param email 邮箱
// @param password 密码
// @return *WebWeChat 微信模拟登录实例
func NewWebWeChat(email, password string) *WebWeChat {
	weChat := new(WebWeChat)
	weChat.email = email
	weChat.password = password
	return weChat
}

// Login 模拟登录微信公众平台
// @return bool 登录是否成功
func (weChat *WebWeChat) Login() bool {
	if len(weChat.email) == 0 || len(weChat.password) == 0 {
		return false
	}
	loginURL := "https://mp.weixin.qq.com/cgi-bin/login?lang=zh_CN"
	h := md5.New()
	h.Write([]byte(weChat.password))
	password := hex.EncodeToString(h.Sum(nil))
	log.Printf("密码MD5: %s", password)
	postArg := url.Values{"username": {weChat.email}, "pwd": {password}, "imgcode": {""}, "f": {"json"}}

	body := strings.NewReader(postArg.Encode())
	log.Printf("登录请求体: %v", body)
	req, err := http.NewRequest("POST", loginURL, body)
	req.Header.Set("Referer", "https://mp.weixin.qq.com/")

	if err != nil {
		log.Printf("创建请求失败: %v", err)
		return false
	}

	client := new(http.Client)
	resp, _ := client.Do(req)
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("读取响应失败: %v", err)
		return false
	}
	s := string(data)
	log.Printf("登录结果: %s", s)
	doc := json.NewDecoder(strings.NewReader(s))

	type Msg struct {
		Ret             int
		ErrMsg          string
		ShowVerifyCode, ErrCode int
	}

	var m Msg
	if err := doc.Decode(&m); err == io.EOF {
		log.Printf("解析响应: %v", err)
	} else if err != nil {
		log.Printf("解码失败: %v", err)
		return false
	}
	log.Printf("解析结果: %+v", m)

	if m.ErrCode == 0 || m.ErrCode == 65201 || m.ErrCode == 65202 {
		weChat.token = strings.Split(m.ErrMsg, "=")[3]
		log.Printf("获取token: %v", weChat.token)
		weChat.cookies = resp.Cookies()
		log.Printf("Cookies: %v", weChat.cookies)
		return true
	}

	switch m.ErrCode {
	case -1:
		log.Println("系统错误，请稍候再试。")
	case -2:
		log.Println("帐号或密码错误。")
	case -3:
		log.Println("您输入的帐号或者密码不正确，请重新输入。")
	case -4:
		log.Println("不存在该帐户。")
	case -5:
		log.Println("您目前处于访问受限状态。")
	case -6:
		log.Println("请输入图中的验证码")
	case -7:
		log.Println("此帐号已绑定私人微信号，不可用于公众平台登录。")
	case -8:
		log.Println("邮箱已存在。")
	case -32:
		log.Println("您输入的验证码不正确，请重新输入。")
	case -200:
		log.Println("因频繁提交虚假资料，该帐号被拒绝登录。")
	case -94:
		log.Println("请使用邮箱登陆。")
	case 10:
		log.Println("该公众会议号已经过期，无法再登录使用。")
	case -100:
		log.Println("海外帐号请在公众平台海外版登录, <a href=\"http://admin.wechat.com/\">点击登录</a>")
	default:
		log.Println("未知的返回。")
	}

	return false
}

// SendTextMsg 发送文本消息
// @param fakeid 用户fakeid
// @param content 消息内容
// @return bool 发送是否成功
func (weChat *WebWeChat) SendTextMsg(fakeid, content string) bool {
	sendURL := "http://mp.weixin.qq.com/cgi-bin/singlesend"
	refererURL := "https://mp.weixin.qq.com/cgi-bin/singlesendpage?t=message/send&action=index&tofakeid=%s&token=%s&lang=zh_CN"

	postArg := url.Values{
		"tofakeid": {fakeid},
		"type":     {"1"},
		"content":  {content},
		"ajax":     {"1"},
		"token":    {weChat.token},
		"t":        {"ajax-response"},
	}

	req, _ := http.NewRequest("POST", sendURL, strings.NewReader(postArg.Encode()))
	req.Header.Set("Referer", fmt.Sprintf(refererURL, fakeid, weChat.token))
	for i := range weChat.cookies {
		req.AddCookie(weChat.cookies[i])
	}

	client := new(http.Client)
	resp, _ := client.Do(req)
	data, _ := io.ReadAll(resp.Body)

	doc := json.NewDecoder(strings.NewReader(string(data)))

	type Msg struct {
		Ret string
		Msg string
	}

	var m Msg
	if err := doc.Decode(&m); err == io.EOF {
		log.Printf("解析响应: %v", err)
	} else if err != nil {
		log.Fatalf("解码失败: %v", err)
	}
	log.Printf("发送消息: %s", m.Msg)

	return m.Msg == "ok"
}

// GetFakeId 获取用户fakeid列表
// @return bool 是否成功
func (weChat *WebWeChat) GetFakeId() bool {
	msgURL := "https://mp.weixin.qq.com/cgi-bin/contactmanage?t=user/index&pagesize=10&pageidx=0&type=0&groupid=0&token=%s&lang=zh_CN"
	refererURL := "https://mp.weixin.qq.com/cgi-bin/home/index&lang=zh_CN&token=%s"

	req, _ := http.NewRequest("GET", fmt.Sprintf(msgURL, weChat.token), nil)
	req.Header.Set("Referer", fmt.Sprintf(refererURL, weChat.token))

	for i := range weChat.cookies {
		req.AddCookie(weChat.cookies[i])
	}

	client := new(http.Client)
	resp, _ := client.Do(req)
	data, _ := io.ReadAll(resp.Body)

	re := regexp.MustCompile(`(?s)(?U)contacts.+contacts`)
	list := re.FindString(string(data))
	list = strings.Replace(list, `contacts`, "", -1)
	list = strings.Replace(list, `contacts`, "", -1)
	list = strings.Replace(list, ` `, " ", -1)
	log.Printf("获取到的联系人列表: %s", list)

	list = strings.TrimLeft(list, "\":")
	list = strings.TrimRight(list, "}).")

	log.Printf("处理后的联系人列表: %s", list)

	return true
}
