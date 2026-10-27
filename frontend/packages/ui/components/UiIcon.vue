<script setup lang="ts">
defineProps<{
  name: "file-plus" | "close" | "download" | "pencil" | "eye" | "arrow-right";
}>();
</script>
<template>
  <svg
    class="ui-icon"
    aria-hidden="true"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    stroke-width="1.8"
    stroke-linecap="round"
    stroke-linejoin="round"
    xmlns="http://www.w3.org/2000/svg"
  >
    <template v-if="name === 'file-plus'">
      <path
        d="M14 3H6.5A1.5 1.5 0 0 0 5 4.5v15A1.5 1.5 0 0 0 6.5 21h11a1.5 1.5 0 0 0 1.5-1.5V8z"
      />
      <path d="M14 3v5h5M12 11v6M9 14h6" />
    </template>
    <path v-else-if="name === 'close'" d="M6 6l12 12M18 6L6 18" />
    <path v-else-if="name === 'arrow-right'" d="M5 12h14M13 6l6 6-6 6" />
    <template v-else-if="name === 'pencil'">
      <path d="M15.5 4.5l4 4L9 19H5v-4z" />
      <path d="M13 7l4 4" />
    </template>
    <template v-else-if="name === 'eye'">
      <path
        d="M2.5 12S6 5.5 12 5.5 21.5 12 21.5 12 18 18.5 12 18.5 2.5 12 2.5 12z"
      />
      <circle cx="12" cy="12" r="3" />
    </template>
    <template v-else>
      <path d="M12 4v11M7.5 10.5 12 15l4.5-4.5" />
      <path d="M5 19h14" />
    </template>
  </svg>
</template>
<style scoped>
.ui-icon {
  width: 18px;
  height: 18px;
  flex: none;
}
</style>
