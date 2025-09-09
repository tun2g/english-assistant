module.exports = {
  env: {
    browser: true,
    es2021: true,
    node: true,
  },
  extends: ['eslint:recommended'],
  parser: '@typescript-eslint/parser',
  parserOptions: {
    ecmaVersion: 'latest',
    sourceType: 'module',
  },
  plugins: ['@typescript-eslint'],
  rules: {
    // TypeScript rules
    '@typescript-eslint/no-unused-vars': ['error', { argsIgnorePattern: '^_' }],
    '@typescript-eslint/no-explicit-any': 'warn',
    '@typescript-eslint/explicit-function-return-type': 'off',
    '@typescript-eslint/explicit-module-boundary-types': 'off',
    '@typescript-eslint/no-non-null-assertion': 'warn',

    // General JavaScript/TypeScript rules
    'prefer-const': 'error',
    'no-var': 'error',
    'no-console': 'off',
    'no-debugger': 'error',
    'no-alert': 'error',

    // Code style - disabled in favor of Prettier
    indent: 'off',
    '@typescript-eslint/indent': 'off',
    quotes: 'off',
    semi: 'off',
    'comma-dangle': 'off',
    'object-curly-spacing': 'off',
    'array-bracket-spacing': 'off',

    // Best practices
    eqeqeq: ['error', 'always'],
    'no-duplicate-imports': 'error',
    'no-unused-expressions': 'error',
    'prefer-template': 'error',
    'prefer-arrow-callback': 'error',
  },
  ignorePatterns: ['node_modules/', 'dist/', 'build/', '*.min.js', 'coverage/'],
};
