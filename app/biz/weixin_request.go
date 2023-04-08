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

type ContentSubscribeOrDefaultModel struct {
	AccountID   int64                         `json:"account_id" form:"account_id"`
	RequestType int32                         `json:"request_type" form:"request_type"`
	TextReply   *model.WeixinRequestContent   `json:"text_reply" form:"text_reply"`
	ImageReply  []*model.WeixinRequestContent `json:"image_reply" form:"image_reply"`
	SoundReply  *model.WeixinRequestContent   `json:"sound_reply" form:"sound_reply"`
}

type RuleFindModel struct {
	RuleID       int64     `gorm:"column:rule_id;type:bigint" json:"rule_id" form:"rule_id"`                                  // 主键
	Name         string    `gorm:"column:name;type:varchar(256)" json:"name" form:"name"`                                     // 规则名称
	Keywords     string    `gorm:"column:keywords;type:varchar(2048)" json:"keywords" form:"keywords"`                        // 请求关键词,逗号分隔
	RequestType  int32     `gorm:"column:request_type;type:int" json:"request_type" form:"request_type"`                      // 请求类型(0默认回复1文字2图片3语音4链接5地理位置6关注7取消关注8扫描带参数二维码事件9上报地理位置事件10自定义菜单事件）
	ResponseType int32     `gorm:"column:response_type;type:int" json:"response_type" form:"response_type"`                   // 回复类型(1文本2图文3语音4视频5第三方接口)
	IsLikeQuery  int32     `gorm:"column:is_like_query;type:tinyint" json:"is_like_query" form:"is_like_query"`               // 是否模糊查询
	IsDefault    int32     `gorm:"column:is_default;type:tinyint" json:"is_default" form:"is_default"`                        // 是否默认回复
	SortID       int32     `gorm:"column:sort_id;type:int" json:"sort_id" form:"sort_id"`                                     // 排序
	CreateTime   time.Time `gorm:"column:create_time;type:int unsigned;autoCreateTime" json:"create_time" form:"create_time"` // 创建时间
	UpdateTime   time.Time `gorm:"column:update_time;type:int unsigned;autoUpdateTime" json:"update_time" form:"update_time"` // 修改时间
	ContentID    int64     `gorm:"column:content_id;type:bigint" json:"content_id" form:"content_id"`                         // 主键
	Title        string    `gorm:"column:title;type:varchar(512)" json:"title" form:"title"`                                  // 回复标题
	Content      string    `gorm:"column:content;type:text" json:"content" form:"content"`                                    // 回复内容
	LinkURL      string    `gorm:"column:link_url;type:varchar(512)" json:"link_url" form:"link_url"`                         // 详情链接地址
	ImgURL       string    `gorm:"column:img_url;type:varchar(512)" json:"img_url" form:"img_url"`                            // 图片地址
	MediaURL     string    `gorm:"column:media_url;type:varchar(512)" json:"media_url" form:"media_url"`                      // 语音或视频地址
	MediaHdURL   string    `gorm:"column:media_hd_url;type:varchar(512)" json:"media_hd_url" form:"media_hd_url"`             // 高清语音或者视频地址
}

