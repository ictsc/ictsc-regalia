<script setup lang="ts">
import { reactive } from "vue";
import RequestState from "@ictsc/ui/components/RequestState.vue";
import type { ContestantProfile } from "~/features/models";

export type ProfileForm = Pick<
  ContestantProfile,
  "displayName" | "selfIntroduction"
>;

const props = defineProps<{
  profile?: ContestantProfile;
  pending?: boolean;
  error?: unknown;
  saving: boolean;
  /** Result of the last save, shown as a status message. */
  message: string;
  signingOut: boolean;
  signOutError: string;
}>();
defineEmits<{ save: [form: ProfileForm]; logout: []; retry: [] }>();

// Initialized once like the page did; edits are kept across re-renders.
const form = reactive<ProfileForm>({
  displayName: props.profile?.displayName ?? "",
  selfIntroduction: props.profile?.selfIntroduction ?? "",
});
</script>
<template>
  <h1 class="page-title">プロフィール</h1>
  <RequestState :error="error" :pending="pending" @retry="$emit('retry')"
    ><form class="form-grid" @submit.prevent="$emit('save', { ...form })">
      <label>競技者名<input :value="profile?.name" disabled /></label
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
      @click="$emit('logout')"
    >
      {{ signingOut ? "処理中…" : "ログアウト" }}
    </button>
  </section>
</template>
