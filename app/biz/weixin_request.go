package biz

import (
	"errors"
	"fmt"
	"time"

	"haedu.gov.cn/cms/app/biz/bmodel"
	"haedu.gov.cn/cms/app/dal/model"
	"haedu.gov.cn/cms/app/dal/query"
)

type WeixinRequest struct{}

func NewWeixinRequest() *WeixinRequest {
	return &WeixinRequest{}
}

// ContentFindSubscribeOrDefault 获取关注回复与默认回复
func (w *WeixinRequest) ContentFindSubscribeOrDefault(accountId int64, requestType int32) (*bmodel.Weixin_ContentSubscribeOrDefaultModel, error) {
	_, ruleDo := query.WeixinRequestRuleDo()
	// contentMdl, contentDo := query.WeixinRequestContentDo()
	sql1 := fmt.Sprintf(`SELECT a.* FROM weixin_request_content a LEFT JOIN weixin_request_rule b ON a.rule_id=b.rule_id WHERE b.account_id=%d AND b.request_type=%d AND b.response_type=%d limit 1`, accountId, requestType, 1)
	sql2 := fmt.Sprintf(`SELECT a.* FROM weixin_request_content a LEFT JOIN weixin_request_rule b ON a.rule_id=b.rule_id WHERE b.account_id=%d AND b.request_type=%d AND b.response_type=%d ORDER BY a.sort_id`, accountId, requestType, 2)
	sql3 := fmt.Sprintf(`SELECT a.* FROM weixin_request_content a LEFT JOIN weixin_request_rule b ON a.rule_id=b.rule_id WHERE b.account_id=%d AND b.request_type=%d AND b.response_type=%d limit 1`, accountId, requestType, 3)
	model := &bmodel.Weixin_ContentSubscribeOrDefaultModel{
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
func (w *WeixinRequest) ContentSaveSubscribeOrDefault(input *bmodel.Weixin_ContentSubscribeOrDefaultModel) error {
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

// RulePictureFind 规则文本回复查询
func (w *WeixinRequest) RulePictureFind(ruleId int64) ([]*model.WeixinRequestContent, error) {
	contentMdl, contentDo := query.WeixinRequestContentDo()
	return contentDo.Where(contentMdl.RuleID.Eq(ruleId)).Order(contentMdl.SortID).Find()
}

// RulePaginate 规则分页查询
func (w *WeixinRequest) RulePaginate(page, limit int, accountId int64, requestType int32) ([]*bmodel.Weixin_RuleModel, int64, error) {
	_, ruleDo := query.WeixinRequestRuleDo()
	field := `a.rule_id,a.account_id,a.name,a.keywords,a.request_type,a.is_like_query,a.is_default,a.sort_id,a.create_time,a.update_time,b.content_id,b.title,b.content,b.link_url,b.img_url,b.media_url,b.media_hd_url`
	sqlCount := fmt.Sprintf(`SELECT COUNT(*) as Count FROM weixin_request_rule a LEFT JOIN weixin_request_content b ON a.rule_id=b.rule_id WHERE a.account_id=? AND a.request_type=?`)
	sqlSearch := fmt.Sprintf(`SELECT %s FROM weixin_request_rule a LEFT JOIN weixin_request_content b ON a.rule_id=b.rule_id WHERE a.account_id=? AND a.request_type=? ORDER BY a.sort_id ASC`, field)
	var count int64
	if err := ruleDo.Debug().UnderlyingDB().Raw(sqlCount, accountId, requestType).Scan(&count).Error; err != nil {
		return nil, 0, err
	}
	lmt := fmt.Sprintf(" LIMIT %d,%d", (page-1)*limit, limit)
	sqlSearch += lmt
	var list []*bmodel.Weixin_RuleModel
	if err := ruleDo.Debug().UnderlyingDB().Raw(sqlSearch, accountId, requestType).Scan(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}

// RulePaginateCount 规则汇总分页
func (w *WeixinRequest) RulePaginateCount(page, limit int, accountId int64, requestType int32) ([]*bmodel.Weixin_RuleCountModel, int64, error) {
	_, ruleDo := query.WeixinRequestRuleDo()
	field := `a.rule_id,a.account_id,a.name,a.keywords,a.request_type,a.is_like_query,a.is_default,a.sort_id,a.create_time,a.update_time,(SELECT COUNT(1) AS count FROM weixin_request_content b WHERE b.rule_id=a.rule_id) AS count`
	sqlCount := fmt.Sprintf(`SELECT COUNT(*) as Count FROM weixin_request_rule a WHERE a.account_id=? AND a.request_type=?`)
	sqlSearch := fmt.Sprintf(`SELECT %s FROM weixin_request_rule a WHERE a.account_id=? AND a.request_type=? ORDER BY a.sort_id ASC`, field)
	var count int64
	if err := ruleDo.Debug().UnderlyingDB().Raw(sqlCount, accountId, requestType).Scan(&count).Error; err != nil {
		return nil, 0, err
	}
	lmt := fmt.Sprintf(" LIMIT %d,%d", (page-1)*limit, limit)
	sqlSearch += lmt
	var list []*bmodel.Weixin_RuleCountModel
	if err := ruleDo.Debug().UnderlyingDB().Raw(sqlSearch, accountId, requestType).Scan(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}

// RuleFind 规则查询
func (w *WeixinRequest) RuleFind(ruleId int64) (*bmodel.Weixin_RuleModel, error) {
	field := `a.rule_id,a.account_id,a.name,a.keywords,a.request_type,a.is_like_query,a.is_default,a.sort_id,a.create_time,a.update_time,b.content_id,b.title,b.content,b.link_url,b.img_url,b.media_url,b.media_hd_url`
	sqlSearch := fmt.Sprintf(`SELECT %s FROM weixin_request_rule a LEFT JOIN weixin_request_content b ON a.rule_id=b.rule_id WHERE a.rule_id=?`, field)
	_, ruleDo := query.WeixinRequestRuleDo()
	var ruleFind *bmodel.Weixin_RuleModel
	if err := ruleDo.Debug().UnderlyingDB().Raw(sqlSearch, ruleId).Scan(&ruleFind).Error; err != nil {
		return nil, err
	}
	return ruleFind, nil
}

// RuleSaveSortId 规则排序
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

// RuleDestory 规则删除
func (w *WeixinRequest) RuleDestory(ruleId int64) error {
	ruleMdl, ruleDo := query.WeixinRequestRuleDo()
	contentMdl, contentDo := query.WeixinRequestContentDo()
	contentDo.Where(contentMdl.RuleID.Eq(ruleId)).Delete()
	ruleDo.Where(ruleMdl.RuleID.Eq(ruleId)).Delete()
	return nil
}

// RuleSave 规则保存
func (w *WeixinRequest) RuleSave(mdl *bmodel.Weixin_RuleModel) error {
	ruleMdl, ruleDo := query.WeixinRequestRuleDo()
	contentMdl, contentDo := query.WeixinRequestContentDo()
	if mdl.RuleID <= 0 {
		ruleModel := &model.WeixinRequestRule{
			AccountID:    mdl.AccountID,
			Keywords:     mdl.Keywords,
			Name:         mdl.Name,
			RequestType:  mdl.RequestType,
			ResponseType: 0,
			IsLikeQuery:  mdl.IsLikeQuery,
			IsDefault:    mdl.IsDefault,
			SortID:       mdl.SortID,
			CreateTime:   time.Now(),
			UpdateTime:   time.Now(),
		}
		if err := ruleDo.Create(ruleModel); err != nil {
			return err
		}
		contentModel := &model.WeixinRequestContent{
			RuleID:     ruleModel.RuleID,
			Title:      mdl.Title,
			Content:    mdl.Content,
			LinkURL:    mdl.LinkURL,
			ImgURL:     mdl.ImgURL,
			MediaURL:   mdl.MediaURL,
			MediaHdURL: mdl.MediaHdURL,
			SortID:     mdl.SortID,
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		}
		if err := contentDo.Create(contentModel); err != nil {
			return err
		}
	} else {
		ruleModel := map[string]interface{}{
			ruleMdl.AccountID.ColumnName().String():   mdl.AccountID,
			ruleMdl.Keywords.ColumnName().String():    mdl.Keywords,
			ruleMdl.Name.ColumnName().String():        mdl.Name,
			ruleMdl.RequestType.ColumnName().String(): mdl.RequestType,
			// ruleMdl.ResponseType.ColumnName().String(): mdl.ResponseType,
			ruleMdl.IsLikeQuery.ColumnName().String(): mdl.IsLikeQuery,
			ruleMdl.IsDefault.ColumnName().String():   mdl.IsDefault,
			ruleMdl.SortID.ColumnName().String():      mdl.SortID,
			ruleMdl.UpdateTime.ColumnName().String():  time.Now(),
		}
		if _, err := ruleDo.Where(ruleMdl.RuleID.Eq(mdl.RuleID)).UpdateColumns(&ruleModel); err != nil {
			return err
		}
		contentModel := map[string]interface{}{
			contentMdl.Title.ColumnName().String():      mdl.Title,
			contentMdl.Content.ColumnName().String():    mdl.Content,
			contentMdl.LinkURL.ColumnName().String():    mdl.LinkURL,
			contentMdl.ImgURL.ColumnName().String():     mdl.ImgURL,
			contentMdl.MediaURL.ColumnName().String():   mdl.MediaURL,
			contentMdl.MediaHdURL.ColumnName().String(): mdl.MediaHdURL,
			contentMdl.SortID.ColumnName().String():     mdl.SortID,
			contentMdl.UpdateTime.ColumnName().String(): time.Now(),
		}
		if _, err := contentDo.Where(contentMdl.ContentID.Eq(mdl.ContentID)).UpdateColumns(&contentModel); err != nil {
			return err
		}
	}
	return nil
}

// PictureSave 图片保存
func (w *WeixinRequest) PictureSave(mdl *bmodel.Weixin_PictureModel) error {
	ruleMdl, ruleDo := query.WeixinRequestRuleDo()
	if mdl.RuleID > 0 {
		// 删除图片回复
		sql := `DELETE FROM weixin_request_content WHERE rule_id=?;`
		if err := ruleDo.Debug().UnderlyingDB().Exec(sql, mdl.RuleID).Error; err != nil {
			return err
		}
	}
	if mdl.RuleID <= 0 {
		ruleModel := &model.WeixinRequestRule{
			AccountID:    mdl.AccountID,
			Name:         mdl.Name,
			Keywords:     mdl.Keywords,
			RequestType:  mdl.RequestType,
			ResponseType: mdl.ResponseType,
			IsLikeQuery:  mdl.IsLikeQuery,
			IsDefault:    mdl.IsDefault,
			SortID:       mdl.SortID,
			CreateTime:   time.Now(),
			UpdateTime:   time.Now(),
		}
		if err := ruleDo.Create(ruleModel); err != nil {
			return err
		}
		mdl.RuleID = ruleModel.RuleID
	} else {
		ruleModel := map[string]interface{}{
			ruleMdl.AccountID.ColumnName().String():    mdl.AccountID,
			ruleMdl.Name.ColumnName().String():         mdl.Name,
			ruleMdl.Keywords.ColumnName().String():     mdl.Keywords,
			ruleMdl.RequestType.ColumnName().String():  mdl.RequestType,
			ruleMdl.ResponseType.ColumnName().String(): mdl.ResponseType,
			ruleMdl.IsLikeQuery.ColumnName().String():  mdl.IsLikeQuery,
			ruleMdl.IsDefault.ColumnName().String():    mdl.IsDefault,
			ruleMdl.SortID.ColumnName().String():       mdl.SortID,
			ruleMdl.UpdateTime.ColumnName().String():   time.Now(),
		}
		if _, err := ruleDo.Where(ruleMdl.RuleID.Eq(mdl.RuleID)).UpdateColumns(&ruleModel); err != nil {
			return err
		}
	}
	{
		// 插入图片回复
		if mdl.ImageReply != nil && len(mdl.ImageReply) > 0 {
			_, contentDo := query.WeixinRequestContentDo()
			for _, v := range mdl.ImageReply {
				contentModel := &model.WeixinRequestContent{
					RuleID:     mdl.RuleID,
					Title:      v.Title,
					Content:    v.Content,
					LinkURL:    v.LinkURL,
					ImgURL:     v.ImgURL,
					MediaURL:   v.MediaURL,
					MediaHdURL: v.MediaHdURL,
					SortID:     v.SortID,
					CreateTime: time.Now(),
					UpdateTime: time.Now(),
				}
				if err := contentDo.Create(contentModel); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
