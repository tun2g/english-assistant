import { VideoApiClient, createVideoApiClient } from '@english/shared/api/video-api';
import type { VideoApiConfig } from '@english/shared/types/video-types';

import { OAuthManager } from '../auth/oauth-manager';
import { NotificationManager } from '../notifications/notification-manager';
import { AuthModal } from '../../ui/modals/auth-modal';
import { TranscriptModal } from '../../ui/modals/transcript-modal';
import { ExtensionPanel } from '../../ui/panels/extension-panel';
import { PlayerControlsManager } from './player-controls-manager';
import { VideoMonitor } from './video-monitor';

import { 
  EXTENSION_STORAGE_KEYS, 
  DEFAULT_API_BASE_URL 
} from '../../../shared/constants/extension-constants';
import { extractVideoIdFromCurrentUrl } from '../../utils/video/video-utils';
import type { VideoTranscript } from '../../../shared/types/extension-types';

export class YouTubeIntegrationService {
  private currentVideoId: string | null = null;
  private isIntegrationActive: boolean = false;
  private apiClient: VideoApiClient | null = null;
  
  // Feature managers
  private oauthManager: OAuthManager;
  private playerControlsManager: PlayerControlsManager;
  private videoMonitor: VideoMonitor;

  constructor() {
    this.oauthManager = new OAuthManager();
    
    this.playerControlsManager = new PlayerControlsManager({
      onButtonClick: () => this.handleExtensionButtonClick()
    });

    this.videoMonitor = new VideoMonitor({
      onVideoChange: (videoId) => this.handleVideoChange(videoId)
    });
  }

  // Initialize YouTube integration
  async init(): Promise<void> {
    console.log('YouTubeIntegrationService: Starting initialization...');
    
    if (this.isIntegrationActive) {
      console.log('YouTubeIntegrationService: Already active, skipping');
      return;
    }

    this.isIntegrationActive = true;
    
    // Initialize API client
    await this.initializeApiClient();
    
    // Start video monitoring
    this.videoMonitor.start();
    
    // Check current video
    this.checkCurrentVideo();
    
    console.log('YouTubeIntegrationService: Initialization complete');
  }

  // Clean up integration
  destroy(): void {
    console.log('YouTubeIntegrationService: Destroying...');
    
    this.isIntegrationActive = false;
    
    // Clean up managers
    this.videoMonitor.stop();
    this.playerControlsManager.removeControls();
    this.oauthManager.destroy();
    
    console.log('YouTubeIntegrationService: Destroyed');
  }

  // Initialize API client
  private async initializeApiClient(): Promise<void> {
    if (this.apiClient) return;

    // Get API base URL from storage
    const baseUrl = await new Promise<string>((resolve) => {
      chrome.storage.local.get([EXTENSION_STORAGE_KEYS.API_BASE_URL], (result) => {
        resolve(result[EXTENSION_STORAGE_KEYS.API_BASE_URL] || DEFAULT_API_BASE_URL);
      });
    });

    const config: VideoApiConfig = {
      baseUrl,
      timeout: 30000,
    };

    this.apiClient = createVideoApiClient(config);
  }

  // Handle video change
  private async handleVideoChange(videoId: string | null): Promise<void> {
    if (videoId !== this.currentVideoId) {
      this.currentVideoId = videoId;
      console.log('YouTubeIntegrationService: Video ID changed to:', videoId);
      
      if (videoId) {
        console.log('YouTubeIntegrationService: Initializing video integration...');
        await this.initializeVideoIntegration();
      } else {
        console.log('YouTubeIntegrationService: No video ID, cleaning up...');
        this.playerControlsManager.removeControls();
      }
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
  }

  // Handle extension button click
  private async handleExtensionButtonClick(): Promise<void> {
    if (!this.currentVideoId) {
      NotificationManager.showError('No video detected!');
      return;
    }

    this.createExtensionPanel();
  }

  // Create extension panel
  private createExtensionPanel(): void {
    const panel = new ExtensionPanel(
      {
        onOAuthConnect: () => this.handleOAuthConnect(),
        onTranscriptRequest: () => this.handleTranscriptRequest(),
        onClose: () => console.log('Panel closed')
      },
      this.oauthManager.authenticated
    );

    panel.show();
  }

  // Handle OAuth connect
  private async handleOAuthConnect(): Promise<void> {
    try {
      await this.oauthManager.connect();
    } catch (error) {
      console.error('OAuth connection failed:', error);
    }
  }

  // Handle transcript request
  private async handleTranscriptRequest(): Promise<void> {
    if (!this.currentVideoId) {
      NotificationManager.showError('No video detected!');
      return;
    }

    // If not authenticated, show authentication prompt
    if (!this.oauthManager.authenticated) {
      this.showAuthenticationPrompt();
      return;
    }

    await this.fetchAndShowTranscript();
  }

  // Show authentication prompt
  private showAuthenticationPrompt(): void {
    const authModal = new AuthModal({
      onConnect: async () => {
        try {
          await this.oauthManager.connect();
          // Remove extension panel
          const panel = document.querySelector('.english-learning-panel');
          panel?.remove();
        } catch (error) {
          NotificationManager.showError('Failed to start authentication: ' + (error as Error).message);
        }
      },
      onCancel: () => console.log('Authentication cancelled')
    });

    authModal.show();
  }

  // Fetch and show transcript
  private async fetchAndShowTranscript(): Promise<void> {
    if (!this.apiClient || !this.currentVideoId) return;

    NotificationManager.showLoading('Fetching transcript...');

    try {
      const transcript = await this.apiClient.getTranscript({ 
        videoUrl: this.currentVideoId 
      });
      
      NotificationManager.hide();
      
      if (transcript && transcript.segments && transcript.segments.length > 0) {
        // Convert to our transcript format
        const videoTranscript: VideoTranscript = {
          videoId: transcript.videoId || this.currentVideoId,
          provider: transcript.provider || 'youtube',
          language: transcript.language || 'unknown',
          segments: transcript.segments.map(segment => ({
            text: segment.text,
            start: segment.startTime,
            end: segment.endTime,
            duration: segment.endTime - segment.startTime
          })),
          available: transcript.available,
          source: transcript.source || 'api'
        };

        TranscriptModal.show(videoTranscript);
      } else {
        NotificationManager.showWarning('No transcript available for this video');
      }
    } catch (error: any) {
      console.error('Failed to fetch transcript:', error);
      NotificationManager.hide();
      
      // Check if it's an authentication error
      if (this.oauthManager.isAuthenticationError(error)) {
        await this.oauthManager.checkAuthStatus();
        this.showAuthenticationPrompt();
      } else {
        NotificationManager.showError('Failed to get transcript: ' + (error.message || 'Unknown error'));
      }
    }
  }
}