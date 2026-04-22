<script setup lang="ts">
import type { DivisionItem, EducationTypeItem } from '../../types/catalogs'
import type { ProgramListItem } from '../../types/program'
import AppButton from '../../components/ui/AppButton.vue'
import AppCard from '../../components/ui/AppCard.vue'
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
const errorMessage = ref('')
const sortField = ref<'individual' | 'group' | 'campus' | null>(null)
const sortOrder = ref<'asc' | 'desc' | null>(null)
let debounceTimer: ReturnType<typeof setTimeout> | null = null

const priceFieldMap = {
  individual: 'individual_price',
  group: 'group_price',
  campus: 'campus_price'
} as const

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

const sortedPrograms = computed(() => {
  const list = [...programs.value]

  if (!sortField.value || !sortOrder.value) {
    return list
  }

  const targetField = priceFieldMap[sortField.value]

  return list.sort((left, right) => {
    const leftValue = Number(left[targetField] || 0)
    const rightValue = Number(right[targetField] || 0)

    return sortOrder.value === 'asc' ? leftValue - rightValue : rightValue - leftValue
  })
})

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

const removeProgram = async (id: string) => {
  if (!window.confirm('Вы точно хотите удалить эту программу?')) {
    return
  }

  try {
    await deleteProgramRequest(id)
    await loadPrograms()
    notifications.success('Программа удалена. Список обновлён.', 'Программы обучения')
  } catch (error) {
    notifications.error(
      error instanceof Error ? error.message : 'Не удалось удалить программу обучения',
      'Программы обучения'
    )
  }
}

const toggleSort = (field: 'individual' | 'group' | 'campus') => {
  if (sortField.value !== field) {
    sortField.value = field
    sortOrder.value = 'asc'
    return
  }

  sortOrder.value = sortOrder.value === 'asc' ? 'desc' : sortOrder.value === 'desc' ? null : 'asc'

  if (!sortOrder.value) {
    sortField.value = null
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

onMounted(() => {
  void loadPrograms()
})
</script>

<template>
  <section class="stack">
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

      <div class="toolbar toolbar--sort">
        <AppButton variant="ghost" @click="toggleSort('individual')">
          Индивидуальная цена {{ sortField === 'individual' && sortOrder ? (sortOrder === 'asc' ? '↑' : '↓') : '' }}
        </AppButton>
        <AppButton variant="ghost" @click="toggleSort('group')">
          Групповая цена {{ sortField === 'group' && sortOrder ? (sortOrder === 'asc' ? '↑' : '↓') : '' }}
        </AppButton>
        <AppButton variant="ghost" @click="toggleSort('campus')">
          Кампусная цена {{ sortField === 'campus' && sortOrder ? (sortOrder === 'asc' ? '↑' : '↓') : '' }}
        </AppButton>
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
              <th>Цена (₽)</th>
              <th>Индивидуально (₽)</th>
              <th>Группа (₽)</th>
              <th>Самообучение (₽)</th>
              <th>Тип обучения</th>
              <th>Подразделение</th>
              <th>Действия</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="!sortedPrograms.length">
              <td colspan="8" class="programs-table__empty">Программы пока не найдены.</td>
            </tr>
            <tr v-for="program in sortedPrograms" :key="program.id_program_education">
              <td class="programs-table__name">{{ program.name_prof_education }}</td>
              <td>{{ program.time_education }}</td>
              <td>{{ formatPrice(program.price) }}</td>
              <td>{{ formatPrice(program.individual_price) }}</td>
              <td>{{ formatPrice(program.group_price) }}</td>
              <td>{{ formatPrice(program.campus_price) }}</td>
              <td>{{ getTypeName(program.id_education_type) }}</td>
              <td>{{ getDivisionName(program.id_divisions_education) }}</td>
              <td>
                <div class="programs-table__actions">
                  <AppButton variant="secondary" @click="router.push(`/programs/edit/${program.id_program_education}`)">
                    Изменить
                  </AppButton>
                  <AppButton variant="ghost" @click="removeProgram(program.id_program_education)">
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

.toolbar--sort {
  margin-top: 1rem;
  justify-content: flex-start;
  align-items: stretch;
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
  min-width: 1100px;
  border-collapse: collapse;
}

.programs-table th,
.programs-table td {
  padding: 0.95rem 0.85rem;
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
  min-width: 18rem;
  max-width: 24rem;
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
