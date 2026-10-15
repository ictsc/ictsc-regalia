declare global {
  interface Window {
    __ICTSC_RUNTIME_CONFIG__?: {
      demoMode?: boolean;
    };
  }
}

const runtimeDemoMode = () => {
  const injected = import.meta.client
    ? window.__ICTSC_RUNTIME_CONFIG__?.demoMode
    : undefined;
  return injected ?? useRuntimeConfig().public.demoMode;
};

export function useDemoMode() {
  return useState<boolean>("demo-mode", runtimeDemoMode);
}

export function loadDemoClock() {
  const now = Date.now();
  if (typeof window === "undefined") return now;
  try {
    const key = "ictsc-demo-now";
    const saved = localStorage.getItem(key);
    const parsed = saved == null ? Number.NaN : Number(saved);
    if (Number.isFinite(parsed) && parsed > 0) return parsed;
    localStorage.setItem(key, String(now));
  } catch {
    // When storage is unavailable, keep the clock fixed until page reload.
  }
  return now;
}

export function useClock(freezeInDemo = true) {
  const now = ref(Date.now());
  const configuredDemoMode = useDemoMode();
  const demoMode = computed(() => freezeInDemo && configuredDemoMode.value);
  const demoNow = useState<number | null>("demo-now", () =>
    runtimeDemoMode() ? loadDemoClock() : null,
  );
  let timer: ReturnType<typeof setInterval> | undefined;

  const updateTimer = () => {
    if (demoMode.value) {
      clearInterval(timer);
      timer = undefined;
    } else if (!timer) {
      now.value = Date.now();
      timer = setInterval(() => {
        now.value = Date.now();
      }, 1000);
    }
  };

  onMounted(() => {
    configuredDemoMode.value = runtimeDemoMode();
    if (configuredDemoMode.value) demoNow.value ??= loadDemoClock();
    updateTimer();
  });
  watch(demoMode, updateTimer);
  onUnmounted(() => {
    clearInterval(timer);
  });
  return computed(() =>
    demoMode.value && demoNow.value != null ? demoNow.value : now.value,
  );
}
