<script setup lang="ts">
import { headerRankingItems } from "~/features/ranking";

const { data, refresh } = await useCompetition();
const { viewer } = useSession();
const emit = defineEmits<{ logout: [] }>();
const route = useRoute();
const mobileMenu = ref<HTMLDetailsElement | null>(null);
const now = useClock();
const demoMode = useDemoMode();
const isImpersonating = computed(
  () =>
    viewer.value?.state === "CONTESTANT" &&
    viewer.value.impersonated_by != null,
);
const accountActionLabel = computed(() =>
  viewer.value?.state === "CONTESTANT" && viewer.value.impersonated_by
    ? "代"
    : "ログアウト",
);
const accountActionAriaLabel = computed(() =>
  viewer.value?.state === "CONTESTANT" && viewer.value.impersonated_by
    ? `${viewer.value.impersonated_by}による代理操作を終了`
    : "ログアウト",
);
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
watch(
  () => route.fullPath,
  () => {
    if (mobileMenu.value) mobileMenu.value.open = false;
  },
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
          ><template
            v-for="item in rankingItems"
            :key="
              item.kind === 'rank'
                ? `rank-${item.entry.teamCode}`
                : `tie-${item.rank}`
            "
            ><span v-if="item.kind === 'tie'" class="is-tie"
              ><small>同率チームあり</small
              ><b>{{ item.count }}チーム同率</b></span
            ><span
              v-else
              :class="{ 'is-team': item.entry.teamCode === own?.teamCode }"
              ><small>{{ item.entry.rank }}位</small
              ><b
                >{{ item.entry.score.toLocaleString()
                }}<em v-if="item.entry.teamCode === own?.teamCode">
                  / {{ maxScore.toLocaleString() }}</em
                ></b
              ></span
            ></template
          ></span
        ><span class="header-ranking-link">順位表を見る →</span></NuxtLink
      >
      <nav class="header-account" aria-label="アカウント">
        <NuxtLink to="/teams">チーム</NuxtLink
        ><NuxtLink to="/rule">ルール</NuxtLink
        ><NuxtLink to="/profile">プロフィール</NuxtLink
        ><button
          :aria-label="accountActionAriaLabel"
          :class="{ 'is-impersonating': isImpersonating }"
          @click="emit('logout')"
        >
          {{ accountActionLabel }}
        </button>
      </nav>
    </div>
    <nav v-else class="simple-header-nav" aria-label="メインメニュー">
      <NuxtLink to="/activity">提出履歴</NuxtLink
      ><NuxtLink to="/ranking">順位表</NuxtLink
      ><span class="header-account">
        <NuxtLink to="/teams">チーム</NuxtLink
        ><NuxtLink to="/rule">ルール</NuxtLink
        ><NuxtLink to="/profile">プロフィール</NuxtLink
        ><button
          :aria-label="accountActionAriaLabel"
          :class="{ 'is-impersonating': isImpersonating }"
          @click="emit('logout')"
        >
          {{ accountActionLabel }}
        </button>
      </span>
    </nav>
    <div class="competition-clock">
      <time>{{ clock.text }}</time
      ><span :class="{ 'is-demo': demoMode }">{{
        demoMode ? "デモモード" : clock.label
      }}</span>
    </div>
    <NuxtLink class="competition-notice" to="/announces">
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
      <span class="notice-link-label">お知らせを確認する →</span>
    </NuxtLink>
    <details ref="mobileMenu" class="mobile-menu">
      <summary>
        <span aria-hidden="true" /><b class="visually-hidden">メニューを開く</b>
      </summary>
      <nav aria-label="モバイルメニュー">
        <span v-if="viewer?.state === 'CONTESTANT'" class="mobile-team-badge">{{
          viewer.team.name
        }}</span
        ><NuxtLink class="mobile-menu-primary" to="/problems">問題一覧</NuxtLink
        ><NuxtLink class="mobile-menu-primary" to="/activity">提出履歴</NuxtLink
        ><NuxtLink class="mobile-menu-primary" to="/ranking">順位表</NuxtLink
        ><NuxtLink class="mobile-menu-primary" to="/announces"
          >お知らせ</NuxtLink
        ><NuxtLink to="/teams">チーム</NuxtLink
        ><NuxtLink to="/rule">ルール</NuxtLink
        ><NuxtLink to="/profile">プロフィール</NuxtLink
        ><button
          :aria-label="accountActionAriaLabel"
          :class="{ 'is-impersonating': isImpersonating }"
          @click="emit('logout')"
        >
          {{ accountActionLabel }}
        </button>
      </nav>
    </details>
  </header>
</template>
