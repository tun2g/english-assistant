import { crx } from '@crxjs/vite-plugin';
import react from '@vitejs/plugin-react';
import { resolve } from 'path';
import { defineConfig, loadEnv, UserConfig } from 'vite';
import { fileURLToPath, URL } from 'node:url';

import manifest from './src/manifest.config';

export default ({ mode = 'development' }: UserConfig) => {
  process.env = { ...process.env, ...loadEnv(mode, process.cwd()) };

  const isProduction = mode === 'production';

  return defineConfig({
    base: './',
    define: {
      'process.env.BUILD_TARGET': JSON.stringify('extension'),
      'process.env.NODE_ENV': JSON.stringify(mode),
      'process.env.INLINE_RUNTIME_CHUNK': JSON.stringify('false'),
      global: 'globalThis',
    },
    server: {
      port: process.env.VITE_PORT ? +process.env.VITE_PORT : 4444,
    },
    plugins: [
      react({
        // Enable automatic JSX runtime
        jsxRuntime: 'automatic',
      }),
      // CRXJS plugin handles manifest and Chrome extension specific bundling
      crx({
        manifest,
        // Enable content script HMR in development
        contentScripts: {
          injectCss: true,
        },
        // Ensure proper handling of popup and options pages
        browser: 'chrome',
      }),
    ],
    build: {
      outDir: 'dist',
      minify: isProduction,
      sourcemap: !isProduction,
      // Enable target for modern browsers (Chrome extensions support ES2020+)
      target: 'es2020',
      rollupOptions: {
        output: {
          // Use consistent naming for easier debugging and web_accessible_resources
          entryFileNames: '[name].js',
          chunkFileNames: '[name]-[hash].js',
          assetFileNames: '[name].[ext]',
        },
        // Disable code splitting for content scripts to avoid CSP issues
        external: id => {
          // Keep React bundled for content scripts to avoid dynamic imports
          return false;
        },
      },
      // Disable CSS code splitting to avoid web_accessible_resources issues
      cssCodeSplit: false,
      // Chrome extensions have different performance constraints
      chunkSizeWarningLimit: 1000,
    },
    resolve: {
      alias: {
        '@': resolve(fileURLToPath(new URL('.', import.meta.url)), 'src'),
        '@english/shared': resolve(fileURLToPath(new URL('.', import.meta.url)), '../../packages/shared/src'),
        '@english/ui': resolve(fileURLToPath(new URL('.', import.meta.url)), '../../packages/ui/src'),
      },
    },
    // Optimize dependencies for faster builds
    optimizeDeps: {
      include: ['react', 'react-dom', '@tanstack/react-query'],
      exclude: ['@english/shared', '@english/ui'],
    },
    // CSS configuration with PostCSS and SCSS support
    css: {
      // PostCSS will use postcss.config.mjs automatically
      preprocessorOptions: {
        scss: {
          // Enable modern API and suppress deprecation warnings
          api: 'modern-compiler',
          silenceDeprecations: ['legacy-js-api'],
        },
      },
    },
  });
};
