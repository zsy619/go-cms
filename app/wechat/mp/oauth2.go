package mp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"haedu.gov.cn/cms/app/wechat/models"
	"haedu.gov.cn/cms/app/wechat/utils"
)

const (
	redirectOauthURL       string = "https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s#wechat_redirect"
	webAppRedirectOauthURL string = "https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s#wechat_redirect"
	accessTokenURL         string = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
	refreshAccessTokenURL  string = "https://api.weixin.qq.com/sns/oauth2/refresh_token?appid=%s&grant_type=refresh_token&refresh_token=%s"
	userInfoURL            string = "https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN"
	checkAccessTokenURL    string = "https://api.weixin.qq.com/sns/auth?access_token=%s&openid=%s"
)

type OAuth2 struct {
	*models.MpConfig
}

func NewOAuth2(config *models.MpConfig) *OAuth2 {
	rt := new(OAuth2)
	rt.MpConfig = config
	return rt
}

// GetRedirectURL 获取跳转的url地址
func (oauth *OAuth2) GetRedirectURL(redirectURI, scope, state string) (string, error) {
	// url encode
	urlStr := url.QueryEscape(redirectURI)
	return fmt.Sprintf(redirectOauthURL, oauth.AppId, urlStr, scope, state), nil
}

// GetWebAppRedirectURL 获取网页应用跳转的url地址
func (oauth *OAuth2) GetWebAppRedirectURL(redirectURI, scope, state string) (string, error) {
	urlStr := url.QueryEscape(redirectURI)
	return fmt.Sprintf(webAppRedirectOauthURL, oauth.AppId, urlStr, scope, state), nil
}

// Redirect 跳转到网页授权
func (oauth *OAuth2) Redirect(writer http.ResponseWriter, req *http.Request, redirectURI, scope, state string) error {
	// 	参数	是否必须	说明
	// appid	是	公众号的唯一标识
	// redirect_uri	是	授权后重定向的回调链接地址， 请使用 urlEncode 对链接进行处理
	// response_type	是	返回类型，请填写code
	// scope	是	应用授权作用域，snsapi_base （不弹出授权页面，直接跳转，只能获取用户openid），snsapi_userinfo （弹出授权页面，可通过openid拿到昵称、性别、所在地。并且， 即使在未关注的情况下，只要用户授权，也能获取其信息 ）
	// state	否	重定向后会带上state参数，开发者可以填写a-zA-Z0-9的参数值，最多128字节
	// #wechat_redirect	是	无论直接打开还是做页面302重定向时候，必须带此参数

	// https://developers.weixin.qq.com/community/develop/doc/000626c6120830c8da2a8ad215bc00

	location, err := oauth.GetRedirectURL(redirectURI, scope, state)
	if err != nil {
		return err
	}
	http.Redirect(writer, req, location, http.StatusFound)
	return nil
}

// GetUserAccessToken 通过网页授权的code 换取access_token(区别于context中的access_token)
func (oauth *OAuth2) GetUserAccessToken(code string) (result models.ResAccessToken, err error) {
	urlStr := fmt.Sprintf(accessTokenURL, oauth.AppId, oauth.AppSecret, code)
	var response []byte
	response, err = utils.HTTPGet(urlStr)
	if err != nil {
		return
	}
	err = json.Unmarshal(response, &result)
	if err != nil {
		return
	}
	if result.ErrCode != 0 {
		err = fmt.Errorf("GetUserAccessToken error : errcode=%v , errmsg=%v", result.ErrCode, result.ErrMsg)
		return
	}
	return
}

func (oauth *OAuth2) GetGlobalUserAccessToken(code string) (string, string, error) {
	// 第二步：通过code换取网页授权access_token
	// 首先请注意，这里通过code换取的是一个特殊的网页授权access_token,与基础支持中的access_token（该access_token用于调用其他接口）不同。公众号可通过下述接口来获取网页授权access_token。
	// 		如果网页授权的作用域为snsapi_base，则本步骤中获取到网页授权access_token的同时，也获取到了openid，snsapi_base式的网页授权流程即到此为止。
	// 尤其注意：由于公众号的secret和获取到的access_token安全级别都非常高，必须只保存在服务器，不允许传给客户端。后续刷新access_token、通过access_token获取用户信息等步骤，也必须从服务器发起。
	// 	获取code后，请求以下链接获取access_token： https://api.weixin.qq.com/sns/oauth2/access_token?appid=APPID&secret=SECRET&code=CODE&grant_type=authorization_code
	// 	参数说明
	// 参数	是否必须	说明
	// appid	是	公众号的唯一标识
	// secret	是	公众号的appsecret
	// code	是	填写第一步获取的code参数
	// grant_type	是	填写为authorization_code
	if len(code) > 0 {
		token, err := oauth.GetUserAccessToken(code)
		if err != nil {
			logs.Error("GetGlobalUserAccessToken", err.Error())
			return "", "", err
		} else {
			SetGlobalToken(models.RefreshToken{AccessToken: token.AccessToken, RefreshToken: token.RefreshToken, ExpiresIn: token.ExpiresIn})
			return token.AccessToken, token.OpenId, nil
		}
	}
	tokenx := GlobalToken()
	oauth.RefreshToken() // 异步刷新token
	return tokenx.AccessToken, "", nil
}

