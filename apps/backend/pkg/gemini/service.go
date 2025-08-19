package gemini

import (
	"context"
	"fmt"
	"strings"
	"time"

	"app-backend/internal/types"
	"github.com/google/generative-ai-go/genai"
	"go.uber.org/zap"
	"google.golang.org/api/option"
)

// Service implements translation functionality using Google Gemini
type Service struct {
	client   *genai.Client
	model    *genai.GenerativeModel
	logger   *zap.Logger
	apiKey   string
}

// Config holds configuration for Gemini service
type Config struct {
	APIKey    string
	ModelName string // Optional, defaults to "gemini-1.5-flash"
	Logger    *zap.Logger
}

// TranslationRequest represents a request to translate text
type TranslationRequest struct {
	Text       string `json:"text"`
	SourceLang string `json:"sourceLang,omitempty"`
	TargetLang string `json:"targetLang"`
	Context    string `json:"context,omitempty"` // Additional context for better translation
}

// TranslationResponse represents the response from translation
type TranslationResponse struct {
	OriginalText   string `json:"originalText"`
	TranslatedText string `json:"translatedText"`
	SourceLang     string `json:"sourceLang"`
	TargetLang     string `json:"targetLang"`
	Confidence     float64 `json:"confidence,omitempty"`
}

// NewServiceWithConfig creates a new Gemini translation service with config
func NewServiceWithConfig(config *Config) (*Service, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("gemini API key is required")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(config.APIKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create gemini client: %w", err)
	}

	modelName := config.ModelName
	if modelName == "" {
		modelName = "gemini-1.5-flash" // Default model
	}

	model := client.GenerativeModel(modelName)
	
	// Configure model for better translation performance
	model.SetTemperature(0.1) // Low temperature for consistent translations
	model.SetTopK(1)
	model.SetTopP(0.1)

	return &Service{
		client: client,
		model:  model,
		logger: config.Logger,
		apiKey: config.APIKey,
	}, nil
}

// NewService creates a new Gemini translation service (for container injection)
func NewService(apiKey string, logger *zap.Logger) *Service {
	if apiKey == "" {
		logger.Error("Gemini API key is required")
		// Return a service that will gracefully handle missing API key
		return &Service{
			client: nil,
			model:  nil,
			logger: logger,
			apiKey: apiKey,
		}
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		logger.Error("Failed to create gemini client", zap.Error(err))
		return &Service{
			client: nil,
			model:  nil,
			logger: logger,
			apiKey: apiKey,
		}
	}

	modelName := "gemini-1.5-flash" // Default model
	model := client.GenerativeModel(modelName)
	
	// Configure model for better translation performance
	model.SetTemperature(0.1) // Low temperature for consistent translations
	model.SetTopK(1)
	model.SetTopP(0.1)

	return &Service{
		client: client,
		model:  model,
		logger: logger,
		apiKey: apiKey,
	}
}

// Close closes the Gemini client
func (s *Service) Close() error {
	return s.client.Close()
}

// TranslateText translates a single text string
func (s *Service) TranslateText(ctx context.Context, req *TranslationRequest) (*TranslationResponse, error) {
	if req.Text == "" {
		return nil, fmt.Errorf("text is required for translation")
	}

	if req.TargetLang == "" {
		return nil, fmt.Errorf("target language is required")
	}

	// Build the translation prompt
	prompt := s.buildTranslationPrompt(req)

	// Generate translation
	resp, err := s.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		s.logger.Error("Failed to generate translation", 
			zap.String("text", req.Text),
			zap.String("targetLang", req.TargetLang),
			zap.Error(err))
		return nil, fmt.Errorf("failed to generate translation: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no translation generated")
	}

	// Extract translated text
	translatedText := fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])
	translatedText = strings.TrimSpace(translatedText)

	return &TranslationResponse{
		OriginalText:   req.Text,
		TranslatedText: translatedText,
		SourceLang:     req.SourceLang,
		TargetLang:     req.TargetLang,
	}, nil
}

// TranslateSegments translates multiple transcript segments efficiently
func (s *Service) TranslateSegments(ctx context.Context, segments []types.TranscriptSegment, targetLang string, sourceLang string) ([]types.TranslatedSegment, error) {
	if len(segments) == 0 {
		return nil, fmt.Errorf("no segments to translate")
	}

	// Process segments in batches for efficiency
	batchSize := 10 // Adjust based on API limits and performance
	var allTranslations []types.TranslatedSegment

	for i := 0; i < len(segments); i += batchSize {
		end := i + batchSize
		if end > len(segments) {
			end = len(segments)
		}

		batch := segments[i:end]
		translations, err := s.translateBatch(ctx, batch, targetLang, sourceLang)
		if err != nil {
			s.logger.Error("Failed to translate batch", 
				zap.Int("batchStart", i),
				zap.Int("batchEnd", end),
				zap.Error(err))
			return nil, fmt.Errorf("failed to translate batch: %w", err)
		}

		allTranslations = append(allTranslations, translations...)

		// Add small delay between batches to respect rate limits
		if end < len(segments) {
			time.Sleep(100 * time.Millisecond)
		}
	}

	return allTranslations, nil
}

