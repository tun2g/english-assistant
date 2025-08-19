module.exports = {
  extends: ['./react.js'],
  env: {
    browser: true,
    webextensions: true,
  },
  globals: {
    chrome: 'readonly',
  },
  rules: {
    // Chrome Extension specific rules
    'no-console': 'warn', // Console logs are useful for debugging extensions
    
    // Security rules for extensions
    'no-eval': 'error',
    'no-implied-eval': 'error',
    'no-new-func': 'error',
    'no-script-url': 'error',
    
    // Chrome API usage rules
    '@typescript-eslint/no-explicit-any': 'warn', // Chrome APIs sometimes use any
    
    // Content Security Policy compliance
    'no-inline-comments': 'off',
    'no-undef': 'error',
  },
  overrides: [
    {
      files: ['**/background/**/*.ts', '**/service-worker.ts'],
      env: {
        browser: true,
        webextensions: true,
        serviceworker: true,
      },
      rules: {
        'no-restricted-globals': [
          'error',
          {
            name: 'window',
            message: 'window is not available in service workers',
          },
          {
            name: 'document',
            message: 'document is not available in service workers',
          },
        ],
      },
    },
    {
      files: ['**/content/**/*.ts'],
      env: {
        browser: true,
        webextensions: true,
      },
      rules: {
        // Content scripts have access to both page and extension contexts
        'no-undef': 'error',
      },
    },
  ],
};