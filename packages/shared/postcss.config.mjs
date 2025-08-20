import autoprefixer from 'autoprefixer';
import tailwindcss from 'tailwindcss';

export default {
  plugins: [
    tailwindcss(),
    autoprefixer({
      // Target modern browsers that support Chrome extensions
      overrideBrowserslist: [
        'chrome >= 88',  // Minimum Chrome version for Manifest V3
        'edge >= 88',    // Minimum Edge version for Manifest V3
      ],
    }),
  ],
};