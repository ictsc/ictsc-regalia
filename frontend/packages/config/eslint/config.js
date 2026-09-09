import tseslint from "typescript-eslint";
export default function config() { return tseslint.config(...tseslint.configs.recommended); }
