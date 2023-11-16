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

type WebWeChat struct {
	email    string
	password string
	token    string
	cookies  []*http.Cookie
}

func NewWebWeChat(email, password string) *WebWeChat {
	weChat := new(WebWeChat)
	weChat.email = email
	weChat.password = password
	return weChat
}

func (weChat *WebWeChat) Login() bool {
	if len(weChat.email) == 0 || len(weChat.password) == 0 {
		return false
	}
	login_url := "https://mp.weixin.qq.com/cgi-bin/login?lang=zh_CN"
	h := md5.New()
	h.Write([]byte(weChat.password))
	password := hex.EncodeToString(h.Sum(nil))
	fmt.Println(password)
	post_arg := url.Values{"username": {weChat.email}, "pwd": {password}, "imgcode": {""}, "f": {"json"}}

	body := strings.NewReader(post_arg.Encode())
	fmt.Println(body)
	req, err := http.NewRequest("POST", login_url, body)
	req.Header.Set("Referer", "https://mp.weixin.qq.com/")

	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	client := new(http.Client)
	resp, _ := client.Do(req)
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	s := string(data)
	fmt.Printf("登录结果 --> %s", s)
	doc := json.NewDecoder(strings.NewReader(s))

	type Msg struct {
		Ret                     int
		ErrMsg                  string
		ShowVerifyCode, ErrCode int
	}

	var m Msg
	if err := doc.Decode(&m); err == io.EOF {
		fmt.Println(err)
	} else if err != nil {
		log.Println(err)
		return false
	}
	fmt.Println(m)

	if m.ErrCode == 0 || m.ErrCode == 65201 || m.ErrCode == 65202 {
		weChat.token = strings.Split(m.ErrMsg, "=")[3]
		fmt.Printf("token:%v\n", weChat.token)
		weChat.cookies = resp.Cookies()
		fmt.Println(weChat.cookies)
		return true
	}

	switch m.ErrCode {
	case -1:
		fmt.Println("系统错误，请稍候再试。")
	case -2:
		fmt.Println("帐号或密码错误。")
	case -3:
		fmt.Println("您输入的帐号或者密码不正确，请重新输入。")
	case -4:
		fmt.Println("不存在该帐户。")
	case -5:
		fmt.Println("您目前处于访问受限状态。")
	case -6:
		fmt.Println("请输入图中的验证码")
	case -7:
		fmt.Println("此帐号已绑定私人微信号，不可用于公众平台登录。")
	case -8:
		fmt.Println("邮箱已存在。")
	case -32:
		fmt.Println("您输入的验证码不正确，请重新输入。")
	case -200:
		fmt.Println("因频繁提交虚假资料，该帐号被拒绝登录。")
	case -94:
		fmt.Println("请使用邮箱登陆。")
	case 10:
		fmt.Println("该公众会议号已经过期，无法再登录使用。")
	case -100:
		fmt.Println("海外帐号请在公众平台海外版登录,<a href=\"http://admin.wechat.com/\">点击登录</a>")
	default:
		fmt.Println("未知的返回。")
	}

	return false
}

func (weChat *WebWeChat) SendTextMsg(fakeid string, content string) bool {
	send_url := "http://mp.weixin.qq.com/cgi-bin/singlesend"
	referer_url := "https://mp.weixin.qq.com/cgi-bin/singlesendpage?t=message/send&action=index&tofakeid=%s&token=%s&lang=zh_CN"

	post_arg := url.Values{
		"tofakeid": {fakeid},
		"type":     {"1"},
		"content":  {content},
		"ajax":     {"1"},
		"token":    {weChat.token},
		"t":        {"ajax-response"},
	}

	req, _ := http.NewRequest("POST", send_url, strings.NewReader(post_arg.Encode()))
	req.Header.Set("Referer", fmt.Sprintf(referer_url, fakeid, weChat.token))
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
		fmt.Println(err)
	} else if err != nil {
		log.Fatal(err)
	}
	fmt.Println(m.Msg)

	return m.Msg == "ok"
}

func (weChat *WebWeChat) GetFakeId() bool {
	msg_url := "https://mp.weixin.qq.com/cgi-bin/contactmanage?t=user/index&pagesize=10&pageidx=0&type=0&groupid=0&token=%s&lang=zh_CN"
	referer_url := "https://mp.weixin.qq.com/cgi-bin/home?t=home/index&lang=zh_CN&token=%s"

	req, _ := http.NewRequest("GET", fmt.Sprintf(msg_url, weChat.token), nil)

	req.Header.Set("Referer", fmt.Sprintf(referer_url, weChat.token))

	for i := range weChat.cookies {
		req.AddCookie(weChat.cookies[i])
	}

	client := new(http.Client)
	resp, _ := client.Do(req)
	data, _ := io.ReadAll(resp.Body)

	// fmt.Println(string(data))
	fmt.Println(string(data))
	re := regexp.MustCompile(`(?s)(?U)contacts.+contacts`)
	list := re.FindString(string(data))
	list = strings.Replace(list, `contacts`, "", -1)
	list = strings.Replace(list, `contacts`, "", -1)
	list = strings.Replace(list, ` `, " ", -1)
	fmt.Println(list)

	list = strings.TrimLeft(list, "\":")
	list = strings.TrimRight(list, "}).")

	fmt.Println(list)

	return true
}
