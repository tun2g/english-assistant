module.exports = {
  extends: ['./base.js'],
  env: {
    node: true,
    browser: false,
  },
  rules: {
    // Node.js specific rules
    'no-console': 'off', // Console is allowed in Node.js
    'no-process-exit': 'error',
    'no-process-env': 'off',
    
    // Import/Export rules for Node.js
    'prefer-const': 'error',
    'no-var': 'error',
    
    // Security rules
    'no-eval': 'error',
    'no-implied-eval': 'error',
    'no-new-func': 'error',
  },
};