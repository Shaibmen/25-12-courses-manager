<script setup lang="ts">
import type { EnrollmentItem } from '../../types/enrollment'
import type { GroupItem } from '../../types/group'
import { formatApiDate } from '../../utils/date'
import AppButton from '../../components/ui/AppButton.vue'
import AppCard from '../../components/ui/AppCard.vue'
import AppInput from '../../components/ui/AppInput.vue'

const router = useRouter()
const enrollments = ref<EnrollmentItem[]>([])
const groups = ref<GroupItem[]>([])
const page = ref(1)
const filter = ref('')
const hasMore = ref(false)
const loading = ref(false)
const errorMessage = ref('')
let debounceTimer: ReturnType<typeof setTimeout> | null = null

const groupNameById = computed(() =>
  new Map(groups.value.map((group) => [group.group, group.name_group]))
)

const resolveGroupName = (id: string) => groupNameById.value.get(id) || id || '—'

const loadEnrollments = async () => {
  loading.value = true
  errorMessage.value = ''

  try {
    const [data, groupData] = await Promise.all([
      getEnrollments(page.value, filter.value),
      getGroups(1, '')
    ])

    enrollments.value = data
    groups.value = groupData
    hasMore.value = data.length === 25
  } catch (error) {
    errorMessage.value =
      error instanceof Error ? error.message : 'Не удалось загрузить записи на курсы'
  } finally {
    loading.value = false
  }
}

watch(filter, () => {
  if (debounceTimer) {
    clearTimeout(debounceTimer)
  }

  debounceTimer = setTimeout(() => {
    page.value = 1
    void loadEnrollments()
  }, 350)
})

onBeforeUnmount(() => {
  if (debounceTimer) {
    clearTimeout(debounceTimer)
  }
})

onMounted(() => {
  void loadEnrollments()
})
</script>

<template>
  <section class="stack content-shell">
    <AppCard title="Записи на курсы">
      <div class="toolbar">
        <div class="toolbar__search">
          <AppInput v-model="filter" label="Поиск" placeholder="Поиск по фамилии" />
        </div>

        <div class="toolbar__actions">
          <AppButton to="/enrollments/by-program">Поиск по программе</AppButton>
          <AppButton to="/enrollments/accurate">Точные записи</AppButton>
          <AppButton variant="secondary" to="/dashboard">Назад</AppButton>
        </div>
      </div>
    </AppCard>

    <AppCard v-if="loading" title="Загрузка">
      <p>Подтягиваю записи на курсы.</p>
    </AppCard>

    <AppCard v-else-if="errorMessage" title="Ошибка">
      <p>{{ errorMessage }}</p>
      <AppButton variant="secondary" @click="loadEnrollments">Повторить запрос</AppButton>
    </AppCard>

    <AppCard v-else title="Каталог записей">
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
            <tr v-if="!enrollments.length">
              <td colspan="8" class="enrollments-table__empty">Записей пока не найдено.</td>
            </tr>
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

.enrollments-table__empty {
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
