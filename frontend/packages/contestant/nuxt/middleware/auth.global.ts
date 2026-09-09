export default defineNuxtRouteMiddleware(async (to) => {
  const path = to.path.replace(/\/$/, "") || "/";
  const { viewer, refresh } = useSession();
  if (!viewer.value) await refresh();
  if (path === "/signin/impersonation") return;
  if (viewer.value?.state === "ANONYMOUS" && path !== "/signin")
    return navigateTo("/signin");
  if (viewer.value?.state === "DISCORD_AUTHENTICATED" && path !== "/signup")
    return navigateTo("/signup");
  if (
    viewer.value?.state === "CONTESTANT" &&
    ["/signin", "/signup"].includes(path)
  )
    return navigateTo("/problems");
});
