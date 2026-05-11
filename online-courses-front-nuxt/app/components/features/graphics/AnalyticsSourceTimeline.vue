<script setup lang="ts">
const props = withDefaults(defineProps<{
  items: Array<{
    monthLabel: string
    total: number
    segments: Array<{
      label: string
      value: number
      color: string
    }>
  }>
  emptyText?: string
}>(), {
  emptyText: 'Нет данных о зачислениях.'
})

const formatValue = (value: number) => new Intl.NumberFormat('ru-RU').format(value)
</script>

<template>
  <div v-if="items.length" class="timeline">
    <article
      v-for="item in items"
      :key="item.monthLabel"
      class="timeline__item"
    >
      <div class="timeline__head">
        <strong>{{ item.monthLabel }}</strong>
        <span>{{ formatValue(item.total) }} зачислений</span>
      </div>

      <div class="timeline__track">
        <div
          v-for="segment in item.segments"
          :key="`${item.monthLabel}-${segment.label}`"
          class="timeline__segment"
          :style="{
            width: `${item.total ? (segment.value / item.total) * 100 : 0}%`,
            background: segment.color
          }"
          :title="`${segment.label}: ${formatValue(segment.value)}`"
        />
      </div>

      <div class="timeline__legend">
        <span
          v-for="segment in item.segments"
          :key="`${item.monthLabel}-${segment.label}-legend`"
          class="timeline__legend-item"
        >
          <i :style="{ background: segment.color }" />
          {{ segment.label }}: {{ formatValue(segment.value) }}
        </span>
      </div>
    </article>
  </div>

  <div v-else class="timeline__empty">
    {{ emptyText }}
  </div>
</template>

<style scoped>
.timeline {
  display: grid;
  gap: 1rem;
}

.timeline__item {
  display: grid;
  gap: 0.7rem;
  padding: 1rem;
  border-radius: 1.1rem;
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.96), rgba(240, 249, 255, 0.92));
  border: 1px solid rgba(14, 116, 144, 0.1);
}

.timeline__head {
  display: flex;
  justify-content: space-between;
  gap: 1rem;
  color: #0f172a;
}

.timeline__head span {
  color: #475569;
}

.timeline__track {
  display: flex;
  min-height: 1rem;
  overflow: hidden;
  border-radius: 999px;
  background: rgba(148, 163, 184, 0.16);
}

.timeline__segment {
  min-width: 0.35rem;
}

.timeline__legend {
  display: flex;
  flex-wrap: wrap;
  gap: 0.75rem;
}

.timeline__legend-item {
  display: inline-flex;
  align-items: center;
  gap: 0.45rem;
  font-size: 0.86rem;
  color: #475569;
}

.timeline__legend-item i {
  display: inline-block;
  width: 0.8rem;
  height: 0.8rem;
  border-radius: 999px;
}

.timeline__empty {
  color: #64748b;
}

@media (max-width: 720px) {
  .timeline__head {
    flex-direction: column;
  }
}
</style>
