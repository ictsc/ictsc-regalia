/** @type {import('tailwindcss').Config} */
const {fontFamily} = require('tailwindcss/defaultTheme');

module.exports = {
  content: [
    "./pages/**/*.{js,ts,jsx,tsx}",
    "./components/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      fontFamily: {
        sans: ['var(--font-noto-sans-jp)', ...fontFamily.sans]
      }
    },
  },
  daisyui: {
    themes: [
      {
        ictsc: {
          "primary": "#f43f5e",
          // "primary-focus": "#f43f5e",
          "primary-content": "#FFFFFF",
          "accent": "#37CDBE",
          // "neutral": "#f43f5e",
          "base-100": "#FFFFFF",
          "info": "#3ABFF8",
          "success": "#36D399",
          "warning": "#FBBD23",
          "error": "#F87272",
        },
      },
    ],
  },
  plugins: [require("daisyui")],
}
