<script setup lang="ts">
import { useQuery } from "@tanstack/vue-query";

const { data, isPending } = useQuery({
  queryKey: ["frontend-health"],
  queryFn: async () => ({
    status: "ready",
    checkedAt: new Date().toISOString(),
  }),
});
</script>

<template>
  <main class="page">
    <section class="hero">
      <p class="eyebrow">ICTSC Regalia</p>
      <h1>Score Server Frontend</h1>
      <p class="lead">Nuxt and TanStack Query are ready to build on.</p>
      <dl class="status">
        <div>
          <dt>Query</dt>
          <dd>{{ isPending ? "loading" : data?.status }}</dd>
        </div>
        <div>
          <dt>Runtime</dt>
          <dd>Bun + Nuxt</dd>
        </div>
      </dl>
    </section>
  </main>
</template>

<style scoped>
:global(body) {
  margin: 0;
  background: #f6f7f9;
  color: #172026;
  font-family:
    Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI",
    sans-serif;
}

.page {
  min-height: 100vh;
  display: grid;
  place-items: center;
  padding: 32px;
}

.hero {
  width: min(100%, 720px);
}

.eyebrow {
  margin: 0 0 12px;
  color: #246b5f;
  font-size: 0.875rem;
  font-weight: 700;
  letter-spacing: 0;
  text-transform: uppercase;
}

h1 {
  margin: 0;
  font-size: clamp(2.5rem, 8vw, 5.5rem);
  line-height: 0.95;
  letter-spacing: 0;
}

.lead {
  max-width: 560px;
  margin: 24px 0 0;
  color: #4d5a62;
  font-size: 1.125rem;
  line-height: 1.7;
}

.status {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 1px;
  margin: 36px 0 0;
  overflow: hidden;
  border: 1px solid #d8dee4;
  border-radius: 8px;
  background: #d8dee4;
}

.status div {
  padding: 18px;
  background: #ffffff;
}

.status dt {
  color: #687782;
  font-size: 0.75rem;
  font-weight: 700;
  text-transform: uppercase;
}

.status dd {
  margin: 8px 0 0;
  font-size: 1rem;
  font-weight: 700;
}

@media (max-width: 560px) {
  .page {
    padding: 24px;
    place-items: start;
  }

  .status {
    grid-template-columns: 1fr;
  }
}
</style>
