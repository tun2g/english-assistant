// Types for transcript functionality
export interface TranscriptSegment {
  text: string;
  start: number;
  end: number;
  duration: number;
  index?: number;
}

export interface TranslatedSegment extends TranscriptSegment {
  originalText: string;
  translatedText: string;
  isTranslated: boolean;
}

export interface VideoTranscript {
  videoId: string;
  provider: string;
  language: string;
  segments: TranscriptSegment[];
  available: boolean;
  source: string;
}

export interface DualLanguageTranscript {
  videoId: string;
  primaryLanguage: string;
  secondaryLanguage: string;
  originalSegments: TranscriptSegment[];
  translatedSegments: TranslatedSegment[];
  currentSegmentIndex: number;
}

export interface TranscriptCache {
  [videoId: string]: {
    transcript: VideoTranscript;
    translations: Record<string, TranslatedSegment[]>;
    timestamp: number;
  };
}

export interface TranslationQueueItem {
  segmentIndex: number;
  segment: TranscriptSegment;
  targetLanguage: string;
  priority: number;
}
