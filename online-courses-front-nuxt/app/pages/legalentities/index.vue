<script setup lang="ts">
import type { LegalEntityListItem } from '../../types/legalentity'
import AppButton from '../../components/ui/AppButton.vue'
import AppCard from '../../components/ui/AppCard.vue'
import AppConfirmDialog from '../../components/ui/AppConfirmDialog.vue'
import AppInput from '../../components/ui/AppInput.vue'

const router = useRouter()
const notifications = useNotifications()
const items = ref<LegalEntityListItem[]>([])
const page = ref(1)
const filter = ref('')
const hasMore = ref(false)
const loading = ref(false)
const deleteLoading = ref(false)
const errorMessage = ref('')
const itemToDelete = ref<LegalEntityListItem | null>(null)
let debounceTimer: ReturnType<typeof setTimeout> | null = null

const loadLegalEntitiesList = async () => {
  loading.value = true
  errorMessage.value = ''

  try {
    const data = await getLegalEntities(page.value, filter.value)
    items.value = data
    hasMore.value = data.length === 25
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : 'Не удалось загрузить юридические лица'
  } finally {
    loading.value = false
  }
}

const confirmDelete = async () => {
  if (!itemToDelete.value) {
    return
  }

  deleteLoading.value = true

  try {
    await deleteLegalEntityRequest(itemToDelete.value.id_legalentity)
    itemToDelete.value = null
    await loadLegalEntitiesList()
    notifications.success('Юридическое лицо удалено. Список обновлён.', 'Юридические лица')
  } catch (error) {
    notifications.error(
      error instanceof Error ? error.message : 'Не удалось удалить юридическое лицо',
      'Юридические лица'
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
    void loadLegalEntitiesList()
  }, 350)
})

onBeforeUnmount(() => {
  if (debounceTimer) {
    clearTimeout(debounceTimer)
  }
})

onMounted(() => {
  void loadLegalEntitiesList()
})
</script>

<template>
  <section class="stack content-shell content-shell--wide">
    <AppCard title="Юридические лица">
      <div class="toolbar">
        <div class="toolbar__search">
          <AppInput v-model="filter" label="Поиск" placeholder="Поиск по названию компании" />
        </div>

        <div class="toolbar__actions">
          <AppButton to="/legalentities/create">Добавить юридическое лицо</AppButton>
          <AppButton variant="secondary" to="/dashboard">Назад</AppButton>
        </div>
      </div>
    </AppCard>

    <AppCard v-if="loading" title="Загрузка">
      <p>Подтягиваю список юридических лиц.</p>
    </AppCard>

    <AppCard v-else-if="errorMessage" title="Ошибка">
      <p>{{ errorMessage }}</p>
      <AppButton variant="secondary" @click="loadLegalEntitiesList">Повторить запрос</AppButton>
    </AppCard>

    <AppCard v-else title="Список">
      <div class="table-wrap">
        <table class="legalentities-table">
          <thead>
            <tr>
              <th>Компания</th>
              <th>ИНН</th>
              <th>КПП</th>
              <th>ОГРН</th>
              <th>Телефон</th>
              <th>Email</th>
              <th>Контактное лицо</th>
              <th>Должность</th>
              <th>Действия</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="!items.length">
              <td colspan="9" class="legalentities-table__empty">Юридические лица пока не найдены.</td>
            </tr>
            <tr v-for="item in items" :key="item.id_legalentity">
              <td class="legalentities-table__company">{{ item.name_company }}</td>
              <td>{{ item.inn || '—' }}</td>
              <td>{{ item.kpp || '—' }}</td>
              <td>{{ item.ogrn || '—' }}</td>
              <td>{{ item.phone || '—' }}</td>
              <td class="legalentities-table__company">{{ item.email || '—' }}</td>
              <td class="legalentities-table__company">
                {{ item.second_name }} {{ item.first_name }} {{ item.middle_name || '' }}
              </td>
              <td class="legalentities-table__company">{{ item.status || '—' }}</td>
              <td>
                <div class="legalentities-table__actions">
                  <AppButton variant="ghost" @click="router.push(`/legalentities/${item.id_legalentity}`)">
                    Подробнее
                  </AppButton>
                  <AppButton variant="secondary" @click="router.push(`/legalentities/edit/${item.id_legalentity}`)">
                    Изменить
                  </AppButton>
                  <AppButton variant="ghost" @click="itemToDelete = item">
                    Удалить
                  </AppButton>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="pagination">
        <AppButton variant="secondary" :disabled="page <= 1" @click="page > 1 && (page--, loadLegalEntitiesList())">
          Назад
        </AppButton>
        <span>Страница {{ page }}</span>
        <AppButton variant="secondary" :disabled="!hasMore" @click="page++, loadLegalEntitiesList()">
          Вперёд
        </AppButton>
      </div>
    </AppCard>

    <AppConfirmDialog
      :open="Boolean(itemToDelete)"
      title="Удаление юридического лица"
      :message="itemToDelete ? `Удалить организацию «${itemToDelete.name_company}»?` : ''"
      :loading="deleteLoading"
      confirm-label="Удалить"
      @cancel="itemToDelete = null"
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
.legalentities-table__actions {
  display: flex;
  gap: 0.75rem;
  flex-wrap: wrap;
}

.table-wrap {
  overflow: visible;
}

.legalentities-table {
  width: 100%;
  border-collapse: collapse;
  table-layout: fixed;
}

.legalentities-table th,
.legalentities-table td {
  padding: 1rem 0.95rem;
  border-bottom: 1px solid rgba(15, 23, 42, 0.08);
  text-align: left;
  vertical-align: middle;
}

.legalentities-table th {
  color: #475569;
  font-size: 0.82rem;
  text-transform: uppercase;
  letter-spacing: 0.06em;
}

.legalentities-table tbody tr:hover {
  background: rgba(219, 234, 254, 0.26);
}

.legalentities-table__company {
  max-width: 18rem;
  overflow-wrap: anywhere;
  word-break: break-word;
}

.legalentities-table__empty {
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

.content-shell--wide {
  width: min(1560px, 100%);
}
</style>
