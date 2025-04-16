// situation/model/situation.go

package model

type MultilingualText map[string]string
type MultilingualArray map[string][]string

type ActionStep struct {
	Step    string   `json:"step" bson:"step"`
	Details []string `json:"details" bson:"details"`
}

type MultiLangActions map[string][]ActionStep

type Situation struct {
	Index       int               `bson:"index" json:"index"`
	Slug        string            `bson:"slug" json:"slug"`
	Emoji       string            `bson:"emoji" json:"emoji"`
	EmerTitle   MultilingualText  `bson:"emer_title" json:"emer_title"`
	Description MultilingualArray `bson:"description" json:"description"`
	Actions     MultiLangActions  `bson:"actions" json:"actions"`
}
