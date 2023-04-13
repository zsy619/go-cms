package models

// https://developers.weixin.qq.com/doc/offiaccount/Basic_Information/Get_access_token.html
type ResAccessToken struct {
	CommonError

	AccessToken  string `json:"access_token"`  //网页授权接口调用凭证,注意：此access_token与基础支持的access_token不同
	ExpiresIn    int    `json:"expires_in"`    //access_token接口调用凭证超时时间，单位（秒）
	RefreshToken string `json:"refresh_token"` //用户刷新access_token
	OpenId       string `json:"openid"`        //用户唯一标识，请注意，在未关注公众号时，用户访问公众号的网页，也会产生一个用户和公众号唯一的OpenID
	Scope        string `json:"scope"`         //用户授权的作用域，使用逗号（,）分隔

	// UnionID 只有在用户将公众号绑定到微信开放平台帐号后，才会出现该字段。
	// 公众号文档 https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1421140842
	UnionID string `json:"unionid"`
}

type RefreshToken struct {
	AccessToken  string `json:"access_token"`  //网页授权接口调用凭证,注意：此access_token与基础支持的access_token不同
	ExpiresIn    int    `json:"expires_in"`    //access_token接口调用凭证超时时间，单位（秒）
	RefreshToken string `json:"refresh_token"` //用户刷新access_token
}
