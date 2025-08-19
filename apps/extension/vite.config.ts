import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import { resolve } from 'path';
import { fileURLToPath, URL } from 'node:url';
import { copyFileSync } from 'fs';

export default defineConfig({
  plugins: [
    react(),
    {
      name: 'copy-manifest',
      writeBundle() {
        copyFileSync(
          resolve(fileURLToPath(new URL('.', import.meta.url)), 'src/manifest.json'),
          resolve(fileURLToPath(new URL('.', import.meta.url)), 'dist/manifest.json')
        );
      }
    }
  ],
  build: {
    outDir: 'dist',
    emptyOutDir: false, // Don't empty dir to preserve content.js built separately
    minify: false,
    rollupOptions: {
      input: {
        popup: resolve(fileURLToPath(new URL('.', import.meta.url)), 'src/popup/popup.html'),
        options: resolve(fileURLToPath(new URL('.', import.meta.url)), 'src/options/options.html'),
        background: resolve(fileURLToPath(new URL('.', import.meta.url)), 'src/background/background.ts'),
        // content script is built separately via rollup
      },
      output: {
        entryFileNames: '[name].js',
        chunkFileNames: '[name].js', 
        assetFileNames: '[name].[ext]',
        format: 'es',
        inlineDynamicImports: false,
        manualChunks: undefined,
      },
    },
  },
  resolve: {
    alias: {
      '@': resolve(fileURLToPath(new URL('.', import.meta.url)), 'src'),
      '@english/shared': resolve(fileURLToPath(new URL('.', import.meta.url)), '../../packages/shared/src'),
      '@english/ui': resolve(fileURLToPath(new URL('.', import.meta.url)), '../../packages/ui/src'),
    },
  },
  define: {
    global: 'globalThis',
  },
});