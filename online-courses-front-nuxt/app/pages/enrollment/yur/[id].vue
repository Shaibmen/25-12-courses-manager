<script setup lang="ts">
import {
  ageCategories,
  legalEntityEnrollmentContracts,
  loadVariantsNotDO,
  optDocumentOptions
} from '../../../types/enrollment'
import AppButton from '../../../components/ui/AppButton.vue'
import AppCard from '../../../components/ui/AppCard.vue'
import AppInput from '../../../components/ui/AppInput.vue'

const route = useRoute()
const router = useRouter()
const notifications = useNotifications()
const legalEntityId = computed(() => String(route.params.id || ''))

const details = ref<Awaited<ReturnType<typeof getLegalEntityDetails>> | null>(null)
const programs = ref<Awaited<ReturnType<typeof getPrograms>>>([])
const executers = ref<Awaited<ReturnType<typeof getExecuters>>>([])
const loading = ref(true)
const saving = ref(false)

const selectedListenerIds = ref<string[]>([])
const selectAll = ref(false)
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
const studyLoadOption = ref('')
const optDocumentSelected = ref('')

const selectedProgram = computed(() =>
  programs.value.find((item) => item.id_program_education === selectedProgramId.value) || null
)

watch(selectAll, (value) => {
  selectedListenerIds.value = value
    ? details.value?.legal_entity.listeners?.map((item) => item.id_listener) || []
    : []
})

watch(selectedListenerIds, (value) => {
  const total = details.value?.legal_entity.listeners?.length || 0
  selectAll.value = total > 0 && value.length === total
})

watch(selectedProgramId, () => {
  if (selectedProgram.value) {
    currentPrice.value = selectedProgram.value.individual_price
  }
})

watch(selectedContractId, (value) => {
  const id = value.toUpperCase()

  if (id.startsWith('PK')) {
    typeOfRetraining.value = 'Повышение квалификации'
  } else if (id.startsWith('PP')) {
    typeOfRetraining.value = 'Профессиональная переподготовка'
  } else {
    typeOfRetraining.value = ''
  }
})

const isFormValid = computed(() =>
  Boolean(
    selectedListenerIds.value.length &&
    selectedProgramId.value &&
    startDate.value &&
    endDate.value &&
    currentPrice.value &&
    selectedExecutorId.value &&
    selectedContractId.value
  )
)

const load = async () => {
  loading.value = true

  try {
    const [detailsData, programData, executerData] = await Promise.all([
      getLegalEntityDetails(legalEntityId.value),
      getPrograms(1, ''),
      getExecuters(1, '')
    ])

    details.value = detailsData
    programs.value = programData
    executers.value = executerData
    selectedExecutorId.value = executerData[0]?.id_executor || ''
    selectedProgramId.value = programData[0]?.id_program_education || ''
    currentPrice.value = programData[0]?.individual_price || ''
  } catch (error) {
    notifications.error(
      error instanceof Error ? error.message : 'Не удалось загрузить запись по юрлицу',
      'Запись на курс'
    )
  } finally {
    loading.value = false
  }
}

const getPaymentText = () => {
  if (paymentOption.value === 'split') {
    return `Оплата осуществляется авансом 50 %, оставшиеся 50 % — в установленный срок${secondPaymentDate.value ? ` ${secondPaymentDate.value}` : ''}`
  }

  if (paymentOption.value === 'full') {
    return 'Оплата осуществляется в следующем порядке: 100% предоплата до начала обучения.'
  }

  if (paymentOption.value === 'halfsplit') {
    return 'Оплата перечисляется на расчётный счёт Исполнителя в течение 5 рабочих дней после подписания акта.'
  }

  return null
}

