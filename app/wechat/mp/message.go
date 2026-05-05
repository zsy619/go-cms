package mp

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/beego/beego/v2/core/logs"
)

type Message struct {
	Request     Request
	AccessToken AccessToken
}

func NewMessage(appId, appSecret string, refreshToken bool) *Message {
	message := &Message{
		Request:     Request{Token: ""},
		AccessToken: AccessToken{AppId: appId, AppSecret: appSecret},
	}
	// 判定是否刷新token
	// 创建菜单需要刷新token
	// 自动消息回复不需要刷新token
	if refreshToken {
		token, _ := message.AccessToken.Fresh()
		message.Request.Token = token
	}
	return message
}

// reply text message
func (message *Message) ReplyTextMsg(rw http.ResponseWriter, content string) error {
	var msg textMsg
	msg.MsgType = "text"
	msg.Content = content
	return message.replyMsg(rw, &msg)
}

// reply image message
func (message *Message) ReplyImageMsg(rw http.ResponseWriter, mediaId string) error {
	var msg imageMsg
	msg.MsgType = "image"
	msg.Image.MediaId = mediaId
	return message.replyMsg(rw, &msg)
}

// reply voice message
func (message *Message) ReplyVoiceMsg(rw http.ResponseWriter, mediaId string) error {
	var msg voiceMsg
	msg.MsgType = "voice"
	msg.Voice.MediaId = mediaId
	return message.replyMsg(rw, &msg)
}

// reply video message
func (message *Message) ReplyVideoMsg(rw http.ResponseWriter, video *Video) error {
	var msg videoMsg
	msg.MsgType = "video"
	msg.Video = video
	return message.replyMsg(rw, &msg)
}

// reply music message
func (message *Message) ReplyMusicMsg(rw http.ResponseWriter, music *Music) error {
	var msg musicMsg
	msg.MsgType = "music"
	msg.Music = music
	return message.replyMsg(rw, &msg)
}

// reply news  message
func (message *Message) ReplyNewsMsg(rw http.ResponseWriter, articles *[]Article) error {
	var msg newsMsg
	msg.MsgType = "news"
	msg.ArticleCount = len(*articles)
	msg.Articles.Item = articles
	return message.replyMsg(rw, &msg)
}

// reply message
func (message *Message) replyMsg(rw http.ResponseWriter, msg any) error {
	v := reflect.ValueOf(msg).Elem()
	v.FieldByName("ToUserName").SetString(message.Request.FromUserName)
	v.FieldByName("FromUserName").SetString(message.Request.ToUserName)
	v.FieldByName("CreateTime").SetInt(time.Now().Unix())
	data, err := xml.Marshal(msg)
	if err != nil {
		return err
	}
	if _, err := rw.Write(data); err != nil {
		return err
	}
	return nil
}

// send text message
func (message *Message) SendTextMsg(touser string, content string) error {
	var msg textMsg
	msg.MsgType = "text"
	msg.Text.Content = content
	return message.sendMsg(touser, &msg)
}

// send image message
func (message *Message) SendImageMsg(touser string, mediaId string) error {
	var msg imageMsg
	msg.MsgType = "image"
	msg.Image.MediaId = mediaId
	return message.sendMsg(touser, &msg)
}

// send voice message
func (message *Message) SendVoiceMsg(touser string, mediaId string) error {
	var msg voiceMsg
	msg.MsgType = "voice"
	msg.Voice.MediaId = mediaId
	return message.sendMsg(touser, &msg)
}

// send video message
func (message *Message) SendVideoMsg(touser string, video *Video) error {
	var msg videoMsg
	msg.MsgType = "video"
	msg.Video = video
	return message.sendMsg(touser, &msg)
}

// send music message
func (message *Message) SendMusicMsg(touser string, music *Music) error {
	var msg musicMsg
	msg.MsgType = "music"
	msg.Music = music
	return message.sendMsg(touser, &msg)
}

// send news message
func (message *Message) SendNewsMsg(touser string, articles *[]Article) error {
	var msg newsMsg
	msg.MsgType = "news"
	msg.Articles.Item = articles
	return message.sendMsg(touser, &msg)
}

// send message
func (message *Message) sendMsg(touser string, msg any) error {
	v := reflect.ValueOf(msg).Elem()
	v.FieldByName("ToUserName").SetString(touser)
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%smessage/custom/send?access_token=", UrlPrefix)
	buf := bytes.NewBuffer(data)
	logs.Debug("sendMsg  ", string(data))
	// retry
	for i := 0; i < retryNum; i++ {
		token := message.Request.Token
		logs.Debug("sendMsg  ", url+token)
		if err != nil {
			logs.Error("sendMsg  ", err.Error())
			if i < retryNum-1 {
				continue
			}
			return err
		}
		if _, err := post(url+token, "text/plain", buf); err != nil {
			logs.Error("sendMsg  ", err.Error())
			if i < retryNum-1 {
				continue
			}
			return err
		}
		break // success
	}
	return nil
}

