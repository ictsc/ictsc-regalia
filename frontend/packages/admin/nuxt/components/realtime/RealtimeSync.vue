<script setup lang="ts">
const { failing } = useRealtimeSync();
</script>
<template>
  <p v-if="failing" role="status" class="realtime-status">
    更新の受信に再接続しています。
  </p>
</template>
