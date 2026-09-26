<script setup lang="ts">
import { headerRankingItems } from "~/features/ranking";
import { useReadAnnouncements } from "~/features/announce-read-status";

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
const { getUnreadNotices } = useReadAnnouncements();
const hasUnreadNotices = computed(
  () => getUnreadNotices(data.value?.notices ?? []).length > 0,
);
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
const clock = computed(() => {
  // Keep the demo presentation independent of the schedule and page-load time.
  if (demoMode.value) return { label: "デモモード", text: "02:00:00" };

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
    >
      <AsciiLogo />
    </NuxtLink>
    <div v-if="!detail" class="header-performance">
      <div class="header-ranking" aria-label="現在の順位">
        <small class="header-ranking-label">順位</small>
        <span class="header-neighbors">
          <template
            v-for="item in rankingItems"
            :key="
              item.kind === 'rank'
                ? `rank-${item.entry.teamCode}`
                : `tie-${item.rank}`
            "
          >
            <span v-if="item.kind === 'tie'" class="is-tie">
              <small>同率チームあり</small>
              <b>{{ item.count }}チーム同率</b>
            </span>
            <span
              v-else
              :class="{ 'is-team': item.entry.teamCode === own?.teamCode }"
            >
              <small>{{ item.entry.rank }}位</small>
              <b>
                {{ item.entry.score.toLocaleString()
                }}<em v-if="item.entry.teamCode === own?.teamCode">
                  / {{ maxScore.toLocaleString() }}</em
                >
              </b>
            </span>
          </template>
        </span>
      </div>
      <nav class="header-account" aria-label="メインメニュー">
        <NuxtLink class="header-nav-primary" to="/problems">問題一覧</NuxtLink>
        <NuxtLink class="header-nav-primary" to="/activity">提出履歴</NuxtLink>
        <NuxtLink class="header-nav-primary" to="/ranking">順位表</NuxtLink>
        <NuxtLink
          class="header-nav-primary header-notifications"
          to="/announces"
          >通知<span v-if="hasUnreadNotices" class="notification-dot"
            ><span class="visually-hidden">未読の通知あり</span></span
          ></NuxtLink
        >
        <NuxtLink class="header-nav-secondary" to="/teams">チーム</NuxtLink>
        <NuxtLink class="header-nav-secondary" to="/rule">ルール</NuxtLink>
        <NuxtLink class="header-nav-secondary" to="/profile"
          >プロフィール</NuxtLink
        >
        <ThemeToggle class="header-nav-secondary" />
        <button
          class="header-nav-secondary"
          :aria-label="accountActionAriaLabel"
          :class="{ 'is-impersonating': isImpersonating }"
          @click="emit('logout')"
        >
          {{ accountActionLabel }}
        </button>
      </nav>
    </div>
    <nav
      v-else
      class="simple-header-nav header-account"
      aria-label="メインメニュー"
    >
      <NuxtLink class="header-nav-primary" to="/problems">問題一覧</NuxtLink>
      <NuxtLink class="header-nav-primary" to="/activity">提出履歴</NuxtLink>
      <NuxtLink class="header-nav-primary" to="/ranking">順位表</NuxtLink>
      <NuxtLink class="header-nav-primary header-notifications" to="/announces"
        >通知<span v-if="hasUnreadNotices" class="notification-dot"
          ><span class="visually-hidden">未読の通知あり</span></span
        ></NuxtLink
      >
      <NuxtLink class="header-nav-secondary" to="/teams">チーム</NuxtLink>
      <NuxtLink class="header-nav-secondary" to="/rule">ルール</NuxtLink>
      <NuxtLink class="header-nav-secondary" to="/profile"
        >プロフィール</NuxtLink
      >
      <ThemeToggle class="header-nav-secondary" />
      <button
        class="header-nav-secondary"
        :aria-label="accountActionAriaLabel"
        :class="{ 'is-impersonating': isImpersonating }"
        @click="emit('logout')"
      >
        {{ accountActionLabel }}
      </button>
    </nav>
    <div class="competition-clock">
      <time>{{ clock.text }}</time
      ><span :class="{ 'is-demo': demoMode }">{{ clock.label }}</span>
    </div>
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
        ><NuxtLink
          class="mobile-menu-primary mobile-notifications"
          to="/announces"
          >通知<span v-if="hasUnreadNotices" class="notification-dot"
            ><span class="visually-hidden">未読の通知あり</span></span
          ></NuxtLink
        ><NuxtLink to="/teams">チーム</NuxtLink
        ><NuxtLink to="/rule">ルール</NuxtLink
        ><NuxtLink to="/profile">プロフィール</NuxtLink>
        <ThemeToggle />
        <button
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