// DetectLanguage detects the language of the given text
func (s *Service) DetectLanguage(ctx context.Context, text string) (string, error) {
	if text == "" {
		return "", fmt.Errorf("text is required for language detection")
	}

	prompt := fmt.Sprintf(`Detect the language of the following text and respond with only the ISO 639-1 language code (e.g., "en", "es", "fr", "de", "ja", "zh", etc.):

Text: "%s"

Response format: Just the 2-letter language code`, text)

	resp, err := s.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		s.logger.Error("Failed to detect language", zap.String("text", text), zap.Error(err))
		return "", fmt.Errorf("failed to detect language: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no language detection result")
	}

	language := fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])
	language = strings.TrimSpace(strings.ToLower(language))

	// Validate that we got a reasonable language code
	if len(language) != 2 {
		return "", fmt.Errorf("invalid language code detected: %s", language)
	}

	return language, nil
}

// translateBatch translates a batch of segments together for efficiency
func (s *Service) translateBatch(ctx context.Context, segments []types.TranscriptSegment, targetLang string, sourceLang string) ([]types.TranslatedSegment, error) {
	// Build a combined prompt with all segments
	var segmentTexts []string
	for i, segment := range segments {
		segmentTexts = append(segmentTexts, fmt.Sprintf("%d: %s", i, segment.Text))
	}

	combinedText := strings.Join(segmentTexts, "\n")
	
	req := &TranslationRequest{
		Text:       combinedText,
		SourceLang: sourceLang,
		TargetLang: targetLang,
		Context:    "This is a video transcript with numbered segments. Maintain the same numbering in your translation.",
	}

	response, err := s.TranslateText(ctx, req)
	if err != nil {
		return nil, err
	}

	// Parse the response to extract individual translations
	translatedLines := strings.Split(response.TranslatedText, "\n")
	var translations []types.TranslatedSegment

	for i, segment := range segments {
		var translatedText string
		
		// Try to find the corresponding translated line
		for _, line := range translatedLines {
			if strings.HasPrefix(line, fmt.Sprintf("%d:", i)) {
				// Remove the number prefix
				translatedText = strings.TrimSpace(strings.TrimPrefix(line, fmt.Sprintf("%d:", i)))
				break
			}
		}

		// Fallback: if we can't match by number, use positional matching
		if translatedText == "" && i < len(translatedLines) {
			translatedText = strings.TrimSpace(translatedLines[i])
		}

		// If still empty, use original text as fallback
		if translatedText == "" {
			translatedText = segment.Text
		}

		translations = append(translations, types.TranslatedSegment{
			Index:          segment.Index,
			OriginalText:   segment.Text,
			TranslatedText: translatedText,
		})
	}

	return translations, nil
}

// buildTranslationPrompt creates an optimized prompt for translation
func (s *Service) buildTranslationPrompt(req *TranslationRequest) string {
	var prompt strings.Builder

	if req.SourceLang != "" {
		prompt.WriteString(fmt.Sprintf("Translate the following text from %s to %s", req.SourceLang, req.TargetLang))
	} else {
		prompt.WriteString(fmt.Sprintf("Translate the following text to %s", req.TargetLang))
	}

	if req.Context != "" {
		prompt.WriteString(fmt.Sprintf(" (%s)", req.Context))
	}

	prompt.WriteString(". Provide only the translation without any additional text, explanations, or formatting:\n\n")
	prompt.WriteString(req.Text)

	return prompt.String()
}

// GetSupportedLanguages returns a list of commonly supported languages
func (s *Service) GetSupportedLanguages() []types.Language {
	return []types.Language{
		{Code: "en", Name: "English"},
		{Code: "es", Name: "Spanish"},
		{Code: "fr", Name: "French"},
		{Code: "de", Name: "German"},
		{Code: "it", Name: "Italian"},
		{Code: "pt", Name: "Portuguese"},
		{Code: "ru", Name: "Russian"},
		{Code: "ja", Name: "Japanese"},
		{Code: "ko", Name: "Korean"},
		{Code: "zh", Name: "Chinese"},
		{Code: "ar", Name: "Arabic"},
		{Code: "hi", Name: "Hindi"},
		{Code: "th", Name: "Thai"},
		{Code: "vi", Name: "Vietnamese"},
		{Code: "nl", Name: "Dutch"},
		{Code: "sv", Name: "Swedish"},
		{Code: "no", Name: "Norwegian"},
		{Code: "da", Name: "Danish"},
		{Code: "fi", Name: "Finnish"},
		{Code: "pl", Name: "Polish"},
	}
}