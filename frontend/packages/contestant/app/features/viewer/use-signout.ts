import { useRouter } from "@tanstack/react-router";

export type SignOutAction = () => Promise<void>;

export function useSignOut(): SignOutAction {
  const router = useRouter();
  return async () => {
    await fetch("/api/auth/signout", {
      method: "POST",
    });
    await router.invalidate();
  };
}
