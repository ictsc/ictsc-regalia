<script setup lang="ts">
/**
 * The top of every admin page, so pages line up: back link, title, lead,
 * page-level actions and the result of the last operation, then the panels.
 */
defineProps<{
  title: string;
  lead?: string;
  back?: { to: string; label: string };
  /** Result of the last operation. */
  message?: string;
}>();
</script>
<template>
  <header class="admin-page-header">
    <NuxtLink v-if="back" class="admin-back" :to="back.to"
      >← {{ back.label }}</NuxtLink
    >
    <h1>{{ title }}</h1>
    <p v-if="lead || $slots.lead" class="admin-lead">
      <slot name="lead">{{ lead }}</slot>
    </p>
    <div v-if="$slots.actions" class="admin-page-actions">
      <slot name="actions" />
    </div>
  </header>
  <p v-if="message" role="status" class="admin-status">{{ message }}</p>
  <slot />
</template>