func (oauth *OAuth2) RefreshToken() {
	GlobalLock.Lock()
	defer GlobalLock.Unlock()
	if GlobalRunning == false {
		GlobalRunning = true
		oauth.RefreshGlobalAccessToken(GlobalToken().RefreshToken) // 刷新token
	}
}

// RefreshAccessToken 刷新access_token
func (oauth *OAuth2) RefreshAccessToken(refreshToken string) (result models.ResAccessToken, err error) {
	urlStr := fmt.Sprintf(refreshAccessTokenURL, oauth.AppId, refreshToken)
	var response []byte
	response, err = utils.HTTPGet(urlStr)
	if err != nil {
		logs.Error("RefreshAccessToken --> ", err.Error())
		return
	}
	err = json.Unmarshal(response, &result)
	if err != nil {
		logs.Error("RefreshAccessToken --> ", err.Error())
		return
	}
	if result.ErrCode != 0 {
		err = fmt.Errorf("GetUserAccessToken error : errcode=%v , errmsg=%v", result.ErrCode, result.ErrMsg)
		logs.Error(err.Error())
		return
	}
	return
}

func (oauth *OAuth2) RefreshGlobalAccessToken(refreshToken string) {
	go func(oauth *OAuth2) {
		expiresIn := 7100                        // 刷新周期，提前100秒刷新
		lastTime := time.Now()                   // 最后一次刷新时间
		if len(GlobalToken().AccessToken) == 0 { // 立即刷新
			expiresIn = 60
		}
	NEW_TICK_DURATION:
		duration := time.Duration(expiresIn) * time.Second
		tick := time.Tick(duration)
		lastTime = time.Now() // 设置下次获取时间

		for {
			select {
			case <-tick:
				logs.Debug("RefreshGlobalAccessToken --> 上次获取时间：", lastTime, " 开始获取AccessToken. ", time.Now())
				result, err := oauth.RefreshAccessToken(GlobalToken().RefreshToken)
				if err != nil {
					logs.Error("RefreshGlobalAccessToken --> 上次获取时间：", lastTime, " 异常信息：", err.Error())
					expiresIn = 60 // 出现异常，60秒后重试
				} else {
					if result.ErrCode != 0 {
						logs.Error("RefreshGlobalAccessToken --> 上次获取时间：", lastTime, " 异常信息：", result.ErrCode, result.ErrMsg)
						expiresIn = 6 // 出现异常，60秒后重试
					} else {
						expiresIn = 7100 // 重新设置刷新周期
						logs.Debug("RefreshGlobalAccessToken --> 上次获取时间：", lastTime, " 获取结果：", result)
						SetGlobalToken(models.RefreshToken{AccessToken: result.AccessToken, RefreshToken: result.RefreshToken, ExpiresIn: result.ExpiresIn})
					}
				}
				goto NEW_TICK_DURATION
			default:
				time.Sleep(90 * time.Second) // 修复90秒
				logs.Debug("RefreshGlobalAccessToken --> 上次获取时间：", lastTime, " No end signal received. ", time.Now())
				continue
			}
		}
	}(oauth)
}

// CheckAccessToken 检验access_token是否有效
func (oauth *OAuth2) CheckAccessToken(accessToken, openID string) (b bool, err error) {
	urlStr := fmt.Sprintf(checkAccessTokenURL, accessToken, openID)
	var response []byte
	response, err = utils.HTTPGet(urlStr)
	if err != nil {
		return
	}
	var result utils.CommonError
	err = json.Unmarshal(response, &result)
	if err != nil {
		return
	}
	if result.ErrCode != 0 {
		b = false
		return
	}
	b = true
	return
}

// GetUserInfo 如果scope为 snsapi_userinfo 则可以通过此方法获取到用户基本信息
func (oauth *OAuth2) GetUserInfo(accessToken, openID string) (result models.UserInfo, err error) {
	// 第四步：拉取用户信息(需scope为 snsapi_userinfo)
	// 如果网页授权作用域为snsapi_userinfo，则此时开发者可以通过access_token和openid拉取用户信息了。
	// 请求方法
	// http：GET（请使用https协议） https://api.weixin.qq.com/sns/userinfo?access_token=ACCESS_TOKEN&openid=OPENID&lang=zh_CN
	// 参数说明
	// 参数	描述
	// access_token	网页授权接口调用凭证,注意：此access_token与基础支持的access_token不同
	// openid	用户的唯一标识
	// lang	返回国家地区语言版本，zh_CN 简体，zh_TW 繁体，en 英语
	urlStr := fmt.Sprintf(userInfoURL, accessToken, openID)
	var response []byte
	response, err = utils.HTTPGet(urlStr)
	if err != nil {
		return
	}
	err = json.Unmarshal(response, &result)
	if err != nil {
		return
	}
	if result.ErrCode != 0 {
		err = fmt.Errorf("GetUserInfo error : errcode=%v , errmsg=%v", result.ErrCode, result.ErrMsg)
		return
	}
	return
}
