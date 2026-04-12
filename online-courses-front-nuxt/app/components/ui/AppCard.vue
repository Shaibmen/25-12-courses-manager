<script setup lang="ts">
defineProps<{
  title: string
}>()
</script>

<template>
  <section class="app-card">
    <header class="app-card__header">
      <h3>{{ title }}</h3>
    </header>
    <div class="app-card__body">
      <slot />
    </div>
  </section>
</template>

<style scoped>
.app-card {
  padding: 1.25rem;
  border-radius: 1.25rem;
  border: 1px solid rgba(15, 23, 42, 0.08);
  background: rgba(255, 255, 255, 0.85);
  box-shadow: 0 20px 60px rgba(15, 23, 42, 0.08);
}

.app-card__header h3 {
  margin: 0;
  font-size: 1rem;
  color: #0f172a;
}

.app-card__body {
  margin-top: 0.85rem;
  color: #334155;
}
</style>
