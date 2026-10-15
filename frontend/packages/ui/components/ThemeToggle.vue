<script setup lang="ts">
import { computed, onMounted } from "vue";
import { initializeTheme, useTheme } from "../theme";

const { preference, theme, toggle } = useTheme();
const label = computed(() =>
  preference.value === "system" ? "自動" : theme.value === "dark" ? "暗" : "明",
);
const actionLabel = computed(() =>
  preference.value === "system"
    ? `現在はOS設定の${theme.value === "dark" ? "ダーク" : "ライト"}テーマです。${theme.value === "dark" ? "ライト" : "ダーク"}テーマに切り替える`
    : `現在は${theme.value === "dark" ? "ダーク" : "ライト"}テーマです。OS設定に戻す`,
);

onMounted(initializeTheme);
</script>

<template>
  <button
    class="theme-toggle"
    type="button"
    :aria-label="actionLabel"
    @click="toggle"
  >
    <span aria-hidden="true">◐ {{ label }}</span>
  </button>
</template>
