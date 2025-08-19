import type { VideoProvider } from '../types/video-types';
import { VIDEO_URL_PATTERNS, VIDEO_PROVIDERS } from '../constants/video-constants';

/**
 * Utility functions for video operations
 */

/**
 * Extract video ID from various video platform URLs
 */
export function extractVideoId(url: string, provider?: VideoProvider): string | null {
  const trimmedUrl = url.trim();
  
  if (provider) {
    // Use specific provider patterns
    const patterns = VIDEO_URL_PATTERNS[provider.toUpperCase() as keyof typeof VIDEO_URL_PATTERNS];
    for (const pattern of patterns) {
      const match = trimmedUrl.match(pattern);
      if (match) {
        return match[1];
      }
    }
    return null;
  }

  // Try all patterns
  for (const [providerName, patterns] of Object.entries(VIDEO_URL_PATTERNS)) {
    for (const pattern of patterns) {
      const match = trimmedUrl.match(pattern);
      if (match) {
        return match[1];
      }
    }
  }

  // If no pattern matches, assume it's already a video ID
  return trimmedUrl;
}

/**
 * Detect video provider from URL
 */
export function detectVideoProvider(url: string): VideoProvider | null {
  const lowerUrl = url.toLowerCase();
  
  if (lowerUrl.includes('youtube.com') || lowerUrl.includes('youtu.be') || lowerUrl.includes('youtube-nocookie.com')) {
    return VIDEO_PROVIDERS.YOUTUBE as VideoProvider;
  }
  
  if (lowerUrl.includes('vimeo.com')) {
    return VIDEO_PROVIDERS.VIMEO as VideoProvider;
  }
  
  if (lowerUrl.includes('twitch.tv')) {
    return VIDEO_PROVIDERS.TWITCH as VideoProvider;
  }
  
  return null;
}

/**
 * Check if URL is a supported video URL
 */
export function isSupportedVideoUrl(url: string): boolean {
  return detectVideoProvider(url) !== null;
}

/**
 * Validate video ID format for specific provider
 */
export function validateVideoId(videoId: string, provider: VideoProvider): boolean {
  switch (provider) {
    case VIDEO_PROVIDERS.YOUTUBE:
      // YouTube video IDs are 11 characters long and contain alphanumeric characters, hyphens, and underscores
      return /^[a-zA-Z0-9_-]{11}$/.test(videoId);
    
    case VIDEO_PROVIDERS.VIMEO:
      // Vimeo video IDs are numeric
      return /^\d+$/.test(videoId);
    
    case VIDEO_PROVIDERS.TWITCH:
      // Twitch video IDs are numeric
      return /^\d+$/.test(videoId);
    
    default:
      return false;
  }
}

/**
 * Normalize video URL to a standard format
 */
export function normalizeVideoUrl(url: string): string {
  const provider = detectVideoProvider(url);
  const videoId = extractVideoId(url, provider || undefined);
  
  if (!provider || !videoId) {
    return url; // Return original if we can't parse it
  }
  
  switch (provider) {
    case VIDEO_PROVIDERS.YOUTUBE:
      return `https://www.youtube.com/watch?v=${videoId}`;
    
    case VIDEO_PROVIDERS.VIMEO:
      return `https://vimeo.com/${videoId}`;
    
    case VIDEO_PROVIDERS.TWITCH:
      return `https://www.twitch.tv/videos/${videoId}`;
    
    default:
      return url;
  }
}

/**
 * Generate thumbnail URL for video
 */
export function generateThumbnailUrl(videoId: string, provider: VideoProvider, quality: 'default' | 'medium' | 'high' = 'medium'): string {
  switch (provider) {
    case VIDEO_PROVIDERS.YOUTUBE:
      const qualityMap = {
        default: 'default',
        medium: 'mqdefault',
        high: 'hqdefault',
      };
      return `https://img.youtube.com/vi/${videoId}/${qualityMap[quality]}.jpg`;
    
    case VIDEO_PROVIDERS.VIMEO:
      // Vimeo thumbnails require API call, return placeholder
      return `https://vumbnail.com/${videoId}.jpg`;
    
    case VIDEO_PROVIDERS.TWITCH:
      // Twitch thumbnails are more complex, return placeholder
      return `https://static-cdn.jtvnw.net/ttv-static/404_preview-480x272.jpg`;
    
    default:
      return '';
  }
}

/**
 * Format duration from seconds to human-readable string
 */
export function formatDuration(seconds: number): string {
  const hours = Math.floor(seconds / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);
  const secs = Math.floor(seconds % 60);
  
  if (hours > 0) {
    return `${hours}:${minutes.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
  } else {
    return `${minutes}:${secs.toString().padStart(2, '0')}`;
  }
}

/**
 * Format timestamp from milliseconds to time string
 */
export function formatTimestamp(milliseconds: number): string {
  const totalSeconds = Math.floor(milliseconds / 1000);
  const minutes = Math.floor(totalSeconds / 60);
  const seconds = totalSeconds % 60;
  
  return `${minutes}:${seconds.toString().padStart(2, '0')}`;
}

/**
 * Parse time string to milliseconds
 */
export function parseTimeToMilliseconds(timeString: string): number {
  const parts = timeString.split(':').map(Number);
  
  if (parts.length === 2) {
    // MM:SS format
    return (parts[0] * 60 + parts[1]) * 1000;
  } else if (parts.length === 3) {
    // HH:MM:SS format
    return (parts[0] * 3600 + parts[1] * 60 + parts[2]) * 1000;
  }
  
  return 0;
}

/**
 * Get language flag emoji
 */
export function getLanguageFlag(languageCode: string): string {
  const flagMap: Record<string, string> = {
    'en': '🇺🇸',
    'es': '🇪🇸',
    'fr': '🇫🇷',
    'de': '🇩🇪',
    'ja': '🇯🇵',
    'ko': '🇰🇷',
    'zh': '🇨🇳',
    'pt': '🇵🇹',
    'ru': '🇷🇺',
    'it': '🇮🇹',
    'ar': '🇸🇦',
    'hi': '🇮🇳',
    'th': '🇹🇭',
    'vi': '🇻🇳',
    'nl': '🇳🇱',
    'sv': '🇸🇪',
    'no': '🇳🇴',
    'da': '🇩🇰',
    'fi': '🇫🇮',
    'pl': '🇵🇱',
  };
  
  return flagMap[languageCode] || '🌐';
}

/**
 * Create video URL for sharing
 */
export function createShareableUrl(videoId: string, provider: VideoProvider, timestamp?: number): string {
  const baseUrl = normalizeVideoUrl(`${provider}:${videoId}`);
  
  if (timestamp && provider === VIDEO_PROVIDERS.YOUTUBE) {
    const seconds = Math.floor(timestamp / 1000);
    return `${baseUrl}&t=${seconds}s`;
  }
  
  return baseUrl;
}

/**
 * Check if video supports live streaming
 */
export function supportsLiveStreaming(provider: VideoProvider): boolean {
  return [VIDEO_PROVIDERS.YOUTUBE, VIDEO_PROVIDERS.TWITCH].includes(provider as any);
}

/**
 * Get provider display name
 */
export function getProviderDisplayName(provider: VideoProvider): string {
  const displayNames: Record<VideoProvider, string> = {
    youtube: 'YouTube',
    vimeo: 'Vimeo',
    twitch: 'Twitch',
  };
  
  return displayNames[provider] || provider;
}