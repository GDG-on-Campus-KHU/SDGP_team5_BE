// db/model/user.go

package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type CountryLang string

const (
	CountryLangKorean	CountryLang = "ko"
	CountryLangEnglish	CountryLang = "en"
	CountryLangJapanese CountryLang = "ja"
	CountryLangChinese  CountryLang = "zh"
	CountryLangGerman   CountryLang = "de"
	CountryLangFrench   CountryLang = "fr"
	CountryLangSpanish  CountryLang = "es"
)


type User struct {
	ID      primitive.ObjectID	`bson:"_id,omitempty" json:"id,omitempty"`	// MongoDB ObjectId
	UserID  int					`bson:"user_id" json:"user_id"`				// auto-increment integer
	Name    string				`bson:"name" json:"name"`
	Email   string				`bson:"email" json:"email"`
	AppLang CountryLang        	`bson:"app_lang" json:"app_lang"`			// app 사용 언어 (ex. "ko", "en")

	CountryCode string			`bson:"country_code" json:"country_code"`	// 여행 중인 국가 코드 (ex. "KR", "US")
	InfoID      int				`bson:"info_id" json:"info_id"`				// ref. struct 'MedicalInfo'

	GroupIDs  []int				`bson:"group_ids" json:"group_ids"`
	Favorites []int				`bson:"favorites" json:"favorites"`			// situation index
}