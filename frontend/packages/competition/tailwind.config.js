/** @type {import('tailwindcss').Config} */
export default {
  content: ["./index.html", "./app/**/*.{js,jsx,ts,tsx}"],
  theme: {
    colors: {
      transparent: "transparent",
      current: "currentColor",
      primary: "#ff4040",
      surface: {
        "0": "#ffffff",
        "1": "#fdf1f1",
        "2": "#f8dcdc",
      },
      text: "#505050"
    },
    spacing: {
      full: "100%",
      screen: "100vw",
      auto: "auto",
      fit: "fit-content",
      0: "0",
    },
  },
  plugins: [],
};