const createEnrollment = async () => {
  if (!isFormValid.value || !details.value || !selectedProgram.value) {
    notifications.error('Заполните обязательные поля записи по юрлицу.', 'Запись на курс')
    return
  }

  saving.value = true

  try {
    const optNagruz = studyLoadOption.value ? Number(studyLoadOption.value) : null

    await Promise.all(
      selectedListenerIds.value.map((listenerId) =>
        createEnrollmentRequest({
          id_listener: listenerId,
          id_program: selectedProgramId.value,
          start_date: startDate.value,
          end_date: endDate.value,
          current_price: Number(currentPrice.value),
          group: group.value || null,
          type_of_retraining: typeOfRetraining.value || null,
          is_active: true,
          opt_nagruz: optNagruz
        })
      )
    )

    const legalEntity = details.value.legal_entity
    const regAddress = details.value.reg_address
    const selectedListeners = legalEntity.listeners?.filter((item) => selectedListenerIds.value.includes(item.id_listener)) || []

    await createEnrollmentDocumentRequest({
      id_listener: selectedListenerIds.value[0],
      id_program: selectedProgramId.value,
      id_executor: selectedExecutorId.value || null,
      front_data: {
        legal_entity: {
          listeners: selectedListeners,
          reg_address: regAddress,
          company_name: legalEntity.name_company,
          zakazchikfio: `${legalEntity.second_name} ${legalEntity.first_name} ${legalEntity.middle_name || ''}`.trim(),
          status: legalEntity.status,
          inn: legalEntity.inn,
          kpp: legalEntity.kpp,
          ogrn: legalEntity.ogrn,
          phone: legalEntity.phone,
          email: legalEntity.email
        },
        start_date: startDate.value,
        end_date: endDate.value,
        program_name: selectedProgram.value.name_prof_education,
        time_education: selectedProgram.value.time_education,
        price_enrollment: Number(currentPrice.value),
        variant: optNagruz,
        dogovor_type: selectedContractId.value,
        opt_nagruz: optNagruz,
        opt_document: optDocumentSelected.value ? Number(optDocumentSelected.value) : null,
        dogovor_age: ageCategory.value || 'EIGHTEEN',
        opt_price: getPaymentText()
      }
    })

    await Promise.all(
      selectedListenerIds.value.map((listenerId) =>
        createEnrollmentDocumentRequest({
          id_listener: listenerId,
          id_program: selectedProgramId.value,
          id_executor: selectedExecutorId.value || null,
          front_data: {
            dogovor_type: selectedContractId.value,
            dogovor_age: ageCategory.value || 'EIGHTEEN',
            opt_document: optDocumentSelected.value ? Number(optDocumentSelected.value) : null,
            opt_price: getPaymentText(),
            opt_nagruz: optNagruz,
            variant: optNagruz
          }
        })
      )
    )

    notifications.success(`Создано записей: ${selectedListenerIds.value.length}.`, 'Запись на курс')
    await router.push(`/legalentities/${legalEntityId.value}`)
  } catch (error) {
    notifications.error(
      error instanceof Error ? error.message : 'Не удалось создать записи по юрлицу',
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
      <p>Подтягиваю данные по юридическому лицу и программам.</p>
    </AppCard>

    <template v-else-if="details">
      <AppCard title="Слушатели юридического лица">
        <div class="listener-actions">
          <label class="checkbox-row">
            <input v-model="selectAll" type="checkbox">
            <span>Выбрать всех</span>
          </label>
          <span class="listener-count">Выбрано: {{ selectedListenerIds.length }}</span>
        </div>

        <div class="listener-list">
          <label v-for="listener in details.legal_entity.listeners || []" :key="listener.id_listener" class="checkbox-card">
            <input v-model="selectedListenerIds" type="checkbox" :value="listener.id_listener">
            <span>{{ listener.second_name }} {{ listener.first_name }} {{ listener.middle_name || '' }}</span>
          </label>
        </div>
      </AppCard>

      <AppCard title="Параметры договора">
        <div class="enrollment-grid">
          <label class="app-select">
            <span class="app-select__label">Договор</span>
            <div class="app-select__field" :class="{ 'app-select__field--placeholder': !selectedContractId }">
              <select v-model="selectedContractId" class="app-select__control">
                <option value="">Выберите договор</option>
                <option v-for="item in legalEntityEnrollmentContracts" :key="item.id_contract" :value="item.id_contract">
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
            <span class="app-select__label">Недельная нагрузка</span>
            <div class="app-select__field" :class="{ 'app-select__field--placeholder': !studyLoadOption }">
              <select v-model="studyLoadOption" class="app-select__control">
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
          <AppInput v-model="typeOfRetraining" label="Тип обучения" disabled />
        </div>
      </AppCard>

      <div class="page-actions">
        <AppButton variant="ghost" @click="router.push(`/legalentities/${legalEntityId}`)">Назад</AppButton>
        <AppButton :disabled="saving" @click="createEnrollment">
          {{ saving ? 'Создаём...' : 'Создать записи' }}
        </AppButton>
      </div>
    </template>
  </section>
</template>

<style scoped>
.listener-actions,
.page-actions {
  display: flex;
  justify-content: space-between;
  gap: 1rem;
  flex-wrap: wrap;
}

.listener-count {
  color: #64748b;
}

.listener-list,
.radio-grid,
.enrollment-grid {
  display: grid;
  gap: 1rem;
}

.listener-list,
.radio-grid {
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
}

.checkbox-row,
.checkbox-card,
.radio-card {
  display: flex;
  gap: 0.75rem;
  align-items: flex-start;
}

.checkbox-card,
.radio-card {
  padding: 1rem;
  border-radius: 1rem;
  border: 1px solid rgba(15, 23, 42, 0.08);
  background: rgba(255, 255, 255, 0.7);
}

.enrollment-grid {
  grid-template-columns: repeat(2, minmax(0, 1fr));
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

@media (max-width: 900px) {
  .enrollment-grid {
    grid-template-columns: 1fr;
  }

  .page-actions {
    flex-direction: column;
  }
}
</style>
