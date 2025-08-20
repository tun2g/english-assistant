# Chrome Extension Development Guide

## Project Setup with React + Vite + Tailwind CSS + Shadcn/UI

This guide covers setting up and configuring a Chrome Extension using modern web technologies
including React 18, Vite, Tailwind CSS, and Shadcn/UI components.

## Current Architecture Overview

The English Learning Assistant Chrome Extension is built using:

- **Build Tool**: Vite 5.0.8 with `@crxjs/vite-plugin` for Chrome extension support
- **Framework**: React 18.2.0 with TypeScript 5.2.2
- **UI Library**: Framework7 React 8.3.4 for mobile-like interface
- **CSS Framework**: Tailwind CSS 3.4.1
- **Component Library**: Custom `@english/ui` package with Radix UI primitives
- **Manifest Version**: Chrome Extension Manifest V3

## Configuration Files

### 1. Package Configuration

```json
{
  "name": "@english/extension",
  "version": "1.0.0",
  "type": "module",
  "dependencies": {
    "@english/shared": "workspace:*",
    "@english/ui": "workspace:*",
    "framework7": "^8.3.4",
    "framework7-react": "^8.3.4",
    "react": "^18.2.0",
    "react-dom": "^18.2.0"
  },
  "devDependencies": {
    "@crxjs/vite-plugin": "^2.2.0",
    "@vitejs/plugin-react": "^4.2.1",
    "tailwindcss": "^3.4.1",
    "vite": "^5.0.8"
  }
}
```

### 2. Vite Configuration

```typescript
// vite.config.ts
import { crx } from '@crxjs/vite-plugin';
import react from '@vitejs/plugin-react';
import { resolve } from 'path';
import { defineConfig, loadEnv, UserConfig } from 'vite';
import { fileURLToPath, URL } from 'node:url';

import manifest from './src/manifest.config';

export default ({ mode = 'development' }: UserConfig) => {
  const isProduction = mode === 'production';

  return defineConfig({
    base: './',
    define: {
      'process.env.BUILD_TARGET': JSON.stringify('extension'),
      'process.env.NODE_ENV': JSON.stringify(mode),
      global: 'globalThis',
    },
    plugins: [
      react({
        jsxRuntime: 'automatic', // Enable automatic JSX runtime
      }),
      crx({
        manifest,
        contentScripts: {
          injectCss: true, // Enable content script HMR in development
        },
      }),
    ],
    build: {
      outDir: 'dist',
      minify: isProduction,
      target: 'es2020', // Modern browsers (Chrome extensions support ES2020+)
      rollupOptions: {
        output: {
          entryFileNames: '[name].js',
          chunkFileNames: '[name]-[hash].js',
          assetFileNames: '[name].[ext]',
        },
      },
      cssCodeSplit: false, // Disable CSS code splitting to avoid CSP issues
      chunkSizeWarningLimit: 1000,
    },
    resolve: {
      alias: {
        '@': resolve(fileURLToPath(new URL('.', import.meta.url)), 'src'),
        '@english/shared': resolve('../../packages/shared/src'),
        '@english/ui': resolve('../../packages/ui/src'),
      },
    },
    optimizeDeps: {
      include: ['react', 'react-dom'],
      exclude: ['@english/shared', '@english/ui'],
    },
  });
};
```

### 3. Tailwind CSS Configuration

```javascript
// tailwind.config.js
export default {
  content: [
    './src/**/*.{js,ts,jsx,tsx,html}',
    '../../packages/ui/src/**/*.{js,ts,jsx,tsx}', // Include shared UI package
  ],
  theme: {
    extend: {
      // Custom theme extensions can be added here
    },
  },
  plugins: [],
};
```

### 4. Manifest Configuration

