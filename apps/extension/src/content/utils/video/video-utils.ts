import { VIDEO_URL_PATTERNS } from '../../../shared/constants';

// Utility function to extract video ID from URL
export function extractVideoId(url: string, provider: string): string | null {
  const trimmedUrl = url.trim();

  if (provider.toUpperCase() === 'YOUTUBE') {
    const patterns = VIDEO_URL_PATTERNS[provider.toUpperCase()];
    if (patterns) {
      for (const pattern of patterns) {
        const match = trimmedUrl.match(pattern);
        if (match) {
          return match[1];
        }
      }
    }
  }

  return null;
}

// Extract video ID from current URL
export function extractVideoIdFromCurrentUrl(): string | null {
  const url = window.location.href;
  const match = url.match(/[?&]v=([^&]+)/);
  return match ? match[1] : null;
}

// Check if current page is YouTube
export function isYouTubePage(): boolean {
  return window.location.hostname.includes('youtube.com');
}

// Check if URL contains video ID
export function hasVideoId(url: string): boolean {
  return extractVideoId(url, 'youtube') !== null;
}
