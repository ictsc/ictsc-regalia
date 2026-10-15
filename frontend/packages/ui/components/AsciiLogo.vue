<script setup lang="ts">
import logo from "../assets/regalia-logo.txt?raw";
</script>

<template>
  <span class="ascii-logo-frame" aria-hidden="true">
    <pre class="ascii-logo" v-text="logo" />
  </span>
</template>
