<script setup lang="ts">
import type { GroupItem } from '../../../types/group'
import {
  ageCategories,
  legalEntityEnrollmentContracts,
  loadVariantsNotDO,
  optDocumentOptions
} from '../../../types/enrollment'
import AppButton from '../../../components/ui/AppButton.vue'
import AppCard from '../../../components/ui/AppCard.vue'
import AppCheckbox from '../../../components/ui/AppCheckbox.vue'
import AppInput from '../../../components/ui/AppInput.vue'
import AppSelect from '../../../components/ui/AppSelect.vue'

const route = useRoute()
const router = useRouter()
const notifications = useNotifications()
const legalEntityId = computed(() => String(route.params.id || ''))

const details = ref<Awaited<ReturnType<typeof getLegalEntityDetails>> | null>(null)
const programs = ref<Awaited<ReturnType<typeof getPrograms>>>([])
const executers = ref<Awaited<ReturnType<typeof getExecuters>>>([])
const groups = ref<GroupItem[]>([])
const loading = ref(true)
const saving = ref(false)

const selectedListenerIds = ref<string[]>([])
const selectAll = ref(false)
const selectedContractId = ref('')
const selectedExecutorId = ref('')
const selectedProgramId = ref('')
const selectedGroupId = ref('')
const startDate = ref('')
const endDate = ref('')
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
    selectedGroupId.value &&
    startDate.value &&
    endDate.value &&
    selectedExecutorId.value &&
    selectedContractId.value
  )
)

const load = async () => {
  loading.value = true

  try {
    const [detailsData, programData, executerData, groupData] = await Promise.all([
      getLegalEntityDetails(legalEntityId.value),
      getPrograms(1, ''),
      getExecuters(1, ''),
      getGroups(1, '')
    ])

    details.value = detailsData
    programs.value = programData
    executers.value = executerData
    groups.value = groupData
    selectedExecutorId.value = executerData[0]?.id_executor || ''
    selectedProgramId.value = programData[0]?.id_program_education || ''
    selectedGroupId.value = groupData[0]?.group || ''
  } catch (error) {
    notifications.error(
      error instanceof Error ? error.message : 'Не удалось открыть запись по юрлицу',
      'Запись на курс'
    )
  } finally {
    loading.value = false
  }
}

