/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './src/entry.js',
    './src/**/*.{js,jsx,ts,tsx}',
    './internal/components/**/*.go',
  ],
  theme: {
    extend: {},
  },
  plugins: [],
}

