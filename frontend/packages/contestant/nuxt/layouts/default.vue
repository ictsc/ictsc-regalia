<script setup lang="ts">
import { teamStyle } from "@ictsc/ui/colors";
const { viewer, signOut } = useSession();
const error = ref<unknown>();
async function logout() {
  try {
    await signOut();
  } catch (e) {
    error.value = e;
  }
}
</script>
<template>
  <div
    :style="
      teamStyle(viewer?.state === 'CONTESTANT' ? viewer.team.color : undefined)
    "
  >
    <CompetitionHeader v-if="viewer?.state === 'CONTESTANT'" />
    <header v-else class="admin-header">
      <NuxtLink class="brand" to="/"
        ><span class="brand-name">ICTSC</span
        ><span class="brand-edition">2026<br />REGALIA</span></NuxtLink
      >
    </header>
    <nav
      v-if="viewer?.state === 'CONTESTANT'"
      class="account-nav"
      aria-label="アカウント"
    >
      <span v-if="viewer.impersonated_by"
        >{{ viewer.impersonated_by }} による代理操作中</span
      ><span class="team-badge">{{ viewer.team.name }}</span
      ><NuxtLink to="/teams">チーム</NuxtLink
      ><NuxtLink to="/rule">ルール</NuxtLink
      ><NuxtLink to="/profile">プロフィール</NuxtLink
      ><button @click="logout">ログアウト</button>
    </nav>
    <RequestState v-if="error" :error="error" @retry="logout" /><slot />
    <footer class="site-footer">
      <span>ICTSC 2026 / REGALIA</span><span>競技システム</span>
    </footer>
  </div>
</template>
