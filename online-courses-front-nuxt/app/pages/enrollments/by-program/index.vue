<script setup lang="ts">
import type { EnrollmentItem } from '../../../types/enrollment'
import type { GroupItem } from '../../../types/group'
import type { ProgramListItem } from '../../../types/program'
import { formatApiDate } from '../../../utils/date'
import AppButton from '../../../components/ui/AppButton.vue'
import AppCard from '../../../components/ui/AppCard.vue'
import AppSelect from '../../../components/ui/AppSelect.vue'

const router = useRouter()
const notifications = useNotifications()
const programs = ref<ProgramListItem[]>([])
const enrollments = ref<EnrollmentItem[]>([])
const groups = ref<GroupItem[]>([])
const selectedProgramId = ref('')
const programPage = ref(1)
const page = ref(1)
const hasMore = ref(false)
const loading = ref(false)
const errorMessage = ref('')

const groupNameById = computed(() =>
  new Map(groups.value.map((group) => [group.group, group.name_group]))
)

const resolveGroupName = (id: string) => groupNameById.value.get(id) || id || '—'

const loadPrograms = async () => {
  try {
    const [programData, groupData] = await Promise.all([
      getPrograms(programPage.value, ''),
      getGroups(1, '')
    ])

    programs.value = programData
    groups.value = groupData

    if (!selectedProgramId.value && programData.length) {
      selectedProgramId.value = programData[0].id_program_education
    }
  } catch (error) {
    notifications.error(
      error instanceof Error ? error.message : 'Не удалось загрузить программы',
      'Поиск по программе'
    )
  }
}

const loadEnrollments = async () => {
  if (!selectedProgramId.value) {
    enrollments.value = []
    return
  }

  loading.value = true
  errorMessage.value = ''

  try {
    const data = await getEnrollmentByProgram(selectedProgramId.value, page.value)
    enrollments.value = data
    hasMore.value = data.length === 25
  } catch (error) {
    errorMessage.value =
      error instanceof Error ? error.message : 'Не удалось загрузить записи по программе'
  } finally {
    loading.value = false
  }
}

watch(selectedProgramId, () => {
  page.value = 1
  void loadEnrollments()
})

onMounted(async () => {
  await loadPrograms()
  await loadEnrollments()
})
</script>

<template>
  <section class="stack content-shell">
    <AppCard title="Поиск записей по программе">
      <div class="toolbar">
        <div class="toolbar__program">
          <AppSelect v-model="selectedProgramId" label="Программа" placeholder="Выберите программу">
            <option v-for="program in programs" :key="program.id_program_education" :value="program.id_program_education">
              {{ program.name_prof_education }}
            </option>
          </AppSelect>

          <div class="program-pagination">
            <AppButton variant="ghost" :disabled="programPage <= 1" @click="programPage--, loadPrograms()">
              Предыдущая страница
            </AppButton>
            <span>Страница программ {{ programPage }}</span>
            <AppButton variant="ghost" @click="programPage++, loadPrograms()">
              Следующая страница
            </AppButton>
          </div>
        </div>

        <div class="toolbar__actions">
          <AppButton variant="secondary" to="/enrollments">Назад</AppButton>
        </div>
      </div>
    </AppCard>

    <AppCard v-if="loading" title="Загрузка">
      <p>Подтягиваю записи по выбранной программе.</p>
    </AppCard>

    <AppCard v-else-if="errorMessage" title="Ошибка">
      <p>{{ errorMessage }}</p>
      <AppButton variant="secondary" @click="loadEnrollments">Повторить запрос</AppButton>
    </AppCard>

    <AppCard v-else-if="!selectedProgramId" title="Выберите программу">
      <p>Сначала выберите программу из списка выше.</p>
    </AppCard>

    <AppCard v-else-if="!enrollments.length" title="Нет записей">
      <p>На этой программе слушателей пока нет.</p>
    </AppCard>

    <AppCard v-else title="Записи по программе">
      <div class="table-wrap">
        <table class="enrollments-table">
          <thead>
            <tr>
              <th>Фамилия</th>
              <th>Имя</th>
              <th>Курс</th>
              <th>Начало</th>
              <th>Окончание</th>
              <th>Группа</th>
              <th>Тип обучения</th>
              <th>Действия</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in enrollments" :key="item.id_listener + item.name_prof_education">
              <td>{{ item.second_name }}</td>
              <td>{{ item.first_name }}</td>
              <td class="enrollments-table__name">{{ item.name_prof_education }}</td>
              <td>{{ formatApiDate(item.start_date) }}</td>
              <td>{{ formatApiDate(item.end_date) }}</td>
              <td>{{ resolveGroupName(item.id_group) }}</td>
              <td>{{ item.type_of_retraining || '—' }}</td>
              <td>
                <AppButton variant="secondary" @click="router.push(`/enrollment/details/${item.id_listener}`)">
                  Подробнее
                </AppButton>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="pagination">
        <AppButton variant="secondary" :disabled="page <= 1" @click="page > 1 && (page--, loadEnrollments())">
          Назад
        </AppButton>
        <span>Страница {{ page }}</span>
        <AppButton variant="secondary" :disabled="!hasMore" @click="page++, loadEnrollments()">
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

.toolbar__program {
  display: grid;
  gap: 1rem;
  flex: 1;
}

.toolbar__actions {
  display: flex;
  gap: 0.75rem;
}

.program-pagination {
  display: flex;
  align-items: center;
  gap: 1rem;
  flex-wrap: wrap;
}

.table-wrap {
  overflow-x: auto;
}

.enrollments-table {
  width: 100%;
  min-width: 1120px;
  border-collapse: collapse;
}

.enrollments-table th,
.enrollments-table td {
  padding: 1rem 0.95rem;
  border-bottom: 1px solid rgba(15, 23, 42, 0.08);
  text-align: left;
  vertical-align: middle;
}

.enrollments-table th {
  color: #475569;
  font-size: 0.82rem;
  text-transform: uppercase;
  letter-spacing: 0.06em;
}

.enrollments-table tbody tr:hover {
  background: rgba(219, 234, 254, 0.26);
}

.enrollments-table__name {
  min-width: 22rem;
  max-width: 30rem;
  overflow-wrap: anywhere;
  word-break: break-word;
}

.pagination {
  margin-top: 1rem;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 1rem;
}
</style>
