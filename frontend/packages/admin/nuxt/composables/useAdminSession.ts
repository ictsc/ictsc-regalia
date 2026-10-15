import { api, expectData, expectNoContent, type components } from "@ictsc/api";
export function useAdminSession() {
  const viewer = useState<components["schemas"]["AdminViewer"] | null>(
    "admin-viewer",
    () => null,
  );
  async function refresh() {
    viewer.value = expectData(await api.GET("/api/v1/admin/viewer")).viewer;
  }
  async function signOut() {
    expectNoContent(await api.POST("/api/v1/admin/auth/signout"));
    viewer.value = null;
    clearNuxtData();
    await refresh();
    await navigateTo("/");
  }
  return { viewer, refresh, signOut };
}
