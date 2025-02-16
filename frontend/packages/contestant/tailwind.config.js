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
      icon: "#1c1b1f",
    },
    width: {
      full: "100%",
      screen: "100vw",
      auto: "auto",
      fit: "fit-content",
      0: "0",
    },
    spacing: {
      0: "0",
      4: "4px",
      8: "8px",
      12: "12px",
      16: "16px",
      20: "20px",
      24: "24px",
      40: "40px",
      64: "64px",
    },
    fontSize: {
      12: ["12px", "20px"],
      14: ["14px", "22px"],
      16: ["16px", "24px"],
      24: ["24px", "32px"],
      32: ["32px", "40px"],
      48: ["48px", "1"],
    },
    borderRadius: {
      none: "0",
      4: "4px",
      8: "8px",
      12: "12px",
      16: "16px",
      full: "9999px",
    },
  },
  plugins: [],
};
