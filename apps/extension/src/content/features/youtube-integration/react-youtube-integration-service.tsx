import React from 'react';
import { getVideoTranscript } from '@english/shared/api/video-api';

import { OAuthManager } from '../auth/oauth-manager';
import { NotificationManager } from '../notifications/notification-manager';
import { PlayerControlsManager } from './player-controls-manager';
import { VideoMonitor } from './video-monitor';

import { EXTENSION_STORAGE_KEYS, EXTENSION_CLASSES } from '../../../shared/constants';
import { captureError, logDebug } from '../../utils/error-handler';
import { extractVideoIdFromCurrentUrl } from '../../utils/video/video-utils';
import type { VideoTranscript, TranscriptCache } from '../../../shared/types/extension-types';

import {
  renderReactComponent,
  unmountReactComponent,
  cleanupAllReactComponents,
} from '../../components/utils/react-renderer';
import { TranscriptOverlayManager } from '../../components/transcript/transcript-overlay-manager';

export class ReactYouTubeIntegrationService {
  private currentVideoId: string | null = null;
  private isIntegrationActive: boolean = false;
  private transcriptCache: TranscriptCache = {};
  private autoFetchEnabled: boolean = false;
  private overlayEnabled: boolean = false;

  // Feature managers
  private oauthManager: OAuthManager;
  private playerControlsManager: PlayerControlsManager;
  private videoMonitor: VideoMonitor;

  constructor() {
    this.oauthManager = new OAuthManager();

    this.playerControlsManager = new PlayerControlsManager({
      onButtonClick: () => this.handleExtensionButtonClick(),
    });

    this.videoMonitor = new VideoMonitor({
      onVideoChange: videoId => this.handleVideoChange(videoId),
    });
  }

  // Initialize YouTube integration
  async init(): Promise<void> {
    try {
      logDebug('Starting initialization...', undefined, 'ReactYouTubeIntegrationService');

      if (this.isIntegrationActive) {
        logDebug('Already active, skipping', undefined, 'ReactYouTubeIntegrationService');
        return;
      }

      this.isIntegrationActive = true;

      // Load user settings
      await this.loadSettings();

      // Start video monitoring
      this.videoMonitor.start();

      // Check current video
      this.checkCurrentVideo();

      logDebug('Initialization complete', undefined, 'ReactYouTubeIntegrationService');
    } catch (error) {
      captureError('Failed to initialize YouTube integration', error, 'init', 'ReactYouTubeIntegrationService');
      this.isIntegrationActive = false;
      throw error;
    }
  }

  // Clean up integration
  destroy(): void {
    console.log('ReactYouTubeIntegrationService: Destroying...');

    this.isIntegrationActive = false;

    // Clean up managers
    this.videoMonitor.stop();
    this.playerControlsManager.removeControls();
    this.oauthManager.destroy();

    // Clean up all React components
    cleanupAllReactComponents();

    console.log('ReactYouTubeIntegrationService: Destroyed');
  }

  // Handle video change
  private async handleVideoChange(videoId: string | null): Promise<void> {
    console.log(
      'ReactYouTubeIntegrationService: handleVideoChange called with:',
      videoId,
      'current:',
      this.currentVideoId
    );

    if (videoId !== this.currentVideoId) {
      this.currentVideoId = videoId;
      console.log('ReactYouTubeIntegrationService: Video ID changed to:', videoId);

      // Clean up previous video's React components
      unmountReactComponent(EXTENSION_CLASSES.TRANSCRIPT_OVERLAY);

      if (videoId) {
        console.log('ReactYouTubeIntegrationService: Initializing video integration...');
        await this.initializeVideoIntegration();
      } else {
        console.log('ReactYouTubeIntegrationService: No video ID, cleaning up...');
        this.playerControlsManager.removeControls();
      }
    } else {
      console.log('ReactYouTubeIntegrationService: Same video ID, no action needed');
    }
  }

  // Check current video
  private checkCurrentVideo(): void {
    const videoId = extractVideoIdFromCurrentUrl();
    this.handleVideoChange(videoId);
  }

  // Initialize integration for current video
  private async initializeVideoIntegration(): Promise<void> {
    if (!this.currentVideoId) return;

    this.playerControlsManager.injectControls();

    // Auto-fetch transcript if enabled and authenticated
    if (this.autoFetchEnabled && this.oauthManager.authenticated) {
      await this.autoFetchTranscript();
    }
  }