```typescript
// src/manifest.config.ts
import { defineManifest } from '@crxjs/vite-plugin';

export default defineManifest({
  manifest_version: 3,
  name: 'English Learning Assistant',
  version: '1.0.0',
  description: 'Chrome extension for English learning platform',
  permissions: ['storage', 'tabs', 'activeTab', 'contextMenus', 'notifications', 'scripting'],
  host_permissions: [
    'http://localhost:*/*',
    'https://*.english-learning.com/*',
    '*://*.youtube.com/*',
  ],
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
    extension_pages: "script-src 'self'; object-src 'self'; style-src 'self' 'unsafe-inline';",
  },
  options_page: 'src/options/options.html',
});
```

## Setting Up Shadcn/UI Components

### 1. Install Shadcn/UI Dependencies

The project uses a custom `@english/ui` package that includes Shadcn/UI-compatible components:

```json
// packages/ui/package.json
{
  "dependencies": {
    "@radix-ui/react-slot": "^1.0.2",
    "class-variance-authority": "^0.7.0",
    "clsx": "^2.0.0",
    "tailwind-merge": "^2.2.0",
    "lucide-react": "^0.307.0"
  }
}
```

### 2. Utility Functions

Create utility functions similar to Shadcn/UI's `cn` helper:

```typescript
// packages/ui/src/lib/utils.ts
import { type ClassValue, clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}
```

### 3. Component Structure

Follow Shadcn/UI patterns for component creation:

```typescript
// packages/ui/src/components/button.tsx
import * as React from "react"
import { Slot } from "@radix-ui/react-slot"
import { cva, type VariantProps } from "class-variance-authority"
import { cn } from "../lib/utils"

const buttonVariants = cva(
  "inline-flex items-center justify-center whitespace-nowrap rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50",
  {
    variants: {
      variant: {
        default: "bg-primary text-primary-foreground shadow hover:bg-primary/90",
        destructive: "bg-destructive text-destructive-foreground shadow-sm hover:bg-destructive/90",
        outline: "border border-input bg-background shadow-sm hover:bg-accent hover:text-accent-foreground",
        secondary: "bg-secondary text-secondary-foreground shadow-sm hover:bg-secondary/80",
        ghost: "hover:bg-accent hover:text-accent-foreground",
        link: "text-primary underline-offset-4 hover:underline",
      },
      size: {
        default: "h-9 px-4 py-2",
        sm: "h-8 rounded-md px-3 text-xs",
        lg: "h-10 rounded-md px-8",
        icon: "h-9 w-9",
      },
    },
    defaultVariants: {
      variant: "default",
      size: "default",
    },
  }
)

export interface ButtonProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement>,
    VariantProps<typeof buttonVariants> {
  asChild?: boolean
}

const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant, size, asChild = false, ...props }, ref) => {
    const Comp = asChild ? Slot : "button"
    return (
      <Comp
        className={cn(buttonVariants({ variant, size, className }))}
        ref={ref}
        {...props}
      />
    )
  }
)
Button.displayName = "Button"

export { Button, buttonVariants }
```

## Extension Development Workflow

### 1. Development Commands

```bash
# Start development with watch mode
pnpm dev

# Build for production
pnpm build

# Build for development
pnpm build:dev

# Type checking
pnpm type-check

# Linting
pnpm lint

# Clean build artifacts
pnpm clean
```

### 2. Loading Extension in Chrome

1. Build the extension: `pnpm build:dev`
2. Open Chrome and navigate to `chrome://extensions/`
3. Enable "Developer mode" in the top right
4. Click "Load unpacked" and select the `dist/` folder
5. The extension will appear in your extensions list

### 3. Debugging

- **Background Script**: Go to `chrome://extensions/` � Click "Service Worker" link
- **Content Script**: Open browser DevTools � Console tab
- **Popup**: Right-click extension icon � "Inspect popup"
- **Options Page**: Right-click extension icon � "Options" � Open DevTools

## Key Features & Patterns

### 1. Component Architecture

```
src/
    background/              # Background service worker
    content/                # Content scripts
       components/         # React components for injection
       features/          # Feature-specific services
    utils/             # Content script utilities
    popup/                 # Popup interface
       components/        # Popup React components
    hooks/            # Custom React hooks
    options/              # Options page
 shared/               # Shared utilities and constants
```

### 2. Message Passing

