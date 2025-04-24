// db/model/country.go

package model

type Country struct {
	CountryCode   string `bson:"country_code" json:"country_code"`
	CountryLang   string `bson:"country_lang" json:"country_lang"`
	CountryName   string `bson:"country_name" json:"country_name"`
	LocalEmerCall string `bson:"local_emer_call" json:"local_emer_call"`
	IntlEmerCall  string `bson:"intl_emer_call" json:"intl_emer_call"`
}
