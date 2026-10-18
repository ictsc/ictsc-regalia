<script setup lang="ts">
defineProps<{
  /** Router path without the /admin base, e.g. "/teams/12". */
  currentPath: string;
}>();

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
  <nav class="admin-nav" aria-label="運営メニュー">
    <p class="admin-nav-heading">運営メニュー</p>
    <NuxtLink
      v-for="[to, label] in links"
      :key="to"
      :to="to!"
      :class="{
        'is-current':
          currentPath === to ||
          (to !== '/' && currentPath.startsWith(`${to}/`)),
      }"
      :aria-current="
        currentPath === to || (to !== '/' && currentPath.startsWith(`${to}/`))
          ? 'page'
          : undefined
      "
      >{{ label }}</NuxtLink
    >
  </nav>
</template>