// ContentFindSubscribeOrDefault 获取关注回复与默认回复
func (w *WeixinRequest) ContentFindSubscribeOrDefault(accountId int64, requestType int32) (*ContentSubscribeOrDefaultModel, error) {
	_, ruleDo := query.WeixinRequestRuleDo()
	// contentMdl, contentDo := query.WeixinRequestContentDo()
	sql1 := fmt.Sprintf(`SELECT a.* FROM weixin_request_content a LEFT JOIN weixin_request_rule b ON a.rule_id=b.rule_id WHERE b.account_id=%d AND b.request_type=%d AND b.response_type=%d limit 1`, accountId, requestType, 1)
	sql2 := fmt.Sprintf(`SELECT a.* FROM weixin_request_content a LEFT JOIN weixin_request_rule b ON a.rule_id=b.rule_id WHERE b.account_id=%d AND b.request_type=%d AND b.response_type=%d ORDER BY a.sort_id`, accountId, requestType, 2)
	sql3 := fmt.Sprintf(`SELECT a.* FROM weixin_request_content a LEFT JOIN weixin_request_rule b ON a.rule_id=b.rule_id WHERE b.account_id=%d AND b.request_type=%d AND b.response_type=%d limit 1`, accountId, requestType, 3)
	model := &ContentSubscribeOrDefaultModel{
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
func (w *WeixinRequest) ContentSaveSubscribeOrDefault(input *ContentSubscribeOrDefaultModel) error {
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
					Name:         "文本回复",
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
				input.TextReply.MediaHdURL = ""
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
						Name:         "图片回复",
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
					Name:         "语音回复",
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

func (w *WeixinRequest) RulePaginate(page, limit int, accountId int64, requestType int32) ([]*RuleFindModel, int64, error) {
	_, ruleDo := query.WeixinRequestRuleDo()
	field := `a.rule_id,a.keywords,a.request_type,a.is_like_query,a.is_default,a.sort_id,a.create_time,a.update_time,b.content_id,b.title,b.content,b.link_url,b.img_url,b.media_url,b.media_hd_url`
	sqlCount := fmt.Sprintf(`SELECT COUNT(*) as Count FROM weixin_request_rule a LEFT JOIN weixin_request_content b ON a.rule_id=b.rule_id WHERE a.account_id=? AND a.request_type=?`)
	sqlSearch := fmt.Sprintf(`SELECT %s FROM weixin_request_rule a LEFT JOIN weixin_request_content b ON a.rule_id=b.rule_id WHERE a.account_id=? AND a.request_type=? AND a.response_type=? ORDER BY a.sort_id ASC`, field)
	var count int64
	if err := ruleDo.Debug().UnderlyingDB().Raw(sqlCount, accountId, requestType).Scan(&count).Error; err != nil {
		return nil, 0, err
	}
	lmt := fmt.Sprintf(" LIMIT %d,%d", (page-1)*limit, limit)
	sqlSearch += lmt
	var list []*RuleFindModel
	if err := ruleDo.Debug().UnderlyingDB().Raw(sqlSearch, accountId, requestType, 1).Scan(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}

func (w *WeixinRequest) RuleFind(ruleId int64) (*RuleFindModel, error) {
	field := `a.rule_id,a.keywords,a.request_type,a.is_like_query,a.is_default,a.sort_id,a.create_time,a.update_time,b.content_id,b.title,b.content,b.link_url,b.img_url,b.media_url,b.media_hd_url`
	sqlSearch := fmt.Sprintf(`SELECT %s FROM weixin_request_rule a LEFT JOIN weixin_request_content b ON a.rule_id=b.rule_id WHERE a.rule_id=?`, field)
	_, ruleDo := query.WeixinRequestRuleDo()
	var ruleFind *RuleFindModel
	if err := ruleDo.Debug().UnderlyingDB().Raw(sqlSearch, ruleId).Scan(&ruleFind).Error; err != nil {
		return nil, err
	}
	return ruleFind, nil
}

func (w *WeixinRequest) RuleSaveSortId(ruleId int64, sortId int32) error {
	ruleMdl, ruleDo := query.WeixinRequestRuleDo()
	_, err := ruleDo.Where(ruleMdl.RuleID.Eq(ruleId)).UpdateColumns(
		map[string]interface{}{
			ruleMdl.SortID.ColumnName().String():     sortId,
			ruleMdl.UpdateTime.ColumnName().String(): time.Now(),
		},
	)
	return err
}

func (w *WeixinRequest) RuleDestory(ruleId int64) error {
	ruleMdl, ruleDo := query.WeixinRequestRuleDo()
	contentMdl, contentDo := query.WeixinRequestContentDo()
	contentDo.Where(contentMdl.RuleID.Eq(ruleId)).Delete()
	ruleDo.Where(ruleMdl.RuleID.Eq(ruleId)).Delete()
	return nil
}
