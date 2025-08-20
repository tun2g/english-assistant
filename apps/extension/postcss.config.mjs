import autoprefixer from 'autoprefixer';
import tailwindcss from 'tailwindcss';
import postcssRemToPx from '@thedutchcoder/postcss-rem-to-px';

export default {
  plugins: [
    tailwindcss(),
    postcssRemToPx({
      // Convert rem to px for better extension compatibility
      baseValue: 16, // 1rem = 16px
      unitPrecision: 5, // Precision for converted values
      propList: ['*'], // Convert all properties
      selectorBlackList: [], // No selectors to ignore
      replace: true, // Replace rem values instead of adding fallbacks
      mediaQuery: false, // Don't convert rem in media queries
      minRemValue: 0, // Convert all rem values
    }),
    autoprefixer(),
  ],
};