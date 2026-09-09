<script setup lang="ts">
const { viewer, refresh, signOut } = useAdminSession();
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
const links = [
  ["/", "概要"],
  ["/teams", "チーム"],
  ["/contestants", "参加者"],
  ["/content", "コンテンツ"],
  ["/submissions", "採点"],
  ["/scores", "得点・順位"],
  ["/deployments", "再展開"],
  ["/settings", "競技設定"],
];
</script>
<template>
  <div>
    <header class="admin-header">
      <NuxtLink class="brand" to="/"
        ><span class="brand-name">ICTSC</span
        ><span class="brand-edition">REGALIA<br />ADMIN</span></NuxtLink
      >
      <div v-if="viewer?.state !== 'ANONYMOUS'" class="actions">
        <span>{{
          viewer?.state === "ADMIN" ? viewer.admin.discord.display_name : ""
        }}</span
        ><button class="button-secondary" @click="logout">ログアウト</button>
      </div>
    </header>
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
        ><nav class="admin-nav" aria-label="運営メニュー">
          <NuxtLink v-for="[to, label] in links" :key="to" :to="to!">{{
            label
          }}</NuxtLink>
        </nav>
        <slot /></template
    ></RequestState>
  </div>
</template>
