package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/beego/beego/v2/core/logs"

	"haedu.gov.cn/cms/app/cms/domain"
	"haedu.gov.cn/cms/app/cms/mapper"
	service_model "haedu.gov.cn/cms/app/cms/service/model"
)

// WeixinRequest 微信请求管理服务
type WeixinRequest struct{}

// NewWeixinRequest 创建微信请求服务实例
func NewWeixinRequest() *WeixinRequest {
	return &WeixinRequest{}
}

// ContentFindSubscribeOrDefault 获取关注回复与默认回复
// @param accountId 公众号ID
// @param requestType 请求类型
// @return *service_model.Weixin_ContentSubscribeOrDefaultModel 回复模型, error
func (*WeixinRequest) ContentFindSubscribeOrDefault(accountId int64, requestType int32) (*service_model.Weixin_ContentSubscribeOrDefaultModel, error) {
	_, ruleDo := mapper.WeixinRequestRuleDo()

	sql1 := fmt.Sprintf(`SELECT a.* FROM weixin_request_content a LEFT JOIN weixin_request_rule b ON a.rule_id=b.rule_id WHERE b.account_id=%d AND b.request_type=%d AND b.response_type=%d limit 1`, accountId, requestType, 1)
	sql2 := fmt.Sprintf(`SELECT a.* FROM weixin_request_content a LEFT JOIN weixin_request_rule b ON a.rule_id=b.rule_id WHERE b.account_id=%d AND b.request_type=%d AND b.response_type=%d ORDER BY a.sort_id`, accountId, requestType, 2)
	sql3 := fmt.Sprintf(`SELECT a.* FROM weixin_request_content a LEFT JOIN weixin_request_rule b ON a.rule_id=b.rule_id WHERE b.account_id=%d AND b.request_type=%d AND b.response_type=%d limit 1`, accountId, requestType, 3)

	model := &service_model.Weixin_ContentSubscribeOrDefaultModel{
		AccountID:   accountId,
		RequestType: requestType,
		TextReply:   &domain.WeixinRequestContent{},
		ImageReply:  []*domain.WeixinRequestContent{},
		SoundReply:  &domain.WeixinRequestContent{},
	}

	ruleDo.UnderlyingDB().Raw(sql1).Scan(&model.TextReply)
	ruleDo.UnderlyingDB().Raw(sql2).Scan(&model.ImageReply)
	ruleDo.UnderlyingDB().Raw(sql3).Scan(&model.SoundReply)

	return model, nil
}