```typescript
// Background to content script
chrome.tabs.sendMessage(tabId, {
  action: 'TOGGLE_TRANSLATION',
  enabled: true,
});

// Content script to background
chrome.runtime.sendMessage({
  action: 'OPEN_TAB',
  url: authUrl,
});

// Popup to content script
const [tab] = await chrome.tabs.query({ active: true, currentWindow: true });
chrome.tabs.sendMessage(tab.id!, { action: 'GET_PAGE_INFO' });
```

### 3. React Component Injection

```typescript
// utils/react-renderer.tsx
import { createRoot, Root } from 'react-dom/client';

const rootMap = new Map<string, Root>();

export async function renderReactComponent(
  component: React.ReactElement,
  containerId: string,
  containerSelector?: string
): Promise<void> {
  let container = document.getElementById(containerId);

  if (!container) {
    container = document.createElement('div');
    container.id = containerId;

    const parent = containerSelector ? document.querySelector(containerSelector) : document.body;

    parent?.appendChild(container);
  }

  if (!rootMap.has(containerId)) {
    const root = createRoot(container);
    rootMap.set(containerId, root);
  }

  const root = rootMap.get(containerId);
  root?.render(component);
}
```

### 4. Framework7 React Integration

```typescript
// popup/popup.tsx
import { App, View, Page, Navbar } from 'framework7-react';
import 'framework7/css/bundle';

const f7params = {
  name: 'English Learning Assistant',
  theme: 'auto',
};

export function PopupApp() {
  return (
    <App {...f7params}>
      <View main>
        <Page>
          <Navbar title="English Learning" />
          {/* Your popup content */}
        </Page>
      </View>
    </App>
  );
}
```

## Best Practices

### 1. File Naming Conventions

- Use kebab-case for all files and directories
- Component files: `component-name.tsx`
- Service files: `service-name-service.ts`
- Utility files: `category-utils.ts`

### 2. TypeScript Patterns

- Prefer `interface` over `type` for object shapes
- Use `as const` instead of `enum`
- Enable strict mode and proper type checking

### 3. Styling Guidelines

- Use Tailwind CSS classes with the `cn()` utility
- Avoid inline styles unless absolutely necessary
- Leverage CSS-in-JS only for dynamic styles

### 4. Performance Considerations

- Disable CSS code splitting for content scripts
- Use ES2020 target for modern browser support
- Optimize chunk sizes for extension constraints
- Lazy load components when possible

### 5. Security

- Follow Content Security Policy guidelines
- Validate all external data
- Use secure communication patterns
- Minimize permissions in manifest

## Troubleshooting

### Common Issues

1. **Content Script Not Loading**: Check manifest permissions and host_permissions
2. **CSS Not Applied**: Ensure web_accessible_resources includes CSS files
3. **React Components Not Rendering**: Verify content script injection timing
4. **Build Errors**: Check Vite configuration and dependency versions
5. **HMR Not Working**: Content scripts require extension reload for changes

### Development Tips

- Use `console.log` with prefixes to identify which context is logging
- Test in incognito mode to avoid cache issues
- Use Chrome DevTools for debugging each extension context
- Monitor the Extensions page for errors and warnings

## Migration to Shadcn/UI

If you want to add more Shadcn/UI components to the existing setup:

1. **Install Additional Radix Components**: Add required Radix UI primitives to `@english/ui`
2. **Create Component Files**: Follow the Shadcn/UI component patterns in
   `packages/ui/src/components/`
3. **Update Tailwind Config**: Add any required theme extensions
4. **Export Components**: Update `packages/ui/src/index.ts` to export new components
5. **Use in Extension**: Import and use components in popup, options, or content scripts

Example component addition:

```bash
# Add to packages/ui/package.json dependencies
"@radix-ui/react-dialog": "^1.0.5"

# Create packages/ui/src/components/dialog.tsx
# Export from packages/ui/src/components/index.ts
# Use in extension code
```

This setup provides a modern, maintainable Chrome extension development environment with full React,
TypeScript, and Tailwind CSS support while following Chrome Extension best practices and security
guidelines.
