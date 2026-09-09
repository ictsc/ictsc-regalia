export function useMutation(refresh?: () => Promise<unknown>) {
  const busy = ref(false),
    message = ref(""),
    failed = ref(false);
  async function run(
    label: string,
    action: () => Promise<unknown>,
    confirmation?: string,
  ) {
    if (busy.value || (confirmation && !window.confirm(confirmation)))
      return false;
    busy.value = true;
    message.value = "";
    failed.value = false;
    try {
      await action();
      await refresh?.();
      message.value = `${label}しました`;
      return true;
    } catch (e) {
      failed.value = true;
      message.value = e instanceof Error ? e.message : String(e);
      return false;
    } finally {
      busy.value = false;
    }
  }
  return { busy, message, failed, run };
}
