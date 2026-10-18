<script setup lang="ts">
import AdminConfirmationDialog from "~/components/confirmation/AdminConfirmationDialog.vue";
import { useAdminConfirmation } from "~/features/confirmation";

const { current: confirmation, settle: closeConfirmation } =
  useAdminConfirmation();

onBeforeUnmount(() => closeConfirmation(false));
</script>

<template>
  <NuxtLayout><NuxtPage /></NuxtLayout>
  <AdminConfirmationDialog
    :request="confirmation"
    @settle="closeConfirmation"
  />
</template>
