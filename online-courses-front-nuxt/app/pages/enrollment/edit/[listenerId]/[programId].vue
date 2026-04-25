<script setup lang="ts">
import type { GroupItem } from '../../../../types/group'
import type { ProgramListItem } from '../../../../types/program'
import AppButton from '../../../../components/ui/AppButton.vue'
import AppCard from '../../../../components/ui/AppCard.vue'
import AppInput from '../../../../components/ui/AppInput.vue'
import AppSelect from '../../../../components/ui/AppSelect.vue'

const route = useRoute()
const router = useRouter()
const notifications = useNotifications()
const listenerId = computed(() => String(route.params.listenerId || ''))
const programId = computed(() => String(route.params.programId || ''))

const loading = ref(false)
const saving = ref(false)
const listenerFio = ref('')
const program = ref<ProgramListItem | null>(null)
const groups = ref<GroupItem[]>([])
const startDate = ref('')
const endDate = ref('')
const selectedGroupId = ref('')
const typeOfRetraining = ref('')
const errorMessage = ref('')

const isValid = computed(() =>
  Boolean(startDate.value && endDate.value && selectedGroupId.value && typeOfRetraining.value.trim())
)

const load = async () => {
  loading.value = true
  errorMessage.value = ''

  try {
    const [listenerDetails, programDetails, enrollmentDetails, groupData] = await Promise.all([
      getListenerDetails(listenerId.value),
      getProgramDetails(programId.value),
      getListenerEnrollments(listenerId.value),
      getGroups(1, '')
    ])

    listenerFio.value = [
      listenerDetails.listener.second_name,
      listenerDetails.listener.first_name,
      listenerDetails.listener.middle_name
    ].filter(Boolean).join(' ')

    program.value = programDetails
    groups.value = groupData

    const currentEnrollment = enrollmentDetails.find((item) => item.id_program_education === programId.value)

    if (!currentEnrollment) {
      throw new Error('Запись по программе не найдена')
    }

    startDate.value = normalizeApiDate(currentEnrollment.start_date)
    endDate.value = normalizeApiDate(currentEnrollment.end_date)
    selectedGroupId.value = currentEnrollment.id_group
    typeOfRetraining.value = currentEnrollment.type_of_retraining || ''
  } catch (error) {
    errorMessage.value =
      error instanceof Error ? error.message : 'Не удалось загрузить запись'
  } finally {
    loading.value = false
  }
}

const submit = async () => {
  if (!isValid.value) {
    notifications.error('Заполните обязательные поля записи.', 'Запись на курс')
    return
  }

  saving.value = true

  try {
    await updateEnrollmentRequest(listenerId.value, programId.value, {
      id_program: programId.value,
      start_date: startDate.value,
      end_date: endDate.value,
      id_group: selectedGroupId.value,
      type_of_retraining: typeOfRetraining.value.trim()
    })

    notifications.success('Запись обновлена.', 'Запись на курс')
    await router.push(`/enrollment/details/${listenerId.value}`)
  } catch (error) {
    notifications.error(
      error instanceof Error ? error.message : 'Не удалось обновить запись',
      'Запись на курс'
    )
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  void load()
})
</script>

<template>
  <section class="stack content-shell">
    <AppCard title="Редактирование записи">
      <div class="summary">
        <p><strong>Слушатель:</strong> {{ listenerFio || '—' }}</p>
        <p><strong>Программа:</strong> {{ program?.name_prof_education || '—' }}</p>
        <p><strong>Цена программы:</strong> {{ program ? `${new Intl.NumberFormat('ru-RU').format(program.price)} ₽` : '—' }}</p>
      </div>
    </AppCard>

    <AppCard v-if="loading" title="Загрузка">
      <p>Подтягиваю данные записи.</p>
    </AppCard>

    <AppCard v-else-if="errorMessage" title="Ошибка">
      <p>{{ errorMessage }}</p>
      <AppButton variant="secondary" @click="load">Повторить запрос</AppButton>
    </AppCard>

    <AppCard v-else title="Параметры записи">
      <form class="form-grid" @submit.prevent="submit">
        <AppInput v-model="startDate" label="Дата начала" type="date" />
        <AppInput v-model="endDate" label="Дата окончания" type="date" />

        <AppSelect v-model="selectedGroupId" label="Группа" placeholder="Выберите группу">
          <option v-for="group in groups" :key="group.group" :value="group.group">
            {{ group.name_group }}
          </option>
        </AppSelect>

        <AppInput
          v-model="typeOfRetraining"
          label="Тип обучения"
          placeholder="Например: Повышение квалификации"
        />

        <div class="form-actions">
          <AppButton type="button" variant="ghost" @click="router.push(`/enrollment/details/${listenerId}`)">Назад</AppButton>
          <AppButton type="submit" :disabled="saving">
            {{ saving ? 'Сохраняем...' : 'Сохранить изменения' }}
          </AppButton>
        </div>
      </form>
    </AppCard>
  </section>
</template>

<style scoped>
.summary,
.form-grid {
  display: grid;
  gap: 1rem;
}

.form-grid {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.form-actions {
  grid-column: 1 / -1;
  display: flex;
  justify-content: space-between;
  gap: 1rem;
}

@media (max-width: 900px) {
  .form-grid {
    grid-template-columns: 1fr;
  }

  .form-actions {
    flex-direction: column;
  }
}
</style>
