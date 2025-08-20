// YouTube selectors for extension integration
export const YOUTUBE_SELECTORS = {
  RIGHT_CONTROLS: [
    '#movie_player .ytp-right-controls',
    '.ytp-chrome-bottom .ytp-right-controls',
    '.ytp-chrome-controls .ytp-right-controls',
    '.html5-video-player .ytp-right-controls',
    'div.ytp-right-controls',
  ],
  SETTINGS_BUTTON: '.ytp-settings-button',
  MOVIE_PLAYER: '#movie_player',
  VIDEO_PLAYER: '.html5-video-player',
  VIDEO_ELEMENT: 'video',
  PLAYER_CONTAINER: ['#movie_player', '.html5-video-player', '#player'],
} as const;
