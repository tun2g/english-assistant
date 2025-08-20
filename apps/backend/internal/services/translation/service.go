package translation

import (
	"context"
	"fmt"
	"strings"

	"app-backend/internal/logger"
	"app-backend/internal/types"
	"app-backend/pkg/gemini"
)

// Service implements translation functionality using Google Gemini
type Service struct {
	geminiService *gemini.Service
	logger        *logger.Logger
}

// Config holds configuration for translation service
type Config struct {
	GeminiAPIKey string
	Logger       *logger.Logger
}

// NewService creates a new translation service
func NewService(config *Config) (*Service, error) {
	if config.GeminiAPIKey == "" {
		return nil, fmt.Errorf("gemini API key is required for translation service")
	}

	// Create Gemini service with config
	geminiConfig := &gemini.Config{
		APIKey: config.GeminiAPIKey,
		Logger: config.Logger.Zap(),
	}

	geminiService, err := gemini.NewServiceWithConfig(geminiConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create gemini service: %w", err)
	}

	return &Service{
		geminiService: geminiService,
		logger:        config.Logger,
	}, nil
}

// TranslateTexts translates an array of texts to the target language
func (s *Service) TranslateTexts(ctx context.Context, texts []string, targetLang string, sourceLang string) ([]string, error) {
	if len(texts) == 0 {
		return []string{}, nil
	}

	// Mock translation implementation - temporarily disabled Gemini service
	translations := make([]string, len(texts))
	for i, text := range texts {
		// Format: [TARGET_LANG] original_text - to clearly show it's mock data
		translations[i] = fmt.Sprintf("[%s] %s", strings.ToUpper(targetLang), text)
	}

	return translations, nil

	// Original Gemini implementation - commented out for reuse later
	// // Convert texts to transcript segments for Gemini service compatibility
	// segments := make([]types.TranscriptSegment, len(texts))
	// for i, text := range texts {
	// 	segments[i] = types.TranscriptSegment{
	// 		Text:      text,
	// 		StartTime: types.MillisecondDuration(0),
	// 		EndTime:   types.MillisecondDuration(0),
	// 	}
	// }

	// // Use Gemini service to translate segments
	// translatedSegments, err := s.geminiService.TranslateSegments(ctx, segments, targetLang, sourceLang)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to translate texts: %w", err)
	// }

	// // Extract translated texts from segments
	// translations := make([]string, len(translatedSegments))
	// for i, segment := range translatedSegments {
	// 	translations[i] = segment.TranslatedText
	// }

	// return translations, nil
}

// DetectLanguage detects the language of the given text
func (s *Service) DetectLanguage(ctx context.Context, text string) (string, error) {
	// Mock language detection - return English as default
	return "en", nil
	
	// Original Gemini implementation - commented out for reuse later
	// return s.geminiService.DetectLanguage(ctx, text)
}

// GetSupportedLanguages returns list of supported translation languages
func (s *Service) GetSupportedLanguages() []types.Language {
	return s.geminiService.GetSupportedLanguages()
}

// Close closes the translation service and cleans up resources
func (s *Service) Close() error {
	if s.geminiService != nil {
		return s.geminiService.Close()
	}
	return nil
}