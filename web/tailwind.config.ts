import type { Config } from 'tailwindcss';

export default {
  content: ['./src/**/*.{astro,html,js,jsx,md,mdx,svelte,ts,tsx,vue}'],
  theme: {
    extend: {
      colors: {
        azure: {
          DEFAULT: '#0078D4',
          dark: '#005A9E',
          light: '#50A0E6',
        },
      },
    },
  },
  plugins: [],
} satisfies Config;
