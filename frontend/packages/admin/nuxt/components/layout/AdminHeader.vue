<script setup lang="ts">
import ThemeToggle from "@ictsc/ui/components/ThemeToggle.vue";

defineProps<{
  /** Not anonymous: show the account name and the logout button. */
  signedIn: boolean;
  displayName?: string;
}>();
defineEmits<{ logout: [] }>();
</script>
<template>
  <header class="admin-header">
    <NuxtLink class="brand" to="/"
      ><span class="brand-name">ICTSC</span
      ><span class="brand-edition">REGALIA<br />ADMIN</span></NuxtLink
    >
    <div class="actions">
      <ThemeToggle />
      <template v-if="signedIn">
        <span>{{ displayName ?? "" }}</span
        ><button class="button-secondary" @click="$emit('logout')">
          ログアウト
        </button>
      </template>
    </div>
  </header>
</template>