// ContentSaveSubscribeOrDefault 保存关注回复与默认回复
// @param input 回复内容模型
// @return error 错误信息
func (*WeixinRequest) ContentSaveSubscribeOrDefault(input *service_model.Weixin_ContentSubscribeOrDefaultModel) error {
	if input.AccountID <= 0 {
		return errors.New("accountId is empty")
	}

	_, ruleDo := mapper.WeixinRequestRuleDo()
	_, contentDo := mapper.WeixinRequestContentDo()

	// 删除文本回复、图片回复、语音回复
	{
		sql := `DELETE FROM weixin_request_content WHERE rule_id IN (SELECT rule_id FROM weixin_request_rule WHERE account_id=? AND request_type=? AND response_type=?);`
		ruleDo.Debug().UnderlyingDB().Exec(sql, input.AccountID, input.RequestType, 1)
		ruleDo.Debug().UnderlyingDB().Exec(sql, input.AccountID, input.RequestType, 2)
		ruleDo.Debug().UnderlyingDB().Exec(sql, input.AccountID, input.RequestType, 3)
	}

	// 保存文本回复
	{
		if input.TextReply != nil && input.TextReply.Content != "" {
			input.TextReply.CreateTime = time.Now()
			input.TextReply.UpdateTime = time.Now()
			if input.TextReply.RuleID <= 0 {
				ruleMdl := &domain.WeixinRequestRule{
					AccountID:    input.AccountID,
					RequestType:  input.RequestType,
					ResponseType: 1,
					Name:         "文本回复",
				}
				if err := ruleDo.Save(ruleMdl); err != nil {
					logs.Error("保存文本回复规则失败: accountId=%d, error=%v", input.AccountID, err)
					return err
				}
				input.TextReply.RuleID = ruleMdl.RuleID
			}
			input.TextReply.LinkURL = ""
			input.TextReply.ImgURL = ""
			input.TextReply.MediaURL = ""
			input.TextReply.MediaHdURL = ""
			if err := contentDo.Save(input.TextReply); err != nil {
				logs.Error("保存文本回复内容失败: accountId=%d, error=%v", input.AccountID, err)
				return err
			}
		} else {
			sql := `DELETE FROM weixin_request_rule WHERE account_id=? AND request_type=? AND response_type=?;`
			ruleDo.Debug().UnderlyingDB().Exec(sql, input.AccountID, input.RequestType, 1)
		}
	}

	// 保存图片回复
	{
		if len(input.ImageReply) > 0 {
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
					ruleMdl := &domain.WeixinRequestRule{
						AccountID:    input.AccountID,
						RequestType:  input.RequestType,
						ResponseType: 2,
						Name:         "图片回复",
					}
					if err := ruleDo.Save(ruleMdl); err != nil {
						logs.Error("保存图片回复规则失败: accountId=%d, error=%v", input.AccountID, err)
						return err
					}
					ruleID = ruleMdl.RuleID
					v.RuleID = ruleID
				}
				if err := contentDo.Save(v); err != nil {
					logs.Error("保存图片回复内容失败: accountId=%d, error=%v", input.AccountID, err)
					return err
				}
			}
		} else {
			sql := `DELETE FROM weixin_request_rule WHERE account_id=? AND request_type=? AND response_type=?;`
			ruleDo.Debug().UnderlyingDB().Exec(sql, input.AccountID, input.RequestType, 2)
		}
	}

	// 保存语音回复
	{
		if input.SoundReply != nil && input.SoundReply.MediaURL != "" {
			input.SoundReply.CreateTime = time.Now()
			input.SoundReply.UpdateTime = time.Now()
			if input.SoundReply.RuleID <= 0 {
				ruleMdl := &domain.WeixinRequestRule{
					AccountID:    input.AccountID,
					RequestType:  input.RequestType,
					ResponseType: 3,
					Name:         "语音回复",
				}
				if err := ruleDo.Save(ruleMdl); err != nil {
					logs.Error("保存语音回复规则失败: accountId=%d, error=%v", input.AccountID, err)
					return err
				}
				input.SoundReply.RuleID = ruleMdl.RuleID
			}
			if err := contentDo.Save(input.SoundReply); err != nil {
				logs.Error("保存语音回复内容失败: accountId=%d, error=%v", input.AccountID, err)
				return err
			}
		} else {
			sql := `DELETE FROM weixin_request_rule WHERE account_id=? AND request_type=? AND response_type=?;`
			ruleDo.Debug().UnderlyingDB().Exec(sql, input.AccountID, input.RequestType, 3)
		}
	}

	return nil
}

// RulePictureFind 规则文本回复查询
// @param ruleId 规则ID
// @return []*domain.WeixinRequestContent 内容列表, error
func (*WeixinRequest) RulePictureFind(ruleId int64) ([]*domain.WeixinRequestContent, error) {
	contentMdl, contentDo := mapper.WeixinRequestContentDo()
	return contentDo.Where(contentMdl.RuleID.Eq(ruleId)).Order(contentMdl.SortID).Find()
}

