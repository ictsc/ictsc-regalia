import { api, expectData, expectNoContent, type components } from "@ictsc/api";
export function useSession() {
  const viewer = useState<components["schemas"]["Viewer"] | null>(
    "viewer",
    () => null,
  );
  async function refresh() {
    viewer.value = expectData(
      await api.GET("/api/v1/viewer", { cache: "no-store" }),
    ).viewer;
  }
  async function signOut() {
    expectNoContent(await api.POST("/api/v1/auth/signout"));
    viewer.value = null;
    clearNuxtData();
    await navigateTo("/signin");
  }
  return { viewer, refresh, signOut };
}
