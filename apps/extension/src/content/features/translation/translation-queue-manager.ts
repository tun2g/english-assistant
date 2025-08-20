import { translateTexts } from '@english/shared/api/translation-api';
import type { TranscriptSegment, TranslatedSegment } from '../../../shared/types/extension-types';
import { captureError, logDebug } from '../../utils/error-handler';

export interface TranslationQueueCallbacks {
  onTranslationComplete: (segment: TranslatedSegment, segmentIndex: number) => void;
  onTranslationError: (segmentIndex: number, error: Error) => void;
  onBatchTranslationComplete: (translatedSegments: TranslatedSegment[]) => void;
}

export class TranslationQueueManager {
  private callbacks: TranslationQueueCallbacks;
  private completed: Map<string, TranslatedSegment> = new Map();
  private isActive: boolean = false;
  private videoId: string | null = null;
  private targetLanguage: string = 'es';
  private processing: boolean = false;

  constructor(callbacks: TranslationQueueCallbacks) {
    this.callbacks = callbacks;
  }

  // Set current video and target language
  setContext(videoId: string, targetLanguage: string): void {
    // Clear previous context if video changed
    if (this.videoId !== videoId) {
      this.completed.clear();
    }

    this.videoId = videoId;
    this.targetLanguage = targetLanguage;
  }

  // Start the translation service
  start(): void {
    if (this.isActive) return;

    this.isActive = true;
    console.log('TranslationQueueManager: Started');
  }

  // Stop the translation service
  stop(): void {
    if (!this.isActive) return;

    this.isActive = false;
    console.log('TranslationQueueManager: Stopped');
  }

  // Translate all segments in a transcript at once
  async translateAllSegments(segments: TranscriptSegment[]): Promise<TranslatedSegment[]> {
    if (!this.isActive || this.processing) {
      logDebug(
        'Cannot translate segments - Queue not active or already processing',
        undefined,
        'TranslationQueueManager'
      );
      return [];
    }

    this.processing = true;

    try {
      // Filter out already translated segments
      const segmentsToTranslate = segments.filter(
        segment => !this.isTranslated(segment.index ?? segments.indexOf(segment))
      );

      if (segmentsToTranslate.length === 0) {
        logDebug('All segments already translated', undefined, 'TranslationQueueManager');
        // Return existing translations
        return segments
          .map((segment, index) => this.getTranslatedSegment(segment.index ?? index))
          .filter(Boolean) as TranslatedSegment[];
      }

      logDebug(
        'Translating batch of segments',
        {
          total: segments.length,
          toTranslate: segmentsToTranslate.length,
          targetLang: this.targetLanguage,
        },
        'TranslationQueueManager'
      );

      // Extract texts to translate
      const textsToTranslate = segmentsToTranslate.map(segment => segment.text);

      // Batch translate all texts at once
      const translationResponse = await translateTexts({
        texts: textsToTranslate,
        targetLang: this.targetLanguage,
        sourceLang: 'auto',
      });

      // Handle both wrapped (ApiResponse) and direct response formats
      const translationData = translationResponse;

      if (!translationData || !translationData.translations) {
        throw new Error('Invalid translation response format');
      }

      if (translationData.translations.length !== textsToTranslate.length) {
        throw new Error(
          `Translation mismatch: expected ${textsToTranslate.length} translations, got ${translationData.translations.length}`
        );
      }

      // Create translated segments
      const translatedSegments: TranslatedSegment[] = [];

      segmentsToTranslate.forEach((segment, index) => {
        const segmentIndex = segment.index ?? segments.indexOf(segment);
        const translatedText = translationData.translations[index];

        const translatedSegment: TranslatedSegment = {
          ...segment,
          index: segmentIndex,
          originalText: segment.text,
          translatedText: translatedText,
          isTranslated: true,
        };

        // Cache the result
        const key = `${segmentIndex}_${this.targetLanguage}`;
        this.completed.set(key, translatedSegment);

        translatedSegments.push(translatedSegment);

        // Notify individual completion callback
        this.callbacks.onTranslationComplete(translatedSegment, segmentIndex);
      });

      logDebug(
        'Batch translation completed successfully',
        {
          translatedCount: translatedSegments.length,
        },
        'TranslationQueueManager'
      );

      // Notify batch completion callback
      this.callbacks.onBatchTranslationComplete(translatedSegments);

      // Return all translated segments (including previously cached ones)
      return segments
        .map((segment, index) => {
          const segmentIndex = segment.index ?? index;
          return this.getTranslatedSegment(segmentIndex);
        })
        .filter(Boolean) as TranslatedSegment[];
    } catch (error) {
      captureError('Batch translation failed', error, 'translateAllSegments', 'TranslationQueueManager');

      // Notify error for all segments
      segments.forEach((segment, index) => {
        const segmentIndex = segment.index ?? index;
        if (!this.isTranslated(segmentIndex)) {
          this.callbacks.onTranslationError(segmentIndex, error as Error);
        }
      });

      return [];
    } finally {
      this.processing = false;
    }
  }

  // Translate single segment immediately (for on-demand translation)
  async translateSegmentImmediate(segment: TranscriptSegment, segmentIndex: number): Promise<TranslatedSegment | null> {
    if (!this.videoId) return null;

    // Check if already translated
    const existing = this.getTranslatedSegment(segmentIndex);
    if (existing) return existing;

    try {
      console.log('TranslationQueueManager: Immediate translation for segment', segmentIndex);

      const translationResponse = await translateTexts({
        texts: [segment.text],
        targetLang: this.targetLanguage,
        sourceLang: 'auto',
      });

      // Handle both wrapped (ApiResponse) and direct response formats
      const translationData = translationResponse;

      if (translationData && translationData.translations && translationData.translations.length > 0) {
        const translatedText = translationData.translations[0];

        const translatedSegment: TranslatedSegment = {
          ...segment,
          index: segmentIndex,
          originalText: segment.text,
          translatedText: translatedText,
          isTranslated: true,
        };

        // Cache the result
        const key = `${segmentIndex}_${this.targetLanguage}`;
        this.completed.set(key, translatedSegment);

        return translatedSegment;
      }
    } catch (error) {
      console.error('TranslationQueueManager: Immediate translation failed for segment', segmentIndex, error);
    }

    return null;
  }

  // Get translated segment if available
  getTranslatedSegment(segmentIndex: number): TranslatedSegment | null {
    const key = `${segmentIndex}_${this.targetLanguage}`;
    return this.completed.get(key) || null;
  }

  // Check if segment is already translated
  isTranslated(segmentIndex: number): boolean {
    const key = `${segmentIndex}_${this.targetLanguage}`;
    return this.completed.has(key);
  }

  // Get statistics
  getStats(): {
    completed: number;
    processing: boolean;
  } {
    return {
      completed: this.completed.size,
      processing: this.processing,
    };
  }

  // Clear all translations (e.g., when language changes)
  clearTranslations(): void {
    this.completed.clear();
  }
}