  // Handle extension button click
  private async handleExtensionButtonClick(): Promise<void> {
    try {
      logDebug('Extension button clicked!', undefined, 'ReactYouTubeIntegrationService');

      if (!this.currentVideoId) {
        NotificationManager.showError('No video detected!');
        captureError(
          'Extension button clicked but no video ID available',
          new Error('No video ID'),
          'button-click',
          'ReactYouTubeIntegrationService'
        );
        return;
      }

      // Check if language selector is currently showing
      const isLanguageSelectorActive = document.querySelector(`.${EXTENSION_CLASSES.TRANSCRIPT_OVERLAY}`) !== null;
      // Also check for any existing modal with language selector content
      const hasLanguageModal = document.querySelector('[data-english-extension="true"]') !== null;

      logDebug(
        'Component states',
        {
          isLanguageSelectorActive,
          hasLanguageModal,
        },
        'ReactYouTubeIntegrationService'
      );

      if (isLanguageSelectorActive || hasLanguageModal) {
        // If any extension modal/overlay is active, close it
        logDebug('Closing active overlay/modal', undefined, 'ReactYouTubeIntegrationService');
        unmountReactComponent(EXTENSION_CLASSES.TRANSCRIPT_OVERLAY);
        return;
      }

      // Always show language selector first (single click behavior)
      logDebug('Showing language selector popup', undefined, 'ReactYouTubeIntegrationService');
      await this.showLanguageSelectorPopup();
    } catch (error) {
      captureError('Failed to handle extension button click', error, 'button-click', 'ReactYouTubeIntegrationService');
      NotificationManager.showError('Something went wrong. Please try again.');
    }
  }

  // Show language selector popup for user to choose dual-language settings
  private async showLanguageSelectorPopup(): Promise<void> {
    // First ensure we have a transcript (fetch if needed)
    const cached = this.getCachedTranscript(this.currentVideoId!);
    let transcript = cached;

    if (!transcript) {
      logDebug('No cached transcript, attempting to fetch...', undefined, 'ReactYouTubeIntegrationService');
      NotificationManager.showLoading('Fetching transcript...');

      try {
        // Get transcript from backend API
        transcript = await this.fetchTranscriptFromAPI(this.currentVideoId!);

        if (!transcript) {
          NotificationManager.hide();
          NotificationManager.showError('No transcript available for this video.');
          return; // Exit early - no transcript means no dual language support
        }

        // Cache the transcript
        await this.cacheTranscript(transcript);
        NotificationManager.hide();
      } catch (error: any) {
        logDebug('All transcript methods failed', { error }, 'ReactYouTubeIntegrationService');
        NotificationManager.hide();
        NotificationManager.showError('Failed to get transcript for this video.');
        return; // Exit early - no transcript available
      }
    }

    // Initialize overlay with language selector visible
    await this.initializeReactOverlayWithLanguageSelector(transcript);
  }

  // Fetch transcript from backend API
  private async fetchTranscriptFromAPI(videoId: string): Promise<VideoTranscript | null> {
    logDebug('Fetching transcript from backend API', { videoId }, 'ReactYouTubeIntegrationService');

    try {
      const transcriptResponse = await getVideoTranscript({
        videoUrl: videoId,
      });

      // Handle both wrapped (ApiResponse) and direct response formats
      const transcript = transcriptResponse;
      if (transcript && transcript.segments && transcript.segments.length > 0) {
        const videoTranscript: VideoTranscript = {
          videoId: transcript.videoId || videoId,
          provider: transcript.provider || 'youtube',
          language: transcript.language || 'unknown',
          segments: transcript.segments.map((segment, index) => {
            // Debug: Log first few segment timings
            if (index < 3) {
              console.log('ReactYouTubeIntegrationService: Segment timing debug', {
                index,
                originalStartTime: segment.startTime,
                originalEndTime: segment.endTime,
                text: segment.text.substring(0, 30) + '...',
              });
            }

            // Temporary fix: Generate timing if missing
            // Use more realistic timing based on text length for faster sync
            const avgWordsPerSecond = 3; // Average speaking speed
            const wordCount = segment.text.split(' ').length;
            const estimatedDuration = Math.max(1.5, wordCount / avgWordsPerSecond); // Min 1.5 seconds

            let estimatedStart = 0;
            if (index > 0) {
              // Calculate cumulative start time from previous segments
              estimatedStart = transcript.segments.slice(0, index).reduce((acc, prevSeg) => {
                if (prevSeg.startTime && prevSeg.startTime > 0) return acc + (prevSeg.endTime - prevSeg.startTime);
                const prevWordCount = prevSeg.text.split(' ').length;
                return acc + Math.max(1.5, prevWordCount / avgWordsPerSecond);
              }, 0);
            }

            const startTime = segment.startTime || estimatedStart;
            const endTime = segment.endTime || startTime + estimatedDuration;

            return {
              text: segment.text,
              start: startTime,
              end: endTime,
              duration: endTime - startTime,
            };
          }),
          available: transcript.available,
          source: transcript.source || 'api',
        };

        return videoTranscript;
      } else {
        return null;
      }
    } catch (error: any) {
      console.error('ReactYouTubeIntegrationService: Backend API transcript fetch failed:', error);
      return null;
    }
  }

