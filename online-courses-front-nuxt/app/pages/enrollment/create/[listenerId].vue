<script setup lang="ts">
import type { ContractorPayload } from '../../../types/enrollment'
import {
  ageCategories,
  loadVariantsDO,
  loadVariantsNotDO,
  optDocumentOptions,
  singleEnrollmentContracts
} from '../../../types/enrollment'
import ContractorForm from '../../../components/features/enrollments/ContractorForm.vue'
import AppButton from '../../../components/ui/AppButton.vue'
import AppCard from '../../../components/ui/AppCard.vue'

const route = useRoute()
const router = useRouter()
const notifications = useNotifications()
const listenerId = computed(() => String(route.params.listenerId || ''))

const listenerContext = ref<Awaited<ReturnType<typeof getListenerEnrollmentContext>> | null>(null)
const programs = ref<Awaited<ReturnType<typeof getPrograms>>>([])
const executers = ref<Awaited<ReturnType<typeof getExecuters>>>([])
const loading = ref(true)
const saving = ref(false)
const contractorSaving = ref(false)
const showContractorForm = ref(false)

const selectedContractId = ref('')
const selectedExecutorId = ref('')
const selectedProgramId = ref('')
const currentPrice = ref<number | ''>('')
const startDate = ref('')
const endDate = ref('')
const group = ref('')
const typeOfRetraining = ref('')
const paymentOption = ref('')
const secondPaymentDate = ref('')
const ageCategory = ref('')
const loadVariant = ref('')
const studyLoadOption = ref('')
const optDocumentSelected = ref('')

const hasContractor = computed(() => Boolean(listenerContext.value?.contractor))
const filteredContracts = computed(() =>
  singleEnrollmentContracts.filter((item) => item.type === (hasContractor.value ? 'trilateral' : 'bilateral'))
)
const selectedProgram = computed(() =>
  programs.value.find((item) => item.id_program_education === selectedProgramId.value) || null
)
const isDO = computed(() => selectedContractId.value.toUpperCase().startsWith('DO'))
const isPKorPP = computed(() => {
  const id = selectedContractId.value.toUpperCase()
  return id.startsWith('PK') || id.startsWith('PP')
})

const isFormValid = computed(() =>
  Boolean(
    selectedProgramId.value &&
    startDate.value &&
    endDate.value &&
    currentPrice.value &&
    selectedExecutorId.value &&
    (!hasContractor.value || selectedContractId.value)
  )
)

const syncProgram = () => {
  if (!selectedProgram.value) {
    currentPrice.value = ''
    return
  }

  currentPrice.value = selectedProgram.value.individual_price
}

watch(selectedProgramId, syncProgram)

watch(selectedContractId, (value) => {
  const id = value.toUpperCase()

  if (id.startsWith('DO')) {
    typeOfRetraining.value = 'Дополнительное образование'
  } else if (id.startsWith('PK')) {
    typeOfRetraining.value = 'Повышение квалификации'
  } else if (id.startsWith('PP')) {
    typeOfRetraining.value = 'Профессиональная переподготовка'
  } else {
    typeOfRetraining.value = ''
  }
})

const load = async () => {
  loading.value = true

  try {
    const [listenerData, programData, executerData] = await Promise.all([
      getListenerEnrollmentContext(listenerId.value),
      getPrograms(1, ''),
      getExecuters(1, '')
    ])

    listenerContext.value = listenerData
    programs.value = programData
    executers.value = executerData
    selectedExecutorId.value = executerData[0]?.id_executor || ''
    selectedProgramId.value = programData[0]?.id_program_education || ''
    syncProgram()
  } catch (error) {
    notifications.error(
      error instanceof Error ? error.message : 'Не удалось загрузить форму записи на курс',
      'Запись на курс'
    )
  } finally {
    loading.value = false
  }
}

const getPaymentText = () => {
  if (paymentOption.value === 'split') {
    return `Оплата осуществляется в следующем порядке: аванс 50 % — предоплата до начала обучения, оставшиеся 50 % — в установленный срок${secondPaymentDate.value ? ` ${secondPaymentDate.value}` : ''}`
  }

  if (paymentOption.value === 'full') {
    return 'Оплата осуществляется в следующем порядке: 100% предоплата до начала обучения.'
  }

  if (paymentOption.value === 'halfsplit') {
    return 'Оплата подлежит перечислению на расчётный счёт Исполнителя в срок до 5 рабочих дней после подписания акта.'
  }

  return null
}

