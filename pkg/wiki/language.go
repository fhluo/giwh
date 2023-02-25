package wiki

import (
	"fmt"
	"github.com/samber/lo"
)

type Language struct {
	Key   string
	Name  string
	Short string
}

func (lang Language) String() string {
	return fmt.Sprintf("%s(%s)", lang.Key, lang.Name)
}

func (lang Language) Wiki() Wiki {
	return Wiki{Language: lang}
}

var (
	SimplifiedChinese  = Language{Key: "zh-cn", Name: "简体中文", Short: "简"}
	TraditionalChinese = Language{Key: "zh-tw", Name: "繁體中文", Short: "繁"}
	German             = Language{Key: "de-de", Name: "Deutsch", Short: "DE"}
	English            = Language{Key: "en-us", Name: "English", Short: "EN"}
	Spanish            = Language{Key: "es-es", Name: "Español", Short: "ES"}
	French             = Language{Key: "fr-fr", Name: "Français", Short: "FR"}
	Indonesian         = Language{Key: "id-id", Name: "Indonesia", Short: "ID"}
	Japanese           = Language{Key: "ja-jp", Name: "日本語", Short: "JP"}
	Korean             = Language{Key: "ko-kr", Name: "한국어", Short: "KR"}
	Portuguese         = Language{Key: "pt-pt", Name: "Português", Short: "PT"}
	Russian            = Language{Key: "ru-ru", Name: "Pусский", Short: "RU"}
	Thai               = Language{Key: "th-th", Name: "ภาษาไทย", Short: "TH"}
	Vietnamese         = Language{Key: "vi-vn", Name: "Tiếng Việt", Short: "VN"}

	AvailableLanguages = []Language{
		SimplifiedChinese, TraditionalChinese, German, English, Spanish, French, Indonesian, Japanese, Korean, Portuguese, Russian, Thai, Vietnamese,
	}

	Wikis = lo.Map(AvailableLanguages, func(lang Language, _ int) Wiki {
		return lang.Wiki()
	})
)
