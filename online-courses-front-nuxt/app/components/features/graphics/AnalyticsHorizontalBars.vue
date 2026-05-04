<script setup lang="ts">
const props = withDefaults(defineProps<{
  items: Array<{
    label: string
    value: number
    hint?: string
  }>
  colorFrom?: string
  colorTo?: string
  valueSuffix?: string
  emptyText?: string
}>(), {
  colorFrom: '#0f172a',
  colorTo: '#38bdf8',
  valueSuffix: '',
  emptyText: 'Нет данных для отображения.'
})

const maxValue = computed(() =>
  props.items.length ? Math.max(...props.items.map(item => item.value), 1) : 1
)

const formatValue = (value: number) =>
  `${new Intl.NumberFormat('ru-RU').format(value)}${props.valueSuffix}`
</script>

<template>
  <div v-if="items.length" class="bars">
    <article
      v-for="(item, index) in items"
      :key="`${item.label}-${index}`"
      class="bars__item"
    >
      <div class="bars__topline">
        <div class="bars__labels">
          <span class="bars__rank">#{{ index + 1 }}</span>
          <strong>{{ item.label }}</strong>
        </div>

        <div class="bars__meta">
          <strong>{{ formatValue(item.value) }}</strong>
          <span v-if="item.hint">{{ item.hint }}</span>
        </div>
      </div>

      <div class="bars__track">
        <div
          class="bars__fill"
          :style="{
            width: `${(item.value / maxValue) * 100}%`,
            background: `linear-gradient(90deg, ${colorFrom}, ${colorTo})`
          }"
        />
      </div>
    </article>
  </div>

  <div v-else class="bars__empty">
    {{ emptyText }}
  </div>
</template>

<style scoped>
.bars {
  display: grid;
  gap: 0.9rem;
}

.bars__item {
  display: grid;
  gap: 0.55rem;
  padding: 1rem;
  border-radius: 1.1rem;
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.94), rgba(241, 245, 249, 0.82));
  border: 1px solid rgba(15, 23, 42, 0.07);
}

.bars__topline {
  display: flex;
  justify-content: space-between;
  gap: 1rem;
  align-items: flex-start;
}

.bars__labels,
.bars__meta {
  display: grid;
  gap: 0.18rem;
  min-width: 0;
}

.bars__labels strong,
.bars__meta strong {
  color: #0f172a;
  line-height: 1.35;
  overflow-wrap: anywhere;
}

.bars__meta {
  flex: 0 0 min(11rem, 100%);
  text-align: right;
  justify-items: end;
}

.bars__meta span,
.bars__rank {
  font-size: 0.8rem;
  color: #64748b;
}

.bars__track {
  height: 0.82rem;
  border-radius: 999px;
  background: rgba(148, 163, 184, 0.18);
  overflow: hidden;
}

.bars__fill {
  height: 100%;
  min-width: 0.9rem;
  border-radius: inherit;
  box-shadow: 0 12px 28px rgba(14, 116, 144, 0.28);
}

.bars__empty {
  color: #64748b;
}

@media (max-width: 720px) {
  .bars__topline {
    flex-direction: column;
  }

  .bars__meta {
    text-align: left;
    justify-items: start;
  }
}
</style>
