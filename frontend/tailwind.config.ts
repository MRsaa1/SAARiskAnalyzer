import type { Config } from 'tailwindcss'

export default {
  darkMode: 'class',
  content: ['./index.html', './src/**/*.{ts,tsx}'],
  theme: {
    extend: {
      colors: {
        bg: 'var(--bg)',
        fg: 'var(--fg)',
        muted: 'var(--muted)',
        card: 'var(--card)',
        border: 'var(--border)',
        accent: 'var(--accent)',
        link: 'var(--link)',
      },
      borderRadius: {
        xl: 'var(--radius)',
      },
      boxShadow: {
        brand: 'var(--shadow)',
      },
      fontFamily: {
        sans: ['var(--font-sans)'],
      },
    },
  },
  plugins: [],
} satisfies Config