const getPaymentText = () => {
  if (paymentOption.value === 'split') {
    return `Оплата в два этапа, вторая часть${secondPaymentDate.value ? ` до ${secondPaymentDate.value}` : ''}`
  }

  if (paymentOption.value === 'full') {
    return 'Полная предоплата до начала обучения'
  }

  if (paymentOption.value === 'halfsplit') {
    return 'Оплата после подписания акта'
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
    const legalEntity = details.value.legal_entity
    const regAddress = details.value.reg_address
    const selectedListeners =
      legalEntity.listeners?.filter((item) => selectedListenerIds.value.includes(item.id_listener)) || []

    await Promise.all(
      selectedListenerIds.value.map(async (listenerId) => {
        await createEnrollmentRequest({
          id_listener: listenerId,
          id_program: selectedProgramId.value,
          start_date: startDate.value,
          end_date: endDate.value,
          id_group: selectedGroupId.value,
          type_of_retraining: typeOfRetraining.value,
          is_active: true
        })

        await createEnrollmentDocumentRequest({
          id_listener: listenerId,
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
            price_enrollment: selectedProgram.value.price,
            dogovor_type: selectedContractId.value,
            dogovor_age: ageCategory.value || 'EIGHTEEN',
            opt_document: optDocumentSelected.value ? Number(optDocumentSelected.value) : null,
            opt_price: getPaymentText(),
            opt_nagruz: optNagruz,
            variant: optNagruz
          }
        })
      })
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
  <section class="stack content-shell">
    <AppCard v-if="loading" title="Загрузка">
      <p>Подтягиваю данные по юридическому лицу и программам.</p>
    </AppCard>

    <template v-else-if="details">
      <AppCard title="Слушатели юридического лица">
        <div class="listener-actions">
          <div class="listener-actions__left">
            <AppCheckbox v-model="selectAll" label="Выбрать всех" />
            <span class="listener-count">Выбрано: {{ selectedListenerIds.length }}</span>
          </div>
        </div>

        <div class="listener-list">
          <button
            v-for="listener in details.legal_entity.listeners || []"
            :key="listener.id_listener"
            type="button"
            class="listener-card"
            :class="{ 'listener-card--selected': selectedListenerIds.includes(listener.id_listener) }"
            @click="
              selectedListenerIds = selectedListenerIds.includes(listener.id_listener)
                ? selectedListenerIds.filter((id) => id !== listener.id_listener)
                : [...selectedListenerIds, listener.id_listener]
            "
          >
            <AppCheckbox
              :model-value="selectedListenerIds.includes(listener.id_listener)"
              :label="`${listener.second_name} ${listener.first_name} ${listener.middle_name || ''}`"
              disabled
            />
            <span class="listener-card__meta">{{ listener.snils }}</span>
          </button>
        </div>
      </AppCard>

      <AppCard title="Параметры договора">
        <div class="enrollment-grid">
          <AppSelect
            v-model="selectedContractId"
            label="Договор"
            placeholder="Выберите договор"
          >
            <option
              v-for="item in legalEntityEnrollmentContracts"
              :key="item.id_contract"
              :value="item.id_contract"
            >
              {{ item.name }}
            </option>
          </AppSelect>

          <AppSelect
            v-model="paymentOption"
            label="Порядок оплаты"
            placeholder="Выберите порядок оплаты"
          >
            <option value="full">100% предоплата</option>
            <option value="split">50/50</option>
            <option value="halfsplit">Оплата после акта</option>
          </AppSelect>

          <AppInput
            v-if="paymentOption === 'split'"
            v-model="secondPaymentDate"
            label="Срок второй оплаты"
            type="date"
          />

          <AppSelect
            v-model="ageCategory"
            label="Возрастная категория"
            placeholder="Выберите категорию"
          >
            <option v-for="(label, key) in ageCategories" :key="key" :value="key">{{ label }}</option>
          </AppSelect>

          <AppSelect
            v-model="optDocumentSelected"
            label="Итоговый документ"
            placeholder="Выберите режим выдачи"
          >
            <option v-for="(label, key) in optDocumentOptions" :key="key" :value="String(key)">
              {{ label }}
            </option>
          </AppSelect>

          <AppSelect
            v-model="studyLoadOption"
            label="Недельная нагрузка"
            placeholder="Выберите вариант"
          >
            <option v-for="(label, key) in loadVariantsNotDO" :key="key" :value="String(key)">
              {{ label }}
            </option>
          </AppSelect>
        </div>
      </AppCard>

      <AppCard title="Исполнитель">
        <div class="executer-grid">
          <button
            v-for="executer in executers"
            :key="executer.id_executor"
            type="button"
            class="executer-card"
            :class="{ 'executer-card--selected': selectedExecutorId === executer.id_executor }"
            @click="selectedExecutorId = executer.id_executor"
          >
            <span class="executer-card__name">
              {{ executer.second_name }} {{ executer.first_name }} {{ executer.middle_name || '' }}
            </span>
          </button>
        </div>
      </AppCard>

      <AppCard title="Программа обучения">
        <div class="enrollment-grid">
          <div class="enrollment-grid__full">
            <AppSelect
              v-model="selectedProgramId"
              label="Программа"
              placeholder="Выберите программу"
            >
              <option
                v-for="program in programs"
                :key="program.id_program_education"
                :value="program.id_program_education"
              >
                {{ program.name_prof_education }}
              </option>
            </AppSelect>
          </div>

          <AppSelect
            v-model="selectedGroupId"
            label="Группа"
            placeholder="Выберите группу"
          >
            <option v-for="groupItem in groups" :key="groupItem.group" :value="groupItem.group">
              {{ groupItem.name_group }}
            </option>
          </AppSelect>

          <AppInput
            :model-value="selectedProgram ? `${selectedProgram.price} ₽` : ''"
            label="Цена"
            disabled
          />

          <AppInput v-model="startDate" label="Дата начала" type="date" />
          <AppInput v-model="endDate" label="Дата окончания" type="date" />
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

.listener-actions__left {
  display: flex;
  align-items: center;
  gap: 1rem;
  flex-wrap: wrap;
}

.listener-count {
  color: #64748b;
}

.listener-list,
.executer-grid,
.enrollment-grid {
  display: grid;
  gap: 1rem;
}

.listener-list,
.executer-grid {
  grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
}

.listener-card,
.executer-card {
  display: grid;
  gap: 0.6rem;
  padding: 1rem;
  border-radius: 1rem;
  border: 1px solid rgba(15, 23, 42, 0.08);
  background: rgba(255, 255, 255, 0.78);
  text-align: left;
  transition:
    border-color 180ms ease,
    box-shadow 180ms ease,
    transform 180ms ease;
}

.listener-card:hover,
.executer-card:hover {
  transform: translateY(-1px);
  border-color: rgba(37, 99, 235, 0.22);
  box-shadow: 0 14px 34px rgba(15, 23, 42, 0.08);
}

.listener-card--selected,
.executer-card--selected {
  border-color: rgba(37, 99, 235, 0.38);
  box-shadow: 0 0 0 4px rgba(59, 130, 246, 0.1);
}

.listener-card__meta {
  font-size: 0.86rem;
  color: #64748b;
}

.executer-card__name {
  font-weight: 600;
  color: #0f172a;
}

.enrollment-grid {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.enrollment-grid__full {
  grid-column: 1 / -1;
}

@media (max-width: 900px) {
  .enrollment-grid {
    grid-template-columns: 1fr;
  }

  .enrollment-grid__full {
    grid-column: auto;
  }

  .page-actions {
    flex-direction: column;
  }
}
</style>
