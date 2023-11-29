package mp

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"

	"haedu.gov.cn/cms/app/wechat/models"
)

type TokenVerify struct {
	Token string
	models.TokenParam
}

func NewTokenVerify(tokenParam models.TokenParam, token string) *TokenVerify {
	tokenVerify := new(TokenVerify)
	tokenVerify.Token = token
	tokenVerify.TokenParam = tokenParam
	return tokenVerify
}

// Verify 有效性验证
// https://developers.weixin.qq.com/doc/offiaccount/Basic_Information/Access_Overview.htm
func (token *TokenVerify) Verify() (string, error) {
	// 开发者通过检验signature对请求进行校验（下面有校验方式）。
	//	若确认此次GET请求来自微信服务器，请原样返回echostr参数内容，则接入生效，成为开发者成功，否则接入失败。
	// 加密/校验流程如下：
	// 1）将token、timestamp、nonce三个参数进行字典序排序
	// 2）将三个参数字符串拼接成一个字符串进行sha1加密
	// 3）开发者获得加密后的字符串可与signature对比，标识该请求来源于微信
	if len(token.Timestamp) == 0 {
		return "", errors.New("timestamp为空")
	}
	if len(token.Nonce) == 0 {
		return "", errors.New("none为空")
	}
	if len(token.Signature) == 0 {
		return "", errors.New("signature为空")
	}
	// if len(token.EchoStr) == 0 {
	// 	return "", errors.New("echostr为空")
	// }
	signatureGen := token.MakeSignature(token.Timestamp, token.Nonce)
	if signatureGen == token.Signature {
		return token.EchoStr, nil
	}
	return "", errors.New("signature验证错误")
}

func (token *TokenVerify) MakeSignature(timestamp, nonce string) string {
	sl := []string{token.Token, timestamp, nonce}
	sort.Strings(sl)
	s := sha1.New()
	io.WriteString(s, strings.Join(sl, ""))
	return fmt.Sprintf("%x", s.Sum(nil))
}
