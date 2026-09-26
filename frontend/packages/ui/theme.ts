import { computed, readonly, ref } from "vue";

export type Theme = "light" | "dark";
export type ThemePreference = "system" | Theme;

const storageKey = "ictsc/theme";
const preference = ref<ThemePreference>("system");
const systemDark = ref(false);
let initialized = false;

export function nextThemePreference(
  current: ThemePreference,
  prefersDark: boolean,
): ThemePreference {
  if (current === "system") return prefersDark ? "light" : "dark";
  return "system";
}

function applyTheme() {
  const theme =
    preference.value === "system"
      ? systemDark.value
        ? "dark"
        : "light"
      : preference.value;
  document.documentElement.dataset.theme = theme;
  document.documentElement.dataset.themePreference = preference.value;
}

function setPreference(value: ThemePreference) {
  preference.value = value;
  if (value === "system") localStorage.removeItem(storageKey);
  else localStorage.setItem(storageKey, value);
  applyTheme();
}

export function initializeTheme() {
  if (initialized || typeof window === "undefined") return;
  initialized = true;

  const media = window.matchMedia("(prefers-color-scheme: dark)");
  systemDark.value = media.matches;
  const saved = localStorage.getItem(storageKey);
  preference.value = saved === "light" || saved === "dark" ? saved : "system";
  applyTheme();

  media.addEventListener("change", (event) => {
    systemDark.value = event.matches;
    if (preference.value === "system") applyTheme();
  });
  window.addEventListener("storage", (event) => {
    if (event.key !== storageKey) return;
    preference.value =
      event.newValue === "light" || event.newValue === "dark"
        ? event.newValue
        : "system";
    applyTheme();
  });
}

export function useTheme() {
  const theme = computed<Theme>(() =>
    preference.value === "system"
      ? systemDark.value
        ? "dark"
        : "light"
      : preference.value,
  );
  const toggle = () =>
    setPreference(nextThemePreference(preference.value, systemDark.value));

  return {
    preference: readonly(preference),
    theme,
    toggle,
  };
}
