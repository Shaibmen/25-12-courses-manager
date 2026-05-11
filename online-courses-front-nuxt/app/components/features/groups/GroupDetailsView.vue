<script setup lang="ts">
import type { GroupItem } from '../../../types/group'
import { formatApiDate, getDateSortValue } from '../../../utils/date'
import AppButton from '../../ui/AppButton.vue'
import AppCard from '../../ui/AppCard.vue'

const props = defineProps<{
  group: GroupItem
  exportLoading?: boolean
}>()

defineEmits<{
  back: []
  edit: []
  export: []
}>()

const schedule = computed(() =>
  [...(props.group.rapspisanie || [])].sort(
    (left, right) => getDateSortValue(left.date) - getDateSortValue(right.date)
  )
)
</script>

<template>
  <section class="stack">
    <AppCard title="Карточка группы">
      <div class="group-hero">
        <div class="group-hero__content">
          <span class="group-hero__label">Учебная группа</span>
          <h2>{{ group.name_group }}</h2>
          <p>Всего занятий в расписании: {{ schedule.length }}</p>
        </div>

        <div class="group-hero__actions">
          <AppButton variant="ghost" @click="$emit('back')">Назад</AppButton>
          <AppButton variant="secondary" :disabled="props.exportLoading" @click="$emit('export')">
            {{ props.exportLoading ? 'Скачиваем...' : 'Экспорт' }}
          </AppButton>
          <AppButton @click="$emit('edit')">Редактировать</AppButton>
        </div>
      </div>
    </AppCard>

    <AppCard title="Расписание занятий">
      <p v-if="!schedule.length" class="empty-state">
        В этой группе пока нет занятий.
      </p>

      <div v-else class="schedule-grid">
        <article
          v-for="(item, index) in schedule"
          :key="`${item.date}-${item.theme}-${index}`"
          class="schedule-card"
        >
          <span class="schedule-card__index">Занятие {{ index + 1 }}</span>
          <strong>{{ item.theme }}</strong>
          <span>{{ formatApiDate(item.date) }}</span>
        </article>
      </div>
    </AppCard>
  </section>
</template>

<style scoped>
.group-hero {
  display: flex;
  justify-content: space-between;
  gap: 1rem;
  align-items: end;
  flex-wrap: wrap;
}

.group-hero__content {
  display: grid;
  gap: 0.5rem;
}

.group-hero__content h2,
.group-hero__content p {
  margin: 0;
}

.group-hero__label {
  font-size: 0.78rem;
  letter-spacing: 0.12em;
  text-transform: uppercase;
  color: #64748b;
}

.group-hero__actions {
  display: flex;
  gap: 1rem;
  flex-wrap: wrap;
}

.schedule-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 1rem;
}

.schedule-card {
  display: grid;
  gap: 0.4rem;
  min-width: 0;
  padding: 1rem;
  border-radius: 1.15rem;
  background: rgba(239, 246, 255, 0.82);
  border: 1px solid rgba(59, 130, 246, 0.12);
}

.schedule-card strong,
.schedule-card span {
  min-width: 0;
  overflow-wrap: anywhere;
  word-break: break-word;
}

.schedule-card__index,
.schedule-card span {
  color: #475569;
}

.empty-state {
  margin: 0;
  color: #64748b;
}

@media (max-width: 720px) {
  .group-hero__actions {
    width: 100%;
  }
}
</style>
