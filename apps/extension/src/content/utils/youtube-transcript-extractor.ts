// YouTube transcript extractor - works directly with YouTube's native transcript API
import type { VideoTranscript, TranscriptSegment } from '../../shared/types/extension-types';
import { captureError, logDebug } from './error-handler';

export class YouTubeTranscriptExtractor {
  /**
   * Extract transcript directly from YouTube's native transcript API
   */
  static async extractTranscript(videoId: string): Promise<VideoTranscript | null> {
    try {
      logDebug('Attempting to extract transcript for video:', { videoId }, 'YouTubeTranscriptExtractor');

      // First try to get transcript from YouTube's player API
      const transcript = await this.getTranscriptFromYouTubeAPI(videoId);
      if (transcript) {
        return transcript;
      }

      // Fallback: try to extract from page DOM
      const domTranscript = await this.getTranscriptFromDOM(videoId);
      if (domTranscript) {
        return domTranscript;
      }

      logDebug('No transcript found for video', { videoId }, 'YouTubeTranscriptExtractor');
      return null;
    } catch (error) {
      captureError('Failed to extract transcript', error, 'extract', 'YouTubeTranscriptExtractor');
      return null;
    }
  }

  /**
   * Get transcript from YouTube's internal API
   */
  private static async getTranscriptFromYouTubeAPI(videoId: string): Promise<VideoTranscript | null> {
    try {
      // Look for YouTube's player data in the page
      const playerData = this.extractPlayerData();
      if (!playerData) {
        logDebug('No player data found', undefined, 'YouTubeTranscriptExtractor');
        return null;
      }

      // Try to find captions/transcript data
      const captionsData = this.extractCaptionsData(playerData);
      if (!captionsData) {
        logDebug('No captions data found in player', undefined, 'YouTubeTranscriptExtractor');
        return null;
      }

      // Parse the captions into transcript segments
      const segments = this.parseCaptionsToSegments(captionsData);
      if (segments.length === 0) {
        logDebug('No valid segments found in captions', undefined, 'YouTubeTranscriptExtractor');
        return null;
      }

      const transcript: VideoTranscript = {
        videoId,
        provider: 'youtube',
        language: captionsData.languageCode || 'unknown',
        segments,
        available: true,
        source: 'youtube-native',
      };

      logDebug(
        'Successfully extracted transcript from YouTube API',
        {
          videoId,
          segmentCount: segments.length,
          language: transcript.language,
        },
        'YouTubeTranscriptExtractor'
      );

      return transcript;
    } catch (error) {
      captureError('Failed to get transcript from YouTube API', error, 'youtube-api', 'YouTubeTranscriptExtractor');
      return null;
    }
  }

  /**
   * Get transcript from DOM elements (fallback)
   */
  private static async getTranscriptFromDOM(videoId: string): Promise<VideoTranscript | null> {
    try {
      logDebug('Attempting DOM transcript extraction', { videoId }, 'YouTubeTranscriptExtractor');

      // Look for transcript button and click it to open transcript panel
      const transcriptButton = document.querySelector(
        '[aria-label*="transcript" i], [aria-label*="Show transcript" i]'
      ) as HTMLElement;
      if (transcriptButton) {
        transcriptButton.click();
        await this.delay(1000); // Wait for transcript panel to load
      }

      // Look for transcript items in the DOM
      const transcriptItems = document.querySelectorAll(
        '[data-params*="transcript"], .ytd-transcript-segment-renderer'
      );

      if (transcriptItems.length === 0) {
        logDebug('No transcript items found in DOM', undefined, 'YouTubeTranscriptExtractor');
        return null;
      }

      const segments: TranscriptSegment[] = [];

      transcriptItems.forEach((item, index) => {
        try {
          const textElement = item.querySelector('.segment-text, [class*="text"]');
          const timeElement = item.querySelector('.segment-timestamp, [class*="time"]');

          if (textElement && timeElement) {
            const text = textElement.textContent?.trim() || '';
            const timeStr = timeElement.textContent?.trim() || '';
            const startTime = this.parseTimeString(timeStr);

            if (text && startTime !== null) {
              segments.push({
                text,
                start: startTime,
                end: startTime + 3, // Default 3 second duration
                duration: 3,
                index,
              });
            }
          }
        } catch (error) {
          // Skip invalid segments
        }
      });

      if (segments.length === 0) {
        logDebug('No valid segments extracted from DOM', undefined, 'YouTubeTranscriptExtractor');
        return null;
      }

      const transcript: VideoTranscript = {
        videoId,
        provider: 'youtube',
        language: 'unknown',
        segments,
        available: true,
        source: 'youtube-dom',
      };

      logDebug(
        'Successfully extracted transcript from DOM',
        {
          videoId,
          segmentCount: segments.length,
        },
        'YouTubeTranscriptExtractor'
      );

      return transcript;
    } catch (error) {
      captureError('Failed to get transcript from DOM', error, 'dom-extract', 'YouTubeTranscriptExtractor');
      return null;
    }
  }

