<script setup lang="ts">
import type { UserDashboardResponse } from '../../../types/dashboard'
import AppButton from '../../ui/AppButton.vue'
import AppCard from '../../ui/AppCard.vue'
import AppStatCard from '../../ui/AppStatCard.vue'

defineProps<{
  data: UserDashboardResponse
}>()

const navigationItems = [
  { title: 'Слушатели', description: 'Карточки слушателей и их история обучения.', to: '/listeners' },
  { title: 'Подразделения обучения', description: 'Справочник подразделений для программ.', to: '/divisions' },
  { title: 'Типы обучения', description: 'Классификация образовательных программ.', to: '/types' },
  { title: 'Уровни обучения', description: 'Уровни и формат образовательной подготовки.', to: '/levels' },
  { title: 'Программы обучения', description: 'Каталог программ и их наполнение.', to: '/programs' },
  { title: 'Запись на курс', description: 'Управление зачислениями и текущими наборами.', to: '/enrollments' },
  { title: 'Юридические лица', description: 'Организации, договоры и привязка слушателей.', to: '/legalentities' },
  { title: 'Исполнители', description: 'Преподаватели и подрядчики по программам.', to: '/executers' }
]
</script>

<template>
  <section class="dashboard-stack">
    <AppCard title="Разделы системы">
      <div class="module-grid">
        <article
          v-for="item in navigationItems"
          :key="item.to"
          class="module-card"
        >
          <div class="module-card__content">
            <strong>{{ item.title }}</strong>
            <p>{{ item.description }}</p>
          </div>

          <AppButton :to="item.to" variant="ghost">
            Открыть раздел
          </AppButton>
        </article>
      </div>
    </AppCard>

    <div class="stats-grid">
      <AppStatCard title="Всего слушателей" :value="data.total_listener" />
      <AppStatCard title="Всего программ" :value="data.total_program" />
      <AppStatCard
        title="Активные зачисления"
        :value="data.active_enrollments"
        description="Актуальная нагрузка по действующим наборам."
      />
    </div>

    <AppCard title="Курсы с ближайшим окончанием">
      <div v-if="!data.program_ending_soon.length" class="empty-state">
        Пока нет записей, которые скоро заканчиваются.
      </div>

      <div v-else class="ending-grid">
        <article
          v-for="program in data.program_ending_soon"
          :key="`${program.name_prof_education}-${program.end_date}`"
          class="ending-card"
        >
          <strong>{{ program.name_prof_education }}</strong>
          <span>Дата окончания: {{ new Date(program.end_date).toLocaleDateString('ru-RU') }}</span>
          <span>Слушателей: {{ program.total_listeners }}</span>
        </article>
      </div>
    </AppCard>
  </section>
</template>

<style scoped>
.dashboard-stack {
  display: grid;
  gap: 1rem;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 1rem;
}

.module-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 1rem;
}

.module-card {
  display: grid;
  gap: 1rem;
  align-content: space-between;
  padding: 1rem;
  border-radius: 1.1rem;
  background: rgba(255, 255, 255, 0.7);
  border: 1px solid rgba(15, 23, 42, 0.08);
}

.module-card__content {
  display: grid;
  gap: 0.45rem;
}

.module-card__content p {
  margin: 0;
  color: #475569;
}

.ending-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 0.9rem;
}

.ending-card {
  display: grid;
  gap: 0.45rem;
  padding: 1rem;
  border-radius: 1.1rem;
  background: rgba(255, 255, 255, 0.85);
  border: 1px solid rgba(15, 23, 42, 0.07);
}

.ending-card span {
  color: #475569;
}

.empty-state {
  color: #64748b;
}

@media (max-width: 900px) {
  .module-grid {
    grid-template-columns: 1fr;
  }

  .stats-grid {
    grid-template-columns: 1fr;
  }
}
</style>
