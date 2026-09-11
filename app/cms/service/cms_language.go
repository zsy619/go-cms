package service

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"haedu.gov.cn/cms/app/cms/domain"
	"haedu.gov.cn/cms/app/db"
)

// CmsLanguage 语言设置服务
type CmsLanguage struct{}

// NewCmsLanguage 创建语言服务实例
func NewCmsLanguage() *CmsLanguage {
	return &CmsLanguage{}
}

// LanguageAll 获取所有语言（不分页）
func (svc *CmsLanguage) LanguageAll(name, code string) ([]*domain.CmsLanguage, error) {
	var list []*domain.CmsLanguage

	query := db.CmsDatabase.Where("deleted = ?", false)

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if code != "" {
		query = query.Where("code LIKE ?", "%"+code+"%")
	}

	err := query.Order("sort_id ASC, language_id ASC").Find(&list).Error
	return list, err
}

// LanguagePaginate 分页查询语言列表（保留兼容）
func (svc *CmsLanguage) LanguagePaginate(page, limit int, name, code string) ([]*domain.CmsLanguage, int64, error) {
	var list []*domain.CmsLanguage
	var total int64

	query := db.CmsDatabase.Where("deleted = ?", false)

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if code != "" {
		query = query.Where("code LIKE ?", "%"+code+"%")
	}

	query.Model(&domain.CmsLanguage{}).Count(&total)

	offset := (page - 1) * limit
	err := query.Order("sort_id ASC, language_id ASC").Offset(offset).Limit(limit).Find(&list).Error

	return list, total, err
}

// LanguageOne 获取单条语言
func (svc *CmsLanguage) LanguageOne(id int64) (*domain.CmsLanguage, error) {
	var lang domain.CmsLanguage
	err := db.CmsDatabase.Where("language_id = ? AND deleted = ?", id, false).First(&lang).Error
	if err != nil {
		return nil, err
	}
	return &lang, nil
}

// LanguageOneByCode 根据代码获取语言
func (svc *CmsLanguage) LanguageOneByCode(code string) (*domain.CmsLanguage, error) {
	var lang domain.CmsLanguage
	err := db.CmsDatabase.Where("code = ? AND deleted = ?", code, false).First(&lang).Error
	if err != nil {
		return nil, err
	}
	return &lang, nil
}

// LanguageSave 保存语言
func (svc *CmsLanguage) LanguageSave(mdl *domain.CmsLanguage) error {
	if mdl.Name == "" {
		return errors.New("语言名称不能为空")
	}
	if mdl.Code == "" {
		return errors.New("语言代码不能为空")
	}

	// 检查代码唯一性
	var count int64
	db.CmsDatabase.Model(&domain.CmsLanguage{}).Where("code = ? AND deleted = ? AND language_id != ?", mdl.Code, false, mdl.LanguageID).Count(&count)
	if count > 0 {
		return errors.New("语言代码已存在")
	}

	// 如果设置为默认，取消其他默认
	if mdl.IsDefault {
		db.CmsDatabase.Model(&domain.CmsLanguage{}).Where("is_default = ?", true).Update("is_default", false)
	}

	mdl.UpdateTime = time.Now()

	if mdl.LanguageID <= 0 {
		mdl.CreateTime = time.Now()
		return db.CmsDatabase.Save(mdl).Error
	}

	return db.CmsDatabase.Save(mdl).Error
}

// LanguageDelete 删除语言
func (svc *CmsLanguage) LanguageDelete(ids string) error {
	if ids == "" {
		return errors.New("参数丢失")
	}
	idArr := strings.Split(ids, ",")
	for _, id := range idArr {
		id64, _ := strconv.ParseInt(id, 0, 64)
		if id != "" {
			db.CmsDatabase.Model(&domain.CmsLanguage{}).Where("language_id = ?", id64).Updates(map[string]interface{}{
				"deleted":     true,
				"update_time": time.Now(),
			})
		}
	}
	return nil
}

// LanguageSort 排序
func (svc *CmsLanguage) LanguageSort(id int64, sortId int32) error {
	return db.CmsDatabase.Model(&domain.CmsLanguage{}).Where("language_id = ?", id).Updates(map[string]interface{}{
		"sort_id":     sortId,
		"update_time": time.Now(),
	}).Error
}

// LanguageStatus 状态切换
func (svc *CmsLanguage) LanguageStatus(id int64, status int32) error {
	return db.CmsDatabase.Model(&domain.CmsLanguage{}).Where("language_id = ?", id).Updates(map[string]interface{}{
		"status":      status,
		"update_time": time.Now(),
	}).Error
}

