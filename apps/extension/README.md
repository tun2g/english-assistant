# English Learning Assistant - Chrome Extension

A Chrome extension template for the English learning platform, built with TypeScript and Vite.

## Features (Template)

- 🎯 **Background Service Worker** - Handles extension lifecycle and messaging
- 📄 **Content Scripts** - Interact with web pages
- 🖼️ **Popup Interface** - Quick access extension UI
- ⚙️ **Options Page** - Extension settings and configuration
- 💾 **Storage Integration** - Uses shared storage utilities
- 🔧 **TypeScript** - Type-safe development
- ⚡ **Vite Build System** - Fast development and building

## Development

### Prerequisites

- Node.js 18+
- pnpm 8+

### Setup

```bash
# Install dependencies (from monorepo root)
pnpm install

# Build extension
cd apps/extension
pnpm build

# Development build (watch mode)
pnpm dev
```

### Loading in Chrome

1. Open Chrome and navigate to `chrome://extensions/`
2. Enable "Developer mode"
3. Click "Load unpacked"
4. Select the `dist` folder

### Project Structure

```
src/
├── background/
│   └── background.ts    # Service worker
├── content/
│   └── content.ts       # Content script
├── popup/
│   ├── popup.html       # Popup UI
│   ├── popup.ts         # Popup logic
│   └── popup.css        # Popup styles
├── options/
│   ├── options.html     # Options page UI
│   ├── options.ts       # Options page logic
│   └── options.css      # Options page styles
├── types/
│   ├── extension.ts     # Extension-specific types
│   └── global.d.ts      # Global type declarations
├── utils/
│   └── storage.ts       # Storage utilities
└── manifest.json        # Extension manifest
```

### Available Scripts

- `pnpm build` - Build for production
- `pnpm dev` - Build and watch for changes
- `pnpm lint` - Run ESLint
- `pnpm type-check` - Run TypeScript type checking
- `pnpm clean` - Clean build directory

## Dependencies

### Runtime Dependencies
- `@english/shared` - Shared utilities and storage

### Development Dependencies
- TypeScript for type safety
- Vite for building
- ESLint for code linting

## Extension Capabilities

The extension template includes support for:

- ✅ Storage API (via shared utilities)
- ✅ Tabs API
- ✅ Runtime messaging
- ✅ Context menus
- ✅ Keyboard shortcuts
- ✅ Notifications
- ✅ Script injection

## TODOs

- [ ] Implement specific learning features
- [ ] Add vocabulary tracking
- [ ] Implement tab management
- [ ] Add translation features
- [ ] Create settings interface
- [ ] Add data import/export
- [ ] Implement offline functionality

## Architecture

The extension follows Chrome Extension Manifest V3 standards and integrates with the monorepo's shared packages for consistent storage and utilities across the platform.