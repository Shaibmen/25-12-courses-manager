<script setup lang="ts">
import AppButton from '../../components/ui/AppButton.vue'
import AppCard from '../../components/ui/AppCard.vue'
import AppInput from '../../components/ui/AppInput.vue'
import type { EducationTypeItem } from '../../types/catalogs'

const router = useRouter()
const notifications = useNotifications()
const filter = ref('')
const items = ref<EducationTypeItem[]>([])
const loading = ref(false)
const errorMessage = ref('')
let debounceTimer: ReturnType<typeof setTimeout> | null = null

const loadTypes = async () => {
  loading.value = true
  errorMessage.value = ''

  try {
    items.value = await getEducationTypes(filter.value)
  } catch (error) {
    errorMessage.value =
      error instanceof Error ? error.message : 'Не удалось загрузить типы обучения'
  } finally {
    loading.value = false
  }
}

const removeType = async (id: string) => {
  if (!window.confirm('Вы точно хотите удалить этот тип обучения?')) {
    return
  }

  try {
    await deleteEducationTypeRequest(id)
    await loadTypes()
    notifications.success('Тип обучения удалён. Справочник обновлён.', 'Типы обучения')
  } catch (error) {
    notifications.error('Не удалось удалить тип обучения, возможно он где то используется')
  }
}

watch(filter, () => {
  if (debounceTimer) {
    clearTimeout(debounceTimer)
  }

  debounceTimer = setTimeout(() => {
    void loadTypes()
  }, 350)
})

onMounted(() => {
  void loadTypes()
})
</script>

<template>
  <section class="stack">
    <AppCard title="Типы обучения">
      <div class="toolbar">
        <div class="toolbar__search">
          <AppInput v-model="filter" label="Поиск" placeholder="Поиск по названию типа" />
        </div>

        <div class="toolbar__actions">
          <AppButton to="/types/create">Добавить тип обучения</AppButton>
          <AppButton variant="secondary" to="/dashboard">Назад</AppButton>
        </div>
      </div>
    </AppCard>

    <AppCard v-if="loading" title="Загрузка">
      <p>Подтягиваю типы обучения.</p>
    </AppCard>

    <AppCard v-else-if="errorMessage" title="Ошибка">
      <p>{{ errorMessage }}</p>
      <AppButton variant="secondary" @click="loadTypes">Повторить запрос</AppButton>
    </AppCard>

    <AppCard v-else title="Справочник">
      <div class="table-wrap">
        <table class="catalog-table">
          <thead>
            <tr>
              <th>Название типа</th>
              <th>Действия</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="!items.length">
              <td colspan="2" class="catalog-table__empty">Типы обучения пока не найдены.</td>
            </tr>
            <tr v-for="item in items" :key="item.id_educationType">
              <td>{{ item.typeName }}</td>
              <td>
                <div class="catalog-table__actions">
                  <AppButton variant="secondary" @click="router.push(`/types/edit/${item.id_educationType}`)">
                    Изменить
                  </AppButton>
                  <AppButton variant="ghost" @click="removeType(item.id_educationType)">
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
