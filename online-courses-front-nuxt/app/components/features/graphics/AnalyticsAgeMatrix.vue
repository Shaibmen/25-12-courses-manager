<script setup lang="ts">
const props = withDefaults(defineProps<{
  rows: Array<{
    label: string
    values: Array<{
      ageRange: string
      listeners: number
    }>
  }>
  ageRanges: string[]
  emptyText?: string
}>(), {
  emptyText: 'Нет данных по возрастным группам.'
})

const maxListeners = computed(() => {
  const values = props.rows.flatMap(row => row.values.map(item => item.listeners))
  return values.length ? Math.max(...values, 1) : 1
})

const intensity = (listeners: number) => {
  if (!listeners) {
    return 0.08
  }

  return 0.22 + (listeners / maxListeners.value) * 0.78
}
</script>

<template>
  <div v-if="rows.length" class="matrix">
    <div class="matrix__header">
      <span>Программа</span>
      <span
        v-for="ageRange in ageRanges"
        :key="ageRange"
        class="matrix__header-cell"
      >
        {{ ageRange }}
      </span>
    </div>

    <article
      v-for="row in rows"
      :key="row.label"
      class="matrix__row"
    >
      <strong class="matrix__label">{{ row.label }}</strong>

      <div
        v-for="cell in row.values"
        :key="`${row.label}-${cell.ageRange}`"
        class="matrix__cell"
        :style="{ background: `rgba(14, 165, 233, ${intensity(cell.listeners)})` }"
      >
        <span>{{ cell.listeners }}</span>
      </div>
    </article>
  </div>

  <div v-else class="matrix__empty">
    {{ emptyText }}
  </div>
</template>

<style scoped>
.matrix {
  display: grid;
  gap: 0.75rem;
}

.matrix__header,
.matrix__row {
  display: grid;
  grid-template-columns: minmax(220px, 1.8fr) repeat(auto-fit, minmax(84px, 1fr));
  gap: 0.65rem;
  align-items: stretch;
}

.matrix__header {
  color: #64748b;
  font-size: 0.8rem;
  text-transform: uppercase;
  letter-spacing: 0.06em;
}

.matrix__header-cell {
  text-align: center;
}

.matrix__label {
  display: flex;
  align-items: center;
  padding: 0.85rem 1rem;
  border-radius: 1rem;
  background: rgba(255, 255, 255, 0.88);
  border: 1px solid rgba(15, 23, 42, 0.06);
  color: #0f172a;
}

.matrix__cell {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 3.2rem;
  border-radius: 1rem;
  border: 1px solid rgba(14, 116, 144, 0.08);
  color: #0f172a;
  font-weight: 700;
}

.matrix__empty {
  color: #64748b;
}

@media (max-width: 900px) {
  .matrix__header,
  .matrix__row {
    grid-template-columns: 1.5fr repeat(auto-fit, minmax(70px, 1fr));
  }
}

@media (max-width: 640px) {
  .matrix {
    overflow-x: auto;
  }

  .matrix__header,
  .matrix__row {
    min-width: 560px;
  }
}
</style>
