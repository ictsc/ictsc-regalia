<script setup lang="ts">
import { ref, watch } from "vue";
import AsciiLogo from "@ictsc/ui/components/AsciiLogo.vue";
import ThemeToggle from "@ictsc/ui/components/ThemeToggle.vue";
import type { HeaderRankingItem } from "~/features/ranking";

const props = defineProps<{
  /** Simplified header used on problem detail pages. */
  compact: boolean;
  rankingItems: HeaderRankingItem[];
  ownTeamCode?: number;
  maxScore: number;
  clock: { label: string; text: string };
  demoMode: boolean;
  hasUnreadNotices: boolean;
  teamName?: string;
  /** Closes the mobile menu whenever it changes. */
  routeKey: string;
}>();

const mobileMenu = ref<HTMLDetailsElement | null>(null);
watch(
  () => props.routeKey,
  () => {
    if (mobileMenu.value) mobileMenu.value.open = false;
  },
);
</script>
<template>
  <header
    class="competition-header"
    :class="{ 'simple-header-version': compact }"
  >
    <NuxtLink
      class="brand"
      to="/problems"
      aria-label="ICTSC 2026 REGALIA ホーム"
    >
      <AsciiLogo />
    </NuxtLink>
    <div v-if="!compact" class="header-performance">
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
              :class="{ 'is-team': item.entry.teamCode === ownTeamCode }"
            >
              <small>{{ item.entry.rank }}位</small>
              <b>
                {{ item.entry.score.toLocaleString()
                }}<em v-if="item.entry.teamCode === ownTeamCode">
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
        <span v-if="teamName" class="mobile-team-badge">{{ teamName }}</span
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
      </nav>
    </details>
  </header>
</template>
