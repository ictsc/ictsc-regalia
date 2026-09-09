<script setup lang="ts">
import { signUp } from "~/features/viewer/signup";
useHead({ title: "参加登録" });
const { viewer } = useSession();
const form = reactive({
    invitationCode: "",
    name:
      viewer.value?.state === "DISCORD_AUTHENTICATED"
        ? viewer.value.discord.username
        : "",
    displayName:
      viewer.value?.state === "DISCORD_AUTHENTICATED"
        ? viewer.value.discord.display_name
        : "",
  }),
  message = ref(""),
  pending = ref(false);
async function submit() {
  pending.value = true;
  try {
    const result = await signUp(form);
    if (result.error) {
      message.value = `登録できませんでした: ${result.invitationCodeError ?? result.nameError ?? result.displayNameError ?? result.error}`;
      return;
    }
    await useSession().refresh();
    await navigateTo("/problems");
  } catch (e) {
    message.value = e instanceof Error ? e.message : String(e);
  } finally {
    pending.value = false;
  }
}
</script>
<template>
  <main class="workspace">
    <h1>参加登録</h1>
    <form class="form-grid" @submit.prevent="submit">
      <label>招待コード<input v-model="form.invitationCode" required /></label
      ><label>競技者名<input v-model="form.name" required /></label
      ><label
        >表示名<input v-model="form.displayName" required maxlength="255"
      /></label>
      <p v-if="message" role="alert">{{ message }}</p>
      <button class="button-primary" :disabled="pending">登録する</button>
    </form>
  </main>
</template>
