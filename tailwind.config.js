/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './src/entry.js',
    './src/**/*.{js,jsx,ts,tsx}',
    './internal/components/**/*.go',
  ],
  theme: {
    extend: {
      screens: {
        'md': '768px',
      },
    },
  },
  plugins: [
    function ({ addBase, theme }) {
      addBase({
        ':root': {
          '--breakpoint-md': theme('screens.md'),
        },
      });
    },
  ],
}

