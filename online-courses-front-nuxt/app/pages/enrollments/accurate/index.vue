<script setup lang="ts">
import type { AccurateEnrollmentItem } from '../../../types/enrollment'
import AppButton from '../../../components/ui/AppButton.vue'
import AppCard from '../../../components/ui/AppCard.vue'
import AppInput from '../../../components/ui/AppInput.vue'

type AccurateEnrollmentView = AccurateEnrollmentItem & {
  listener_fio: string
}

const router = useRouter()
const enrollments = ref<AccurateEnrollmentView[]>([])
const loading = ref(false)
const errorMessage = ref('')
const filter = ref('')

const filteredEnrollments = computed(() => {
  if (!filter.value) {
    return enrollments.value
  }

  const normalized = filter.value.toLowerCase().trim()
  return enrollments.value.filter((item) =>
    item.listener_fio.toLowerCase().includes(normalized)
  )
})

const loadEnrollments = async () => {
  loading.value = true
  errorMessage.value = ''

  try {
    const data = await getAccurateEnrollments()

    const listenerMeta = await Promise.all(
      data.map(async (item) => {
        try {
          const details = await getListenerDetails(item.id_listener)
          const listener = details.listener

          return {
            ...item,
            listener_fio: [listener.second_name, listener.first_name, listener.middle_name]
              .filter(Boolean)
              .join(' ')
          }
        } catch {
          return {
            ...item,
            listener_fio: item.id_listener
          }
        }
      })
    )

    enrollments.value = listenerMeta
  } catch (error) {
    errorMessage.value =
      error instanceof Error ? error.message : 'Не удалось загрузить точные записи'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  void loadEnrollments()
})
</script>

<template>
  <section class="stack content-shell">
    <AppCard title="Точные записи на курс">
      <div class="toolbar">
        <div class="toolbar__search">
          <AppInput v-model="filter" label="Поиск" placeholder="Фильтр по ФИО" />
        </div>

        <div class="toolbar__actions">
          <AppButton variant="secondary" @click="loadEnrollments">Обновить</AppButton>
          <AppButton variant="secondary" to="/enrollments">Назад</AppButton>
        </div>
      </div>
    </AppCard>

    <AppCard v-if="loading" title="Загрузка">
      <p>Подтягиваю точные записи.</p>
    </AppCard>

    <AppCard v-else-if="errorMessage" title="Ошибка">
      <p>{{ errorMessage }}</p>
      <AppButton variant="secondary" @click="loadEnrollments">Повторить запрос</AppButton>
    </AppCard>

    <AppCard v-else title="Каталог accurate записей">
      <div class="table-wrap">
        <table class="enrollments-table">
          <thead>
            <tr>
              <th>Слушатель</th>
              <th>Курс</th>
              <th>Часы</th>
              <th>Цена</th>
              <th>Тип обучения</th>
              <th>Подразделение</th>
              <th>Действия</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="!filteredEnrollments.length">
              <td colspan="7" class="enrollments-table__empty">Точных записей пока не найдено.</td>
            </tr>
            <tr v-for="item in filteredEnrollments" :key="item.id_listener + item.name_prof_education">
              <td>{{ item.listener_fio }}</td>
              <td class="enrollments-table__name">{{ item.name_prof_education }}</td>
              <td>{{ item.time_education }}</td>
              <td>{{ new Intl.NumberFormat('ru-RU').format(item.price) }} ₽</td>
              <td>{{ item.education_type || '—' }}</td>
              <td>{{ item.division || '—' }}</td>
              <td>
                <AppButton variant="secondary" @click="router.push(`/enrollment/details/${item.id_listener}`)">
                  Подробнее
                </AppButton>
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
</style>
