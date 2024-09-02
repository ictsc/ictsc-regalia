/**
 * @type{import("prettier").Config}
 */
const config = {
  plugins: ["prettier-plugin-tailwindcss"],
  tailwindFunctions: ["clsx", "cva"],
};

export default config;
