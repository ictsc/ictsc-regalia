<script setup lang="ts">
const { data, refresh } = await useCompetition();
const { viewer } = useSession();
const route = useRoute();
const now = useClock();
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
const own = computed(() =>
  data.value?.ranking.ranking.find(
    (r) =>
      viewer.value?.state === "CONTESTANT" &&
      r.teamCode === viewer.value.team.code,
  ),
);
const neighbors = computed(
  () =>
    data.value?.ranking.ranking.filter(
      (r) =>
        r.rank === 1 || (own.value && Math.abs(r.rank - own.value.rank) <= 1),
    ) ?? [],
);
const maxScore = computed(
  () => data.value?.problems.reduce((sum, p) => sum + p.maxScore, 0) ?? 0,
);
const solved = computed(
  () => data.value?.problems.filter((p) => p.score != null).length ?? 0,
);
const notice = computed(
  () =>
    data.value?.notices.toSorted(
      (a, b) => Date.parse(b.effectiveFrom) - Date.parse(a.effectiveFrom),
    )[0],
);
const clock = computed(() => {
  const sections = data.value?.schedule?.entries ?? [];
  const active = sections.find(
    (s) =>
      Date.parse(s.startAt) <= now.value && now.value < Date.parse(s.endAt),
  );
  const next = sections.find((s) => Date.parse(s.startAt) > now.value);
  const end = active?.endAt ?? next?.startAt;
  const seconds = end
    ? Math.max(0, Math.floor((Date.parse(end) - now.value) / 1000))
    : 0;
  return {
    label: active ? "残り時間" : next ? "開始まで" : "競技終了",
    text: [
      Math.floor(seconds / 3600),
      Math.floor(seconds / 60) % 60,
      seconds % 60,
    ]
      .map((n) => String(n).padStart(2, "0"))
      .join(":"),
  };
});
</script>
<template>
  <header
    class="competition-header"
    :class="{ 'simple-header-version': detail }"
  >
    <NuxtLink
      class="brand"
      to="/problems"
      aria-label="ICTSC 2026 REGALIA ホーム"
      ><span class="brand-name">ICTSC</span
      ><span class="brand-edition">2026<br />REGALIA</span></NuxtLink
    >
    <div v-if="!detail" class="header-performance">
      <NuxtLink class="header-solved" to="/activity"
        ><span class="header-solved-value"
          ><strong>{{ solved }}</strong
          ><span>/ {{ data?.problems.length ?? "—" }}問</span></span
        ><small>提出履歴を見る →</small></NuxtLink
      ><NuxtLink class="header-ranking" to="/ranking"
        ><span class="header-neighbors"
          ><span
            v-for="r in neighbors"
            :key="r.teamCode"
            :class="{ 'is-team': r.teamCode === own?.teamCode }"
            ><small>{{ r.rank }}位</small
            ><b
              >{{ r.score.toLocaleString()
              }}<em v-if="r.teamCode === own?.teamCode">
                / {{ maxScore.toLocaleString() }}</em
              ></b
            ></span
          ></span
        ><span class="header-ranking-link">順位表を見る →</span></NuxtLink
      >
    </div>
    <nav v-else class="simple-header-nav" aria-label="メインメニュー">
      <NuxtLink to="/activity">提出履歴</NuxtLink
      ><NuxtLink to="/ranking">順位表</NuxtLink>
    </nav>
    <div class="competition-clock">
      <time>{{ clock.text }}</time
      ><span>{{ clock.label }}</span>
    </div>
    <aside class="competition-notice">
      <strong>お知らせ</strong
      ><time>{{
        notice
          ? new Date(notice.effectiveFrom).toLocaleTimeString("ja-JP", {
              hour: "2-digit",
              minute: "2-digit",
            })
          : "—"
      }}</time>
      <p>{{ notice?.title ?? "現在お知らせはありません" }}</p>
      <NuxtLink to="/announces">お知らせを確認する →</NuxtLink>
    </aside>
    <details class="mobile-menu">
      <summary>
        <span aria-hidden="true" /><b class="visually-hidden">メニューを開く</b>
      </summary>
      <nav aria-label="モバイルメニュー">
        <NuxtLink to="/problems">問題一覧</NuxtLink
        ><NuxtLink to="/activity">提出履歴</NuxtLink
        ><NuxtLink to="/ranking">順位表</NuxtLink
        ><NuxtLink to="/announces">お知らせ</NuxtLink>
      </nav>
    </details>
  </header>
</template>