  /**
   * Extract YouTube player data from page
   */
  private static extractPlayerData(): any {
    try {
      // Look for ytInitialPlayerResponse in page
      const scripts = document.querySelectorAll('script');
      for (const script of scripts) {
        const content = script.textContent || '';
        if (content.includes('ytInitialPlayerResponse')) {
          const match = content.match(/ytInitialPlayerResponse\s*=\s*({.+?});/);
          if (match) {
            return JSON.parse(match[1]);
          }
        }
      }

      // Try to get from window object
      // @ts-ignore
      if (window.ytInitialPlayerResponse) {
        // @ts-ignore
        return window.ytInitialPlayerResponse;
      }

      return null;
    } catch (error) {
      captureError('Failed to extract player data', error, 'player-data', 'YouTubeTranscriptExtractor');
      return null;
    }
  }

  /**
   * Extract captions data from player data
   */
  private static extractCaptionsData(playerData: any): any {
    try {
      const captions = playerData?.captions?.playerCaptionsTracklistRenderer?.captionTracks;
      if (!captions || !Array.isArray(captions)) {
        return null;
      }

      // Prefer auto-generated captions or first available
      const autoCaption = captions.find((c: any) => c.kind === 'asr');
      const firstCaption = captions[0];

      return autoCaption || firstCaption;
    } catch (error) {
      return null;
    }
  }

  /**
   * Parse captions data to transcript segments
   */
  private static parseCaptionsToSegments(captionsData: any): TranscriptSegment[] {
    try {
      if (!captionsData.baseUrl) {
        return [];
      }

      // Note: In a real implementation, we would fetch the baseUrl to get the actual captions XML
      // For now, return empty array as we need to handle CORS issues
      logDebug(
        'Caption URL found but CORS prevents direct fetch',
        {
          baseUrl: captionsData.baseUrl,
        },
        'YouTubeTranscriptExtractor'
      );

      return [];
    } catch (error) {
      return [];
    }
  }

  /**
   * Parse time string to seconds
   */
  private static parseTimeString(timeStr: string): number | null {
    try {
      // Handle formats like "1:23", "12:34", "1:23:45"
      const parts = timeStr.split(':').reverse();
      let seconds = 0;

      for (let i = 0; i < parts.length; i++) {
        const value = parseInt(parts[i], 10);
        if (isNaN(value)) return null;

        seconds += value * Math.pow(60, i);
      }

      return seconds;
    } catch (error) {
      return null;
    }
  }

  /**
   * Simple delay utility
   */
  private static delay(ms: number): Promise<void> {
    return new Promise(resolve => setTimeout(resolve, ms));
  }

  /**
   * Check if transcript is likely available for a video
   */
  static async isTranscriptAvailable(videoId: string): Promise<boolean> {
    try {
      // Check if transcript button exists
      const transcriptButton = document.querySelector(
        '[aria-label*="transcript" i], [aria-label*="Show transcript" i]'
      );
      if (transcriptButton) {
        return true;
      }

      // Check player data
      const playerData = this.extractPlayerData();
      if (playerData?.captions?.playerCaptionsTracklistRenderer?.captionTracks?.length > 0) {
        return true;
      }

      return false;
    } catch (error) {
      return false;
    }
  }
}

/**
 * Create a mock transcript for testing when no real transcript is available
 */
export function createMockTranscript(videoId: string): VideoTranscript {
  const mockSegments: TranscriptSegment[] = [
    { text: 'Welcome to this video', start: 0, end: 3, duration: 3, index: 0 },
    {
      text: "Today we're going to learn something interesting",
      start: 3,
      end: 8,
      duration: 5,
      index: 1,
    },
    {
      text: "Let's start with the basics",
      start: 8,
      end: 12,
      duration: 4,
      index: 2,
    },
    {
      text: 'This is a demonstration of dual language support',
      start: 12,
      end: 18,
      duration: 6,
      index: 3,
    },
    {
      text: 'You can see both original and translated text',
      start: 18,
      end: 24,
      duration: 6,
      index: 4,
    },
  ];

  return {
    videoId,
    provider: 'mock',
    language: 'en',
    segments: mockSegments,
    available: true,
    source: 'mock',
  };
}
