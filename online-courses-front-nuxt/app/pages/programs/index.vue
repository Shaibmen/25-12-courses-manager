<script setup lang="ts">
import type { DivisionItem, EducationTypeItem } from '../../types/catalogs'
import type { ProgramListItem } from '../../types/program'
import AppButton from '../../components/ui/AppButton.vue'
import AppCard from '../../components/ui/AppCard.vue'
import AppConfirmDialog from '../../components/ui/AppConfirmDialog.vue'
import AppInput from '../../components/ui/AppInput.vue'

const router = useRouter()
const notifications = useNotifications()
const programs = ref<ProgramListItem[]>([])
const educationTypes = ref<EducationTypeItem[]>([])
const divisions = ref<DivisionItem[]>([])
const page = ref(1)
const filter = ref('')
const hasMore = ref(false)
const loading = ref(false)
const deleteLoading = ref(false)
const errorMessage = ref('')
const programToDelete = ref<ProgramListItem | null>(null)
let debounceTimer: ReturnType<typeof setTimeout> | null = null

const formatPrice = (value?: number | null) => {
  if (value === undefined || value === null) {
    return '—'
  }

  return `${new Intl.NumberFormat('ru-RU').format(value)} ₽`
}

const getTypeName = (id: string) =>
  educationTypes.value.find((item) => item.id_educationType === id)?.typeName || '—'

const getDivisionName = (id: string) =>
  divisions.value.find((item) => item.id_divisionsEducation === id)?.divisions || '—'

const loadPrograms = async () => {
  loading.value = true
  errorMessage.value = ''

  try {
    const [programData, typeData, divisionData] = await Promise.all([
      getPrograms(page.value, filter.value),
      getEducationTypes(''),
      getDivisions('')
    ])

    programs.value = programData
    educationTypes.value = typeData
    divisions.value = divisionData
    hasMore.value = programData.length === 25
  } catch (error) {
    errorMessage.value =
      error instanceof Error ? error.message : 'Не удалось загрузить программы обучения'
  } finally {
    loading.value = false
  }
}

const confirmDelete = async () => {
  if (!programToDelete.value) {
    return
  }

  deleteLoading.value = true

  try {
    await deleteProgramRequest(programToDelete.value.id_program_education)
    programToDelete.value = null
    await loadPrograms()
    notifications.success('Программа удалена. Список обновлён.', 'Программы обучения')
  } catch (error) {
    notifications.error(
      error instanceof Error ? error.message : 'Не удалось удалить программу обучения',
      'Программы обучения'
    )
  } finally {
    deleteLoading.value = false
  }
}

watch(filter, () => {
  if (debounceTimer) {
    clearTimeout(debounceTimer)
  }

  debounceTimer = setTimeout(() => {
    page.value = 1
    void loadPrograms()
  }, 350)
})

onBeforeUnmount(() => {
  if (debounceTimer) {
    clearTimeout(debounceTimer)
  }
})

onMounted(() => {
  void loadPrograms()
})
</script>

<template>
  <section class="stack content-shell">
    <AppCard title="Программы обучения">
      <div class="toolbar">
        <div class="toolbar__search">
          <AppInput v-model="filter" label="Поиск" placeholder="Поиск по названию программы" />
        </div>

        <div class="toolbar__actions">
          <AppButton to="/programs/create">Добавить программу</AppButton>
          <AppButton variant="secondary" to="/dashboard">Назад</AppButton>
        </div>
      </div>
    </AppCard>

    <AppCard v-if="loading" title="Загрузка">
      <p>Подтягиваю программы и связанные справочники.</p>
    </AppCard>

    <AppCard v-else-if="errorMessage" title="Ошибка">
      <p>{{ errorMessage }}</p>
      <AppButton variant="secondary" @click="loadPrograms">Повторить запрос</AppButton>
    </AppCard>

    <AppCard v-else title="Каталог программ">
      <div class="table-wrap">
        <table class="programs-table">
          <thead>
            <tr>
              <th>Название программы</th>
              <th>Длительность (часы)</th>
              <th>Цена</th>
              <th>Тип обучения</th>
              <th>Подразделение</th>
              <th>Действия</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="!programs.length">
              <td colspan="6" class="programs-table__empty">Программы пока не найдены.</td>
            </tr>
            <tr v-for="program in programs" :key="program.id_program_education">
              <td class="programs-table__name">{{ program.name_prof_education }}</td>
              <td>{{ program.time_education }}</td>
              <td>{{ formatPrice(program.price) }}</td>
              <td>{{ getTypeName(program.id_education_type) }}</td>
              <td>{{ getDivisionName(program.id_divisions_education) }}</td>
              <td>
                <div class="programs-table__actions">
                  <AppButton variant="secondary" @click="router.push(`/programs/edit/${program.id_program_education}`)">
                    Изменить
                  </AppButton>
                  <AppButton variant="ghost" @click="programToDelete = program">
                    Удалить
                  </AppButton>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="pagination">
        <AppButton variant="secondary" :disabled="page <= 1" @click="page > 1 && (page--, loadPrograms())">
          Назад
        </AppButton>
        <span>Страница {{ page }}</span>
        <AppButton variant="secondary" :disabled="!hasMore" @click="page++, loadPrograms()">
          Вперёд
        </AppButton>
      </div>
    </AppCard>

    <AppConfirmDialog
      :open="Boolean(programToDelete)"
      title="Удаление программы"
      :message="programToDelete ? `Удалить программу «${programToDelete.name_prof_education}»?` : ''"
      :loading="deleteLoading"
      confirm-label="Удалить"
      @cancel="programToDelete = null"
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
.programs-table__actions {
  display: flex;
  gap: 0.75rem;
  flex-wrap: wrap;
}

.table-wrap {
  overflow-x: auto;
}

.programs-table {
  width: 100%;
  min-width: 1040px;
  border-collapse: collapse;
}

.programs-table th,
.programs-table td {
  padding: 1rem 0.95rem;
  border-bottom: 1px solid rgba(15, 23, 42, 0.08);
  text-align: left;
  vertical-align: middle;
}

.programs-table th {
  color: #475569;
  font-size: 0.82rem;
  text-transform: uppercase;
  letter-spacing: 0.06em;
}

.programs-table tbody tr:hover {
  background: rgba(219, 234, 254, 0.26);
}

.programs-table__name {
  min-width: 22rem;
  max-width: 30rem;
  overflow-wrap: anywhere;
  word-break: break-word;
}

.programs-table__empty {
  text-align: center;
  color: #64748b;
}

.pagination {
  margin-top: 1rem;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 1rem;
}
</style>