  // Initialize React overlay with language selector visible initially
  private async initializeReactOverlayWithLanguageSelector(transcript: VideoTranscript): Promise<void> {
    if (!this.currentVideoId) return;

    try {
      console.log('ReactYouTubeIntegrationService: Initializing React transcript overlay with language selector');

      // Create the React component with language selector initially visible
      const overlayComponent = (
        <TranscriptOverlayManager
          videoId={this.currentVideoId}
          transcript={transcript}
          showLanguageSelectorInitially={true}
          onReady={() => console.log('React overlay with language selector ready')}
          onError={(error: any) => console.error('React overlay error:', error)}
          onSegmentSeek={(seekTime: number) => this.handleSegmentSeek(seekTime)}
        />
      );

      // Render the React component
      await renderReactComponent(
        overlayComponent,
        EXTENSION_CLASSES.TRANSCRIPT_OVERLAY,
        'react-transcript-overlay-with-selector'
      );

      console.log('ReactYouTubeIntegrationService: React transcript overlay with language selector started');
    } catch (error) {
      console.error(
        'ReactYouTubeIntegrationService: Failed to initialize React overlay with language selector:',
        error
      );
    }
  }

  // Load user settings from storage
  private async loadSettings(): Promise<void> {
    try {
      const result = await new Promise<any>(resolve => {
        chrome.storage.local.get(
          [
            EXTENSION_STORAGE_KEYS.AUTO_FETCH_TRANSCRIPT,
            EXTENSION_STORAGE_KEYS.TRANSCRIPT_CACHE,
            EXTENSION_STORAGE_KEYS.TRANSCRIPT_OVERLAY_ENABLED,
          ],
          resolve
        );
      });

      this.autoFetchEnabled = result[EXTENSION_STORAGE_KEYS.AUTO_FETCH_TRANSCRIPT] || false;
      this.transcriptCache = result[EXTENSION_STORAGE_KEYS.TRANSCRIPT_CACHE] || {};
      this.overlayEnabled = result[EXTENSION_STORAGE_KEYS.TRANSCRIPT_OVERLAY_ENABLED] || false;

      console.log('ReactYouTubeIntegrationService: Settings loaded', {
        autoFetchEnabled: this.autoFetchEnabled,
        overlayEnabled: this.overlayEnabled,
        cacheSize: Object.keys(this.transcriptCache).length,
      });
    } catch (error) {
      console.error('Failed to load settings:', error);
    }
  }

  // Auto-fetch transcript for current video
  private async autoFetchTranscript(): Promise<void> {
    if (!this.currentVideoId) return;

    // Check cache first
    const cached = this.getCachedTranscript(this.currentVideoId);
    if (cached) {
      console.log('ReactYouTubeIntegrationService: Using cached transcript for', this.currentVideoId);

      // Don't auto-initialize overlay - let user click to show language selector
      console.log('ReactYouTubeIntegrationService: Transcript cached and ready for user interaction');
      return;
    }

    console.log('ReactYouTubeIntegrationService: Auto-fetching transcript for', this.currentVideoId);

    try {
      const transcriptResponse = await getVideoTranscript({
        videoUrl: this.currentVideoId,
      });

      const transcript = transcriptResponse;
      if (transcript && transcript.segments && transcript.segments.length > 0) {
        const videoTranscript: VideoTranscript = {
          videoId: transcript.videoId || this.currentVideoId,
          provider: transcript.provider || 'youtube',
          language: transcript.language || 'unknown',
          segments: transcript.segments.map(segment => ({
            text: segment.text,
            start: segment.startTime,
            end: segment.endTime,
            duration: segment.endTime - segment.startTime,
          })),
          available: transcript.available,
          source: transcript.source || 'api',
        };

        // Cache the transcript
        await this.cacheTranscript(videoTranscript);

        console.log('ReactYouTubeIntegrationService: Transcript auto-fetched and cached');

        // Don't auto-initialize overlay - let user click to show language selector
        console.log('ReactYouTubeIntegrationService: Transcript auto-fetched and cached, ready for user interaction');
      }
    } catch (error) {
      console.error('ReactYouTubeIntegrationService: Auto-fetch failed:', error);
      // Silently fail for auto-fetch - don't show error to user
    }
  }

