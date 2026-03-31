/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./App.{js,jsx,ts,tsx}",
    "./src/**/*.{js,jsx,ts,tsx}"
  ],
  theme: {
    extend: {
      colors: {
        // Primary colors (exact from web frontend)
        primary: {
          50: '#f0f9ff',
          100: '#e0f2fe',
          200: '#bae6fd',
          300: '#7dd3fc',
          400: '#38bdf8',
          500: '#0ea5e9',
          600: '#0284c7',
          700: '#0369a1',
          800: '#075985',
          900: '#0c4a6e',
          DEFAULT: '#0ea5e9',
        },
        // Secondary colors
        secondary: {
          50: '#f0fdfa',
          100: '#ccfbf1',
          200: '#99f6e4',
          300: '#5eead4',
          400: '#2dd4bf',
          500: '#14b8a6',
          600: '#0d9488',
          700: '#0f766e',
          800: '#115e59',
          900: '#134e4a',
          DEFAULT: '#14b8a6',
        },
        // Accent colors
        accent: {
          100: '#ffedd5',
          200: '#fed7aa',
          DEFAULT: '#f97316',
        },
        // Semantic colors
        success: '#22c55e',
        warning: '#eab308',
        destructive: '#dc2626',
        // Base colors
        background: '#ffffff',
        foreground: '#1e293b',
        card: '#ffffff',
        'card-foreground': '#1e293b',
        popover: '#ffffff',
        'popover-foreground': '#1e293b',
        muted: '#f1f5f9',
        'muted-foreground': '#334155',
        border: '#cbd5e1',
        input: '#cbd5e1',
        ring: '#0ea5e9',
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'],
      },
      borderRadius: {
        DEFAULT: '0.5rem',
      },
    },
  },
  plugins: [],
}
