import uiConfig from '../../packages/ui/tailwind.config.js';

/** @type {import('tailwindcss').Config} */
export default {
  // Extend the UI package configuration
  ...uiConfig,
  content: [
    // Extension specific content
    './src/**/*.{js,ts,jsx,tsx,html}',
    // Include all shared packages that use Tailwind
    '../../packages/shared/src/**/*.{js,ts,jsx,tsx}',
    '../../packages/ui/src/**/*.{js,ts,jsx,tsx}',
  ],
  theme: {
    ...uiConfig.theme,
    extend: {
      ...uiConfig.theme?.extend,
      // Extension-specific theme extensions
      colors: {
        ...uiConfig.theme?.extend?.colors,
      },
      fontFamily: {
        ...uiConfig.theme?.extend?.fontFamily,
        inter: ['Inter', 'system-ui', '-apple-system', 'BlinkMacSystemFont', '"Segoe UI"', 'Roboto', 'sans-serif'],
      },
      zIndex: {
        ...uiConfig.theme?.extend?.zIndex,
        // Chrome extension specific z-index values
        'extension-base': '2147483640',
        'extension-overlay': '2147483645',
        'extension-modal': '2147483647',
        'extension-tooltip': '2147483648',
      },
      animation: {
        ...uiConfig.theme?.extend?.animation,
        // Extension-specific animations
        'fade-in': 'fadeIn 0.15s ease-out',
        'fade-out': 'fadeOut 0.15s ease-in',
        'slide-up': 'slideUp 0.2s ease-out',
        'scale-in': 'scaleIn 0.15s ease-out',
      },
      keyframes: {
        ...uiConfig.theme?.extend?.keyframes,
        // Extension-specific keyframes
        fadeIn: {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' },
        },
        fadeOut: {
          '0%': { opacity: '1' },
          '100%': { opacity: '0' },
        },
        slideUp: {
          '0%': { transform: 'translateY(10px)', opacity: '0' },
          '100%': { transform: 'translateY(0)', opacity: '1' },
        },
        scaleIn: {
          '0%': { transform: 'scale(0.95)', opacity: '0' },
          '100%': { transform: 'scale(1)', opacity: '1' },
        },
      },
    },
  },
  // Extension-specific safelist for dynamic classes
  safelist: [
    ...(uiConfig.safelist || []),
    // Extension-specific safelist
    'z-extension-base',
    'z-extension-overlay',
    'z-extension-modal',
    'z-extension-tooltip',
    'animate-fade-in',
    'animate-fade-out',
    'animate-slide-up',
    'animate-scale-in',
    // Dynamic z-index classes
    'z-[2147483647]',
    'z-[2147483648]',
  ],
};
