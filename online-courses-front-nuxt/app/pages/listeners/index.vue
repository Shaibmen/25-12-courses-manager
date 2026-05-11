<script setup lang="ts">
import type { ListenerListItem } from '../../types/listener'
import { formatApiDate, getDateSortValue } from '../../utils/date'
import AppButton from '../../components/ui/AppButton.vue'
import AppCard from '../../components/ui/AppCard.vue'
import AppConfirmDialog from '../../components/ui/AppConfirmDialog.vue'
import AppInput from '../../components/ui/AppInput.vue'

const router = useRouter()
const notifications = useNotifications()

const listeners = ref<ListenerListItem[]>([])
const page = ref(1)
const filter = ref('')
const loading = ref(false)
const deleteLoading = ref(false)
const hasMore = ref(false)
const sortOrder = ref<'asc' | 'desc'>('asc')
const errorMessage = ref('')
const listenerToDelete = ref<ListenerListItem | null>(null)
let debounceTimer: ReturnType<typeof setTimeout> | null = null

const sortedListeners = computed(() =>
  [...listeners.value].sort((left, right) => {
    const leftTime = getDateSortValue(left.date_of_birth)
    const rightTime = getDateSortValue(right.date_of_birth)

    return sortOrder.value === 'asc' ? leftTime - rightTime : rightTime - leftTime
  })
)

const loadListeners = async () => {
  loading.value = true
  errorMessage.value = ''

  try {
    const data = await getListeners(page.value, filter.value)
    listeners.value = data
    hasMore.value = data.length === 25
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : 'Не удалось загрузить список слушателей'
  } finally {
    loading.value = false
  }
}

const confirmDelete = async () => {
  if (!listenerToDelete.value) {
    return
  }

  deleteLoading.value = true

  try {
    await deleteListenerRequest(listenerToDelete.value.id_listener)
    listenerToDelete.value = null
    await loadListeners()
    notifications.success('Слушатель удалён. Список обновлён.', 'Слушатели')
  } catch (error) {
    notifications.error(
      error instanceof Error ? error.message : 'Не удалось удалить слушателя',
      'Слушатели'
    )
  } finally {
    deleteLoading.value = false
  }
}

const toggleSort = () => {
  sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc'
}

watch(filter, () => {
  if (debounceTimer) {
    clearTimeout(debounceTimer)
  }

  debounceTimer = setTimeout(() => {
    page.value = 1
    void loadListeners()
  }, 350)
})

onBeforeUnmount(() => {
  if (debounceTimer) {
    clearTimeout(debounceTimer)
  }
})

onMounted(() => {
  void loadListeners()
})
</script>

<template>
  <section class="stack content-shell">
    <AppCard title="Слушатели">
      <div class="toolbar">
        <div class="toolbar__search">
          <AppInput
            v-model="filter"
            label="Поиск"
            placeholder="Поиск по фамилии"
          />
        </div>

        <div class="toolbar__actions">
          <AppButton to="/listeners/create">Добавить слушателя</AppButton>
          <AppButton variant="ghost" @click="toggleSort">
            По дате рождения {{ sortOrder === 'asc' ? '↑' : '↓' }}
          </AppButton>
          <AppButton variant="secondary" to="/dashboard">Назад</AppButton>
        </div>
      </div>
    </AppCard>

    <AppCard v-if="loading" title="Загрузка">
      <p>Подтягиваю список слушателей.</p>
    </AppCard>

    <AppCard v-else-if="errorMessage" title="Ошибка">
      <p>{{ errorMessage }}</p>
      <AppButton variant="secondary" @click="loadListeners">
        Повторить запрос
      </AppButton>
    </AppCard>

    <AppCard v-else title="Список">
      <div class="table-wrap">
        <table class="listeners-table">
          <thead>
            <tr>
              <th>Фамилия</th>
              <th>Имя</th>
              <th>Отчество</th>
              <th>Email</th>
              <th>Телефон</th>
              <th>Дата рождения</th>
              <th>СНИЛС</th>
              <th>Действия</th>
            </tr>
          </thead>

          <tbody>
            <tr v-if="!sortedListeners.length">
              <td colspan="8" class="listeners-table__empty">
                Список пока пуст.
              </td>
            </tr>

            <tr v-for="listener in sortedListeners" :key="listener.id_listener">
              <td class="listeners-table__text-cell listeners-table__text-cell--name">{{ listener.second_name }}</td>
              <td class="listeners-table__text-cell listeners-table__text-cell--name">{{ listener.first_name }}</td>
              <td class="listeners-table__text-cell listeners-table__text-cell--name">{{ listener.middle_name || '—' }}</td>
              <td class="listeners-table__text-cell">{{ listener.email }}</td>
              <td class="listeners-table__text-cell">{{ listener.contact_phone }}</td>
              <td>{{ formatApiDate(listener.date_of_birth) }}</td>
              <td class="listeners-table__text-cell">{{ listener.snils }}</td>
              <td>
                <div class="listeners-table__actions">
                  <AppButton variant="ghost" @click="router.push(`/listeners/${listener.id_listener}`)">
                    Подробнее
                  </AppButton>
                  <AppButton variant="secondary" @click="router.push(`/listeners/edit/${listener.id_listener}`)">
                    Изменить
                  </AppButton>
                  <AppButton variant="ghost" @click="listenerToDelete = listener">
                    Удалить
                  </AppButton>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="pagination">
        <AppButton variant="secondary" :disabled="page <= 1" @click="page > 1 && (page--, loadListeners())">
          Назад
        </AppButton>
        <span>Страница {{ page }}</span>
        <AppButton variant="secondary" :disabled="!hasMore" @click="page++, loadListeners()">
          Вперёд
        </AppButton>
      </div>
    </AppCard>

    <AppConfirmDialog
      :open="Boolean(listenerToDelete)"
      title="Удаление слушателя"
      :message="listenerToDelete ? `Удалить слушателя ${listenerToDelete.second_name} ${listenerToDelete.first_name}?` : ''"
      :loading="deleteLoading"
      confirm-label="Удалить"
      @cancel="listenerToDelete = null"
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

.toolbar__actions {
  display: flex;
  gap: 0.75rem;
  flex-wrap: wrap;
}

.table-wrap {
  overflow-x: auto;
}

.listeners-table {
  width: 100%;
  min-width: 1180px;
  border-collapse: collapse;
}

.listeners-table th,
.listeners-table td {
  padding: 1rem 0.95rem;
  border-bottom: 1px solid rgba(15, 23, 42, 0.08);
  text-align: left;
  vertical-align: middle;
}

.listeners-table__text-cell {
  max-width: 14rem;
  overflow-wrap: anywhere;
  word-break: break-word;
}

.listeners-table__text-cell--name {
  max-width: 11rem;
  font-size: clamp(0.82rem, 0.76rem + 0.2vw, 0.95rem);
  line-height: 1.35;
}

.listeners-table th {
  color: #475569;
  font-size: 0.82rem;
  text-transform: uppercase;
  letter-spacing: 0.06em;
}

.listeners-table tbody tr:hover {
  background: rgba(219, 234, 254, 0.26);
}

.listeners-table__actions {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.listeners-table__empty {
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
