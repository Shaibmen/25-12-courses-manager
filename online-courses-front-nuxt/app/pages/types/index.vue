<script setup lang="ts">
import AppButton from '../../components/ui/AppButton.vue'
import AppCard from '../../components/ui/AppCard.vue'
import AppConfirmDialog from '../../components/ui/AppConfirmDialog.vue'
import AppInput from '../../components/ui/AppInput.vue'
import type { EducationTypeItem } from '../../types/catalogs'

const router = useRouter()
const notifications = useNotifications()
const filter = ref('')
const items = ref<EducationTypeItem[]>([])
const loading = ref(false)
const deleteLoading = ref(false)
const errorMessage = ref('')
const typeToDelete = ref<EducationTypeItem | null>(null)
let debounceTimer: ReturnType<typeof setTimeout> | null = null

const loadTypes = async () => {
  loading.value = true
  errorMessage.value = ''

  try {
    items.value = await getEducationTypes(filter.value)
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : 'Не удалось загрузить типы обучения'
  } finally {
    loading.value = false
  }
}

const confirmDelete = async () => {
  if (!typeToDelete.value) {
    return
  }

  deleteLoading.value = true

  try {
    await deleteEducationTypeRequest(typeToDelete.value.id_educationType)
    typeToDelete.value = null
    await loadTypes()
    notifications.success('Тип обучения удалён. Справочник обновлён.', 'Типы обучения')
  } catch {
    notifications.error('Не удалось удалить тип обучения. Возможно, он используется в системе.', 'Типы обучения')
  } finally {
    deleteLoading.value = false
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

onBeforeUnmount(() => {
  if (debounceTimer) {
    clearTimeout(debounceTimer)
  }
})

onMounted(() => {
  void loadTypes()
})
</script>

<template>
  <section class="stack content-shell">
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
                  <AppButton variant="ghost" @click="typeToDelete = item">
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
      :open="Boolean(typeToDelete)"
      title="Удаление типа обучения"
      :message="typeToDelete ? `Удалить тип «${typeToDelete.typeName}»?` : ''"
      :loading="deleteLoading"
      confirm-label="Удалить"
      @cancel="typeToDelete = null"
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
