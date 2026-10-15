export default defineNuxtPlugin(() => {
  if (!("serviceWorker" in navigator)) return;
  navigator.serviceWorker.addEventListener("message", (event) => {
    if (event.data?.type === "announcement-push") {
      void refreshNuxtData("competition");
      void refreshNuxtData("notices");
    }
  });
});
