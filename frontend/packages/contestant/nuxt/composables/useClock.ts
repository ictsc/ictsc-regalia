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

export function setDemoClock(referenceTime: string | number = Date.now()) {
  const demoMode = useDemoMode();
  const time =
    typeof referenceTime === "number"
      ? referenceTime
      : Date.parse(referenceTime);
  if (demoMode.value && Number.isFinite(time)) {
    useState<number | null>("demo-now", () => null).value = time;
  }
}

export function useClock() {
  const now = ref(Date.now());
  const demoMode = useDemoMode();
  const demoNow = useState<number | null>("demo-now", () =>
    runtimeDemoMode() ? Date.now() : null,
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
    demoMode.value = runtimeDemoMode();
    demoNow.value = demoMode.value ? (demoNow.value ?? Date.now()) : null;
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
