package docs

import (
	"app-backend/internal/dto"
)

// NewTranslationDocs creates instances of translation-related DTOs for swagger documentation
// This function is never called but ensures the DTOs are considered "used" by the linter
func NewTranslationDocs() {
	_ = dto.TranslateTextsRequest{}
	_ = dto.TranslateTextsResponse{}
}

// TranslateTexts godoc
// @Summary Translate texts
// @Description Translate array of texts to target language using AI. Source language is auto-detected if not provided.
// @Tags translation
// @Accept json
// @Produce json
// @Param request body dto.TranslateTextsRequest true "Translation request"
// @Success 200 {object} dto.TranslateTextsResponse "Translated texts"
// @Failure 400 {object} dto.ErrorResponse "Invalid request"
// @Failure 500 {object} dto.ErrorResponse "Translation service error"
// @Router /api/v1/translate [post]
// @Security BearerAuth
func TranslateTexts() {}