// RulePaginate 规则分页查询
// @param page 页码
// @param limit 每页数量
// @param accountId 公众号ID
// @param requestType 请求类型
// @return []*service_model.Weixin_RuleModel 规则列表, int64 总数, error
func (*WeixinRequest) RulePaginate(page, limit int, accountId int64, requestType int32) ([]*service_model.Weixin_RuleModel, int64, error) {
	_, ruleDo := mapper.WeixinRequestRuleDo()

	field := `a.rule_id,a.account_id,a.name,a.keywords,a.request_type,a.is_like_query,a.is_default,a.sort_id,a.create_time,a.update_time,b.content_id,b.title,b.content,b.link_url,b.img_url,b.media_url,b.media_hd_url`
	sqlCount := `SELECT COUNT(*) as Count FROM weixin_request_rule a LEFT JOIN weixin_request_content b ON a.rule_id=b.rule_id WHERE a.account_id=? AND a.request_type=?`
	sqlSearch := fmt.Sprintf(`SELECT %s FROM weixin_request_rule a LEFT JOIN weixin_request_content b ON a.rule_id=b.rule_id WHERE a.account_id=? AND a.request_type=? ORDER BY a.sort_id ASC`, field)

	var count int64
	if err := ruleDo.Debug().UnderlyingDB().Raw(sqlCount, accountId, requestType).Scan(&count).Error; err != nil {
		return nil, 0, err
	}

	lmt := fmt.Sprintf(" LIMIT %d,%d", (page-1)*limit, limit)
	sqlSearch += lmt

	var list []*service_model.Weixin_RuleModel
	if err := ruleDo.Debug().UnderlyingDB().Raw(sqlSearch, accountId, requestType).Scan(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, count, nil
}

// RulePaginateCount 规则汇总分页
// @param page 页码
// @param limit 每页数量
// @param accountId 公众号ID
// @param requestType 请求类型
// @return []*service_model.Weixin_RuleCountModel 规则列表, int64 总数, error
func (*WeixinRequest) RulePaginateCount(page, limit int, accountId int64, requestType int32) ([]*service_model.Weixin_RuleCountModel, int64, error) {
	_, ruleDo := mapper.WeixinRequestRuleDo()

	field := `a.rule_id,a.account_id,a.name,a.keywords,a.request_type,a.is_like_query,a.is_default,a.sort_id,a.create_time,a.update_time,(SELECT COUNT(1) AS count FROM weixin_request_content b WHERE b.rule_id=a.rule_id) AS count`
	sqlCount := `SELECT COUNT(*) as Count FROM weixin_request_rule a WHERE a.account_id=? AND a.request_type=?`
	sqlSearch := fmt.Sprintf(`SELECT %s FROM weixin_request_rule a WHERE a.account_id=? AND a.request_type=? ORDER BY a.sort_id ASC`, field)

	var count int64
	if err := ruleDo.Debug().UnderlyingDB().Raw(sqlCount, accountId, requestType).Scan(&count).Error; err != nil {
		return nil, 0, err
	}

	lmt := fmt.Sprintf(" LIMIT %d,%d", (page-1)*limit, limit)
	sqlSearch += lmt

	var list []*service_model.Weixin_RuleCountModel
	if err := ruleDo.Debug().UnderlyingDB().Raw(sqlSearch, accountId, requestType).Scan(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, count, nil
}

// RuleFind 规则查询
// @param ruleId 规则ID
// @return *service_model.Weixin_RuleModel 规则模型, error
func (*WeixinRequest) RuleFind(ruleId int64) (*service_model.Weixin_RuleModel, error) {
	field := `a.rule_id,a.account_id,a.name,a.keywords,a.request_type,a.is_like_query,a.is_default,a.sort_id,a.create_time,a.update_time,b.content_id,b.title,b.content,b.link_url,b.img_url,b.media_url,b.media_hd_url`
	sqlSearch := fmt.Sprintf(`SELECT %s FROM weixin_request_rule a LEFT JOIN weixin_request_content b ON a.rule_id=b.rule_id WHERE a.rule_id=?`, field)

	_, ruleDo := mapper.WeixinRequestRuleDo()
	var ruleFind *service_model.Weixin_RuleModel
	if err := ruleDo.Debug().UnderlyingDB().Raw(sqlSearch, ruleId).Scan(&ruleFind).Error; err != nil {
		return nil, err
	}
	return ruleFind, nil
}

// RuleSaveSortId 规则排序
// @param ruleId 规则ID
// @param sortId 排序值
// @return error
func (*WeixinRequest) RuleSaveSortId(ruleId int64, sortId int32) error {
	ruleMdl, ruleDo := mapper.WeixinRequestRuleDo()
	_, err := ruleDo.Where(ruleMdl.RuleID.Eq(ruleId)).UpdateColumns(
		map[string]interface{}{
			ruleMdl.SortID.ColumnName().String():     sortId,
			ruleMdl.UpdateTime.ColumnName().String(): time.Now(),
		},
	)
	return err
}

// RuleDestory 规则删除
// @param ruleId 规则ID
// @return error
func (*WeixinRequest) RuleDestory(ruleId int64) error {
	ruleMdl, ruleDo := mapper.WeixinRequestRuleDo()
	contentMdl, contentDo := mapper.WeixinRequestContentDo()

	_, err := contentDo.Where(contentMdl.RuleID.Eq(ruleId)).Delete()
	if err != nil {
		logs.Error("删除微信回复内容失败: ruleId=%d, error=%v", ruleId, err)
	}

	_, err = ruleDo.Where(ruleMdl.RuleID.Eq(ruleId)).Delete()
	if err != nil {
		logs.Error("删除微信回复规则失败: ruleId=%d, error=%v", ruleId, err)
	}
	return err
}

// RuleSave 规则保存
// @param mdl 规则模型
// @return error
func (*WeixinRequest) RuleSave(mdl *service_model.Weixin_RuleModel) error {
	ruleMdl, ruleDo := mapper.WeixinRequestRuleDo()
	contentMdl, contentDo := mapper.WeixinRequestContentDo()

	if mdl.RuleID <= 0 {
		ruleModel := &domain.WeixinRequestRule{
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
			logs.Error("创建微信回复规则失败: accountId=%d, error=%v", mdl.AccountID, err)
			return err
		}
		contentModel := &domain.WeixinRequestContent{
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
			logs.Error("创建微信回复内容失败: accountId=%d, error=%v", mdl.AccountID, err)
			return err
		}
	} else {
		ruleModel := map[string]interface{}{
			ruleMdl.AccountID.ColumnName().String():     mdl.AccountID,
			ruleMdl.Keywords.ColumnName().String():      mdl.Keywords,
			ruleMdl.Name.ColumnName().String():         mdl.Name,
			ruleMdl.RequestType.ColumnName().String():   mdl.RequestType,
			ruleMdl.IsLikeQuery.ColumnName().String():   mdl.IsLikeQuery,
			ruleMdl.IsDefault.ColumnName().String():     mdl.IsDefault,
			ruleMdl.SortID.ColumnName().String():        mdl.SortID,
			ruleMdl.UpdateTime.ColumnName().String():    time.Now(),
		}
		if _, err := ruleDo.Where(ruleMdl.RuleID.Eq(mdl.RuleID)).UpdateColumns(&ruleModel); err != nil {
			logs.Error("更新微信回复规则失败: ruleId=%d, error=%v", mdl.RuleID, err)
			return err
		}

		contentModel := map[string]interface{}{
			contentMdl.Title.ColumnName().String():      mdl.Title,
			contentMdl.Content.ColumnName().String():    mdl.Content,
			contentMdl.LinkURL.ColumnName().String():    mdl.LinkURL,
			contentMdl.ImgURL.ColumnName().String():     mdl.ImgURL,
			contentMdl.MediaURL.ColumnName().String():    mdl.MediaURL,
			contentMdl.MediaHdURL.ColumnName().String(): mdl.MediaHdURL,
			contentMdl.SortID.ColumnName().String():     mdl.SortID,
			contentMdl.UpdateTime.ColumnName().String():  time.Now(),
		}
		if _, err := contentDo.Where(contentMdl.ContentID.Eq(mdl.ContentID)).UpdateColumns(&contentModel); err != nil {
			logs.Error("更新微信回复内容失败: contentId=%d, error=%v", mdl.ContentID, err)
			return err
		}
	}
	return nil
}

// PictureSave 图片保存
// @param mdl 图片模型
// @return error
func (*WeixinRequest) PictureSave(mdl *service_model.Weixin_PictureModel) error {
	ruleMdl, ruleDo := mapper.WeixinRequestRuleDo()

	if mdl.RuleID > 0 {
		sql := `DELETE FROM weixin_request_content WHERE rule_id=?;`
		if err := ruleDo.Debug().UnderlyingDB().Exec(sql, mdl.RuleID).Error; err != nil {
			logs.Error("清空图片回复内容失败: ruleId=%d, error=%v", mdl.RuleID, err)
			return err
		}
	}

	if mdl.RuleID <= 0 {
		ruleModel := &domain.WeixinRequestRule{
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
			logs.Error("创建图片回复规则失败: accountId=%d, error=%v", mdl.AccountID, err)
			return err
		}
		mdl.RuleID = ruleModel.RuleID
	} else {
		ruleModel := map[string]interface{}{
			ruleMdl.AccountID.ColumnName().String():     mdl.AccountID,
			ruleMdl.Name.ColumnName().String():         mdl.Name,
			ruleMdl.Keywords.ColumnName().String():      mdl.Keywords,
			ruleMdl.RequestType.ColumnName().String():   mdl.RequestType,
			ruleMdl.ResponseType.ColumnName().String():  mdl.ResponseType,
			ruleMdl.IsLikeQuery.ColumnName().String():   mdl.IsLikeQuery,
			ruleMdl.IsDefault.ColumnName().String():     mdl.IsDefault,
			ruleMdl.SortID.ColumnName().String():        mdl.SortID,
			ruleMdl.UpdateTime.ColumnName().String():    time.Now(),
		}
		if _, err := ruleDo.Where(ruleMdl.RuleID.Eq(mdl.RuleID)).UpdateColumns(&ruleModel); err != nil {
			logs.Error("更新图片回复规则失败: ruleId=%d, error=%v", mdl.RuleID, err)
			return err
		}
	}

	// 插入图片回复
	if len(mdl.ImageReply) > 0 {
		_, contentDo := mapper.WeixinRequestContentDo()
		for _, v := range mdl.ImageReply {
			contentModel := &domain.WeixinRequestContent{
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
				logs.Error("创建图片回复内容失败: ruleId=%d, error=%v", mdl.RuleID, err)
				return err
			}
		}
	}

	return nil
}