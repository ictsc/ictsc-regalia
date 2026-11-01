<script lang="ts">
export type AdminFact = {
  label: string;
  value: string | number;
  /** Commit SHAs, IDs and other values read character by character. */
  mono?: boolean;
};
</script>
<script setup lang="ts">
/** Label and value pairs, such as a problem's settings or a commit. */
defineProps<{ items: AdminFact[] }>();
</script>
<template>
  <dl class="admin-facts">
    <div v-for="item in items" :key="item.label">
      <dt>{{ item.label }}</dt>
      <dd :class="{ 'admin-code': item.mono }">{{ item.value }}</dd>
    </div>
  </dl>
</template>
