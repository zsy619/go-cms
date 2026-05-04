package mp

import (
	"fmt"
	"sync"
	"time"

	"haedu.gov.cn/cms/app/wechat/models"
)

var (
	GlobalRunning = false         // 是否启动
	GlobalLock    = &sync.Mutex{} // 全局锁
	_GlobalToken  = models.RefreshToken{
		AccessToken:  "",
		RefreshToken: "",
		ExpiresIn:    7200,
	}
	_GlobalTokenGetTime = time.Now()
)

// 获取token
func GlobalToken() models.RefreshToken {
	now := time.Now()
	seconds := now.Sub(_GlobalTokenGetTime).Seconds()
	if seconds > 7000 && seconds <= 7100 { // 异步刷新token
		go func() {
		}()
	}
	if seconds > 7100 { // 同步刷新token
		fmt.Println("同步刷新token")
	}
	return _GlobalToken
}

// 设置token
func SetGlobalToken(token models.RefreshToken) {
	GlobalLock.Lock()
	defer GlobalLock.Unlock()
	_GlobalToken = token
	_GlobalTokenGetTime = time.Now()
}

const (
	TemplateURL = "https://api.weixin.qq.com/cgi-bin/template/"

	// request message types
	MsgTypeText       = "text"
	MsgTypeImage      = "image"
	MsgTypeVoice      = "voice"
	MsgTypeVideo      = "video"
	MsgTypeShortVideo = "shortvideo"
	MsgTypeLocation   = "location" // 地理位置消息
	MsgTypeLink       = "link"
	MsgTypeEvent      = "event"
	// event types
	EventSubscribe   = "subscribe"
	EventUnsubscribe = "unsubscribe"
	EventScan        = "SCAN"
	EventLocation    = "LOCATION"
	EventClick       = "CLICK"
	EventView        = "VIEW"
	// media types
	MediaTypeImage      = "image"
	MediaTypeVoice      = "voice"      //
	MediaTypeVideo      = "video"      // 视频消息
	MediaTypeShortVideo = "shortvideo" // 小视频消息
	MediaTypeThumb      = "thumb"
	// button types
	ButtonTypeClick = "click"
	ButtonTypeView  = "view"
	// environment constants
	UrlPrefix      = "https://api.weixin.qq.com/cgi-bin/"
	MediaUrlPrefix = "http://file.api.weixin.qq.com/cgi-bin/media/"
	retryNum       = 1 // 3
)
