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


func CountryCodeToLangCode(countryCode string) string {
	switch countryCode {
		case "KR":
			return "ko-KR"
		case "US":
			return "en-US"
		case "JP":
			return "ja-JP"
		case "CN":
			return "zh-CN"
		case "DE":
			return "de-DE"
		case "FR":
			return "fr-FR"
		case "ES":
			return "es-ES"
		default:
			return "en-US"
	}
}