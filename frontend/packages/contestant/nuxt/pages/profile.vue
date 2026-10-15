<script setup lang="ts">
import { api } from "@ictsc/api";
import { fetchProfile, updateProfile } from "~/features/profile";
useHead({ title: "プロフィール" });
const { data, error, pending, refresh } = await useAsyncData("profile", () =>
  fetchProfile(api),
);
const form = reactive({
  displayName: data.value?.displayName ?? "",
  selfIntroduction: data.value?.selfIntroduction ?? "",
});
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
async function save() {
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
    <h1 class="page-title">プロフィール</h1>
    <RequestState :error="error" :pending="pending" @retry="refresh"
      ><form class="form-grid" @submit.prevent="save">
        <label>競技者名<input :value="data?.name" disabled /></label
        ><label
          >表示名<input
            v-model="form.displayName"
            required
            maxlength="255" /></label
        ><label
          >自己紹介<textarea v-model="form.selfIntroduction" maxlength="2000" />
        </label>
        <p role="status">{{ message }}</p>
        <button class="button-primary" :disabled="saving">更新する</button>
      </form></RequestState
    >
    <section class="profile-account" aria-labelledby="profile-account-heading">
      <h2 id="profile-account-heading">アカウント</h2>
      <p v-if="signOutError" role="alert" class="error-message">
        {{ signOutError }}
      </p>
      <button
        class="button-secondary profile-logout"
        :disabled="signingOut"
        @click="logout"
      >
        {{ signingOut ? "処理中…" : "ログアウト" }}
      </button>
    </section>
  </main>
</template>
