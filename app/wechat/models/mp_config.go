package models

type MpConfig struct {
	AppId          string `db:"app_id" json:"app_id"`                     // 开发者ID
	AppSecret      string `db:"app_secret" json:"app_secret"`             // 开发者密码
	Token          string `db:"token" json:"token"`                       // token
	EncodingAesKey string `db:"encoding_aes_key" json:"encoding_aes_key"` // 消息加密密钥
}
