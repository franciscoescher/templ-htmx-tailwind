import autoprefixer from 'autoprefixer';
import tailwindcss from 'tailwindcss';
import { defineConfig } from 'vite';

export default defineConfig({
  server: {
    open: true,
  },

  css: {
    postcss: {
      plugins: [
        tailwindcss(),
        autoprefixer(),
      ],
    },
  },

  build: {
    rollupOptions: {
      input: {
        entry: './src/entry.js',
      },
      output: {
        dir: './src/dist/',
        entryFileNames: '[name].js',
        assetFileNames: '[name].[ext]',
        chunkFileNames: '[name].[hash].js',
      },
    },
  },
});
