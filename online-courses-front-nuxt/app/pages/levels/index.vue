<script setup lang="ts">
import AppButton from '../../components/ui/AppButton.vue'
import AppCard from '../../components/ui/AppCard.vue'
import AppInput from '../../components/ui/AppInput.vue'
import type { LevelEducationItem } from '../../types/listener'

const filter = ref('')
const levels = ref<LevelEducationItem[]>([])
const loading = ref(false)
const errorMessage = ref('')
let debounceTimer: ReturnType<typeof setTimeout> | null = null

const loadLevels = async () => {
  loading.value = true
  errorMessage.value = ''

  try {
    levels.value = await getEducationLevelsCatalog(filter.value)
  } catch (error) {
    errorMessage.value =
      error instanceof Error ? error.message : 'Не удалось загрузить уровни обучения'
  } finally {
    loading.value = false
  }
}

watch(filter, () => {
  if (debounceTimer) {
    clearTimeout(debounceTimer)
  }

  debounceTimer = setTimeout(() => {
    void loadLevels()
  }, 350)
})

onMounted(() => {
  void loadLevels()
})
</script>

<template>
  <section class="stack">
    <AppCard title="Уровни обучения">
      <div class="toolbar">
        <div class="toolbar__search">
          <AppInput v-model="filter" label="Поиск" placeholder="Поиск по названию уровня" />
        </div>

        <div class="toolbar__actions">
          <AppButton variant="secondary" to="/dashboard">Назад</AppButton>
        </div>
      </div>
    </AppCard>

    <AppCard v-if="loading" title="Загрузка">
      <p>Подтягиваю уровни обучения.</p>
    </AppCard>

    <AppCard v-else-if="errorMessage" title="Ошибка">
      <p>{{ errorMessage }}</p>
      <AppButton variant="secondary" @click="loadLevels">Повторить запрос</AppButton>
    </AppCard>

    <AppCard v-else title="Справочник">
      <div class="table-wrap">
        <table class="catalog-table">
          <thead>
            <tr>
              <th>Название уровня образования</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="!levels.length">
              <td class="catalog-table__empty">Нет данных</td>
            </tr>
            <tr v-for="level in levels" :key="level.id_level_education">
              <td>{{ level.education }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </AppCard>
  </section>
</template>

<style scoped>
.toolbar {
  display: flex;
  justify-content: space-between;
  gap: 1rem;
  align-items: end;
  flex-wrap: wrap;
}

.toolbar__search {
  min-width: min(100%, 24rem);
  flex: 1 1 24rem;
}

.toolbar__actions {
  display: flex;
  gap: 0.75rem;
  flex-wrap: wrap;
}

.table-wrap {
  overflow-x: auto;
}

.catalog-table {
  width: 100%;
  border-collapse: collapse;
}

.catalog-table th,
.catalog-table td {
  padding: 0.95rem 0.85rem;
  border-bottom: 1px solid rgba(15, 23, 42, 0.08);
  text-align: left;
}

.catalog-table__empty {
  text-align: center;
  color: #64748b;
}
</style>
