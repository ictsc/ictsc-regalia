<script setup lang="ts">
import AdminHeader from "~/components/layout/AdminHeader.vue";
import AdminNav from "~/components/layout/AdminNav.vue";
const { viewer, refresh, signOut } = useAdminSession();
const route = useRoute();
const error = ref<unknown>();
async function load() {
  try {
    await refresh();
    error.value = undefined;
  } catch (e) {
    error.value = e;
  }
}
await load();
async function logout() {
  try {
    await signOut();
  } catch (e) {
    error.value = e;
  }
}
</script>
<template>
  <div class="admin-app">
    <AdminHeader
      :signed-in="viewer?.state !== 'ANONYMOUS'"
      :display-name="
        viewer?.state === 'ADMIN'
          ? viewer.admin.discord.display_name
          : undefined
      "
      @logout="logout"
    />
    <RequestState :error="error" @retry="load"
      ><main v-if="viewer?.state === 'ANONYMOUS'" class="workspace">
        <h1>ICTSC Admin</h1>
        <p>Discordの運営ロールでログインしてください。</p>
        <a
          class="button-primary"
          href="/api/v1/admin/auth/discord?next=%2Fadmin%2F"
          >Discordでログイン</a
        >
      </main>
      <template v-else-if="viewer"
        ><div class="admin-shell">
          <AdminNav :current-path="route.path" />
          <slot /></div></template
    ></RequestState>
  </div>
</template>
