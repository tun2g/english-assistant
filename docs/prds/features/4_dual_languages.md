# Dual Language Transcripts Feature PRD

## Overview

The Dual Language Transcripts feature provides synchronized bilingual video transcripts by combining
original transcript extraction with AI-powered translation services. This system enables users to
consume video content with side-by-side original and translated text segments, maintaining precise
timing synchronization for enhanced language learning experiences.

## Product Scope

### In Scope

- Original transcript extraction from multiple video providers
- AI-powered translation of transcript segments
- Synchronized dual-language data structure
- Language detection and validation
- Batch processing for translation efficiency
- Error handling and fallback mechanisms
- Support for 20+ languages

### Out of Scope (Phase 1)

- Real-time translation streaming
- Translation caching/persistence
- User-customizable translation models
- Manual translation corrections
- Translation quality scoring

## Technical Architecture

### Core Components

#### Video Service Layer (`services/video/`)

**Main Service** (`service.go`):

- Orchestrates transcript retrieval and translation workflow
- Implements `GetDualLanguageTranscript()` method
- Handles provider detection and video ID extraction
- Manages language detection fallback logic

**Service Interface** (`interface.go`):

```go
GetDualLanguageTranscript(ctx context.Context, provider VideoProvider,
    videoID string, sourceLang string, targetLang string) (*DualLanguageTranscript, error)
```

#### Translation Service (`pkg/gemini/`)

**Google Gemini Integration** (`service.go`):

- Implements Google Gemini 1.5 Flash model for translations
- Batch processing with configurable size (default: 10 segments)
- Rate limiting with 100ms delays between batches
- Language detection capabilities
- Optimized prompts for video transcript context

#### Data Models (`types/video.go`)

**Core Structures**:

```go
type DualLanguageTranscript struct {
    VideoID      string               `json:"videoId"`
    Provider     VideoProvider        `json:"provider"`
    SourceLang   string               `json:"sourceLang"`
    TargetLang   string               `json:"targetLang"`
    Segments     []TranscriptSegment  `json:"segments"`
    Translations []TranslatedSegment  `json:"translations"`
    Cached       bool                 `json:"cached"`
}

type TranslatedSegment struct {
    Index          int    `json:"index"`
    OriginalText   string `json:"originalText"`
    TranslatedText string `json:"translatedText"`
}
```

### Data Flow Architecture

#### Primary Workflow

1. **Provider Detection**: Extract video provider and ID from URL
2. **Transcript Retrieval**: Fetch original transcript using multi-provider system
3. **Language Detection**: Auto-detect source language if not provided
4. **Batch Translation**: Process segments in batches of 10 via Gemini API
5. **Data Synchronization**: Align translated segments with original timing
6. **Response Construction**: Build unified dual-language response

#### Detailed Implementation (`GetDualLanguageTranscript`)

```go
// apps/backend/internal/services/video/service.go:121-172
func (s *Service) GetDualLanguageTranscript(ctx context.Context,
    provider types.VideoProvider, videoID string,
    sourceLang string, targetLang string) (*types.DualLanguageTranscript, error) {

    // 1. Validation
    if s.translator == nil {
        return nil, fmt.Errorf("translation service not available")
    }

    // 2. Get original transcript
    transcript, err := s.GetTranscript(ctx, provider, videoID, sourceLang)
    if err != nil {
        return nil, fmt.Errorf("failed to get transcript: %w", err)
    }

    // 3. Handle empty transcripts
    if !transcript.Available || len(transcript.Segments) == 0 {
        return &types.DualLanguageTranscript{
            VideoID:  videoID,
            Provider: provider,
            Segments: []types.TranscriptSegment{},
        }, nil
    }

    // 4. Language detection fallback
    detectedSourceLang := transcript.Language
    if sourceLang == "" && len(transcript.Segments) > 0 {
        sampleText := ""
        for i, segment := range transcript.Segments {
            if i >= 3 { break } // Use first 3 segments
            sampleText += segment.Text + " "
        }

        if detectedLang, err := s.translator.DetectLanguage(ctx, sampleText); err == nil {
            detectedSourceLang = detectedLang
        }
    }

    // 5. Batch translation
    translations, err := s.translator.TranslateSegments(ctx,
        transcript.Segments, targetLang, detectedSourceLang)
    if err != nil {
        return nil, fmt.Errorf("failed to translate segments: %w", err)
    }

    // 6. Construct response
    return &types.DualLanguageTranscript{
        VideoID:      videoID,
        Provider:     provider,
        SourceLang:   detectedSourceLang,
        TargetLang:   targetLang,
        Segments:     transcript.Segments,
        Translations: translations,
        Cached:       false, // TODO: implement caching
    }, nil
}
```

### Translation Service Implementation

#### Batch Processing Strategy (`pkg/gemini/service.go`)

**Segment Translation** (`TranslateSegments` method):

- Processes segments in configurable batches (default: 10)
- Implements rate limiting between batches (100ms delay)
- Maintains segment indexing throughout process
- Provides fallback text on translation failures

**Batch Translation Logic** (`translateBatch` method):

```go
// Build numbered segment format
var segmentTexts []string
for i, segment := range segments {
    segmentTexts = append(segmentTexts, fmt.Sprintf("%d: %s", i, segment.Text))
}
combinedText := strings.Join(segmentTexts, "\n")

// Translation with context
req := &TranslationRequest{
    Text:       combinedText,
    SourceLang: sourceLang,
    TargetLang: targetLang,
    Context:    "This is a video transcript with numbered segments. Maintain the same numbering in your translation.",
}
```