  // Initialize the React transcript overlay
  private async initializeReactOverlay(transcript: VideoTranscript): Promise<void> {
    if (!this.currentVideoId) return;

    try {
      console.log('ReactYouTubeIntegrationService: Initializing React transcript overlay');

      // Create the React component using JSX syntax
      const overlayComponent = (
        <TranscriptOverlayManager
          videoId={this.currentVideoId}
          transcript={transcript}
          onReady={() => console.log('React overlay ready')}
          onError={(error: any) => console.error('React overlay error:', error)}
          onSegmentSeek={(seekTime: number) => this.handleSegmentSeek(seekTime)}
        />
      );

      // Render the React component
      await renderReactComponent(overlayComponent, EXTENSION_CLASSES.TRANSCRIPT_OVERLAY, 'react-transcript-overlay');

      console.log('ReactYouTubeIntegrationService: React transcript overlay started');
    } catch (error) {
      console.error('ReactYouTubeIntegrationService: Failed to initialize React overlay:', error);
    }
  }

  // Toggle React overlay
  private async toggleReactOverlay(): Promise<void> {
    if (!this.currentVideoId) return;

    const isCurrentlyMounted = document.querySelector(`.${EXTENSION_CLASSES.TRANSCRIPT_OVERLAY}`);

    if (isCurrentlyMounted) {
      unmountReactComponent(EXTENSION_CLASSES.TRANSCRIPT_OVERLAY);
      console.log('ReactYouTubeIntegrationService: React overlay hidden');
    } else {
      // Try to get cached transcript and show overlay
      const cached = this.getCachedTranscript(this.currentVideoId);
      if (cached) {
        await this.initializeReactOverlay(cached);
      } else {
        console.log('No cached transcript available for overlay');
      }
    }
  }

  // Handle segment seeking from React overlay
  private handleSegmentSeek(seekTime: number): void {
    // Find the video element and seek to the specified time
    const video = document.querySelector('video') as HTMLVideoElement;
    if (video) {
      video.currentTime = seekTime;
      console.log('ReactYouTubeIntegrationService: Seeked to', seekTime);
    }
  }

  // Get cached transcript for video
  private getCachedTranscript(videoId: string): VideoTranscript | null {
    const cached = this.transcriptCache[videoId];
    if (!cached) return null;

    // Check if cache is expired (1 hour)
    const now = Date.now();
    const CACHE_EXPIRY = 3600000; // 1 hour in ms

    if (now - cached.timestamp > CACHE_EXPIRY) {
      // Remove expired cache
      delete this.transcriptCache[videoId];
      this.saveTranscriptCache();
      return null;
    }

    return cached.transcript;
  }

  // Cache transcript for video
  private async cacheTranscript(transcript: VideoTranscript): Promise<void> {
    this.transcriptCache[transcript.videoId] = {
      transcript,
      translations: {},
      timestamp: Date.now(),
    };

    await this.saveTranscriptCache();
  }

  // Save transcript cache to storage
  private async saveTranscriptCache(): Promise<void> {
    try {
      await new Promise<void>(resolve => {
        chrome.storage.local.set(
          {
            [EXTENSION_STORAGE_KEYS.TRANSCRIPT_CACHE]: this.transcriptCache,
          },
          resolve
        );
      });
    } catch (error) {
      console.error('Failed to save transcript cache:', error);
    }
  }

  // Public methods
  get isActive(): boolean {
    return this.isIntegrationActive;
  }

  get overlayActive(): boolean {
    return document.querySelector(`.${EXTENSION_CLASSES.TRANSCRIPT_OVERLAY}`) !== null;
  }

  async toggleOverlay(): Promise<void> {
    await this.toggleReactOverlay();
  }

  async setOverlayEnabled(enabled: boolean): Promise<void> {
    this.overlayEnabled = enabled;

    // Save to storage
    try {
      await new Promise<void>(resolve => {
        chrome.storage.local.set(
          {
            [EXTENSION_STORAGE_KEYS.TRANSCRIPT_OVERLAY_ENABLED]: enabled,
          },
          resolve
        );
      });
    } catch (error) {
      console.error('Failed to save overlay setting:', error);
    }

    // Apply the setting
    if (enabled && this.currentVideoId) {
      const cached = this.getCachedTranscript(this.currentVideoId);
      if (cached) {
        await this.initializeReactOverlay(cached);
      }
    } else {
      unmountReactComponent(EXTENSION_CLASSES.TRANSCRIPT_OVERLAY);
    }
  }
}
