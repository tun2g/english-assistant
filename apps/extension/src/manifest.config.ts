import { defineManifest } from '@crxjs/vite-plugin';

export default defineManifest({
  manifest_version: 3,
  name: 'English Learning Assistant',
  version: '1.0.0',
  description: 'Chrome extension for English learning platform with tab management and vocabulary tracking',
  permissions: ['storage', 'tabs', 'activeTab', 'contextMenus', 'notifications', 'scripting'],
  host_permissions: ['http://localhost:*/*', 'https://*.english-learning.com/*', '*://*.youtube.com/*'],
  background: {
    service_worker: 'src/background/background.ts',
    type: 'module',
  },
  action: {
    default_popup: 'src/popup/popup.html',
  },
  content_scripts: [
    {
      matches: ['*://*.youtube.com/*'],
      js: ['src/content/content-script.ts'],
      run_at: 'document_idle',
    },
  ],
  web_accessible_resources: [
    {
      resources: ['*.js', '*.css', '*.html', 'assets/*', 'chunks/*.js'],
      matches: ['*://*.youtube.com/*'],
    },
  ],
  content_security_policy: {
    extension_pages:
      "script-src 'self'; object-src 'self'; style-src 'self' 'unsafe-inline'; font-src 'self' data:; img-src 'self' data: https:;",
  },
  options_page: 'src/options/options.html',
  commands: {
    'save-word': {
      suggested_key: {
        default: 'Ctrl+Shift+S',
        mac: 'Command+Shift+S',
      },
      description: 'Save selected word to vocabulary',
    },
    'toggle-extension': {
      suggested_key: {
        default: 'Ctrl+Shift+E',
        mac: 'Command+Shift+E',
      },
      description: 'Toggle extension on/off',
    },
  },
});
