# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this
repository.

## Project Overview

This is a Chrome Extension V3 project for the English Learning Assistant platform, built with
TypeScript, React, Framework7 React, and Vite. The extension provides YouTube integration for
language learning with features like transcript extraction, translation, and OAuth authentication.

## Development Commands

### Extension Development

```bash
# Development build with watch mode
pnpm dev

# Production build
pnpm build

# Development build only
pnpm build:dev

# Type checking
pnpm type-check

# Linting with auto-fix
pnpm lint

# Clean build artifacts
pnpm clean
```

### Testing Extension in Chrome

```bash
# 1. Build extension
pnpm build:dev

# 2. Load extension for testing:
# - Open Chrome -> chrome://extensions/
# - Enable "Developer mode"
# - Click "Load unpacked" -> select dist/ folder

# 3. Debug extension:
# - Background script: chrome://extensions/ -> "Service Worker" link
# - Content script: Browser DevTools -> Console tab
# - Popup: Right-click extension icon -> "Inspect popup"
```

## Architecture

### Extension Structure

The extension follows Chrome Extension Manifest V3 standards with these main components:

- **Background Service Worker** (`src/background/background.ts`) - Handles extension lifecycle,
  messaging, and OAuth tab management
- **Content Scripts** (`src/content/content-script.ts`) - YouTube page integration and DOM
  manipulation
- **Popup UI** (`src/popup/`) - Framework7 React-based popup interface
- **Options Page** (`src/options/`) - Extension settings and configuration

### Key Architectural Patterns

**Service-Oriented Architecture**: Content script functionality is organized into feature-specific
services:

```text
content/features/
├── auth/                    # OAuth authentication
├── notifications/           # User notifications
├── transcript-sync/         # Transcript synchronization
├── translation/             # Translation services
├── video-tracking/          # Video progress tracking
└── youtube-integration/     # Main YouTube integration service
```

**Component-Based UI**: React components with Framework7 React for mobile-like UI:

```text
popup/components/
├── navigation/              # Quick actions and settings
├── oauth/                   # OAuth status and controls
└── youtube/                # YouTube-specific features
```

**Shared Architecture**: Integrates with monorepo shared packages:

- `@english/shared` - API clients, storage utilities, types, constants
- `@english/ui` - Reusable UI components

### Core Services

**ReactYouTubeIntegrationService**: Main orchestrator that manages:

- OAuth authentication flow
- Video monitoring and transcript extraction
- Player controls injection
- Translation overlay management
- Content script to popup communication

**OAuthManager**: Handles YouTube API authentication using OAuth 2.0 flow with background script
coordination.

**PlayerControlsManager**: Injects custom UI controls into YouTube player interface.

**VideoMonitor**: Tracks video changes and manages video-specific state.

## File Naming Conventions

**ALWAYS use kebab-case for ALL filenames and directories:**

```text
✅ oauth-manager.ts, transcript-overlay.tsx, video-utils.ts
❌ OAuthManager.ts, transcriptOverlay.tsx, video_utils.ts
```

**Component Structure Patterns:**

- **Feature services**: `{feature-name}-{service-type}.ts`
  - `oauth-manager.ts`, `notification-manager.ts`
- **React components**: `{descriptive-name}.tsx`
  - `transcript-overlay-manager.tsx`, `youtube-section.tsx`
- **Utilities**: `{category}/{name}-{category}.ts`
  - `dom/dom-utils.ts`, `video/video-utils.ts`

## TypeScript Conventions

**Interface over Type**: Use `interface` for object shapes, `type` for unions/primitives:

```typescript
// ✅ Use interface for objects
interface OAuthConfig {
  clientId: string;
  scopes: string[];
}

// ✅ Use type for unions
type ExtensionMessage = 'GET_PAGE_INFO' | 'TOGGLE_TRANSLATION';
```

**Constants Pattern**: Use `as const` instead of `enum`:

```typescript
// ✅ Preferred pattern
const EXTENSION_MESSAGES = {
  GET_PAGE_INFO: 'GET_PAGE_INFO',
  TOGGLE_TRANSLATION: 'TOGGLE_TRANSLATION',
} as const;
```

