// Main Content Script for English Learning Assistant - YouTube Integration
import '../styles/globals.scss';
import { ReactYouTubeIntegrationService } from './features/youtube-integration/react-youtube-integration-service';

import { extractVideoId } from './utils/video/video-utils';
import type { ExtensionMessage, MessageResponse } from '../shared/types/extension-types';
import { EXTENSION_MESSAGES } from '@/shared/constants';

// YouTube Integration Service instance
let youtubeIntegration: ReactYouTubeIntegrationService | null = null;
let initializationInProgress: boolean = false;

// Message listener for communication with popup and background
chrome.runtime.onMessage.addListener(
  (
    request: ExtensionMessage,
    _sender: chrome.runtime.MessageSender,
    sendResponse: (response?: MessageResponse) => void
  ) => {
    switch (request.action) {
      case EXTENSION_MESSAGES.GET_PAGE_INFO:
        sendResponse({
          success: true,
          data: {
            url: window.location.href,
            title: document.title,
            isYouTube: window.location.hostname.includes('youtube.com'),
            videoId: extractVideoId(window.location.href, 'youtube'),
          },
        });
        break;

      case EXTENSION_MESSAGES.TOGGLE_TRANSLATION:
        handleToggleTranslation(request.enabled || false);
        sendResponse({ success: true });
        break;

      case EXTENSION_MESSAGES.GET_TRANSCRIPT_WITH_AUTH:
        // This message is handled internally by the YouTube integration
        sendResponse({
          success: true,
          message: 'Transcript request forwarded',
        });
        break;

      case EXTENSION_MESSAGES.TOGGLE_TRANSCRIPT_OVERLAY:
        if (youtubeIntegration) {
          youtubeIntegration
            .toggleOverlay()
            .then(() => {
              sendResponse({ success: true, message: 'Overlay toggled' });
            })
            .catch(error => {
              console.error('Failed to toggle overlay:', error);
              sendResponse({
                success: false,
                message: 'Failed to toggle overlay',
              });
            });
        } else {
          sendResponse({
            success: false,
            message: 'YouTube integration not active',
          });
        }
        break;

      default:
        sendResponse({ success: true, message: 'Unknown action' });
    }

    return true;
  }
);

// Handle translation toggle
function handleToggleTranslation(enabled: boolean): void {
  if (enabled) {
    initializeYouTubeIntegration();
  } else {
    destroyYouTubeIntegration();
  }
}

// Initialize content script
function initializeContentScript(): void {
  if (window.location.hostname.includes('youtube.com')) {
    initializeYouTubeIntegration();
  }
}

// Initialize YouTube-specific functionality
async function initializeYouTubeIntegration(): Promise<void> {
  // Prevent concurrent initialization
  if (initializationInProgress) {
    return;
  }

  // If service exists and is active, don't recreate it
  if (youtubeIntegration && youtubeIntegration.isActive) {
    return;
  }

  try {
    initializationInProgress = true;

    // Only destroy if we're creating a new one
    if (youtubeIntegration) {
      youtubeIntegration.destroy();
    }

    youtubeIntegration = new ReactYouTubeIntegrationService();
    await youtubeIntegration.init();
  } catch (error) {
    // Silent error handling
  } finally {
    initializationInProgress = false;
  }
}

// Destroy YouTube integration
function destroyYouTubeIntegration(): void {
  if (youtubeIntegration) {
    youtubeIntegration.destroy();
    youtubeIntegration = null;
  }
}

// Wait for DOM to be ready
if (document.readyState === 'loading') {
  document.addEventListener('DOMContentLoaded', initializeContentScript);
} else {
  initializeContentScript();
}
