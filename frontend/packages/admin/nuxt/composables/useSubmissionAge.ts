import { onMounted, onUnmounted, ref } from "vue";

export function useSubmissionAge() {
  const now = ref(Date.now());
  let timer: ReturnType<typeof setInterval> | undefined;
  const updateNow = () => {
    now.value = Date.now();
  };

  onMounted(() => {
    updateNow();
    timer = setInterval(updateNow, 15_000);
    window.addEventListener("focus", updateNow);
  });
  onUnmounted(() => {
    clearInterval(timer);
    window.removeEventListener("focus", updateNow);
  });

  function elapsedMinutes(submittedAt: string) {
    return Math.max(
      0,
      Math.floor((now.value - Date.parse(submittedAt)) / 60_000),
    );
  }
  return { now, elapsedMinutes };
}
