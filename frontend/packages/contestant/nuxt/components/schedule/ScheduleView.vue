<script setup lang="ts">
import { computed } from "vue";
import RequestState from "@ictsc/ui/components/RequestState.vue";
import type { ScheduleEntry } from "~/features/models";
import {
  scheduleEntryState,
  scheduleEntryStateLabels,
  sortedScheduleEntries,
} from "~/features/schedule/entries";

const props = defineProps<{
  entries?: ScheduleEntry[];
  now: number;
  pending?: boolean;
  error?: unknown;
}>();
defineEmits<{ retry: [] }>();

const rows = computed(() =>
  sortedScheduleEntries(props.entries ?? []).map((entry) => ({
    entry,
    state: scheduleEntryState(entry, props.now),
  })),
);
const format = (iso: string) =>
  new Date(iso).toLocaleString("ja-JP", {
    timeZone: "Asia/Tokyo",
    month: "numeric",
    day: "numeric",
    weekday: "short",
    hour: "2-digit",
    minute: "2-digit",
  });
</script>
<template>
  <header class="linked-page-heading">
    <h1 class="page-title">スケジュール</h1>
    <p>競技時間（日本時間）</p>
  </header>
  <RequestState :error="error" :pending="pending" @retry="$emit('retry')"
    ><div class="table-scroll schedule-table">
      <table class="data-table">
        <thead>
          <tr>
            <th>枠</th>
            <th>開始</th>
            <th>終了</th>
            <th>状態</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="{ entry, state } in rows"
            :key="entry.name"
            :class="{ 'is-running': state === 'running' }"
          >
            <td data-label="枠">{{ entry.name }}</td>
            <td class="numeric" data-label="開始">
              <time :datetime="entry.startAt">{{ format(entry.startAt) }}</time>
            </td>
            <td class="numeric" data-label="終了">
              <time :datetime="entry.endAt">{{ format(entry.endAt) }}</time>
            </td>
            <td data-label="状態">
              <span class="schedule-state" :class="`is-${state}`">{{
                scheduleEntryStateLabels[state]
              }}</span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <p v-if="!rows.length">
      スケジュールはまだ公開されていません。
    </p></RequestState
  >
</template>
<style scoped>
.schedule-table tr.is-running td:first-child {
  box-shadow: inset 3px 0 0 var(--accent);
}
.schedule-state {
  font-size: 12px;
  font-weight: 600;
}
.schedule-state.is-running {
  color: var(--accent);
}
.schedule-state.is-ended {
  color: var(--muted);
}
</style>
