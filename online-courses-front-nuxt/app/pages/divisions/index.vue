<script setup lang="ts">
import AppButton from '../../components/ui/AppButton.vue'
import AppCard from '../../components/ui/AppCard.vue'
import AppConfirmDialog from '../../components/ui/AppConfirmDialog.vue'
import AppInput from '../../components/ui/AppInput.vue'
import type { DivisionItem } from '../../types/catalogs'

const router = useRouter()
const notifications = useNotifications()
const filter = ref('')
const divisions = ref<DivisionItem[]>([])
const loading = ref(false)
const deleteLoading = ref(false)
const errorMessage = ref('')
const divisionToDelete = ref<DivisionItem | null>(null)
let debounceTimer: ReturnType<typeof setTimeout> | null = null

const loadDivisions = async () => {
  loading.value = true
  errorMessage.value = ''

  try {
    divisions.value = await getDivisions(filter.value)
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : 'Не удалось загрузить подразделения'
  } finally {
    loading.value = false
  }
}

const confirmDelete = async () => {
  if (!divisionToDelete.value) {
    return
  }

  deleteLoading.value = true

  try {
    await deleteDivisionRequest(divisionToDelete.value.id_divisionsEducation)
    divisionToDelete.value = null
    await loadDivisions()
    notifications.success('Подразделение удалено. Список обновлён.', 'Подразделения')
  } catch {
    notifications.error('Не удалось удалить подразделение. Возможно, оно используется в системе.', 'Подразделения')
  } finally {
    deleteLoading.value = false
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

onBeforeUnmount(() => {
  if (debounceTimer) {
    clearTimeout(debounceTimer)
  }
})

onMounted(() => {
  void loadDivisions()
})
</script>

<template>
  <section class="stack content-shell">
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
                  <AppButton variant="ghost" @click="divisionToDelete = division">
                    Удалить
                  </AppButton>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </AppCard>

    <AppConfirmDialog
      :open="Boolean(divisionToDelete)"
      title="Удаление подразделения"
      :message="divisionToDelete ? `Удалить подразделение «${divisionToDelete.divisions}»?` : ''"
      :loading="deleteLoading"
      confirm-label="Удалить"
      @cancel="divisionToDelete = null"
      @confirm="confirmDelete"
    />
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
  min-width: 760px;
  border-collapse: collapse;
}

.catalog-table th,
.catalog-table td {
  padding: 1rem 0.95rem;
  border-bottom: 1px solid rgba(15, 23, 42, 0.08);
  text-align: left;
}

.catalog-table__empty {
  text-align: center;
  color: #64748b;
}
</style>
