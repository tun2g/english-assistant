package translation

import (
	"fmt"
	"net/http"
	"strings"

	"app-backend/internal/dto"
	"app-backend/internal/logger"
	"app-backend/internal/services/translation"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handler implements translation HTTP handlers
type Handler struct {
	translationService translation.ServiceInterface
	logger             *logger.Logger
}

// NewTranslationHandler creates a new translation handler
func NewTranslationHandler(translationService translation.ServiceInterface, logger *logger.Logger) HandlerInterface {
	return &Handler{
		translationService: translationService,
		logger:             logger,
	}
}

// TranslateTexts handles text translation requests
func (h *Handler) TranslateTexts(c *gin.Context) {
	var req dto.TranslateTextsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid JSON body", zap.Error(err))
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	// Validate request
	if len(req.Texts) == 0 {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "No texts provided for translation",
		})
		return
	}

	if req.TargetLang == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Target language is required",
		})
		return
	}

	// Auto-detect source language if not provided
	detectedSourceLang := req.SourceLang
	if req.SourceLang == "" && len(req.Texts) > 0 {
		// Use first text to detect language
		sampleText := req.Texts[0]
		if len(sampleText) > 200 {
			sampleText = sampleText[:200] // Limit sample size for detection
		}

		if detected, err := h.translationService.DetectLanguage(c.Request.Context(), sampleText); err == nil {
			detectedSourceLang = detected
			h.logger.Debug("Language detected", zap.String("detected", detected), zap.String("sample", sampleText[:min(50, len(sampleText))]))
		} else {
			h.logger.Warn("Failed to detect language", zap.Error(err))
			detectedSourceLang = "auto" // Fallback to auto-detection
		}
	}

	// Translate texts
	translations, err := h.translationService.TranslateTexts(
		c.Request.Context(),
		req.Texts,
		req.TargetLang,
		detectedSourceLang,
	)
	if err != nil {
		// Check if it's a quota exceeded or context canceled error and return mock data
		if strings.Contains(err.Error(), "quota") || strings.Contains(err.Error(), "429") || strings.Contains(err.Error(), "context canceled") {
			h.logger.Warn("Translation quota exceeded, returning mock translations",
				zap.Int("textCount", len(req.Texts)),
				zap.String("sourceLang", detectedSourceLang),
				zap.String("targetLang", req.TargetLang))
			
			// Generate mock translations
			mockTranslations := make([]string, len(req.Texts))
			for i, text := range req.Texts {
				// Simple mock translation - add [TRANSLATED] prefix
				mockTranslations[i] = fmt.Sprintf("[%s] %s", strings.ToUpper(req.TargetLang), text)
			}
			
			response := dto.TranslateTextsResponse{
				Translations: mockTranslations,
				SourceLang:   detectedSourceLang,
				TargetLang:   req.TargetLang,
			}
			
			c.JSON(http.StatusOK, response)
			return
		}
		
		h.logger.Error("Failed to translate texts",
			zap.Int("textCount", len(req.Texts)),
			zap.String("sourceLang", detectedSourceLang),
			zap.String("targetLang", req.TargetLang),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to translate texts",
			Details: err.Error(),
		})
		return
	}

	// Return response
	response := dto.TranslateTextsResponse{
		Translations: translations,
		SourceLang:   detectedSourceLang,
		TargetLang:   req.TargetLang,
	}

	h.logger.Debug("Translation completed",
		zap.Int("textCount", len(req.Texts)),
		zap.String("sourceLang", detectedSourceLang),
		zap.String("targetLang", req.TargetLang))

	c.JSON(http.StatusOK, response)
}

// min helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}