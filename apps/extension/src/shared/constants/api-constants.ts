// Video URL patterns for different providers
export const VIDEO_URL_PATTERNS: Record<string, RegExp[]> = {
  YOUTUBE: [/(?:youtube\.com\/watch\?v=|youtu\.be\/)([a-zA-Z0-9_-]{11})/],
};
