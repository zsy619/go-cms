package biz

import (
	"errors"
	"fmt"
	"time"

	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/app/dal/query"
)

type WeixinRequest struct{}

func NewWeixinRequest() *WeixinRequest {
	return &WeixinRequest{}
}

type ContentSubscribeOrDefault struct {
	AccountID   int64                         `json:"account_id" form:"account_id"`
	RequestType int32                         `json:"request_type" form:"request_type"`
	TextReply   *model.WeixinRequestContent   `json:"text_reply" form:"text_reply"`
	ImageReply  []*model.WeixinRequestContent `json:"image_reply" form:"image_reply"`
	SoundReply  *model.WeixinRequestContent   `json:"sound_reply" form:"sound_reply"`
}

// ContentFindSubscribeOrDefault 获取关注回复与默认回复
func (w *WeixinRequest) ContentFindSubscribeOrDefault(accountId int64, requestType int32) (*ContentSubscribeOrDefault, error) {
	_, ruleDo := query.WeixinRequestRuleDo()
	// contentMdl, contentDo := query.WeixinRequestContentDo()
	sql1 := fmt.Sprintf(`SELECT a.* FROM weixin_request_content a LEFT JOIN weixin_request_rule b ON a.rule_id=b.rule_id WHERE b.account_id=%d AND b.request_type=%d AND b.response_type=%d limit 1`, accountId, requestType, 1)
	sql2 := fmt.Sprintf(`SELECT a.* FROM weixin_request_content a LEFT JOIN weixin_request_rule b ON a.rule_id=b.rule_id WHERE b.account_id=%d AND b.request_type=%d AND b.response_type=%d ORDER BY a.sort_id`, accountId, requestType, 2)
	sql3 := fmt.Sprintf(`SELECT a.* FROM weixin_request_content a LEFT JOIN weixin_request_rule b ON a.rule_id=b.rule_id WHERE b.account_id=%d AND b.request_type=%d AND b.response_type=%d limit 1`, accountId, requestType, 3)
	model := &ContentSubscribeOrDefault{
		AccountID:   accountId,
		RequestType: requestType,
		TextReply:   &model.WeixinRequestContent{},
		ImageReply:  []*model.WeixinRequestContent{},
		SoundReply:  &model.WeixinRequestContent{},
	}
	ruleDo.UnderlyingDB().Raw(sql1).Scan(&model.TextReply)
	ruleDo.UnderlyingDB().Raw(sql2).Scan(&model.ImageReply)
	ruleDo.UnderlyingDB().Raw(sql3).Scan(&model.SoundReply)
	return model, nil
}

// ContentSaveSubscribeOrDefault 保存关注回复与默认回复
func (w *WeixinRequest) ContentSaveSubscribeOrDefault(input *ContentSubscribeOrDefault) error {
	if input.AccountID <= 0 {
		return errors.New("accountId is empty")
	}
	_, ruleDo := query.WeixinRequestRuleDo()
	_, contentDo := query.WeixinRequestContentDo()

	// 删除文本回复、图片回复、语音回复
	{
		sql := `DELETE FROM weixin_request_content WHERE rule_id IN (SELECT rule_id FROM weixin_request_rule WHERE account_id=? AND request_type=? AND response_type=?);`
		ruleDo.Debug().UnderlyingDB().Exec(sql, input.AccountID, input.RequestType, 1)
		ruleDo.Debug().UnderlyingDB().Exec(sql, input.AccountID, input.RequestType, 2)
		ruleDo.Debug().UnderlyingDB().Exec(sql, input.AccountID, input.RequestType, 3)
	}

	{
		if input.TextReply != nil && input.TextReply.Content != "" {
			// 保存文本回复
			input.TextReply.CreateTime = time.Now()
			input.TextReply.UpdateTime = time.Now()
			if input.TextReply.RuleID <= 0 {
				ruleMdl := &model.WeixinRequestRule{
					AccountID:    input.AccountID,
					RequestType:  input.RequestType,
					ResponseType: 1,
				}
				if err := ruleDo.Save(ruleMdl); err != nil {
					return err
				}
				input.TextReply.RuleID = ruleMdl.RuleID
			}
			{
				input.TextReply.LinkURL = ""
				input.TextReply.ImgURL = ""
				input.TextReply.MediaURL = ""
				input.TextReply.MeidaHdURL = ""
			}
			if err := contentDo.Save(input.TextReply); err != nil {
				return err
			}
		} else {
			// 删除文本回复
			sql := `DELETE FROM weixin_request_rule WHERE account_id=? AND request_type=? AND response_type=?;`
			ruleDo.Debug().UnderlyingDB().Exec(sql, input.AccountID, input.RequestType, 1)
		}
	}

	{
		if input.ImageReply != nil && len(input.ImageReply) > 0 {
			// 保存图片回复
			var ruleID int64
			for _, v := range input.ImageReply {
				if v.RuleID > 0 {
					ruleID = v.RuleID
					break
				}
			}
			for _, v := range input.ImageReply {
				v.CreateTime = time.Now()
				v.UpdateTime = time.Now()
				v.RuleID = ruleID
				if v.RuleID <= 0 {
					ruleMdl := &model.WeixinRequestRule{
						AccountID:    input.AccountID,
						RequestType:  input.RequestType,
						ResponseType: 2,
					}
					if err := ruleDo.Save(ruleMdl); err != nil {
						return err
					}
					ruleID = ruleMdl.RuleID
					v.RuleID = ruleID
				}
				if err := contentDo.Save(v); err != nil {
					return err
				}
			}
		} else {
			// 删除图片回复
			sql := `DELETE FROM weixin_request_rule WHERE account_id=? AND request_type=? AND response_type=?;`
			ruleDo.Debug().UnderlyingDB().Exec(sql, input.AccountID, input.RequestType, 2)
		}
	}

	{
		if input.SoundReply != nil && input.SoundReply.MediaURL != "" {
			// 保存语音回复
			input.SoundReply.CreateTime = time.Now()
			input.SoundReply.UpdateTime = time.Now()
			if input.SoundReply.RuleID <= 0 {
				ruleMdl := &model.WeixinRequestRule{
					AccountID:    input.AccountID,
					RequestType:  input.RequestType,
					ResponseType: 3,
				}
				if err := ruleDo.Save(ruleMdl); err != nil {
					return err
				}
				input.SoundReply.RuleID = ruleMdl.RuleID
			}
			if err := contentDo.Save(input.SoundReply); err != nil {
				return err
			}
		} else {
			// 删除语音回复
			sql := `DELETE FROM weixin_request_rule WHERE account_id=? AND request_type=? AND response_type=?;`
			ruleDo.Debug().UnderlyingDB().Exec(sql, input.AccountID, input.RequestType, 3)
		}
	}

	return nil
}
