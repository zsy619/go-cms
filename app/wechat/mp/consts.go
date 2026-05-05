package mp

import (
	"sync"
	"time"

	"github.com/beego/beego/v2/core/logs"

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

// GlobalToken 获取token
// @return models.RefreshToken token信息
func GlobalToken() models.RefreshToken {
	now := time.Now()
	seconds := now.Sub(_GlobalTokenGetTime).Seconds()
	if seconds > 7000 && seconds <= 7100 {
		go func() {
		}()
	}
	if seconds > 7100 {
		logs.Debug("同步刷新token")
	}
	return _GlobalToken
}

// SetGlobalToken 设置token
// @param token models.RefreshToken token信息
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
	MediaTypeVoice      = "voice"
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