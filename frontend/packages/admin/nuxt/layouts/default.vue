<script setup lang="ts">
const { viewer, refresh, signOut } = useAdminSession();
const route = useRoute();
const error = ref<unknown>();
async function load() {
  try {
    await refresh();
    error.value = undefined;
  } catch (e) {
    error.value = e;
  }
}
await load();
async function logout() {
  try {
    await signOut();
  } catch (e) {
    error.value = e;
  }
}
const links = [
  ["/", "概要"],
  ["/teams", "チーム"],
  ["/contestants", "参加者"],
  ["/content", "コンテンツ"],
  ["/submissions", "採点"],
  ["/scores", "得点・順位"],
  ["/deployments", "再展開"],
  ["/settings", "競技設定"],
];
</script>
<template>
  <div class="admin-app">
    <header class="admin-header">
      <NuxtLink class="brand" to="/"
        ><span class="brand-name">ICTSC</span
        ><span class="brand-edition">REGALIA<br />ADMIN</span></NuxtLink
      >
      <div class="actions">
        <ThemeToggle />
        <template v-if="viewer?.state !== 'ANONYMOUS'">
          <span>{{
            viewer?.state === "ADMIN" ? viewer.admin.discord.display_name : ""
          }}</span
          ><button class="button-secondary" @click="logout">ログアウト</button>
        </template>
      </div>
    </header>
    <RequestState :error="error" @retry="load"
      ><main v-if="viewer?.state === 'ANONYMOUS'" class="workspace">
        <h1>ICTSC Admin</h1>
        <p>Discordの運営ロールでログインしてください。</p>
        <a
          class="button-primary"
          href="/api/v1/admin/auth/discord?next=%2Fadmin%2F"
          >Discordでログイン</a
        >
      </main>
      <template v-else-if="viewer"
        ><div class="admin-shell">
          <nav class="admin-nav" aria-label="運営メニュー">
            <p class="admin-nav-heading">運営メニュー</p>
            <NuxtLink
              v-for="[to, label] in links"
              :key="to"
              :to="to!"
              :class="{
                'is-current':
                  route.path === to ||
                  (to !== '/' && route.path.startsWith(`${to}/`)),
              }"
              :aria-current="
                route.path === to ||
                (to !== '/' && route.path.startsWith(`${to}/`))
                  ? 'page'
                  : undefined
              "
              >{{ label }}</NuxtLink
            >
          </nav>
          <slot /></div></template
    ></RequestState>
  </div>
