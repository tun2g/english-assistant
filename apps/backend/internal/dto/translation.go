package dto

// TranslateTextsRequest represents a request to translate multiple texts
type TranslateTextsRequest struct {
	Texts      []string `json:"texts" binding:"required"`
	SourceLang string   `json:"sourceLang"` // auto-detect if empty
	TargetLang string   `json:"targetLang" binding:"required"`
}

// TranslateTextsResponse represents the response with translated texts
type TranslateTextsResponse struct {
	Translations []string `json:"translations"`
	SourceLang   string   `json:"sourceLang"` // detected or provided
	TargetLang   string   `json:"targetLang"`
}