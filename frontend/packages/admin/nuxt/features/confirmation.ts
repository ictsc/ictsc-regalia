import { useState } from "#imports";

export type AdminConfirmation = {
  title: string;
  message: string;
  confirmLabel?: string;
  destructive?: boolean;
};

let finishPending: ((accepted: boolean) => void) | null = null;

export function useAdminConfirmation() {
  const current = useState<AdminConfirmation | null>(
    "admin-confirmation",
    () => null,
  );
  function ask(request: AdminConfirmation): Promise<boolean> {
    if (current.value) return Promise.resolve(false);
    return new Promise((resolve) => {
      finishPending = resolve;
      current.value = request;
    });
  }
  function settle(accepted: boolean) {
    const resolve = finishPending;
    if (!resolve) return;
    finishPending = null;
    current.value = null;
    resolve(accepted);
  }
  return { current, ask, settle };
}
