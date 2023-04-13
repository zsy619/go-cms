package mp

import (
	"encoding/xml"

	"haedu.gov.cn/cms/app/wechat/models"
)

// message structs
type msgHeader struct {
	XMLName      xml.Name `xml:"xml" json:"-"`
	ToUserName   string   `json:"touser"`
	FromUserName string   `json:"-"`
	CreateTime   int64    `json:"-"`
	MsgType      string   `json:"msgtype"`
}

type textMsg struct {
	msgHeader
	Content string `json:"-"`
	Text    struct {
		Content string `xml:"-" json:"content"`
	} `xml:"-" json:"text"`
}

type imageMsg struct {
	msgHeader
	Image struct {
		MediaId string `json:"media_id"`
	} `json:"image"`
}

type voiceMsg struct {
	msgHeader
	Voice struct {
		MediaId string `json:"media_id"`
	} `json:"voice"`
}

type videoMsg struct {
	msgHeader
	Video *Video `json:"video"`
}

type musicMsg struct {
	msgHeader
	Music *Music `json:"music"`
}

type newsMsg struct {
	msgHeader
	ArticleCount int `json:"-"`
	Articles     struct {
		Item *[]Article `xml:"item" json:"articles"`
	} `json:"news"`
}

// 群发图文消息
type newsGroupMsg struct {
	Filter struct {
		IsToAll bool   `json:"is_to_all"`
		GroupId string `json:"group_id"`
	} `json:"filter"`

	Mpnews struct {
		MediaId string `json:"media_id"`
	} `json:"mpnews"`
	MsgType string `json:"msgtype"`
}

type Video struct {
	MediaId     string `json:"media_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Music struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	MusicUrl     string `json:"musicurl"`
	HQMusicUrl   string `json:"hqmusicurl"`
	ThumbMediaId string `json:"thumb_media_id"`
}

type Article struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	PicUrl      string `json:"picurl"`
	Url         string `json:"url"`
}

type Button struct {
	Type      string   `json:"type,omitempty"`
	Name      string   `json:"name"`
	Key       string   `json:"key,omitempty"`
	Url       string   `json:"url,omitempty"`
	SubButton []Button `json:"sub_button,omitempty"`
}

type qrScene struct {
	ExpireSeconds int64  `json:"expire_seconds,omitempty"`
	ActionName    string `json:"action_name"`
	ActionInfo    struct {
		Scene struct {
			SceneId int64 `json:"scene_id"`
		} `json:"scene"`
	} `json:"action_info"`
}

type (
	UserInfo struct {
		models.CommonError
		Subscribe      int64   `json:"subscribe"`       // 用户是否订阅该公众号标识，值为0时，代表此用户没有关注该公众号，拉取不到其余信息。
		Openid         string  `json:"openid"`          // 用户的标识，对当前公众号唯一
		Nickname       string  `json:"nickname"`        // 用户的昵称
		Sex            int     `json:"sex"`             // 用户的性别，值为1时是男性，值为2时是女性，值为0时是未知
		Language       string  `json:"language"`        // 用户的语言，简体中文为zh_CN
		City           string  `json:"city"`            // 用户所在城市
		Province       string  `json:"province"`        // 用户所在省份
		Country        string  `json:"country"`         // 用户所在国家
		Headimgurl     string  `json:"headimgurl"`      // 用户头像，最后一个数值代表正方形头像大小（有0、46、64、96、132数值可选，0代表640*640正方形头像），用户没有头像时该项为空。若用户更换头像，原有头像URL将失效。
		SubscribeTime  int64   `json:"subscribe_time"`  // 用户关注时间，为时间戳。如果用户曾多次关注，则取最后关注时间
		UnionId        string  `json:"unionid"`         // 只有在用户将公众号绑定到微信开放平台帐号后，才会出现该字段。
		Remark         string  `json:"remark"`          // 公众号运营者对粉丝的备注，公众号运营者可在微信公众平台用户管理界面对粉丝添加备注
		GroupId        int64   `json:"groupid"`         // 用户所在的分组ID（暂时兼容用户分组旧接口）
		TagIdList      []int64 `json:"tagid_list"`      // 用户被打上的标签ID列表
		SubscribeScene string  `json:"subscribe_scene"` // 返回用户关注的渠道来源，ADD_SCENE_SEARCH 公众号搜索，ADD_SCENE_ACCOUNT_MIGRATION 公众号迁移，ADD_SCENE_PROFILE_CARD 名片分享，ADD_SCENE_QR_CODE 扫描二维码，ADD_SCENE_PROFILE_LINK 图文页内名称点击，ADD_SCENE_PROFILE_ITEM 图文页右上角菜单，ADD_SCENE_PAID 支付后关注，ADD_SCENE_WECHAT_ADVERTISEMENT 微信广告，ADD_SCENE_OTHERS 其他
		QrScene        int64   `json:"qr_scene"`        // 二维码扫码场景（开发者自定义）
		QrSceneStr     string  `json:"qr_scene_str"`    // 二维码扫码场景描述（开发者自定义）
	}

	UserListRequest struct {
		UserList []struct {
			OpenId string `json:"openid"` // 用户的标识，对当前公众号唯一
			Lang   string `json:"lang"`   // 国家地区语言版本，zh_CN 简体，zh_TW 繁体，en 英语，默认为zh-CN
		} `json:"user_list"`
	}

	UserListReponse struct {
		models.CommonError
		UserInfoList []UserInfo `json:"user_info_list"`
	}
)

// 标签
type (
	Tag struct {
		Id    int    `json:"id"`    // 标签id，由微信分配
		Name  string `json:"name"`  // 标签名，UTF8编码
		Count int    `json:"count"` // 此标签下粉丝数
	}
)

type Fans struct {
	models.CommonError
	Total      int    `json:"total"`       // 关注该公众账号的总用户数
	Count      int    `json:"count"`       // 拉取的OPENID个数，最大值为10000
	NextOpenId string `json:"next_openid"` // 拉取列表的最后一个用户的OPENID
	Data       struct {
		OpenId []string `json:"openid"`
	} `json:"data"`
}
