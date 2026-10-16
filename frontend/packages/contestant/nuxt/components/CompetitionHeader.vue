<script setup lang="ts">
import CompetitionHeaderView from "~/components/header/CompetitionHeaderView.vue";
import { headerRankingItems } from "~/features/ranking";
import { useReadAnnouncements } from "~/features/announce-read-status";
import { contestClock } from "~/features/schedule/clock";

const { data, refresh } = await useCompetition();
const { viewer } = useSession();
const route = useRoute();
const now = useClock();
const demoMode = useDemoMode();
const boundary = computed(
  () =>
    data.value?.schedule?.entries
      .flatMap((s) => [Date.parse(s.startAt), Date.parse(s.endAt)])
      .filter((t) => t > now.value)
      .sort((a, b) => a - b)[0],
);
let nextBoundary: number | undefined;
watch(now, (value) => {
  if (nextBoundary != null && value >= nextBoundary) {
    nextBoundary = undefined;
    void refresh();
  } else if (nextBoundary == null) nextBoundary = boundary.value;
});
const detail = computed(() => /^\/problems\/.+/.test(route.path));
const { getUnreadNotices } = useReadAnnouncements();
const hasUnreadNotices = computed(
  () => getUnreadNotices(data.value?.notices ?? []).length > 0,
);
const own = computed(() =>
  data.value?.ranking.ranking.find(
    (r) =>
      viewer.value?.state === "CONTESTANT" &&
      r.teamCode === viewer.value.team.code,
  ),
);
const rankingItems = computed(() =>
  headerRankingItems(data.value?.ranking.ranking ?? [], own.value?.teamCode),
);
const maxScore = computed(
  () => data.value?.problems.reduce((sum, p) => sum + p.maxScore, 0) ?? 0,
);
const clock = computed(() => {
  // Keep the demo presentation independent of the schedule and page-load time.
  if (demoMode.value) return { label: "デモモード", text: "02:00:00" };

  return contestClock(data.value?.schedule?.entries ?? [], now.value);
});
</script>
<template>
  <CompetitionHeaderView
    :compact="detail"
    :ranking-items="rankingItems"
    :own-team-code="own?.teamCode"
    :max-score="maxScore"
    :clock="clock"
    :demo-mode="demoMode"
    :has-unread-notices="hasUnreadNotices"
    :team-name="viewer?.state === 'CONTESTANT' ? viewer.team.name : undefined"
    :route-key="route.fullPath"
  />
</template>