const getNagruzValue = () => {
  if (isDO.value) {
    return loadVariant.value ? Number(loadVariant.value) : null
  }

  if (isPKorPP.value) {
    return studyLoadOption.value ? Number(studyLoadOption.value) : null
  }

  return null
}

const saveContractor = async (payload: ContractorPayload) => {
  contractorSaving.value = true

  try {
    await upsertContractorRequest(listenerId.value, payload)
    notifications.success('Заказчик сохранён.', 'Запись на курс')
    showContractorForm.value = false
    await load()
  } catch (error) {
    notifications.error(
      error instanceof Error ? error.message : 'Не удалось сохранить заказчика',
      'Запись на курс'
    )
  } finally {
    contractorSaving.value = false
  }
}

const removeContractor = async () => {
  const contractorId = listenerContext.value?.contractor?.contractor.id_contractor

  if (!contractorId || !window.confirm('Удалить заказчика у слушателя?')) {
    return
  }

  try {
    await deleteContractorRequest(contractorId)
    notifications.success('Заказчик удалён.', 'Запись на курс')
    await load()
  } catch (error) {
    notifications.error(
      error instanceof Error ? error.message : 'Не удалось удалить заказчика',
      'Запись на курс'
    )
  }
}

const createEnrollment = async () => {
  if (!isFormValid.value || !listenerContext.value || !selectedProgram.value) {
    notifications.error('Заполните обязательные поля записи.', 'Запись на курс')
    return
  }

  saving.value = true

  try {
    const optNagruz = getNagruzValue()

    await createEnrollmentRequest({
      id_listener: listenerId.value,
      id_program: selectedProgramId.value,
      start_date: startDate.value,
      end_date: endDate.value,
      current_price: Number(currentPrice.value),
      group: group.value || null,
      type_of_retraining: typeOfRetraining.value || null,
      is_active: true,
      opt_nagruz: optNagruz
    })

    await createEnrollmentDocumentRequest({
      id_listener: listenerId.value,
      id_program: selectedProgramId.value,
      id_executor: selectedExecutorId.value || null,
      front_data: {
        dogovor_type: selectedContractId.value || null,
        dogovor_age: ageCategory.value || null,
        opt_document: optDocumentSelected.value ? Number(optDocumentSelected.value) : null,
        opt_price: getPaymentText(),
        opt_nagruz: optNagruz
      }
    })

    notifications.success('Запись и документы успешно созданы.', 'Запись на курс')
    await router.push(`/listeners/${listenerId.value}`)
  } catch (error) {
    notifications.error(
      error instanceof Error ? error.message : 'Не удалось создать запись на курс',
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
  <section class="stack">
    <AppCard v-if="loading" title="Загрузка">
      <p>Подтягиваю данные для записи на курс.</p>
    </AppCard>

    <template v-else-if="listenerContext">
      <AppCard title="Слушатель">
        <p class="summary-line">
          {{ listenerContext.listener.second_name }} {{ listenerContext.listener.first_name }} {{ listenerContext.listener.middle_name || '' }}
        </p>
      </AppCard>

      <AppCard title="Заказчик">
        <div v-if="listenerContext.contractor" class="summary-stack">
          <p class="summary-line">
            {{ listenerContext.contractor.contractor.second_name }} {{ listenerContext.contractor.contractor.first_name }} {{ listenerContext.contractor.contractor.middle_name }}
          </p>
          <p class="summary-subline">{{ listenerContext.contractor.contractor.contact_phone }} · {{ listenerContext.contractor.contractor.email }}</p>
          <div class="summary-actions">
            <AppButton variant="secondary" @click="showContractorForm = !showContractorForm">
              {{ showContractorForm ? 'Скрыть форму' : 'Изменить заказчика' }}
            </AppButton>
            <AppButton variant="ghost" @click="removeContractor">Удалить заказчика</AppButton>
          </div>
        </div>

        <div v-else class="summary-actions">
          <p class="summary-subline">Заказчик пока не добавлен.</p>
          <AppButton @click="showContractorForm = !showContractorForm">
            {{ showContractorForm ? 'Скрыть форму' : 'Добавить заказчика' }}
          </AppButton>
        </div>
      </AppCard>

      <ContractorForm
        v-if="showContractorForm"
        :initial-state="listenerContext.contractor || undefined"
        :loading="contractorSaving"
        @submit="saveContractor"
        @cancel="showContractorForm = false"
      />

      <AppCard title="Параметры договора">
        <div class="enrollment-grid">
          <label class="app-select">
            <span class="app-select__label">Договор</span>
            <div class="app-select__field" :class="{ 'app-select__field--placeholder': !selectedContractId }">
              <select v-model="selectedContractId" class="app-select__control">
                <option value="">Выберите договор</option>
                <option v-for="item in filteredContracts" :key="item.id_contract" :value="item.id_contract">
                  {{ item.name }}
                </option>
              </select>
              <span class="app-select__icon" aria-hidden="true">⌄</span>
            </div>
          </label>

          <label class="app-select">
            <span class="app-select__label">Порядок оплаты</span>
            <div class="app-select__field" :class="{ 'app-select__field--placeholder': !paymentOption }">
              <select v-model="paymentOption" class="app-select__control">
                <option value="">Выберите порядок оплаты</option>
                <option value="full">100% предоплата</option>
                <option value="split">50/50</option>
                <option value="halfsplit">Оплата после акта</option>
              </select>
              <span class="app-select__icon" aria-hidden="true">⌄</span>
            </div>
          </label>

          <AppInput v-if="paymentOption === 'split'" v-model="secondPaymentDate" label="Срок второй оплаты" type="date" />

          <label class="app-select">
            <span class="app-select__label">Возрастная категория</span>
            <div class="app-select__field" :class="{ 'app-select__field--placeholder': !ageCategory }">
              <select v-model="ageCategory" class="app-select__control">
                <option value="">Выберите категорию</option>
                <option v-for="(label, key) in ageCategories" :key="key" :value="key">{{ label }}</option>
              </select>
              <span class="app-select__icon" aria-hidden="true">⌄</span>
            </div>
          </label>

          <label class="app-select">
            <span class="app-select__label">Итоговый документ</span>
            <div class="app-select__field" :class="{ 'app-select__field--placeholder': !optDocumentSelected }">
              <select v-model="optDocumentSelected" class="app-select__control">
                <option value="">Выберите режим выдачи</option>
                <option v-for="(label, key) in optDocumentOptions" :key="key" :value="String(key)">{{ label }}</option>
              </select>
              <span class="app-select__icon" aria-hidden="true">⌄</span>
            </div>
          </label>

          <label class="app-select">
            <span class="app-select__label">Нагрузка ДО</span>
            <div class="app-select__field" :class="{ 'app-select__field--placeholder': !loadVariant }">
              <select v-model="loadVariant" class="app-select__control" :disabled="!isDO">
                <option value="">Выберите вариант</option>
                <option v-for="(label, key) in loadVariantsDO" :key="key" :value="String(key)">{{ label }}</option>
              </select>
              <span class="app-select__icon" aria-hidden="true">⌄</span>
            </div>
          </label>

          <label class="app-select">
            <span class="app-select__label">Нагрузка ПК/ПП</span>
            <div class="app-select__field" :class="{ 'app-select__field--placeholder': !studyLoadOption }">
              <select v-model="studyLoadOption" class="app-select__control" :disabled="!isPKorPP">
                <option value="">Выберите вариант</option>
                <option v-for="(label, key) in loadVariantsNotDO" :key="key" :value="String(key)">{{ label }}</option>
              </select>
              <span class="app-select__icon" aria-hidden="true">⌄</span>
            </div>
          </label>
        </div>
      </AppCard>

      <AppCard title="Исполнитель">
        <div class="radio-grid">
          <label v-for="executer in executers" :key="executer.id_executor" class="radio-card">
            <input v-model="selectedExecutorId" type="radio" :value="executer.id_executor">
            <span>{{ executer.second_name }} {{ executer.first_name }} {{ executer.middle_name || '' }}</span>
          </label>
        </div>
      </AppCard>

      <AppCard title="Программа обучения">
        <div class="enrollment-grid">
          <label class="app-select app-select--full">
            <span class="app-select__label">Программа</span>
            <div class="app-select__field" :class="{ 'app-select__field--placeholder': !selectedProgramId }">
              <select v-model="selectedProgramId" class="app-select__control">
                <option value="">Выберите программу</option>
                <option v-for="program in programs" :key="program.id_program_education" :value="program.id_program_education">
                  {{ program.name_prof_education }}
                </option>
              </select>
              <span class="app-select__icon" aria-hidden="true">⌄</span>
            </div>
          </label>

          <AppInput v-model="startDate" label="Дата начала" type="date" />
          <AppInput v-model="endDate" label="Дата окончания" type="date" />

          <label class="app-select">
            <span class="app-select__label">Цена</span>
            <div class="app-select__field" :class="{ 'app-select__field--placeholder': !currentPrice }">
              <select v-model="currentPrice" class="app-select__control">
                <option :value="selectedProgram?.individual_price || ''">Индивидуальное: {{ selectedProgram?.individual_price ?? '—' }} ₽</option>
                <option :value="selectedProgram?.group_price || ''">Групповое: {{ selectedProgram?.group_price ?? '—' }} ₽</option>
                <option :value="selectedProgram?.campus_price || ''">Кампус: {{ selectedProgram?.campus_price ?? '—' }} ₽</option>
              </select>
              <span class="app-select__icon" aria-hidden="true">⌄</span>
            </div>
          </label>

          <AppInput v-model="group" label="Группа" placeholder="Например: 0" />

          <AppInput v-model="typeOfRetraining" label="Тип обучения" placeholder="Подставляется автоматически" disabled />
        </div>
      </AppCard>

      <div class="page-actions">
        <AppButton variant="ghost" @click="router.push(`/listeners/${listenerId}`)">Назад</AppButton>
        <AppButton :disabled="saving" @click="createEnrollment">
          {{ saving ? 'Создаём...' : 'Создать запись' }}
        </AppButton>
      </div>
    </template>
  </section>
</template>

<style scoped>
.summary-line {
  margin: 0;
  font-weight: 700;
  color: #0f172a;
}

.summary-subline {
  margin: 0;
  color: #64748b;
}

.summary-stack,
.summary-actions,
.enrollment-grid,
.radio-grid {
  display: grid;
  gap: 1rem;
}

.summary-actions {
  display: flex;
  flex-wrap: wrap;
}

.enrollment-grid {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.radio-grid {
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
}

.radio-card {
  display: flex;
  gap: 0.75rem;
  align-items: flex-start;
  padding: 1rem;
  border-radius: 1rem;
  border: 1px solid rgba(15, 23, 42, 0.08);
  background: rgba(255, 255, 255, 0.7);
}

.app-select {
  display: grid;
  gap: 0.45rem;
}

.app-select--full {
  grid-column: 1 / -1;
}

.app-select__label {
  font-size: 0.92rem;
  font-weight: 600;
  color: #0f172a;
}

.app-select__field {
  position: relative;
  border-radius: 1rem;
  border: 1px solid rgba(15, 23, 42, 0.12);
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.97) 0%, rgba(248, 250, 252, 0.94) 100%);
}

.app-select__field--placeholder .app-select__control {
  color: #64748b;
}

.app-select__control {
  appearance: none;
  width: 100%;
  min-height: 3rem;
  padding: 0.8rem 2.75rem 0.8rem 1rem;
  border: none;
  border-radius: 1rem;
  background: transparent;
  color: #0f172a;
  font: inherit;
  outline: none;
}

.app-select__icon {
  position: absolute;
  top: 50%;
  right: 0.95rem;
  transform: translateY(-50%);
  color: #475569;
  pointer-events: none;
}

.page-actions {
  display: flex;
  justify-content: space-between;
  gap: 1rem;
}

@media (max-width: 900px) {
  .enrollment-grid {
    grid-template-columns: 1fr;
  }

  .page-actions {
    flex-direction: column;
  }
}
</style>
