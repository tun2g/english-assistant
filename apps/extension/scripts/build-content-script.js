import { rollup } from 'rollup';
import { nodeResolve } from '@rollup/plugin-node-resolve';
import commonjs from '@rollup/plugin-commonjs';
import typescript from '@rollup/plugin-typescript';
import alias from '@rollup/plugin-alias';
import { fileURLToPath } from 'url';
import { resolve, dirname } from 'path';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

// Build content script as a single IIFE bundle without imports
async function buildContentScript() {
  try {
    const bundle = await rollup({
      input: resolve(__dirname, '../src/content/content-script.ts'),
      plugins: [
        typescript({
          tsconfig: resolve(__dirname, '../tsconfig.json'),
          compilerOptions: {
            declaration: false,
            declarationMap: false,
            target: 'es2020',
            module: 'esnext',
            moduleResolution: 'node',
            allowSyntheticDefaultImports: true,
            esModuleInterop: true,
            skipLibCheck: true,
            isolatedModules: false,
            paths: {
              '@english/shared/*': [resolve(__dirname, '../../../packages/shared/src/*')],
            },
          },
          include: ['**/*.ts', '**/*.tsx'],
          exclude: ['node_modules/**'],
        }),
        alias({
          entries: [
            {
              find: '@english/shared',
              replacement: resolve(__dirname, '../../../packages/shared/src'),
            },
          ],
        }),
        nodeResolve({
          browser: true,
          preferBuiltins: false,
          extensions: ['.ts', '.js', '.json'],
          exportConditions: ['browser', 'default', 'module', 'main'],
          mainFields: ['browser', 'module', 'main'],
        }),
        commonjs({
          include: ['node_modules/**'],
        }),
      ],
      external: ['chrome'], // Only externalize chrome API
    });

    await bundle.write({
      file: resolve(__dirname, '../dist/content.js'),
      format: 'iife',
      name: 'EnglishLearningContentScript',
      globals: {
        chrome: 'chrome',
      },
    });

    console.log('Content script built successfully as IIFE bundle');
  } catch (error) {
    console.error('Failed to build content script:', error);
    process.exit(1);
  }
}

buildContentScript();