**Response Parsing**:

- Extracts numbered translations from batch response
- Falls back to positional matching if numbering fails
- Uses original text as final fallback for missing translations

#### Model Configuration

**Gemini 1.5 Flash Settings**:

- Temperature: 0.1 (consistent translations)
- TopK: 1 (focused responses)
- TopP: 0.1 (deterministic output)
- Model: `gemini-1.5-flash` (optimized for speed and accuracy)

### Error Handling & Resilience

#### Translation Service Availability

- Graceful degradation when Gemini API unavailable
- Service initialization continues without translation capabilities
- Clear error messages for missing API key scenarios

#### Segment Processing Failures

- Individual segment fallback to original text
- Batch retry logic for temporary API failures
- Comprehensive logging for debugging translation issues

#### Language Detection Fallback

- Uses first 3 segments for language detection sampling
- Falls back to transcript metadata language
- Continues processing even with detection failures

## API Specifications

### Endpoints

**Current Implementation**: No dedicated dual-language endpoint exists yet in the video handler.

**Proposed Endpoint**:

```http
GET /api/v1/videos/{videoUrl}/dual-transcript?sourceLang={lang}&targetLang={lang}
```

**Response Format**:

```json
{
  "videoId": "dQw4w9WgXcQ",
  "provider": "youtube",
  "sourceLang": "en",
  "targetLang": "es",
  "segments": [
    {
      "startTime": 0,
      "endTime": 3000,
      "text": "Hello world",
      "index": 1
    }
  ],
  "translations": [
    {
      "index": 1,
      "originalText": "Hello world",
      "translatedText": "Hola mundo"
    }
  ],
  "cached": false
}
```

### Integration Points

**Video Service Integration**: Leverages existing multi-provider transcript system for consistent
data sourcing.

**Translation Service Integration**: Direct integration with Google Gemini service for AI-powered
translations.

**Handler Layer**: Currently uses video handler for transcript-related endpoints, dual-language
endpoint would follow same pattern.

## Performance Optimizations

### Batch Processing

- **Batch Size**: 10 segments per API call (configurable)
- **Rate Limiting**: 100ms delays between batches
- **Parallel Processing**: Potential for concurrent batch processing

### Language Detection Optimization

- **Sampling Strategy**: Uses first 3 segments only for detection
- **Detection Caching**: Language detected once per transcript
- **Fallback Chain**: Metadata � Detection � Default

### Memory Management

- **Stream Processing**: Processes segments in batches to control memory usage
- **Garbage Collection**: Minimal object allocation during translation loops
- **Response Building**: Efficient slice operations for segment alignment

## Future Enhancements

### Caching Implementation

**Database Schema**:

```sql
CREATE TABLE dual_language_transcripts (
    id SERIAL PRIMARY KEY,
    video_id VARCHAR(255) NOT NULL,
    provider VARCHAR(50) NOT NULL,
    source_lang VARCHAR(10) NOT NULL,
    target_lang VARCHAR(10) NOT NULL,
    transcript_data JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP,
    INDEX idx_video_lang (video_id, source_lang, target_lang)
);
```

### Performance Monitoring

- **Translation Latency**: Track batch processing times
- **API Usage**: Monitor Gemini API quota consumption
- **Error Rates**: Track translation failure patterns
- **Cache Hit Rates**: Measure caching effectiveness

### Advanced Features

- **Translation Quality Scoring**: Confidence metrics for translations
- **User Corrections**: Manual translation improvement system
- **Multiple Translation Providers**: Support for alternative translation APIs
- **Real-time Translation**: Streaming translation for live content

## Implementation Roadmap

### Phase 1: Core Functionality (Current)

-  Dual language data structures
-  Translation service integration
-  Batch processing implementation
-  Error handling and fallbacks
- L HTTP endpoint implementation
- L Frontend integration

### Phase 2: Performance & Caching

- Database caching implementation
- Response time optimization
- Memory usage optimization
- API rate limiting improvements
- Performance monitoring setup

## Dependencies

### External Services

- **Google Gemini 1.5 Flash**: Primary translation provider
- **YouTube Transcript Providers**: Source transcript data
- **Video Provider APIs**: Video metadata and capabilities

### Internal Dependencies

- **Video Service**: Transcript extraction and provider management
- **Types Package**: Shared data structures
- **Configuration System**: API key management
- **Logging System**: Error tracking and debugging

### Required Configurations

```yaml
external:
  gemini_api_key: 'your-gemini-api-key'
  youtube_api_key: 'your-youtube-api-key'

translation:
  batch_size: 10
  rate_limit_delay: 100ms
  model_name: 'gemini-1.5-flash'
  temperature: 0.1
```

## Success Metrics

### Functional Metrics

- **Translation Accuracy**: >90% contextually appropriate translations
- **Language Coverage**: Support for 20+ language pairs
- **Error Rate**: <5% translation failures

### Performance Metrics

- **Response Time**: <10 seconds for 5-minute video transcripts
- **API Efficiency**: <50 Gemini API calls per video
- **Memory Usage**: <100MB per concurrent translation

### User Experience Metrics

- **Feature Adoption**: % of users using dual-language feature
- **Session Duration**: Increased engagement with translated content
- **Language Learning**: User progression tracking with bilingual content