## Extension Development Patterns

### Message Passing Architecture

The extension uses Chrome's message passing API for communication between components:

```typescript
// Content script -> Background
chrome.runtime.sendMessage({ action: 'OPEN_TAB', url: authUrl });

// Popup -> Content script
chrome.tabs.sendMessage(tabId, { action: 'TOGGLE_TRANSLATION', enabled: true });
```

### Storage Patterns

Uses Chrome storage API with shared storage utilities:

```typescript
// Get storage values
chrome.storage.local.get([EXTENSION_STORAGE_KEYS.AUTO_TRANSLATE_ENABLED]);

// Set storage values
chrome.storage.local.set({ [EXTENSION_STORAGE_KEYS.OAUTH_TOKEN]: token });
```

### React Component Integration

Content scripts use `react-renderer.tsx` utility for injecting React components into web pages:

```typescript
// Render React component in content script
await renderReactComponent(
  <TranscriptOverlay transcript={transcript} />,
  'transcript-overlay-root'
);

// Cleanup on destroy
cleanupAllReactComponents();
```

### Error Handling Pattern

Centralized error handling with structured logging:

```typescript
import { captureError, logDebug } from './utils/error-handler';

try {
  await apiCall();
} catch (error) {
  captureError('Operation failed', error, { context: 'additional-info' });
}
```

## Build Configuration

### Vite + CRXJS Setup

The extension uses Vite with CRXJS plugin for Chrome extension specific bundling:

- **Manifest Configuration**: `src/manifest.config.ts` - Defines permissions, content scripts, and
  extension metadata
- **Build Output**: Optimized for Chrome extension with proper chunk splitting and CSP compliance
- **Development**: Hot module reloading for popup and options page, content script changes require
  extension reload

### Key Build Features

- ES2020 target for modern Chrome extension support
- Disabled CSS code splitting to avoid web_accessible_resources issues
- Custom chunk naming for consistent asset references
- Framework7 React integration with automatic JSX runtime

## Extension Capabilities

The extension has permissions for:

- **Storage API**: Local data persistence
- **Tabs API**: Tab management and communication
- **ActiveTab**: Current tab interaction
- **ContextMenus**: Right-click menu integration
- **Notifications**: User notifications
- **Scripting**: Dynamic script injection
- **Host Permissions**: YouTube.com integration, localhost development

## Framework7 React Integration

The popup uses Framework7 React for mobile-like UI components:

```typescript
import { App, View, Page, Navbar, Block } from 'framework7-react';
import 'framework7/css/bundle';

// Initialize Framework7 app
const f7params = {
  name: 'English Learning Assistant',
  theme: 'auto',
};
```

## Development Workflow

1. **Start Development**: Run `pnpm dev` for watch mode building
2. **Load Extension**: Load `dist/` folder in Chrome developer mode
3. **Test Changes**: Popup/options changes hot reload, content script requires extension reload
4. **Debug**: Use Chrome DevTools for each extension context (background, content, popup)
5. **Build Production**: Run `pnpm build` for optimized production bundle

## Common Patterns

### Content Script Initialization

```typescript
// Wait for DOM ready
if (document.readyState === 'loading') {
  document.addEventListener('DOMContentLoaded', initializeContentScript);
} else {
  initializeContentScript();
}
```

### Service Management

```typescript
// Singleton pattern for service instances
let youtubeIntegration: ReactYouTubeIntegrationService | null = null;

async function initializeService() {
  if (youtubeIntegration?.isActive) return;

  youtubeIntegration = new ReactYouTubeIntegrationService();
  await youtubeIntegration.init();
}
```

### React Hook Patterns

```typescript
// Custom hooks for extension state
function usePageInfo() {
  const [pageInfo, setPageInfo] = useState(null);

  useEffect(() => {
    // Get page info from content script
    chrome.tabs.query({ active: true, currentWindow: true }, ([tab]) => {
      chrome.tabs.sendMessage(tab.id, { action: 'GET_PAGE_INFO' });
    });
  }, []);

  return { pageInfo };
}
```