// 向全部用户群发图文消息
func (message *Message) sendNewsToALl(mediaId string) error {
	var news newsGroupMsg
	news.MsgType = "mpnews"
	news.Filter.IsToAll = true
	news.Mpnews.MediaId = mediaId
	return message.sendGroupMsg(news)
}

// 向特定GroupId用户群发图文消息
func (message *Message) sendNewsToGroup(groupId string, mediaId string) error {
	var news newsGroupMsg
	news.MsgType = "mpnews"
	news.Filter.IsToAll = false
	news.Filter.GroupId = groupId
	news.Mpnews.MediaId = mediaId
	return message.sendGroupMsg(news)
}

// 群发消息
func (message *Message) sendGroupMsg(msg interface{}) error {
	url := fmt.Sprintf("%smessage/mass/sendall?access_token=", UrlPrefix)
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(data)
	// retry
	for i := 0; i < retryNum; i++ {
		token := message.Request.Token
		if err != nil {
			if i < retryNum-1 {
				continue
			}
			return err
		}
		if _, err := post(url+token, "text/plain", buf); err != nil {
			if i < retryNum-1 {
				continue
			}
			return err
		}
		break // success
	}
	return nil
}

// get qrcode url
func (message *Message) GetQRCodeURL(ticket string) string {
	return "https://mp.weixin.qq.com/cgi-bin/showqrcode?ticket=" + ticket
}

// create permanent qrcode
func (message *Message) CreateQRScene(sceneId int64) (string, error) {
	var inf qrScene
	inf.ActionName = "QR_SCENE"
	inf.ActionInfo.Scene.SceneId = sceneId
	return message.createQRCode(&inf)
}

// create temporary qrcode
func (message *Message) CreateQRLimitScene(expireSeconds, sceneId int64) (string, error) {
	var inf qrScene
	inf.ExpireSeconds = expireSeconds
	inf.ActionName = "QR_LIMIT_SCENE"
	inf.ActionInfo.Scene.SceneId = sceneId
	return message.createQRCode(&inf)
}

func (message *Message) createQRCode(inf *qrScene) (string, error) {
	data, err := json.Marshal(inf)
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("%sqrcode/create?access_token=", UrlPrefix)
	buf := bytes.NewBuffer(data)
	ticket := ""
	// retry
	for i := 0; i < retryNum; i++ {
		token := message.Request.Token
		if err != nil {
			if i < retryNum-1 {
				continue
			}
			return "", err
		}
		rtn, err := post(url+token, "text/plain", buf)
		if err != nil {
			if i < retryNum-1 {
				continue
			}
			return "", err
		}
		ticket = rtn.Ticket
		break // success
	}
	return ticket, nil
}

// download media to file
func (message *Message) DownloadMediaFile(mediaId, fileName string) error {
	url := fmt.Sprintf("%sget?media_id=%s&access_token=", MediaUrlPrefix, mediaId)
	// retry
	for i := 0; i < retryNum; i++ {
		token := message.Request.Token
		// if err != nil {
		// 	if i < retryNum-1 {
		// 		continue
		// 	}
		// 	return err
		// }
		resp, err := http.Get(url + token)
		if err != nil {
			if i < retryNum-1 {
				continue
			}
			return err
		}
		defer resp.Body.Close()
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			if i < retryNum-1 {
				continue
			}
			return err
		}
		// json
		if resp.Header.Get("Content-Type") == "text/plain" {
			var rtn response
			if err := json.Unmarshal(data, &rtn); err != nil {
				if i < retryNum-1 {
					continue
				}
				return err
			}
			if i < retryNum-1 {
				continue
			}
			return fmt.Errorf("%d %s", rtn.ErrCode, rtn.ErrMsg)
		}
		// media
		f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, os.ModePerm)
		if err != nil {
			if i < retryNum-1 {
				continue
			}
			return err
		}
		defer f.Close()
		if _, err := f.Write(data); err != nil {
			if i < retryNum-1 {
				continue
			}
			return err
		}
		break // success
	}
	return nil
}

