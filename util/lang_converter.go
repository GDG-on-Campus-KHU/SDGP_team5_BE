// util/lang_converter.go

package util

import (
	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/db/model"
)


func AppLangToLangCode(appLang model.CountryLang) string {
	switch appLang {
		case model.CountryLangKorean:
			return "ko-KR"
		case model.CountryLangEnglish:
			return "en-US"
		case model.CountryLangJapanese:
			return "ja-JP"
		case model.CountryLangChinese:
			return "zh-CN"
		case model.CountryLangGerman:
			return "de-DE"
		case model.CountryLangFrench:
			return "fr-FR"
		case model.CountryLangSpanish:
			return "es-ES"
		default:
			return "en-US"
	}
}