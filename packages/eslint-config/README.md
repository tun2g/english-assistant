# @english/eslint-config

Shared ESLint configurations for the English learning platform monorepo.

## Usage

### Base Configuration

For general TypeScript projects:

```json
{
  "extends": ["@english/eslint-config"]
}
```

### React Configuration

For React + TypeScript projects:

```json
{
  "extends": ["@english/eslint-config/react"]
}
```

### Node.js Configuration

For Node.js backend projects:

```json
{
  "extends": ["@english/eslint-config/node"]
}
```

### Chrome Extension Configuration

For Chrome extensions with React:

```json
{
  "extends": ["@english/eslint-config/chrome-extension"]
}
```

## Available Configurations

- **`@english/eslint-config`** - Base configuration with TypeScript support
- **`@english/eslint-config/base`** - Same as above (explicit)
- **`@english/eslint-config/react`** - React + TypeScript configuration
- **`@english/eslint-config/node`** - Node.js backend configuration
- **`@english/eslint-config/chrome-extension`** - Chrome extension configuration

## Features

### Base Configuration
- TypeScript support with `@typescript-eslint`
- Modern ES2021+ syntax
- Consistent code formatting rules
- Security best practices
- Performance optimizations

### React Configuration
- React 17+ support (no need for React import)
- React Hooks rules
- JSX best practices
- Component naming conventions
- React Refresh support for Vite

### Node.js Configuration
- Node.js specific environment
- Server-side best practices
- Security rules for backend code

### Chrome Extension Configuration
- Chrome extension environment support
- Service worker specific rules
- Content script validation
- CSP compliance rules
- Chrome API usage patterns

## Rules Overview

### Code Quality
- Prefer `const` over `let` and `var`
- No unused variables (with underscore prefix exception)
- Consistent indentation (2 spaces)
- Single quotes for strings
- Semicolons required

### Security
- No `eval()` usage
- No implied eval
- No script URLs
- No dangerous HTML operations

### Performance
- Prefer template literals
- Prefer arrow functions
- No duplicate imports
- Efficient React patterns

## Installation

This package is automatically available in the monorepo workspace. For external use:

```bash
npm install @english/eslint-config --save-dev
```

## Peer Dependencies

Make sure to install the required peer dependencies:

```bash
npm install eslint @typescript-eslint/eslint-plugin @typescript-eslint/parser --save-dev
```

For React projects, also install:

```bash
npm install eslint-plugin-react eslint-plugin-react-hooks eslint-plugin-react-refresh --save-dev
```

## Customization

You can extend or override rules in your local `.eslintrc.json`:

```json
{
  "extends": ["@english/eslint-config/react"],
  "rules": {
    "no-console": "off",
    "prefer-const": "warn"
  }
}
```

## Development

To modify configurations:

1. Edit the relevant `.js` file in this package
2. Test changes across the monorepo
3. Update version in `package.json`
4. Document changes in this README