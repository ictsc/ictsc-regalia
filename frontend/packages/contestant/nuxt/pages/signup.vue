<script setup lang="ts">
import SignupView, {
  type SignupForm,
} from "~/components/signup/SignupView.vue";
import { signUp } from "~/features/viewer/signup";
useHead({ title: "参加登録" });
const { viewer } = useSession();
const registrationTeam = computed(() =>
  viewer.value?.state === "DISCORD_AUTHENTICATED"
    ? viewer.value.registration_team
    : undefined,
);
const message = ref(""),
  pending = ref(false);
async function submit(form: SignupForm) {
  pending.value = true;
  try {
    const result = await signUp({
      ...form,
      invitationCode: registrationTeam.value ? undefined : form.invitationCode,
    });
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
    <SignupView
      :registration-team-name="registrationTeam?.name"
      :initial-name="
        viewer?.state === 'DISCORD_AUTHENTICATED' ? viewer.discord.username : ''
      "
      :initial-display-name="
        viewer?.state === 'DISCORD_AUTHENTICATED'
          ? viewer.discord.display_name
          : ''
      "
      :pending="pending"
      :message="message"
      @submit="submit"
    />
  </main>
</template>
