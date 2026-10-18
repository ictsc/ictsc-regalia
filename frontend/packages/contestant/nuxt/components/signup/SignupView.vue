<script setup lang="ts">
import { reactive } from "vue";

export type SignupForm = {
  invitationCode: string;
  name: string;
  displayName: string;
};

const props = defineProps<{
  /** Team assigned by the Discord role; the invitation code is not needed. */
  registrationTeamName?: string;
  initialName: string;
  initialDisplayName: string;
  pending: boolean;
  message: string;
}>();
defineEmits<{ submit: [form: SignupForm] }>();

const form = reactive<SignupForm>({
  invitationCode: "",
  name: props.initialName,
  displayName: props.initialDisplayName,
});
</script>
<template>
  <h1 class="page-title">参加登録</h1>
  <form class="form-grid" @submit.prevent="$emit('submit', { ...form })">
    <p v-if="registrationTeamName">参加チーム: {{ registrationTeamName }}</p>
    <label v-else
      >招待コード<input v-model="form.invitationCode" required /></label
    ><label>競技者名<input v-model="form.name" required /></label
    ><label
      >表示名<input v-model="form.displayName" required maxlength="255"
    /></label>
    <p v-if="message" role="alert">{{ message }}</p>
    <button class="button-primary" :disabled="pending">登録する</button>
  </form>
</template>
