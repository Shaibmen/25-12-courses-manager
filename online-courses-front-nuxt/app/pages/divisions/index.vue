<script setup lang="ts">
import AppButton from '../../components/ui/AppButton.vue'
import AppCard from '../../components/ui/AppCard.vue'
import AppInput from '../../components/ui/AppInput.vue'
import type { DivisionItem } from '../../types/catalogs'

const router = useRouter()
const notifications = useNotifications()
const filter = ref('')
const divisions = ref<DivisionItem[]>([])
const loading = ref(false)
const errorMessage = ref('')
let debounceTimer: ReturnType<typeof setTimeout> | null = null

const loadDivisions = async () => {
  loading.value = true
  errorMessage.value = ''

  try {
    divisions.value = await getDivisions(filter.value)
  } catch (error) {
    errorMessage.value =
      error instanceof Error ? error.message : 'Не удалось загрузить подразделения'
  } finally {
    loading.value = false
  }
}

const removeDivision = async (id: string) => {
  if (!window.confirm('Вы точно хотите удалить это подразделение?')) {
    return
  }

  try {
    await deleteDivisionRequest(id)
    await loadDivisions()
    notifications.success('Подразделение удалено. Список обновлён.', 'Подразделения')
  } catch (error) {
    notifications.error('Не удалось удалить подразделение, возможно оно где то используется')
  }
}

watch(filter, () => {
  if (debounceTimer) {
    clearTimeout(debounceTimer)
  }

  debounceTimer = setTimeout(() => {
    void loadDivisions()
  }, 350)
})

onMounted(() => {
  void loadDivisions()
})
</script>

<template>
  <section class="stack">
    <AppCard title="Подразделения обучения">
      <div class="toolbar">
        <div class="toolbar__search">
          <AppInput v-model="filter" label="Поиск" placeholder="Поиск по названию" />
        </div>

        <div class="toolbar__actions">
          <AppButton to="/divisions/create">Добавить подразделение</AppButton>
          <AppButton variant="secondary" to="/dashboard">Назад</AppButton>
        </div>
      </div>
    </AppCard>

    <AppCard v-if="loading" title="Загрузка">
      <p>Подтягиваю подразделения обучения.</p>
    </AppCard>

    <AppCard v-else-if="errorMessage" title="Ошибка">
      <p>{{ errorMessage }}</p>
      <AppButton variant="secondary" @click="loadDivisions">Повторить запрос</AppButton>
    </AppCard>

    <AppCard v-else title="Справочник">
      <div class="table-wrap">
        <table class="catalog-table">
          <thead>
            <tr>
              <th>Название подразделения</th>
              <th>Действия</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="!divisions.length">
              <td colspan="2" class="catalog-table__empty">Подразделения пока не найдены.</td>
            </tr>
            <tr v-for="division in divisions" :key="division.id_divisionsEducation">
              <td>{{ division.divisions }}</td>
              <td>
                <div class="catalog-table__actions">
                  <AppButton variant="secondary" @click="router.push(`/divisions/edit/${division.id_divisionsEducation}`)">
                    Изменить
                  </AppButton>
                  <AppButton variant="ghost" @click="removeDivision(division.id_divisionsEducation)">
                    Удалить
                  </AppButton>
                </div>
              </td>
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

.toolbar__actions,
.catalog-table__actions {
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
