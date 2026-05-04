<script setup lang="ts">
import type { DivisionItem, EducationTypeItem } from '../../../types/catalogs'
import type { GroupItem } from '../../../types/group'
import { getDateSortValue, normalizeApiDate } from '../../../utils/date'
import {
  ageCategories,
  legalEntityEnrollmentContracts,
  loadVariantsNotDO,
  optDocumentOptions
} from '../../../types/enrollment'
import GroupSearchField from '../../../components/features/enrollments/GroupSearchField.vue'
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
const divisions = ref<DivisionItem[]>([])
const educationTypes = ref<EducationTypeItem[]>([])
const groupOptions = ref<GroupItem[]>([])
const loading = ref(true)
const saving = ref(false)
const groupsLoading = ref(false)

const selectedListenerIds = ref<string[]>([])
const selectAll = ref(false)
const selectedDivisionId = ref('')
const selectedEducationTypeId = ref('')
const selectedContractId = ref('')
const selectedExecutorId = ref('')
const selectedProgramId = ref('')
const selectedGroupId = ref('')
const groupSearch = ref('')
const startDate = ref('')
const endDate = ref('')
const typeOfRetraining = ref('')
const paymentOption = ref('')
const secondPaymentDate = ref('')
const ageCategory = ref('')
const studyLoadOption = ref('')
const optDocumentSelected = ref('')

let groupSearchTimer: ReturnType<typeof setTimeout> | null = null

const filteredPrograms = computed(() =>
  programs.value.filter((item) =>
    item.id_divisions_education === selectedDivisionId.value &&
    item.id_education_type === selectedEducationTypeId.value
  )
)

const selectedProgram = computed(() =>
  filteredPrograms.value.find((item) => item.id_program_education === selectedProgramId.value) || null
)

const selectedGroup = computed(() =>
  groupOptions.value.find((item) => item.group === selectedGroupId.value) || null
)

const selectedGroupSchedule = computed(() =>
  [...(selectedGroup.value?.rapspisanie || [])]
    .map((item) => ({
      ...item,
      date: normalizeApiDate(item.date)
    }))
    .filter((item) => item.date)
    .sort((left, right) => getDateSortValue(left.date) - getDateSortValue(right.date))
)

const toggleListener = (listenerId: string) => {
  selectedListenerIds.value = selectedListenerIds.value.includes(listenerId)
    ? selectedListenerIds.value.filter((id) => id !== listenerId)
    : [...selectedListenerIds.value, listenerId]
}

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

watch(filteredPrograms, (items) => {
  if (selectedProgramId.value && !items.some((item) => item.id_program_education === selectedProgramId.value)) {
    selectedProgramId.value = ''
  }
}, { immediate: true })

watch(selectedProgramId, () => {
  selectedGroupId.value = ''
  groupSearch.value = ''
  startDate.value = ''
  endDate.value = ''
})

const syncSelectedGroup = (groups: GroupItem[]) => {
  if (!selectedGroupId.value || groups.some((item) => item.group === selectedGroupId.value)) {
    return groups
  }

  if (!groupSearch.value.trim()) {
    selectedGroupId.value = ''
    return groups
  }

  return [
    {
      group: selectedGroupId.value,
      name_group: groupSearch.value,
      rapspisanie: []
    },
    ...groups
  ]
}

const loadGroups = async (filter: string) => {
  groupsLoading.value = true

  try {
    const groups = await getGroups(1, filter.trim())
    groupOptions.value = syncSelectedGroup(groups)
  } catch (error) {
    notifications.error(
      error instanceof Error ? error.message : 'Не удалось загрузить список групп',
      'Запись на курс'
    )
  } finally {
    groupsLoading.value = false
  }
}

watch(groupSearch, (value) => {
  if (groupSearchTimer) {
    clearTimeout(groupSearchTimer)
  }

  if (selectedGroupId.value && value.trim() !== (selectedGroup.value?.name_group || '')) {
    selectedGroupId.value = ''
  }

  groupSearchTimer = setTimeout(() => {
    void loadGroups(value)
  }, 300)
})

watch(selectedGroupSchedule, (schedule) => {
  startDate.value = schedule[0]?.date || ''
  endDate.value = schedule[schedule.length - 1]?.date || ''
}, { immediate: true })

