// Import UI config as the base
import uiConfig from '../ui/tailwind.config.js';

/** @type {import('tailwindcss').Config} */
export default {
  // Extend the UI package configuration
  ...uiConfig,
  content: [
    './src/**/*.{html,js,ts,jsx,tsx}',
    // Include UI package content
    '../ui/src/**/*.{js,ts,jsx,tsx}',
  ],
  theme: {
    ...uiConfig.theme,
    extend: {
      ...uiConfig.theme?.extend,
      // Add shared semantic colors that complement shadcn/ui
      colors: {
        ...uiConfig.theme?.extend?.colors,
        // Semantic colors using HSL like shadcn/ui
        success: {
          DEFAULT: "hsl(142.1 76.2% 36.3%)",
          foreground: "hsl(355.7 100% 97.3%)",
        },
        warning: {
          DEFAULT: "hsl(32.5 94.6% 43.7%)",
          foreground: "hsl(355.7 100% 97.3%)",
        },
        error: {
          DEFAULT: "hsl(0 84.2% 60.2%)",
          foreground: "hsl(210 40% 98%)",
        },
        info: {
          DEFAULT: "hsl(221.2 83.2% 53.3%)",
          foreground: "hsl(210 40% 98%)",
        },
      },
      // Additional font families
      fontFamily: {
        inter: [
          'Inter', 
          'system-ui', 
          '-apple-system', 
          'BlinkMacSystemFont', 
          '"Segoe UI"', 
          'Roboto', 
          'sans-serif'
        ],
        mono: [
          '"JetBrains Mono"', 
          'Consolas', 
          'Monaco', 
          '"Courier New"', 
          'monospace'
        ],
      },
      // Custom shadows for extension UI
      boxShadow: {
        'custom-sm': '0 1px 2px 0 rgba(0, 0, 0, 0.05)',
        'custom-md': '0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06)',
        'custom-lg': '0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05)',
        'custom-xl': '0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04)',
      },
    },
  },
  plugins: [...(uiConfig.plugins || [])],
};