// upload media to file
func (message *Message) UploadMediaFile(mediaType, fileName string) (string, error) {
	var buf bytes.Buffer
	bw := multipart.NewWriter(&buf)
	defer bw.Close()
	f, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer f.Close()
	fw, err := bw.CreateFormFile("filename", f.Name())
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(fw, f); err != nil {
		return "", err
	}
	f.Close()
	bw.Close()
	url := fmt.Sprintf("%supload?type=%s&access_token=", MediaUrlPrefix, mediaType)
	mime := bw.FormDataContentType()
	mediaId := ""
	// retry
	for i := 0; i < retryNum; i++ {
		token := message.Request.Token
		if err != nil {
			if i < retryNum-1 {
				continue
			}
			return "", err
		}
		rtn, err := post(url+token, mime, &buf)
		if err != nil {
			if i < retryNum-1 {
				continue
			}
			return "", err
		}
		mediaId = rtn.MediaId
		break // success
	}
	return mediaId, nil
}

// create custom menu
func (message *Message) CreateCustomMenu(btn *[]Button) error {
	var menu struct {
		Button *[]Button `json:"button"`
	}
	menu.Button = btn
	data, err := json.Marshal(&menu)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(data)
	logs.Debug(buf)
	url := fmt.Sprintf("%smenu/create?access_token=", UrlPrefix)
	// retry
	for i := 0; i < retryNum; i++ {
		token := message.Request.Token
		if err != nil {
			if i < retryNum-1 {
				continue
			}
			return err
		}
		if _, err := post(url+token, "text/plain", buf); err != nil {
			if i < retryNum-1 {
				continue
			}
			return err
		}
		break // success
	}
	return nil
}

// get custom menu
func (message *Message) GetCustomMenu() ([]Button, error) {
	var menu struct {
		Menu struct {
			Button []Button `json:"button"`
		} `json:"menu"`
	}
	url := fmt.Sprintf("%smenu/get?access_token=", UrlPrefix)
	// retry
	for i := 0; i < retryNum; i++ {
		token := message.Request.Token
		// if err != nil {
		// 	if i < retryNum-1 {
		// 		continue
		// 	}
		// 	return nil, err
		// }
		resp, err := http.Get(url + token)
		if err != nil {
			if i < retryNum-1 {
				continue
			}
			return nil, err
		}
		defer resp.Body.Close()
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			if i < retryNum-1 {
				continue
			}
			return nil, err
		}
		// has error?
		var rtn response
		if err := json.Unmarshal(data, &rtn); err != nil {
			if i < retryNum-1 {
				continue
			}
			return nil, err
		}
		// yes
		if rtn.ErrCode != 0 {
			if i < retryNum-1 {
				continue
			}
			return nil, fmt.Errorf("%d %s", rtn.ErrCode, rtn.ErrMsg)
		}
		// no
		if err := json.Unmarshal(data, &menu); err != nil {
			if i < retryNum-1 {
				continue
			}
			return nil, err
		}
		break // success
	}
	return menu.Menu.Button, nil
}

// delete custom menu
func (message *Message) DeleteCustomMenu() error {
	url := UrlPrefix + "menu/delete?access_token="
	// retry
	for i := 0; i < retryNum; i++ {
		token := message.Request.Token
		// if err != nil {
		// 	if i < retryNum-1 {
		// 		continue
		// 	}
		// 	return err
		// }
		if _, err := get(url + token); err != nil {
			if i < retryNum-1 {
				continue
			}
			return err
		}
		break // success
	}
	return nil
}

// get user info
func (message *Message) GetUserInfo(openId string) (UserInfo, error) {
	var uinf UserInfo
	url := fmt.Sprintf("%suser/info?lang=zh_CN&openid=%s&access_token=", UrlPrefix, openId)
	// retry
	for i := 0; i < retryNum; i++ {
		token := message.Request.Token
		// if err != nil {
		// 	if i < retryNum-1 {
		// 		continue
		// 	}
		// 	return uinf, err
		// }
		resp, err := http.Get(url + token)
		if err != nil {
			if i < retryNum-1 {
				continue
			}
			return uinf, err
		}
		defer resp.Body.Close()
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			if i < retryNum-1 {
				continue
			}
			return uinf, err
		}
		// has error?
		var rtn response
		if err := json.Unmarshal(data, &rtn); err != nil {
			if i < retryNum-1 {
				continue
			}
			return uinf, err
		}
		// yes
		if rtn.ErrCode != 0 {
			if i < retryNum-1 {
				continue
			}
			return uinf, fmt.Errorf("%d %s", rtn.ErrCode, rtn.ErrMsg)
		}
		// no
		if err := json.Unmarshal(data, &uinf); err != nil {
			if i < retryNum-1 {
				continue
			}
			return uinf, err
		}
		break // success
	}
	return uinf, nil
}
