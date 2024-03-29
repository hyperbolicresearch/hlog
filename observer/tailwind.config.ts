import type { Config } from 'tailwindcss'

export default {
  content: ['./app/**/*.{js,jsx,ts,tsx}'],
  theme: {
    extend: {
      maxWidth: {
        '200': '1600px',
      }
    },
  },
  plugins: [],
} satisfies Config