// LanguageInitData 初始化语言数据
func (svc *CmsLanguage) LanguageInitData() (int, error) {
	languages := []struct {
		Name        string
		Code        string
		Icon        string
		SortID      int32
		Description string
	}{
		// ========== 亚洲语言 ==========
		// 中国大陆及港澳台
		{"简体中文", "zh-CN", "🇨🇳", 1, "简体中文"},
		{"繁體中文", "zh-TW", "🇭🇰", 2, "繁体中文"},
		{"廣東話", "zh-HK", "🇭🇰", 3, "粤语"},
		{" Português", "zh-MO", "🇲🇴", 4, "澳门语"},
		// 日本韩国
		{"日本語", "ja", "🇯🇵", 5, "日语"},
		{"한국어", "ko", "🇰🇷", 6, "韩语"},
		{" 조선어", "ko-KP", "🇰🇵", 7, "朝鲜语"},
		// 东南亚
		{"Tiếng Việt", "vi", "🇻🇳", 8, "越南语"},
		{"ภาษาไทย", "th", "🇹🇭", 9, "泰语"},
		{"Bahasa Indonesia", "id", "🇮🇩", 10, "印尼语"},
		{"Bahasa Melayu", "ms", "🇲🇾", 11, "马来语"},
		{"Filipino", "fil", "🇵🇭", 12, "菲律宾语"},
		{"Tagalog", "tl", "🇵🇭", 13, "他加禄语"},
		{"Cebuano", "ceb", "🇵🇭", 14, "宿务语"},
		{"Ilocano", "ilo", "🇵🇭", 15, "伊洛卡诺语"},
		{"Myanmar", "my", "🇲🇲", 16, "缅甸语"},
		{" Khmer", "km", "🇰🇭", 17, "高棉语"},
		{"Lao", "lo", "🇱🇦", 18, "老挝语"},
		{" Brunei", "ms-BN", "🇧🇳", 19, "文莱马来语"},
		// 南亚
		{"हिन्दी", "hi", "🇮🇳", 20, "印地语"},
		{"বাংলা", "bn", "🇧🇩", 21, "孟加拉语"},
		{"தமிழ்", "ta", "🇮🇳", 22, "泰米尔语"},
		{"తెలుగు", "te", "🇮🇳", 23, "泰卢固语"},
		{"मराठी", "mr", "🇮🇳", 24, "马拉地语"},
		{"ಕನ್ನಡ", "kn", "🇮🇳", 25, "卡纳达语"},
		{"മലയാളം", "ml", "🇮🇳", 26, "马拉雅拉姆语"},
		{"অসমীয়া", "as", "🇮🇳", 27, "阿萨姆语"},
		{"ଓଡ଼ା", "or", "🇮🇳", 28, "奥里亚语"},
		{"পাঞ্জাবী", "pa", "🇮🇳", 29, "旁遮普语"},
		{"نیپالی", "ne", "🇳🇵", 30, "尼泊尔语"},
		{"संस्कृत", "sa", "🇮🇳", 31, "梵语"},
		{"සිංහල", "si", "🇱🇰", 32, "僧伽罗语"},
		{"دَری", "prs", "🇦🇫", 33, "达里语"},
		{"پښتو", "ps", "🇦🇫", 34, "普什图语"},
		{"اردو", "ur", "🇵🇰", 35, "乌尔都语"},
		{"سنڌي", "sd", "🇵🇰", 36, "信德语"},
		{"কাশ্মীরি", "ks", "🇮🇳", 37, "克什米尔语"},
		{"মৈতৈ", "mni", "🇮🇳", 38, "曼尼普尔语"},
		// 中亚
		{"Қазақша", "kk", "🇰🇿", 39, "哈萨克语"},
		{"Кыргызча", "ky", "🇰🇬", 40, "吉尔吉斯语"},
		{"Тоҷикӣ", "tg", "🇹🇯", 41, "塔吉克语"},
		{"Ўзбекча", "uz", "🇺🇿", 42, "乌兹别克语"},
		{"Түркменче", "tk", "🇹🇲", 43, "土库曼语"},
		{"Монгол", "mn", "🇲🇳", 44, "蒙古语"},
		{"Монгол", "mn-MN", "🇲🇳", 45, "蒙古语"},
		// 其他亚洲语言
		{"العربية", "ar", "🇸🇦", 46, "阿拉伯语"},
		{"עברית", "he", "🇮🇱", 47, "希伯来语"},
		{"فارسی", "fa", "🇮🇷", 48, "波斯语"},
		{"Türkçe", "tr", "🇹🇷", 49, "土耳其语"},
		{"Azərbaycan", "az", "🇦🇿", 50, "阿塞拜疆语"},
		{"Հայերեն", "hy", "🇦🇲", 51, "亚美尼亚语"},
		{"ქართული", "ka", "🇬🇪", 52, "格鲁吉亚语"},
		// ========== 欧洲语言 ==========
		// 西欧
		{"English", "en", "🇺🇸", 53, "英语"},
		{"Français", "fr", "🇫🇷", 54, "法语"},
		{"Deutsch", "de", "🇩🇪", 55, "德语"},
		{"Nederlands", "nl", "🇳🇱", 56, "荷兰语"},
		{"Italiano", "it", "🇮🇹", 57, "意大利语"},
		{"Español", "es", "🇪🇸", 58, "西班牙语"},
		{"Português", "pt", "🇵🇹", 59, "葡萄牙语"},
		{"Català", "ca", "🇪🇸", 60, "加泰罗尼亚语"},
		{"Galego", "gl", "🇪🇸", 61, "加利西亚语"},
		{"Euskara", "eu", "🇪🇸", 62, "巴斯克语"},
		{"Occitan", "oc", "🇫🇷", 63, "奥克语"},
		{"Breton", "br", "🇫🇷", 64, "布列塔尼语"},
		{"Corsu", "co", "🇫🇷", 65, "科西嘉语"},
		{"Ladino", "lad", "🇹🇷", 66, "拉迪诺语"},
		{"Walon", "wa", "🇧🇪", 67, "瓦隆语"},
		{"Frysk", "fy", "🇳🇱", 68, "弗里斯兰语"},
		{"Lëtzebuergesch", "lb", "🇱🇺", 69, "卢森堡语"},
		// 北欧
		{"Svenska", "sv", "🇸🇪", 70, "瑞典语"},
		{"Norsk", "no", "🇳🇴", 71, "挪威语"},
		{"Dansk", "da", "🇩🇰", 72, "丹麦语"},
		{"Suomi", "fi", "🇫🇮", 73, "芬兰语"},
		{"Íslenska", "is", "🇮🇸", 74, "冰岛语"},
		{"Føroyskt", "fo", "🇫🇴", 75, "法罗语"},
		{"Grønlandsk", "kl", "🇬🇱", 76, "格陵兰语"},
		{"Sami", "se", "🇳🇴", 77, "萨米语"},
		{"North Sami", "se-NO", "🇳🇴", 78, "北萨米语"},
		{"Inari Sami", "smn", "🇫🇮", 79, "伊纳里萨米语"},
		{"Skolt Sami", "sms", "🇫🇮", 80, "斯科尔特萨米语"},
		// 东欧
		{"Русский", "ru", "🇷🇺", 81, "俄语"},
		{"Українська", "uk", "🇺🇦", 82, "乌克兰语"},
		{"Беларуская", "be", "🇧🇾", 83, "白俄罗斯语"},
		{"Polski", "pl", "🇵🇱", 84, "波兰语"},
		{"Čeština", "cs", "🇨🇿", 85, "捷克语"},
		{"Slovenčina", "sk", "🇸🇰", 86, "斯洛伐克语"},
		{"Македонски", "mk", "🇲🇰", 87, "马其顿语"},
		{"Български", "bg", "🇧🇬", 88, "保加利亚语"},
		{"Српски", "sr", "🇷🇸", 89, "塞尔维亚语"},
		{"Hrvatski", "hr", "🇭🇷", 90, "克罗地亚语"},
		{"Bosanski", "bs", "🇧🇦", 91, "波斯尼亚语"},
		{"Crnogorski", "me", "🇲🇪", 92, "黑山语"},
		{"Slovenščina", "sl", "🇸🇮", 93, "斯洛文尼亚语"},
		{"Română", "ro", "🇷🇴", 94, "罗马尼亚语"},
		{"Magyar", "hu", "🇭🇺", 95, "匈牙利语"},
		{"Lithuanian", "lt", "🇱🇹", 96, "立陶宛语"},
		{"Latviešu", "lv", "🇱🇻", 97, "拉脱维亚语"},
		{"Eesti", "et", "🇪🇪", 98, "爱沙尼亚语"},
		{"Greek", "el", "🇬🇷", 99, "希腊语"},
		{"Albanian", "sq", "🇦🇱", 100, "阿尔巴尼亚语"},
		// ========== 非洲语言 ==========
		{"العربية", "ar", "🇸🇦", 101, "阿拉伯语"},
		{"المصرى", "ar-EG", "🇪🇬", 102, "埃及阿拉伯语"},
		{"ال مغربية", "ar-MA", "🇲🇦", 103, "摩洛哥阿拉伯语"},
		{"Kiswahili", "sw", "🇹🇿", 104, "斯瓦希里语"},
		{"Hausa", "ha", "🇳🇬", 105, "豪萨语"},
		{"Yorùbá", "yo", "🇳🇬", 106, "约鲁巴语"},
		{"Igbo", "ig", "🇳🇬", 107, "伊博语"},
		{"Zulu", "zu", "🇿🇦", 108, "祖鲁语"},
		{"Xhosa", "xh", "🇿🇦", 109, "科萨语"},
		{"Afrikaans", "af", "🇿🇦", 110, "南非荷兰语"},
		{"Sesotho", "st", "🇱🇸", 111, "塞索托语"},
		{"Setswana", "tn", "🇧🇼", 112, "茨瓦纳语"},
		{"Northern Sotho", "nso", "🇿🇦", 113, "北索托语"},
		{"Shona", "sn", "🇿🇼", 114, "绍纳语"},
		{"Ndebele", "nd", "🇿🇼", 115, "恩德贝莱语"},
		{"Somali", "so", "🇸🇴", 116, "索马里语"},
		{"Amharic", "am", "🇪🇹", 117, "阿姆哈拉语"},
		{"Oromo", "om", "🇪🇹", 118, "奥罗莫语"},
		{"Tigrinya", "ti", "🇪🇷", 119, "提格里尼亚语"},
		{"Wolof", "wo", "🇸🇳", 120, "沃洛夫语"},
		{"Fula", "ff", "🇸🇳", 121, "富拉尼语"},
		{"Lingala", "ln", "🇨🇩", 122, "林加拉语"},
		{"Kongo", "kg", "🇨🇩", 123, "刚果语"},
		{"Luba-Katanga", "lu", "🇨🇩", 124, "卢巴语"},
		{"Swahili", "sw-TZ", "🇹🇿", 125, "坦桑尼亚斯瓦希里语"},
		{"Kinyarwanda", "rw", "🇷🇼", 126, "卢旺达语"},
		{"Kirundi", "rn", "🇧🇮", 127, "基隆迪语"},
		{"Chichewa", "ny", "🇲🇼", 128, "奇切瓦语"},
		{"Malagasy", "mg", "🇲🇬", 129, "马达加斯加语"},
		{"Sotho", "st", "🇱🇸", 130, "塞苏陀语"},
		{"Twi", "tw", "🇬🇭", 131, "契维语"},
		{"Wolof", "wo", "🇸🇳", 132, "沃洛夫语"},
		{"Mossi", "mos", "🇧🇫", 133, "莫西语"},
		{"Zarma", "dje", "🇳🇪", 134, "扎尔马语"},
		{"Kanuri", "kr", "🇳🇬", 135, "卡努里语"},
		{"Ibibio", "ibb", "🇳🇬", 136, "伊比比奥语"},
		{"Yoruba", "yo", "🇳🇬", 137, "约鲁巴语"},
		{"Pular", "ff", "🇬🇳", 138, "普拉尔语"},
		// ========== 美洲语言 ==========
		{"English", "en-US", "🇺🇸", 139, "美式英语"},
		{"Spanish", "es-US", "🇺🇸", 140, "美式西班牙语"},
		{"Français", "fr-CA", "🇨🇦", 141, "加拿大法语"},
		{"Português", "pt-BR", "🇧🇷", 142, "巴西葡萄牙语"},
		{"Nahuatl", "nah", "🇲🇽", 143, "纳瓦特尔语"},
		{"Quechua", "qu", "🇵🇪", 144, "克丘亚语"},
		{"Aymara", "ay", "🇧🇴", 145, "艾马拉语"},
		{"Guarani", "gn", "🇵🇾", 146, "瓜拉尼语"},
		{"Maya", "yua", "🇬🇹", 147, "玛雅语"},
		{"Mixtec", "mix", "🇲🇽", 148, "米斯特克语"},
		{"Zapotec", "zap", "🇲🇽", 149, "萨帕泰克语"},
		{"Inuktitut", "iu", "🇨🇦", 150, "因纽特语"},
		{"Cree", "cr", "🇨🇦", 151, "克里语"},
		{"Ojibwe", "oj", "🇺🇸", 152, "奥吉布韦语"},
		{"Navajo", "nv", "🇺🇸", 153, "纳瓦霍语"},
		{"Cherokee", "chr", "🇺🇸", 154, "切罗基语"},
		{"Lakota", "lkt", "🇺🇸", 155, "拉科塔语"},
		{"Hawaiian", "haw", "🇺🇸", 156, "夏威夷语"},
		{"Samoan", "sm", "🇼🇸", 157, "萨摩亚语"},
		{"Tongan", "to", "🇹🇴", 158, "汤加语"},
		{"Maori", "mi", "🇳🇿", 159, "毛利语"},
		{"Tahitian", "ty", "🇵🇫", 160, "塔希提语"},
		{"Fijian", "fj", "🇫🇯", 161, "斐济语"},
		// ========== 大洋洲语言 ==========
		{"Papuan", "tpi", "🇵🇬", 162, "巴布亚语"},
		{"Hiri Motu", "ho", "🇵🇬", 163, "希里莫图语"},
		{"Chamorro", "ch", "🇬🇺", 164, "查莫罗语"},
		{"Palauan", "pau", "🇵🇼", 165, "帕劳语"},
		{"Marshallese", "mh", "🇲🇭", 166, "马绍尔语"},
		{"Carolinian", "cal", "🇲🇵", 167, "卡罗琳语"},
		{"Nauru", "na", "🇳🇷", 168, "瑙鲁语"},
		{"Kiribati", "gil", "🇰🇮", 169, "基里巴斯语"},
		// ========== 编程语言/特殊用途 ==========
		{"English (Unicode)", "en-US", "🇺🇸", 170, "Unicode英语"},
		{"简体中文 (Unicode)", "zh-Hans", "🇨🇳", 171, "简体中文Unicode"},
		{"繁體中文 (Unicode)", "zh-Hant", "🇭🇰", 172, "繁体中文Unicode"},
		{"日本語 (Unicode)", "ja", "🇯🇵", 173, "日语Unicode"},
		{"한국어 (Unicode)", "ko", "🇰🇷", 174, "韩语Unicode"},
		// ========== 古典语言 ==========
		{"Latin", "la", "🇻🇦", 179, "拉丁语"},
		{"Ancient Greek", "grc", "🇬🇷", 180, "古希腊语"},
		{"Sanskrit", "sa", "🇮🇳", 181, "梵语"},
		{"Classical Arabic", "ar-SA", "🇸🇦", 182, "古典阿拉伯语"},
		{"Biblical Hebrew", "he", "🇮🇱", 183, "圣经希伯来语"},
		{"Old Church Slavonic", "cu", "🇷🇺", 184, "古教会斯拉夫语"},
		{"Classical Armenian", "xcl", "🇦🇲", 185, "古典亚美尼亚语"},
		{"Middle English", "enm", "🏴󠁧󠁢󠁥󠁮󠁧󠁿", 186, "中古英语"},
		{"Old French", "fro", "🇫🇷", 187, "古法语"},
		{"Old High German", "goh", "🇩🇪", 188, "古高地德语"},
		{"Old Norse", "non", "🇳🇴", 189, "古诺尔斯语"},
		{"Middle High German", "gmh", "🇩🇪", 190, "中古高地德语"},
		{"Middle Low German", "gml", "🇩🇪", 191, "中古低地德语"},
		{"Old Church Slavonic", "chu", "🇧🇬", 192, "古教会斯拉夫语"},
		{"Biblical Latin", "la", "🇻🇦", 193, "圣经拉丁语"},
		{"Vulgar Latin", "la", "🇻🇦", 194, "通俗拉丁语"},
		{"Sanskrit", "sa", "🇮🇳", 195, "梵语"},
		{"Pali", "pi", "🇮🇳", 196, "巴利语"},
		{"Prakrit", "pra", "🇮🇳", 197, "普拉克里特语"},
	}

	count := 0
	for _, lang := range languages {
		// 检查是否已存在
		var existing domain.CmsLanguage
		err := db.CmsDatabase.Where("code = ? AND deleted = ?", lang.Code, false).First(&existing).Error
		if err == nil {
			continue // 已存在，跳过
		}

		mdl := &domain.CmsLanguage{
			Name:        lang.Name,
			Code:        lang.Code,
			Icon:        lang.Icon,
			SortID:      lang.SortID,
			Description: lang.Description,
			Status:      1,
			IsDefault:   lang.SortID == 1,
		}
		mdl.CreateTime = time.Now()
		mdl.UpdateTime = time.Now()

		if err := db.CmsDatabase.Save(mdl).Error; err == nil {
			count++
		}
	}
	return count, nil
}
