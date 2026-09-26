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
    <CompetitionHeader v-if="viewer?.state === 'CONTESTANT'" @logout="logout" />
    <header v-else class="admin-header">
      <NuxtLink class="brand" to="/" aria-label="ICTSC 2026 REGALIA ホーム">
        <AsciiLogo />
      </NuxtLink>
      <ThemeToggle />
    </header>
    <RequestState v-if="error" :error="error" @retry="logout" /><slot />
    <footer class="site-footer">
      <span>ICTSC 2026 / REGALIA</span><span>競技システム</span>
    </footer>
  </div>
</template>
