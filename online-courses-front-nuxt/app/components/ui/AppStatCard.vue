<script setup lang="ts">
defineProps<{
  title: string
  value: string | number
  description?: string
}>()
</script>

<template>
  <article class="stat-card">
    <p class="stat-card__title">{{ title }}</p>
    <strong class="stat-card__value">{{ value }}</strong>
    <p v-if="description" class="stat-card__description">
      {{ description }}
    </p>
  </article>
</template>

<style scoped>
.stat-card {
  padding: 1.2rem;
  border-radius: 1.25rem;
  background:
    linear-gradient(180deg, rgba(255, 255, 255, 0.96) 0%, rgba(239, 246, 255, 0.92) 100%);
  border: 1px solid rgba(15, 23, 42, 0.08);
  box-shadow: 0 18px 40px rgba(15, 23, 42, 0.08);
}

.stat-card__title {
  margin: 0;
  color: #475569;
  font-size: 0.9rem;
}

.stat-card__value {
  display: block;
  margin-top: 0.55rem;
  font-size: clamp(1.6rem, 3vw, 2.3rem);
  color: #0f172a;
}

.stat-card__description {
  margin: 0.65rem 0 0;
  color: #64748b;
  font-size: 0.88rem;
}
</style>
