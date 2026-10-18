<script setup lang="ts">
import ProfileView, {
  type ProfileForm,
} from "~/components/profile/ProfileView.vue";
import { api } from "@ictsc/api";
import { fetchProfile, updateProfile } from "~/features/profile";
useHead({ title: "プロフィール" });
const { data, error, pending, refresh } = await useAsyncData("profile", () =>
  fetchProfile(api),
);
const saving = ref(false),
  message = ref("");
const { signOut } = useSession();
const signingOut = ref(false);
const signOutError = ref("");
async function logout() {
  signingOut.value = true;
  signOutError.value = "";
  try {
    await signOut();
  } catch (e) {
    signOutError.value = e instanceof Error ? e.message : String(e);
  } finally {
    signingOut.value = false;
  }
}
async function save(form: ProfileForm) {
  saving.value = true;
  try {
    await updateProfile(api, form);
    await useSession().refresh();
    message.value = "プロフィールを更新しました";
  } catch (e) {
    message.value = e instanceof Error ? e.message : String(e);
  } finally {
    saving.value = false;
  }
}
</script>
<template>
  <main class="workspace">
    <ProfileView
      :profile="data ?? undefined"
      :pending="pending"
      :error="error"
      :saving="saving"
      :message="message"
      :signing-out="signingOut"
      :sign-out-error="signOutError"
      @save="save"
      @logout="logout"
      @retry="refresh"
    />
  </main>
</template>
