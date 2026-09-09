<script setup lang="ts">
import { api, expectData } from "@ictsc/api";
import { createAdminActions } from "~/features/admin-api";
useHead({ title: "競技設定" });
const { data, error, pending, refresh } = await useAsyncData(
  "settings",
  async () => {
    const [rule, schedule] = await Promise.all([
      api.GET("/api/v1/admin/rule").then(expectData),
      api.GET("/api/v1/admin/dashboard-schedule").then(expectData),
    ]);
    return { rule: rule.rule, schedule: schedule.dashboard_schedule };
  },
);
const markdown = ref(data.value?.rule.markdown ?? "");
const freeze = ref(
  data.value?.schedule.ranking_freeze_at
    ? new Date(
        Date.parse(data.value.schedule.ranking_freeze_at) -
          new Date().getTimezoneOffset() * 60000,
      )
        .toISOString()
        .slice(0, 16)
    : "",
);
const { busy, message, run } = useMutation(refresh);
const actions = createAdminActions(api);
</script>
<template>
  <main class="workspace">
    <h1>競技設定</h1>
    <p role="status">{{ message }}</p>
    <RequestState :error="error" :pending="pending" @retry="refresh"
      ><h2>ルール編集</h2>
      <form
        class="form-grid"
        @submit.prevent="run('ルール更新', () => actions.replaceRule(markdown))"
      >
        <label>Markdown<textarea v-model="markdown" rows="18" /></label
        ><button class="button-primary" :disabled="busy">ルールを更新</button>
      </form>
      <h2>プレビュー</h2>
      <MarkdownContent :source="markdown" />
      <h2>ランキング凍結</h2>
      <p>
        ブラウザのローカル時刻で入力します。空欄で保存すると凍結を解除します。
      </p>
      <form
        class="form-grid"
        @submit.prevent="
          run('凍結時刻更新', () =>
            actions.replaceFreeze(
              freeze ? new Date(freeze).toISOString() : null,
            ),
          )
        "
      >
        <label>凍結時刻<input v-model="freeze" type="datetime-local" /></label
        ><button class="button-primary" :disabled="busy">凍結時刻を保存</button>
      </form></RequestState
    >
  </main>
</template>