</template>
<style>
.admin-app {
  min-height: 100vh;
}
.admin-app .admin-header {
  max-width: none;
  min-height: 72px;
  padding-inline: clamp(16px, 3vw, 48px);
}
.admin-app .admin-header .actions {
  margin: 0;
  font-size: 13px;
}
.admin-app .admin-header .actions > span {
  font-weight: 600;
}
.admin-app .admin-header .button-secondary {
  min-height: 38px;
  padding: 7px 12px;
}
.admin-shell {
  display: grid;
  grid-template-columns: 216px minmax(0, 1fr);
  max-width: 1680px;
  margin-inline: auto;
}
.admin-shell:has(> .grading-detail-workspace) {
  max-width: none;
}
.admin-shell .admin-nav {
  align-self: stretch;
  display: flex;
  flex-direction: column;
  flex-wrap: nowrap;
  align-items: stretch;
  gap: 3px;
  min-height: calc(100vh - 72px);
  padding: 30px 14px;
  border-right: 1px solid var(--line-light);
  border-bottom: 0;
}
.admin-nav-heading {
  margin: 0 10px 12px;
  color: var(--muted);
  font-family: var(--mono);
  font-size: 11px;
  letter-spacing: 0.08em;
}
.admin-shell .admin-nav a {
  display: block;
  padding: 11px 12px;
  border-left: 3px solid transparent;
  font-weight: 600;
  line-height: 1.4;
}
.admin-shell .admin-nav a:hover {
  background: var(--surface-hover);
}
.admin-shell .admin-nav .is-current {
  border-left-color: var(--accent);
  background: var(--surface-hover);
}
.admin-shell > .workspace {
  width: 100%;
  min-width: 0;
  max-width: none;
  padding: 34px clamp(20px, 3vw, 48px) 80px;
}
.admin-shell .workspace h1 {
  margin-bottom: 10px;
  font-size: clamp(27px, 3vw, 36px);
}
.admin-shell .workspace h2 {
  margin: 0 0 16px;
  font-size: 19px;
}
.admin-lead {
  margin: 0 0 30px;
  color: var(--muted);
  line-height: 1.7;
}
.admin-panel {
  min-width: 0;
  margin: 0 0 28px;
  padding: 24px;
  border: 1px solid var(--line-light);
}
.admin-panel > :first-child {
  margin-top: 0;
}
.admin-panel > :last-child {
  margin-bottom: 0;
}
.admin-panel-heading {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  gap: 16px;
  margin-bottom: 14px;
}
.admin-panel-heading h2 {
  margin: 0 !important;
}
.admin-panel-heading a {
  font-size: 13px;
}
.admin-panel .data-table {
  margin: 0;
}
.admin-panel .form-grid {
  margin: 0;
  max-width: none;
}
.admin-panel .table-scroll {
  max-width: 100%;
}
.admin-panel .data-table th,
.admin-panel .data-table td {
  border-color: var(--line-light);
}
.admin-panel .data-table th:first-child,
.admin-panel .data-table td:first-child {
  padding-left: 0;
}
.admin-panel .data-table th:last-child,
.admin-panel .data-table td:last-child {
  padding-right: 0;
}
.admin-panel .data-table tr:last-child td {
  border-bottom: 0;
}
.admin-summary-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 14px;
  margin-bottom: 28px;
}
.admin-summary-card {
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 14px;
  padding: 20px;
  border: 1px solid var(--line-light);
}
.admin-summary-card:hover,
.admin-summary-card:focus-visible {
  border-color: var(--accent);
  background: var(--surface-hover);
}
.admin-summary-label {
  font-size: 13px;
  font-weight: 600;
}
.admin-summary-number {
  font-family: var(--mono);
  font-size: clamp(30px, 3vw, 45px);
  font-weight: 500;
  line-height: 1;
}
.admin-summary-footer {
  margin-top: auto;
  color: var(--accent);
  font-size: 12px;
}
.admin-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 8px 20px;
  align-items: center;
  margin: 0 0 18px;
}
.admin-meta code,
.admin-code {
  font-family: var(--mono);
  font-size: 12px;
  overflow-wrap: anywhere;
}
.admin-badge {
  display: inline-flex;
  align-items: center;
  padding: 4px 8px;
  border: 1px solid var(--line-light);
  font-size: 12px;
  font-weight: 600;
  line-height: 1.3;
}
.admin-badge.is-active {
  border-color: var(--accent);
  color: var(--accent);
}
.admin-empty {
  padding: 28px 12px;
  color: var(--muted);
  text-align: center;
}
.admin-grid-two {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 20px;
}
.admin-grid-two > .admin-panel {
  margin-bottom: 0;
}
.admin-filter-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
}
.admin-filter-grid label {
  min-width: 0;
  display: grid;
  gap: 8px;
  font-weight: 600;
}
.admin-filter-grid select,
.admin-filter-grid input {
  width: 100%;
  min-width: 0;
  padding: 10px 12px;
  border: 1px solid var(--line-light);
  background: var(--paper);
  color: var(--ink);
}
.admin-panel .actions {
  margin: 16px 0 0;
  gap: 10px;
}
.admin-panel .button-primary,
.admin-panel .button-secondary {
  min-height: 42px;
  padding: 9px 16px;
}
.admin-panel .button-primary {
  width: fit-content;
}
.admin-panel .form-grid label {
  min-width: 0;
}
.admin-panel .form-grid input,
.admin-panel .form-grid textarea,
.admin-panel .form-grid select {
  width: 100%;
  border-color: var(--line-light);
}
.admin-panel .form-grid textarea {
  resize: vertical;
}
.admin-panel .markdown {
  min-width: 0;
  overflow-wrap: anywhere;
}
.admin-panel .mono,
.admin-panel pre {
  overflow-wrap: anywhere;
}
.admin-panel pre {
  max-width: 100%;
  overflow-x: auto;
}
.admin-panel .table-scroll .data-table {
  min-width: 580px;
}
.admin-panel .actions label {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  cursor: pointer;
}
.admin-panel .actions input[type="checkbox"] {
  accent-color: var(--accent);
}
.admin-danger {
  border-color: var(--accent);
}
.admin-danger h2 {
  color: var(--accent);
}
.admin-danger .button-primary {
  background: var(--accent-surface);
}
.admin-panel p[role="status"]:empty {
  display: none;
}
@media (max-width: 980px) {
  .admin-shell {
    display: block;
  }
  .admin-shell .admin-nav {
    min-height: 0;
    overflow-x: auto;
    flex-direction: row;
    gap: 2px;
    padding: 8px 14px;
    border-right: 0;
    border-bottom: 1px solid var(--line-light);
    scrollbar-width: thin;
  }
  .admin-shell .admin-nav a {
    white-space: nowrap;
    border-left: 0;
    border-bottom: 3px solid transparent;
  }
  .admin-shell .admin-nav .is-current {
    border-bottom-color: var(--accent);
  }
  .admin-nav-heading {
    display: none;
  }
  .admin-shell > .workspace {
    padding: 28px 20px 64px;
  }
}
@media (max-width: 650px) {
  .admin-app .admin-header {
    gap: 8px;
    padding-block: 10px;
  }
  .admin-app .admin-header .actions {
    gap: 8px;
  }
  .admin-app .admin-header .actions > span {
    display: none;
  }
  .admin-shell > .workspace {
    padding: 24px 14px 60px;
  }
  .admin-panel {
    padding: 18px 14px;
  }
  .admin-summary-grid,
  .admin-grid-two {
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 10px;
  }
  .admin-summary-card {
    padding: 14px;
  }
  .admin-summary-number {
    font-size: 30px;
  }
  .admin-filter-grid {
    grid-template-columns: 1fr;
  }
}
@media (max-width: 410px) {
  .admin-summary-grid,
  .admin-grid-two {
    grid-template-columns: 1fr;
  }
}
</style>
