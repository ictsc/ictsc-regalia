/**
 * @type{import("prettier").Config}
 */
const config = {
  plugins: [import.meta.resolve("prettier-plugin-tailwindcss")],
  tailwindFunctions: ["clsx", "cva"],
};

export default config;
