package translation

import (
	"context"

	"app-backend/internal/types"
)

// ServiceInterface defines the contract for translation services
type ServiceInterface interface {
	// TranslateTexts translates an array of texts to the target language
	TranslateTexts(ctx context.Context, texts []string, targetLang string, sourceLang string) ([]string, error)
	
	// DetectLanguage detects the language of the given text
	DetectLanguage(ctx context.Context, text string) (string, error)
	
	// GetSupportedLanguages returns list of supported translation languages
	GetSupportedLanguages() []types.Language
	
	// Close closes the translation service and cleans up resources
	Close() error
}