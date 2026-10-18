import { onMounted, onUnmounted, ref } from "vue";
import { minutesSince } from "~/features/grading";

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
    return minutesSince(submittedAt, now.value);
  }
  return { now, elapsedMinutes };
}
