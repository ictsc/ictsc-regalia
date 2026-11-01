<script setup lang="ts">
/**
 * One section of an admin page: a top rule, a heading with an optional
 * count, badge or link on the right, then the content.
 */
defineProps<{
  title?: string;
  /** Irreversible operations: accent rule and heading. */
  danger?: boolean;
}>();
</script>
<template>
  <section class="admin-panel" :class="{ 'admin-danger': danger }">
    <div v-if="title" class="admin-panel-heading">
      <h2>{{ title }}</h2>
      <span v-if="$slots.aside" class="admin-panel-aside"
        ><slot name="aside"
      /></span>
    </div>
    <slot />
  </section>
</template>
