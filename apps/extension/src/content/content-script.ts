// Main Content Script for English Learning Assistant - YouTube Integration
import { YouTubeIntegrationService } from './features/youtube-integration/youtube-integration-service';
import { 
  EXTENSION_STORAGE_KEYS, 
  EXTENSION_MESSAGES 
} from '../shared/constants/extension-constants';
import { extractVideoId } from './utils/video/video-utils';
import type { ExtensionMessage, MessageResponse } from '../shared/types/extension-types';

console.log('English Learning Assistant: Content script loaded');

// YouTube Integration Service instance
let youtubeIntegration: YouTubeIntegrationService | null = null;

// Message listener for communication with popup and background
chrome.runtime.onMessage.addListener((
  request: ExtensionMessage, 
  _sender: chrome.runtime.MessageSender, 
  sendResponse: (response?: MessageResponse) => void
) => {
  console.log('Content script received message:', request);
  
  switch (request.action) {
    case EXTENSION_MESSAGES.GET_PAGE_INFO:
      sendResponse({
        success: true,
        data: {
          url: window.location.href,
          title: document.title,
          isYouTube: window.location.hostname.includes('youtube.com'),
          videoId: extractVideoId(window.location.href, 'youtube'),
        }
      });
      break;
      
    case EXTENSION_MESSAGES.TOGGLE_TRANSLATION:
      handleToggleTranslation(request.enabled || false);
      sendResponse({ success: true });
      break;

    case EXTENSION_MESSAGES.GET_TRANSCRIPT_WITH_AUTH:
      // This message is handled internally by the YouTube integration
      sendResponse({ success: true, message: 'Transcript request forwarded' });
      break;
      
    default:
      sendResponse({ success: true, message: 'Unknown action' });
  }
  
  return true;
});

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
  console.log('Content script initialized');
  
  if (window.location.hostname.includes('youtube.com')) {
    initializeYouTubeIntegration();
  }
}

// Initialize YouTube-specific functionality
async function initializeYouTubeIntegration(): Promise<void> {
  console.log('Initializing YouTube integration');
  
  if (youtubeIntegration) {
    console.log('YouTube integration already exists, destroying first');
    youtubeIntegration.destroy();
  }

  youtubeIntegration = new YouTubeIntegrationService();
  await youtubeIntegration.init();
  
  // Check auto-translate setting
  chrome.storage.local.get([EXTENSION_STORAGE_KEYS.AUTO_TRANSLATE_ENABLED], (result) => {
    console.log('Auto-translate enabled:', result[EXTENSION_STORAGE_KEYS.AUTO_TRANSLATE_ENABLED]);
  });
  
  // Listen for storage changes
  chrome.storage.onChanged.addListener((changes) => {
    if (changes[EXTENSION_STORAGE_KEYS.AUTO_TRANSLATE_ENABLED]) {
      console.log('Auto-translate setting changed:', changes[EXTENSION_STORAGE_KEYS.AUTO_TRANSLATE_ENABLED].newValue);
    }
  });
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

// Handle page navigation in YouTube (SPA)
let currentUrl = location.href;
new MutationObserver(() => {
  if (location.href !== currentUrl) {
    currentUrl = location.href;
    
    if (window.location.hostname.includes('youtube.com')) {
      setTimeout(initializeYouTubeIntegration, 1000);
    }
  }
}).observe(document, { subtree: true, childList: true });