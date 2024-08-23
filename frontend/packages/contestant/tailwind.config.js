/** @type {import('tailwindcss').Config} */
export default {
  content: ["./index.html", "./app/**/*.{js,jsx,ts,tsx}"],
  theme: {
    colors: {
      transparent: "transparent",
      current: "currentColor",
      primary: "#ff4040",
      surface: {
        0: "#ffffff",
        1: "#fdf1f1",
        2: "#f8dcdc",
      },
      text: "#505050",
    },
    width: {
      full: "100%",
      screen: "100vw",
      auto: "auto",
      fit: "fit-content",
      0: "0",
    },
    fontSize: {
      12: ["0.75rem", "1.25rem"],
      14: ["0.875rem", "1.375rem"],
      16: ["1rem", "1.5rem"],
      24: ["1.5rem", "2rem"],
    },
  },
  plugins: [],
};
