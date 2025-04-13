// language/models.go

package language

// request
type TranslationRequest struct {
	Text string `json:"text"`
	// Text       []string `json:"text"
	TargetLang string `json:"target_lang"` // ex: "ko", "en", "fr"
}

// response
type TranslationResponse struct {
	TranslatedText string `json:"translated_text"`
	// TranslatedText []string `json:"translated_text"`
}