const isFormValid = computed(() =>
  Boolean(
    selectedListenerIds.value.length &&
    selectedDivisionId.value &&
    selectedEducationTypeId.value &&
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
    const [detailsData, programData, executerData, divisionData, educationTypeData] = await Promise.all([
      getLegalEntityDetails(legalEntityId.value),
      getPrograms(1, ''),
      getExecuters(1, ''),
      getDivisions(''),
      getEducationTypes('')
    ])

    details.value = detailsData
    programs.value = programData
    executers.value = executerData
    divisions.value = divisionData
    educationTypes.value = educationTypeData
    selectedExecutorId.value = executerData[0]?.id_executor || ''
    await loadGroups('')
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

const buildRegAddressPayload = (address: Record<string, unknown> | null | undefined) => ({
  mail_index: String(address?.mail_index ?? ''),
  region: String(address?.region ?? ''),
  city: String(address?.city ?? ''),
  street: String(address?.street ?? ''),
  house: String(address?.house ?? ''),
  building: String(address?.building ?? ''),
  apartment: String(address?.apartment ?? '')
})

const selectGroup = (group: GroupItem) => {
  selectedGroupId.value = group.group
  groupSearch.value = group.name_group
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
    const selectedListeners =
      legalEntity.listeners?.filter((item) => selectedListenerIds.value.includes(item.id_listener)) || []

    const contractFrontData = {
      legal_entity: {
        listeners: selectedListeners,
        reg_address: buildRegAddressPayload(details.value?.reg_address as Record<string, unknown> | null | undefined),
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

    const listenerFrontData = {
      dogovor_type: selectedContractId.value,
      dogovor_age: ageCategory.value || 'EIGHTEEN',
      opt_document: optDocumentSelected.value ? Number(optDocumentSelected.value) : null,
      opt_price: getPaymentText(),
      opt_nagruz: optNagruz,
      variant: optNagruz
    }

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
      })
    )

    await createEnrollmentDocumentRequest({
      id_listener: selectedListenerIds.value[0],
      id_program: selectedProgramId.value,
      id_executor: selectedExecutorId.value || null,
      front_data: contractFrontData
    })

    await Promise.all(
      selectedListenerIds.value.map(async (listenerId) => {
        await createEnrollmentDocumentRequest({
          id_listener: listenerId,
          id_program: selectedProgramId.value,
          id_executor: selectedExecutorId.value || null,
          front_data: listenerFrontData
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

onBeforeUnmount(() => {
  if (groupSearchTimer) {
    clearTimeout(groupSearchTimer)
  }
})

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
          <article
            v-for="listener in details.legal_entity.listeners || []"
            :key="listener.id_listener"
            role="button"
            tabindex="0"
            class="listener-card"
            :class="{ 'listener-card--selected': selectedListenerIds.includes(listener.id_listener) }"
            @click="toggleListener(listener.id_listener)"
            @keydown.enter.prevent="toggleListener(listener.id_listener)"
            @keydown.space.prevent="toggleListener(listener.id_listener)"
          >
            <div class="listener-card__top">
              <span class="listener-card__checkbox" :class="{ 'listener-card__checkbox--selected': selectedListenerIds.includes(listener.id_listener) }">
                <svg viewBox="0 0 16 16" fill="none">
                  <path d="M3.5 8.2L6.6 11.3L12.5 4.9" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" />
                </svg>
              </span>
              <span class="listener-card__name">
                {{ listener.second_name }} {{ listener.first_name }} {{ listener.middle_name || '' }}
              </span>
            </div>
            <span class="listener-card__meta">СНИЛС: {{ listener.snils || '—' }}</span>
          </article>
        </div>
      </AppCard>

      <AppCard title="Параметры договора">
        <div class="enrollment-grid">
          <AppSelect v-model="selectedContractId" label="Договор" placeholder="Выберите договор">
            <option v-for="item in legalEntityEnrollmentContracts" :key="item.id_contract" :value="item.id_contract">
              {{ item.name }}
            </option>
          </AppSelect>

          <AppSelect v-model="paymentOption" label="Порядок оплаты" placeholder="Выберите порядок оплаты">
             <option value="full">Оплата осуществляется в следующем порядке: 100% предоплата до начала обучения.</option>
            <option value="split">Оплата осуществляется в следующем порядке: аванс 50%, оставшиеся 50% — в установленный срок</option>
            <option value="halfsplit">Оплата подлежит перечислению на расчётный счёт Исполнителя в срок до 5 (пяти) рабочих дней, считая с момента (даты) подписания Сторонами Акта оказанных услуг.</option>
          </AppSelect>

          <AppInput
            v-if="paymentOption === 'split'"
            v-model="secondPaymentDate"
            label="Срок второй оплаты"
            type="date"
          />

          <AppSelect v-model="ageCategory" label="Возрастная категория" placeholder="Выберите категорию">
            <option v-for="(label, key) in ageCategories" :key="key" :value="key">{{ label }}</option>
          </AppSelect>

          <AppSelect v-model="optDocumentSelected" label="Итоговый документ" placeholder="Выберите режим выдачи">
            <option v-for="(label, key) in optDocumentOptions" :key="key" :value="String(key)">
              {{ label }}
            </option>
          </AppSelect>

          <AppSelect v-model="studyLoadOption" label="Недельная нагрузка" placeholder="Выберите вариант">
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
          <AppSelect v-model="selectedDivisionId" label="Подразделение" placeholder="Сначала выберите подразделение">
            <option v-for="division in divisions" :key="division.id_divisionsEducation" :value="division.id_divisionsEducation">
              {{ division.divisions }}
            </option>
          </AppSelect>

          <AppSelect
            v-model="selectedEducationTypeId"
            label="Тип обучения"
            placeholder="Выберите тип обучения"
            :disabled="!selectedDivisionId"
          >
            <option v-for="type in educationTypes" :key="type.id_educationType" :value="type.id_educationType">
              {{ type.typeName }}
            </option>
          </AppSelect>

          <AppSelect
            class="enrollment-grid__wide"
            v-model="selectedProgramId"
            label="Программа"
            placeholder="Выберите программу"
            :disabled="!selectedDivisionId || !selectedEducationTypeId || !filteredPrograms.length"
            :help="selectedDivisionId && selectedEducationTypeId && !filteredPrograms.length ? 'По выбранным фильтрам программ не найдено.' : ''"
          >
            <option v-for="program in filteredPrograms" :key="program.id_program_education" :value="program.id_program_education">
              {{ program.name_prof_education }}
            </option>
          </AppSelect>

          <GroupSearchField
            class="enrollment-grid__wide"
            v-model="groupSearch"
            label="Группа"
            placeholder="Введите часть названия группы"
            :options="groupOptions"
            :loading="groupsLoading"
            :disabled="!selectedProgramId"
            :help="selectedGroupId ? `Выбрана группа: ${selectedGroup?.name_group || ''}. Даты подставлены автоматически.` : 'Поиск показывает похожие группы по введённому тексту.'"
            @select="selectGroup"
          />

          <AppInput :model-value="selectedProgram ? `${selectedProgram.price} ₽` : ''" label="Цена" disabled />
          <AppInput :model-value="startDate" label="Дата начала" type="date" disabled />
          <AppInput :model-value="endDate" label="Дата окончания" type="date" disabled />
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

.listener-card__top {
  display: flex;
  align-items: flex-start;
  gap: 0.8rem;
}

.listener-card__checkbox {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 1.35rem;
  height: 1.35rem;
  flex: 0 0 1.35rem;
  border-radius: 0.45rem;
  border: 1px solid rgba(15, 23, 42, 0.14);
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.98) 0%, rgba(241, 245, 249, 0.92) 100%);
  box-shadow: 0 10px 22px rgba(15, 23, 42, 0.08);
  color: transparent;
}

.listener-card__checkbox svg {
  width: 0.9rem;
  height: 0.9rem;
}

.listener-card__checkbox--selected {
  border-color: rgba(37, 99, 235, 0.4);
  background: linear-gradient(135deg, #0f172a 0%, #2563eb 100%);
  color: #eff6ff;
}

.listener-card__name {
  font-weight: 600;
  color: #0f172a;
}

.executer-card__name {
  font-weight: 600;
  color: #0f172a;
}

.enrollment-grid {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.enrollment-grid__wide {
  grid-column: 1 / -1;
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
