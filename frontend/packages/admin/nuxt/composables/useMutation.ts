import {
  useAdminConfirmation,
  type AdminConfirmation,
} from "~/features/confirmation";

export function useMutation(refresh?: () => Promise<unknown>) {
  const { ask: askConfirmation } = useAdminConfirmation();
  const busy = ref(false),
    message = ref(""),
    failed = ref(false);
  async function run(
    label: string,
    action: () => Promise<unknown>,
    confirmation?: string | Omit<AdminConfirmation, "title">,
  ) {
    if (busy.value) return false;
    if (confirmation) {
      const options =
        typeof confirmation === "string"
          ? { title: label, message: confirmation, confirmLabel: label }
          : { title: label, confirmLabel: label, ...confirmation };
      if (!(await askConfirmation(options))) return false;
    }
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
