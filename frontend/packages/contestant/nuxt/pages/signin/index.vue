<script setup lang="ts">
useHead({ title: "ログイン" });
</script>
<template>
  <main class="workspace">
    <h1 class="page-title">ログイン</h1>
    <p>競技に参加するDiscordアカウントでログインしてください。</p>
    <div class="actions">
      <a class="button-primary" href="/api/v1/auth/discord?next=%2F"
        >Discordでログイン</a
      >
    </div>
  </main>
</template>
