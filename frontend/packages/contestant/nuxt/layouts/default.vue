<script setup lang="ts">
import { teamStyle } from "@ictsc/ui/colors";
const { viewer } = useSession();
</script>
<template>
  <div
    :style="
      teamStyle(viewer?.state === 'CONTESTANT' ? viewer.team.color : undefined)
    "
  >
    <CompetitionHeader v-if="viewer?.state === 'CONTESTANT'" />
    <header v-else class="admin-header">
      <NuxtLink class="brand" to="/" aria-label="ICTSC 2026 REGALIA ホーム">
        <AsciiLogo />
      </NuxtLink>
      <ThemeToggle />
    </header>
    <slot />
    <footer class="site-footer">
      <span>ICTSC 2026 / REGALIA</span><span>競技システム</span>
    </footer>
  </div>
</template>
