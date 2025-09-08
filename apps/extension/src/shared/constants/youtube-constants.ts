// YouTube selectors for extension integration
export const YOUTUBE_SELECTORS = {
  RIGHT_CONTROLS: [
    '#movie_player .ytp-right-controls',
    '.ytp-chrome-bottom .ytp-right-controls',
    '.ytp-chrome-controls .ytp-right-controls',
    '.html5-video-player .ytp-right-controls',
    'div.ytp-right-controls',
    '.ytp-chrome-bottom div[class*="controls"]',
  ],
  SETTINGS_BUTTON: [
    '.ytp-settings-button',
    'button[title*="Settings"]',
    'button[aria-label*="Settings"]',
    'button[aria-label*="Cài đặt"]', // Vietnamese for Settings
    '.ytp-chrome-controls button[class*="settings"]',
    '.ytp-right-controls .ytp-settings-button',
  ],
  FULLSCREEN_BUTTON: [
    '.ytp-fullscreen-button',
    'button[title*="Fullscreen"]',
    'button[aria-label*="Fullscreen"]',
    'button[aria-label*="Toàn màn hình"]', // Vietnamese for Fullscreen
    '.ytp-right-controls .ytp-fullscreen-button',
  ],
  MOVIE_PLAYER: '#movie_player',
  VIDEO_PLAYER: '.html5-video-player',
  VIDEO_ELEMENT: 'video',
  PLAYER_CONTAINER: ['#movie_player', '.html5-video-player', '#player'],
  // Additional selectors for better positioning
  RIGHT_CONTROLS_LEFT: '.ytp-right-controls-left',
  RIGHT_CONTROLS_RIGHT: '.ytp-right-controls-right',
} as const;
